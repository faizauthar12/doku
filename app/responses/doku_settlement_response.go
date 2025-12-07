package responses

type DokuSettlementResultResponse struct {
	PaymentMethod  string  `json:"payment_method"`
	GrossAmount    float64 `json:"gross_amount"`
	TransactionFee float64 `json:"transaction_fee"`
	Tax            float64 `json:"tax"`
	TotalDeduction float64 `json:"total_deduction"`
	NetAmount      float64 `json:"net_amount"`
}
