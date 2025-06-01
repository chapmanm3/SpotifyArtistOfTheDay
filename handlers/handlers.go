package handlers

import (
	"gorm.io/gorm"
)

type DBHandlerService struct {
	DB *gorm.DB
}

type APIHandlerService interface {
}

func New(db *gorm.DB) DBHandlerService {
	return DBHandlerService{DB: db}
}
