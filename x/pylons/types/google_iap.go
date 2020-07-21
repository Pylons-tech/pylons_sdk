package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GoogleIAPOrder is a struct that contains all the metadata of a google iap order
type GoogleIAPOrder struct {
	ProductID     string
	PurchaseToken string
	ReceiptData   string
	Signature     string
	Sender        sdk.AccAddress
}

// NewGoogleIAPOrder return a new Google IAP Order
func NewGoogleIAPOrder(ProductID, PurchaseToken, ReceiptData, Signature string, Sender sdk.AccAddress) GoogleIAPOrder {
	cb := GoogleIAPOrder{
		ProductID:     ProductID,
		PurchaseToken: PurchaseToken,
		ReceiptData:   ReceiptData,
		Signature:     Signature,
		Sender:        Sender,
	}

	return cb
}
