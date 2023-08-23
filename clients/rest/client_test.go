package rest_test

import (
	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/mflow/clients/rest"
)

var (
	client *rest.ClientSet
)

func init() {
	zap.DevelopmentSetup()
	conf := rest.NewDefaultConfig()
	client = rest.NewClient(conf)
}
