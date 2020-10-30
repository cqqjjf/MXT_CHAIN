// Copyright 2019 The go-mxt Authors
// This file is part of the go-mxt library.
//
// The go-mxt library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-mxt library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-mxt library. If not, see <http://www.gnu.org/licenses/>.

package mxt

import (
	"github.com/mxt/go-mxt/core"
	"github.com/mxt/go-mxt/core/forkid"
	"github.com/mxt/go-mxt/p2p"
	"github.com/mxt/go-mxt/p2p/dnsdisc"
	"github.com/mxt/go-mxt/p2p/enode"
	"github.com/mxt/go-mxt/rlp"
)

// mxtEntry is the "mxt" ENR entry which advertises mxt protocol
// on the discovery network.
type mxtEntry struct {
	ForkID forkid.ID // Fork identifier per EIP-2124

	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

// ENRKey implements enr.Entry.
func (e mxtEntry) ENRKey() string {
	return "mxt"
}

// startEthEntryUpdate starts the ENR updater loop.
func (mxt *Ethereum) startEthEntryUpdate(ln *enode.LocalNode) {
	var newHead = make(chan core.ChainHeadEvent, 10)
	sub := mxt.blockchain.SubscribeChainHeadEvent(newHead)

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-newHead:
				ln.Set(mxt.currentEthEntry())
			case <-sub.Err():
				// Would be nice to sync with mxt.Stop, but there is no
				// good way to do that.
				return
			}
		}
	}()
}

func (mxt *Ethereum) currentEthEntry() *mxtEntry {
	return &mxtEntry{ForkID: forkid.NewID(mxt.blockchain.Config(), mxt.blockchain.Genesis().Hash(),
		mxt.blockchain.CurrentHeader().Number.Uint64())}
}

// setupDiscovery creates the node discovery source for the mxt protocol.
func (mxt *Ethereum) setupDiscovery(cfg *p2p.Config) (enode.Iterator, error) {
	if cfg.NoDiscovery || len(mxt.config.DiscoveryURLs) == 0 {
		return nil, nil
	}
	client := dnsdisc.NewClient(dnsdisc.Config{})
	return client.NewIterator(mxt.config.DiscoveryURLs...)
}
