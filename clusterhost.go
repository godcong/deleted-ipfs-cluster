package ipfscluster

import (
	"context"
	"encoding/hex"

	ma "gx/ipfs/QmNTCey11oxhb1AxDnQBRHtdhap6Ctud872NjAYPYYXPuc/go-multiaddr"
	ipnet "gx/ipfs/QmW7Ump7YyBMr712Ta3iEVh3ZYcfVvJaPryfbCnyE826b4/go-libp2p-interface-pnet"
	pnet "gx/ipfs/QmY4Q5JC4vxLEi8EpVxJM4rcRryEVtH1zRKVTAm6BKV1pg/go-libp2p-pnet"
	host "gx/ipfs/QmaoXrM4Z41PD48JY36YqQGKQpLGjyLA2cKcLsES7YddAq/go-libp2p-host"
	libp2p "gx/ipfs/QmdJdFQc5U3RAKgJQGmWR7SSM7TLuER5FWz5Wq6Tzs2CnS/go-libp2p"
)

// NewClusterHost creates a libp2p Host with the options from the
// provided cluster configuration.
func NewClusterHost(ctx context.Context, cfg *Config) (host.Host, error) {
	var prot ipnet.Protector
	var err error

	// Create protector if we have a secret.
	if cfg.Secret != nil && len(cfg.Secret) > 0 {
		var key [32]byte
		copy(key[:], cfg.Secret)
		prot, err = pnet.NewV1ProtectorFromBytes(&key)
		if err != nil {
			return nil, err
		}
	}

	return libp2p.New(
		ctx,
		libp2p.Identity(cfg.PrivateKey),
		libp2p.ListenAddrs([]ma.Multiaddr{cfg.ListenAddr}...),
		libp2p.PrivateNetwork(prot),
		libp2p.NATPortMap(),
	)
}

// EncodeProtectorKey converts a byte slice to its hex string representation.
func EncodeProtectorKey(secretBytes []byte) string {
	return hex.EncodeToString(secretBytes)
}
