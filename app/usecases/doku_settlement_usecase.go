package usecases

import (
	"errors"
	"math"

	"github.com/faizauthar12/doku/app/config"
	"github.com/faizauthar12/doku/app/constants"
	"github.com/faizauthar12/doku/app/responses"
)

type DokuSettlementUseCaseInterface interface {
	CalculateSettlementFee(paymentMethod string, amount float64) (*responses.DokuSettlementResultResponse, error)
	CalculateGrossAmount(paymentMethod string, desiredNetAmount float64) (*responses.DokuSettlementResultResponse, error)
}

type dokuSettlementUseCase struct {
	cfg config.Configuration
}

func NewDokuSettlementUseCase() DokuSettlementUseCaseInterface {
	return &dokuSettlementUseCase{
		cfg: config.Get(),
	}
}

func (u *dokuSettlementUseCase) calculateFeeAndTax(paymentMethod string, amount float64) (transactionFee float64, tax float64, err error) {
	taxRate := float64(u.cfg.TransactionFee.Tax) / 100

	switch paymentMethod {
	// Cards
	case constants.CREDIT_CARD, constants.KKI:
		percentageRate := u.cfg.TransactionFee.Cards.PercentageRate / 100
		flatFee := float64(u.cfg.TransactionFee.Cards.FlatFee)
		transactionFee = (amount * percentageRate) + flatFee
		tax = transactionFee * taxRate

	// Virtual Account (Transfer Bank)
	case constants.VIRTUAL_ACCOUNT, constants.BCA_VA, constants.Mandiri_VA, constants.BSI_VA, constants.BRI_VA,
		constants.BNI_VA, constants.DOKU_VA, constants.PERMATA_VA, constants.CIMB_VA,
		constants.DANAMON_VA, constants.BTN_VA, constants.BNC_VA:
		transactionFee = float64(u.cfg.TransactionFee.VirtualAccount.FlatFee)
		tax = transactionFee * taxRate

	// Convenience Store - Alfamart
	case constants.ALFA_GROUP:
		transactionFee = float64(u.cfg.TransactionFee.ConvenienceStore.Alfamart.FlatFee)
		tax = transactionFee * taxRate

	// Convenience Store - Indomaret
	case constants.INDOMARET:
		transactionFee = float64(u.cfg.TransactionFee.ConvenienceStore.Indomaret.FlatFee)
		tax = transactionFee * taxRate

	// QRIS - No tax
	case constants.QRIS:
		transactionFee = float64(u.cfg.TransactionFee.QR.FlatFee)
		tax = 0

	// E-Wallet - ShopeePay
	case constants.SHOPEEPAY:
		percentageRate := u.cfg.TransactionFee.EWallet.ShopeePay.PercentageRate / 100
		transactionFee = amount * percentageRate
		tax = transactionFee * taxRate

	// E-Wallet - OVO
	case constants.OVO:
		percentageRate := u.cfg.TransactionFee.EWallet.OVO.PercentageRate / 100
		transactionFee = amount * percentageRate
		tax = transactionFee * taxRate

	// E-Wallet - LinkAja
	case constants.LINKAJA:
		percentageRate := u.cfg.TransactionFee.EWallet.LinkAja.PercentageRate / 100
		transactionFee = amount * percentageRate
		tax = transactionFee * taxRate

	// E-Wallet - DOKU Wallet
	case constants.DOKU_WALLET:
		percentageRate := u.cfg.TransactionFee.EWallet.Doku.PercentageRate / 100
		transactionFee = amount * percentageRate
		tax = transactionFee * taxRate

	// E-Wallet - DANA
	case constants.DANA:
		percentageRate := u.cfg.TransactionFee.EWallet.Dana.PercentageRate / 100
		transactionFee = amount * percentageRate
		tax = transactionFee * taxRate

	// Direct Debit - BRI
	case constants.DIRECT_DEBIT_BRI:
		percentageRate := u.cfg.TransactionFee.DirectDebit.BRI.PercentageRate / 100
		transactionFee = amount * percentageRate
		tax = transactionFee * taxRate

	// Digital Banking - Jenius Pay
	case constants.JENIUS_PAY:
		percentageRate := u.cfg.TransactionFee.DigitalBanking.JeniusPay.PercentageRate / 100
		transactionFee = amount * percentageRate
		tax = transactionFee * taxRate

	// PayLater - Akulaku
	case constants.PAYLATER_AKULAKU:
		percentageRate := u.cfg.TransactionFee.PayLater.Akulaku.PercentageRate / 100
		transactionFee = amount * percentageRate
		tax = transactionFee * taxRate

	// PayLater - Kredivo
	case constants.PAYLATER_KREDIVO:
		percentageRate := u.cfg.TransactionFee.PayLater.Kredivo.PercentageRate / 100
		transactionFee = amount * percentageRate
		tax = transactionFee * taxRate

	// PayLater - Indodana
	case constants.PAYLATER_INDODANA:
		percentageRate := u.cfg.TransactionFee.PayLater.Indodana.PercentageRate / 100
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

// roundUpToWholeNumber rounds up a value to the nearest whole number (ceiling)
// This ensures the merchant always receives at least the desired net amount
func roundUpToWholeNumber(value float64) float64 {
	return math.Ceil(value)
}

// getFeeParameters returns the percentage rate (as decimal), flat fee, and whether tax applies
// for a given payment method
func (u *dokuSettlementUseCase) getFeeParameters(paymentMethod string) (percentageRate float64, flatFee float64, hasTax bool, err error) {
	switch paymentMethod {
	// Cards: Percentage + Flat Fee
	case constants.CREDIT_CARD, constants.KKI:
		percentageRate = u.cfg.TransactionFee.Cards.PercentageRate / 100
		flatFee = float64(u.cfg.TransactionFee.Cards.FlatFee)
		hasTax = true

	// Virtual Account: Flat Fee Only
	case constants.VIRTUAL_ACCOUNT, constants.BCA_VA, constants.Mandiri_VA, constants.BSI_VA, constants.BRI_VA,
		constants.BNI_VA, constants.DOKU_VA, constants.PERMATA_VA, constants.CIMB_VA,
		constants.DANAMON_VA, constants.BTN_VA, constants.BNC_VA:
		percentageRate = 0
		flatFee = float64(u.cfg.TransactionFee.VirtualAccount.FlatFee)
		hasTax = true

	// Convenience Store - Alfamart
	case constants.ALFA_GROUP:
		percentageRate = 0
		flatFee = float64(u.cfg.TransactionFee.ConvenienceStore.Alfamart.FlatFee)
		hasTax = true

	// Convenience Store - Indomaret
	case constants.INDOMARET:
		percentageRate = 0
		flatFee = float64(u.cfg.TransactionFee.ConvenienceStore.Indomaret.FlatFee)
		hasTax = true

	// QRIS - No tax
	case constants.QRIS:
		percentageRate = 0
		flatFee = float64(u.cfg.TransactionFee.QR.FlatFee)
		hasTax = false

	// E-Wallet - ShopeePay
	case constants.SHOPEEPAY:
		percentageRate = u.cfg.TransactionFee.EWallet.ShopeePay.PercentageRate / 100
		flatFee = 0
		hasTax = true

	// E-Wallet - OVO
	case constants.OVO:
		percentageRate = u.cfg.TransactionFee.EWallet.OVO.PercentageRate / 100
		flatFee = 0
		hasTax = true

	// E-Wallet - LinkAja
	case constants.LINKAJA:
		percentageRate = u.cfg.TransactionFee.EWallet.LinkAja.PercentageRate / 100
		flatFee = 0
		hasTax = true

	// E-Wallet - DOKU Wallet
	case constants.DOKU_WALLET:
		percentageRate = u.cfg.TransactionFee.EWallet.Doku.PercentageRate / 100
		flatFee = 0
		hasTax = true

	// E-Wallet - DANA
	case constants.DANA:
		percentageRate = u.cfg.TransactionFee.EWallet.Dana.PercentageRate / 100
		flatFee = 0
		hasTax = true

	// Direct Debit - BRI
	case constants.DIRECT_DEBIT_BRI:
		percentageRate = u.cfg.TransactionFee.DirectDebit.BRI.PercentageRate / 100
		flatFee = 0
		hasTax = true

	// Digital Banking - Jenius Pay
	case constants.JENIUS_PAY:
		percentageRate = u.cfg.TransactionFee.DigitalBanking.JeniusPay.PercentageRate / 100
		flatFee = 0
		hasTax = true

	// PayLater - Akulaku
	case constants.PAYLATER_AKULAKU:
		percentageRate = u.cfg.TransactionFee.PayLater.Akulaku.PercentageRate / 100
		flatFee = 0
		hasTax = true

	// PayLater - Kredivo
	case constants.PAYLATER_KREDIVO:
		percentageRate = u.cfg.TransactionFee.PayLater.Kredivo.PercentageRate / 100
		flatFee = 0
		hasTax = true

	// PayLater - Indodana
	case constants.PAYLATER_INDODANA:
		percentageRate = u.cfg.TransactionFee.PayLater.Indodana.PercentageRate / 100
		flatFee = 0
		hasTax = true

	default:
		return 0, 0, false, errors.New("unknown payment method: " + paymentMethod)
	}

	return percentageRate, flatFee, hasTax, nil
}

func (u *dokuSettlementUseCase) CalculateSettlementFee(paymentMethod string, amount float64) (*responses.DokuSettlementResultResponse, error) {
	if paymentMethod == "" {
		return nil, errors.New("payment method is empty")
	}

	if amount <= 0 {
		return nil, errors.New("invalid amount: must be greater than 0")
	}

	transactionFee, tax, err := u.calculateFeeAndTax(paymentMethod, amount)
	if err != nil {
		return nil, err
	}

	totalDeduction := transactionFee + tax
	netAmount := amount - totalDeduction

	return &responses.DokuSettlementResultResponse{
		PaymentMethod:  paymentMethod,
		GrossAmount:    amount,
		TransactionFee: roundToTwoDecimals(transactionFee),
		Tax:            roundToTwoDecimals(tax),
		TotalDeduction: roundToTwoDecimals(totalDeduction),
		NetAmount:      roundToTwoDecimals(netAmount),
	}, nil
}

// CalculateGrossAmount calculates the gross amount (what customer pays) given a desired net amount
// that the merchant wants to receive after all fees and taxes are deducted.
//
// Formula derivation:
// For percentage-based fees with tax:
//
//	netAmount = grossAmount - fee - tax
//	netAmount = grossAmount - (grossAmount * rate) - (grossAmount * rate * taxRate)
//	netAmount = grossAmount * (1 - rate - rate * taxRate)
//	netAmount = grossAmount * (1 - rate * (1 + taxRate))
//	grossAmount = netAmount / (1 - rate * (1 + taxRate))
//
// For flat fee with tax:
//
//	netAmount = grossAmount - flatFee - (flatFee * taxRate)
//	grossAmount = netAmount + flatFee * (1 + taxRate)
//
// For mixed (percentage + flat) with tax:
//
//	grossAmount = (netAmount + flatFee * (1 + taxRate)) / (1 - rate * (1 + taxRate))
func (u *dokuSettlementUseCase) CalculateGrossAmount(paymentMethod string, desiredNetAmount float64) (*responses.DokuSettlementResultResponse, error) {
	if paymentMethod == "" {
		return nil, errors.New("payment method is empty")
	}

	if desiredNetAmount <= 0 {
		return nil, errors.New("invalid desired net amount: must be greater than 0")
	}

	// Get fee parameters for this payment method
	percentageRate, flatFee, hasTax, err := u.getFeeParameters(paymentMethod)
	if err != nil {
		return nil, err
	}

	taxRate := float64(0)
	if hasTax {
		taxRate = float64(u.cfg.TransactionFee.Tax) / 100
	}

	// Calculate gross amount using the inverse formula
	var grossAmount float64

	// taxMultiplier = (1 + taxRate)
	taxMultiplier := 1 + taxRate

	// For percentage-based: divisor = 1 - percentageRate * (1 + taxRate)
	divisor := 1 - (percentageRate * taxMultiplier)

	if divisor <= 0 {
		return nil, errors.New("invalid fee configuration: fees exceed 100%")
	}

	// grossAmount = (netAmount + flatFee * (1 + taxRate)) / (1 - rate * (1 + taxRate))
	grossAmount = (desiredNetAmount + (flatFee * taxMultiplier)) / divisor

	// Round up to ensure merchant receives at least the desired net amount
	grossAmount = roundUpToWholeNumber(grossAmount)

	// Now verify by calculating the actual settlement
	transactionFee, tax, err := u.calculateFeeAndTax(paymentMethod, grossAmount)
	if err != nil {
		return nil, err
	}

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
