package policy

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

func (e *VRFElectionEvidence) GetNonce() []byte      { return e.Nonce }
func (e *VRFElectionEvidence) GetWeights() [][]byte  { return e.Weights }
func (e *VRFElectionEvidence) GetElectionAt() uint64 { return e.ElectionAt }
