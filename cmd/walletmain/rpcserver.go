package walletmain

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	
	"github.com/p9c/pod/app/conte"
	
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/wallet"
)

type listenFunc func(net string, laddr string) (net.Listener, error)

// GenerateRPCKeyPair generates a new RPC TLS keypair and writes the cert and possibly also the key in PEM format to the
// paths specified by the config. If successful, the new keypair is returned.
func GenerateRPCKeyPair(config *pod.Config, writeKey bool) (tls.Certificate, error) {
	Debug("generating TLS certificates")
	// Create directories for cert and key files if they do not yet exist.
	Debug("rpc tls ", *config.RPCCert, " ", *config.RPCKey)
	certDir, _ := filepath.Split(*config.RPCCert)
	keyDir, _ := filepath.Split(*config.RPCKey)
	err := os.MkdirAll(certDir, 0700)
	if err != nil {
		Error(err)
		return tls.Certificate{}, err
	}
	err = os.MkdirAll(keyDir, 0700)
	if err != nil {
		Error(err)
		return tls.Certificate{}, err
	}
	// Generate cert pair.
	org := "pod/wallet autogenerated cert"
	validUntil := time.Now().Add(time.Hour * 24 * 365 * 10)
	cert, key, err := util.NewTLSCertPair(org, validUntil, nil)
	if err != nil {
		Error(err)
		return tls.Certificate{}, err
	}
	keyPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		Error(err)
		return tls.Certificate{}, err
	}
	// Write cert and (potentially) the key files.
	err = ioutil.WriteFile(*config.RPCCert, cert, 0600)
	if err != nil {
		rmErr := os.Remove(*config.RPCCert)
		if rmErr != nil {
			Error("cannot remove written certificates:", rmErr)
		}
		return tls.Certificate{}, err
	}
	err = ioutil.WriteFile(*config.CAFile, cert, 0600)
	if err != nil {
		rmErr := os.Remove(*config.RPCCert)
		if rmErr != nil {
			Error("cannot remove written certificates:", rmErr)
		}
		return tls.Certificate{}, err
	}
	if writeKey {
		err = ioutil.WriteFile(*config.RPCKey, key, 0600)
		if err != nil {
			Error(err)
			rmErr := os.Remove(*config.RPCCert)
			if rmErr != nil {
				Error("cannot remove written certificates:", rmErr)
			}
			rmErr = os.Remove(*config.CAFile)
			if rmErr != nil {
				Error("cannot remove written certificates:", rmErr)
			}
			return tls.Certificate{}, err
		}
	}
	Info("done generating TLS certificates")
	return keyPair, nil
}

// makeListeners splits the normalized listen addresses into IPv4 and IPv6 addresses and creates new net.Listeners for
// each with the passed listen func. Invalid addresses are logged and skipped.
func makeListeners(normalizedListenAddrs []string, listen listenFunc) []net.Listener {
	ipv4Addrs := make([]string, 0, len(normalizedListenAddrs)*2)
	// ipv6Addrs := make([]string, 0, len(normalizedListenAddrs)*2)
	for _, addr := range normalizedListenAddrs {
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			Error(err)
			// Shouldn't happen due to already being normalized.
			Errorf(
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
			Warnf("`%s` is not a valid IP address", host)
		case ip.To4() == nil:
			// ipv6Addrs = append(ipv6Addrs, addr)
		default:
			ipv4Addrs = append(ipv4Addrs, addr)
		}
	}
	listeners := make([]net.Listener, 0,
		// len(ipv6Addrs)+
		len(ipv4Addrs))
	for _, addr := range ipv4Addrs {
		listener, err := listen("tcp4", addr)
		if err != nil {
			Error(err)
			Warnf(
				"Can't listen on %s: %v", addr, err,
			)
			continue
		}
		listeners = append(listeners, listener)
	}
	// for _, addr := range ipv6Addrs {
	// 	listener, err := listen("tcp6", addr)
	// 	if err != nil {
	// 		Warnf(
	// 			"Can't listen on %s: %v", addr, err,
	// 		)
	// 		continue
	// 	}
	// 	listeners = append(listeners, listener)
	// }
	return listeners
}

