package model

import "time"

type SCBPayment struct {
	Payeeproxyid           string    `json:"payeeProxyId" bson:"payeeProxyId"`
	Payeeproxytype         string    `json:"payeeProxyType" bson:"payeeProxyType"`
	Payeeaccountnumber     string    `json:"payeeAccountNumber" bson:"payeeAccountNumber"`
	Payeename              string    `json:"payeeName" bson:"payeeName"`
	Payerproxyid           string    `json:"payerProxyId" bson:"payerProxyId"`
	Payerproxytype         string    `json:"payerProxyType" bson:"payerProxyType"`
	Payeraccountnumber     string    `json:"payerAccountNumber" bson:"payerAccountNumber"`
	Payername              string    `json:"payerName" bson:"payerName"`
	Sendingbankcode        string    `json:"sendingBankCode" bson:"sendingBankCode"`
	Receivingbankcode      string    `json:"receivingBankCode" bson:"receivingBankCode"`
	Amount                 string    `json:"amount" bson:"amount"`
	Channelcode            string    `json:"channelCode" bson:"channelCode"`
	Transactionid          string    `json:"transactionId" bson:"transactionId"`
	Transactiondateandtime time.Time `json:"transactionDateandTime" bson:"transactionDateandTime"`
	Billpaymentref1        string    `json:"billPaymentRef1" bson:"billPaymentRef1" binding:"required"`
	Billpaymentref2        string    `json:"billPaymentRef2" bson:"billPaymentRef2" binding:"required"`
	Billpaymentref3        string    `json:"billPaymentRef3" bson:"billPaymentRef3"`
	Currencycode           string    `json:"currencyCode" bson:"currencyCode"`
	Transactiontype        string    `json:"transactionType" bson:"transactionType"`
	CreatedAt              time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" bson:"updated_at"`
}
