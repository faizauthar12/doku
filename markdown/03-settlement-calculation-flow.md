# Settlement Calculation Flow - Business Logic Documentation

## Overview

The Settlement Calculation flow handles the calculation of transaction fees and net amounts for different payment methods. When DOKU settles funds to your account, they deduct fees based on the payment method used by the customer.

This module provides two key functions:
1. **CalculateSettlementFee**: Given a gross amount, calculate the net amount after fees
2. **CalculateGrossAmount (Upsert)**: Given a desired net amount, calculate what the customer must pay so the merchant receives exactly that amount

---

## Fee Structure Overview

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         SETTLEMENT CALCULATION FLOW                              │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│   Gross Amount (Customer Pays)                                                  │
│         │                                                                       │
│         ▼                                                                       │
│   ┌─────────────────────────────────────────────────┐                          │
│   │  Transaction Fee (varies by payment method)     │                          │
│   │  - Flat fee (e.g., IDR 4,000 for VA)           │                          │
│   │  - Percentage fee (e.g., 2% for E-Wallet)       │                          │
│   │  - Combination (e.g., 2.8% + IDR 2,000 for CC)  │                          │
│   └─────────────────────────────────────────────────┘                          │
│         │                                                                       │
│         ▼                                                                       │
│   ┌─────────────────────────────────────────────────┐                          │
│   │  Tax (PPN 11% of Transaction Fee)              │                          │
│   │  Note: QRIS has no tax                          │                          │
│   └─────────────────────────────────────────────────┘                          │
│         │                                                                       │
│         ▼                                                                       │
│   Total Deduction = Transaction Fee + Tax                                       │
│   Net Amount = Gross Amount - Total Deduction                                   │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Response Structure

### DokuSettlementResultResponse

```go
type DokuSettlementResultResponse struct {
    PaymentMethod  string  `json:"payment_method"`
    GrossAmount    float64 `json:"gross_amount"`
    TransactionFee float64 `json:"transaction_fee"`
    Tax            float64 `json:"tax"`
    TotalDeduction float64 `json:"total_deduction"`
    NetAmount      float64 `json:"net_amount"`
}
```

| Field | Type | Description |
|-------|------|-------------|
| `PaymentMethod` | string | Payment method constant |
| `GrossAmount` | float64 | Original amount (customer paid) |
| `TransactionFee` | float64 | Fee charged by DOKU |
| `Tax` | float64 | PPN tax (11% of fee) |
| `TotalDeduction` | float64 | Total fee + tax |
| `NetAmount` | float64 | Amount received after deductions |

---

## Fee Calculation Rules by Payment Method

### Credit/Debit Cards
| Payment Method | Fee Rate | Tax | Formula |
|----------------|----------|-----|---------|
| `CREDIT_CARD` | 2.8% + IDR 2,000 | 11% | fee = (amount × 0.028) + 2000; tax = fee × 0.11 |
| `KARTU_KREDIT_INDONESIA` | 2.8% + IDR 2,000 | 11% | fee = (amount × 0.028) + 2000; tax = fee × 0.11 |

**Example (IDR 100,000):**
```
Transaction Fee = (100,000 × 2.8%) + 2,000 = 2,800 + 2,000 = 4,800
Tax = 4,800 × 11% = 528
Total Deduction = 4,800 + 528 = 5,328
Net Amount = 100,000 - 5,328 = 94,672
```

---

