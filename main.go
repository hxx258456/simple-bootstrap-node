package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

var logger = log.Logger("bootsrap")

// var
var prvKey string

func main() {
	log.SetAllLoggers(log.LevelInfo)
	log.SetLogLevel("bootsrap", "info")
	help := flag.Bool("help", false, "Display Help")
	listenHost := flag.String("host", "0.0.0.0", "The bootstrap node host listen address\n")
	port := flag.Int("port", 4001, "The bootstrap node listen port")
	flag.Parse()

	if *help {
		fmt.Printf("This is a simple bootstrap node for kad-dht application using libp2p\n\n")
		fmt.Printf("Usage: \n   Run './bootnode'\nor Run './bootnode -host [host] -port [port]'\n")

		os.Exit(0)
	}

	fmt.Printf("[*] Listening on: %s with port: %d\n", *listenHost, *port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if prvKey == "" {
		logger.Info("No private key provided. Generating a new one...")
		// Creates a new RSA key pair for this host.
		host, err := libp2p.New()

		if err != nil {
			panic(err)
		}
		//节点私钥
		keyBytes, err := crypto.MarshalPrivateKey(host.Peerstore().PrivKey(host.ID()))
		if err != nil {
			panic(err)
		}

		prvKey = string(keyBytes)
		WriteConfig()
	}

	// Unmarshal Private Key
	privateKey, err := crypto.UnmarshalPrivateKey([]byte(prvKey))
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", *listenHost, *port))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	opts := []libp2p.Option{
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(privateKey),
		// libp2p.NoSecurity,
	}
	host, err := libp2p.New(
		opts...)

	if err != nil {
		panic(err)
	}
	logger.Info("Host created. We are:", host.ID())
	logger.Info(host.Addrs())

	// dht.ProtocolPrefix("primihub")
	_, err = dht.New(ctx, host,
		dht.Mode(dht.ModeServer))
	// dht.ProtocolPrefix("primihub"))
	if err != nil {
		panic(err)
	}
	// err1 := kad.Bootstrap(ctx)
	// if err1 != nil {
	// 	panic(err1)
	// }
	fmt.Println("")
	logger.Info("Your Bootstrap Address is: ", host.Addrs())
	fmt.Println("")
	select {}
}
