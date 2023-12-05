package rpc

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/tools/pretty"
)

func init() {
	ioc.Config().Registry(&Mflow{})
}

type Mflow struct {
	ioc.ObjectImpl
	cs *ClientSet
}

func (m *Mflow) String() string {
	return pretty.ToJSON(m)
}

func (m *Mflow) Name() string {
	return MFLOW
}

func (m *Mflow) Init() error {
	cs, err := NewClient()
	if err != nil {
		return err
	}
	m.cs = cs
	return nil
}

func (c *Mflow) Close(ctx context.Context) error {
	if c.cs != nil {
		c.cs.conn.Close()
	}

	return nil
}
