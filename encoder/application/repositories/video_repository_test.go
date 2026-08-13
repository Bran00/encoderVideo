package repositories_test

import (
	"enconder/framework/database"
	"testing"
)

func TestVideoRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()
}
