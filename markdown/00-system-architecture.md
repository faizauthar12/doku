# DOKU Payment Gateway Module - System Architecture Documentation

## Overview

The DOKU module is a Go package designed to integrate with the DOKU Payment Gateway API. It provides a clean interface for payment processing, account management, and settlement calculations. This module is intended to be used by the `setter-service` (backend service) to handle all DOKU-related operations.

## System Purpose

- Create and manage DOKU Sub-Accounts (SAC) for merchants
- Generate payment links for customers (Accept Payment)
- Handle payment notifications/webhooks from DOKU
- Calculate settlement fees based on payment methods
- Query account balances from DOKU
- Perform bank account inquiries for disbursements ("KIRIM DOKU")

---

## Core Entities

### 1. DokuOrder
Represents an order/transaction in the payment flow.

| Field | Type | Description |
|-------|------|-------------|
| `invoice_number` | string | Unique invoice/order ID |
| `amount` | int64 | Payment amount |
| `currency` | string | Currency code (e.g., "IDR") |
| `session_id` | string | Session ID from DOKU |

### 2. DokuCustomer
Represents customer information for payment.

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Customer ID |
| `name` | string | Customer name |
| `email` | string | Customer email |

### 3. DokuPayment
Represents payment configuration and response data.

| Field | Type | Description |
|-------|------|-------------|
| `payment_method_types` | []string | Allowed payment methods |
| `payment_due_date` | int64 | Payment expiration (in minutes) |
| `token_id` | string | DOKU session token |
| `url` | string | Checkout URL for customer |
| `expired_date` | string | Payment expiration datetime |

### 4. DokuTransaction
Represents transaction status from DOKU webhook.

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Transaction type |
| `status` | string | SUCCESS, FAILED, PENDING |
| `original_request_id` | string | Original request ID for matching |

### 5. DokuBalance
Represents account balance from DOKU.

| Field | Type | Description |
|-------|------|-------------|
| `pending` | string | Pending balance (waiting for settlement) |
| `available` | string | Available balance |

### 6. DokuSettlement
Represents settlement information from DOKU.

| Field | Type | Description |
|-------|------|-------------|
| `bank_account_settlement_id` | string | Settlement bank account ID |
| `value` | float64 | Settlement amount |
| `type` | string | Settlement type |

---

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              DOKU MODULE                                         │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                         DokuUseCaseInterface                             │   │
│  ├─────────────────────────────────────────────────────────────────────────┤   │
│  │  • CreateAccount()        - Create DOKU Sub-Account                     │   │
│  │  • AcceptPayment()        - Generate payment link                       │   │
│  │  • HandleNotification()   - Process DOKU webhooks                       │   │
│  │  • GetBalance()           - Query account balance                       │   │
│  │  • GetToken()             - Get access token for SNAP API               │   │
│  │  • BankAccountInquiry()   - Verify bank account for disbursement        │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                    DokuSettlementUseCaseInterface                        │   │
│  ├─────────────────────────────────────────────────────────────────────────┤   │
│  │  • CalculateSettlementFee()  - Calculate net amount after fees          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Integration with Setter-Service

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         SETTER-SERVICE INTEGRATION                               │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌──────────────┐                ┌──────────────┐          ┌──────────────┐    │
│  │   Frontend   │                │   Setter     │          │    DOKU      │    │
│  │   (Client)   │                │   Service    │          │   Module     │    │
│  └──────┬───────┘                └──────┬───────┘          └──────┬───────┘    │
│         │                               │                         │            │
│         │  1. Create Payment Request    │                         │            │
│         │ ─────────────────────────────▶│                         │            │
│         │                               │                         │            │
│         │                               │  2. AcceptPayment()     │            │
│         │                               │ ────────────────────────▶            │
│         │                               │                         │            │
│         │                               │                    3. Call DOKU API  │
│         │                               │                         │ ─────────▶ │
│         │                               │                         │            │
│         │                               │  4. Return payment URL  │            │
│         │                               │ ◀────────────────────────            │
│         │                               │                         │            │
│         │  5. Redirect to payment URL   │                         │            │
│         │ ◀─────────────────────────────│                         │            │
│         │                               │                         │            │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Payment Methods Supported

### Virtual Account (Transfer Bank)
| Constant | Description |
|----------|-------------|
| `VIRTUAL_ACCOUNT_BCA` | Virtual Account BCA |
| `VIRTUAL_ACCOUNT_BANK_MANDIRI` | Virtual Account Mandiri |
| `VIRTUAL_ACCOUNT_BANK_SYARIAH_MANDIRI` | Virtual Account BSI |
| `VIRTUAL_ACCOUNT_BRI` | Virtual Account BRI |
| `VIRTUAL_ACCOUNT_BNI` | Virtual Account BNI |
| `VIRTUAL_ACCOUNT_DOKU` | Virtual Account DOKU |
| `VIRTUAL_ACCOUNT_BANK_PERMATA` | Virtual Account Permata |
| `VIRTUAL_ACCOUNT_BANK_CIMB` | Virtual Account CIMB |
| `VIRTUAL_ACCOUNT_BANK_DANAMON` | Virtual Account Danamon |
| `VIRTUAL_ACCOUNT_BTN` | Virtual Account BTN |
| `VIRTUAL_ACCOUNT_BNC` | Virtual Account BNC |

