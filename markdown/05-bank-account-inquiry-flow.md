# Bank Account Inquiry Flow - Business Logic Documentation

## Overview

The Bank Account Inquiry flow handles the verification of destination bank accounts before initiating a disbursement ("KIRIM DOKU"). This is a critical step to ensure funds are sent to valid, verified bank accounts and to display the correct beneficiary name to the user for confirmation.

---

## Bank Account Inquiry Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        BANK ACCOUNT INQUIRY FLOW                                 │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌────────────┐ │
│  │    User      │     │   Setter     │     │    DOKU      │     │   Bank     │ │
│  │   (Client)   │     │   Service    │     │   Module     │     │   System   │ │
│  └──────┬───────┘     └──────┬───────┘     └──────┬───────┘     └──────┬─────┘ │
│         │                    │                    │                    │        │
│         │ 1. Enter bank      │                    │                    │        │
│         │    account details │                    │                    │        │
│         │ ──────────────────▶│                    │                    │        │
│         │                    │                    │                    │        │
│         │                    │ 2. GetToken()      │                    │        │
│         │                    │ ──────────────────▶│                    │        │
│         │                    │                    │                    │        │
│         │                    │ 3. Return          │                    │        │
│         │                    │    access_token    │                    │        │
│         │                    │ ◀──────────────────│                    │        │
│         │                    │                    │                    │        │
│         │                    │ 4. BankAccount     │                    │        │
│         │                    │    Inquiry()       │                    │        │
│         │                    │ ──────────────────▶│                    │        │
│         │                    │                    │                    │        │
│         │                    │                    │ 5. Verify with     │        │
│         │                    │                    │    bank system     │        │
│         │                    │                    │ ──────────────────▶│        │
│         │                    │                    │                    │        │
│         │                    │                    │ 6. Return account  │        │
│         │                    │                    │    holder name     │        │
│         │                    │                    │ ◀──────────────────│        │
│         │                    │                    │                    │        │
│         │                    │ 7. Return inquiry  │                    │        │
│         │                    │    result          │                    │        │
│         │                    │ ◀──────────────────│                    │        │
│         │                    │                    │                    │        │
│         │ 8. Show beneficiary│                    │                    │        │
│         │    name for        │                    │                    │        │
│         │    confirmation    │                    │                    │        │
│         │ ◀──────────────────│                    │                    │        │
│         │                    │                    │                    │        │
│         │ 9. User confirms   │                    │                    │        │
│         │    disbursement    │                    │                    │        │
│         │                    │                    │                    │        │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## SNAP API Authentication

Before calling Bank Account Inquiry, you need to obtain an access token using the SNAP API authentication flow.

### GetToken Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            GET TOKEN FLOW (SNAP API)                             │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  1. Generate timestamp: 2025-01-15T10:30:00+00:00                              │
│                                                                                 │
│  2. Create signature string:                                                    │
│     "{client_id}|{timestamp}"                                                   │
│                                                                                 │
│  3. Sign with RSA-SHA256 using private key                                     │
│                                                                                 │
│  4. POST to /authorization/v1/access-token/b2b                                 │
│     Headers:                                                                    │
│       - X-Timestamp: {timestamp}                                               │
│       - X-Signature: {rsa_signature}                                           │
│       - X-Client-Key: {client_id}                                              │
│     Body:                                                                       │
│       - grantType: "client_credentials"                                        │
│                                                                                 │
│  5. Receive access_token (valid for limited time)                              │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Request Structure

### DokuBankAccountInquiryRequest

