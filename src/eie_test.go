// cereate test package and file eie_test.go
package main_test

import (
	"petpal-backend/src/configs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateServicesHandler(t *testing.T) {

	// Set test environment
	env := configs.GetInstance()
	env.SetTestEnv()

	// assert env with real env
	assert.Equal(t, "Test", env.GetName())
	assert.Equal(t, "8000", env.GetPort())
	assert.Equal(t, "mongodb://inwza:strongpassword@localhost:27017", env.GetDB_URI())

}

func TestCreateServicesHandler2(t *testing.T) {

	// Set test environment
	env := configs.GetInstance()
	env.SetProductionEnv()

	// assert env with real env
	assert.Equal(t, "Production", env.GetName())
	assert.Equal(t, "8000", env.GetPort())
	assert.Equal(t, "mongodb://inwza:strongpassword@localhost:27017", env.GetDB_URI())

}
