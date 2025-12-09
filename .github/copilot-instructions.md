# DOKU Payment Gateway Module - Copilot Instructions

## Project Overview

This is the **DOKU Payment Gateway Module**, a Go package that integrates with the DOKU Payment Gateway API. It provides payment processing, account management, settlement calculations, and bank account inquiry functionalities. This module is designed to be used by the `setter-service` backend.

---

## 📚 Documentation Reference

**IMPORTANT**: Before implementing any feature or making changes, always consult the business logic documentation in the `markdown/` folder:

| Document | Purpose |
|----------|---------|
| `markdown/00-system-architecture.md` | System overview, core entities, project structure, and configuration |
| `markdown/01-create-payment-flow.md` | Payment link generation, request/response structures, and integration examples |
| `markdown/02-payment-notification-flow.md` | Webhook handling, signature verification, and payment status updates |
| `markdown/03-settlement-calculation-flow.md` | Fee calculation rules by payment method, tax handling, and gross amount calculation |
| `markdown/04-account-management-flow.md` | Sub-account creation, balance queries, and account management |
| `markdown/05-bank-account-inquiry-flow.md` | Bank account verification for disbursements, SNAP API authentication |

---

## 🏗️ Project Structure

```
doku/
├── app/
│   ├── cmd/                          # CLI commands
│   ├── config/                       # Configuration management
│   ├── constants/                    # Payment method constants
│   ├── models/                       # Data models (DokuOrder, DokuCustomer, etc.)
│   ├── requests/                     # API request structures
│   ├── responses/                    # API response structures
│   ├── usecases/                     # Business logic (doku_usecase.go, doku_settlement_usecase.go)
│   └── utils/                        # HTTP client, logging, helpers
├── markdown/                         # Business logic documentation
├── go.mod
└── go.sum
```

---

## 🔑 Core Interfaces

### DokuUseCaseInterface
Main interface for DOKU operations:
- `CreateAccount()` - Create DOKU Sub-Account (SAC)
- `AcceptPayment()` - Generate payment link for customers
- `HandleNotification()` - Process webhooks from DOKU
- `GetBalance()` - Query account balance
- `GetToken()` - Get access token for SNAP API
- `BankAccountInquiry()` - Verify bank account for disbursement

### DokuSettlementUseCaseInterface
Settlement and fee calculation:
- `CalculateSettlementFee()` - Calculate net amount after fees (given gross amount)
- `CalculateGrossAmount()` - Calculate gross amount needed (given desired net amount)

---

## 💳 Supported Payment Methods

When working with payment methods, use the constants defined in `app/constants/doku_payment_method_constant.go`:

### Virtual Account
`VIRTUAL_ACCOUNT_BCA`, `VIRTUAL_ACCOUNT_BANK_MANDIRI`, `VIRTUAL_ACCOUNT_BANK_SYARIAH_MANDIRI`, `VIRTUAL_ACCOUNT_BRI`, `VIRTUAL_ACCOUNT_BNI`, `VIRTUAL_ACCOUNT_DOKU`, `VIRTUAL_ACCOUNT_BANK_PERMATA`, `VIRTUAL_ACCOUNT_BANK_CIMB`, `VIRTUAL_ACCOUNT_BANK_DANAMON`, `VIRTUAL_ACCOUNT_BTN`, `VIRTUAL_ACCOUNT_BNC`

### Credit/Debit Cards
`CREDIT_CARD`, `KARTU_KREDIT_INDONESIA`

### Convenience Store
`ONLINE_TO_OFFLINE_ALFA`, `ONLINE_TO_OFFLINE_INDOMARET`

### QRIS
`QRIS`

### E-Wallet
`EMONEY_OVO`, `EMONEY_SHOPEE_PAY`, `EMONEY_DOKU`, `EMONEY_LINKAJA`, `EMONEY_DANA`

### PayLater
`PEER_TO_PEER_AKULAKU`, `PEER_TO_PEER_KREDIVO`, `PEER_TO_PEER_INDODANA`

### Direct Debit & Digital Banking
`DIRECT_DEBIT_BRI`, `JENIUS_PAY`

---

## 💰 Fee Calculation Rules

Refer to `markdown/03-settlement-calculation-flow.md` for complete details. Quick reference:

