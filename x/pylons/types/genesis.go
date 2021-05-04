package types

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Cookbooks: []Cookbook{},
		Recipes:   []Recipe{},
		Items:     []Item{},
	}
}

func (gs GenesisState) Validate() error {
	return nil
}
