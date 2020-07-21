package msgs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgGoogleIAPGetPylons defines a GetPylons message
type MsgGoogleIAPGetPylons struct {
	ProductID         string
	PurchaseToken     string
	ReceiptDataBase64 string
	Signature         string
	Requester         sdk.AccAddress
}

// NewMsgGoogleIAPGetPylons is a function to get MsgGetPylons msg from required params
func NewMsgGoogleIAPGetPylons(ProductID, PurchaseToken, ReceiptDataBase64, Signature string, requester sdk.AccAddress) MsgGoogleIAPGetPylons {
	return MsgGoogleIAPGetPylons{
		ProductID:         ProductID,
		PurchaseToken:     PurchaseToken,
		ReceiptDataBase64: ReceiptDataBase64,
		Signature:         Signature,
		Requester:         requester,
	}
}

// Route should return the name of the module
func (msg MsgGoogleIAPGetPylons) Route() string { return RouterKey }

// Type should return the action
func (msg MsgGoogleIAPGetPylons) Type() string { return "get_pylons" }

// ValidateBasic is a function to validate MsgGoogleIAPGetPylons msg
func (msg MsgGoogleIAPGetPylons) ValidateBasic() error {

	if msg.Requester.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Requester.String())
	}

	var jsonData map[string]interface{}

	receiptData, err := base64.StdEncoding.DecodeString(msg.ReceiptDataBase64)
	if err != nil {
		return err
	}

	err = json.Unmarshal(receiptData, &jsonData)
	if err != nil {
		return err
	}
	if msg.PurchaseToken != jsonData["purchaseToken"] {
		return fmt.Errorf("purchaseToken does not match with receipt data")
	}
	if msg.ProductID != jsonData["productId"] {
		return fmt.Errorf("productId does not match with receipt data")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgGoogleIAPGetPylons) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners encodes the message for signing
func (msg MsgGoogleIAPGetPylons) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Requester}
}
