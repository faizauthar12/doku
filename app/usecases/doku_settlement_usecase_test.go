package usecases

import (
	"fmt"
	"testing"

	"github.com/faizauthar12/doku/app/config"
	"github.com/faizauthar12/doku/app/constants"
)

func init() {
	// Initialize config with default values for testing
	config.InitConfigFromEnv()
}

func TestCalculateSettlementFee(t *testing.T) {
	usecase := NewDokuSettlementUseCase()

	testCases := []struct {
		name              string
		paymentMethod     string
		amount            float64
		expectedNetAmount float64
		description       string
	}{
		{
			name:              "Transfer Bank - Virtual Account (Generic)",
			paymentMethod:     constants.VIRTUAL_ACCOUNT,
			amount:            100000,
			expectedNetAmount: 95560,
			description:       "Fee: 4000, Tax: 440, Net: 95560",
		},
		{
			name:              "Transfer Bank - BCA VA",
			paymentMethod:     constants.BCA_VA,
			amount:            100000,
			expectedNetAmount: 95560,
			description:       "Fee: 4000, Tax: 440, Net: 95560",
		},
		{
			name:              "Transfer Bank - Mandiri VA",
			paymentMethod:     constants.Mandiri_VA,
			amount:            100000,
			expectedNetAmount: 95560,
			description:       "Fee: 4000, Tax: 440, Net: 95560",
		},
		{
			name:              "Alfamart",
			paymentMethod:     constants.ALFA_GROUP,
			amount:            100000,
			expectedNetAmount: 94450,
			description:       "Fee: 5000, Tax: 550, Net: 94450",
		},
		{
			name:              "Indomaret",
			paymentMethod:     constants.INDOMARET,
			amount:            100000,
			expectedNetAmount: 92785,
			description:       "Fee: 6500, Tax: 715, Net: 92785",
		},
		{
			name:              "QRIS",
			paymentMethod:     constants.QRIS,
			amount:            100000,
			expectedNetAmount: 99300,
			description:       "Fee: 700, No Tax, Net: 99300",
		},
		{
			name:              "ShopeePay",
			paymentMethod:     constants.SHOPEEPAY,
			amount:            100000,
			expectedNetAmount: 97780,
			description:       "Fee: 2%, Tax: 11%, Net: 97780",
		},
		{
			name:              "OVO",
			paymentMethod:     constants.OVO,
			amount:            100000,
			expectedNetAmount: 97780,
			description:       "Fee: 2%, Tax: 11%, Net: 97780",
		},
		{
			name:              "LinkAja",
			paymentMethod:     constants.LINKAJA,
			amount:            100000,
			expectedNetAmount: 97780,
			description:       "Fee: 2%, Tax: 11%, Net: 97780",
		},
		{
			name:              "DOKU Wallet",
			paymentMethod:     constants.DOKU_WALLET,
			amount:            100000,
			expectedNetAmount: 98335,
			description:       "Fee: 1.5%, Tax: 11%, Net: 98335",
		},
		{
			name:              "DANA Wallet",
			paymentMethod:     constants.DANA,
			amount:            100000,
			expectedNetAmount: 98335,
			description:       "Fee: 1.5%, Tax: 11%, Net: 98335",
		},
		{
			name:              "Paylater Akulaku",
			paymentMethod:     constants.PAYLATER_AKULAKU,
			amount:            100000,
			expectedNetAmount: 98335,
			description:       "Fee: 1.5%, Tax: 11%, Net: 98335",
		},
		{
			name:              "Paylater Kredivo",
			paymentMethod:     constants.PAYLATER_KREDIVO,
			amount:            100000,
			expectedNetAmount: 97447,
			description:       "Fee: 2.3%, Tax: 11%, Net: 97447",
		},
		{
			name:              "Paylater Indodana",
			paymentMethod:     constants.PAYLATER_INDODANA,
			amount:            100000,
			expectedNetAmount: 97447,
			description:       "Fee: 2.3%, Tax: 11%, Net: 97447",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := usecase.CalculateSettlementFee(tc.paymentMethod, tc.amount)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.NetAmount != tc.expectedNetAmount {
				t.Errorf("%s: expected net amount %.2f, got %.2f", tc.description, tc.expectedNetAmount, result.NetAmount)
			}

			if result.GrossAmount != tc.amount {
				t.Errorf("expected gross amount %.2f, got %.2f", tc.amount, result.GrossAmount)
			}

			if result.PaymentMethod != tc.paymentMethod {
				t.Errorf("expected payment method %s, got %s", tc.paymentMethod, result.PaymentMethod)
			}

			// Verify total deduction calculation
			expectedDeduction := result.GrossAmount - result.NetAmount
			if result.TotalDeduction != expectedDeduction {
				t.Errorf("total deduction mismatch: expected %.2f, got %.2f", expectedDeduction, result.TotalDeduction)
			}

			// Verify fee + tax = total deduction
			feeAndTax := result.TransactionFee + result.Tax
			if roundToTwoDecimals(feeAndTax) != result.TotalDeduction {
				t.Errorf("fee + tax mismatch: expected %.2f, got %.2f", result.TotalDeduction, feeAndTax)
			}

			fmt.Printf("Amount: %.2f, Payment Method: %s, Gross: %.2f, Fee: %.2f, Tax: %.2f, Net: %.2f\n", tc.amount, tc.paymentMethod, result.GrossAmount, result.TransactionFee, result.Tax, result.NetAmount)
		})
	}
}

