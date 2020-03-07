package config

import (
	"github.com/bianjieai/irita/config"
)

type (
	IritaAddrPrefixConfig struct {
		Conf *config.Config
	}
)

func GetIritaAddrPrefixConfig(network string) IritaAddrPrefixConfig {
	config.SetNetworkType(network)
	c := config.GetConfig()
	return IritaAddrPrefixConfig{
		Conf: c,
	}
}

// GetBech32AccountAddrPrefix returns the Bech32 prefix for account address
func (c IritaAddrPrefixConfig) GetBech32AccountAddrPrefix2() string {
	return c.Conf.GetBech32AccountPubPrefix()
}

// GetBech32ValidatorAddrPrefix returns the Bech32 prefix for validator address
func (c IritaAddrPrefixConfig) GetBech32ValidatorAddrPrefix() string {
	return c.Conf.GetBech32ValidatorAddrPrefix()
}

// GetBech32ConsensusAddrPrefix returns the Bech32 prefix for consensus node address
func (c IritaAddrPrefixConfig) GetBech32ConsensusAddrPrefix() string {
	return c.Conf.GetBech32ConsensusAddrPrefix()
}

// GetBech32AccountPubPrefix returns the Bech32 prefix for account public key
func (c IritaAddrPrefixConfig) GetBech32AccountPubPrefix() string {
	return c.Conf.GetBech32AccountPubPrefix()
}

// GetBech32ValidatorPubPrefix returns the Bech32 prefix for validator public key
func (c IritaAddrPrefixConfig) GetBech32ValidatorPubPrefix() string {
	return c.Conf.GetBech32ValidatorPubPrefix()
}

// GetBech32ConsensusPubPrefix returns the Bech32 prefix for consensus node public key
func (c IritaAddrPrefixConfig) GetBech32ConsensusPubPrefix() string {
	return c.Conf.GetBech32ConsensusPubPrefix()
}
