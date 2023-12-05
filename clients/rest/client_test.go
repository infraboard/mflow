package rest_test

import (
	"github.com/infraboard/mflow/clients/rest"
)

var (
	client *rest.ClientSet
)

func init() {
	conf := rest.NewDefaultConfig()
	client = rest.NewClient(conf)
}
