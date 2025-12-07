# Account Management Flow - Business Logic Documentation

## Overview

The Account Management flow handles the creation and management of DOKU Sub-Accounts (SAC) and balance queries. Each merchant/user in your system needs a DOKU Sub-Account to receive payments and track balances.

---

## Sub-Account Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         SUB-ACCOUNT LIFECYCLE                                    │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐                    │
│  │   New User   │────▶│   Create     │────▶│   Active     │                    │
│  │   Signup     │     │   SAC        │     │   Account    │                    │
│  └──────────────┘     └──────────────┘     └──────────────┘                    │
│                                                   │                             │
│                                                   ▼                             │
│                                            ┌──────────────┐                    │
│                                            │   Receive    │                    │
│                                            │   Payments   │                    │
│                                            └──────────────┘                    │
│                                                   │                             │
│                                                   ▼                             │
│                                            ┌──────────────┐                    │
│                                            │   Query      │                    │
│                                            │   Balance    │                    │
│                                            └──────────────┘                    │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Create Sub-Account

### Request Structure

```go
type DokuCreateSubAccountRequest struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Email` | string | Yes | Unique email for the sub-account |
| `Name` | string | Yes | Account holder name |

### Response Structure

```go
type DokuCreateSubAccountAccountResponse struct {
    CreatedDate *time.Time  `json:"created_date"`
    UpdatedDate *time.Time  `json:"updated_date"`
    Name        null.String `json:"name"`
    Type        null.String `json:"type"`
    Status      null.String `json:"status"`
    ID          null.String `json:"id"`
}
```

| Field | Description |
|-------|-------------|
| `ID` | Sub-Account ID (e.g., "SAC-7327-1764507463535") |
| `Name` | Account name |
| `Type` | Account type (always "STANDARD") |
| `Status` | Account status |
| `CreatedDate` | Account creation timestamp |

---

## Create Sub-Account Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        CREATE SUB-ACCOUNT FLOW                                   │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌──────────────┐          ┌──────────────┐          ┌──────────────┐          │
│  │   Setter     │          │    DOKU      │          │    DOKU      │          │
│  │   Service    │          │   Module     │          │    API       │          │
│  └──────┬───────┘          └──────┬───────┘          └──────┬───────┘          │
│         │                         │                         │                   │
│         │ 1. User registers       │                         │                   │
│         │                         │                         │                   │
│         │ 2. CreateAccount()      │                         │                   │
│         │ ───────────────────────▶│                         │                   │
│         │                         │                         │                   │
│         │                         │ 3. Build request with   │                   │
│         │                         │    HMAC signature       │                   │
│         │                         │                         │                   │
│         │                         │ 4. POST /sac-merchant/  │                   │
│         │                         │    v1/accounts          │                   │
│         │                         │ ───────────────────────▶│                   │
│         │                         │                         │                   │
│         │                         │ 5. Return SAC ID        │                   │
│         │                         │ ◀───────────────────────│                   │
│         │                         │                         │                   │
│         │ 6. Return SAC response  │                         │                   │
│         │ ◀───────────────────────│                         │                   │
│         │                         │                         │                   │
│         │ 7. Store SAC ID in      │                         │                   │
│         │    user profile         │                         │                   │
│         │                         │                         │                   │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Implementation Logic

### CreateAccount Method

