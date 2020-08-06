package tx

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// Factory defines a client transaction factory that facilitates generating and
// signing an application-specific transaction.
type Factory struct {
	txGenerator        Generator
	accountNumber      uint64
	sequence           uint64
	gas                uint64
	gasAdjustment      float64
	simulateAndExecute bool
	chainID            string
	memo               string
	fees               sdk.Coins
	gasPrices          sdk.DecCoins
	keyManager         sdk.KeyManager
	from               string
	password           string
}

func (f Factory) AccountNumber() uint64   { return f.accountNumber }
func (f Factory) Sequence() uint64        { return f.sequence }
func (f Factory) Gas() uint64             { return f.gas }
func (f Factory) GasAdjustment() float64  { return f.gasAdjustment }
func (f Factory) ChainID() string         { return f.chainID }
func (f Factory) Memo() string            { return f.memo }
func (f Factory) Fees() sdk.Coins         { return f.fees }
func (f Factory) GasPrices() sdk.DecCoins { return f.gasPrices }

// SimulateAndExecute returns the option to simulate and then execute the transaction
// using the gas from the simulation results
func (f Factory) SimulateAndExecute() bool { return f.simulateAndExecute }

// WithTxGenerator returns a copy of the Factory with an updated Generator.
func (f Factory) WithTxGenerator(g Generator) Factory {
	f.txGenerator = g
	return f
}

// WithChainID returns a copy of the Factory with an updated chainID.
func (f Factory) WithChainID(chainID string) Factory {
	f.chainID = chainID
	return f
}

// WithGas returns a copy of the Factory with an updated gas value.
func (f Factory) WithGas(gas uint64) Factory {
	f.gas = gas
	return f
}

// WithFees returns a copy of the Factory with an updated fee.
func (f Factory) WithFees(fees string) Factory {
	parsedFees, err := sdk.ParseCoins(fees)
	if err != nil {
		panic(err)
	}

	f.fees = parsedFees
	return f
}

// WithGasPrices returns a copy of the Factory with updated gas prices.
func (f Factory) WithGasPrices(gasPrices string) Factory {
	parsedGasPrices, err := sdk.ParseDecCoins(gasPrices)
	if err != nil {
		panic(err)
	}

	f.gasPrices = parsedGasPrices
	return f
}

// WithSequence returns a copy of the Factory with an updated sequence number.
func (f Factory) WithSequence(sequence uint64) Factory {
	f.sequence = sequence
	return f
}

// WithMemo returns a copy of the Factory with an updated memo.
func (f Factory) WithMemo(memo string) Factory {
	f.memo = memo
	return f
}

// WithAccountNumber returns a copy of the Factory with an updated account number.
func (f Factory) WithAccountNumber(accnum uint64) Factory {
	f.accountNumber = accnum
	return f
}

// WithGasAdjustment returns a copy of the Factory with an updated gas adjustment.
func (f Factory) WithGasAdjustment(gasAdj float64) Factory {
	f.gasAdjustment = gasAdj
	return f
}

// WithSimulateAndExecute returns a copy of the Factory with an updated gas
// simulation value.
func (f Factory) WithSimulateAndExecute(sim bool) Factory {
	f.simulateAndExecute = sim
	return f
}
