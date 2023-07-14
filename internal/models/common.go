package models

import (
	"fmt"

	"github.com/hibare/DomainHQ/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDBUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Current.DB.Username, config.Current.DB.Password, config.Current.DB.Host, config.Current.DB.Port, config.Current.DB.Name)
}

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(getDBUrl()), &gorm.Config{})
	if err != nil {
		return db, err
	}

	db.AutoMigrate(&GPGPubKeyStore{}, &GPGUsers{})
	return db, nil
}
