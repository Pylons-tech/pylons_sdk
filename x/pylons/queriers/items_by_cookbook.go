package queriers

import (
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
)

// query endpoints supported by the nameservice Querier
const (
	KeyItemsByCookbook = "items_by_cookbook"
)

// ItemResponse is the response for Items
type ItemResponse struct {
	Items []types.Item
}
