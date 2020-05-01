package intTest

func GetRecipeGUIDFromName(name string, account string) (string, error) {
	rcpList, err := ListRecipesViaCLI(account)
	if err != nil {
		return "", err
	}
	rcp, _ := FindRecipeFromArrayByName(rcpList, name)
	return rcp.ID, nil
}
