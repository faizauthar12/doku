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
	case constants.BCA_VA, constants.Mandiri_VA, constants.BSI_VA, constants.BRI_VA,
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
