package repository_test

import (
	"testing"

	"github.com/runetale/runevision/di"
	"github.com/runetale/runevision/domain/entity"
	"gopkg.in/go-playground/assert.v1"
)

func TestDashboardRepository(t *testing.T) {
	setup()

	repo := di.InitializeDashboardRepository()
	err := repo.Create(db, entity.NewDashboardHistory("test", "ok", 30))
	if err != nil {
		t.Fatal(err)
	}

	histories, err := repo.GetHistories(db)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(histories))
}