### Credit/Debit Cards
| Constant | Description |
|----------|-------------|
| `CREDIT_CARD` | Credit Card (Visa, Mastercard, JCB, Amex) |
| `KARTU_KREDIT_INDONESIA` | Kartu Kredit Indonesia (KKI) |

### Convenience Store
| Constant | Description |
|----------|-------------|
| `ONLINE_TO_OFFLINE_ALFA` | Alfamart/Alfa Group |
| `ONLINE_TO_OFFLINE_INDOMARET` | Indomaret |

### QRIS
| Constant | Description |
|----------|-------------|
| `QRIS` | QR Code Indonesia Standard |

### E-Wallet
| Constant | Description |
|----------|-------------|
| `EMONEY_OVO` | OVO |
| `EMONEY_SHOPEE_PAY` | ShopeePay |
| `EMONEY_DOKU` | DOKU Wallet |
| `EMONEY_LINKAJA` | LinkAja |
| `EMONEY_DANA` | DANA |

### PayLater
| Constant | Description |
|----------|-------------|
| `PEER_TO_PEER_AKULAKU` | Akulaku PayLater |
| `PEER_TO_PEER_KREDIVO` | Kredivo PayLater |
| `PEER_TO_PEER_INDODANA` | Indodana PayLater |

### Direct Debit
| Constant | Description |
|----------|-------------|
| `DIRECT_DEBIT_BRI` | Direct Debit BRI |

### Digital Banking
| Constant | Description |
|----------|-------------|
| `JENIUS_PAY` | Jenius Pay |

---

## Project Structure

```
doku/
├── app/
│   ├── cmd/                          # CLI commands (if any)
│   │
│   ├── config/
│   │   └── config.go                 # Configuration management
│   │
│   ├── constants/
│   │   └── doku_payment_method_constant.go  # Payment method constants
│   │
│   ├── models/
│   │   ├── base_model.go             # Base model with error handling
│   │   └── doku_model.go             # DOKU-specific models
│   │
│   ├── requests/
│   │   └── doku_request.go           # Request structures
│   │
│   ├── responses/
│   │   ├── doku_response.go          # Response structures
│   │   └── doku_settlement_response.go  # Settlement response
│   │
│   ├── usecases/
│   │   ├── doku_usecase.go           # Main DOKU use case
│   │   ├── doku_settlement_usecase.go # Settlement fee calculation
│   │   └── snap.go                   # SNAP API signature helpers
│   │
│   └── utils/
│       └── helper/                   # HTTP client, logging, etc.
│
├── markdown/
│   ├── 00-system-architecture.md     # This file
│   ├── 01-create-payment-flow.md     # Payment creation flow
│   ├── 02-payment-notification-flow.md # Webhook handling
│   ├── 03-settlement-calculation-flow.md # Fee calculation
│   ├── 04-account-management-flow.md # Account & balance
│   └── 05-bank-account-inquiry-flow.md # Disbursement prep
│
├── go.mod
├── go.sum
└── .gitignore
```

---

## Configuration

The module requires the following environment variables:

### DOKU API Credentials
| Variable | Description |
|----------|-------------|
| `DOKU_API_CLIENT_ID` | DOKU Client ID |
| `DOKU_API_SECRET_KEY` | DOKU Secret Key |
| `DOKU_API_PRIVATE_KEY` | DOKU Private Key (for SNAP API) |

### Transaction Fee Configuration
| Variable | Default | Description |
|----------|---------|-------------|
| `TRANSACTION_FEE_CARDS_PERCENTAGE_RATE` | 2.8 | Card fee percentage |
| `TRANSACTION_FEE_CARDS_FLAT_FEE` | 2000 | Card flat fee (IDR) |
| `TRANSACTION_FEE_VIRTUAL_ACCOUNT_FLAT_FEE` | 4000 | VA flat fee (IDR) |
| `TRANSACTION_FEE_ALFAMART_FLAT_FEE` | 5000 | Alfamart fee (IDR) |
| `TRANSACTION_FEE_INDOMARET_FLAT_FEE` | 6500 | Indomaret fee (IDR) |
| `TRANSACTION_FEE_QR_FLAT_FEE` | 700 | QRIS fee (IDR) |
| `TRANSACTION_FEE_SHOPEEPAY_PERCENTAGE_RATE` | 2.0 | ShopeePay percentage |
| `TRANSACTION_FEE_OVO_PERCENTAGE_RATE` | 2.0 | OVO percentage |
| `TRANSACTION_FEE_LINKAJA_PERCENTAGE_RATE` | 2.0 | LinkAja percentage |
| `TRANSACTION_FEE_DOKU_WALLET_PERCENTAGE_RATE` | 1.5 | DOKU Wallet percentage |
| `TRANSACTION_FEE_DANA_PERCENTAGE_RATE` | 1.5 | DANA percentage |
| `TRANSACTION_FEE_BRI_DIRECT_DEBIT_PERCENTAGE_RATE` | 2.0 | Direct Debit BRI |
| `TRANSACTION_FEE_JENIUSPAY_PERCENTAGE_RATE` | 1.5 | Jenius Pay percentage |
| `TRANSACTION_FEE_AKULAKU_PERCENTAGE_RATE` | 1.5 | Akulaku percentage |
| `TRANSACTION_FEE_KREDIVO_PERCENTAGE_RATE` | 2.3 | Kredivo percentage |
| `TRANSACTION_FEE_INDODANA_PERCENTAGE_RATE` | 2.3 | Indodana percentage |
| `TRANSACTION_FEE_TAX` | 11 | Tax percentage (PPN) |

