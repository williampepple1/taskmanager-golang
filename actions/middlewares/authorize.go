package middlewares

import (
	"errors"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/golang-jwt/jwt"
)

func AuthorizeMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(envy.Get("JWT_SECRET", "")), nil
		})

		if err != nil {
			return c.Error(401, errors.New("Unauthorized"))
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
			return next(c)
		}

		return c.Error(401, errors.New("Unauthorized"))
	}
}
