package transaction

import "time"

type KostTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"creted_at"`
}

func FormatKostTransaction(transaction Transaction) KostTransactionFormatter {
	formatter := KostTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt

	return formatter
}

//Mengubah list of transaction menjadi slice / list kost transaction formatter
func FormatKostTransactions(transactions []Transaction) []KostTransactionFormatter {
	if len(transactions) == 0 {
		return []KostTransactionFormatter{}
	}
	var transactionsFormatter []KostTransactionFormatter

	for _, transaction := range transactions {
		formatter := FormatKostTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}

	return transactionsFormatter
}

type UserTransactionFormatter struct {
	ID        int           `json:"id"`
	Amount    int           `json:"amount"`
	Status    string        `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	Kost      KostFormatter `json:"kost"`
}

type KostFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {
	formatter := UserTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	kostFormatter := KostFormatter{}
	kostFormatter.Name = transaction.Kost.Name
	kostFormatter.ImageURL = ""

	if len(transaction.Kost.KostImages) > 0 {

		kostFormatter.ImageURL = transaction.Kost.KostImages[0].FileName
	}

	formatter.Kost = kostFormatter

	return formatter
}

func FormatUserTransactions(transactions []Transaction) []UserTransactionFormatter {
	if len(transactions) == 0 {
		return []UserTransactionFormatter{}
	}

	//Hasilnya Single Object
	//di Append lalu masuk ke single object
	var transactionsFormatter []UserTransactionFormatter

	for _, transaction := range transactions {
		formatter := FormatUserTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}

	return transactionsFormatter
}

type TransactionFormatter struct {
	ID         int    `json:"id"`
	KostID     int    `json:"kost_id"`
	UserID     int    `json:"user_id"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

func FormatTransaction(transaction Transaction) TransactionFormatter {
	formatter := TransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.KostID = transaction.KostID
	formatter.UserID = transaction.UserID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.Code = transaction.Code
	formatter.PaymentURL = transaction.PaymentURL
	return formatter
}