```go
type DokuBankAccountInquiryRequest struct {
    PartnerReferenceNo string `json:"partnerReferenceNo"`
    CustomerNumber     string `json:"customerNumber"`
    Amount             struct {
        Value    string `json:"value"`
        Currency string `json:"currency"`
    } `json:"amount"`
    BeneficiaryAccountNumber string `json:"beneficiaryAccountNumber"`
    AdditionalInfo           struct {
        BeneficiaryBankCode    string `json:"beneficiaryBankCode"`
        BeneficiaryAccountName string `json:"beneficiaryAccountName"`
        SenderCountryCode      string `json:"senderCountryCode"`
    } `json:"additionalInfo"`
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `PartnerReferenceNo` | string | Yes | Unique reference number from your system |
| `CustomerNumber` | string | Yes | Customer identifier |
| `Amount.Value` | string | Yes | Transfer amount (e.g., "100000.00") |
| `Amount.Currency` | string | Yes | Currency code (e.g., "IDR") |
| `BeneficiaryAccountNumber` | string | Yes | Destination bank account number |
| `AdditionalInfo.BeneficiaryBankCode` | string | Yes | Bank code (e.g., "002" for BRI) |
| `AdditionalInfo.BeneficiaryAccountName` | string | No | Expected account holder name |
| `AdditionalInfo.SenderCountryCode` | string | Yes | Sender country (e.g., "ID") |

---

## Response Structure

### BankAccountInquiryResponse

```go
type BankAccountInquiryResponse struct {
    ResponseCode             int    `json:"responseCode"`
    ResponseMessage          string `json:"responseMessage"`
    ReferenceNo              string `json:"referenceNo"`
    PartnerReferenceNo       string `json:"partnerReferenceNo"`
    BeneficiaryAccountNumber string `json:"beneficiaryAccountNumber"`
    BeneficiaryAccountName   string `json:"beneficiaryAccountName"`
    BeneficiaryBankCode      string `json:"beneficiaryBankCode"`
    BeneficiaryBankShortName string `json:"beneficiaryBankShortName"`
    BeneficiaryBankName      string `json:"beneficiaryBankName"`
    Amount                   struct {
        Value    string `json:"value"`
        Currency string `json:"currency"`
    } `json:"amount"`
    SessionID      string `json:"sessionId"`
    AdditionalInfo struct {
        SenderCountryCode   string `json:"senderCountryCode"`
        ForexRate           string `json:"forexRate"`
        ForexOriginCurrency string `json:"forexOriginCurrency"`
        FeeAmount           string `json:"feeAmount"`
        FeeCurrency         string `json:"feeCurrency"`
        BeneficiaryAmount   string `json:"beneficiaryAmount"`
        ReferenceNumber     string `json:"referenceNumber"`
    } `json:"additionalInfo"`
}
```

### Key Response Fields

| Field | Description |
|-------|-------------|
| `BeneficiaryAccountName` | Verified account holder name from bank |
| `BeneficiaryBankName` | Full bank name |
| `BeneficiaryBankShortName` | Short bank name |
| `SessionID` | Session ID for subsequent disbursement |
| `AdditionalInfo.FeeAmount` | Transfer fee (if applicable) |

---

## GetToken Response Structure

### GetTokenResponse

```go
type GetTokenResponse struct {
    ResponseCode    string `json:"responseCode"`
    ResponseMessage string `json:"responseMessage"`
    AccessToken     string `json:"accessToken"`
    TokenType       string `json:"tokenType"`
    ExpiresIn       int    `json:"expiresIn"`
    AdditionalInfo  string `json:"additionalInfo"`
}
```

| Field | Description |
|-------|-------------|
| `AccessToken` | Bearer token for SNAP API calls |
| `TokenType` | Token type (usually "Bearer") |
| `ExpiresIn` | Token validity in seconds |

---

## Implementation Logic

### GetToken Method

```go
func (u *dokuUseCase) GetToken() (*responses.GetTokenResponse, *models.ErrorLog) {
    
    // 1. Generate timestamp
    xTimestamp := time.Now().UTC().Format("2006-01-02T15:04:05-07:00")
    
    // 2. Generate RSA-SHA256 signature
    xSignature, err := generateGetTokenSignature(
        u.DokuPrivateKey, 
        xTimestamp, 
        u.DokuAPIClientID,
    )
    if err != nil {
        return nil, helper.WriteLog(err, http.StatusInternalServerError, 
            "Failed to generate get token signature")
    }

    // 3. Build request body
    requestBody := map[string]string{
        "grantType": "client_credentials",
    }
    requestBodyJson, _ := json.Marshal(requestBody)

    // 4. Call DOKU API
    response := helper.POST(&helper.Options{
        Method: "POST",
        URL:    "https://api-sandbox.doku.com/authorization/v1/access-token/b2b",
        Headers: map[string]string{
            "X-Timestamp":  xTimestamp,
            "X-Signature":  xSignature,
            "X-Client-Key": u.DokuAPIClientID,
            "Content-Type": "application/json",
        },
        Body: requestBodyJson,
    })

    // 5. Handle errors
    if response.StatusCode != http.StatusOK {
        return nil, helper.WriteLog(
            fmt.Errorf("Get Token Error: %v", response.Error),
            response.StatusCode,
            "Failed to get token from Doku",
        )
    }

    // 6. Parse response
    var getTokenResponse *responses.GetTokenResponse
    err = json.Unmarshal(response.Body, &getTokenResponse)
    if err != nil {
        return nil, helper.WriteLog(err, http.StatusInternalServerError, 
            "Failed to unmarshal get token response")
    }

    return getTokenResponse, nil
}
```

### BankAccountInquiry Method

```go
func (u *dokuUseCase) BankAccountInquiry(
    request *requests.DokuBankAccountInquiryRequest, 
    accessToken string,
) (*responses.BankAccountInquiryResponse, *models.ErrorLog) {
    
    // 1. Generate timestamp
    xTimestamp := time.Now().UTC().Format(time.RFC3339)

    // 2. Marshal request body
    requestBodyBytes, err := json.Marshal(request)
    if err != nil {
        return nil, helper.WriteLog(err, http.StatusInternalServerError, 
            "Failed to marshal bank account inquiry request body")
    }

    // 3. Generate HMAC-SHA512 signature for SNAP API
    xSignature, err := generateKirimDokuRequestSignature(
        u.DokuAPISecretKey, 
        "POST", 
        "/snap/v1.1/emoney/bank-account-inquiry", 
        accessToken, 
        xTimestamp, 
        requestBodyBytes,
    )
    if err != nil {
        return nil, helper.WriteLog(err, http.StatusInternalServerError, 
            "Failed to generate bank account inquiry signature")
    }

    // 4. Generate unique external ID
    xExternalID := uuid.NewString()

    // 5. Call DOKU API
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
        Body: requestBodyBytes,
    })

    // 6. Handle errors
    if response.StatusCode != http.StatusOK {
        dokuErrorResponse := &responses.DokuErrorHTTPResponse{}
        json.Unmarshal(response.Body, &dokuErrorResponse)
        return nil, helper.WriteLog(
            fmt.Errorf("Bank Account Inquiry Error: %v", dokuErrorResponse.Message),
            response.StatusCode,
            "Failed to get bank account inquiry from Doku",
        )
    }

    // 7. Parse response
    var bankAccountInquiryResponse *responses.BankAccountInquiryResponse
    err = json.Unmarshal(response.Body, &bankAccountInquiryResponse)
    if err != nil {
        return nil, helper.WriteLog(err, http.StatusInternalServerError, 
            "Failed to unmarshal bank account inquiry response")
    }

    return bankAccountInquiryResponse, nil
}
```

---

## Setter-Service Integration Example

```go
// In setter-service for disbursement flow
func (s *disbursementService) VerifyBankAccount(
    ctx context.Context, 
    req *VerifyBankAccountRequest,
) (*VerifyBankAccountResponse, error) {
    
    // 1. Get access token
    tokenResponse, err := s.dokuUseCase.GetToken()
    if err != nil {
        return nil, fmt.Errorf("failed to get token: %w", err)
    }
    
    // 2. Prepare inquiry request
    inquiryRequest := &requests.DokuBankAccountInquiryRequest{
        PartnerReferenceNo:       fmt.Sprintf("REF-%s-%d", req.UserID, time.Now().Unix()),
        CustomerNumber:           req.UserID,
        BeneficiaryAccountNumber: req.BankAccountNumber,
    }
    inquiryRequest.Amount.Value = fmt.Sprintf("%.2f", req.Amount)
    inquiryRequest.Amount.Currency = "IDR"
    inquiryRequest.AdditionalInfo.BeneficiaryBankCode = req.BankCode
    inquiryRequest.AdditionalInfo.SenderCountryCode = "ID"
    
    // 3. Call bank account inquiry
    inquiryResponse, err := s.dokuUseCase.BankAccountInquiry(
        inquiryRequest, 
        tokenResponse.AccessToken,
    )
    if err != nil {
        return nil, fmt.Errorf("bank account inquiry failed: %w", err)
    }
    
    // 4. Return verification result
    return &VerifyBankAccountResponse{
        BankAccountNumber: inquiryResponse.BeneficiaryAccountNumber,
        BankAccountName:   inquiryResponse.BeneficiaryAccountName,
        BankName:          inquiryResponse.BeneficiaryBankName,
        BankCode:          inquiryResponse.BeneficiaryBankCode,
        SessionID:         inquiryResponse.SessionID,
        IsVerified:        true,
    }, nil
}

