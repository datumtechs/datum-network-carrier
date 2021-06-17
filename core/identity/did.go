package identity

import "github.com/ethereum/go-ethereum/crypto"

const (
	PREFIX_DID = "did:"
	PREFIX_CA = "ca:"
)

type Identity interface {
	Apply(name string) (string, error)
	Revoke(identity string) error
}

type IdentityByDID struct {
	Name  string
	Identity string
}
type IdentityByCA struct {
	Name  string
	Identity string
}
func NewIdentityByDID () *IdentityByDID {return &IdentityByDID{}}
func (did *IdentityByDID) Apply(name string) (string, error) {
	priv, err := crypto.GenerateKey()
	if nil != err {
		log.Error("Failed to generate ecdsa privateKey, err:", err)
		return "", err
	}
	addr := crypto.PubkeyToAddress(priv.PublicKey).Hex()
	did.Name = name
	did.Identity = PREFIX_DID + addr
	return did.Identity, nil
}

func NewIdentityByCA () *IdentityByCA {return &IdentityByCA{}}