package core

import (
	"fmt"
	"github.com/Metisnetwork/Metis-Carrier/types"
)

// StateProcessor implements Processor.
type CenterProcessor struct {
	config *types.CarrierChainConfig
	bc     *DataCenter
}

// NewStateProcessor initialises a new StateProcessor.
func NewCenterProcessor(config *types.CarrierChainConfig, bc *DataCenter) *CenterProcessor {
	return &CenterProcessor{
		config: config,
		bc:     bc,
	}
}

func (p *CenterProcessor) Process(block *types.Block, config *types.CarrierChainConfig) error {
	for _, metadata := range block.Metadatas() {
		response, err := p.bc.client.SaveMetadata(p.bc.ctx, types.NewMetadataSaveRequest(metadata))
		if err != nil {
			log.WithError(err).WithField("hash", metadata.Hash()).Errorf("save metadata failed")
			return err
		}
		if response.Status != 0 {
			return fmt.Errorf("error: %s", response.Msg)
		}
	}
	for _, resource := range block.Resources() {
		response, err := p.bc.client.SaveResource(p.bc.ctx, types.NewPublishPowerRequest(resource))
		if err != nil {
			log.WithError(err).WithField("hash", resource.Hash()).Errorf("save resource failed")
			return err
		}
		if response.Status != 0 {
			return fmt.Errorf("error: %s", response.Msg)
		}
	}
	for _, identity := range block.Identities() {
		response, err := p.bc.client.SaveIdentity(p.bc.ctx, types.NewSaveIdentityRequest(identity))
		if err != nil {
			log.WithError(err).WithField("hash", identity.Hash()).Errorf("save identity failed")
			return err
		}
		if response.Status != 0 {
			return fmt.Errorf("error: %s", response.Msg)
		}
	}
	for _, task := range block.TaskDatas() {
		response, err := p.bc.client.SaveTask(p.bc.ctx, types.NewSaveTaskRequest(task))
		if err != nil {
			log.WithError(err).WithField("hash", task.Hash()).Errorf("save identity failed")
			return err
		}
		if response.Status != 0 {
			return fmt.Errorf("error: %s", response.Msg)
		}
	}
	return nil
}
