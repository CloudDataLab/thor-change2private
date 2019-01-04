// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/inconshreveable/log15"
	"github.com/mattn/go-isatty"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/vechain/thor/api"
	"github.com/vechain/thor/cmd/thor/node"
	"github.com/vechain/thor/cmd/thor/solo"
	"github.com/vechain/thor/genesis"
	"github.com/vechain/thor/logdb"
	"github.com/vechain/thor/lvldb"
	"github.com/vechain/thor/state"
	"github.com/vechain/thor/thor"
	"github.com/vechain/thor/txpool"
	"gopkg.in/urfave/cli.v1"
	"github.com/vechain/thor/p2psrv"
	"github.com/vechain/thor/p2psrv/discv5"
	"github.com/ethereum/go-ethereum/p2p/netutil"
	"math/rand"
)

var (
	version   string
	gitCommit string
	gitTag    string
	//edit by sion
	newNode   *node.Node
	serverAddress *p2psrv.Server
	netAddress *discv5.Network
	//edit by sion
	log       = log15.New()

	defaultTxPoolOptions = txpool.Options{
		Limit:           10000,
		LimitPerAccount: 16,
		MaxLifetime:     20 * time.Minute,
	}
)

func fullVersion() string {
	versionMeta := "release"
	if gitTag == "" {
		versionMeta = "dev"
	}
	return fmt.Sprintf("%s-%s-%s", version, gitCommit, versionMeta)
}

