package db

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(dsn string) (*gorm.DB, error) {
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	return gdb, nil
}

func WithContext(ctx context.Context, gdb *gorm.DB) *gorm.DB {
	return gdb.WithContext(ctx)
}