### Virtual Account (Transfer Bank)
| Payment Method | Fee Rate | Tax | Formula |
|----------------|----------|-----|---------|
| `VIRTUAL_ACCOUNT_BCA` | IDR 4,000 | 11% | fee = 4000; tax = fee × 0.11 |
| `VIRTUAL_ACCOUNT_BANK_MANDIRI` | IDR 4,000 | 11% | fee = 4000; tax = fee × 0.11 |
| `VIRTUAL_ACCOUNT_BANK_SYARIAH_MANDIRI` | IDR 4,000 | 11% | fee = 4000; tax = fee × 0.11 |
| `VIRTUAL_ACCOUNT_BRI` | IDR 4,000 | 11% | fee = 4000; tax = fee × 0.11 |
| `VIRTUAL_ACCOUNT_BNI` | IDR 4,000 | 11% | fee = 4000; tax = fee × 0.11 |
| `VIRTUAL_ACCOUNT_DOKU` | IDR 4,000 | 11% | fee = 4000; tax = fee × 0.11 |
| `VIRTUAL_ACCOUNT_BANK_PERMATA` | IDR 4,000 | 11% | fee = 4000; tax = fee × 0.11 |
| `VIRTUAL_ACCOUNT_BANK_CIMB` | IDR 4,000 | 11% | fee = 4000; tax = fee × 0.11 |
| `VIRTUAL_ACCOUNT_BANK_DANAMON` | IDR 4,000 | 11% | fee = 4000; tax = fee × 0.11 |
| `VIRTUAL_ACCOUNT_BTN` | IDR 4,000 | 11% | fee = 4000; tax = fee × 0.11 |
| `VIRTUAL_ACCOUNT_BNC` | IDR 4,000 | 11% | fee = 4000; tax = fee × 0.11 |

**Example (IDR 100,000):**
```
Transaction Fee = 4,000
Tax = 4,000 × 11% = 440
Total Deduction = 4,000 + 440 = 4,440
Net Amount = 100,000 - 4,440 = 95,560
```

---

### Convenience Store
| Payment Method | Fee Rate | Tax | Formula |
|----------------|----------|-----|---------|
| `ONLINE_TO_OFFLINE_ALFA` | IDR 5,000 | 11% | fee = 5000; tax = fee × 0.11 |
| `ONLINE_TO_OFFLINE_INDOMARET` | IDR 6,500 | 11% | fee = 6500; tax = fee × 0.11 |

**Example Alfamart (IDR 100,000):**
```
Transaction Fee = 5,000
Tax = 5,000 × 11% = 550
Total Deduction = 5,000 + 550 = 5,550
Net Amount = 100,000 - 5,550 = 94,450
```

**Example Indomaret (IDR 100,000):**
```
Transaction Fee = 6,500
Tax = 6,500 × 11% = 715
Total Deduction = 6,500 + 715 = 7,215
Net Amount = 100,000 - 7,215 = 92,785
```

---

### QRIS
| Payment Method | Fee Rate | Tax | Formula |
|----------------|----------|-----|---------|
| `QRIS` | IDR 700 | **No Tax** | fee = 700; tax = 0 |

**Example (IDR 100,000):**
```
Transaction Fee = 700
Tax = 0 (QRIS exempt from tax)
Total Deduction = 700
Net Amount = 100,000 - 700 = 99,300
```

---

### E-Wallet
| Payment Method | Fee Rate | Tax | Formula |
|----------------|----------|-----|---------|
| `EMONEY_SHOPEE_PAY` | 2% | 11% | fee = amount × 0.02; tax = fee × 0.11 |
| `EMONEY_OVO` | 2% | 11% | fee = amount × 0.02; tax = fee × 0.11 |
| `EMONEY_LINKAJA` | 2% | 11% | fee = amount × 0.02; tax = fee × 0.11 |
| `EMONEY_DOKU` | 1.5% | 11% | fee = amount × 0.015; tax = fee × 0.11 |
| `EMONEY_DANA` | 1.5% | 11% | fee = amount × 0.015; tax = fee × 0.11 |

**Example ShopeePay/OVO/LinkAja (IDR 100,000):**
```
Transaction Fee = 100,000 × 2% = 2,000
Tax = 2,000 × 11% = 220
Total Deduction = 2,000 + 220 = 2,220
Net Amount = 100,000 - 2,220 = 97,780
```

**Example DOKU/DANA Wallet (IDR 100,000):**
```
Transaction Fee = 100,000 × 1.5% = 1,500
Tax = 1,500 × 11% = 165
Total Deduction = 1,500 + 165 = 1,665
Net Amount = 100,000 - 1,665 = 98,335
```

---

### PayLater
| Payment Method | Fee Rate | Tax | Formula |
|----------------|----------|-----|---------|
| `PEER_TO_PEER_AKULAKU` | 1.5% | 11% | fee = amount × 0.015; tax = fee × 0.11 |
| `PEER_TO_PEER_KREDIVO` | 2.3% | 11% | fee = amount × 0.023; tax = fee × 0.11 |
| `PEER_TO_PEER_INDODANA` | 2.3% | 11% | fee = amount × 0.023; tax = fee × 0.11 |

