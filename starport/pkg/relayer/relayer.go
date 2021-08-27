package relayer

import (
	"context"
	"io"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	tsrelayer "github.com/tendermint/starport/starport/pkg/nodetime/programs/ts-relayer"
	relayerconfig "github.com/tendermint/starport/starport/pkg/relayer/config"
	"github.com/tendermint/starport/starport/pkg/xurl"
)

type Relayer struct {
	ca cosmosaccount.Registry
}

func NewRelayer(ca cosmosaccount.Registry) Relayer {
	r := Relayer{
		ca: ca,
	}

	return r
}

func fixRPCAddress(rpcAddress string) string {
	return strings.TrimSuffix(xurl.HTTPEnsurePort(rpcAddress), "/")
}

func rpcClient(rpcAddress string) (*rpchttp.HTTP, error) {
	rpcAddress = fixRPCAddress(rpcAddress)
	return rpchttp.New(rpcAddress, "/websocket")
}

func clientCtx(rpcAddress string) (client.Context, error) {
	rpcClient, err := rpcClient(rpcAddress)
	if err != nil {
		return nil, err
	}
	cc := cosmosclient.NewContext(rpcClient, io.Discard, "", "")
	return cc, nil
}

type LinkStatus struct {
	ID       string `json:"pathName"`
	ErrorMsg string `json:"error"`
}

// Link links all chains that has a path to each other.
// paths are optional and acts as a filter to only link some chains.
// calling Link multiple times for the same paths does not have any side effects.
func (r Relayer) Link(ctx context.Context, pathIDs ...string) (err error) {
	conf, err := relayerconfig.Get()
	if err != nil {
		return err
	}

	for _, id := range pathIDs {
		path, err := conf.PathByID(id)
		if err != nil {
			return err
		}

		if path.Src.ChannelID != "" { // already linked.
			continue
		}

		chainSrc, err := conf.ChainByID(path.Src.ID)
		if err != nil {
			return err
		}

		chainDst, err := conf.ChainByID(path.Dst.ID)
		if err != nil {
			return err
		}

		var reply interface{}
		tsrelayer.Call(ctx, "link", []interface{}{chainSrc, chainDst}, &reply)

	}

	return
}

func (r Relayer) checkBalance(rpcAddress, account, addressPrefix string) error {
	rpcAddress = strings.TrimSuffix(xurl.HTTPEnsurePort(rpcAddress), "/")
	rpcClient, err := rpchttp.New(rpcAddress, "/websocket")
	if err != nil {
		return nil, err
	}

	acc, err := r.ca.GetByName(account)
	if err != nil {
		return err
	}
	addr := acc.Address(chainSrc.options.AddressPrefix)

	return nil
}

// Start relays tx packets for paths until ctx is canceled.
func (r Relayer) Start(ctx context.Context, paths ...string) error {
	var reply interface{}
	return tsrelayer.Call(ctx, "start", []interface{}{paths}, &reply)
}

// Path represents a path between two chains.
type Path struct {
	// ID is id of the path.
	ID string `json:"id"`

	// IsLinked indicates that chains of these paths are linked or not.
	IsLinked bool `json:"isLinked"`

	// Src end of the path.
	Src PathEnd `json:"src"`

	// Dst end of the path.
	Dst PathEnd `json:"dst"`
}

// PathEnd represents the chain at one side of a Path.
type PathEnd struct {
	ChannelID string `json:"channelID"`
	ChainID   string `json:"chainID"`
	PortID    string `json:"portID"`
}

// GetPath returns a path by its id.
func (r Relayer) GetPath(ctx context.Context, id string) (Path, error) {
	var path Path
	err := tsrelayer.Call(ctx, "getPath", []interface{}{id}, &path)
	return path, err
}

// ListPaths list all the paths.
func (r Relayer) ListPaths(ctx context.Context) ([]Path, error) {
	var paths []Path
	err := tsrelayer.Call(ctx, "listPaths", nil, &paths)
	return paths, err
}
