package handler

import (
	"backendEkost/helper"
	"backendEkost/transaction"
	"backendEkost/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetKostTransactions(c *gin.Context) {
	var input transaction.GetKostTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to Get Kost's Transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	transactions, err := h.service.GetTransactionsByKostID(input)
	if err != nil {
		response := helper.APIResponse("Failed to Get Kost's Transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Kost's Transactions", http.StatusOK, "success", transaction.FormatKostTransactions(transactions))
	c.JSON(http.StatusOK, response)

}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions, err := h.service.GetTransactionsByUserID(userID)
	if err != nil {
		response := helper.APIResponse("Failed to Get User's Transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("User's Transactions", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	//Mapping
	var input transaction.CreateTransactionInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed To Create Transaction", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)

	if err != nil {
		response := helper.APIResponse("Failed To Create Transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success To Create Transaction", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetNotification(c *gin.Context) {
	var input transaction.TransactioNotificationInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := helper.APIResponse("Failed to Process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	err = h.service.ProcessPayment(input)
	if err != nil {
		response := helper.APIResponse("Failed to Process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)

		return
	}
	c.JSON(http.StatusOK, input)
}

// parameter di uri
// ketika di handler, tangkap parameter mapping	input struct
// panggil service, input struct sebagai parameter
// service, berbekal kost id bisa panggil repo
// repo mencari data transaction suatu kost

//GetUserTransaction
//1. perlu tau siapa user yang login -> handler user (ambil nilai user dari jwt atau middleware)
//2. Service
//3. repo => ambil data transaction (preload -> data kost)

//Integrasi dengan midtrans
//1. Ada Inputan dari user
//2. Handler menangkap input , terus di-mapping ke input struct
//3. Memanggil Service untuk transaksi , perlu memanggil sistem midtrans-> snap gateway
//4. Memanggil repository create New transaction Data
