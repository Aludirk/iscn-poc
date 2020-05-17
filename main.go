package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/plugin/plugins/cosmosds"
	"github.com/ipfs/go-ipfs/repo/fsrepo"

	cosmos "github.com/cosmos/cosmos-sdk/types"
	config "github.com/ipfs/go-ipfs-config"
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	icore "github.com/ipfs/interface-go-ipfs-core"
	abci "github.com/tendermint/tendermint/abci/types"
	tlog "github.com/tendermint/tendermint/libs/log"
)

func setupCosmosStore(plugins *loader.PluginLoader) cosmos.CommitMultiStore {
	pl, err := plugins.GetPlugin("ds-cosmos")
	if err != nil {
		log.Panicf("Cannot retrieve \"ds-cosmos\" plugin: %s", err)
	}

	cosmosDSPlugin, ok := pl.(*cosmosds.Plugin)
	if !ok {
		log.Panic("The plugin is not a \"*cosmosds.Plugin\"")
	}

	dataDir := filepath.Join(".", "cosmos")
	db, err := cosmos.NewLevelDB("application", dataDir)
	if err != nil {
		log.Panicf("Failed to create LevelDB: %s", err)
	}

	key := cosmos.NewKVStoreKey("StoreKey")
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, cosmos.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := cosmos.NewContext(cms, abci.Header{}, false, tlog.NewNopLogger())
	kv := ctx.KVStore(key)

	err = cosmosDSPlugin.SetCosmosStore(kv)
	if err != nil {
		log.Panicf("Cannot setup Cosmos store: %s", err)
	}

	return cms
}

func setupDefaultDatastoreConfig(cfg *config.Config) *config.Config {
	cfg.Datastore.Spec = map[string]interface{}{
		"mountpoint": "/",
		"type":       "measure",
		"prefix":     "cosmossdk.datastore",
		"child": map[string]interface{}{
			"type":        "cosmosds",
			"path":        "datastore",
			"compression": "none",
		},
	}
	return cfg
}

func setupNode(ctx context.Context) (
	*loader.PluginLoader,
	icore.CoreAPI,
	cosmos.CommitMultiStore,
	error) {
	rootPath, err := filepath.Abs("./ipfs")
	if err != nil {
		log.Printf("Cannot parse path: %s", err)
		return nil, nil, nil, err
	}

	plugins, err := loader.NewPluginLoader(rootPath)
	if err != nil {
		log.Printf("error loading plugins: %s", err)
		return nil, nil, nil, err
	}

	if err := plugins.Initialize(); err != nil {
		log.Printf("error initializing plugins: %s", err)
		return nil, nil, nil, err
	}

	if err := plugins.Inject(); err != nil {
		log.Printf("error initializing plugins: %s", err)
		return nil, nil, nil, err
	}

	rootPath, err = config.Path(rootPath, "")
	if err != nil {
		log.Printf("Cannot set config path: %s", err)
		return nil, nil, nil, err
	}

	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		log.Printf("Cannot init config: %s", err)
		return nil, nil, nil, err
	}
	cfg = setupDefaultDatastoreConfig(cfg)

	err = fsrepo.Init(rootPath, cfg)
	if err != nil {
		log.Printf("Cannot init repo: %s", "ab")
		return nil, nil, nil, err
	}

	repo, err := fsrepo.Open(rootPath)
	if err != nil {
		log.Printf("Cannot open repo: %s", err)
		return nil, nil, nil, err
	}

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption,
		Repo:    repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		log.Printf("Cannot create node: %s", err)
		return nil, nil, nil, err
	}
	node.IsDaemon = true

	ipfs, err := coreapi.NewCoreAPI(node)
	if err != nil {
		log.Println("No IPFS repo available on the default path")
		return nil, nil, nil, err
	}

	log.Println("Start plugin")
	err = plugins.Start(node)
	if err != nil {
		log.Printf("Cannot start plugins: %s", err)
		return nil, nil, nil, err
	}

	cms := setupCosmosStore(plugins)

	return plugins, ipfs, cms, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("\nGet signal: \"%v\"", sig)
		done <- true
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Setting up IPFS node ...")
	plugins, ipfs, cms, err := setupNode(ctx)
	if err != nil {
		log.Panicf("Failed to set up IPFS node")
	}
	log.Println("IPFS node is created")

	entities := testEntity(ctx, ipfs)
	stakeholders := testStakeholders(ctx, ipfs, entities)
	content := testContent(ctx, ipfs)
	testIscnKernel(ctx, ipfs, stakeholders, content)

	cms.Commit()

	<-done
	log.Println("Close plugin")
	err = plugins.Close()
	if err != nil {
		log.Panicf("Cannot close plugins: %s", err)
	}
}
