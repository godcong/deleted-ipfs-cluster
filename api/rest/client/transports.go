package client

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/ipfs/ipfs-cluster/api"

	p2phttp "gx/ipfs/QmPEnz47VKuAt2yhAiWEpvDtMYdxcfAVV5opEsRC3pHMkB/go-libp2p-http"
	peerstore "gx/ipfs/QmPiemjiKBC9VA7vZF82m4x1oygtg2c2YVqag8PX7dN1BD/go-libp2p-peerstore"
	madns "gx/ipfs/QmQc7jbDUsxUJZyFJzxVrnrWeECCct6fErEpMqtjyWvCX8/go-multiaddr-dns"
	ipnet "gx/ipfs/QmW7Ump7YyBMr712Ta3iEVh3ZYcfVvJaPryfbCnyE826b4/go-libp2p-interface-pnet"
	pnet "gx/ipfs/QmY4Q5JC4vxLEi8EpVxJM4rcRryEVtH1zRKVTAm6BKV1pg/go-libp2p-pnet"
	libp2p "gx/ipfs/QmdJdFQc5U3RAKgJQGmWR7SSM7TLuER5FWz5Wq6Tzs2CnS/go-libp2p"
)

// This is essentially a http.DefaultTransport. We should not mess
// with it since it's a global variable, and we don't know who else uses
// it, so we create our own.
// TODO: Allow more configuration options.
func (c *defaultClient) defaultTransport() {
	c.transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	c.net = "http"
}

func (c *defaultClient) enableLibp2p() error {
	c.defaultTransport()

	pid, addr, err := api.Libp2pMultiaddrSplit(c.config.APIAddr)
	if err != nil {
		return err
	}

	var prot ipnet.Protector
	if c.config.ProtectorKey != nil && len(c.config.ProtectorKey) > 0 {
		if len(c.config.ProtectorKey) != 32 {
			return errors.New("length of ProtectorKey should be 32")
		}
		var key [32]byte
		copy(key[:], c.config.ProtectorKey)

		prot, err = pnet.NewV1ProtectorFromBytes(&key)
		if err != nil {
			return err
		}
	}

	h, err := libp2p.New(c.ctx, libp2p.PrivateNetwork(prot))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.ctx, ResolveTimeout)
	defer cancel()
	resolvedAddrs, err := madns.Resolve(ctx, addr)
	if err != nil {
		return err
	}

	h.Peerstore().AddAddrs(pid, resolvedAddrs, peerstore.PermanentAddrTTL)
	c.transport.RegisterProtocol("libp2p", p2phttp.NewTransport(h))
	c.net = "libp2p"
	c.p2p = h
	c.hostname = pid.Pretty()
	return nil
}

func (c *defaultClient) enableTLS() error {
	c.defaultTransport()
	// based on https://github.com/denji/golang-tls
	c.transport.TLSClientConfig = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		InsecureSkipVerify: c.config.NoVerifyCert,
	}
	c.net = "https"
	return nil
}
