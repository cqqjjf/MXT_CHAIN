package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/mxt/go-mxt/p2p"
	"github.com/mxt/go-mxt/rpc"
)

type gmxtrpc struct {
	name     string
	rpc      *rpc.Client
	gmxt     *testgmxt
	nodeInfo *p2p.NodeInfo
}

func (g *gmxtrpc) killAndWait() {
	g.gmxt.Kill()
	g.gmxt.WaitExit()
}

func (g *gmxtrpc) callRPC(result interface{}, mmxtod string, args ...interface{}) {
	if err := g.rpc.Call(&result, mmxtod, args...); err != nil {
		g.gmxt.Fatalf("callRPC %v: %v", mmxtod, err)
	}
}

func (g *gmxtrpc) addPeer(peer *gmxtrpc) {
	g.gmxt.Logf("%v.addPeer(%v)", g.name, peer.name)
	enode := peer.getNodeInfo().Enode
	peerCh := make(chan *p2p.PeerEvent)
	sub, err := g.rpc.Subscribe(context.Background(), "admin", peerCh, "peerEvents")
	if err != nil {
		g.gmxt.Fatalf("subscribe %v: %v", g.name, err)
	}
	defer sub.Unsubscribe()
	g.callRPC(nil, "admin_addPeer", enode)
	dur := 14 * time.Second
	timeout := time.After(dur)
	select {
	case ev := <-peerCh:
		g.gmxt.Logf("%v received event: type=%v, peer=%v", g.name, ev.Type, ev.Peer)
	case err := <-sub.Err():
		g.gmxt.Fatalf("%v sub error: %v", g.name, err)
	case <-timeout:
		g.gmxt.Error("timeout adding peer after", dur)
	}
}

// Use this function instead of `g.nodeInfo` directly
func (g *gmxtrpc) getNodeInfo() *p2p.NodeInfo {
	if g.nodeInfo != nil {
		return g.nodeInfo
	}
	g.nodeInfo = &p2p.NodeInfo{}
	g.callRPC(&g.nodeInfo, "admin_nodeInfo")
	return g.nodeInfo
}

func (g *gmxtrpc) waitSynced() {
	// Check if it's synced now
	var result interface{}
	g.callRPC(&result, "mxt_syncing")
	syncing, ok := result.(bool)
	if ok && !syncing {
		g.gmxt.Logf("%v already synced", g.name)
		return
	}

	// Actually wait, subscribe to the event
	ch := make(chan interface{})
	sub, err := g.rpc.Subscribe(context.Background(), "mxt", ch, "syncing")
	if err != nil {
		g.gmxt.Fatalf("%v syncing: %v", g.name, err)
	}
	defer sub.Unsubscribe()
	timeout := time.After(4 * time.Second)
	select {
	case ev := <-ch:
		g.gmxt.Log("'syncing' event", ev)
		syncing, ok := ev.(bool)
		if ok && !syncing {
			break
		}
		g.gmxt.Log("Other 'syncing' event", ev)
	case err := <-sub.Err():
		g.gmxt.Fatalf("%v notification: %v", g.name, err)
		break
	case <-timeout:
		g.gmxt.Fatalf("%v timeout syncing", g.name)
		break
	}
}

func startGmxtWithIpc(t *testing.T, name string, args ...string) *gmxtrpc {
	g := &gmxtrpc{name: name}
	args = append([]string{"--networkid=42", "--port=0", "--nousb"}, args...)
	t.Logf("Starting %v with rpc: %v", name, args)
	g.gmxt = runGmxt(t, args...)
	// wait before we can attach to it. TODO: probe for it properly
	time.Sleep(1 * time.Second)
	var err error
	ipcpath := filepath.Join(g.gmxt.Datadir, "gmxt.ipc")
	g.rpc, err = rpc.Dial(ipcpath)
	if err != nil {
		t.Fatalf("%v rpc connect: %v", name, err)
	}
	return g
}

func initGmxt(t *testing.T) string {
	g := runGmxt(t, "--nousb", "--networkid=42", "init", "./testdata/clique.json")
	datadir := g.Datadir
	g.WaitExit()
	return datadir
}

func startLightServer(t *testing.T) *gmxtrpc {
	datadir := initGmxt(t)
	runGmxt(t, "--nousb", "--datadir", datadir, "--password", "./testdata/password.txt", "account", "import", "./testdata/key.prv").WaitExit()
	account := "0x02f0d131f1f97aef08aec6e3291b957d9efe7105"
	server := startGmxtWithIpc(t, "lightserver", "--allow-insecure-unlock", "--datadir", datadir, "--password", "./testdata/password.txt", "--unlock", account, "--mine", "--light.serve=100", "--light.maxpeers=1", "--nodiscover", "--nat=extip:127.0.0.1")
	return server
}

func startClient(t *testing.T, name string) *gmxtrpc {
	datadir := initGmxt(t)
	return startGmxtWithIpc(t, name, "--datadir", datadir, "--nodiscover", "--syncmode=light", "--nat=extip:127.0.0.1")
}

func TestPriorityClient(t *testing.T) {
	lightServer := startLightServer(t)
	defer lightServer.killAndWait()

	// Start client and add lightServer as peer
	freeCli := startClient(t, "freeCli")
	defer freeCli.killAndWait()
	freeCli.addPeer(lightServer)

	var peers []*p2p.PeerInfo
	freeCli.callRPC(&peers, "admin_peers")
	if len(peers) != 1 {
		t.Errorf("Expected: # of client peers == 1, actual: %v", len(peers))
		return
	}

	// Set up priority client, get its nodeID, increase its balance on the lightServer
	prioCli := startClient(t, "prioCli")
	defer prioCli.killAndWait()
	// 3_000_000_000 once we move to Go 1.13
	tokens := 3000000000
	lightServer.callRPC(nil, "les_addBalance", prioCli.getNodeInfo().ID, tokens)
	prioCli.addPeer(lightServer)

	// Check if priority client is actually syncing and the regular client got kicked out
	prioCli.callRPC(&peers, "admin_peers")
	if len(peers) != 1 {
		t.Errorf("Expected: # of prio peers == 1, actual: %v", len(peers))
	}

	nodes := map[string]*gmxtrpc{
		lightServer.getNodeInfo().ID: lightServer,
		freeCli.getNodeInfo().ID:     freeCli,
		prioCli.getNodeInfo().ID:     prioCli,
	}
	time.Sleep(1 * time.Second)
	lightServer.callRPC(&peers, "admin_peers")
	peersWithNames := make(map[string]string)
	for _, p := range peers {
		peersWithNames[nodes[p.ID].name] = p.ID
	}
	if _, freeClientFound := peersWithNames[freeCli.name]; freeClientFound {
		t.Error("client is still a peer of lightServer", peersWithNames)
	}
	if _, prioClientFound := peersWithNames[prioCli.name]; !prioClientFound {
		t.Error("prio client is not among lightServer peers", peersWithNames)
	}
}