**Example Akulaku (IDR 100,000):**
```
Transaction Fee = 100,000 × 1.5% = 1,500
Tax = 1,500 × 11% = 165
Total Deduction = 1,500 + 165 = 1,665
Net Amount = 100,000 - 1,665 = 98,335
```

**Example Kredivo/Indodana (IDR 100,000):**
```
Transaction Fee = 100,000 × 2.3% = 2,300
Tax = 2,300 × 11% = 253
Total Deduction = 2,300 + 253 = 2,553
Net Amount = 100,000 - 2,553 = 97,447
```

---

### Direct Debit
| Payment Method | Fee Rate | Tax | Formula |
|----------------|----------|-----|---------|
| `DIRECT_DEBIT_BRI` | 2% | 11% | fee = amount × 0.02; tax = fee × 0.11 |

**Example (IDR 100,000):**
```
Transaction Fee = 100,000 × 2% = 2,000
Tax = 2,000 × 11% = 220
Total Deduction = 2,000 + 220 = 2,220
Net Amount = 100,000 - 2,220 = 97,780
```

---

### Digital Banking
| Payment Method | Fee Rate | Tax | Formula |
|----------------|----------|-----|---------|
| `JENIUS_PAY` | 1.5% | 11% | fee = amount × 0.015; tax = fee × 0.11 |

**Example (IDR 100,000):**
```
Transaction Fee = 100,000 × 1.5% = 1,500
Tax = 1,500 × 11% = 165
Total Deduction = 1,500 + 165 = 1,665
Net Amount = 100,000 - 1,665 = 98,335
```

---

## Implementation Logic

### CalculateSettlementFee Method

```go
func (u *dokuSettlementUseCase) CalculateSettlementFee(
    paymentMethod string, 
    amount float64,
) (*responses.DokuSettlementResultResponse, error) {
    
    // 1. Validate inputs
    if paymentMethod == "" {
        return nil, errors.New("payment method is empty")
    }
    if amount <= 0 {
        return nil, errors.New("invalid amount: must be greater than 0")
    }

    // 2. Calculate fee and tax based on payment method
    transactionFee, tax, err := u.calculateFeeAndTax(paymentMethod, amount)
    if err != nil {
        return nil, err
    }

    // 3. Calculate totals
    totalDeduction := transactionFee + tax
    netAmount := amount - totalDeduction

    // 4. Return result with rounded values
    return &responses.DokuSettlementResultResponse{
        PaymentMethod:  paymentMethod,
        GrossAmount:    amount,
        TransactionFee: roundToTwoDecimals(transactionFee),
        Tax:            roundToTwoDecimals(tax),
        TotalDeduction: roundToTwoDecimals(totalDeduction),
        NetAmount:      roundToTwoDecimals(netAmount),
    }, nil
}
```

### Fee Calculation Logic

