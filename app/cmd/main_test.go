package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/faizauthar12/doku/app/config"
	"github.com/faizauthar12/doku/app/requests"
	"github.com/faizauthar12/doku/app/usecases"
)

// func init() {
// 	config.InitConfig()
//
// 	cfg := config.Get()
// 	dokuUseCase = usecases.NewDokuUseCase(cfg.Doku.ClientID, cfg.Doku.SecretKey, cfg.Doku.PrivateKey)
// }

func LoadEnv() {
	re := regexp.MustCompile(`^(.*` + "doku" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	fmt.Println("Root Path:", string(rootPath))

	config.InitConfig(string(rootPath) + `/.env`)
}

func TestCreatePayment(t *testing.T) {
	LoadEnv()
	dokuUseCase := usecases.NewDokuUseCase(config.Get().Doku.ClientID, config.Get().Doku.SecretKey, config.Get().Doku.PrivateKey)

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
	LoadEnv()
	dokuUseCase := usecases.NewDokuUseCase(config.Get().Doku.ClientID, config.Get().Doku.SecretKey, config.Get().Doku.PrivateKey)

	sacID := "SAC-8760-1762081713175"

	resultGetBalance, logData := dokuUseCase.GetBalance(sacID)
	if logData != nil {
		t.Fatalf("Error getting balance: %+v", logData)
	}

	resultGetBalanceJson, _ := json.Marshal(resultGetBalance)

	t.Logf("Result Get Balance: %s\n", resultGetBalanceJson)
}

func TestGetToken(t *testing.T) {
	LoadEnv()
	dokuService := usecases.NewDokuUseCase(config.Get().Doku.ClientID, config.Get().Doku.SecretKey, config.Get().Doku.PrivateKey)

	resultGetToken, logData := dokuService.GetToken()
	if logData != nil {
		t.Fatalf("Error getting token: %+v", logData)
	}

	t.Logf("Result Get Token Struct: %+v\n", resultGetToken)

	resultGetTokenJson, err := json.Marshal(resultGetToken)
	if err != nil {
		t.Fatalf("Error marshaling result to JSON: %v", err)
	}

	t.Logf("Result Get Token: %s\n", resultGetTokenJson)
}

func TestBankAccountInquiry(t *testing.T) {
	LoadEnv()
	dokuUseCase := usecases.NewDokuUseCase(config.Get().Doku.ClientID, config.Get().Doku.SecretKey, config.Get().Doku.PrivateKey)

	accessToken, logData := dokuUseCase.GetToken()
	if logData != nil {
		t.Fatalf("Error getting access token: %+v", logData)
	}
	t.Logf("Access Token: %+v\n", accessToken)

	requestBody := `{"partnerReferenceNo":"hsjkans284b2he54","customerNumber":"628115678890","amount":{"value":"200000.00","currency":"IDR"},"beneficiaryAccountNumber":"8377388292","additionalInfo":{"beneficiaryBankCode":"014","beneficiaryAccountName":"FHILEA HERMANUS","senderCountryCode":"ID"}}`

	fmt.Println("Request Body:", requestBody)
	var dokuBankAccountInquiryRequest requests.DokuBankAccountInquiryRequest
	err := json.Unmarshal([]byte(requestBody), &dokuBankAccountInquiryRequest)
	if err != nil {
		t.Fatalf("Error unmarshaling request body: %v", err)
	}

	resultBankAccountInquiry, logData := dokuUseCase.BankAccountInquiry(&dokuBankAccountInquiryRequest, accessToken.AccessToken)
	if logData != nil {
		t.Fatalf("Error in bank account inquiry: %+v", logData)
	}

	t.Logf("Result Bank Account Inquiry: %+v\n", resultBankAccountInquiry)

}
