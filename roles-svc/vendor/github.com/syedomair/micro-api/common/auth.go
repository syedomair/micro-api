package common

import (
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
)

func CheckAuth(tokenString string) (string, string, error) {

	type Claims struct {
		CurrentUserId string `json:"current_user_id"`
		ClientId     string `json:"client_id"`
		IsAdmin       string `json:"is_admin"`
		jwt.StandardClaims
	}
	tokenClaims := Claims{}

	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SIGNING_KEY), nil
	})
	if err != nil {
		return "", "", errors.New(err.Error())
	}
	if token.Valid {
		return tokenClaims.CurrentUserId, tokenClaims.ClientId, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return "", "", errors.New("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return "", "", errors.New("Timing is everything")
		} else {
			return "", "", errors.New("Couldn't handle this token")
		}
	} else {
		return "", "", errors.New("Couldn't handle this token")
	}
}
func GetAPIKey(tokenString string) (string, error) {

	token, _ := ValidateJWTToken(tokenString, "")
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Claim error")
	}
	return claims["api_key"].(string), nil
}

func ValidateJWTToken(tokenString string, tokenSecret string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if token.Valid {
		return token, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return token, errors.New("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return token, errors.New("Timing is everything")
		} else {
			return token, errors.New("Couldn't handle this token:")
		}
	} else {
		return token, errors.New("Couldn't handle this token:")
	}

}
