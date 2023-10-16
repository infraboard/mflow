package impl

import (
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/mflow/apps/log"
	"github.com/infraboard/mflow/conf"
)

func init() {
	ioc.RegistryController(&service{})
}

type service struct {
	log *zerolog.Logger
	db  *mongo.Database
	ioc.ObjectImpl
}

func (s *service) Init() error {
	db, err := conf.C().Mongo.GetDB()
	if err != nil {
		return err
	}

	s.db = db
	s.log = logger.Sub("storage")
	return nil
}

func (s *service) Name() string {
	return log.AppName
}
