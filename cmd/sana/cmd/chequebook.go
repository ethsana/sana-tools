package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethersphere/bee/pkg/node"
	"github.com/ethersphere/bee/pkg/storage"
	"github.com/ethsana/sana/pkg/logging"
	"github.com/spf13/cobra"
)

const (
	chequebookKey           = "swap_chequebook"
	chequebookDeploymentKey = "swap_chequebook_transaction_deployment"
	deployedTopic           = `0xc0ffc525a1c7689549d7f79b49eca900e61ac49b43d977f680bcc3b36224c004`
)

func (c *command) initChequebookCmd() {
	cmd := &cobra.Command{
		Use:   "chequebook txhash",
		Short: "Repair the checkbook",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) < 1 {
				return cmd.Help()
			}

			dataDir, err := cmd.Flags().GetString(optionNameDataDir)
			if err != nil {
				return fmt.Errorf("get data-dir: %v", err)
			}
			if dataDir == "" {
				return errors.New("no data-dir provided")
			}

			logger := logging.New(ioutil.Discard, 0)

			stateStore, err := node.InitStateStore(logger, dataDir)
			if err != nil {
				fmt.Printf(`init statestore fail: %v`, err)
				return
			}
			defer stateStore.Close()

			var chequebook common.Address
			err = stateStore.Get(chequebookKey, &chequebook)
			if err != nil && err != storage.ErrNotFound {
				fmt.Printf(`get chequebook fail: %v`, err)
				return
			}

			if err == storage.ErrNotFound {
				endpoint, err := cmd.Flags().GetString(optionNameSwapEndpoint)
				if err != nil {
					return fmt.Errorf("get swap-endpoint: %v", err)
				}
				if endpoint == "" {
					return errors.New("no swap-endpoint provided")
				}

				client, err := ethclient.Dial(endpoint)
				if err != nil {
					return fmt.Errorf("ethclient dail fail: %v", err)
				}

				receipt, err := client.TransactionReceipt(context.TODO(), common.HexToHash(args[0]))
				if err != nil {
					return fmt.Errorf(`get transaction receipt fail: %v`, err)
				}

				for _, l := range receipt.Logs {
					if l.Topics[0].Hex() == deployedTopic {
						chequebook = common.BytesToAddress(l.Data)
						break
					}
				}

				if (chequebook == common.Address{}) {
					return fmt.Errorf(`not found chequebook with transaction %v`, args[0])
				}

				err = stateStore.Put(chequebookKey, chequebook)
				if err != nil {
					return fmt.Errorf(`put chequebook fail: %v`, err)
				}

				err = stateStore.Put(chequebookDeploymentKey, common.HexToHash(os.Args[2]))
				if err != nil {
					return fmt.Errorf(`put chequebook deploy transaction fail: %v`, err)
				}
			}
			fmt.Printf("chequebook: %v\n", chequebook.Hex())
			return nil
		},
	}

	cmd.Flags().String(optionNameDataDir, "", "data directory")
	cmd.Flags().String(optionNameSwapEndpoint, "http://127.0.0.1:8545", "swap ethereum blockchain endpoint")
	c.root.AddCommand(cmd)
}
