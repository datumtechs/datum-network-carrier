package handler

import (
	taskmngpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/taskmng"
	"github.com/libp2p/go-libp2p-core/peer"
)

func (s *Service) validateTaskResultMsg(pid peer.ID, r *taskmngpb.TaskResultMsg) error {
	return s.cfg.TaskManager.ValidateTaskResultMsg(pid, r)
}

func (s *Service) onTaskResultMsg(pid peer.ID, r *taskmngpb.TaskResultMsg) error {
	return s.cfg.TaskManager.OnTaskResultMsg(pid, r)
}

