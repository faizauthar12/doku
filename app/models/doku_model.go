package models

import (
	"time"

	"github.com/guregu/null/v6"
)

type DokuSignatureComponent struct {
	ClientId         null.String `json:"Client-Id"`
	RequestId        null.String `json:"Request-Id"`
	RequestTimestamp null.String `json:"Request-Timestamp"`
	RequestTarget    null.String `json:"Request-Target"`
	Digest           null.String `json:"Digest"`
}

type DokuOrder struct {
	InvoiceNumber null.String `json:"invoice_number"`
	Amount        null.Int    `json:"amount"`
	Currency      null.String `json:"currency,omitempty"` // for response object
	SessionID     null.String `json:"session_id,omitempty"`
}

type DokuVirtualAccountInfo struct {
	ExpiredTime    null.Int    `json:"expired_time"`
	ReusableStatus null.Bool   `json:"reusable_status"`
	Info1          null.String `json:"info1,omitempty"`
	Info2          null.String `json:"info2,omitempty"`
	Info3          null.String `json:"info3,omitempty"`
}

type DokuCustomer struct {
	ID    null.String `json:"id,omitempty"`
	Name  null.String `json:"name"`
	Email null.String `json:"email"`
}

// will be used as SubAccount ID.
type DokuAdditionalInfo struct {
	Account DokuAccount `json:"account"`
}

// SubAccount ID
type DokuAccount struct {
	ID null.String `json:"id"`
}

type DokuPayment struct {
	PaymentMethodTypes []string    `json:"payment_method_types,omitempty"`
	PaymentDueDate     null.Int    `json:"payment_due_date,omitempty"`
	TokenID            null.String `json:"token_id,omitempty"`
	URL                null.String `json:"url,omitempty"`
	ExpiredDate        null.String `json:"expired_date,omitempty"`
}

type DokuPaymentAdditionalInfo struct {
	Origin struct {
		Product   null.String `json:"product,omitempty"`
		System    null.String `json:"system,omitempty"`
		ApiFormat null.String `json:"apiFormat,omitempty"`
		Source    null.String `json:"source,omitempty"`
	}
}

type DokuHeader struct {
	RequestID null.String `json:"request_id"`
	Signature null.String `json:"signature"`
	Date      *time.Time  `json:"date,omitempty"`
	ClientID  null.String `json:"client_id"`
}

type DokuBalance struct {
	Pending   null.String `json:"pending"`
	Available null.String `json:"available"`
}

type DokuTransaction struct {
	Type              null.String `json:"type,omitempty"`
	Status            null.String `json:"status"`
	Date              *time.Time  `json:"date"`
	OriginalRequestID null.String `json:"original_request_id"`
}

type DokuPaymentIdentifier struct {
	Name  null.String `json:"name"`
	Value null.String `json:"value"`
}
type DokuVirtualAccountpayment struct {
	ReferenceNumber null.String              `json:"reference_number"`
	Date            *time.Time               `json:"date"`
	Identifier      []*DokuPaymentIdentifier `json:"identifier"`
}

type DokuCardPayment struct {
	MaskedCardNumber null.String `json:"masked_card_number"`
	ApprovalCode     null.String `json:"approval_code"`
	ResponseCode     null.String `json:"response_code"`
	ResponseMessage  null.String `json:"response_message"`
	Issuer           null.String `json:"issuer"`
	PaymentID        null.String `json:"payment_id"`
}

type DokuSettlement struct {
	BankAccountSettlementID null.String `json:"bank_account_settlement_id"`
	Value                   null.Float  `json:"value"`
	Type                    null.String `json:"type"`
}
