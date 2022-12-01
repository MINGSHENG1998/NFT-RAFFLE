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
	expenseRouter.POST("/get-expenses", authMiddleware.Authenticate, expenseController.GetExpenses)
	expenseRouter.POST("/get-expenses-by-type", authMiddleware.Authenticate, expenseController.GetExpensesByType)
	expenseRouter.PATCH("/update-expense", authMiddleware.Authenticate, expenseController.UpdateExpense)
}