```go
func (u *dokuUseCase) CreateAccount(
    request *requests.DokuCreateSubAccountRequest,
) (*responses.DokuCreateSubAccountAccountResponse, *models.ErrorLog) {

    // 1. Build DOKU API request payload
    createAccountPayload := &requests.DokuCreateSubAccountHTTPRequest{
        Account: requests.DokuCreateSubAccountAccount{
            Email: request.Email,
            Type:  "STANDARD",
            Name:  request.Name,
        },
    }

    // 2. Marshal to JSON
    createAccountPayloadJson, err := json.Marshal(createAccountPayload)
    if err != nil {
        return nil, helper.WriteLog(err, http.StatusInternalServerError, 
            "Failed to marshal create account payload")
    }

    // 3. Generate signature components
    requestId := uuid.NewString()
    requestTimeStamp := time.Now().UTC()
    requestTarget := "/sac-merchant/v1/accounts"

    signature, logData := u.createSignatureComponent(
        requestId, 
        &requestTimeStamp, 
        requestTarget, 
        createAccountPayloadJson,
    )
    if logData != nil {
        return nil, logData
    }

    // 4. Prepare request headers
    requestHeader := map[string]string{
        "Client-Id":         u.DokuAPIClientID,
        "Request-Id":        requestId,
        "Request-Timestamp": requestTimeStamp.Format("2006-01-02T15:04:05Z"),
        "Signature":         signature,
    }

    // 5. Call DOKU API
    createAccountAPI := helper.POST(&helper.Options{
        Method:      "POST",
        URL:         "https://api-sandbox.doku.com/sac-merchant/v1/accounts",
        Body:        createAccountPayloadJson,
        Headers:     requestHeader,
        Timeout:     30 * time.Second,
        ContentType: "application/json",
    })

    // 6. Handle errors and existing accounts
    if createAccountAPI.StatusCode != http.StatusOK {
        // Check if email already registered
        if strings.Contains(errorMessage, "email already registered") {
            // Extract SAC ID from error message
            parts := strings.Split(errorMessage, "account id: ")
            if len(parts) == 2 {
                sacID := strings.TrimSpace(parts[1])
                return &responses.DokuCreateSubAccountAccountResponse{
                    ID: null.StringFrom(sacID),
                }, nil
            }
        }
        return nil, logData
    }

    // 7. Parse successful response
    var createAccountResponse *responses.DokuCreateSubAccountHTTPResponse
    err = json.Unmarshal(createAccountAPI.Body, &createAccountResponse)
    if err != nil {
        return nil, helper.WriteLog(err, http.StatusInternalServerError, 
            "Failed to unmarshal create account response")
    }

    return createAccountResponse.Account, nil
}
```

---

## Setter-Service Integration Example

```go
// In setter-service during user registration
func (s *userService) RegisterUser(ctx context.Context, req *RegisterRequest) (*User, error) {
    
    // 1. Create user in database
    user := &User{
        Name:  req.Name,
        Email: req.Email,
    }
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // 2. Create DOKU Sub-Account
    sacResponse, err := s.dokuUseCase.CreateAccount(&requests.DokuCreateSubAccountRequest{
        Email: req.Email,
        Name:  req.Name,
    })
    if err != nil {
        // Log error but don't fail registration
        log.Printf("Failed to create DOKU SAC: %v", err)
        return user, nil
    }
    
    // 3. Store SAC ID in user profile
    user.DokuSubAccountID = sacResponse.ID.String
    if err := s.userRepo.Update(ctx, user); err != nil {
        log.Printf("Failed to update user with SAC ID: %v", err)
    }
    
    // 4. Create Ledger Account
    ledgerAccount, err := s.ledgerAccountUseCase.CreateAccount(tx, &ledger.LedgerAccountCreateRequest{
        RandID: user.ID,
        Name:   user.Name,
        Email:  user.Email,
    })
    if err != nil {
        log.Printf("Failed to create Ledger account: %v", err)
    }
    
    user.LedgerAccountUUID = ledgerAccount.UUID
    
    return user, nil
}
```

---

## Get Balance

### Request Structure

```go
// Only requires the SAC ID as a path parameter
GetBalance(sacID string)
```

### Response Structure

```go
type DokuGetBalanceHTTPResponse struct {
    Balance *models.DokuBalance `json:"balance"`
}

type DokuBalance struct {
    Pending   null.String `json:"pending"`
    Available null.String `json:"available"`
}
```

| Field | Description |
|-------|-------------|
| `Pending` | Balance waiting for settlement (string, e.g., "50000.00") |
| `Available` | Available balance for disbursement (string, e.g., "95560.00") |

---

