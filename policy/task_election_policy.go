package policy

import (
	"encoding/json"
	"fmt"
)

// VRF election evidence
type VRFElectionEvidence struct {
	Nonce      []byte   `json:"nonce"`
	Weights    [][]byte `json:"weights"`
	ElectionAt uint64   `json:"electionAt"`
}

func NewVRFElectionEvidence(nonce []byte, weights [][]byte, electionAt uint64) *VRFElectionEvidence {
	return &VRFElectionEvidence{
		Nonce:      nonce,
		Weights:    weights,
		ElectionAt: electionAt,
	}
}

// MarshalJSON marshals the original value
func (e *VRFElectionEvidence) MarshalJSON() ([]byte, error) {
	if nil == e {
		return nil, fmt.Errorf("VRFElectionEvidence is nil when MarshalJSON")
	}
	return json.Marshal(&e)
}

// UnmarshalJSON parses MixedcasePlatONAddress
func (e *VRFElectionEvidence) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return fmt.Errorf("input is nil when UnmarshalJSON")
	}
	return json.Unmarshal(input, &e)
}



func (e *VRFElectionEvidence) GetNonce() []byte      { return e.Nonce }
func (e *VRFElectionEvidence) GetWeights() [][]byte  { return e.Weights }
func (e *VRFElectionEvidence) GetElectionAt() uint64 { return e.ElectionAt }
