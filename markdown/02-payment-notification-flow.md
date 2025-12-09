# Payment Notification Flow - Business Logic Documentation

## Overview

The Payment Notification flow handles incoming webhooks from DOKU when a customer completes a payment. This flow is critical for updating payment status and triggering balance updates in the Ledger system.

---

## Notification Flow Diagram

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Customer   │     │     DOKU     │     │   Setter     │     │    Ledger    │
│   Pays       │     │   Gateway    │     │   Service    │     │    System    │
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                    │                    │                    │
       │ 1. Complete        │                    │                    │
       │    payment         │                    │                    │
       │ ──────────────────▶│                    │                    │
       │                    │                    │                    │
       │                    │ 2. Send webhook    │                    │
       │                    │    notification    │                    │
       │                    │ ──────────────────▶│                    │
       │                    │                    │                    │
       │                    │                    │ 3. Verify          │
       │                    │                    │    signature       │
       │                    │                    │                    │
       │                    │                    │ 4. Parse           │
       │                    │                    │    notification    │
       │                    │                    │                    │
       │                    │                    │ 5. Update          │
       │                    │                    │    LedgerPayment   │
       │                    │                    │ ──────────────────▶│
       │                    │                    │                    │
       │                    │ 6. Return 200 OK   │                    │
       │                    │ ◀──────────────────│                    │
       │                    │                    │                    │