## Get Balance Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           GET BALANCE FLOW                                       │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌──────────────┐          ┌──────────────┐          ┌──────────────┐          │
│  │   Setter     │          │    DOKU      │          │    DOKU      │          │
│  │   Service    │          │   Module     │          │    API       │          │
│  └──────┬───────┘          └──────┬───────┘          └──────┬───────┘          │
│         │                         │                         │                   │
│         │ 1. User requests        │                         │                   │
│         │    balance              │                         │                   │
│         │                         │                         │                   │
│         │ 2. GetBalance(sacID)    │                         │                   │
│         │ ───────────────────────▶│                         │                   │
│         │                         │                         │                   │
│         │                         │ 3. Build request with   │                   │
│         │                         │    HMAC signature       │                   │
│         │                         │    (no body for GET)    │                   │
│         │                         │                         │                   │
│         │                         │ 4. GET /sac-merchant/   │                   │
│         │                         │    v1/balances/{sacID}  │                   │
│         │                         │ ───────────────────────▶│                   │
│         │                         │                         │                   │
│         │                         │ 5. Return balance       │                   │
│         │                         │ ◀───────────────────────│                   │
│         │                         │                         │                   │
│         │ 6. Return balance       │                         │                   │
│         │    response             │                         │                   │
│         │ ◀───────────────────────│                         │                   │
│         │                         │                         │                   │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Implementation Logic

### GetBalance Method

```go
func (u *dokuUseCase) GetBalance(
    sacID string,
) (*responses.DokuGetBalanceHTTPResponse, *models.ErrorLog) {

    // 1. Generate signature components (no body for GET requests)
    requestId := uuid.NewString()
    requestTimeStamp := time.Now().UTC()
    requestTarget := fmt.Sprintf("/sac-merchant/v1/balances/%s", sacID)

    signature, logData := u.createSignatureComponent(
        requestId, 
        &requestTimeStamp, 
        requestTarget, 
        nil, // No body for GET requests
    )
    if logData != nil {
        return nil, logData
    }

    // 2. Prepare request headers
    requestHeader := map[string]string{
        "Client-Id":         u.DokuAPIClientID,
        "Request-Id":        requestId,
        "Request-Timestamp": requestTimeStamp.Format("2006-01-02T15:04:05Z"),
        "Signature":         signature,
    }

    // 3. Call DOKU API
    getBalanceAPI := helper.GET(&helper.Options{
        Method:      "GET",
        URL:         fmt.Sprintf("https://api-sandbox.doku.com/sac-merchant/v1/balances/%s", sacID),
        Headers:     requestHeader,
        Timeout:     30 * time.Second,
        ContentType: "application/json",
    })

    // 4. Handle errors
    if getBalanceAPI.Error != nil {
        return nil, helper.WriteLog(getBalanceAPI.Error, 
            getBalanceAPI.StatusCode, 
            helper.DefaultStatusText[getBalanceAPI.StatusCode])
    }

    if getBalanceAPI.StatusCode != http.StatusOK {
        dokuErrorResponse := &responses.DokuErrorHTTPResponse{}
        json.Unmarshal(getBalanceAPI.Body, &dokuErrorResponse)
        return nil, helper.WriteLog(
            fmt.Errorf("DOKU Error: %v", dokuErrorResponse.Message),
            getBalanceAPI.StatusCode,
            fmt.Sprintf("Doku Get Balance API Error: %v", dokuErrorResponse.Message),
        )
    }

    // 5. Parse successful response
    var getBalanceResponse *responses.DokuGetBalanceHTTPResponse
    err := json.Unmarshal(getBalanceAPI.Body, &getBalanceResponse)
    if err != nil {
        return nil, helper.WriteLog(err, http.StatusInternalServerError, 
            "Failed to unmarshal get balance response")
    }

    return getBalanceResponse, nil
}
```

---

## Setter-Service Balance Query Example

