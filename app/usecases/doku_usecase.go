package usecases

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/faizauthar12/doku/app/models"
	"github.com/faizauthar12/doku/app/requests"
	"github.com/faizauthar12/doku/app/responses"
	"github.com/faizauthar12/doku/app/utils/helper"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type DokuUseCaseInterface interface {
	CreateAccount(request *requests.DokuCreateSubAccountRequest) (*responses.DokuCreateSubAccountAccountResponse, *models.ErrorLog)
	AcceptPayment(request *requests.DokuCreatePaymentRequest) (*responses.DokuCreatePaymentHTTPResponse, *models.ErrorLog)
	GetBalance(sacID string) (*responses.DokuGetBalanceHTTPResponse, *models.ErrorLog)
	HandleNotification(request *requests.DokuNotificationRequest) (*responses.DokuPostNotificationHTTPResponse, *models.ErrorLog)
	GetToken() (*responses.GetTokenResponse, *models.ErrorLog)
	BankAccountInquiry(request *requests.DokuBankAccountInquiryRequest, accessToken string) (*responses.BankAccountInquiryResponse, *models.ErrorLog)
}

type dokuUseCase struct {
	DokuAPIClientID  string
	DokuAPISecretKey string
	DokuPrivateKey   string
}

func NewDokuUseCase(
	dokuAPIClientID string,
	dokuAPISecretKey string,
	dokuPrivateKey string,
) DokuUseCaseInterface {
	return &dokuUseCase{
		DokuAPIClientID:  dokuAPIClientID,
		DokuAPISecretKey: dokuAPISecretKey,
		DokuPrivateKey:   dokuPrivateKey,
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
		logData := helper.WriteLog(errors.New(errorMessage), http.StatusInternalServerError, errorMessage)
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

func (u *dokuUseCase) verifySignatureComponent(
	signature string,
	requestId string,
	requestTimestamp string,
	requestTarget string,
	jsonBody []byte,
) (bool, *models.ErrorLog) {

	// Format timestamp
	timestamp, err := time.Parse("2006-01-02T15:04:05Z", requestTimestamp)
	if err != nil {
		errorMessage := fmt.Sprintf("Invalid request timestamp format")
		logData := helper.WriteLog(err, http.StatusBadRequest, errorMessage)
		return false, logData
	}

	// Calculate Digest - only for POST/PUT methods with body
	// For GET requests, jsonBody will be nil and digest will be empty
	var digest string
	if jsonBody != nil && len(jsonBody) > 0 {
		hash := sha256.Sum256(jsonBody)
		digest = base64.StdEncoding.EncodeToString(hash[:])
	}

	signatureComponents := fmt.Sprintf(
		"Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s\nDigest:%s",
		u.DokuAPIClientID,
		requestId,
		timestamp,
		requestTarget,
		digest,
	)

	if signature != signatureComponents {
		errorMessage := fmt.Sprintf("Signature does not match")
		logData := helper.WriteLog(errors.New(errorMessage), http.StatusUnauthorized, errorMessage)
		return false, logData
	}

	return true, nil
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
		"Client-Id":         u.DokuAPIClientID,
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
	if createAccountAPI.StatusCode != http.StatusOK {
		dokuErrorResponse := &responses.DokuErrorHTTPResponse{}
		err = json.Unmarshal(createAccountAPI.Body, &dokuErrorResponse)
		if err != nil {
			logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to unmarshal create payment error response")
			return nil, logData
		}

		if dokuErrorResponse.Message != nil {
			if strings.Contains(dokuErrorResponse.Message[0], "email already registered") {
				// extract the account id from error message
				// message format: "email already registered with account id: SAC-6604-1764297586731"
				parts := strings.Split(dokuErrorResponse.Message[0], "account id: ")
				if len(parts) == 2 {
					sacID := strings.TrimSpace(parts[1])
					return &responses.DokuCreateSubAccountAccountResponse{
						ID: null.StringFrom(sacID),
					}, nil
				}
			}
		}

		errorMessage := fmt.Sprintf("Doku Create Sub Account API Error: %v", dokuErrorResponse.Message)
		logData := helper.WriteLog(fmt.Errorf("Doku Create Sub Account API Error: %v", dokuErrorResponse.Message), createAccountAPI.StatusCode, errorMessage)
		return nil, logData
	}

	err = json.Unmarshal(createAccountAPI.Body, &createAccountResponse)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to unmarshal create account response")
		return nil, logData
	}

	return createAccountResponse.Account, nil
}

func (u *dokuUseCase) AcceptPayment(request *requests.DokuCreatePaymentRequest) (*responses.DokuCreatePaymentHTTPResponse, *models.ErrorLog) {

	createPaymentPayload := &requests.DokuCreatePaymentHTTPRequest{
		Order: &models.DokuOrder{
			InvoiceNumber: null.StringFrom(request.InvoiceNumber),
			Amount:        null.IntFrom(request.Amount),
		},
		Payment: &models.DokuPayment{
			PaymentDueDate: null.IntFrom(request.PaymentDueDate),
		},
		Customer: &models.DokuCustomer{
			Name:  null.StringFrom(request.CustomerName),
			Email: null.StringFrom(request.CustomerEmail),
		},
		AdditionalInfo: &models.DokuAdditionalInfo{
			Account: models.DokuAccount{
				ID: null.StringFrom(request.SacID),
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
		"Client-Id":         u.DokuAPIClientID,
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

	if createPaymentAPI.StatusCode != http.StatusOK {
		dokuErrorResponse := &responses.DokuErrorHTTPResponse{}
		err = json.Unmarshal(createPaymentAPI.Body, &dokuErrorResponse)
		if err != nil {
			logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to unmarshal create payment error response")
			return nil, logData
		}

		errorMessage := fmt.Sprintf("Doku Create Payment API Error: %v", dokuErrorResponse.Message)
		logData := helper.WriteLog(fmt.Errorf("Doku Create Payment API Error: %v", dokuErrorResponse.Message), createPaymentAPI.StatusCode, errorMessage)
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
		"Client-Id":         u.DokuAPIClientID,
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

	if getBalanceAPI.StatusCode != http.StatusOK {
		dokuErrorResponse := &responses.DokuErrorHTTPResponse{}
		err := json.Unmarshal(getBalanceAPI.Body, &dokuErrorResponse)
		if err != nil {
			logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to unmarshal get balance error response")
			return nil, logData
		}

		errorMessage := fmt.Sprintf("Doku Get Balance API Error: %v", dokuErrorResponse.Message)
		logData := helper.WriteLog(fmt.Errorf("Doku Get Balance API Error: %v", dokuErrorResponse.Message), getBalanceAPI.StatusCode, errorMessage)
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

func (u *dokuUseCase) HandleNotification(request *requests.DokuNotificationRequest) (*responses.DokuPostNotificationHTTPResponse, *models.ErrorLog) {

	// verifying signature components
	isValid, logData := u.verifySignatureComponent(
		request.Signature,
		request.RequestID,
		request.RequestTimestamp,
		"/checkout/v1/notification",
		request.JsonBody,
	)

	if logData != nil {
		return nil, logData
	}

	if !isValid {
		errorMessage := fmt.Sprintf("Invalid signature in notification")
		logData := helper.WriteLog(errors.New(errorMessage), http.StatusUnauthorized, errorMessage)
		return nil, logData
	}

	notificationResponse := &responses.DokuPostNotificationHTTPResponse{}
	err := json.Unmarshal(request.JsonBody, &notificationResponse)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to unmarshal notification body")
		return nil, logData
	}

	return notificationResponse, nil
}

func (u *dokuUseCase) GetToken() (*responses.GetTokenResponse, *models.ErrorLog) {
	xTimestamp := time.Now().UTC().Format("2006-01-02T15:04:05-07:00")
	fmt.Printf("X-Timestamp: %s\n", xTimestamp)
	xSignature, err := generateGetTokenSignature(u.DokuPrivateKey, xTimestamp, u.DokuAPIClientID)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to generate get token signature")
		return nil, logData
	}

	requestBody := map[string]string{
		"grantType": "client_credentials",
	}
	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to marshal get token request body")
		return nil, logData
	}

	response := helper.POST(&helper.Options{
		Method: "POST",
		URL:    "https://api-sandbox.doku.com/authorization/v1/access-token/b2b",
		Headers: map[string]string{
			"X-Timestamp":  xTimestamp,
			"X-Signature":  xSignature,
			"X-Client-Key": u.DokuAPIClientID,
			"Content-Type": "application/json",
		},
		IsPrintCurl: true,
		Body:        requestBodyJson,
	})

	if response.StatusCode != http.StatusOK {
		dokuErrorResponse := &responses.DokuErrorHTTPResponse{}
		err = json.Unmarshal(response.Body, &dokuErrorResponse)
		if err != nil {
			logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to unmarshal create payment error response")
			return nil, logData
		}
		logData := helper.WriteLog(fmt.Errorf("Get Token Response Error: %v, Body: %v", response.Error, dokuErrorResponse.Message), response.StatusCode, "Failed to get token from Doku")
		return nil, logData
	}

	fmt.Printf("Get Token Response Body: %s\n", string(response.Body))
	var getTokenResponse *responses.GetTokenResponse

	err = json.Unmarshal(response.Body, &getTokenResponse)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to unmarshal get token response")
		return nil, logData
	}

	return getTokenResponse, nil
}

func (u *dokuUseCase) BankAccountInquiry(request *requests.DokuBankAccountInquiryRequest, accessToken string) (*responses.BankAccountInquiryResponse, *models.ErrorLog) {
	xTimestamp := time.Now().UTC().Format(time.RFC3339)
	xTimestamp = strings.TrimSpace(xTimestamp)
	fmt.Printf("X-Timestamp: '%s'\n", xTimestamp)

	requestBodyBytes, err := json.Marshal(request)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to marshal bank account inquiry request body")
		return nil, logData
	}
	fmt.Printf("Request Body: %s\n", string(requestBodyBytes))

	xSignature, err := generateKirimDokuRequestSignature(u.DokuAPISecretKey, "POST", "/snap/v1.1/emoney/bank-account-inquiry", accessToken, xTimestamp, requestBodyBytes)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to generate bank account inquiry signature")
		return nil, logData
	}

	xExternalID := uuid.NewString()

	response := helper.POST(&helper.Options{
		Method: "POST",
		URL:    "https://api-sandbox.doku.com/snap/v1.1/emoney/bank-account-inquiry",
		Headers: map[string]string{
			"Authorization": "Bearer " + accessToken,
			"X-TIMESTAMP":   xTimestamp,
			"X-SIGNATURE":   xSignature,
			"X-EXTERNAL-ID": xExternalID,
			"CHANNEL-ID":    "H2H",
			"X-PARTNER-ID":  u.DokuAPIClientID,
			"Content-Type":  "application/json",
		},
		IsPrintCurl: true,
		Body:        requestBodyBytes,
	})
	if response.StatusCode != http.StatusOK {
		dokuErrorResponse := &responses.DokuErrorHTTPResponse{}
		err = json.Unmarshal(response.Body, &dokuErrorResponse)
		if err != nil {
			logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to unmarshal bank account inquiry error response")
			return nil, logData
		}
		logData := helper.WriteLog(fmt.Errorf("Bank Account Inquiry Response Error: %v, Body: %v", response.Error, dokuErrorResponse.Message), response.StatusCode, "Failed to get bank account inquiry from Doku")
		return nil, logData
	}

	var bankAccountInquiryResponse *responses.BankAccountInquiryResponse

	err = json.Unmarshal(response.Body, &bankAccountInquiryResponse)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, "Failed to unmarshal bank account inquiry response")
		return nil, logData
	}

	return bankAccountInquiryResponse, nil
}
