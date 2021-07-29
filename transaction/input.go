package transaction

import "backendEkost/user"

type GetKostTransactionsInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}

//Integrasi dengan midtrans

type CreateTransactionInput struct {
	KostID int `json:"kost_id"`
	Amount int `json:"amount" binding:"required"`
	User   user.User
}

type TransactioNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
