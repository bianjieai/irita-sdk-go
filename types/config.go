package types

import (
	"fmt"
	"google.golang.org/grpc"
	"os"

	"github.com/bianjieai/irita-sdk-go/types/store"
)

const (
	defaultGas           = 200000
	defaultFees          = "4point"
	defaultTimeout       = 5
	defaultLevel         = "info"
	defaultMaxTxsBytes   = 1073741824
	defaultAlgo          = "sm2"
	defaultMode          = Sync
	defaultPath          = "$HOME/irita-sdk-go/leveldb"
	defaultGasAdjustment = 1.0
)

type ClientConfig struct {
	// irita node rpc address
	NodeURI string

	// irita grpc address
	GRPCAddr string

	// grpc dial options
	GRPCOptions []grpc.DialOption

	// irita chain-id
	ChainID string

	// max gas limit
	Gas uint64

	// Fee amount of point
	Fee DecCoins

	// PrivKeyArmor DAO Implements
	KeyDAO store.KeyDAO

	// Private key generation algorithm(sm2,secp256k1)
	Algo string

	// Transaction broadcast Mode
	Mode BroadcastMode

	//Transaction broadcast timeout(seconds)
	Timeout uint

	//log level(trace|debug|info|warn|error|fatal|panic)
	Level string

	//maximum bytes of a transaction
	MaxTxBytes uint64

	//adjustment factor to be multiplied against the estimate returned by the tx simulation;
	GasAdjustment float64

	//whether to enable caching
	Cached bool
}

func NewClientConfig(uri, grpcAddr, chainID string, options ...Option) (ClientConfig, error) {
	cfg := ClientConfig{
		NodeURI:  uri,
		ChainID:  chainID,
		GRPCAddr: grpcAddr,
	}
	for _, optionFn := range options {
		if err := optionFn(&cfg); err != nil {
			return ClientConfig{}, err
		}
	}

	if err := cfg.checkAndSetDefault(); err != nil {
		return ClientConfig{}, err
	}
	return cfg, nil
}

func (cfg *ClientConfig) checkAndSetDefault() error {
	if len(cfg.NodeURI) == 0 {
		return fmt.Errorf("nodeURI is required")
	}

	if len(cfg.ChainID) == 0 {
		return fmt.Errorf("chainID is required")
	}

	if err := GasOption(cfg.Gas)(cfg); err != nil {
		return err
	}

	if err := FeeOption(cfg.Fee)(cfg); err != nil {
		return err
	}

	if err := AlgoOption(cfg.Algo)(cfg); err != nil {
		return err
	}

	if err := KeyDAOOption(cfg.KeyDAO)(cfg); err != nil {
		return err
	}

	if err := ModeOption(cfg.Mode)(cfg); err != nil {
		return err
	}

	if err := TimeoutOption(cfg.Timeout)(cfg); err != nil {
		return err
	}

	if err := LevelOption(cfg.Level)(cfg); err != nil {
		return err
	}

	if err := MaxTxBytesOption(cfg.MaxTxBytes)(cfg); err != nil {
		return err
	}

	return GasAdjustmentOption(cfg.GasAdjustment)(cfg)
}

type Option func(cfg *ClientConfig) error

func FeeOption(fee DecCoins) Option {
	return func(cfg *ClientConfig) error {
		if fee == nil || fee.Empty() || !fee.IsValid() {
			fees, _ := ParseDecCoins(defaultFees)
			fee = fees
		}
		cfg.Fee = fee
		return nil
	}
}

func KeyDAOOption(dao store.KeyDAO) Option {
	return func(cfg *ClientConfig) error {
		if dao == nil {
			defaultPath := os.ExpandEnv(defaultPath)
			levelDB, err := store.NewLevelDB(defaultPath, nil)
			if err != nil {
				return err
			}
			dao = levelDB
		}
		cfg.KeyDAO = dao
		return nil
	}
}

func GasOption(gas uint64) Option {
	return func(cfg *ClientConfig) error {
		if gas <= 0 {
			gas = defaultGas
		}
		cfg.Gas = gas
		return nil
	}
}

func AlgoOption(algo string) Option {
	return func(cfg *ClientConfig) error {
		if algo == "" {
			algo = defaultAlgo
		}
		cfg.Algo = algo
		return nil
	}
}

func ModeOption(mode BroadcastMode) Option {
	return func(cfg *ClientConfig) error {
		if mode == "" {
			mode = defaultMode
		}
		cfg.Mode = mode
		return nil
	}
}

func TimeoutOption(timeout uint) Option {
	return func(cfg *ClientConfig) error {
		if timeout <= 0 {
			timeout = defaultTimeout
		}
		cfg.Timeout = timeout
		return nil
	}
}

func LevelOption(level string) Option {
	return func(cfg *ClientConfig) error {
		if level == "" {
			level = defaultLevel
		}
		cfg.Level = level
		return nil
	}
}

func MaxTxBytesOption(maxTxBytes uint64) Option {
	return func(cfg *ClientConfig) error {
		if maxTxBytes <= 0 {
			maxTxBytes = defaultMaxTxsBytes
		}
		cfg.MaxTxBytes = maxTxBytes
		return nil
	}
}

func GasAdjustmentOption(gasAdjustment float64) Option {
	return func(cfg *ClientConfig) error {
		if gasAdjustment <= 0 {
			gasAdjustment = defaultGasAdjustment
		}
		cfg.GasAdjustment = gasAdjustment
		return nil
	}
}

func CachedOption(enabled bool) Option {
	return func(cfg *ClientConfig) error {
		cfg.Cached = enabled
		return nil
	}
}