```go
func (u *dokuSettlementUseCase) calculateFeeAndTax(
    paymentMethod string, 
    amount float64,
) (transactionFee float64, tax float64, err error) {
    
    taxRate := float64(u.cfg.TransactionFee.Tax) / 100 // 11% = 0.11

    switch paymentMethod {
    
    // Cards: Percentage + Flat Fee
    case constants.CREDIT_CARD, constants.KKI:
        percentageRate := u.cfg.TransactionFee.Cards.PercentageRate / 100
        flatFee := float64(u.cfg.TransactionFee.Cards.FlatFee)
        transactionFee = (amount * percentageRate) + flatFee
        tax = transactionFee * taxRate

    // Virtual Account: Flat Fee Only
    case constants.BCA_VA, constants.Mandiri_VA, constants.BSI_VA, 
         constants.BRI_VA, constants.BNI_VA, constants.DOKU_VA,
         constants.PERMATA_VA, constants.CIMB_VA, constants.DANAMON_VA,
         constants.BTN_VA, constants.BNC_VA:
        transactionFee = float64(u.cfg.TransactionFee.VirtualAccount.FlatFee)
        tax = transactionFee * taxRate

    // Convenience Store: Flat Fee
    case constants.ALFA_GROUP:
        transactionFee = float64(u.cfg.TransactionFee.ConvenienceStore.Alfamart.FlatFee)
        tax = transactionFee * taxRate

    case constants.INDOMARET:
        transactionFee = float64(u.cfg.TransactionFee.ConvenienceStore.Indomaret.FlatFee)
        tax = transactionFee * taxRate

    // QRIS: Flat Fee, No Tax
    case constants.QRIS:
        transactionFee = float64(u.cfg.TransactionFee.QR.FlatFee)
        tax = 0 // QRIS exempt from tax

    // E-Wallet: Percentage Fee
    case constants.SHOPEEPAY:
        percentageRate := u.cfg.TransactionFee.EWallet.ShopeePay.PercentageRate / 100
        transactionFee = amount * percentageRate
        tax = transactionFee * taxRate

    case constants.OVO:
        percentageRate := u.cfg.TransactionFee.EWallet.OVO.PercentageRate / 100
        transactionFee = amount * percentageRate
        tax = transactionFee * taxRate

    case constants.LINKAJA:
        percentageRate := u.cfg.TransactionFee.EWallet.LinkAja.PercentageRate / 100
        transactionFee = amount * percentageRate
        tax = transactionFee * taxRate

    case constants.DOKU_WALLET:
        percentageRate := u.cfg.TransactionFee.EWallet.Doku.PercentageRate / 100
        transactionFee = amount * percentageRate
        tax = transactionFee * taxRate

    case constants.DANA:
        percentageRate := u.cfg.TransactionFee.EWallet.Dana.PercentageRate / 100
        transactionFee = amount * percentageRate
        tax = transactionFee * taxRate

    // PayLater: Percentage Fee
    case constants.PAYLATER_AKULAKU:
        percentageRate := u.cfg.TransactionFee.PayLater.Akulaku.PercentageRate / 100
        transactionFee = amount * percentageRate
        tax = transactionFee * taxRate

    case constants.PAYLATER_KREDIVO:
        percentageRate := u.cfg.TransactionFee.PayLater.Kredivo.PercentageRate / 100
        transactionFee = amount * percentageRate
        tax = transactionFee * taxRate

    case constants.PAYLATER_INDODANA:
        percentageRate := u.cfg.TransactionFee.PayLater.Indodana.PercentageRate / 100
        transactionFee = amount * percentageRate
        tax = transactionFee * taxRate

    // Direct Debit
    case constants.DIRECT_DEBIT_BRI:
        percentageRate := u.cfg.TransactionFee.DirectDebit.BRI.PercentageRate / 100
        transactionFee = amount * percentageRate
        tax = transactionFee * taxRate

    // Digital Banking
    case constants.JENIUS_PAY:
        percentageRate := u.cfg.TransactionFee.DigitalBanking.JeniusPay.PercentageRate / 100
        transactionFee = amount * percentageRate
        tax = transactionFee * taxRate

    default:
        return 0, 0, errors.New("unknown payment method: " + paymentMethod)
    }

    return transactionFee, tax, nil
}

func roundToTwoDecimals(value float64) float64 {
    return math.Round(value*100) / 100
}
```

---

## Usage Example

### Calculate Settlement Fee (Standard)

```go
// Initialize use case
settlementUseCase := usecases.NewDokuSettlementUseCase()

// Calculate settlement for a Virtual Account payment
result, err := settlementUseCase.CalculateSettlementFee(
    constants.BCA_VA,  // Payment method
    100000,            // Gross amount
)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Payment Method: %s\n", result.PaymentMethod)
fmt.Printf("Gross Amount: %.2f\n", result.GrossAmount)
fmt.Printf("Transaction Fee: %.2f\n", result.TransactionFee)
fmt.Printf("Tax: %.2f\n", result.Tax)
fmt.Printf("Total Deduction: %.2f\n", result.TotalDeduction)
fmt.Printf("Net Amount: %.2f\n", result.NetAmount)

// Output:
// Payment Method: VIRTUAL_ACCOUNT_BCA
// Gross Amount: 100000.00
// Transaction Fee: 4000.00
// Tax: 440.00
// Total Deduction: 4440.00
// Net Amount: 95560.00
```

---

## Gross Amount Calculation (Upsert Feature)

### Overview

The **CalculateGrossAmount** function solves the reverse problem: given a desired net amount that the merchant wants to receive, calculate the gross amount (what the customer must pay) so that after all fees and taxes are deducted, the merchant receives exactly (or slightly more due to rounding) the desired net amount.

