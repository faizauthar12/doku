# Create Payment Flow - Business Logic Documentation

## Overview

The Create Payment flow handles the generation of payment links for customers through the DOKU payment gateway. When a payment link is created, customers can use various payment methods (Virtual Account, E-Wallet, QRIS, etc.) to complete their transaction.

---

## Payment Creation Flow Diagram

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│ CREATED  │────▶│ PENDING  │────▶│   PAID   │     │  FAILED  │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                      │                                  ▲
                      │                                  │
                      └──────────────────────────────────┘
                      │
                      ▼
                 ┌──────────┐
                 │ EXPIRED  │
                 └──────────┘
```

---

## Request Structure

### DokuCreatePaymentRequest

```go
type DokuCreatePaymentRequest struct {
    Amount         int64  `json:"amount"`
    CustomerName   string `json:"customer_name"`
    CustomerEmail  string `json:"customer_email"`
    SacID          string `json:"SacID"`
    PaymentDueDate int64  `json:"payment_due_date,omitempty"`
    InvoiceNumber  string `json:"invoice_number"`
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Amount` | int64 | Yes | Payment amount in smallest currency unit |
| `CustomerName` | string | Yes | Customer's full name |
| `CustomerEmail` | string | Yes | Customer's email address |
| `SacID` | string | Yes | DOKU Sub-Account ID (merchant identifier) |
| `PaymentDueDate` | int64 | No | Payment expiration in minutes |
| `InvoiceNumber` | string | Yes | Unique invoice/order ID from your system |

---

## Response Structure

### DokuCreatePaymentHTTPResponse

```go
type DokuCreatePaymentHTTPResponse struct {
    Response struct {
        Order          *models.DokuOrder                 `json:"order"`
        Payment        *models.DokuPayment               `json:"payment"`
        AdditionalInfo *models.DokuPaymentAdditionalInfo `json:"additional_info"`
        Headers        *models.DokuHeader                `json:"headers"`
    } `json:"response"`
}
```

### Key Response Fields

| Field | Description |
|-------|-------------|
| `response.order.invoice_number` | Your invoice number |
| `response.order.amount` | Payment amount |
| `response.payment.token_id` | DOKU session token |
| `response.payment.url` | Checkout URL for customer |
| `response.payment.expired_date` | Payment expiration datetime |
| `response.headers.request_id` | DOKU request ID (for webhook matching) |

---

## Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           CREATE PAYMENT FLOW                                    │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌──────────────┐          ┌──────────────┐          ┌──────────────┐          │
│  │   Frontend   │          │   Setter     │          │    DOKU      │          │
│  │   (Client)   │          │   Service    │          │   Module     │          │
│  └──────┬───────┘          └──────┬───────┘          └──────┬───────┘          │
│         │                         │                         │                   │
│         │ 1. Request payment      │                         │                   │
│         │    (product, amount)    │                         │                   │
│         │ ───────────────────────▶│                         │                   │
│         │                         │                         │                   │
│         │                         │ 2. AcceptPayment()      │                   │
│         │                         │ ───────────────────────▶│                   │
│         │                         │                         │                   │
│         │                         │                         │ 3. Build request  │
│         │                         │                         │    with signature │
│         │                         │                         │                   │
│         │                         │                         │ 4. POST to DOKU   │
│         │                         │                         │    /checkout/v1/  │
│         │                         │                         │    payment        │
│         │                         │                         │ ─────────────────▶│
│         │                         │                         │                   │
│         │                         │                         │ 5. Receive        │
│         │                         │                         │    payment_url,   │
│         │                         │                         │    token_id       │
│         │                         │                         │ ◀─────────────────│
│         │                         │                         │                   │
│         │                         │ 6. Return response      │                   │
│         │                         │ ◀───────────────────────│                   │
│         │                         │                         │                   │
│         │                         │ 7. Create LedgerPayment │                   │
│         │                         │    (status: PENDING)    │                   │
│         │                         │    Store gateway refs   │                   │
│         │                         │                         │                   │
│         │ 8. Redirect to          │                         │                   │
│         │    payment_url          │                         │                   │
│         │ ◀───────────────────────│                         │                   │
│         │                         │                         │                   │
│         │ 9. Customer completes   │                         │                   │
│         │    payment on DOKU      │                         │                   │
│         │    checkout page        │                         │                   │
│         │                         │                         │                   │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Implementation Logic

### AcceptPayment Method

```go
func (u *dokuUseCase) AcceptPayment(
    request *requests.DokuCreatePaymentRequest,
) (*responses.DokuCreatePaymentHTTPResponse, *models.ErrorLog) {

    // 1. Build DOKU API request payload
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

    // 2. Marshal to JSON
    createPaymentPayloadJson, err := json.Marshal(createPaymentPayload)
    if err != nil {
        return nil, helper.WriteLog(err, http.StatusInternalServerError, 
            "Failed to marshal create payment payload")
    }

    // 3. Generate signature components
    requestId := uuid.NewString()
    requestTimeStamp := time.Now().UTC()
    requestTarget := "/checkout/v1/payment"

    signature, logData := u.createSignatureComponent(
        requestId, 
        &requestTimeStamp, 
        requestTarget, 
        createPaymentPayloadJson,
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
    createPaymentAPI := helper.POST(&helper.Options{
        Method:      "POST",
        URL:         "https://api-sandbox.doku.com/checkout/v1/payment",
        Body:        createPaymentPayloadJson,
        Headers:     requestHeader,
        Timeout:     30 * time.Second,
        ContentType: "application/json",
    })

    // 6. Handle errors
    if createPaymentAPI.Error != nil {
        return nil, helper.WriteLog(createPaymentAPI.Error, 
            createPaymentAPI.StatusCode, 
            helper.DefaultStatusText[createPaymentAPI.StatusCode])
    }

    if createPaymentAPI.StatusCode != http.StatusOK {
        // Parse error response
        dokuErrorResponse := &responses.DokuErrorHTTPResponse{}
        json.Unmarshal(createPaymentAPI.Body, &dokuErrorResponse)
        return nil, helper.WriteLog(
            fmt.Errorf("DOKU Error: %v", dokuErrorResponse.Message),
            createPaymentAPI.StatusCode,
            fmt.Sprintf("Doku Create Payment API Error: %v", dokuErrorResponse.Message),
        )
    }

    // 7. Parse successful response
    var createPaymentResponse *responses.DokuCreatePaymentHTTPResponse
    err = json.Unmarshal(createPaymentAPI.Body, &createPaymentResponse)
    if err != nil {
        return nil, helper.WriteLog(err, http.StatusInternalServerError, 
            "Failed to unmarshal create payment response")
    }

    return createPaymentResponse, nil
}
```

---

## Setter-Service Integration Example

```go
// In setter-service
func (s *paymentService) CreatePayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
    
    // 1. Get or create DOKU Sub-Account
    sacID := user.DokuSubAccountID
    
    // 2. Generate invoice number
    invoiceNumber := fmt.Sprintf("INV-%s-%d", user.ID, time.Now().Unix())
    
    // 3. Call DOKU module to create payment
    dokuResponse, err := s.dokuUseCase.AcceptPayment(&requests.DokuCreatePaymentRequest{
        Amount:         req.Amount,
        CustomerName:   user.Name,
        CustomerEmail:  user.Email,
        SacID:          sacID,
        PaymentDueDate: 1440, // 24 hours
        InvoiceNumber:  invoiceNumber,
    })
    if err != nil {
        return nil, err
    }
    
    // 4. Create LedgerPayment record
    ledgerPayment, err := s.ledgerPaymentUseCase.CreatePayment(tx, &ledger.LedgerPaymentCreatePaymentRequest{
        LedgerAccountUUID: user.LedgerAccountUUID,
        InvoiceNumber:     invoiceNumber,
        Amount:            req.Amount,
        Currency:          "IDR",
        GatewayRequestId:  dokuResponse.Response.Headers.RequestID.String,
        GatewayTokenId:    dokuResponse.Response.Payment.TokenID.String,
        GatewayPaymentUrl: dokuResponse.Response.Payment.URL.String,
        ExpiresAt:         parseExpiredDate(dokuResponse.Response.Payment.ExpiredDate.String),
    })
    if err != nil {
        return nil, err
    }
    
    // 5. Return payment URL to frontend
    return &PaymentResponse{
        InvoiceNumber: invoiceNumber,
        PaymentURL:    dokuResponse.Response.Payment.URL.String,
        ExpiredAt:     ledgerPayment.ExpiresAt,
    }, nil
}
```

---

## DOKU API Request Example

### HTTP Request
```
POST /checkout/v1/payment HTTP/1.1
Host: api-sandbox.doku.com
Content-Type: application/json
Client-Id: MCH-0001-1234567890
Request-Id: 550e8400-e29b-41d4-a716-446655440000
Request-Timestamp: 2025-01-15T10:30:00Z
Signature: HMACSHA256=xxxxxxxxxxxxxxxxxxxxxx
```

### Request Body
```json
{
  "order": {
    "invoice_number": "INV-USER001-1736939400",
    "amount": 100000
  },
  "payment": {
    "payment_due_date": 1440
  },
  "customer": {
    "name": "John Doe",
    "email": "john.doe@example.com"
  },
  "additional_info": {
    "account": {
      "id": "SAC-7327-1764507463535"
    }
  }
}
```

### Success Response
```json
{
  "order": {
    "invoice_number": "INV-USER001-1736939400",
    "amount": 100000,
    "currency": "IDR",
    "session_id": "abc123xyz"
  },
  "payment": {
    "token_id": "tok_xyz123abc",
    "url": "https://checkout.doku.com/pay?token=tok_xyz123abc",
    "expired_date": "2025-01-16T10:30:00Z"
  },
  "headers": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

---

## Error Handling

### Common Error Responses

| Status Code | Error | Description |
|-------------|-------|-------------|
| 400 | Invalid Amount | Amount must be greater than 0 |
| 400 | Invalid Email | Email format is invalid |
| 401 | Invalid Signature | Signature verification failed |
| 404 | Account Not Found | Sub-Account ID not found |
| 409 | Duplicate Invoice | Invoice number already exists |
| 500 | Internal Server Error | DOKU server error |

### Error Response Structure
```json
{
  "message": ["Invalid amount: must be greater than 0"]
}
```

---

## Best Practices

1. **Invoice Number Uniqueness**: Generate unique invoice numbers using a combination of user ID, timestamp, and random string.

2. **Payment Expiration**: Set appropriate `payment_due_date` based on business requirements (e.g., 1440 minutes = 24 hours).

3. **Store Gateway References**: Always store `request_id`, `token_id`, and `url` for webhook matching and debugging.

4. **Handle Duplicate Requests**: If customer requests payment again for same order, return existing payment URL instead of creating new one.

5. **Timeout Handling**: Set appropriate timeout (30 seconds recommended) for DOKU API calls.

---

## Wallet Impact Summary

| Action | pending_balance | balance | income_accumulation |
|--------|-----------------|---------|---------------------|
| Create Payment | - | - | - |
| Payment Confirmed (via webhook) | +amount | - | +amount |
| Payment Failed | - | - | - |
| Payment Expired | - | - | - |

Note: Balance changes occur during the notification flow, not during payment creation.
