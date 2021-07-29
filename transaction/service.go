package transaction

import (
	"backendEkost/kost"
	"backendEkost/payment"
	"errors"
	"strconv"
)

type service struct {
	repository     Repository
	kostRepository kost.Repository
	paymentService payment.Service
}

type Service interface {
	GetTransactionsByKostID(input GetKostTransactionsInput) ([]Transaction, error)
	GetTransactionsByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
	ProcessPayment(input TransactioNotificationInput) error
}

func NewService(repository Repository, kostRepository kost.Repository, paymentService payment.Service) *service {
	return &service{repository, kostRepository, paymentService}
}

func (s *service) GetTransactionsByKostID(input GetKostTransactionsInput) ([]Transaction, error) {
	//get Kost
	//Check Kost.UserID != user_id_yang_melakukan_request

	kost, err := s.kostRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if kost.UserID != input.User.ID {
		return []Transaction{}, errors.New("Not an owner of the Kost")
	}

	transactions, err := s.repository.GetByKostID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionsByUserID(userID int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.KostID = input.KostID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "Pending"
	transaction.Code = "E-Kost-Order-01"

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)

	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

func (s *service) ProcessPayment(input TransactioNotificationInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.repository.GetByID(transaction_id)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "Paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "Paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "Cancelled"
	}

	updatedTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return err
	}

	kost, err := s.kostRepository.FindByID(updatedTransaction.KostID)

	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		// kost. = kost.LiverCount + 1
		kost.CurrentSpaceCount = kost.CurrentSpaceCount - 1

		_, err := s.kostRepository.Update(kost)

		if err != nil {
			return err
		}

	}
	return nil
}