func TestCalculateSettlementFee_Cards(t *testing.T) {
	usecase := NewDokuSettlementUseCase()

	result, err := usecase.CalculateSettlementFee(constants.CREDIT_CARD, 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Cards: 2.8% + 2000 = 4800, Tax: 528, Total: 5328, Net: 94672
	expectedFee := 100000*0.028 + 2000
	expectedTax := expectedFee * 0.11
	expectedNet := 100000 - (expectedFee + expectedTax)

	if roundToTwoDecimals(result.TransactionFee) != roundToTwoDecimals(expectedFee) {
		t.Errorf("expected transaction fee %.2f, got %.2f", expectedFee, result.TransactionFee)
	}

	if roundToTwoDecimals(result.Tax) != roundToTwoDecimals(expectedTax) {
		t.Errorf("expected tax %.2f, got %.2f", expectedTax, result.Tax)
	}

	if roundToTwoDecimals(result.NetAmount) != roundToTwoDecimals(expectedNet) {
		t.Errorf("expected net amount %.2f, got %.2f", expectedNet, result.NetAmount)
	}
}

func TestCalculateSettlementFee_ErrorCases(t *testing.T) {
	usecase := NewDokuSettlementUseCase()

	t.Run("Empty payment method", func(t *testing.T) {
		_, err := usecase.CalculateSettlementFee("", 100000)
		if err == nil {
			t.Error("expected error for empty payment method")
		}
	})

	t.Run("Zero amount", func(t *testing.T) {
		_, err := usecase.CalculateSettlementFee(constants.BCA_VA, 0)
		if err == nil {
			t.Error("expected error for zero amount")
		}
	})

	t.Run("Negative amount", func(t *testing.T) {
		_, err := usecase.CalculateSettlementFee(constants.BCA_VA, -100000)
		if err == nil {
			t.Error("expected error for negative amount")
		}
	})

	t.Run("Unknown payment method", func(t *testing.T) {
		_, err := usecase.CalculateSettlementFee("UNKNOWN_PAYMENT_METHOD", 100000)
		if err == nil {
			t.Error("expected error for unknown payment method")
		}
	})
}

func TestLargeAmountCalculations(t *testing.T) {
	usecase := NewDokuSettlementUseCase()

	// Test with a large amount (10 million IDR)
	result, err := usecase.CalculateSettlementFee(constants.SHOPEEPAY, 10000000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ShopeePay: 2% fee + 11% tax
	// Fee: 10000000 * 0.02 = 200000
	// Tax: 200000 * 0.11 = 22000
	// Net: 10000000 - 222000 = 9778000
	expectedNet := float64(9778000)
	if result.NetAmount != expectedNet {
		t.Errorf("expected net amount %.2f, got %.2f", expectedNet, result.NetAmount)
	}
}

func TestCalculateGrossAmount(t *testing.T) {
	usecase := NewDokuSettlementUseCase()

	testCases := []struct {
		name             string
		paymentMethod    string
		desiredNetAmount float64
		description      string
	}{
		{
			name:             "Transfer Bank - Virtual Account (Generic)",
			paymentMethod:    constants.VIRTUAL_ACCOUNT,
			desiredNetAmount: 100000,
			description:      "Flat fee IDR 4000 + 11% tax",
		},
		{
			name:             "Transfer Bank - BCA VA",
			paymentMethod:    constants.BCA_VA,
			desiredNetAmount: 100000,
			description:      "Flat fee IDR 4000 + 11% tax",
		},
		{
			name:             "Transfer Bank - Mandiri VA",
			paymentMethod:    constants.Mandiri_VA,
			desiredNetAmount: 100000,
			description:      "Flat fee IDR 4000 + 11% tax",
		},
		{
			name:             "Alfamart",
			paymentMethod:    constants.ALFA_GROUP,
			desiredNetAmount: 100000,
			description:      "Flat fee IDR 5000 + 11% tax",
		},
		{
			name:             "Indomaret",
			paymentMethod:    constants.INDOMARET,
			desiredNetAmount: 100000,
			description:      "Flat fee IDR 6500 + 11% tax",
		},
		{
			name:             "QRIS",
			paymentMethod:    constants.QRIS,
			desiredNetAmount: 100000,
			description:      "Flat fee IDR 700, no tax",
		},
		{
			name:             "ShopeePay",
			paymentMethod:    constants.SHOPEEPAY,
			desiredNetAmount: 100000,
			description:      "2% fee + 11% tax",
		},
		{
			name:             "OVO",
			paymentMethod:    constants.OVO,
			desiredNetAmount: 100000,
			description:      "2% fee + 11% tax",
		},
		{
			name:             "LinkAja",
			paymentMethod:    constants.LINKAJA,
			desiredNetAmount: 100000,
			description:      "2% fee + 11% tax",
		},
		{
			name:             "DOKU Wallet",
			paymentMethod:    constants.DOKU_WALLET,
			desiredNetAmount: 100000,
			description:      "1.5% fee + 11% tax",
		},
		{
			name:             "DANA Wallet",
			paymentMethod:    constants.DANA,
			desiredNetAmount: 100000,
			description:      "1.5% fee + 11% tax",
		},
		{
			name:             "Credit Card",
			paymentMethod:    constants.CREDIT_CARD,
			desiredNetAmount: 100000,
			description:      "2.8% + IDR 2000 + 11% tax",
		},
		{
			name:             "Paylater Akulaku",
			paymentMethod:    constants.PAYLATER_AKULAKU,
			desiredNetAmount: 100000,
			description:      "1.5% fee + 11% tax",
		},
		{
			name:             "Paylater Kredivo",
			paymentMethod:    constants.PAYLATER_KREDIVO,
			desiredNetAmount: 100000,
			description:      "2.3% fee + 11% tax",
		},
		{
			name:             "Paylater Indodana",
			paymentMethod:    constants.PAYLATER_INDODANA,
			desiredNetAmount: 100000,
			description:      "2.3% fee + 11% tax",
		},
		{
			name:             "Direct Debit BRI",
			paymentMethod:    constants.DIRECT_DEBIT_BRI,
			desiredNetAmount: 100000,
			description:      "2% fee + 11% tax",
		},
		{
			name:             "Jenius Pay",
			paymentMethod:    constants.JENIUS_PAY,
			desiredNetAmount: 100000,
			description:      "1.5% fee + 11% tax",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := usecase.CalculateGrossAmount(tc.paymentMethod, tc.desiredNetAmount)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// The net amount should be >= desired net amount (due to rounding up)
			if result.NetAmount < tc.desiredNetAmount {
				t.Errorf("%s: net amount %.2f is less than desired %.2f", tc.description, result.NetAmount, tc.desiredNetAmount)
			}

			// Verify the calculation is correct
			expectedDeduction := result.GrossAmount - result.NetAmount
			if roundToTwoDecimals(result.TotalDeduction) != roundToTwoDecimals(expectedDeduction) {
				t.Errorf("total deduction mismatch: expected %.2f, got %.2f", expectedDeduction, result.TotalDeduction)
			}

			// Verify fee + tax = total deduction
			feeAndTax := result.TransactionFee + result.Tax
			if roundToTwoDecimals(feeAndTax) != roundToTwoDecimals(result.TotalDeduction) {
				t.Errorf("fee + tax mismatch: expected %.2f, got %.2f", result.TotalDeduction, feeAndTax)
			}

			fmt.Printf("CalculateGrossAmount: Desired Net: %.2f, Payment Method: %s, Gross: %.2f, Fee: %.2f, Tax: %.2f, Actual Net: %.2f\n",
				tc.desiredNetAmount, tc.paymentMethod, result.GrossAmount, result.TransactionFee, result.Tax, result.NetAmount)
		})
	}
}

func TestCalculateGrossAmount_RoundTrip(t *testing.T) {
	usecase := NewDokuSettlementUseCase()

	// Test that CalculateGrossAmount produces a gross that gives back at least the desired net
	testCases := []struct {
		paymentMethod    string
		desiredNetAmount float64
	}{
		{constants.VIRTUAL_ACCOUNT, 100000},
		{constants.BCA_VA, 100000},
		{constants.SHOPEEPAY, 100000},
		{constants.CREDIT_CARD, 100000},
		{constants.QRIS, 100000},
		{constants.ALFA_GROUP, 100000},
		{constants.PAYLATER_KREDIVO, 100000},
	}

	for _, tc := range testCases {
		t.Run(tc.paymentMethod, func(t *testing.T) {
			// First, calculate the gross amount needed
			grossResult, err := usecase.CalculateGrossAmount(tc.paymentMethod, tc.desiredNetAmount)
			if err != nil {
				t.Fatalf("unexpected error calculating gross: %v", err)
			}

			// Then, verify by calculating settlement with that gross
			settlementResult, err := usecase.CalculateSettlementFee(tc.paymentMethod, grossResult.GrossAmount)
			if err != nil {
				t.Fatalf("unexpected error calculating settlement: %v", err)
			}

			// The net amount from settlement should be >= desired net
			if settlementResult.NetAmount < tc.desiredNetAmount {
				t.Errorf("round trip failed: desired net %.2f, got %.2f", tc.desiredNetAmount, settlementResult.NetAmount)
			}

			// Both calculations should produce the same results
			if grossResult.GrossAmount != settlementResult.GrossAmount {
				t.Errorf("gross amount mismatch: %.2f vs %.2f", grossResult.GrossAmount, settlementResult.GrossAmount)
			}

			if grossResult.TransactionFee != settlementResult.TransactionFee {
				t.Errorf("transaction fee mismatch: %.2f vs %.2f", grossResult.TransactionFee, settlementResult.TransactionFee)
			}

			if grossResult.Tax != settlementResult.Tax {
				t.Errorf("tax mismatch: %.2f vs %.2f", grossResult.Tax, settlementResult.Tax)
			}

			if grossResult.NetAmount != settlementResult.NetAmount {
				t.Errorf("net amount mismatch: %.2f vs %.2f", grossResult.NetAmount, settlementResult.NetAmount)
			}

			fmt.Printf("Round trip OK: Desired Net: %.2f, Gross: %.2f, Actual Net: %.2f\n",
				tc.desiredNetAmount, grossResult.GrossAmount, settlementResult.NetAmount)
		})
	}
}

func TestCalculateGrossAmount_ErrorCases(t *testing.T) {
	usecase := NewDokuSettlementUseCase()

	t.Run("Empty payment method", func(t *testing.T) {
		_, err := usecase.CalculateGrossAmount("", 100000)
		if err == nil {
			t.Error("expected error for empty payment method")
		}
	})

	t.Run("Zero amount", func(t *testing.T) {
		_, err := usecase.CalculateGrossAmount(constants.BCA_VA, 0)
		if err == nil {
			t.Error("expected error for zero amount")
		}
	})

	t.Run("Negative amount", func(t *testing.T) {
		_, err := usecase.CalculateGrossAmount(constants.BCA_VA, -100000)
		if err == nil {
			t.Error("expected error for negative amount")
		}
	})

	t.Run("Unknown payment method", func(t *testing.T) {
		_, err := usecase.CalculateGrossAmount("UNKNOWN_PAYMENT_METHOD", 100000)
		if err == nil {
			t.Error("expected error for unknown payment method")
		}
	})
}

func TestCalculateGrossAmount_LargeAmounts(t *testing.T) {
	usecase := NewDokuSettlementUseCase()

	// Test with a large amount (10 million IDR desired net)
	result, err := usecase.CalculateGrossAmount(constants.SHOPEEPAY, 10000000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify merchant receives at least 10 million
	if result.NetAmount < 10000000 {
		t.Errorf("expected net amount >= 10000000, got %.2f", result.NetAmount)
	}

	// Verify calculation by reverse
	settlementResult, err := usecase.CalculateSettlementFee(constants.SHOPEEPAY, result.GrossAmount)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if settlementResult.NetAmount < 10000000 {
		t.Errorf("round trip: expected net amount >= 10000000, got %.2f", settlementResult.NetAmount)
	}

	fmt.Printf("Large amount test: Desired Net: 10000000, Gross: %.2f, Actual Net: %.2f\n",
		result.GrossAmount, result.NetAmount)
}

func TestCalculateGrossAmount_SpecificValues(t *testing.T) {
	usecase := NewDokuSettlementUseCase()

	// Test specific expected values for Virtual Account (flat fee IDR 4000 + 11% tax)
	// desiredNet = 100000
	// grossAmount = netAmount + flatFee * (1 + taxRate)
	// grossAmount = 100000 + 4000 * 1.11 = 100000 + 4440 = 104440
	result, err := usecase.CalculateGrossAmount(constants.VIRTUAL_ACCOUNT, 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedGross := float64(104440)
	if result.GrossAmount != expectedGross {
		t.Errorf("VIRTUAL_ACCOUNT: expected gross %.2f, got %.2f", expectedGross, result.GrossAmount)
	}

	// Test specific expected values for BCA VA (flat fee IDR 4000 + 11% tax)
	// Same calculation as VIRTUAL_ACCOUNT since pricing is the same
	result, err = usecase.CalculateGrossAmount(constants.BCA_VA, 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedGross = float64(104440)
	if result.GrossAmount != expectedGross {
		t.Errorf("BCA VA: expected gross %.2f, got %.2f", expectedGross, result.GrossAmount)
	}

	// For QRIS (flat fee IDR 700, no tax)
	// grossAmount = netAmount + flatFee = 100000 + 700 = 100700
	result, err = usecase.CalculateGrossAmount(constants.QRIS, 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedGross = float64(100700)
	if result.GrossAmount != expectedGross {
		t.Errorf("QRIS: expected gross %.2f, got %.2f", expectedGross, result.GrossAmount)
	}

	fmt.Printf("Specific values test passed\n")
}

func TestCalculateSettlementFee_VirtualAccount(t *testing.T) {
	usecase := NewDokuSettlementUseCase()

	// Test that VIRTUAL_ACCOUNT constant works the same as specific VA constants
	result, err := usecase.CalculateSettlementFee(constants.VIRTUAL_ACCOUNT, 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Virtual Account: Flat fee IDR 4000 + 11% tax
	// Fee: 4000, Tax: 440, Net: 95560
	expectedFee := float64(4000)
	expectedTax := float64(440)
	expectedNet := float64(95560)

	if result.TransactionFee != expectedFee {
		t.Errorf("VIRTUAL_ACCOUNT: expected transaction fee %.2f, got %.2f", expectedFee, result.TransactionFee)
	}

	if result.Tax != expectedTax {
		t.Errorf("VIRTUAL_ACCOUNT: expected tax %.2f, got %.2f", expectedTax, result.Tax)
	}

	if result.NetAmount != expectedNet {
		t.Errorf("VIRTUAL_ACCOUNT: expected net amount %.2f, got %.2f", expectedNet, result.NetAmount)
	}

	// Verify it matches BCA_VA (same pricing)
	bcaResult, err := usecase.CalculateSettlementFee(constants.BCA_VA, 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.TransactionFee != bcaResult.TransactionFee {
		t.Errorf("VIRTUAL_ACCOUNT vs BCA_VA fee mismatch: %.2f vs %.2f", result.TransactionFee, bcaResult.TransactionFee)
	}

	if result.Tax != bcaResult.Tax {
		t.Errorf("VIRTUAL_ACCOUNT vs BCA_VA tax mismatch: %.2f vs %.2f", result.Tax, bcaResult.Tax)
	}

	if result.NetAmount != bcaResult.NetAmount {
		t.Errorf("VIRTUAL_ACCOUNT vs BCA_VA net mismatch: %.2f vs %.2f", result.NetAmount, bcaResult.NetAmount)
	}

	fmt.Printf("VIRTUAL_ACCOUNT test passed: Fee: %.2f, Tax: %.2f, Net: %.2f\n", result.TransactionFee, result.Tax, result.NetAmount)
}

func TestCalculateGrossAmount_VirtualAccount(t *testing.T) {
	usecase := NewDokuSettlementUseCase()

	// Test that VIRTUAL_ACCOUNT constant works the same as specific VA constants
	result, err := usecase.CalculateGrossAmount(constants.VIRTUAL_ACCOUNT, 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Virtual Account: Flat fee IDR 4000 + 11% tax
	// grossAmount = netAmount + flatFee * (1 + taxRate)
	// grossAmount = 100000 + 4000 * 1.11 = 104440
	expectedGross := float64(104440)
	expectedNet := float64(100000)

	if result.GrossAmount != expectedGross {
		t.Errorf("VIRTUAL_ACCOUNT: expected gross %.2f, got %.2f", expectedGross, result.GrossAmount)
	}

	if result.NetAmount < expectedNet {
		t.Errorf("VIRTUAL_ACCOUNT: expected net >= %.2f, got %.2f", expectedNet, result.NetAmount)
	}

	// Verify it matches BCA_VA (same pricing)
	bcaResult, err := usecase.CalculateGrossAmount(constants.BCA_VA, 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.GrossAmount != bcaResult.GrossAmount {
		t.Errorf("VIRTUAL_ACCOUNT vs BCA_VA gross mismatch: %.2f vs %.2f", result.GrossAmount, bcaResult.GrossAmount)
	}

	if result.TransactionFee != bcaResult.TransactionFee {
		t.Errorf("VIRTUAL_ACCOUNT vs BCA_VA fee mismatch: %.2f vs %.2f", result.TransactionFee, bcaResult.TransactionFee)
	}

	if result.Tax != bcaResult.Tax {
		t.Errorf("VIRTUAL_ACCOUNT vs BCA_VA tax mismatch: %.2f vs %.2f", result.Tax, bcaResult.Tax)
	}

	fmt.Printf("VIRTUAL_ACCOUNT gross amount test passed: Gross: %.2f, Fee: %.2f, Tax: %.2f, Net: %.2f\n",
		result.GrossAmount, result.TransactionFee, result.Tax, result.NetAmount)
}
