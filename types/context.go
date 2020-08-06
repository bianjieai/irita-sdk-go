package types

import (
	"fmt"

	"github.com/bianjieai/irita-sdk-go/codec"
)

// TxBuilder implements a transaction context created in SDK modules.
type TxBuilder struct {
	address       string
	chainID       string
	memo          string
	password      string
	accountNumber uint64
	sequence      uint64
	gas           uint64

	fee        Coins
	mode       BroadcastMode
	simulate   bool
	codec      codec.Marshaler
	txEncoder  TxEncoder
	keyManager KeyManager
}

func NewTxBuilder(txEncoder TxEncoder) *TxBuilder {
	return &TxBuilder{txEncoder: txEncoder}
}

// WithCodec returns a pointer of the context with an updated codec.
func (builder *TxBuilder) WithCodec(cdc codec.Marshaler) *TxBuilder {
	builder.codec = cdc
	return builder
}

// Codec returns codec.
func (builder *TxBuilder) Codec() codec.Marshaler {
	return builder.codec
}

// WithChainID returns a pointer of the context with an updated ChainID.
func (builder *TxBuilder) WithChainID(chainID string) *TxBuilder {
	builder.chainID = chainID
	return builder
}

// ChainID returns the chainID of the current chain.
func (builder *TxBuilder) ChainID() string {
	return builder.chainID
}

// WithGas returns a pointer of the context with an updated Gas.
func (builder *TxBuilder) WithGas(gas uint64) *TxBuilder {
	builder.gas = gas
	return builder
}

// Gas returns the gas of the transaction.
func (builder *TxBuilder) Gas() uint64 {
	return builder.gas
}

// WithFee returns a pointer of the context with an updated Fee.
func (builder *TxBuilder) WithFee(fee Coins) *TxBuilder {
	builder.fee = fee
	return builder
}

// Fee returns the fee of the transaction.
func (builder *TxBuilder) Fee() Coins {
	return builder.fee
}

// WithSequence returns a pointer of the context with an updated sequence number.
func (builder *TxBuilder) WithSequence(sequence uint64) *TxBuilder {
	builder.sequence = sequence
	return builder
}

// Sequence returns the sequence of the account.
func (builder *TxBuilder) Sequence() uint64 {
	return builder.sequence
}

// WithMemo returns a pointer of the context with an updated memo.
func (builder *TxBuilder) WithMemo(memo string) *TxBuilder {
	builder.memo = memo
	return builder
}

// Memo returns memo.
func (builder *TxBuilder) Memo() string {
	return builder.memo
}

// WithAccountNumber returns a pointer of the context with an account number.
func (builder *TxBuilder) WithAccountNumber(accnum uint64) *TxBuilder {
	builder.accountNumber = accnum
	return builder
}

// AccountNumber returns accountNumber.
func (builder *TxBuilder) AccountNumber() uint64 {
	return builder.accountNumber
}

// WithAccountNumber returns a pointer of the context with a keyDao.
func (builder *TxBuilder) WithKeyManager(keyManager KeyManager) *TxBuilder {
	builder.keyManager = keyManager
	return builder
}

// KeyManager returns keyManager.
func (builder *TxBuilder) KeyManager() KeyManager {
	return builder.keyManager
}

// WithMode returns a pointer of the context with a Mode.
func (builder *TxBuilder) WithMode(mode BroadcastMode) *TxBuilder {
	builder.mode = mode
	return builder
}

// Mode returns mode.
func (builder *TxBuilder) Mode() BroadcastMode {
	return builder.mode
}

// WithRPC returns a pointer of the context with a simulate.
func (builder *TxBuilder) WithSimulate(simulate bool) *TxBuilder {
	builder.simulate = simulate
	return builder
}

// Simulate returns simulate.
func (builder *TxBuilder) Simulate() bool {
	return builder.simulate
}

// WithRPC returns a pointer of the context with a password.
func (builder *TxBuilder) WithPassword(password string) *TxBuilder {
	builder.password = password
	return builder
}

// Password returns password.
func (builder *TxBuilder) Password() string {
	return builder.password
}

// WithAddress returns a pointer of the context with a password.
func (builder *TxBuilder) WithAddress(address string) *TxBuilder {
	builder.address = address
	return builder
}

// Address returns the address.
func (builder *TxBuilder) Address() string {
	return builder.address
}

func (builder *TxBuilder) BuildAndSign(name string, msgs []Msg) ([]byte, error) {
	msg, err := builder.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}
	return builder.Sign(name, msg)
}

// BuildTxForSim creates a StdSignMsg and encodes a transaction with the
// StdSignMsg with a single empty StdSignature for tx simulation.
func (builder TxBuilder) BuildTxForSim(msgs []Msg) ([]byte, error) {
	signMsg, err := builder.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}

	// the ante handler will populate with a sentinel pubkey
	sigs := []StdSignature{{}}
	return builder.txEncoder(NewStdTx(signMsg.Msgs, signMsg.Fee, sigs, signMsg.Memo))
}

// BuildSignMsg builds a single message to be signed from a TxBuilder given a
// set of messages. It returns an error if a fee is supplied but cannot be
// parsed.
func (builder TxBuilder) BuildSignMsg(msgs []Msg) (StdSignMsg, error) {
	if builder.chainID == "" {
		return StdSignMsg{}, fmt.Errorf("chain ID required but not specified")
	}

	return StdSignMsg{
		ChainID:       builder.chainID,
		AccountNumber: builder.accountNumber,
		Sequence:      builder.sequence,
		Memo:          builder.memo,
		Msgs:          msgs,
		Fee:           NewStdFee(builder.gas, builder.fee...),
	}, nil
}

// Sign signs a transaction given a name, passphrase, and a single message to
// signed. An error is returned if signing fails.
func (builder *TxBuilder) Sign(name string, msg StdSignMsg) ([]byte, error) {
	sig, err := builder.makeSignature(name, msg)
	if err != nil {
		return nil, err
	}

	return builder.MarshallTx([]StdSignature{sig}, msg)
}

// Sign signs a transaction given a name, passphrase, and a single message to
// signed. An error is returned if signing fails.
func (builder *TxBuilder) MarshallTx(sigs []StdSignature, msg StdSignMsg) ([]byte, error) {
	return builder.txEncoder(NewStdTx(msg.Msgs, msg.Fee, sigs, msg.Memo))
}

func (builder *TxBuilder) makeSignature(name string, msg StdSignMsg) (sig StdSignature, err error) {
	signature, pubKey, err := builder.keyManager.Sign(name, builder.password, msg.Bytes(builder.codec))
	if err != nil {
		return sig, err
	}
	return StdSignature{
		PubKey:    pubKey.Bytes(),
		Signature: signature,
	}, nil
}