This is essential when:
- Merchant wants to receive exactly IDR 100,000 after all fees
- Customer should bear all transaction costs
- Pricing must guarantee merchant receives the full product/service value

### Formula Derivation

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    GROSS AMOUNT CALCULATION (UPSERT)                            │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│   Desired Net Amount (Merchant Receives)                                        │
│         │                                                                       │
│         ▼                                                                       │
│   ┌─────────────────────────────────────────────────┐                          │
│   │  Reverse Calculation                            │                          │
│   │                                                 │                          │
│   │  For flat fee with tax:                        │                          │
│   │    gross = net + flatFee × (1 + taxRate)       │                          │
│   │                                                 │                          │
│   │  For percentage with tax:                       │                          │
│   │    gross = net / (1 - rate × (1 + taxRate))    │                          │
│   │                                                 │                          │
│   │  For mixed (% + flat) with tax:                │                          │
│   │    gross = (net + flat×(1+tax))                │                          │
│   │            / (1 - rate×(1+tax))                │                          │
│   └─────────────────────────────────────────────────┘                          │
│         │                                                                       │
│         ▼                                                                       │
│   Gross Amount (Customer Pays)                                                  │
│   (Rounded up to ensure merchant gets at least desired net)                     │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### Mathematical Formula

**For percentage-based fees with tax:**
```
netAmount = grossAmount × (1 - percentageRate × (1 + taxRate))

Solving for grossAmount:
grossAmount = netAmount / (1 - percentageRate × (1 + taxRate))
```

**For flat fee with tax:**
```
netAmount = grossAmount - flatFee × (1 + taxRate)

Solving for grossAmount:
grossAmount = netAmount + flatFee × (1 + taxRate)
```

**For mixed (percentage + flat fee) with tax:**
```
grossAmount = (netAmount + flatFee × (1 + taxRate)) / (1 - percentageRate × (1 + taxRate))
```

### Implementation

```go
func (u *dokuSettlementUseCase) CalculateGrossAmount(
    paymentMethod string, 
    desiredNetAmount float64,
) (*responses.DokuSettlementResultResponse, error) {
    
    // 1. Validate inputs
    if paymentMethod == "" {
        return nil, errors.New("payment method is empty")
    }
    if desiredNetAmount <= 0 {
        return nil, errors.New("invalid desired net amount: must be greater than 0")
    }

    // 2. Get fee parameters for this payment method
    percentageRate, flatFee, hasTax, err := u.getFeeParameters(paymentMethod)
    if err != nil {
        return nil, err
    }

    taxRate := float64(0)
    if hasTax {
        taxRate = float64(u.cfg.TransactionFee.Tax) / 100
    }

    // 3. Calculate gross amount using inverse formula
    taxMultiplier := 1 + taxRate
    divisor := 1 - (percentageRate * taxMultiplier)

    if divisor <= 0 {
        return nil, errors.New("invalid fee configuration: fees exceed 100%")
    }

    grossAmount := (desiredNetAmount + (flatFee * taxMultiplier)) / divisor

    // 4. Round up to ensure merchant receives at least desired net
    grossAmount = math.Ceil(grossAmount)

    // 5. Verify by calculating actual settlement
    transactionFee, tax, _ := u.calculateFeeAndTax(paymentMethod, grossAmount)
    totalDeduction := transactionFee + tax
    actualNetAmount := grossAmount - totalDeduction

    return &responses.DokuSettlementResultResponse{
        PaymentMethod:  paymentMethod,
        GrossAmount:    grossAmount,
        TransactionFee: roundToTwoDecimals(transactionFee),
        Tax:            roundToTwoDecimals(tax),
        TotalDeduction: roundToTwoDecimals(totalDeduction),
        NetAmount:      roundToTwoDecimals(actualNetAmount),
    }, nil
}
```

### Usage Example

