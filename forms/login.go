package forms

import "fmt"

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (f *LoginForm) Validate() error {
	if f.Username == "" {
		return fmt.Errorf("Username is required")
	}

	if f.Password == "" {
		return fmt.Errorf("Password is required")
	}

	return nil
}
