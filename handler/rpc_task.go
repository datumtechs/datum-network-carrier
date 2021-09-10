package handler

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/libp2p/go-libp2p-core/peer"
)

func (s *Service) validateTaskResultMsg(pid peer.ID, r *pb.TaskResultMsg) error {
	return s.cfg.TaskManager.ValidateTaskResultMsg(pid, r)
}

func (s *Service) onTaskResultMsg(pid peer.ID, r *pb.TaskResultMsg) error {
	return s.cfg.TaskManager.OnTaskResultMsg(pid, r)
}

