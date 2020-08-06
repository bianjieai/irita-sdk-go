package types

const (
	// prefixChain defines the prefix of this chain
	prefixChain = "i"

	// PrefixAcc is the prefix for account
	prefixAccount = "a"

	// prefixValidator is the prefix for validator keys
	prefixValidator = "v"

	// prefixConsensus is the prefix for consensus keys
	prefixConsensus = "c"

	// prefixPublic is the prefix for public
	prefixPublic = "p"

	// prefixAddress is the prefix for address
	prefixAddress = "a"
)

var (
	prefixCfg = &AddrPrefixCfg{
		bech32AddressPrefix: map[string]string{
			"account_addr":   prefixChain + prefixAccount + prefixAddress,
			"validator_addr": prefixChain + prefixValidator + prefixAddress,
			"consensus_addr": prefixChain + prefixConsensus + prefixAddress,
			"account_pub":    prefixChain + prefixAccount + prefixPublic,
			"validator_pub":  prefixChain + prefixValidator + prefixPublic,
			"consensus_pub":  prefixChain + prefixConsensus + prefixPublic,
		},
	}
)

type AddrPrefixCfg struct {
	bech32AddressPrefix map[string]string
}

// GetAddrPrefixCfg returns the config instance for the corresponding Network type
func GetAddrPrefixCfg() *AddrPrefixCfg {
	return prefixCfg
}

// GetBech32AccountAddrPrefix returns the Bech32 prefix for account address
func (config *AddrPrefixCfg) GetBech32AccountAddrPrefix() string {
	return config.bech32AddressPrefix["account_addr"]
}

// GetBech32ValidatorAddrPrefix returns the Bech32 prefix for validator address
func (config *AddrPrefixCfg) GetBech32ValidatorAddrPrefix() string {
	return config.bech32AddressPrefix["validator_addr"]
}

// GetBech32ConsensusAddrPrefix returns the Bech32 prefix for consensus node address
func (config *AddrPrefixCfg) GetBech32ConsensusAddrPrefix() string {
	return config.bech32AddressPrefix["consensus_addr"]
}

// GetBech32AccountPubPrefix returns the Bech32 prefix for account public key
func (config *AddrPrefixCfg) GetBech32AccountPubPrefix() string {
	return config.bech32AddressPrefix["account_pub"]
}

// GetBech32ValidatorPubPrefix returns the Bech32 prefix for validator public key
func (config *AddrPrefixCfg) GetBech32ValidatorPubPrefix() string {
	return config.bech32AddressPrefix["validator_pub"]
}

// GetBech32ConsensusPubPrefix returns the Bech32 prefix for consensus node public key
func (config *AddrPrefixCfg) GetBech32ConsensusPubPrefix() string {
	return config.bech32AddressPrefix["consensus_pub"]
}