```go
// Merchant wants to receive exactly IDR 100,000
desiredNet := float64(100000)

// Calculate what customer must pay
result, err := settlementUseCase.CalculateGrossAmount(
    constants.BCA_VA,  // Payment method
    desiredNet,        // Desired net amount
)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Desired Net: %.2f\n", desiredNet)
fmt.Printf("Customer Pays (Gross): %.2f\n", result.GrossAmount)
fmt.Printf("Transaction Fee: %.2f\n", result.TransactionFee)
fmt.Printf("Tax: %.2f\n", result.Tax)
fmt.Printf("Merchant Receives (Net): %.2f\n", result.NetAmount)

// Output:
// Desired Net: 100000.00
// Customer Pays (Gross): 104440.00
// Transaction Fee: 4000.00
// Tax: 440.00
// Merchant Receives (Net): 100000.00
```

### Gross Amount Summary Table

| Payment Method | Desired Net | Customer Pays (Gross) | Fee | Tax | Actual Net |
|----------------|-------------|----------------------|-----|-----|------------|
| Virtual Account (BCA, etc.) | 100,000 | 104,440 | 4,000 | 440 | 100,000 |
| Alfamart | 100,000 | 105,550 | 5,000 | 550 | 100,000 |
| Indomaret | 100,000 | 107,215 | 6,500 | 715 | 100,000 |
| QRIS | 100,000 | 100,700 | 700 | 0 | 100,000 |
| ShopeePay/OVO/LinkAja | 100,000 | 102,271 | 2,045.42 | 225.00 | 100,000.58 |
| DOKU/DANA Wallet | 100,000 | 101,694 | 1,525.41 | 167.80 | 100,000.79 |
| Credit Card | 100,000 | 105,499 | 4,953.97 | 544.94 | 100,000.09 |
| Kredivo/Indodana | 100,000 | 102,620 | 2,360.26 | 259.63 | 100,000.11 |
| Akulaku | 100,000 | 101,694 | 1,525.41 | 167.80 | 100,000.79 |

> **Note:** For percentage-based fees, the gross amount is rounded up to the nearest whole number to ensure the merchant always receives at least the desired net amount. This may result in the actual net being slightly higher than requested.

---

## Integration with Ledger System

```go
// In setter-service when processing settlement
func (s *settlementService) ProcessSettlement(
    ctx context.Context,
    paymentMethod string,
    grossAmount float64,
    ledgerAccountUUID string,
) error {
    
    // 1. Calculate fee using DOKU module
    settlementResult, err := s.dokuSettlementUseCase.CalculateSettlementFee(
        paymentMethod,
        grossAmount,
    )
    if err != nil {
        return err
    }
    
    // 2. Create settlement in Ledger
    _, err = s.ledgerSettlementUseCase.CreateSettlement(tx, &ledger.LedgerSettlementCreateRequest{
        LedgerAccountUUID: ledgerAccountUUID,
        GrossAmount:       int64(settlementResult.GrossAmount),
        NetAmount:         int64(settlementResult.NetAmount),
        FeeAmount:         int64(settlementResult.TotalDeduction),
    })
    if err != nil {
        return err
    }
    
    // 3. Complete settlement (moves pending to available balance)
    _, err = s.ledgerSettlementUseCase.CompleteSettlement(tx, settlementUUID)
    if err != nil {
        return err
    }
    
    return nil
}
```

---

## Fee Summary Table

| Payment Method | Fee Type | Fee Rate | Tax | Net from 100K Gross | Gross for 100K Net |
|----------------|----------|----------|-----|---------------------|-------------------|
| Credit Card | % + Flat | 2.8% + 2,000 | 11% | 94,672 | 105,499 |
| Virtual Account | Flat | 4,000 | 11% | 95,560 | 104,440 |
| Alfamart | Flat | 5,000 | 11% | 94,450 | 105,550 |
| Indomaret | Flat | 6,500 | 11% | 92,785 | 107,215 |
| QRIS | Flat | 700 | 0% | 99,300 | 100,700 |
| ShopeePay | % | 2% | 11% | 97,780 | 102,271 |
| OVO | % | 2% | 11% | 97,780 | 102,271 |
| LinkAja | % | 2% | 11% | 97,780 | 102,271 |
| DOKU Wallet | % | 1.5% | 11% | 98,335 | 101,694 |
| DANA | % | 1.5% | 11% | 98,335 | 101,694 |
| Akulaku | % | 1.5% | 11% | 98,335 | 101,694 |
| Kredivo | % | 2.3% | 11% | 97,447 | 102,620 |
| Indodana | % | 2.3% | 11% | 97,447 | 102,620 |
| Direct Debit BRI | % | 2% | 11% | 97,780 | 102,271 |
| Jenius Pay | % | 1.5% | 11% | 98,335 | 101,694 |

