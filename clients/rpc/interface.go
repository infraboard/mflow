package rpc

import "github.com/infraboard/mcube/ioc"

const (
	MFLOW = "mflow_client"
)

func C() *ClientSet {
	return ioc.Config().Get(MFLOW).(*Mflow).cs
}
