package util

import (
	"github.com/ipfs/ipfs-cluster/api"

	peer "gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
)

// IsRemotePin determines whether a Pin's ReplicationFactor has
// been met, so as to either pin or unpin it from the peer.
func IsRemotePin(c api.Pin, pid peer.ID) bool {
	if c.ReplicationFactorMax < 0 {
		return false
	}

	for _, p := range c.Allocations {
		if p == pid {
			return false
		}
	}
	return true
}
