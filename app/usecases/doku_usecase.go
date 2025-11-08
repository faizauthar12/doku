package usecases

import (
	"crypto/hmac"
	"crypto/sha256"
	"doku/app/config"
	"doku/app/models"
	"doku/app/requests"
	"doku/app/responses"
	"doku/app/utils/helper"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type DokuUseCaseInterface interface {
	CreateAccount(request *requests.DokuCreateSubAccountRequest) (*responses.DokuCreateSubAccountAccountResponse, *models.ErrorLog)
	AcceptPayment(request *requests.DokuCreatePaymentRequest) (*responses.DokuCreatePaymentHTTPResponse, *models.ErrorLog)
	GetBalance(sacID string) (*responses.DokuGetBalanceHTTPResponse, *models.ErrorLog)
}

type dokuUseCase struct {
	DokuAPIClientID  string
	DokuAPISecretKey string
}

func NewDokuUseCase(
	dokuAPIClientID string,
	dokuAPISecretKey string,
) DokuUseCaseInterface {
	return &dokuUseCase{
		DokuAPIClientID:  dokuAPIClientID,
		DokuAPISecretKey: dokuAPISecretKey,
	}
}

func (u *dokuUseCase) createSignatureComponent(
	requestId string,
	requestTimestamp *time.Time,
	requestTarget string,
	jsonBody []byte,
) (string, *models.ErrorLog) {

	// Format timestamp
	timestamp := requestTimestamp.Format("2006-01-02T15:04:05Z")

	// Calculate Digest - only for POST/PUT methods with body
	// For GET requests, jsonBody will be nil and digest will be empty
	var digest string
	if jsonBody != nil && len(jsonBody) > 0 {
		hash := sha256.Sum256(jsonBody)
		digest = base64.StdEncoding.EncodeToString(hash[:])
	}

	if u.DokuAPIClientID == "" {
		errorMessage := fmt.Sprintf("DokuClientId is empty")
		logData := helper.WriteLog(fmt.Errorf(errorMessage), http.StatusInternalServerError, errorMessage)
		return "", logData
	}

	// Build signature components string with \n separator
	// For GET requests (no body), Digest will be empty
	var signatureComponents string
	if digest != "" {
		signatureComponents = fmt.Sprintf(
			"Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s\nDigest:%s",
			u.DokuAPIClientID,
			requestId,
			timestamp,
			requestTarget,
			digest,
		)
	} else {
		signatureComponents = fmt.Sprintf(
			"Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s",
			u.DokuAPIClientID,
			requestId,
			timestamp,
			requestTarget,
		)
	}

	// Get secret key from environment
	secretKey := u.DokuAPISecretKey
	if secretKey == "" {
		logData := helper.WriteLog(fmt.Errorf("DOKU_SECRET_KEY not configured"), http.StatusInternalServerError, "Secret key is required for signature generation")
		return "", logData
	}

	// Calculate HMAC-SHA256 base64
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(signatureComponents))
	signatureHash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Format final signature with HMACSHA256= prefix
	signature := fmt.Sprintf("HMACSHA256=%s", signatureHash)

	return signature, nil
}

func (u *dokuUseCase) CreateAccount(request *requests.DokuCreateSubAccountRequest) (*responses.DokuCreateSubAccountAccountResponse, *models.ErrorLog) {

	createAccountPayload := &requests.DokuCreateSubAccountHTTPRequest{
		Account: requests.DokuCreateSubAccountAccount{
			Email: request.Email,
			Type:  "STANDARD",
			Name:  request.Name,
		},
	}

	createAccountPayloadJson, err := json.Marshal(createAccountPayload)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to marshal create account payload")
		return nil, logData
	}

	// preparing signature components
	requestId := uuid.NewString()
	requestTimeStamp := time.Now().UTC()
	requestTarget := "/sac-merchant/v1/accounts"

	signature, logData := u.createSignatureComponent(requestId, &requestTimeStamp, requestTarget, createAccountPayloadJson)
	if logData != nil {
		return nil, logData
	}

	// preparing request headers
	requestHeader := map[string]string{
		"Client-Id":         config.Get().Doku.ClientID,
		"Request-Id":        requestId,
		"Request-Timestamp": requestTimeStamp.Format("2006-01-02T15:04:05Z"),
		"Signature":         signature,
	}

	createAccountAPI := helper.POST(&helper.Options{
		Method:      "POST",
		URL:         "https://api-sandbox.doku.com/sac-merchant/v1/accounts",
		Body:        createAccountPayloadJson,
		Headers:     requestHeader,
		Timeout:     30 * time.Second,
		ContentType: "application/json",
		QueryParams: nil,
		IsPrintCurl: true,
	})

	if createAccountAPI.Error != nil {
		logData := helper.WriteLog(createAccountAPI.Error, createAccountAPI.StatusCode, helper.DefaultStatusText[createAccountAPI.StatusCode])
		return nil, logData
	}

	var createAccountResponse *responses.DokuCreateSubAccountHTTPResponse
	err = json.Unmarshal(createAccountAPI.Body, &createAccountResponse)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to unmarshal create account response")
		return nil, logData
	}

	return createAccountResponse.Account, nil
}