// OpenRPCKeyPair creates or loads the RPC TLS keypair specified by the application config. This function respects the
// pod.Config.OneTimeTLSKey setting.
func OpenRPCKeyPair(config *pod.Config) (tls.Certificate, error) {
	// Check for existence of the TLS key file. If one time TLS keys are enabled but a key already exists, this function
	// should error since it's possible that a persistent certificate was copied to a remote machine. Otherwise,
	// generate a new keypair when the key is missing. When generating new persistent keys, overwriting an existing cert
	// is acceptable if the previous execution used a one time TLS key. Otherwise, both the cert and key should be read
	// from disk. If the cert is missing, the read error will occur in LoadX509KeyPair.
	_, e := os.Stat(*config.RPCKey)
	keyExists := !os.IsNotExist(e)
	switch {
	case *config.OneTimeTLSKey && keyExists:
		err := fmt.Errorf(
			"one time TLS keys are enabled, "+
				"but TLS key `%s` already exists", *config.RPCKey,
		)
		return tls.Certificate{}, err
	case *config.OneTimeTLSKey:
		return GenerateRPCKeyPair(config, false)
	case !keyExists:
		return GenerateRPCKeyPair(config, true)
	default:
		return tls.LoadX509KeyPair(*config.RPCCert, *config.RPCKey)
	}
}
func startRPCServers(cx *conte.Xt, walletLoader *wallet.Loader) (*legacy.Server, error) {
	Trace("startRPCServers")
	var (
		legacyServer *legacy.Server
		walletListen = net.Listen
		keyPair      tls.Certificate
		err          error
	)
	if !*cx.Config.TLS {
		Info("server TLS is disabled - only legacy RPC may be used")
	} else {
		keyPair, err = OpenRPCKeyPair(cx.Config)
		if err != nil {
			Error(err)
			return nil, err
		}
		// Change the standard net.Listen function to the tls one.
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{keyPair},
			MinVersion:         tls.VersionTLS12,
			NextProtos:         []string{"h2"}, // HTTP/2 over TLS
			InsecureSkipVerify: *cx.Config.TLSSkipVerify,
		}
		walletListen = func(net string, laddr string) (net.Listener, error) {
			return tls.Listen(net, laddr, tlsConfig)
		}
	}
	if *cx.Config.Username == "" || *cx.Config.Password == "" {
		Info("legacy RPC server disabled (requires username and password)")
	} else if len(*cx.Config.WalletRPCListeners) != 0 {
		listeners := makeListeners(*cx.Config.WalletRPCListeners, walletListen)
		if len(listeners) == 0 {
			err := errors.New("failed to create listeners for legacy RPC server")
			return nil, err
		}
		opts := legacy.Options{
			Username:            *cx.Config.Username,
			Password:            *cx.Config.Password,
			MaxPOSTClients:      int64(*cx.Config.WalletRPCMaxClients),
			MaxWebsocketClients: int64(*cx.Config.WalletRPCMaxWebsockets),
		}
		legacyServer = legacy.NewServer(&opts, walletLoader, listeners, nil)
	}
	// Error when no legacy RPC servers can be started.
	if legacyServer == nil {
		return nil, errors.New("no suitable RPC services can be started")
	}
	return legacyServer, nil
}

// startWalletRPCServices associates each of the (optionally-nil) RPC servers with a wallet to enable remote wallet
// access. For the legacy JSON-RPC server it enables methods that require a loaded wallet.
func startWalletRPCServices(wallet *wallet.Wallet, legacyServer *legacy.Server) {
	if legacyServer != nil {
		Debug("starting legacy wallet rpc server")
		legacyServer.RegisterWallet(wallet)
	}
}
