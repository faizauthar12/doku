package models

import "time"

type DokuSignatureComponent struct {
	ClientId         string `json:"Client-Id"`
	RequestId        string `json:"Request-Id"`
	RequestTimestamp string `json:"Request-Timestamp"`
	RequestTarget    string `json:"Request-Target"`
	Digest           string `json:"Digest"`
}

type DokuOrder struct {
	InvoiceNumber string `json:"invoice_number"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency,omitempty"` // for response object
	SessionID     string `json:"session_id,omitempty"`
}

type DokuVirtualAccountInfo struct {
	ExpiredTime    int64  `json:"expired_time"`
	ReusableStatus bool   `json:"reusable_status"`
	Info1          string `json:"info1,omitempty"`
	Info2          string `json:"info2,omitempty"`
	Info3          string `json:"info3,omitempty"`
}

type DokuCustomer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// will be used as SubAccount ID.
type DokuAdditionalInfo struct {
	Account DokuAccount `json:"account"`
}

// SubAccount ID
type DokuAccount struct {
	ID string `json:"id"`
}

type DokuPayment struct {
	PaymentMethodTypes []string   `json:"payment_method_types,omitempty"`
	PaymentDueDate     int64      `json:"payment_due_date,omitempty"`
	TokenID            string     `json:"token_id,omitempty"`
	URL                string     `json:"url,omitempty"`
	ExpiredDate        *time.Time `json:"expired_date,omitempty"`
}

type DokuPaymentAdditionalInfo struct {
	Origin struct {
		Product   string `json:"product,omitempty"`
		System    string `json:"system,omitempty"`
		ApiFormat string `json:"apiFormat,omitempty"`
		Source    string `json:"source,omitempty"`
	}
}

type DokuHeader struct {
	RequestID string     `json:"request_id"`
	Signature string     `json:"signature"`
	Date      *time.Time `json:"date,omitempty"`
	ClientID  string     `json:"client_id"`
}

type DokuBalance struct {
	Pending   string `json:"pending"`
	Available string `json:"available"`
}
