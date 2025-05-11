package jwtauth

import "github.com/golang-jwt/jwt/v5"

type (
	AuthFactory interface {
		SignToken() string
	}

	Claims struct {
		Id       string `json:"id"`
		RoleCode int    `json:"role_code"`
	}

	AuthMapClaims struct {
		*Claims
		jwt.RegisteredClaims
	}

	authConcrete struct {
		Secret []byte
		Claims *AuthMapClaims `json:"claims"`
	}

	access_token struct {
		*authConcrete
	}

	refresh_token struct{ *authConcrete }

	apiKey struct{ *authConcrete }
)

func (a *authConcrete) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, a.Claims)
	ss, _ := token.SignedString(a.Secret)
	return ss
}
