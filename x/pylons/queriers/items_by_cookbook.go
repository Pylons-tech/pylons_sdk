package queriers

import (
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
)

// query endpoints supported by the nameservice Querier
const (
	KeyItemsByCookbook = "items_by_cookbook"
)

// ItemResp is the response for Items
type ItemResp struct {
	Items []types.Item
}
