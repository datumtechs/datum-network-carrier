package policy

import "encoding/json"

// VRF election
type VRFElectionEvidence struct {
	nonce      []byte
	weights    [][]byte
	electionAt uint64
}

func NewVRFElectionEvidence(nonce []byte, weights [][]byte, electionAt uint64) *VRFElectionEvidence {
	return &VRFElectionEvidence{
		nonce:      nonce,
		weights:    weights,
		electionAt: electionAt,
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

func (e *VRFElectionEvidence) GetNonce() []byte { return e.nonce }
func (e *VRFElectionEvidence) GetWeights() [][]byte { return e.weights }
func (e *VRFElectionEvidence) GetElectionAt() uint64 { return e.electionAt }