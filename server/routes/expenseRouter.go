package routes

import (
	"nft-raffle/controllers"

	"github.com/gin-gonic/gin"
)

var (
	expenseController controllers.IExpenseController = controllers.ExpenseController
)

func ExpenseRoutes(superRoute *gin.RouterGroup) {
	expenseRouter := superRoute.Group("/expense")

	expenseRouter.POST("/create-new-expense", authMiddleware.Authenticate, expenseController.CreateNewExpense)
}
