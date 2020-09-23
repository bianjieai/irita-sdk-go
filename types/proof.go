package types

import (
	"github.com/tendermint/tendermint/crypto/merkle"
)

type ProofValue struct {
	Proof []byte   `json:"proof"`
	Path  []string `json:"path"`
	Value []byte   `json:"value"`
}

type MerkleProof struct {
	Proof *merkle.Proof `json:"proof"`
}
