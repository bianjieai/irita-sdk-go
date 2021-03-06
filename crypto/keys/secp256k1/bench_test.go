package secp256k1

import (
	"io"
	"testing"

	"github.com/tendermint/tendermint/crypto"

	"github.com/bianjieai/irita-sdk-go/crypto/keys/internal/benchmarking"
)

func BenchmarkKeyGeneration(b *testing.B) {
	benchmarkKeygenWrapper := func(reader io.Reader) crypto.PrivKey {
		return &PrivKey{Key: genPrivKey(crypto.CReader())}
	}
	benchmarking.BenchmarkKeyGeneration(b, benchmarkKeygenWrapper)
}

func BenchmarkSigning(b *testing.B) {
	priv := GenPrivKey()
	benchmarking.BenchmarkSigning(b, priv)
}

func BenchmarkVerification(b *testing.B) {
	priv := GenPrivKey()
	benchmarking.BenchmarkVerification(b, priv)
}