---

## Wallet Impact on Settlement

| Action | pending_balance | balance (available) | Description |
|--------|-----------------|---------------------|-------------|
| Payment Confirmed | +gross_amount | - | Funds waiting for settlement |
| Settlement Completed | -gross_amount | +net_amount | Funds moved to available, fees deducted |

---

## Configuration

Fee rates can be configured via environment variables:

```bash
# Cards
TRANSACTION_FEE_CARDS_PERCENTAGE_RATE=2.8
TRANSACTION_FEE_CARDS_FLAT_FEE=2000

# Virtual Account
TRANSACTION_FEE_VIRTUAL_ACCOUNT_FLAT_FEE=4000

# Convenience Store
TRANSACTION_FEE_ALFAMART_FLAT_FEE=5000
TRANSACTION_FEE_INDOMARET_FLAT_FEE=6500

# QRIS
TRANSACTION_FEE_QR_FLAT_FEE=700

# E-Wallet
TRANSACTION_FEE_SHOPEEPAY_PERCENTAGE_RATE=2.0
TRANSACTION_FEE_OVO_PERCENTAGE_RATE=2.0
TRANSACTION_FEE_LINKAJA_PERCENTAGE_RATE=2.0
TRANSACTION_FEE_DOKU_WALLET_PERCENTAGE_RATE=1.5
TRANSACTION_FEE_DANA_PERCENTAGE_RATE=1.5

# PayLater
TRANSACTION_FEE_AKULAKU_PERCENTAGE_RATE=1.5
TRANSACTION_FEE_KREDIVO_PERCENTAGE_RATE=2.3
TRANSACTION_FEE_INDODANA_PERCENTAGE_RATE=2.3

# Direct Debit
TRANSACTION_FEE_BRI_DIRECT_DEBIT_PERCENTAGE_RATE=2.0

# Digital Banking
TRANSACTION_FEE_JENIUSPAY_PERCENTAGE_RATE=1.5

# Tax Rate (PPN)
TRANSACTION_FEE_TAX=11
```

---

## API Reference

### Interface Methods

```go
type DokuSettlementUseCaseInterface interface {
    // CalculateSettlementFee calculates the net amount from a given gross amount
    // Use when you know what the customer paid and need to calculate merchant's net
    CalculateSettlementFee(paymentMethod string, amount float64) (*responses.DokuSettlementResultResponse, error)
    
    // CalculateGrossAmount calculates the gross amount needed for a desired net amount
    // Use when you need to determine what customer should pay for merchant to receive exact amount
    CalculateGrossAmount(paymentMethod string, desiredNetAmount float64) (*responses.DokuSettlementResultResponse, error)
}
```

### When to Use Each Method

| Scenario | Method | Example |
|----------|--------|---------|
| Customer paid IDR 100,000, what does merchant receive? | `CalculateSettlementFee` | Input: 100,000 → Net: 95,560 (VA) |
| Merchant needs exactly IDR 100,000, what should customer pay? | `CalculateGrossAmount` | Input: 100,000 → Gross: 104,440 (VA) |
| Standard checkout flow | `CalculateSettlementFee` | Use after payment notification |
| Fee-inclusive pricing | `CalculateGrossAmount` | Use when creating payment request |

---

## Balance Reconciliation Integration

### Overview

DOKU settles payments daily at **1PM on weekdays**, but **does not provide a webhook** for settlement completion. To detect when settlements have been processed, the consuming service (setter-service) uses an **on-demand reconciliation** approach by calling the DOKU `GetBalance` API.

### GetBalance API Response

```go
type DokuGetBalanceHTTPResponse struct {
    Balance *DokuBalance `json:"balance"`
}

type DokuBalance struct {
    Pending   null.String `json:"pending"`   // Gross amount waiting for settlement
    Available null.String `json:"available"` // Net amount available for disbursement
}
```

