package database

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var (

    DB *gorm.DB
    err error
)

func NewPostgresDB(dsn string) (*gorm.DB, error) {
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return DB, nil
}