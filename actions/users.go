// users.go in the actions folder

package actions

import (
	"taskmanager/models"

	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v6"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// UserRegister creates a new user account
func UserRegister(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}

	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

	user.PasswordHash = string(hashedPassword)
	verrs, err := tx.ValidateAndCreate(user)

	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	return c.Render(201, r.JSON(user))
}

// UserLogin logs in a user and returns a JWT token
func UserLogin(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}

	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	err := tx.Where("email = ?", user.Email).First(user)
	if err != nil {
		return c.Error(404, errors.New("User not found"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(user.PasswordHash))
	if err != nil {
		return c.Error(401, errors.New("Invalid password"))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
	})

	secretKey := envy.Get("JWT_SECRET", "")
	if secretKey == "" {
		return c.Error(500, errors.New("JWT secret not configured"))
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(map[string]string{"token": tokenString}))
}

// UserUpdate updates user's information
func UserUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}

	// Fetch user from context (this would require the user to be set in context, e.g., via a middleware)
	userID := c.Value("user_id").(string)
	if userID == "" {
		return c.Error(401, errors.New("User not authenticated"))
	}

	if err := tx.Find(user, userID); err != nil {
		return c.Error(404, errors.New("User not found"))
	}

	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	return c.Render(200, r.JSON(user))
}

// UserDelete deletes user's account
func UserDelete(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}

	userID := c.Value("user_id").(string)
	if userID == "" {
		return c.Error(401, errors.New("User not authenticated"))
	}

	if err := tx.Find(user, userID); err != nil {
		return c.Error(404, errors.New("User not found"))
	}

	if err := tx.Destroy(user); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(map[string]string{"message": "User deleted successfully"}))
}
