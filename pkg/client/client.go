package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"dogecoin.org/balance-master/pkg/rpc"
	"dogecoin.org/balance-master/pkg/store"
)

type BalanceMasterClientConfig struct {
	RpcServerHost string
	RpcServerPort string
}

type BalanceMasterClient struct {
	cfg *BalanceMasterClientConfig
}

func NewBalanceMasterClient(cfg *BalanceMasterClientConfig) *BalanceMasterClient {
	return &BalanceMasterClient{cfg: cfg}
}

func (c *BalanceMasterClient) TrackAddress(address string) error {
	request := rpc.PostTrackersRequest{
		Address: address,
	}
	json, err := json.Marshal(request)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s:%s/trackers", c.cfg.RpcServerHost, c.cfg.RpcServerPort), "application/json", bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to track address: %s", resp.Status)
	}

	return nil
}

func (c *BalanceMasterClient) GetUtxos(address string) ([]store.UTXO, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/utxos?address=%s", c.cfg.RpcServerHost, c.cfg.RpcServerPort, address))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get utxos: %s", resp.Status)
	}

	var utxos []store.UTXO
	err = json.NewDecoder(resp.Body).Decode(&utxos)
	if err != nil {
		return nil, err
	}

	return utxos, nil
}
