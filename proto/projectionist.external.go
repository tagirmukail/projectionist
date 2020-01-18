package projectionist

import "fmt"

func (lr *LoginRequest) Validate() error {
	if lr.Username == "" {
		return fmt.Errorf("Username is required")
	}

	if lr.Password == "" {
		return fmt.Errorf("Password is required")
	}

	return nil
}
