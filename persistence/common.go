package persistence

import "gorm.io/gorm"

func FindAll[T any](db *gorm.DB) []T {
	var result []T
	db.Find(&result)
	return result
}

type ErrNotFound error
