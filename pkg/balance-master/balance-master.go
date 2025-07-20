package balance_master

import (
	"log"

	bmcfg "dogecoin.org/balance-master/pkg/config"
	bmrpc "dogecoin.org/balance-master/pkg/rpc"
	"dogecoin.org/balance-master/pkg/store"
	"github.com/dogecoinfoundation/chainfollower/pkg/chainfollower"
	cfcfg "github.com/dogecoinfoundation/chainfollower/pkg/config"
	"github.com/dogecoinfoundation/chainfollower/pkg/messages"
	"github.com/dogecoinfoundation/chainfollower/pkg/rpc"
	"github.com/dogecoinfoundation/chainfollower/pkg/state"
	"github.com/golang-migrate/migrate"
)

type BalanceMaster struct {
	cfg           *bmcfg.Config
	rpcClient     *rpc.RpcTransport
	bmStore       *store.Store
	chainfollower *chainfollower.ChainFollower
	rpcServer     *bmrpc.RpcServer
}

func NewBalanceMaster(cfg *bmcfg.Config) *BalanceMaster {
	rpcClient := rpc.NewRpcTransport(&cfcfg.Config{
		RpcUrl:  cfg.DogeScheme + "://" + cfg.DogeHost + ":" + cfg.DogePort,
		RpcUser: cfg.DogeUser,
		RpcPass: cfg.DogePassword,
	})

	bmStore, err := store.NewStore(cfg.DatabaseURL, cfg.MigrationsPath)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	err = bmStore.Migrate()
	if err != nil && err.Error() != migrate.ErrNoChange.Error() {
		log.Fatalf("Failed to migrate store: %v", err)
	}

	rpcServer := bmrpc.NewRpcServer(cfg, bmStore)

	chainfollower := chainfollower.NewChainFollower(rpcClient)

	return &BalanceMaster{cfg: cfg, rpcClient: rpcClient, bmStore: bmStore, chainfollower: chainfollower, rpcServer: rpcServer}
}

func (bm *BalanceMaster) Start() {
	go bm.rpcServer.Start()

	blockHeight, blockHash, _, err := bm.bmStore.GetChainPosition()
	if err != nil {
		log.Fatalf("Failed to get chain position: %v", err)
	}

	msgChan := bm.chainfollower.Start(&state.ChainPos{
		BlockHeight: blockHeight,
		BlockHash:   blockHash,
	})

	for msg := range msgChan {
		switch msg := msg.(type) {
		case messages.BlockMessage:
			for _, tx := range msg.Block.Tx {
				inputs := []store.InputRef{}
				outputs := []store.UTXO{}
				for _, vin := range tx.VIn {
					inputs = append(inputs, store.InputRef{TxID: vin.TxID, VOut: uint32(vin.VOut)})
				}
				for _, vout := range tx.VOut {
					if len(vout.ScriptPubKey.Addresses) > 0 {
						log.Printf("Checking if address is tracking: %s", vout.ScriptPubKey.Addresses[0])

						isTracking, err := bm.bmStore.IsTracking(vout.ScriptPubKey.Addresses[0])
						if err != nil {
							log.Printf("Failed to check if address is tracking: %v", err)
						}

						if isTracking {
							outputs = append(outputs, store.UTXO{TxID: tx.TxID, VOut: int(vout.N), Address: vout.ScriptPubKey.Addresses[0], Amount: vout.Value.InexactFloat64()})
						}
					}
				}
				err := bm.bmStore.UpdateUtxos(inputs, outputs)
				if err != nil {
					log.Printf("Failed to update utxos: %v", err)
				}
			}
		}
	}
}

func (bm *BalanceMaster) Stop() {
	bm.rpcServer.Stop()
	bm.bmStore.Close()
	bm.chainfollower.Stop()
}