func main() {
	app := cli.App{
		Version:   fullVersion(),
		Name:      "Thor",
		Usage:     "Node of VeChain Thor Network",
		Copyright: "2018 VeChain Foundation <https://vechain.org/>",
		Flags: []cli.Flag{
			networkFlag,
			configDirFlag,
			dataDirFlag,
			beneficiaryFlag,
			apiAddrFlag,
			apiCorsFlag,
			verbosityFlag,
			maxPeersFlag,
			p2pPortFlag,
			natFlag,
		},
		Action: defaultAction,
		Commands: []cli.Command{
			{
				Name:  "solo",
				Usage: "client runs in solo mode for test & dev",
				Flags: []cli.Flag{
					dataDirFlag,
					apiAddrFlag,
					apiCorsFlag,
					onDemandFlag,
					persistFlag,
					verbosityFlag,
				},
				Action: soloAction,
			},
			{
				Name:  "master-key",
				Usage: "import and export master key",
				Flags: []cli.Flag{
					configDirFlag,
					importMasterKeyFlag,
					exportMasterKeyFlag,
				},
				Action: masterKeyAction,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func defaultAction(ctx *cli.Context) error {
	exitSignal := handleExitSignal()

	defer func() { log.Info("exited") }()

	initLogger(ctx)
	gene := selectGenesis(ctx)
	instanceDir := makeInstanceDir(ctx, gene)

	mainDB := openMainDB(ctx, instanceDir)
	defer func() { log.Info("closing main database..."); mainDB.Close() }()

	logDB := openLogDB(ctx, instanceDir)
	defer func() { log.Info("closing log database..."); logDB.Close() }()

	chain := initChain(gene, mainDB, logDB)
	master := loadNodeMaster(ctx)

	txPool := txpool.New(chain, state.NewCreator(mainDB), defaultTxPoolOptions)
	defer func() { log.Info("closing tx pool..."); txPool.Close() }()

	p2pcom := newP2PComm(ctx, chain, txPool, instanceDir)
	//edit by sion
	serverAddress=p2pcom.p2pSrv

	apiSrv, apiURL := startAPIServer(ctx, api.New(chain, state.NewCreator(mainDB), txPool, logDB, p2pcom.comm))
	defer func() { log.Info("stopping API server..."); apiSrv.Close() }()

	printStartupMessage(gene, chain, master, instanceDir, apiURL)

	// by kasper
	// should be pre-set before p2pcom.start()
	//test4p2p(serverAddress)

	p2pcom.Start()
	defer p2pcom.Stop()

    //edit by sion
    newNode=node.New(
		master,
		chain,
		state.NewCreator(mainDB),
		logDB,
		txPool,
		filepath.Join(instanceDir, "tx.stash"),
		p2pcom.comm)

	go func() {
		for{
			restrict:=newNode.GetRestrictList()
			serverAddress.SetRestrict(restrict)
			serverAddress.DisconnectPeers()
			//serverAddress.PickPeerInfo()
			time.Sleep(15*time.Second)
		}
  }()
    // by kasper
    // test blockinterval value
    log.Debug("get block interval", "BlockInterval", newNode.GetBlockInterval())

	return newNode.Run(exitSignal)
}

//edit by sion
// pre-set the net restrict
func geneRestrict() *netutil.Netlist{
	var netrestricts []string
	netrestricts=make([]string,4)
	// add real ip network from thor testnet
	netrestricts[0]="128.1.34.86/32"
	netrestricts[1]="35.226.220.193/32"
	netrestricts[2]="49.51.195.89/32"
	netrestricts[3]="107.155.60.54/32"
	i:=rand.Intn(4)
	netrestrict := netrestricts[i]
	var restrictList *netutil.Netlist
	if netrestrict != "" {
		restrictList, _ = netutil.ParseNetlist(netrestrict)

	}
	log.Debug("Set Net Restrict", "value", *restrictList)

	return restrictList
}

// by kasper
// test for p2p control
func test4p2p(serverAddress *p2psrv.Server){
	go func() {
		//for{
		// TODO : consider that if a node does not connect all peers in the whitelist
		//
		//time.Sleep(5*time.Second)
		restrict := preSetNetRestrict4Inbound() // geneRestrict()
		serverAddress.SetRestrict(restrict)
		//serverAddress.SetNetRestrict(restrict) //useless
		serverAddress.DisconnectPeers()
		//}

	}()
}

// by kasper
// pre-set whitelist not change during running for testing outbound
func preSetNetRestrict4Outbound() *netutil.Netlist{
	netrestrict := "128.1.34.86/32, 35.226.220.193/32, 49.51.195.89/32, 107.155.60.54/32, 203.195.230.202/32"
	list, err := netutil.ParseNetlist(netrestrict)
	if err != nil {
		//Fatalf("Option %q: %v", NetrestrictFlag.Name, err)
		log.Error("Option ",  "Set NetRestrict: ", err)
		return nil
	}
	return list
}

// by kasper
// for inbound test
func preSetNetRestrict4Inbound() *netutil.Netlist{
	netrestrict := "120.78.83.87/32, 112.74.47.241/32"
	list, err := netutil.ParseNetlist(netrestrict)
	if err != nil {
		//Fatalf("Option %q: %v", NetrestrictFlag.Name, err)
		log.Error("Option ",  "Set NetRestrict: ", err)
		return nil
	}
	return list
}

func soloAction(ctx *cli.Context) error {
	defer func() { log.Info("exited") }()

	initLogger(ctx)
	gene := genesis.NewDevnet()

	var mainDB *lvldb.LevelDB
	var logDB *logdb.LogDB
	var instanceDir string

	if ctx.Bool("persist") {
		instanceDir = makeInstanceDir(ctx, gene)
		mainDB = openMainDB(ctx, instanceDir)
		logDB = openLogDB(ctx, instanceDir)
	} else {
		instanceDir = "Memory"
		mainDB = openMemMainDB()
		logDB = openMemLogDB()
	}

	defer func() { log.Info("closing main database..."); mainDB.Close() }()
	defer func() { log.Info("closing log database..."); logDB.Close() }()

	chain := initChain(gene, mainDB, logDB)

	txPool := txpool.New(chain, state.NewCreator(mainDB), defaultTxPoolOptions)
	defer func() { log.Info("closing tx pool..."); txPool.Close() }()

	soloContext := solo.New(chain, state.NewCreator(mainDB), logDB, txPool, ctx.Bool("on-demand"))

	apiSrv, apiURL := startAPIServer(ctx, api.New(chain, state.NewCreator(mainDB), txPool, logDB, solo.Communicator{}))
	defer func() { log.Info("stopping API server..."); apiSrv.Close() }()

	printSoloStartupMessage(gene, chain, instanceDir, apiURL)

	return soloContext.Run(handleExitSignal())
}

func masterKeyAction(ctx *cli.Context) error {
	hasImportFlag := ctx.Bool(importMasterKeyFlag.Name)
	hasExportFlag := ctx.Bool(exportMasterKeyFlag.Name)
	if hasImportFlag && hasExportFlag {
		return fmt.Errorf("flag %s and %s are exclusive", importMasterKeyFlag.Name, exportMasterKeyFlag.Name)
	}

	if !hasImportFlag && !hasExportFlag {
		return fmt.Errorf("missing flag, either %s or %s", importMasterKeyFlag.Name, exportMasterKeyFlag.Name)
	}

	if hasImportFlag {
		if isatty.IsTerminal(os.Stdin.Fd()) {
			fmt.Println("Input JSON keystore (end with ^d):")
		}
		keyjson, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(keyjson, &map[string]interface{}{}); err != nil {
			return errors.WithMessage(err, "unmarshal")
		}
		password, err := readPasswordFromNewTTY("Enter passphrase: ")
		if err != nil {
			return err
		}

		key, err := keystore.DecryptKey(keyjson, password)
		if err != nil {
			return errors.WithMessage(err, "decrypt")
		}

		if err := crypto.SaveECDSA(masterKeyPath(ctx), key.PrivateKey); err != nil {
			return err
		}
		fmt.Println("Master key imported:", thor.Address(key.Address))
		return nil
	}

	if hasExportFlag {
		masterKey, err := loadOrGeneratePrivateKey(masterKeyPath(ctx))
		if err != nil {
			return err
		}

		password, err := readPasswordFromNewTTY("Enter passphrase: ")
		if err != nil {
			return err
		}
		if password == "" {
			return errors.New("non-empty passphrase required")
		}
		confirm, err := readPasswordFromNewTTY("Confirm passphrase: ")
		if err != nil {
			return err
		}

		if password != confirm {
			return errors.New("passphrase confirmation mismatch")
		}

		keyjson, err := keystore.EncryptKey(&keystore.Key{
			PrivateKey: masterKey,
			Address:    crypto.PubkeyToAddress(masterKey.PublicKey),
			Id:         uuid.NewRandom()},
			password, keystore.StandardScryptN, keystore.StandardScryptP)
		if err != nil {
			return err
		}
		if isatty.IsTerminal(os.Stdout.Fd()) {
			fmt.Println("=== JSON keystore ===")
		}
		_, err = fmt.Println(string(keyjson))
		return err
	}
	return nil
}
