package main

import (
	"encoding/json"
	"github.com/faizauthar12/doku/app/config"
	"github.com/faizauthar12/doku/app/requests"
	"github.com/faizauthar12/doku/app/usecases"
	"testing"
)

var dokuUseCase usecases.DokuUseCaseInterface

func init() {
	config.InitConfig()

	cfg := config.Get()
	dokuUseCase = usecases.NewDokuUseCase(cfg.Doku.ClientID, cfg.Doku.SecretKey)
}

func TestCreatePayment(t *testing.T) {
	dokuCreatePaymentRequest := &requests.DokuCreatePaymentRequest{
		Amount:        100000,
		CustomerName:  "Faiz Authar",
		CustomerEmail: "faiz+customer1@gmail.com",
		SacID:         "SAC-8760-1762081713175",
	}

	resultCreatePayment, logData := dokuUseCase.AcceptPayment(dokuCreatePaymentRequest)
	if logData != nil {
		t.Fatalf("Error creating payment: %+v", logData)
	}

	t.Logf("Result Create Payment: %+v\n", resultCreatePayment)
}

func TestGetBalance(t *testing.T) {
	sacID := "SAC-8760-1762081713175"

	resultGetBalance, logData := dokuUseCase.GetBalance(sacID)
	if logData != nil {
		t.Fatalf("Error getting balance: %+v", logData)
	}

	resultGetBalanceJson, _ := json.Marshal(resultGetBalance)

	t.Logf("Result Get Balance: %s\n", resultGetBalanceJson)
}
