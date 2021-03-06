package starportcmd

import (
	relayercmd "github.com/cosmos/relayer/cmd"
	"github.com/spf13/cobra"
)

// NewRelayer returns a new relayer command.
func NewRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:   "relayer",
		Short: "Connects blockchains via IBC protocol",
	}

	rlyCmd := relayercmd.NewRootCmd()
	rlyCmd.Short = "Low-level commands from github.com/cosmos/relayer"

	c.AddCommand(NewRelayerConfigure())
	c.AddCommand(NewRelayerConnect())
	c.AddCommand(NewRelayerLowLevel())
	c.AddCommand(rlyCmd)

	return c
}