```

---

## Request Structure

### DokuNotificationRequest

```go
type DokuNotificationRequest struct {
    RequestID        string `json:"Request-Id"`
    RequestTimestamp string `json:"Request-Timestamp"`
    Signature        string `json:"Signature"`
    RequestTarget    string `json:"Request-Target"`
    JsonBody         []byte `json:"Json-Body"`
}
```

| Field | Type | Description |
|-------|------|-------------|
| `RequestID` | string | DOKU request ID from headers |
| `RequestTimestamp` | string | Timestamp from headers |
| `Signature` | string | HMAC-SHA256 signature from headers |
| `RequestTarget` | string | Webhook endpoint path |
| `JsonBody` | []byte | Raw JSON body of notification |

---

## Response Structure

### DokuPostNotificationHTTPResponse

```go
type DokuPostNotificationHTTPResponse struct {
    Service struct {
        ID null.String `json:"id"`
    } `json:"service"`
    Acquirer struct {
        ID null.String `json:"id"`
    } `json:"acquirer"`
    Channel struct {
        ID null.String `json:"id"`
    } `json:"channel"`
    Transaction           *models.DokuTransaction           `json:"transaction"`
    Order                 *models.DokuOrder                 `json:"order"`
    Customer              *models.DokuCustomer              `json:"customer"`
    VirtualAccountInfo    *VirtualAccountInfo               `json:"virtual_account_info,omitempty"`
    VirtualAccountPayment *models.DokuVirtualAccountpayment `json:"virtual_account_payment,omitempty"`
    CardPayment           *models.DokuCardPayment           `json:"card_payment,omitempty"`
    // ... other payment method specific fields
}
```

### Key Response Fields

| Field | Description |
|-------|-------------|
| `service.id` | Service type (e.g., "VIRTUAL_ACCOUNT") |
| `acquirer.id` | Bank/provider (e.g., "BCA") |
| `channel.id` | Payment channel (e.g., "VIRTUAL_ACCOUNT_BCA") |
| `transaction.status` | "SUCCESS", "FAILED", or "PENDING" |
| `transaction.original_request_id` | Matches `gateway_request_id` from CreatePayment |
| `order.invoice_number` | Your invoice number |
| `order.amount` | Payment amount |

---

## Detailed Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        PAYMENT NOTIFICATION FLOW                                 │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  1. DOKU sends POST webhook to your notification URL                           │
│                                                                                 │
│     Headers:                                                                    │
│     - Client-Id: {your_client_id}                                              │
│     - Request-Id: {unique_request_id}                                          │
│     - Request-Timestamp: {ISO8601_timestamp}                                   │
│     - Signature: HMACSHA256={signature_hash}                                   │
│                                                                                 │
│  2. Setter-service extracts headers and body                                   │
│                                                                                 │
│  3. Setter-service calls DOKU module:                                          │
│     dokuUseCase.HandleNotification(&DokuNotificationRequest{                   │
│         RequestID:        headers["Request-Id"],                               │
│         RequestTimestamp: headers["Request-Timestamp"],                        │
│         Signature:        headers["Signature"],                                │
│         RequestTarget:    "/your/webhook/endpoint",                            │
│         JsonBody:         rawBody,                                             │
│     })                                                                          │
│                                                                                 │
│  4. DOKU module verifies signature                                             │
│     - Reconstructs signature from components                                   │
│     - Compares with received signature                                         │
│     - Rejects if signature doesn't match (401 Unauthorized)                    │
│                                                                                 │
│  5. DOKU module parses JSON body into structured response                      │
│                                                                                 │
│  6. Setter-service processes based on transaction.status:                      │
│                                                                                 │
│     IF status == "SUCCESS":                                                    │
│       - Find LedgerPayment by original_request_id                              │
│       - Update status to PAID                                                  │
│       - Add amount to wallet.pending_balance                                   │
│       - Create LedgerTransaction (type: PAYMENT)                               │
│                                                                                 │
│     IF status == "FAILED":                                                     │
│       - Find LedgerPayment by original_request_id                              │
│       - Update status to FAILED                                                │
│       - Log failure reason                                                     │
│                                                                                 │
│  7. Return HTTP 200 OK to DOKU                                                 │
│     (Important: Always return 200 after processing)                            │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Implementation Logic

### HandleNotification Method

```go
func (u *dokuUseCase) HandleNotification(
    request *requests.DokuNotificationRequest,
) (*responses.DokuPostNotificationHTTPResponse, *models.ErrorLog) {

    // 1. Verify signature components
    isValid, logData := u.verifySignatureComponent(
        request.Signature,
        request.RequestID,
        request.RequestTimestamp,
        request.RequestTarget,
        request.JsonBody,
    )

    if logData != nil {
        return nil, logData
    }

    if !isValid {
        errorMessage := "Invalid signature in notification"
        return nil, helper.WriteLog(
            errors.New(errorMessage), 
            http.StatusUnauthorized, 
            errorMessage,
        )
    }

    // 2. Parse notification body
    notificationResponse := &responses.DokuPostNotificationHTTPResponse{}
    err := json.Unmarshal(request.JsonBody, &notificationResponse)
    if err != nil {
        return nil, helper.WriteLog(
            err, 
            http.StatusInternalServerError, 
            "Failed to unmarshal notification body",
        )
    }

    return notificationResponse, nil
}
```

### Signature Verification

```go
func (u *dokuUseCase) verifySignatureComponent(
    signature string,
    requestId string,
    requestTimestamp string,
    requestTarget string,
    jsonBody []byte,
) (bool, *models.ErrorLog) {

    // 1. Calculate Digest from body
    var digest string
    if jsonBody != nil && len(jsonBody) > 0 {
        hash := sha256.Sum256(jsonBody)
        digest = base64.StdEncoding.EncodeToString(hash[:])
    }

    // 2. Build signature components string
    signatureComponents := fmt.Sprintf(
        "Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s\nDigest:%s",
        u.DokuAPIClientID,
        requestId,
        requestTimestamp,
        requestTarget,
        digest,
    )

    // 3. Calculate expected signature
    h := hmac.New(sha256.New, []byte(u.DokuAPISecretKey))
    h.Write([]byte(signatureComponents))
    signatureHash := base64.StdEncoding.EncodeToString(h.Sum(nil))
    expectedSignature := fmt.Sprintf("HMACSHA256=%s", signatureHash)

    // 4. Compare signatures
    if signature != expectedSignature {
        return false, helper.WriteLog(
            errors.New("Signature does not match"),
            http.StatusUnauthorized,
            "Signature does not match",
        )
    }

    return true, nil
}
```

---

## Setter-Service Integration Example

```go
// Webhook handler in setter-service
func (h *webhookHandler) HandleDokuNotification(c *gin.Context) {
    
    // 1. Extract headers
    requestID := c.GetHeader("Request-Id")
    requestTimestamp := c.GetHeader("Request-Timestamp")
    signature := c.GetHeader("Signature")
    
    // 2. Read raw body
    rawBody, err := io.ReadAll(c.Request.Body)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
        return
    }
    
    // 3. Call DOKU module to handle notification
    notification, errLog := h.dokuUseCase.HandleNotification(&requests.DokuNotificationRequest{
        RequestID:        requestID,
        RequestTimestamp: requestTimestamp,
        Signature:        signature,
        RequestTarget:    "/api/v1/webhooks/doku/payment",
        JsonBody:         rawBody,
    })
    
    if errLog != nil {
        c.JSON(errLog.StatusCode, gin.H{"error": errLog.Message})
        return
    }
    
    // 4. Process based on transaction status
    if notification.Transaction.Status.String == "SUCCESS" {
        err = h.processSuccessfulPayment(c, notification)
    } else if notification.Transaction.Status.String == "FAILED" {
        err = h.processFailedPayment(c, notification)
    }
    
    if err != nil {
        // Log error but still return 200 to DOKU
        log.Printf("Error processing notification: %v", err)
    }
    
    // 5. Always return 200 OK to DOKU
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *webhookHandler) processSuccessfulPayment(
    c *gin.Context, 
    notification *responses.DokuPostNotificationHTTPResponse,
) error {
    
    // 1. Start transaction
    tx, err := h.db.BeginTx(c, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 2. Get payment method from notification
    paymentMethod := notification.Channel.ID.String
    
    // 3. Confirm payment in Ledger
    confirmedPayment, err := h.ledgerPaymentUseCase.ConfirmPayment(tx, &ledger.LedgerPaymentConfirmPaymentRequest{
        GatewayRequestId:       notification.Transaction.OriginalRequestID.String,
        PaymentMethod:          paymentMethod,
        PaymentDate:            time.Now(), // Or parse from notification
        GatewayReferenceNumber: getGatewayReference(notification),
    })
    if err != nil {
        return err
    }
    
    // 4. Calculate settlement fee using actual payment method
    // This gives us accurate net amount after DOKU fees
    settlementResult, err := h.dokuSettlementUseCase.CalculateSettlementFee(
        paymentMethod, 
        float64(confirmedPayment.Amount),
    )
    if err != nil {
        return err
    }
    
    // 5. Create settlement record with IN_PROGRESS status
    // Use invoice number as batch number for idempotency
    // DOKU settles daily on weekdays at ~1 PM
    estimatedSettlementDate := time.Now().AddDate(0, 0, 1) // Next day estimate
    
    _, err = h.ledgerSettlementUseCase.CreateSettlement(
        tx,
        confirmedPayment.LedgerAccountUUID,
        confirmedPayment.InvoiceNumber, // batch_number for idempotency
        estimatedSettlementDate,
        confirmedPayment.Currency,
        int64(settlementResult.GrossAmount),
        int64(settlementResult.NetAmount),
        "", // bankName - filled during disbursement
        "", // bankAccountNumber - filled during disbursement
        ledger_models.AccountTypeSubAccount,
    )
    if err != nil {
        return err
    }
    
    // 6. Commit transaction
    return tx.Commit()
}
```

---

## Webhook Payload Examples

### Virtual Account (BCA) - Success

```json
{
  "service": {
    "id": "VIRTUAL_ACCOUNT"
  },
  "acquirer": {
    "id": "BCA"
  },
  "channel": {
    "id": "VIRTUAL_ACCOUNT_BCA"
  },
  "order": {
    "invoice_number": "INV-USER001-1736939400",
    "amount": 100000
  },
  "virtual_account_info": {
    "virtual_account_number": "1900800000208690"
  },
  "virtual_account_payment": {
    "date": "20251204224523",
    "systrace_number": "116245",
    "reference_number": "00933",
    "identifier": [
      {
        "name": "REQUEST_ID",
        "value": "855183"
      },
      {
        "name": "REFERENCE",
        "value": "00933"
      }
    ]
  },
  "transaction": {
    "status": "SUCCESS",
    "date": "2025-12-04T15:45:23Z",
    "original_request_id": "550e8400-e29b-41d4-a716-446655440000"
  },
  "additional_info": {
    "origin": {
      "source": "direct",
      "system": "mid-jokul-checkout-system",
      "product": "CHECKOUT",
      "apiFormat": "JOKUL"
    },
    "account": {
      "id": "SAC-7327-1764507463535"
    }
  }
}
```

### E-Wallet (ShopeePay) - Success

```json
{
  "service": {
    "id": "EMONEY"
  },
  "acquirer": {
    "id": "SHOPEEPAY"
  },
  "channel": {
    "id": "EMONEY_SHOPEE_PAY"
  },
  "order": {
    "invoice_number": "INV-USER001-1736939500",
    "amount": 50000
  },
  "shopeepay_configuration": {
    "merchant_ext_id": "MID12345",
    "store_ext_id": "STORE001"
  },
  "shopeepay_payment": {
    "transaction_status": "00",
    "transaction_message": "Success",
    "identifier": [
      {
        "name": "SHOPEEPAY_REF_ID",
        "value": "SPY123456789"
      }
    ]
  },
  "transaction": {
    "status": "SUCCESS",
    "original_request_id": "661e9511-b30c-42d8-b789-123456789abc"
  }
}
```

### QRIS - Success

```json
{
  "service": {
    "id": "QRIS"
  },
  "acquirer": {
    "id": "QRIS"
  },
  "channel": {
    "id": "QRIS"
  },
  "order": {
    "invoice_number": "INV-USER001-1736939600",
    "amount": 75000
  },
  "emoney_payment": {
    "account_id": "ACC123456",
    "approval_code": "APR789"
  },
  "settlement": [
    {
      "bank_account_settlement_id": "SETTLE001",
      "value": 74300,
      "type": "NETT"
    }
  ],
  "transaction": {
    "status": "SUCCESS",
    "original_request_id": "772f0622-c41d-53e9-c89a-234567890def"
  }
}
```

### Credit Card - Success

```json
{
  "service": {
    "id": "CREDIT_CARD"
  },
  "acquirer": {
    "id": "BCA"
  },
  "channel": {
    "id": "CREDIT_CARD"
  },
  "order": {
    "invoice_number": "INV-USER001-1736939700",
    "amount": 200000
  },
  "card_payment": {
    "masked_card_number": "4111XXXXXXXX1111",
    "approval_code": "ABC123",
    "response_code": "00",
    "response_message": "Approved",
    "issuer": "BCA",
    "payment_id": "PAY12345"
  },
  "authorized_id": "AUTH98765",
  "transaction": {
    "status": "SUCCESS",
    "original_request_id": "883g1733-d52e-64fa-d9ab-345678901efg"
  }
}
```

---

## Extracting Payment Reference

```go
func getGatewayReference(notification *responses.DokuPostNotificationHTTPResponse) string {
    // Virtual Account
    if notification.VirtualAccountPayment != nil {
        return notification.VirtualAccountPayment.ReferenceNumber.String
    }
    
    // Credit Card
    if notification.CardPayment != nil {
        return notification.CardPayment.PaymentID.String
    }
    
    // ShopeePay
    if notification.ShopeepayPayment != nil {
        for _, id := range notification.ShopeepayPayment.Identifier {
            if id.Name.String == "SHOPEEPAY_REF_ID" {
                return id.Value.String
            }
        }
    }
    
    // OVO
    if notification.Wallet.TokenID.Valid {
        return notification.Wallet.TokenID.String
    }
    
    // QRIS
    if notification.EmoneyPayment != nil {
        return notification.EmoneyPayment.ApprovalCode.String
    }
    
    return ""
}
```

---

## Error Handling

### Signature Verification Failed

```json
{
  "status_code": 401,
  "message": "Invalid signature in notification"
}
```

**Action**: Do not process the notification. Log the event for security review.

### Payment Not Found

```json
{
  "status_code": 404,
  "message": "Payment not found for request_id: xxx"
}
```

**Action**: Log the event. This may indicate a race condition or duplicate notification.

### Payment Already Processed

```json
{
  "status_code": 409,
  "message": "Payment already in PAID status"
}
```

**Action**: Return 200 OK to DOKU (idempotent). Log as duplicate notification.

---

## Idempotency Handling

DOKU may send the same notification multiple times. Implement idempotent handling:

```go
func (h *webhookHandler) processSuccessfulPayment(
    notification *responses.DokuPostNotificationHTTPResponse,
) error {
    
    // 1. Find existing payment
    payment, err := h.ledgerPaymentRepo.GetByGatewayRequestId(
        notification.Transaction.OriginalRequestID.String,
    )
    if err != nil {
        return err
    }
    
    // 2. Check if already processed (idempotent)
    if payment.Status == models.PaymentStatusPaid {
        log.Printf("Payment already processed: %s", payment.UUID)
        return nil // Success - already processed
    }
    
    // 3. Validate current status
    if payment.Status != models.PaymentStatusPending {
        return fmt.Errorf("invalid payment status: %s", payment.Status)
    }
    
    // 4. Process payment...
    return h.updatePaymentToPaid(payment, notification)
}
```

---

## Wallet & Settlement Impact Summary

| Notification Status | pending_balance | balance | income_accumulation | Settlement Created |
|---------------------|-----------------|---------|---------------------|-------------------|
| SUCCESS | +amount | - | +amount | Yes (IN_PROGRESS) |
| FAILED | - | - | - | No |

### Why Create Settlement on Payment Confirmation?

**Important**: Settlements should be created when payment is confirmed (webhook SUCCESS), NOT when the payment link is created.

**Correct Flow:**
```
Payment Link Created → Customer pays → Webhook SUCCESS → ConfirmPayment + CreateSettlement
```

**Why NOT at Payment Link Creation?**
- Customer may abandon payment → orphaned settlement records
- Customer may choose different payment method → wrong fee calculation
- Payment may expire → cleanup needed for unused settlements

**Benefits of Webhook-based Settlement Creation:**
1. **Accurate Fee Calculation**: Uses actual payment method from DOKU (e.g., customer chose QRIS instead of VA)
2. **No Orphaned Records**: Settlement only exists when money is actually received
3. **Idempotency**: Uses invoice number as batch_number to prevent duplicates on webhook retries
4. **Clean Ledger**: Only tracks money that actually moves

---

## Best Practices

1. **Always Return 200**: Return HTTP 200 OK to DOKU after processing (even if internal errors occur) to prevent retry floods.

2. **Verify Signature First**: Always verify the HMAC signature before processing any notification data.

3. **Idempotent Processing**: Handle duplicate notifications gracefully. Check payment status before processing.

4. **Log Everything**: Log all incoming notifications with full payload for debugging and audit purposes.

5. **Use Database Transactions**: Wrap payment confirmation, settlement creation, and wallet balance update in a single transaction.

6. **Async Processing**: For high-volume systems, consider queueing notifications for async processing while returning 200 immediately.

7. **Create Settlement on Confirmation**: Always create the settlement record in the webhook handler after payment confirmation, not at payment link creation. This ensures:
   - Accurate fee calculation based on actual payment method used
   - No orphaned settlement records for abandoned payments
   - Clean ledger state

---

## Webhook Configuration

Register your webhook URL in DOKU Dashboard:

| Environment | URL Pattern |
|-------------|-------------|
| Sandbox | `https://your-sandbox-domain.com/api/v1/webhooks/doku/payment` |
| Production | `https://your-production-domain.com/api/v1/webhooks/doku/payment` |

Ensure your webhook endpoint:
- Accepts POST requests
- Has valid SSL certificate
- Is publicly accessible
- Responds within 30 seconds
