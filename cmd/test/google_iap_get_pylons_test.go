package inttest

import (
	"encoding/base64"
	"fmt"
	originT "testing"
	"time"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/evtesting"

	inttestSDK "github.com/Pylons-tech/pylons_sdk/cmd/test_utils"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGoogleIAPGetPylonsViaCLI(originT *originT.T) {
	t := testing.NewT(originT)
	t.Parallel()

	tests := []struct {
		name            string
		productID       string
		purchaseToken   string
		receiptData     string
		signature       string
		fromAddress     sdk.AccAddress
		showError       bool
		desiredError    string
		reqAmount       int64
		tryReuseOrderID bool
		tryReuseErr     string
	}{
		{
			name:          "successful check",
			productID:     "pylons_1000",
			purchaseToken: "hafokgmjfkcpdnbffanijckj.AO-J1OxXkrKdM8q14T49Qo5a723VG_8h_4MCY_M2Tqn91L0e7FjiVXsZ2Qxc1SnvoFzHN9jBCJpjZqD4ErYIquMG6Li_jUfcuKuXti_wsa7r48eWNA1Oh0o",
			receiptData:   `{"productId":"pylons_1000","purchaseToken":"hafokgmjfkcpdnbffanijckj.AO-J1OxXkrKdM8q14T49Qo5a723VG_8h_4MCY_M2Tqn91L0e7FjiVXsZ2Qxc1SnvoFzHN9jBCJpjZqD4ErYIquMG6Li_jUfcuKuXti_wsa7r48eWNA1Oh0o","purchaseTime":1595031050407,"developerPayload":null}`,
			// Correct signature
			signature:       "HEo0RYQeH0+8nmYa6ETKP9f3S/W/cUuQTBme7VSh3Lzm+1+1GwJIl1pdF1dh32YGhd3BtyMoLVGzr9ZajfHhhznIvbowS/XIlyJJCE6dI+zg68mKo5rDt0wB2BY8azk0+WCkc5XT5y8biRNXe5RyvmuqYKPXmEsgHaYKo6x3mHs6oXrECckKv/c9T9MHCvdAqVFrml9W7K41sRHbpOdFmYnO33bkNITCCaf/C1PDGMVOItxvq7uXi+F0DpjXwXko9AU6L3pK6zDICcD38HblbzumOg6LGsuWCjOw8QwNobYOUNtrdj01fEXqkKhfYzFZcwxM6xsphN38gnO0ksDdyw==",
			showError:       false,
			desiredError:    "",
			reqAmount:       1000,
			tryReuseOrderID: true,
			tryReuseErr:     "the iap order ID is already being used",
		},
		{
			name:          "productId does not match receipt data",
			productID:     "pylons_55000",
			purchaseToken: "hafokgmjfkcpdnbffanijckj.AO-J1OxXkrKdM8q14T49Qo5a723VG_8h_4MCY_M2Tqn91L0e7FjiVXsZ2Qxc1SnvoFzHN9jBCJpjZqD4ErYIquMG6Li_jUfcuKuXti_wsa7r48eWNA1Oh0o",
			receiptData:   `{"productId":"pylons_1000","purchaseToken":"hafokgmjfkcpdnbffanijckj.AO-J1OxXkrKdM8q14T49Qo5a723VG_8h_4MCY_M2Tqn91L0e7FjiVXsZ2Qxc1SnvoFzHN9jBCJpjZqD4ErYIquMG6Li_jUfcuKuXti_wsa7r48eWNA1Oh0o","purchaseTime":1595031050407,"developerPayload":null}`,
			signature:     "HEo0RYQeH0+8nmYa6ETKP9f3S/W/cUuQTBme7VSh3Lzm+1+1GwJIl1pdF1dh32YGhd3BtyMoLVGzr9ZajfHhhznIvbowS/XIlyJJCE6dI+zg68mKo5rDt0wB2BY8azk0+WCkc5XT5y8biRNXe5RyvmuqYKPXmEsgHaYKo6x3mHs6oXrECckKv/c9T9MHCvdAqVFrml9W7K41sRHbpOdFmYnO33bkNITCCaf/C1PDGMVOItxvq7uXi+F0DpjXwXko9AU6L3pK6zDICcD38HblbzumOg6LGsuWCjOw8QwNobYOUNtrdj01fEXqkKhfYzFZcwxM6xsphN38gnO0ksDdyw==",
			showError:     true,
			desiredError:  "productId does not match with receipt data", // this validation is done before signature verification
		},
		{
			name:          "wrong signature base64 encoding",
			productID:     "pylons_1000",
			purchaseToken: "hafokgmjfkcpdnbffanijckj.AO-J1OxXkrKdM8q14T49Qo5a723VG_8h_4MCY_M2Tqn91L0e7FjiVXsZ2Qxc1SnvoFzHN9jBCJpjZqD4ErYIquMG6Li_jUfcuKuXti_wsa7r48eWNA1Oh0o",
			receiptData:   `{"productId":"pylons_1000","purchaseToken":"hafokgmjfkcpdnbffanijckj.AO-J1OxXkrKdM8q14T49Qo5a723VG_8h_4MCY_M2Tqn91L0e7FjiVXsZ2Qxc1SnvoFzHN9jBCJpjZqD4ErYIquMG6Li_jUfcuKuXti_wsa7r48eWNA1Oh0o","purchaseTime":1595031050407,"developerPayload":null}`,
			// Incorrect signature
			signature:    "Invalid signature", // does not fit the signature base64 format
			showError:    true,
			desiredError: "msg signature base64 decoding failure",
		},
		{
			name:          "wrong signature check",
			productID:     "pylons_55000",
			purchaseToken: "hafokgmjfkcpdnbffanijckj.AO-J1OxXkrKdM8q14T49Qo5a723VG_8h_4MCY_M2Tqn91L0e7FjiVXsZ2Qxc1SnvoFzHN9jBCJpjZqD4ErYIquMG6Li_jUfcuKuXti_wsa7r48eWNA1Oh0o",
			receiptData:   `{"productId":"pylons_55000","purchaseToken":"hafokgmjfkcpdnbffanijckj.AO-J1OxXkrKdM8q14T49Qo5a723VG_8h_4MCY_M2Tqn91L0e7FjiVXsZ2Qxc1SnvoFzHN9jBCJpjZqD4ErYIquMG6Li_jUfcuKuXti_wsa7r48eWNA1Oh0o","purchaseTime":1595031050407,"developerPayload":null}`,
			// Incorrect signature
			signature:    "HEo0RYQeH0+8nmYa6ETKP9f3S/W/cUuQTBme7VSh3Lzm+1+1GwJIl1pdF1dh32YGhd3BtyMoLVGzr9ZajfHhhznIvbowS/XIlyJJCE6dI+zg68mKo5rDt0wB2BY8azk0+WCkc5XT5y8biRNXe5RyvmuqYKPXmEsgHaYKo6x3mHs6oXrECckKv/c9T9MHCvdAqVFrml9W7K41sRHbpOdFmYnO33bkNITCCaf/C1PDGMVOItxvq7uXi+F0DpjXwXko9AU6L3pK6zDICcD38HblbzumOg6LGsuWCjOw8QwNobYOUNtrdj01fEXqkKhfYzFZcwxM6xsphN38gnO0ksDdyw==",
			showError:    true,
			desiredError: "crypto/rsa: verification error",
		},
	}

	for tcNum, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			getPylonsKey := fmt.Sprintf("TestGoogleIAPGetPylonsViaCLI%d_%d", tcNum, time.Now().Unix())
			MockAccount(getPylonsKey, t) // mock account with initial balance

			getPylonsSdkAddr := GetAccountAddress(getPylonsKey, t)
			originBalance := inttestSDK.GetAccountBalanceFromAddr(getPylonsSdkAddr.String(), t)

			receiptDataBase64 := base64.StdEncoding.EncodeToString([]byte(tc.receiptData))

			msgGoogleIAPGetPylons := types.NewMsgGoogleIAPGetPylons(tc.productID, tc.purchaseToken, receiptDataBase64, tc.signature, getPylonsSdkAddr.String())
			txhash, err := inttestSDK.TestTxWithMsgWithNonce(t,
				&msgGoogleIAPGetPylons,
				getPylonsKey,
				false,
			)
			if err != nil {
				TxBroadcastErrorExpected(txhash, err, tc.desiredError, t)
				return
			}

			GetTxHandleResult(txhash, t)
			if tc.showError {
			} else {
				accBalance := inttestSDK.GetAccountBalanceFromAddr(getPylonsSdkAddr.String(), t)
				balanceOk := accBalance.Coins.AmountOf(types.Pylon).Equal(sdk.NewInt(originBalance.Coins.AmountOf(types.Pylon).Int64() + tc.reqAmount))
				t.WithFields(testing.Fields{
					"get_pylons_address": getPylonsSdkAddr.String(),
					"target_increase":    tc.reqAmount,
					"actual_amount":      accBalance.Coins.AmountOf(types.Pylon).Int64(),
					"origin_amount":      originBalance.Coins.AmountOf(types.Pylon).Int64(),
				}).MustTrue(balanceOk, "pylons requestor should get correct revenue")
			}

			if tc.tryReuseOrderID {
				hash, err := inttestSDK.TestTxWithMsgWithNonce(t,
					&msgGoogleIAPGetPylons,
					getPylonsKey,
					false,
				)
				t.MustNil(err, hash)
				errdata, err := inttestSDK.WaitAndGetTxError(hash, inttestSDK.GetMaxWaitBlock(), t)
				t.MustNil(err, hash)
				t.MustContain(string(errdata), tc.tryReuseErr)
			}
		})
	}
}
