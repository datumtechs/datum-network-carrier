package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/cmd/utils"
	gcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/urfave/cli/v2"
)

type outputGenkeypair struct {
	PrivateKey string
	PublicKey  string
	PeerId     string
}

var commandGenkeypair = &cli.Command{
	Name:      "genkeypair",
	Usage:     "generate new private key pair",
	ArgsUsage: "[ ]",
	Description: `
Generate a new private key pair.
.
`,
	Flags: []cli.Flag{
		jsonFlag,
	},
	Action: func(ctx *cli.Context) error {
		// Check if keyfile path given and make sure it doesn't already exist.
		var privateKey *ecdsa.PrivateKey
		var err error
		// generate random.
		privateKey, err = gcrypto.GenerateKey()
		if err != nil {
			utils.Fatalf("Failed to generate random private key: %v", err)
		}

		typeAssertedKey := crypto.PubKey((*crypto.Secp256k1PublicKey)(&privateKey.PublicKey))
		id, err := peer.IDFromPublicKey(typeAssertedKey)

		// Output some information.
		out := outputGenkeypair{
			PublicKey:  hex.EncodeToString(gcrypto.FromECDSAPub(&privateKey.PublicKey)[1:]),
			PrivateKey: hex.EncodeToString(gcrypto.FromECDSA(privateKey)),
			PeerId:     id.String(),
		}
		if ctx.Bool(jsonFlag.Name) {
			mustPrintJSON(out)
		} else {
			fmt.Println("PrivateKey: ", out.PrivateKey)
			fmt.Println("PublicKey : ", out.PublicKey)
			fmt.Println("PeerId : ", out.PeerId)
		}
		return nil
	},
}
