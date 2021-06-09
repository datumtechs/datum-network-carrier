package service

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/p2p"
)

type Carrier struct {
		Cfg 			*common.CarrierConfig

}

func New(cfg *common.CarrierConfig) *Carrier {

	carrier := &Carrier{
		Cfg:  cfg,
	}


	return carrier
}

func (c *Carrier) Engine() {

}

func (c *Carrier) Protocols() []p2p.Protocol {
	return nil
}

func (c *Carrier) Start() error {
	return nil
}

func (c *Carrier) Stop() error {
	return nil
}