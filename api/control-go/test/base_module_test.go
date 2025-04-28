package test

import (
	"control-go/model"
	"control-go/seed"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestSeedBaseModule(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&model.BaseModule{})

	if err := seed.SeedBaseModule(db); err != nil {
		t.Fatalf("failed to seed BaseModule: %v", err)
	}

	var count int64
	db.Model(&model.BaseModule{}).Count(&count)
	if count != 3 {
		t.Errorf("expected 3 BaseModule records, got %d", count)
	}
}
