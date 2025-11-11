package main

import (
	"github.com/faizauthar12/doku/app/config"
)

func init() {
	config.InitConfig()
}

func main() {
	//cfg := config.Get()
	//dokuUseCase := usecases.NewDokuUseCase(cfg.Doku.ClientID, cfg.Doku.SecretKey)

	// // Create Sub Account
	//dokuCreateSubAccountRequest := &requests.DokuCreateSubAccountRequest{
	//	Name:  "Faiz",
	//	Email: "faizauthar+subaccount1@gmail.com",
	//}
	//
	//resultCreateAccount, logData := dokuUseCase.CreateAccount(dokuCreateSubAccountRequest)
	//if logData != nil {
	//	panic(logData)
	//}
	//
	//fmt.Printf("Result Create Account: %+v\n", resultCreateAccount)

	//// Create Payment
	//dokuCreatePaymentRequest := &requests.DokuCreatePaymentRequest{
	//	Amount:         100000,
	//	CustomerName:   "Faiz Authar",
	//	CustomerEmail:  "faizauthar+customer1@gmail.com",
	//	SacID:          "SAC-8760-1762081713175",
	//	PaymentDueDate: 60,
	//}
	//
	//resultCreatePayment, logData := dokuUseCase.AcceptPayment(dokuCreatePaymentRequest)
	//if logData != nil {
	//	panic(logData)
	//}
	//
	//resultCreatePaymentJson, _ := json.Marshal(resultCreatePayment)
	//
	//fmt.Printf("Result Create Payment: %s\n", string(resultCreatePaymentJson))
	//
	//// Get Balance
	//dokuSacID := "SAC-8760-1762081713175"
	//
	//resultGetBalance, logData := dokuUseCase.GetBalance(dokuSacID)
	//if logData != nil {
	//	panic(logData)
	//}
	//
	//resultGetBalanceJson, _ := json.Marshal(resultGetBalance)
	//
	//fmt.Printf("Result Get Balance: %s\n", string(resultGetBalanceJson))
}
