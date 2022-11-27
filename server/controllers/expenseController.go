package controllers

import (
	"context"
	"net/http"
	"nft-raffle/logger"
	"nft-raffle/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

type IExpenseController interface {
	CreateNewExpense(c *gin.Context)
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
	newExpense.Expense_time, err = timeHelper.ConvertDateTimeStringToSingaporeTime(request.ExpenseTime)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newExpense.Created_at, err = timeHelper.GetCurrentTimeSingapore()

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing created_at"})
		return
	}

	newExpense.Updated_at, err = timeHelper.GetCurrentTimeSingapore()

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
