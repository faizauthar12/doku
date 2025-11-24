package responses

import (
	"time"

	"github.com/faizauthar12/doku/app/models"

	"github.com/guregu/null/v6"
)

type DokuErrorHTTPResponse struct {
	Message []string `json:"message"`
}

type DokuCreateSubAccountHTTPResponse struct {
	Account *DokuCreateSubAccountAccountResponse `json:"account"`
}

type DokuCreateSubAccountAccountResponse struct {
	CreatedDate *time.Time  `json:"created_date"`
	UpdatedDate *time.Time  `json:"updated_date"`
	Name        null.String `json:"name"`
	Type        null.String `json:"type"`
	Status      null.String `json:"status"`
	ID          null.String `json:"id"`
}

type DokuCreatePaymentHTTPResponse struct {
	Response struct {
		Order          *models.DokuOrder                 `json:"order"`
		Payment        *models.DokuPayment               `json:"payment"`
		AdditionalInfo *models.DokuPaymentAdditionalInfo `json:"additional_info"`
		//UUID           string                            `json:"uuid"`
		Headers *models.DokuHeader `json:"headers"`
	} `json:"response"`
}

type DokuGetBalanceHTTPResponse struct {
	Balance *models.DokuBalance `json:"balance"`
}

type DokuPostNotificationHTTPResponse struct {
	Service struct {
		ID null.String `json:"id"`
	} `json:"service"`
	Acquirer struct {
		ID null.String `json:"id"`
	} `json:"acquirer"`
	Channel struct {
		ID null.String `json:"id"`
	} `json:"channel"`
	Transaction *models.DokuTransaction `json:"transaction"`
	Order       *models.DokuOrder       `json:"order"`
	Customer    *models.DokuCustomer    `json:"customer"`

	VirtualAccountInfo struct {
		VirtualAccountNumber null.String `json:"virtual_account_number"`
	} `json:"virtual_account_info,omitempty"`
	VirtualAccountPayment *models.DokuVirtualAccountpayment `json:"virtual_account_payment,omitempty"`

	// Credit Card
	CardPayment  *models.DokuCardPayment `json:"card_payment,omitempty"`
	AuthorizedID null.String             `json:"authorized_id,omitempty"`

	// Convenience Store
	OnlineToOfflineInfo struct {
		PaymentCode null.String `json:"payment_code"`
	} `json:"online_to_offline_info,omitempty"`
	OnlineToOfflinePayment struct {
		Identifier []*models.DokuPaymentIdentifier `json:"identifier"`
	} `json:"online_to_offline_payment,,omitempty"`

	// E-Wallet - Shopeepay
	ShopeepayConfiguration struct {
		MerchantExtID null.String `json:"merchant_ext_id"`
		StoreExtID    null.String `json:"store_ext_id"`
	} `json:"shopeepay_configuration,omitempty"`
	ShopeepayPayment struct {
		TransactionStatus  null.String                     `json:"transaction_status"`
		TransactionMessage null.String                     `json:"transaction_message"`
		Identifier         []*models.DokuPaymentIdentifier `json:"identifier"`
	} `json:"shopeepay_payment,omitempty"`

	// E-Wallet - OVO
	Wallet struct {
		Issuer            null.String `json:"issuer"`
		TokenID           null.String `json:"token_id"`
		MaskedPhoneNumber null.String `json:"masked_phone_number"`
		Status            null.String `json:"status"`
	} `json:"wallet,omitempty"`

	// Paylater
	PeerToPeerInfo struct {
		VirtualAccountNumber    string                          `json:"virtual_account_number"`
		CreatedDate             string                          `json:"created_date"`
		ExpiredDate             string                          `json:"expired_date"`
		Status                  string                          `json:"status"`
		MerchantUniqueReference string                          `json:"merchant_unique_reference"`
		Identifier              []*models.DokuPaymentIdentifier `json:"identifier"`
	} `json:"peer_to_peer_info,omitempty"`
	Payment struct {
		MerchantUniqueReference string `json:"merchant_unique_reference"`
	} `json:"payment,omitempty"`

	// Qris
	EmoneyPayment struct {
		AccountID    null.String `json:"account_id"`
		ApprovalCode null.String `json:"approval_code"`
	} `json:"emoney_payment,omitempty"`
	AditionalInfo struct {
		PostalCode null.String `json:"postalCode,omitempty"`
		FeeType    null.String `json:"feeType,omitempty"`
	} `json:"aditional_info,omitempty"`
	Settlement []*models.DokuSettlement `json:"settlement,omitempty"`
	Origin     struct {
		Product   null.String `json:"product,omitempty"`
		System    null.String `json:"system,omitempty"`
		ApiFormat null.String `json:"apiFormat,omitempty"`
		Source    null.String `json:"source,omitempty"`
	}
}

type GetTokenResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	AccessToken     string `json:"accessToken"`
	TokenType       string `json:"tokenType"`
	ExpiresIn       int    `json:"expiresIn"`
	AdditionalInfo  string `json:"additionalInfo"`
}

type BankAccountInquiryResponse struct {
	ResponseCode             int    `json:"responseCode"`
	ResponseMessage          string `json:"responseMessage"`
	ReferenceNo              string `json:"referenceNo"`
	PartnerReferenceNo       string `json:"partnerReferenceNo"`
	BeneficiaryAccountNumber string `json:"beneficiaryAccountNumber"`
	BeneficiaryAccountName   string `json:"beneficiaryAccountName"`
	BeneficiaryBankCode      string `json:"beneficiaryBankCode"`
	BeneficiaryBankShortName string `json:"beneficiaryBankShortName"`
	BeneficiaryBankName      string `json:"beneficiaryBankName"`
	Amount                   struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	SessionID      string `json:"sessionId"`
	AdditionalInfo struct {
		SenderCountryCode   string `json:"senderCountryCode"`
		ForexRate           string `json:"forexRate"`
		ForexOriginCurrency string `json:"forexOriginCurrency"`
		FeeAmount           string `json:"feeAmount"`
		FeeCurrency         string `json:"feeCurrency"`
		BeneficiaryAmount   string `json:"beneficiaryAmount"`
		ReferenceNumber     string `json:"referenceNumber"`
	} `json:"additionalInfo"`
}
