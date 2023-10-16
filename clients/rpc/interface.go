package rpc

import "github.com/infraboard/mcube/ioc"

const (
	MFLOW = "mflow"
)

func C() *ClientSet {
	return ioc.Config().Get(MFLOW).(*Mflow).cs
}