```go
// In setter-service for dashboard balance display
func (s *dashboardService) GetUserBalance(ctx context.Context, userID string) (*BalanceResponse, error) {
    
    // 1. Get user with DOKU SAC ID
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 2. Query DOKU balance
    dokuBalance, err := s.dokuUseCase.GetBalance(user.DokuSubAccountID)
    if err != nil {
        log.Printf("Failed to get DOKU balance: %v", err)
        // Fallback to Ledger balance
    }
    
    // 3. Query Ledger balance (internal records)
    ledgerBalance, err := s.ledgerWalletUseCase.GetCurrentBalanceByAccount(
        user.LedgerAccountUUID,
        "IDR",
    )
    if err != nil {
        return nil, err
    }
    
    // 4. Return combined balance view
    return &BalanceResponse{
        // From Ledger (source of truth for your system)
        AvailableBalance: ledgerBalance.AvailableBalance,
        PendingBalance:   ledgerBalance.PendingBalance,
        TotalIncome:      ledgerBalance.TotalIncome,
        TotalWithdrawn:   ledgerBalance.TotalWithdrawn,
        
        // From DOKU (for verification/comparison)
        DokuPending:   dokuBalance.Balance.Pending.String,
        DokuAvailable: dokuBalance.Balance.Available.String,
    }, nil
}
```

---

## API Request Examples

### Create Sub-Account Request

```
POST /sac-merchant/v1/accounts HTTP/1.1
Host: api-sandbox.doku.com
Content-Type: application/json
Client-Id: MCH-0001-1234567890
Request-Id: 550e8400-e29b-41d4-a716-446655440000
Request-Timestamp: 2025-01-15T10:30:00Z
Signature: HMACSHA256=xxxxxxxxxxxxxxxxxxxxxx
```

```json
{
  "account": {
    "email": "merchant@example.com",
    "type": "STANDARD",
    "name": "John's Store"
  }
}
```

### Create Sub-Account Response (Success)

```json
{
  "account": {
    "created_date": "2025-01-15T10:30:00Z",
    "updated_date": "2025-01-15T10:30:00Z",
    "name": "John's Store",
    "type": "STANDARD",
    "status": "ACTIVE",
    "id": "SAC-7327-1764507463535"
  }
}
```

### Create Sub-Account Response (Email Already Exists)

```json
{
  "error": {
    "message": "email already registered with account id: SAC-7327-1764507463535"
  }
}
```

### Get Balance Request

```
GET /sac-merchant/v1/balances/SAC-7327-1764507463535 HTTP/1.1
Host: api-sandbox.doku.com
Client-Id: MCH-0001-1234567890
Request-Id: 661e9511-b30c-42d8-b789-123456789abc
Request-Timestamp: 2025-01-15T11:00:00Z
Signature: HMACSHA256=yyyyyyyyyyyyyyyyyyyy
```

### Get Balance Response

```json
{
  "balance": {
    "pending": "50000.00",
    "available": "95560.00"
  }
}
```

---

## Balance Explanation

| DOKU Balance Field | Description | Ledger Equivalent |
|--------------------|-------------|-------------------|
| `pending` | Payments confirmed but not yet settled | `pending_balance` |
| `available` | Funds settled and ready for disbursement | `balance` |

---

## Error Handling

### Account Creation Errors

| Error | Status Code | Description | Action |
|-------|-------------|-------------|--------|
| Email already registered | 400 | Email exists in DOKU | Extract existing SAC ID |
| Invalid email format | 400 | Email format invalid | Validate before calling |
| Invalid signature | 401 | Signature mismatch | Check credentials |

### Balance Query Errors

| Error | Status Code | Description | Action |
|-------|-------------|-------------|--------|
| Account not found | 404 | SAC ID doesn't exist | Verify SAC ID |
| Invalid signature | 401 | Signature mismatch | Check credentials |

---

## Best Practices

1. **Store SAC ID**: Always store the DOKU Sub-Account ID in your user database for future API calls.

2. **Handle Existing Accounts**: DOKU returns an error if email is already registered, but includes the existing SAC ID. Parse this for idempotent account creation.

3. **Balance Comparison**: Query both DOKU and Ledger balances for reconciliation and verification.

4. **Signature for GET Requests**: GET requests don't include a body, so the signature components exclude the Digest field.

5. **Create Account During Registration**: Create DOKU Sub-Account during user registration to ensure it's ready when they receive their first payment.
