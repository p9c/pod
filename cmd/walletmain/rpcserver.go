package walletmain

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/p9c/pod/pkg/opts"
	walletrpc2 "github.com/p9c/pod/pkg/walletrpc"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/wallet"
)

type listenFunc func(net string, laddr string) (net.Listener, error)

// GenerateRPCKeyPair generates a new RPC TLS keypair and writes the cert and possibly also the key in PEM format to the
// paths specified by the config. If successful, the new keypair is returned.
func GenerateRPCKeyPair(config *opts.Config, writeKey bool) (tls.Certificate, error) {
	D.Ln("generating TLS certificates")
	// Create directories for cert and key files if they do not yet exist.
	D.Ln("rpc tls ", *config.RPCCert, " ", *config.RPCKey)
	certDir, _ := filepath.Split(config.RPCCert.V())
	keyDir, _ := filepath.Split(config.RPCKey.V())
	e := os.MkdirAll(certDir, 0700)
	if e != nil {
		return tls.Certificate{}, e
	}
	e = os.MkdirAll(keyDir, 0700)
	if e != nil {
		return tls.Certificate{}, e
	}
	// Generate cert pair.
	org := "pod/wallet autogenerated cert"
	validUntil := time.Now().Add(time.Hour * 24 * 365 * 10)
	cert, key, e := util.NewTLSCertPair(org, validUntil, nil)
	if e != nil {
		return tls.Certificate{}, e
	}
	keyPair, e := tls.X509KeyPair(cert, key)
	if e != nil {
		return tls.Certificate{}, e
	}
	// Write cert and (potentially) the key files.
	e = ioutil.WriteFile(config.RPCCert.V(), cert, 0600)
	if e != nil {
		rmErr := os.Remove(config.RPCCert.V())
		if rmErr != nil {
			E.Ln("cannot remove written certificates:", rmErr)
		}
		return tls.Certificate{}, e
	}
	e = ioutil.WriteFile(config.CAFile.V(), cert, 0600)
	if e != nil {
		rmErr := os.Remove(config.RPCCert.V())
		if rmErr != nil {
			E.Ln("cannot remove written certificates:", rmErr)
		}
		return tls.Certificate{}, e
	}
	if writeKey {
		e = ioutil.WriteFile(config.RPCKey.V(), key, 0600)
		if e != nil {
			rmErr := os.Remove(config.RPCCert.V())
			if rmErr != nil {
				E.Ln("cannot remove written certificates:", rmErr)
			}
			rmErr = os.Remove(config.CAFile.V())
			if rmErr != nil {
				E.Ln("cannot remove written certificates:", rmErr)
			}
			return tls.Certificate{}, e
		}
	}
	I.Ln("done generating TLS certificates")
	return keyPair, nil
}

// makeListeners splits the normalized listen addresses into IPv4 and IPv6 addresses and creates new net.Listeners for
// each with the passed listen func. Invalid addresses are logged and skipped.
func makeListeners(normalizedListenAddrs []string, listen listenFunc) []net.Listener {
	ipv4Addrs := make([]string, 0, len(normalizedListenAddrs)*2)
	// ipv6Addrs := make([]string, 0, len(normalizedListenAddrs)*2)
	for _, addr := range normalizedListenAddrs {
		var host string
		var e error
		host, _, e = net.SplitHostPort(addr)
		if e != nil {
			// Shouldn't happen due to already being normalized.
			E.F(
				"`%s` is not a normalized listener address", addr,
			)
			continue
		}
		// Empty host or host of * on plan9 is both IPv4 and IPv6.
		if host == "" || (host == "*" && runtime.GOOS == "plan9") {
			ipv4Addrs = append(ipv4Addrs, addr)
			// ipv6Addrs = append(ipv6Addrs, addr)
			continue
		}
		// Remove the IPv6 zone from the host, if present. The zone prevents ParseIP from correctly parsing the IP
		// address. ResolveIPAddr is intentionally not used here due to the possibility of leaking a DNS query over Tor
		// if the host is a hostname and not an IP address.
		zoneIndex := strings.Index(host, "%")
		if zoneIndex != -1 {
			host = host[:zoneIndex]
		}
		ip := net.ParseIP(host)
		switch {
		case ip == nil:
			W.F("`%s` is not a valid IP address", host)
		case ip.To4() == nil:
			// ipv6Addrs = append(ipv6Addrs, addr)
		default:
			ipv4Addrs = append(ipv4Addrs, addr)
		}
	}
	listeners := make(
		[]net.Listener, 0,
		// len(ipv6Addrs)+
		len(ipv4Addrs),
	)
	for _, addr := range ipv4Addrs {
		listener, e := listen("tcp4", addr)
		if e != nil {
			W.F(
				"Can't listen on %s: %v", addr, e,
			)
			continue
		}
		listeners = append(listeners, listener)
	}
	// for _, addr := range ipv6Addrs {
	// 	listener, e := listen("tcp6", addr)
	// 	if e != nil  {
	// 		Warnf(
	// 			"Can't listen on %s: %v", addr, e,
	// 		)
	// 		continue
	// 	}
	// 	listeners = append(listeners, listener)
	// }
	return listeners
}