func (u *dokuUseCase) AcceptPayment(request *requests.DokuCreatePaymentRequest) (*responses.DokuCreatePaymentHTTPResponse, *models.ErrorLog) {

	invoiceNumber := fmt.Sprintf("INV-%d", time.Now().UnixNano())

	createPaymentPayload := &requests.DokuCreatePaymentHTTPRequest{
		Order: &models.DokuOrder{
			InvoiceNumber: invoiceNumber,
			Amount:        request.Amount,
		},
		Payment: &models.DokuPayment{
			PaymentDueDate: 60,
		},
		Customer: &models.DokuCustomer{
			Name:  request.CustomerName,
			Email: request.CustomerEmail,
		},
		AdditionalInfo: &models.DokuAdditionalInfo{
			Account: models.DokuAccount{
				ID: request.SacID,
			},
		},
	}

	createPaymentPayloadJson, err := json.Marshal(createPaymentPayload)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to marshal create payment payload")
		return nil, logData
	}

	// preparing signature components
	requestId := uuid.NewString()
	requestTimeStamp := time.Now().UTC()
	requestTarget := "/checkout/v1/payment"

	signature, logData := u.createSignatureComponent(requestId, &requestTimeStamp, requestTarget, createPaymentPayloadJson)
	if logData != nil {
		return nil, logData
	}

	// preparing request headers
	requestHeader := map[string]string{
		"Client-Id":         config.Get().Doku.ClientID,
		"Request-Id":        requestId,
		"Request-Timestamp": requestTimeStamp.Format("2006-01-02T15:04:05Z"),
		"Signature":         signature,
	}

	createPaymentAPI := helper.POST(&helper.Options{
		Method:      "POST",
		URL:         "https://api-sandbox.doku.com/checkout/v1/payment",
		Body:        createPaymentPayloadJson,
		Headers:     requestHeader,
		Timeout:     30 * time.Second,
		ContentType: "application/json",
		QueryParams: nil,
		IsPrintCurl: true,
	})

	if createPaymentAPI.Error != nil {
		logData := helper.WriteLog(createPaymentAPI.Error, createPaymentAPI.StatusCode, helper.DefaultStatusText[createPaymentAPI.StatusCode])
		return nil, logData
	}

	var createPaymentResponse *responses.DokuCreatePaymentHTTPResponse
	err = json.Unmarshal(createPaymentAPI.Body, &createPaymentResponse)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to unmarshal create payment response")
		return nil, logData
	}

	return createPaymentResponse, nil
}

func (u *dokuUseCase) GetBalance(sacID string) (*responses.DokuGetBalanceHTTPResponse, *models.ErrorLog) {

	// preparing signature components
	requestId := uuid.NewString()
	requestTimeStamp := time.Now().UTC()
	requestTarget := fmt.Sprintf("/sac-merchant/v1/balances/%s", sacID)

	signature, logData := u.createSignatureComponent(requestId, &requestTimeStamp, requestTarget, nil)
	if logData != nil {
		return nil, logData
	}

	// preparing request headers
	requestHeader := map[string]string{
		"Client-Id":         config.Get().Doku.ClientID,
		"Request-Id":        requestId,
		"Request-Timestamp": requestTimeStamp.Format("2006-01-02T15:04:05Z"),
		"Signature":         signature,
	}

	getBalanceAPI := helper.GET(&helper.Options{
		Method:      "GET",
		URL:         fmt.Sprintf("https://api-sandbox.doku.com/sac-merchant/v1/balances/%s", sacID),
		Headers:     requestHeader,
		Timeout:     30 * time.Second,
		ContentType: "application/json",
		QueryParams: nil,
		IsPrintCurl: true,
	})

	if getBalanceAPI.Error != nil {
		fmt.Printf("Get Balance API Error: %v\n", getBalanceAPI.Error)
		logData := helper.WriteLog(getBalanceAPI.Error, getBalanceAPI.StatusCode, helper.DefaultStatusText[getBalanceAPI.StatusCode])
		return nil, logData
	}

	var getBalanceResponse *responses.DokuGetBalanceHTTPResponse
	err := json.Unmarshal(getBalanceAPI.Body, &getBalanceResponse)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to unmarshal get balance response")
		return nil, logData
	}

	//fmt.Printf("Get Balance Response: %+v\n", getBalanceResponse)

	return getBalanceResponse, nil
}
