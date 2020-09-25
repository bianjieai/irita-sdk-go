package types

import (
	"encoding/base64"
	"github.com/tendermint/tendermint/crypto/encoding"

	"github.com/bianjieai/irita-sdk-go/codec"

	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Block struct {
	tmtypes.Header `json:"header"`
	Data           `json:"data"`
	Evidence       tmtypes.EvidenceData `json:"evidence"`
	LastCommit     *tmtypes.Commit      `json:"last_commit"`
}

type Data struct {
	Txs []StdTx `json:"txs"`
}

func ParseBlock(cdc *codec.LegacyAmino, block *tmtypes.Block) Block {
	var txs []StdTx
	for _, tx := range block.Txs {
		var stdTx StdTx
		if err := cdc.UnmarshalBinaryBare(tx, &stdTx); err == nil {
			txs = append(txs, stdTx)
		}
	}
	return Block{
		Header: block.Header,
		Data: Data{
			Txs: txs,
		},
		Evidence:   block.Evidence,
		LastCommit: block.LastCommit,
	}
}

type BlockResult struct {
	Height  int64         `json:"height"`
	Results ABCIResponses `json:"results"`
}

type BlockDetail struct {
	BlockID     tmtypes.BlockID `json:"block_id"`
	Block       Block           `json:"block"`
	BlockResult BlockResult     `json:"block_result"`
}

type ABCIResponses struct {
	DeliverTx  []TxResult
	EndBlock   ResultEndBlock
	BeginBlock ResultBeginBlock
}

type ResultBeginBlock struct {
	Events StringEvents `json:"events"`
}

type ResultEndBlock struct {
	Events           StringEvents      `json:"events"`
	ValidatorUpdates []ValidatorUpdate `json:"validator_updates"`
}

func ParseValidatorUpdate(updates []abci.ValidatorUpdate) []ValidatorUpdate {
	var vUpdates []ValidatorUpdate
	for _, v := range updates {
		pubkey, _ := encoding.PubKeyFromProto(v.PubKey)
		vUpdates = append(vUpdates, ValidatorUpdate{
			PubKey: PubKey{
				Type:  pubkey.Type(),
				Value: base64.StdEncoding.EncodeToString(pubkey.Bytes()),
			},
			Power: v.Power,
		})
	}
	return vUpdates
}

func ParseBlockResult(res *ctypes.ResultBlockResults) BlockResult {
	var txResults = make([]TxResult, len(res.TxsResults))
	for i, r := range res.TxsResults {
		txResults[i] = TxResult{
			Code:      r.Code,
			Log:       r.Log,
			GasWanted: r.GasWanted,
			GasUsed:   r.GasUsed,
			Events:    StringifyEvents(r.Events),
		}
	}
	return BlockResult{
		Height: res.Height,
		Results: ABCIResponses{
			DeliverTx: txResults,
			EndBlock: ResultEndBlock{
				Events:           StringifyEvents(res.EndBlockEvents),
				ValidatorUpdates: ParseValidatorUpdate(res.ValidatorUpdates),
			},
			BeginBlock: ResultBeginBlock{
				Events: StringifyEvents(res.BeginBlockEvents),
			},
		},
	}
}

func ParseValidators(cdc *codec.LegacyAmino, vs []*tmtypes.Validator) []Validator {
	var validators = make([]Validator, len(vs))
	for i, v := range vs {
		bech32Addr, _ := ConsAddressFromHex(v.Address.String())
		bech32PubKey, _ := Bech32ifyPubKey(Bech32PubKeyTypeConsPub, v.PubKey)

		var pubKey PubKey
		if bz, err := cdc.MarshalJSON(v.PubKey); err == nil {
			_ = cdc.UnmarshalJSON(bz, &pubKey)
		}
		validators[i] = Validator{
			Bech32Address:    bech32Addr.String(),
			Bech32PubKey:     bech32PubKey,
			Address:          v.Address.String(),
			PubKey:           pubKey,
			VotingPower:      v.VotingPower,
			ProposerPriority: v.ProposerPriority,
		}
	}
	return validators
}
