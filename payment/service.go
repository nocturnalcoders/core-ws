package payment

import (
	"backendEkost/user"
	"strconv"

	//tidak boleh pake - minus dalam import package
	//harus pake alias midtrans
	midtrans "github.com/veritrans/go-midtrans"
)

//Integrasi Payment dengan Midtrans
type service struct {
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
	// ProcessPayment(input transaction.TransactioNotificationInput) error
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-sjkyPOba5goq8cwRKjTqHaGx"
	midclient.ClientKey = "SB-Mid-client-FryC8709HMXHtV5h"
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil

}
