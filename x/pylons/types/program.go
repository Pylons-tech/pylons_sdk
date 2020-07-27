package types

import "errors"

// ProgramValidateBasic validate program
func ProgramValidateBasic(program string) error {
	if len(program) == 0 {
		return errors.New("length of program code shouldn't be 0")
	}
	return nil
}