### Settlement Detection Logic

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                     SETTLEMENT DETECTION LOGIC                                   │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  DOKU Balance (Source of Truth)          Ledger Balance (Local Cache)           │
│  ───────────────────────────────         ─────────────────────────────          │
│  pending: 0                              pending_balance: 100,000               │
│  available: 95,560                       balance: 0                             │
│                                                                                 │
│                              ▼                                                  │
│                                                                                 │
│  Delta Calculation:                                                             │
│  ─────────────────                                                              │
│  delta = ledger_pending - doku_pending                                          │
│  delta = 100,000 - 0 = 100,000                                                  │
│                                                                                 │
│  If delta > 0:                                                                  │
│    → Settlement(s) have been processed!                                         │
│    → Process IN_PROGRESS settlements (FIFO)                                     │
│    → Update ledger: pending -= gross, balance += net                            │
│                                                                                 │
│  If delta == 0:                                                                 │
│    → Balances in sync, no action needed                                         │
│                                                                                 │
│  If delta < 0:                                                                  │
│    → DOKU has more pending than ledger (data integrity issue)                   │
│    → Log warning, no action                                                     │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### Integration Flow (setter-service)

When user visits the balance page:

```go
// 1. Get DOKU balance (real-time from DOKU API)
dokuBalance, err := dokuUseCase.GetBalance(userDoku.SubAccountID)

// 2. Compare with ledger pending balance
dokuPending := parseDokuBalance(dokuBalance.Balance.Pending)
ledgerPending := ledgerWallet.PendingBalance

delta := ledgerPending - dokuPending

// 3. If delta > 0, reconcile settlements
if delta > 0 {
    // Get IN_PROGRESS settlements (FIFO order)
    settlements := ledgerSettlementUseCase.GetSettlementsByAccountAndStatus(
        ledgerAccountUUID,
        SettlementStatusInProgress,
    )
    
    // Process settlements until delta is satisfied
    for _, settlement := range settlements {
        // Update status to TRANSFERRED
        ledgerSettlementUseCase.UpdateSettlementStatus(tx, settlement.UUID, "TRANSFERRED", &now)
        
        // Update wallet: -pending (gross), +available (net)
        ledgerWalletUseCase.SettlePendingBalance(tx, walletUUID, settlement.GrossAmount, settlement.NetAmount)
    }
}

// 4. Return updated balance to user
```

### Settlement Timeline

```
Day 1 (Monday) 10:00 AM - Customer pays IDR 100,000 via Virtual Account
├── DOKU pending: +100,000 (gross)
├── DOKU available: unchanged
├── Ledger pending_balance: +100,000
├── Ledger balance: unchanged
└── LedgerSettlement: status = IN_PROGRESS, gross = 100,000, net = 95,560

Day 2 (Tuesday) 1:00 PM - DOKU processes settlement (NO WEBHOOK)
├── DOKU pending: 100,000 → 0
├── DOKU available: 0 → 95,560
├── Ledger: UNCHANGED (we don't know yet!)
└── LedgerSettlement: still IN_PROGRESS

Day 2 (Tuesday) 3:00 PM - User visits balance page
├── Backend calls GetBalance API → pending: 0, available: 95,560
├── Delta = 100,000 - 0 = 100,000 (settlement detected!)
├── Process settlement: status → TRANSFERRED
├── Ledger pending_balance: 100,000 → 0
├── Ledger balance: 0 → 95,560
└── User sees: available = 95,560, pending = 0
```

### Key Points

| Aspect | Details |
|--------|---------|
| **Trigger** | User accessing balance page |
| **Source of Truth** | DOKU GetBalance API |
| **Detection Method** | Compare DOKU pending vs Ledger pending_balance |
| **Processing Order** | FIFO (oldest settlements first) |
| **Atomicity** | Single database transaction for all updates |
| **Failure Handling** | If DOKU API fails, return cached ledger balance |

### Fee Calculation During Reconciliation

When processing settlements, the fee has already been calculated and stored:

```
LedgerSettlement record:
├── gross_amount: 100,000  (what customer paid)
├── net_amount: 95,560     (what merchant receives)
└── fee_amount: 4,440      (DOKU fee + tax)

Wallet update:
├── pending_balance -= gross_amount (100,000)
└── balance += net_amount (95,560)
```

The fee calculation functions (`CalculateSettlementFee`, `CalculateGrossAmount`) are used when **creating** the settlement record, not during reconciliation.
