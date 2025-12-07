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