// OpenRPCKeyPair creates or loads the RPC TLS keypair specified by the
// application config. This function respects the pod.Config.OneTimeTLSKey
// setting.
func OpenRPCKeyPair(config *opts.Config) (tls.Certificate, error) {
	// Chk for existence of the TLS key file. If one time TLS keys are enabled but a
	// key already exists, this function should error since it's possible that a
	// persistent certificate was copied to a remote machine. Otherwise, generate a
	// new keypair when the key is missing. When generating new persistent keys,
	// overwriting an existing cert is acceptable if the previous execution used a
	// one time TLS key. Otherwise, both the cert and key should be read from disk.
	// If the cert is missing, the read error will occur in LoadX509KeyPair.
	_, e := os.Stat(config.RPCKey.V())
	keyExists := !os.IsNotExist(e)
	switch {
	case config.OneTimeTLSKey.True() && keyExists:
		e := fmt.Errorf(
			"one time TLS keys are enabled, "+
				"but TLS key `%s` already exists", *config.RPCKey,
		)
		return tls.Certificate{}, e
	case config.OneTimeTLSKey.True():
		return GenerateRPCKeyPair(config, false)
	case !keyExists:
		return GenerateRPCKeyPair(config, true)
	default:
		return tls.LoadX509KeyPair(config.RPCCert.V(), config.RPCKey.V())
	}
}
func startRPCServers(cx *pod.State, walletLoader *wallet.Loader) (*walletrpc2.Server, error) {
	T.Ln("startRPCServers")
	var (
		legacyServer *walletrpc2.Server
		walletListen = net.Listen
		keyPair      tls.Certificate
		e            error
	)
	if !cx.Config.TLS.True() {
		I.Ln("server TLS is disabled - only legacy RPC may be used")
	} else {
		keyPair, e = OpenRPCKeyPair(cx.Config)
		if e != nil {
			return nil, e
		}
		// Change the standard net.Listen function to the tls one.
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{keyPair},
			MinVersion:         tls.VersionTLS12,
			NextProtos:         []string{"h2"}, // HTTP/2 over TLS
			InsecureSkipVerify: cx.Config.TLSSkipVerify.True(),
		}
		walletListen = func(net string, laddr string) (net.Listener, error) {
			return tls.Listen(net, laddr, tlsConfig)
		}
	}
	if cx.Config.Username.V() == "" || cx.Config.Password.V() == "" {
		I.Ln("legacy RPC server disabled (requires username and password)")
	} else if len(cx.Config.WalletRPCListeners.S()) != 0 {
		listeners := makeListeners(cx.Config.WalletRPCListeners.S(), walletListen)
		if len(listeners) == 0 {
			e := errors.New("failed to create listeners for legacy RPC server")
			return nil, e
		}
		opts := walletrpc2.Options{
			Username:            cx.Config.Username.V(),
			Password:            cx.Config.Password.V(),
			MaxPOSTClients:      int64(cx.Config.WalletRPCMaxClients.V()),
			MaxWebsocketClients: int64(cx.Config.WalletRPCMaxWebsockets.V()),
		}
		legacyServer = walletrpc2.NewServer(&opts, walletLoader, listeners, nil)
	}
	// Error when no legacy RPC servers can be started.
	if legacyServer == nil {
		return nil, errors.New("no suitable RPC services can be started")
	}
	return legacyServer, nil
}

// startWalletRPCServices associates each of the (optionally-nil) RPC servers with a wallet to enable remote wallet
// access. For the legacy JSON-RPC server it enables methods that require a loaded wallet.
func startWalletRPCServices(wallet *wallet.Wallet, legacyServer *walletrpc2.Server) {
	if legacyServer != nil {
		D.Ln("starting legacy wallet rpc server")
		legacyServer.RegisterWallet(wallet)
	}
}
