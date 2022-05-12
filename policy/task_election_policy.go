package policy

import "encoding/json"

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

func (e *VRFElectionEvidence) EncodeJson() (string, error) {
	b, err := json.Marshal(&e)
	if nil != err {
		return "", err
	}
	return string(b), nil
}

func (e *VRFElectionEvidence) DecodeJson(jsonStr string) error {
	if err := json.Unmarshal([]byte(jsonStr), &e); nil != err {
		return err
	}
	return nil
}

func (e *VRFElectionEvidence) GetNonce() []byte      { return e.Nonce }
func (e *VRFElectionEvidence) GetWeights() [][]byte  { return e.Weights }
func (e *VRFElectionEvidence) GetElectionAt() uint64 { return e.ElectionAt }
