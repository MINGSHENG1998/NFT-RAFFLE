package controllers

import (
	"context"
	"net/http"
	"nft-raffle/logger"
	"nft-raffle/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ExpenseController IExpenseController = NewExpenseController()

	expenseCollection *mongo.Collection = nftRaffleDb.OpenCollection(nftRaffleDbClient, "expense")
)

type NewExpenseRequest struct {
	ExpenseType   string `validate:"required"`
	ExpenseLabel  string
	ExpenseTime   string `validate:"required"`
	ExpenseAmount int64  `validate:"required"`
}

type MonthlyExpenseRequest struct {
	FromDate string `validate:"required"`
	ToDate   string `validate:"required"`
}

type ExpenseUpdateDTO struct {
	ExpenseId     string `validate:"required"`
	ExpenseType   *string
	ExpenseLabel  *string
	ExpenseTime   *string
	ExpenseAmount *int64
}

type IExpenseController interface {
	CreateNewExpense(c *gin.Context)
	GetExpenses(c *gin.Context)
	GetExpensesByType(c *gin.Context)
	UpdateExpense(c *gin.Context)
}

type expenseControllerStruct struct{}

func NewExpenseController() IExpenseController {
	return &expenseControllerStruct{}
}

func (e expenseControllerStruct) CreateNewExpense(c *gin.Context) {
	userId := c.GetString("uid")

	if userId == "" {
		logger.Logger.Error("User ID is missing in the claim to create new expense")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID is missing in the claim to create new expense"})
		return
	}

	var request NewExpenseRequest

	err := c.BindJSON(&request)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validateErr := validate.Struct(request)
	if validateErr != nil {
		logger.Logger.Error(validateErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
		return
	}

	if request.ExpenseAmount < 0 {
		logger.Logger.Warn("expense amount cannnot be less than 0")
		c.JSON(http.StatusBadRequest, gin.H{"error": "expense amount cannnot be less than 0"})
		return
	}

	if request.ExpenseLabel == "" {
		request.ExpenseLabel = request.ExpenseType
	}

	var newExpense models.Expense

	newExpense.ID = primitive.NewObjectID()
	newExpense.Expense_id = newExpense.ID.Hex()
	newExpense.User_id = userId
	newExpense.Expense_label = request.ExpenseLabel
	newExpense.Expense_type = request.ExpenseType
	newExpense.Expense_amount = request.ExpenseAmount
	newExpense.Expense_time, err = timeHelper.ConvertDateTimeStringToCurrentLocationTime(request.ExpenseTime)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newExpense.Created_at, err = timeHelper.GetCurrentLocationTime()

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing created_at"})
		return
	}

	newExpense.Updated_at, err = timeHelper.GetCurrentLocationTime()

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing updated_at"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resultInsertionNumber, err := expenseCollection.InsertOne(ctx, newExpense)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resultInsertionNumber)
}

func (e expenseControllerStruct) GetExpenses(c *gin.Context) {
	userId := c.GetString("uid")

	if userId == "" {
		logger.Logger.Error("User ID is missing in the claim to create new expense")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID is missing in the claim to create new expense"})
		return
	}

	var request MonthlyExpenseRequest

	err := c.BindJSON(&request)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = validate.Struct(request)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fromDate, err := timeHelper.ConvertDateTimeStringToCurrentLocationTime(request.FromDate)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	toDate, err := timeHelper.ConvertDateTimeStringToCurrentLocationTime(request.ToDate)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var andQuery []bson.M
	andQuery = append(andQuery, bson.M{"user_id": userId})
	andQuery = append(andQuery, bson.M{"expense_time": bson.M{"$gte": fromDate}})
	andQuery = append(andQuery, bson.M{"expense_time": bson.M{"$lte": toDate}})

	matchStage := bson.D{
		{Key: "$match", Value: bson.M{
			"$and": andQuery,
		}},
	}

	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "expense_time", Value: -1},
		}},
	}

	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
			{Key: "totalCount", Value: bson.M{
				"$count": bson.M{},
			}},
			{Key: "totalExpenseAmount", Value: bson.M{
				"$sum": "$expense_amount",
			}},
			{Key: "data", Value: bson.M{
				"$push": "$$ROOT",
			}},
		}},
	}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "totalCount", Value: 1},
			{Key: "totalExpenseAmount", Value: 1},
			{Key: "data", Value: 1},
		}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := expenseCollection.Aggregate(ctx, mongo.Pipeline{matchStage, sortStage, groupStage, projectStage})

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var data []bson.M

	err = result.All(ctx, &data)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(data) > 0 {
		c.JSON(http.StatusOK, data[0])
		return
	}

	c.Status(http.StatusOK)
}

