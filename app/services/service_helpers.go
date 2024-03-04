package services

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func InitService(log zerolog.Logger, dbInt interface{}) BaseService {
	db := dbInt.(*gorm.DB)

	return BaseService{
		db: db,
		log: log,
	}
}