| Payment Method | Fee | Tax |
|----------------|-----|-----|
| Virtual Account | IDR 4,000 flat | 11% PPN |
| Credit Card | 2.8% + IDR 2,000 | 11% PPN |
| QRIS | IDR 700 flat | **No Tax** |
| E-Wallet | 2% | 11% PPN |
| PayLater | 3% | 11% PPN |
| Convenience Store (Alfa) | IDR 5,000 | 11% PPN |
| Convenience Store (Indomaret) | IDR 6,500 | 11% PPN |

**Formula**: `Net Amount = Gross Amount - Transaction Fee - Tax`

---

## 🔐 Security Guidelines

### Signature Generation (for outgoing API calls)
- Use HMAC-SHA256 with the Secret Key
- Components: `Client-Id`, `Request-Id`, `Request-Timestamp`, `Request-Target`, `Digest`
- See `app/usecases/snap.go` for implementation

### Webhook Verification (for incoming notifications)
- Always verify signature before processing webhooks
- Match `original_request_id` with stored `gateway_request_id`
- Implement idempotency to handle duplicate notifications

---

## 📋 Coding Guidelines

### When Creating New Features

1. **Check documentation first**: Review relevant markdown files before implementing
2. **Follow existing patterns**: Look at existing usecases for structure
3. **Use proper error handling**: Return `*models.ErrorLog` for errors
4. **Include signature handling**: All DOKU API calls require proper signature

### Request/Response Handling

1. Request structures go in `app/requests/`
2. Response structures go in `app/responses/`
3. Use `null.String`, `null.Int` for nullable JSON fields
4. Always validate required fields before API calls

### Wallet/Ledger Integration

When integrating with setter-service ledger system:

| Action | pending_balance | balance | income_accumulation |
|--------|-----------------|---------|---------------------|
| Payment Created | - | - | - |
| Payment Confirmed | +amount | - | +amount |
| Settlement Processed | -amount | +net_amount | - |
| Withdrawal | - | -amount | - |

---

## 🧪 Testing Guidelines

1. Use DOKU Sandbox environment for testing
2. Sandbox URL: `https://api-sandbox.doku.com`
3. Production URL: `https://api.doku.com`
4. Test all payment flows with various payment methods
5. Verify webhook signature handling works correctly

---

## ⚠️ Common Pitfalls

1. **Invoice Number Uniqueness**: Always generate unique invoice numbers
2. **Payment Expiration**: Set appropriate `payment_due_date` (in minutes)
3. **Timezone**: Use UTC for all timestamps (`2006-01-02T15:04:05Z` format)
4. **Amount Precision**: Use `int64` for amounts in smallest currency unit
5. **QRIS Tax Exception**: QRIS has no tax, unlike other payment methods
6. **Duplicate Webhooks**: Implement idempotency checks for webhook handling

---

## 🔧 Environment Variables

Required configuration:
- `DOKU_API_CLIENT_ID` - DOKU Client ID
- `DOKU_API_SECRET_KEY` - DOKU Secret Key
- `DOKU_API_PARTNER_ID` - DOKU Partner ID (for SNAP API)
- `DOKU_API_PRIVATE_KEY` - Private Key (for SNAP API asymmetric signature)
- `DOKU_API_PUBLIC_KEY` - Public Key (for SNAP API)

---

## 📖 Quick Reference Commands

```go
// Create Sub-Account
dokuUseCase.CreateAccount(&requests.DokuCreateSubAccountRequest{...})

// Create Payment
dokuUseCase.AcceptPayment(&requests.DokuCreatePaymentRequest{...})

// Handle Webhook
dokuUseCase.HandleNotification(&requests.DokuNotificationRequest{...})

// Get Balance
dokuUseCase.GetBalance(sacID)

// Calculate Settlement Fee
settlementUseCase.CalculateSettlementFee(grossAmount, paymentMethod)

// Calculate Gross Amount (for upsert)
settlementUseCase.CalculateGrossAmount(desiredNetAmount, paymentMethod)

// Bank Account Inquiry
dokuUseCase.BankAccountInquiry(&requests.DokuBankAccountInquiryRequest{...})
```

---

## 🆘 Troubleshooting

When encountering issues:

1. Check the corresponding markdown documentation
2. Verify signature generation is correct
3. Ensure all required headers are present
4. Check API response error messages
5. Verify Sub-Account ID (SAC) is valid
6. Confirm environment variables are set correctly

For webhook issues:
1. Verify signature verification logic
2. Check if `original_request_id` matches
3. Ensure idempotency handling is working
4. Review payment status transitions