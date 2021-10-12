// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
)

const (
	optionNameDataDir      = "data-dir"
	optionNameSwapEndpoint = "swap-endpoint"
)

func init() {
	cobra.EnableCommandSorting = false
}

type command struct {
	root *cobra.Command
}

type option func(*command)

func newCommand(opts ...option) (c *command, err error) {
	c = &command{
		root: &cobra.Command{
			Use:           "sana",
			Short:         "Ethereum Sana tools",
			SilenceErrors: true,
			SilenceUsage:  true,
		},
	}

	for _, o := range opts {
		o(c)
	}

	c.initNonceCmd()
	c.initChequebookCmd()
	c.initVersionCmd()

	return c, nil
}

func (c *command) Execute() (err error) {
	return c.root.Execute()
}

// Execute parses command line arguments and runs appropriate functions.
func Execute() (err error) {
	c, err := newCommand()
	if err != nil {
		return err
	}
	return c.Execute()
}
