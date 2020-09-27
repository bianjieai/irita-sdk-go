package types

import "github.com/tendermint/tendermint/proto/tendermint/crypto"

type ProofValue struct {
	Proof []byte   `json:"proof"`
	Path  []string `json:"path"`
	Value []byte   `json:"value"`
}

type MerkleProof struct {
	Proof *crypto.ProofOps `json:"proof"`
}
