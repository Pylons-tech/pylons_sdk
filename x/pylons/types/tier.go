package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Tier defines the kind of cookbook this is
const (
	// Basic is the free int64 which does allow developers to use pylons ( paid currency ) in their
	// games
	Basic int64 = iota
	Premium
)

// ValidateLevel validates the level
func ValidateLevel(level int64) error {
	if level == Basic || level == Premium {
		return nil
	}

	return errors.New("Invalid cookbook plan")
}

// tier fee types
var (
	// BasicFee is the fee charged to create a basic cookbook
	BasicFee = NewPylon(10000)
	// PremiumFee is the fee charged to create a premium cookbook
	PremiumFee = NewPylon(50000)
)

// Tier defines the kind of cookbook this is
type Tier struct {
	Level int64
	Fee   sdk.Coins
}

// BasicTier is the cookbook tier which doesn't allow paid recipes which means
// the developers cannot have recipes where they can actually carge a fee in pylons
var BasicTier = Tier{
	Level: Basic,
	Fee:   BasicFee,
}

// PremiumTier the cookbook tier which does allow paid recipes
var PremiumTier = Tier{
	Level: Premium,
	Fee:   PremiumFee,
}
