package test

  import (
      "control-go/model"
      "control-go/seed"
      "github.com/glebarez/sqlite"
      "gorm.io/gorm"
      "testing"
  )

  func TestSeedProduct(t *testing.T) {
      db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
      if err != nil {
          t.Fatalf("failed to connect database: %v", err)
      }

      db.AutoMigrate(&model.Product{})

      if err := seed.SeedProduct(db); err != nil {
          t.Fatalf("failed to seed Product: %v", err)
      }

      var count int64
      db.Model(&model.Product{}).Count(&count)
      if count != 3 {
          t.Errorf("expected 3 Product records, got %d", count)
      }
  }