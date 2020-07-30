package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GoogleIAPOrder is a struct that contains all the metadata of a google iap order
type GoogleIAPOrder struct {
	NodeVersion   SemVer
	ProductID     string
	PurchaseToken string
	ReceiptData   string
	Signature     string
	Sender        sdk.AccAddress
}

// NewGoogleIAPOrder return a new Google IAP Order
func NewGoogleIAPOrder(ProductID, PurchaseToken, ReceiptData, Signature string, Sender sdk.AccAddress) GoogleIAPOrder {
	cb := GoogleIAPOrder{
		NodeVersion:   SemVer("0.0.1"),
		ProductID:     ProductID,
		PurchaseToken: PurchaseToken,
		ReceiptData:   ReceiptData,
		Signature:     Signature,
		Sender:        Sender,
	}

	return cb
}
