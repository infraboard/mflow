package rpc

import "github.com/infraboard/mcube/v2/ioc"

const (
	MFLOW = "mflow"
)

func C() *ClientSet {
	return ioc.Config().Get(MFLOW).(*Mflow).cs
}
