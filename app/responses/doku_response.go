package responses

import (
	"doku/app/models"
	"time"
)

type DokuCreateSubAccountHTTPResponse struct {
	Account *DokuCreateSubAccountAccountResponse `json:"account"`
}

type DokuCreateSubAccountAccountResponse struct {
	CreatedDate *time.Time `json:"created_date"`
	UpdatedDate *time.Time `json:"updated_date"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	ID          string     `json:"id"`
}

type DokuCreatePaymentHTTPResponse struct {
	Response struct {
		Order          *models.DokuOrder                 `json:"order"`
		Payment        *models.DokuPayment               `json:"payment"`
		AdditionalInfo *models.DokuPaymentAdditionalInfo `json:"additional_info"`
		UUID           int64                             `json:"uuid"`
		Headers        *models.DokuHeader                `json:"headers"`
	} `json:"response"`
}

type DokuGetBalanceHTTPResponse struct {
	Balance *models.DokuBalance `json:"balance"`
}
