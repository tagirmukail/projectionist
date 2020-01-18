package validate

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"projectionist/models"
	"strings"
)

func ValidateToken(token string, tokenSecretKey string) error {
	var splitted = strings.Split(token, " ")
	if len(splitted) != 2 {
		return fmt.Errorf("invalid/Malformed authorization token")
	}

	var tokenPart = splitted[1]
	var tokenM = &models.Token{}

	tokenJWT, err := jwt.ParseWithClaims(tokenPart, tokenM, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(tokenSecretKey), nil
	})
	if err != nil {
		return fmt.Errorf("malformed authentication token")
	}

	if !tokenJWT.Valid {
		return fmt.Errorf("authentication token is not valid")
	}

	return nil
}