func (e expenseControllerStruct) GetExpensesByType(c *gin.Context) {
	userId := c.GetString("uid")

	if userId == "" {
		logger.Logger.Error("User ID is missing in the claim to create new expense")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID is missing in the claim to create new expense"})
		return
	}

	var request MonthlyExpenseRequest

	err := c.BindJSON(&request)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = validate.Struct(request)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fromDate, err := timeHelper.ConvertDateTimeStringToCurrentLocationTime(request.FromDate)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	toDate, err := timeHelper.ConvertDateTimeStringToCurrentLocationTime(request.ToDate)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var andQuery []bson.M
	andQuery = append(andQuery, bson.M{"user_id": userId})
	andQuery = append(andQuery, bson.M{"expense_time": bson.M{"$gte": fromDate}})
	andQuery = append(andQuery, bson.M{"expense_time": bson.M{"$lte": toDate}})

	matchStage := bson.D{
		{Key: "$match", Value: bson.M{
			"$and": andQuery,
		}},
	}

	groupStageByExpenseType := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$expense_type"},
			{Key: "expenseAmountByType", Value: bson.M{
				"$sum": "$expense_amount",
			}},
			{Key: "expenseCountByType", Value: bson.M{
				"$count": bson.M{},
			}},
			{Key: "expenseTypeData", Value: bson.M{
				"$push": "$$ROOT",
			}},
		}},
	}

	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "expenseAmountByType", Value: -1},
		}},
	}

	groupStage2 := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
			{Key: "totalExpenseTypes", Value: bson.M{
				"$count": bson.M{},
			}},
			{Key: "data", Value: bson.M{
				"$push": "$$ROOT",
			}},
		}},
	}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "totalExpenseTypes", Value: 1},
			{Key: "data", Value: 1},
		}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := expenseCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStageByExpenseType, sortStage, groupStage2, projectStage})

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var data []bson.M

	err = result.All(ctx, &data)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(data) > 0 {
		c.JSON(http.StatusOK, data[0])
		return
	}

	c.Status(http.StatusOK)
}

func (e expenseControllerStruct) UpdateExpense(c *gin.Context) {
	var updateDto ExpenseUpdateDTO

	err := c.BindJSON(&updateDto)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updateObj bson.D

	if updateDto.ExpenseType != nil {
		updateObj = append(updateObj, bson.E{Key: "expense_type", Value: *updateDto.ExpenseType})
	}

	if updateDto.ExpenseLabel != nil {
		updateObj = append(updateObj, bson.E{Key: "expense_label", Value: *updateDto.ExpenseLabel})
	}

	if updateDto.ExpenseAmount != nil {
		if *updateDto.ExpenseAmount < 0 {
			logger.Logger.Warn("expense amount cannnot be less than 0")
			c.JSON(http.StatusBadRequest, gin.H{"error": "expense amount cannnot be less than 0"})
			return
		}

		updateObj = append(updateObj, bson.E{Key: "expense_amount", Value: *updateDto.ExpenseAmount})
	}

	if updateDto.ExpenseTime != nil {
		newExpenseTime, err := timeHelper.ConvertDateTimeStringToCurrentLocationTime(*updateDto.ExpenseTime)

		if err != nil {
			logger.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		updateObj = append(updateObj, bson.E{Key: "expense_time", Value: newExpenseTime})
	}

	if len(updateObj) < 1 {
		logger.Logger.Warn("update dto cannot be empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "update dto cannot be empty"})
		return
	}

	updated_at, err := timeHelper.GetCurrentLocationTime()

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: updated_at})

	upsert := false
	filter := bson.M{"expense_id": updateDto.ExpenseId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err = expenseCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var expense models.Expense

	err = expenseCollection.FindOne(ctx, bson.M{"expense_id": updateDto.ExpenseId}).Decode(&expense)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expense)
}