---

## Initialization

### As a Standalone Module
```go
import "github.com/faizauthar12/doku/app/config"

func main() {
    // Load .env file and initialize config
    config.InitConfig(".env")
    
    // Create use case
    cfg := config.Get()
    dokuUseCase := usecases.NewDokuUseCase(
        cfg.Doku.ClientID,
        cfg.Doku.SecretKey,
        cfg.Doku.PrivateKey,
    )
}
```

### As a Module in Parent Project
```go
import "github.com/faizauthar12/doku/app/config"

func main() {
    // Parent project already loaded .env
    godotenv.Load(".env")
    
    // Initialize config from existing env vars
    config.InitConfigFromEnv()
    
    // Or set config directly
    config.InitConfigWithStruct(config.Configuration{
        Doku: struct{...}{
            ClientID:   "your-client-id",
            SecretKey:  "your-secret-key",
            PrivateKey: "your-private-key",
        },
        // ... other config
    })
}
```

---

## Complete Money Flow with Ledger System

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         COMPLETE INTEGRATION FLOW                                │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌─────────────┐                                                                │
│  │  Customer   │                                                                │
│  │   Pays      │                                                                │
│  └──────┬──────┘                                                                │
│         │                                                                       │
│         ▼                                                                       │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │ STEP 1: CREATE PAYMENT (DOKU Module)                                    │   │
│  │                                                                          │   │
│  │   dokuUseCase.AcceptPayment()                                           │   │
│  │   → Returns payment_url for customer checkout                           │   │
│  │   → Setter-service creates LedgerPayment (status: PENDING)              │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│         │                                                                       │
│         │ Customer completes payment                                            │
│         ▼                                                                       │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │ STEP 2: PAYMENT NOTIFICATION (DOKU Module)                              │   │
│  │                                                                          │   │
│  │   dokuUseCase.HandleNotification()                                      │   │
│  │   → Validates webhook signature                                         │   │
│  │   → Returns transaction details                                         │   │
│  │   → Setter-service updates LedgerPayment (status: PAID)                 │   │
│  │   → Ledger adds to pending_balance                                      │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│         │                                                                       │
│         │ DOKU settles funds (T+1 or T+2)                                       │
│         ▼                                                                       │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │ STEP 3: SETTLEMENT CALCULATION (DOKU Module)                            │   │
│  │                                                                          │   │
│  │   dokuSettlementUseCase.CalculateSettlementFee()                        │   │
│  │   → Calculates transaction fee based on payment method                  │   │
│  │   → Calculates tax (11% of fee, except QRIS)                           │   │
│  │   → Returns net amount                                                  │   │
│  │   → Ledger creates Settlement (moves pending to available)              │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│         │                                                                       │
│         │ User initiates withdrawal                                             │
│         ▼                                                                       │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │ STEP 4: DISBURSEMENT "KIRIM DOKU" (DOKU Module)                         │   │
│  │                                                                          │   │
│  │   dokuUseCase.BankAccountInquiry()                                      │   │
│  │   → Validates destination bank account                                  │   │
│  │   → Setter-service initiates disbursement via DOKU SNAP API             │   │
│  │   → Money arrives in user's bank account                                │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## API Endpoints Summary (for setter-service)

| Use Case Method | DOKU API Endpoint | Description |
|-----------------|-------------------|-------------|
| `CreateAccount()` | `POST /sac-merchant/v1/accounts` | Create Sub-Account |
| `AcceptPayment()` | `POST /checkout/v1/payment` | Create payment link |
| `GetBalance()` | `GET /sac-merchant/v1/balances/{sacID}` | Get account balance |
| `HandleNotification()` | Webhook handler | Process payment webhook |
| `GetToken()` | `POST /authorization/v1/access-token/b2b` | Get SNAP access token |
| `BankAccountInquiry()` | `POST /snap/v1.1/emoney/bank-account-inquiry` | Verify bank account |

---

## Security

### Signature Generation
All API requests to DOKU require HMAC-SHA256 signature:

```
Signature Components (POST/PUT):
Client-Id:{client_id}
Request-Id:{uuid}
Request-Timestamp:{ISO8601}
Request-Target:{endpoint}
Digest:{SHA256(body)}

Signature Components (GET):
Client-Id:{client_id}
Request-Id:{uuid}
Request-Timestamp:{ISO8601}
Request-Target:{endpoint}
```

### Webhook Verification
Incoming webhooks from DOKU are verified using the same signature mechanism to ensure authenticity.