// Complete disbursement after user confirmation
func (s *disbursementService) CreateDisbursement(
    ctx context.Context, 
    req *CreateDisbursementRequest,
) (*DisbursementResponse, error) {
    
    // 1. Verify bank account first
    verification, err := s.VerifyBankAccount(ctx, &VerifyBankAccountRequest{
        UserID:            req.UserID,
        BankAccountNumber: req.BankAccountNumber,
        BankCode:          req.BankCode,
        Amount:            req.Amount,
    })
    if err != nil {
        return nil, err
    }
    
    // 2. Create disbursement in Ledger
    _, err = s.ledgerDisbursementUseCase.CreateDisbursement(tx, &ledger.LedgerDisbursementCreateRequest{
        LedgerAccountUUID:     req.LedgerAccountUUID,
        LedgerWalletUUID:      req.LedgerWalletUUID,
        LedgerAccountBankUUID: req.LedgerAccountBankUUID,
        Amount:                int64(req.Amount),
        Currency:              "IDR",
    })
    if err != nil {
        return nil, err
    }
    
    // 3. Initiate actual transfer via DOKU SNAP API
    // (Implementation depends on DOKU's disbursement API)
    
    return &DisbursementResponse{
        Status:            "PENDING",
        BeneficiaryName:   verification.BankAccountName,
        BeneficiaryBank:   verification.BankName,
        BeneficiaryAccount: verification.BankAccountNumber,
    }, nil
}
```

---

## Bank Code Reference

Common Indonesian bank codes:

| Bank | Code | Short Name |
|------|------|------------|
| Bank Central Asia (BCA) | 014 | BCA |
| Bank Rakyat Indonesia (BRI) | 002 | BRI |
| Bank Mandiri | 008 | MANDIRI |
| Bank Negara Indonesia (BNI) | 009 | BNI |
| Bank CIMB Niaga | 022 | CIMB |
| Bank Danamon | 011 | DANAMON |
| Bank Permata | 013 | PERMATA |
| Bank Syariah Indonesia (BSI) | 451 | BSI |
| Bank OCBC NISP | 028 | OCBC |
| Bank Tabungan Negara (BTN) | 200 | BTN |
| Bank Mega | 426 | MEGA |
| Bank Jago | 542 | JAGO |
| Bank Jenius (BTPN) | 213 | BTPN |
| SeaBank | 535 | SEABANK |
| Bank Digital BCA (Blu) | 501 | BLU |

---

## API Request Examples

### Get Token Request

```
POST /authorization/v1/access-token/b2b HTTP/1.1
Host: api-sandbox.doku.com
Content-Type: application/json
X-Timestamp: 2025-01-15T10:30:00+00:00
X-Signature: {rsa_sha256_signature}
X-Client-Key: MCH-0001-1234567890
```

```json
{
  "grantType": "client_credentials"
}
```

### Get Token Response

```json
{
  "responseCode": "2007300",
  "responseMessage": "Successful",
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "tokenType": "Bearer",
  "expiresIn": 900
}
```

### Bank Account Inquiry Request

```
POST /snap/v1.1/emoney/bank-account-inquiry HTTP/1.1
Host: api-sandbox.doku.com
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
X-TIMESTAMP: 2025-01-15T10:35:00Z
X-SIGNATURE: {hmac_sha512_signature}
X-EXTERNAL-ID: 550e8400-e29b-41d4-a716-446655440000
CHANNEL-ID: H2H
X-PARTNER-ID: MCH-0001-1234567890
```

```json
{
  "partnerReferenceNo": "REF-USER001-1736939700",
  "customerNumber": "USER001",
  "amount": {
    "value": "100000.00",
    "currency": "IDR"
  },
  "beneficiaryAccountNumber": "1234567890",
  "additionalInfo": {
    "beneficiaryBankCode": "014",
    "beneficiaryAccountName": "",
    "senderCountryCode": "ID"
  }
}
```

### Bank Account Inquiry Response (Success)

```json
{
  "responseCode": 200,
  "responseMessage": "Successful",
  "referenceNo": "DOKU-REF-123456",
  "partnerReferenceNo": "REF-USER001-1736939700",
  "beneficiaryAccountNumber": "1234567890",
  "beneficiaryAccountName": "JOHN DOE",
  "beneficiaryBankCode": "014",
  "beneficiaryBankShortName": "BCA",
  "beneficiaryBankName": "BANK CENTRAL ASIA",
  "amount": {
    "value": "100000.00",
    "currency": "IDR"
  },
  "sessionId": "SESSION-789xyz",
  "additionalInfo": {
    "senderCountryCode": "ID",
    "forexRate": "1.00",
    "forexOriginCurrency": "IDR",
    "feeAmount": "0.00",
    "feeCurrency": "IDR",
    "beneficiaryAmount": "100000.00",
    "referenceNumber": "REF-123"
  }
}
```

---

## Error Handling

### Common Errors

| Error Code | Description | Action |
|------------|-------------|--------|
| 400 | Invalid request parameters | Validate input before calling |
| 401 | Invalid/expired token | Refresh access token |
| 404 | Account not found | Bank account doesn't exist |
| 500 | Bank system error | Retry or inform user |

### Error Response Example

```json
{
  "responseCode": 404,
  "responseMessage": "Account not found",
  "referenceNo": "",
  "partnerReferenceNo": "REF-USER001-1736939700"
}
```

---

## Complete Disbursement Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        COMPLETE DISBURSEMENT FLOW                                │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  1. User enters bank details (account number, bank code)                       │
│         │                                                                       │
│         ▼                                                                       │
│  2. Get SNAP Access Token                                                       │
│     dokuUseCase.GetToken()                                                      │
│         │                                                                       │
│         ▼                                                                       │
│  3. Verify Bank Account                                                         │
│     dokuUseCase.BankAccountInquiry()                                           │
│         │                                                                       │
│         ▼                                                                       │
│  4. Display beneficiary name to user for confirmation                          │
│     "Transfer to: JOHN DOE - BCA 1234567890"                                   │
│         │                                                                       │
│         ▼                                                                       │
│  5. User confirms disbursement                                                  │
│         │                                                                       │
│         ▼                                                                       │
│  6. Create Ledger Disbursement (status: PENDING)                               │
│     - Deduct from wallet.balance                                               │
│         │                                                                       │
│         ▼                                                                       │
│  7. Initiate transfer via DOKU SNAP API                                        │
│     - Use sessionId from bank account inquiry                                  │
│         │                                                                       │
│         ▼                                                                       │
│  8. Update Ledger Disbursement (status: PROCESSING)                            │
│         │                                                                       │
│         ▼                                                                       │
│  9. Receive transfer callback from DOKU                                        │
│         │                                                                       │
│         ▼                                                                       │
│  10. Update Ledger Disbursement (status: SUCCESS or FAILED)                    │
│      - If FAILED, refund to wallet.balance                                     │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Wallet Impact Summary

| Action | balance | withdraw_accumulation |
|--------|---------|----------------------|
| Create Disbursement | -amount | - |
| Disbursement Success | - | +amount |
| Disbursement Failed | +amount (refund) | - |

---

## Best Practices

1. **Always Verify First**: Never initiate a disbursement without verifying the bank account first.

2. **Show Beneficiary Name**: Display the verified account holder name to the user for confirmation before proceeding.

3. **Token Caching**: Cache the access token and refresh only when expired (check `expiresIn`).

4. **Session ID**: Store the `sessionId` from bank account inquiry for use in the actual disbursement API call.

5. **Idempotent Reference Numbers**: Use unique `partnerReferenceNo` for each inquiry to prevent duplicates.

6. **Handle Bank Downtime**: Some banks may have maintenance windows. Implement retry logic with exponential backoff.

7. **Validate Bank Code**: Ensure the bank code is valid before calling the API to avoid unnecessary API calls.

---

## Security Considerations

1. **Private Key Protection**: The RSA private key for SNAP API authentication must be stored securely (e.g., environment variable, secrets manager).

2. **Token Security**: Never expose access tokens in client-side code or logs.

3. **Input Validation**: Validate bank account numbers (typically 10-16 digits) and bank codes before API calls.

4. **Rate Limiting**: Implement rate limiting to prevent abuse of the bank account inquiry API.
