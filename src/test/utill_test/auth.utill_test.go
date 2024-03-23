package main

import (
	"testing"

	configs "petpal-backend/src/configs"
	"petpal-backend/src/models"
	"petpal-backend/src/utills"
	auth_utills "petpal-backend/src/utills/auth"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	// create env
	env := configs.GetInstance()
	env.SetTestEnv()

	// create db
	db, _ := utills.NewMongoDB()

	t.Run("Login with svcp login type", func(t *testing.T) {
		req := &models.LoginReq{
			Email:     "0@svcp.com",
			Password:  "password",
			LoginType: "svcp",
		}

		res, err := auth_utills.Login(db, req)

		assert.NoError(t, err)
		assert.Equal(t, "svcp", res.LoginType)
		assert.Equal(t, "0@svcp.com", res.UserEmail)
	})

	t.Run("Login with user login type", func(t *testing.T) {
		req := &models.LoginReq{
			Email:     "0@user.com",
			Password:  "password",
			LoginType: "user",
		}

		res, err := auth_utills.Login(db, req)

		assert.NoError(t, err)
		assert.Equal(t, "user", res.LoginType)
		assert.Equal(t, "0@user.com", res.UserEmail)
	})

	t.Run("Login with admin login type", func(t *testing.T) {
		req := &models.LoginReq{
			Email:     "0@admin.com",
			Password:  "password",
			LoginType: "admin",
		}

		res, err := auth_utills.Login(db, req)

		assert.NoError(t, err)
		assert.Equal(t, "admin", res.LoginType)
		assert.Equal(t, "0@admin.com", res.UserEmail)
	})

	t.Run("Login with invalid login type", func(t *testing.T) {
		req := &models.LoginReq{
			Email:     "test@test.com",
			Password:  "password",
			LoginType: "invalid",
		}

		_, err := auth_utills.Login(db, req)

		assert.Error(t, err)
	})

	t.Run("Login with invalid email", func(t *testing.T) {
		req := &models.LoginReq{
			Email:     "invalid",
			Password:  "password",
			LoginType: "svcp",
		}

		_, err := auth_utills.Login(db, req)

		assert.Error(t, err)
	})

	t.Run("Login with invalid password", func(t *testing.T) {
		req := &models.LoginReq{
			Email:     "0@admin.com",
			Password:  "invalid",
			LoginType: "admin",
		}

		_, err := auth_utills.Login(db, req)

		assert.Error(t, err)
	})
}
