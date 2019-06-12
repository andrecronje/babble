package mobile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/andrecronje/babble/src/babble"
	"github.com/andrecronje/babble/src/crypto/keys"
	"github.com/andrecronje/babble/src/node"
	"github.com/andrecronje/babble/src/peers"
	"github.com/andrecronje/babble/src/proxy"
	"github.com/andrecronje/babble/src/proxy/inmem"
	"github.com/sirupsen/logrus"
)

type Node struct {
	nodeID uint32
	node   *node.Node
	proxy  proxy.AppProxy
	logger *logrus.Logger
}

// New initializes Node struct
func New(privKey string,
	nodeAddr string,
	jsonPeers string,
	commitHandler CommitHandler,
	exceptionHandler ExceptionHandler,
	config *MobileConfig) *Node {

	babbleConfig := config.toBabbleConfig()

	babbleConfig.Logger.WithFields(logrus.Fields{
		"nodeAddr": nodeAddr,
		"peers":    jsonPeers,
		"config":   fmt.Sprintf("%v", config),
	}).Debug("New Mobile Node")

	babbleConfig.BindAddr = nodeAddr

	//Check private key
	key, err := keys.ParsePrivateKey([]byte(privKey))
	if err != nil {
		exceptionHandler.OnException(fmt.Sprintf("Failed to parse private key: %s", err))
		return nil
	}

	babbleConfig.Key = key

	// Decode the peers
	var ps []*peers.Peer
	dec := json.NewDecoder(strings.NewReader(jsonPeers))
	if err := dec.Decode(&ps); err != nil {
		exceptionHandler.OnException(fmt.Sprintf("Failed to parse PeerSet: %s", err))
		return nil
	}

	peerSet := peers.NewPeerSet(ps)

	babbleConfig.LoadPeers = false

	//mobileApp implements the ProxyHandler interface, and we use it to
	//instantiate an InmemProxy
	mobileApp := newMobileApp(commitHandler, exceptionHandler, babbleConfig.Logger)

	babbleConfig.Proxy = inmem.NewInmemProxy(mobileApp, babbleConfig.Logger)

	engine := babble.NewBabble(babbleConfig)

	engine.Peers = peerSet

	if err := engine.Init(); err != nil {
		exceptionHandler.OnException(fmt.Sprintf("Cannot initialize engine: %s", err))
		return nil
	}

	return &Node{
		node:   engine.Node,
		proxy:  babbleConfig.Proxy,
		nodeID: engine.Node.GetID(),
		logger: babbleConfig.Logger,
	}
}

func (n *Node) Run(async bool) {
	if async {
		n.node.RunAsync(true)
	} else {
		n.node.Run(true)
	}
}

func (n *Node) Leave() {
	n.node.Leave()
}

func (n *Node) Shutdown() {
	n.node.Shutdown()
}

func (n *Node) SubmitTx(tx []byte) {
	//have to make a copy or the tx will be garbage collected and weird stuff
	//happens in transaction pool
	t := make([]byte, len(tx), len(tx))
	copy(t, tx)
	n.proxy.SubmitCh() <- t
}

func (n *Node) GetPeers() string {
	peers := n.node.GetPeers()

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(peers); err != nil {
		return ""
	}

	return buf.String()
}

func (n *Node) GetStats() string {
	stats := n.node.GetStats()

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(stats); err != nil {
		return ""
	}

	return buf.String()
}
