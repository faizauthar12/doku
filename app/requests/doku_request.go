package requests

import "github.com/faizauthar12/doku/app/models"

type DokuCreateSubAccountHTTPRequest struct {
	Account DokuCreateSubAccountAccount `json:"account"`
}

type DokuCreateSubAccountAccount struct {
	Email string `json:"email"`
	Type  string `json:"type"` // STANDARD
	Name  string `json:"name"`
}

type DokuCreateSubAccountRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type DokuCreatePaymentHTTPRequest struct {
	Order              *models.DokuOrder              `json:"order"`
	VirtualAccountInfo *models.DokuVirtualAccountInfo `json:"virtual_account_info,omitempty"`
	Customer           *models.DokuCustomer           `json:"customer"`
	AdditionalInfo     *models.DokuAdditionalInfo     `json:"additional_info"`
	Payment            *models.DokuPayment            `json:"payment"`
}

type DokuCreatePaymentRequest struct {
	Amount         int64  `json:"amount"`
	CustomerName   string `json:"customer_name"`
	CustomerEmail  string `json:"customer_email"`
	SacID          string `json:"SacID"`
	PaymentDueDate int64  `json:"payment_due_date,omitempty"`
}

type DokuNotificationRequest struct {
	RequestID        string `json:"Request-Id"`
	RequestTimestamp string `json:"Request-Timestamp"`
	Signature        string `json:"Signature"`
	JsonBody         []byte `json:"Json-Body"`
}
