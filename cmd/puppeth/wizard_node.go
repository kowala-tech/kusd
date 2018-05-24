// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kowala-tech/kcoin/common"
	"github.com/kowala-tech/kcoin/log"
)

// deployNode creates a new node configuration based on some user input.
func (w *wizard) deployNode(boot bool) {
	// Do some sanity check before the user wastes time on input
	if w.conf.genesis == nil {
		log.Error("No genesis block configured")
		return
	}
	if w.conf.stats == "" {
		log.Error("No ethstats server configured")
		return
	}
	// Select the server to interact with
	server := w.selectServer()
	if server == "" {
		return
	}
	client := w.servers[server]

	// Retrieve any active ethstats configurations from the server
	infos, err := checkNode(client, w.network, boot)
	if err != nil {
		if boot {
			infos = &nodeInfos{portFull: 22334, peersTotal: 512, peersLight: 256}
		} else {
			infos = &nodeInfos{portFull: 22334, peersTotal: 50, peersLight: 0, gasTarget: 4.7, gasPrice: 18}
		}
	}
	infos.genesis, _ = json.MarshalIndent(w.conf.genesis, "", "  ")
	infos.network = w.conf.genesis.Config.ChainID.Int64()

	// Figure out where the user wants to store the persistent data
	fmt.Println()
	if infos.datadir == "" {
		fmt.Printf("Where should data be stored on the remote machine?\n")
		infos.datadir = w.readString()
	} else {
		fmt.Printf("Where should data be stored on the remote machine? (default = %s)\n", infos.datadir)
		infos.datadir = w.readDefaultString(infos.datadir)
	}
	// Figure out which port to listen on
	fmt.Println()
	fmt.Printf("Which TCP/UDP port to listen on? (default = %d)\n", infos.portFull)
	infos.portFull = w.readDefaultInt(infos.portFull)

	// Figure out how many peers to allow (different based on node type)
	fmt.Println()
	fmt.Printf("How many peers to allow connecting? (default = %d)\n", infos.peersTotal)
	infos.peersTotal = w.readDefaultInt(infos.peersTotal)

	// Figure out how many light peers to allow (different based on node type)
	fmt.Println()
	fmt.Printf("How many light peers to allow connecting? (default = %d)\n", infos.peersLight)
	infos.peersLight = w.readDefaultInt(infos.peersLight)

	// Set a proper name to report on the stats page
	fmt.Println()
	if infos.stats == "" {
		fmt.Printf("What should the node be called on the stats page?\n")
		infos.stats = w.readString() + ":" + w.conf.stats
	} else {
		fmt.Printf("What should the node be called on the stats page? (default = %s)\n", infos.stats)
		infos.stats = w.readDefaultString(infos.stats) + ":" + w.conf.stats
	}
	// If the node is a validator, load up needed credentials
	if !boot {
		if w.conf.genesis.Config.Tendermint != nil {
			// Tendermint based validators only need a coinbase to validate
			fmt.Println()
			if infos.coinbase == "" {
				fmt.Printf("What address should the validator user?\n")
				for {
					if address := w.readAddress(); address != nil {
						infos.coinbase = address.Hex()
						break
					}
				}
			} else {
				fmt.Printf("What address should the miner user? (default = %s)\n", infos.coinbase)
				infos.coinbase = w.readDefaultAddress(common.HexToAddress(infos.coinbase)).Hex()
			}
		}

		// Establish the gas dynamics to be enforced by the signer
		fmt.Println()
		fmt.Printf("What gas limit should empty blocks target (MGas)? (default = %0.3f)\n", infos.gasTarget)
		infos.gasTarget = w.readDefaultFloat(infos.gasTarget)

		fmt.Println()
		fmt.Printf("What gas price should the signer require (GWei)? (default = %0.3f)\n", infos.gasPrice)
		infos.gasPrice = w.readDefaultFloat(infos.gasPrice)
	}
	// Try to deploy the full node on the host
	if out, err := deployNode(client, w.network, w.conf.bootFull, w.conf.bootLight, infos); err != nil {
		log.Error("Failed to deploy Ethereum node container", "err", err)
		if len(out) > 0 {
			fmt.Printf("%s\n", out)
		}
		return
	}
	// All ok, run a network scan to pick any changes up
	log.Info("Waiting for node to finish booting")
	time.Sleep(3 * time.Second)

	w.networkStats(false)
}