package ipfscluster

import (
	"fmt"

	semver "gx/ipfs/QmYRGECuvQnRX73fcvPnGbYijBcGN2HbKZQ7jh26qmLiHG/semver"
	protocol "gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
)

// Version is the current cluster version. Version alignment between
// components, apis and tools ensures compatibility among them.
var Version = semver.MustParse("0.7.0")

// RPCProtocol is used to send libp2p messages between cluster peers
var RPCProtocol = protocol.ID(
	fmt.Sprintf("/ipfscluster/%d.%d/rpc", Version.Major, Version.Minor),
)
