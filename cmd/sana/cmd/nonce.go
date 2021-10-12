package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ethsana/sana/pkg/logging"
	"github.com/ethsana/sana/pkg/statestore/leveldb"
	"github.com/spf13/cobra"
)

func (c *command) initNonceCmd() {
	cmd := &cobra.Command{
		Use:   "nonce",
		Short: "Reset nonce with Sana node",
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			dataDir, err := cmd.Flags().GetString(optionNameDataDir)
			if err != nil {
				return fmt.Errorf("get data-dir: %v", err)
			}
			if dataDir == "" {
				return errors.New("no data-dir provided")
			}

			store, err := leveldb.NewStateStore(filepath.Join(dataDir, "statestore"), logging.New(ioutil.Discard, 0))
			if err != nil {
				return fmt.Errorf("new statestore fail: %v", err)
			}
			defer store.Close()

			return store.Iterate("transaction_nonce_", func(key, _ []byte) (stop bool, err error) {
				fmt.Printf("remove key %s", string(key))
				return false, store.Delete(string(key))
			})
		},
	}

	cmd.Flags().String(optionNameDataDir, "", "data directory")
	c.root.AddCommand(cmd)
}
