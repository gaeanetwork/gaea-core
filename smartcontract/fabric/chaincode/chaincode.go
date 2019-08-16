package chaincode

import "fmt"

// CheckArgsEmpty check chaincode args are non-empty
func CheckArgsEmpty(args []string, length int) error {
	if l := len(args); l < length {
		return fmt.Errorf("Incorrect number of arguments. Expecting be greater than or equal to %d, Actual: %d(%v)", length, l, args)
	}

	for index := 0; index < length; index++ {
		if len(args[index]) <= 0 {
			return fmt.Errorf("The index %d argument must be a non-empty string, args: %v", index, args)
		}
	}
	return nil
}
