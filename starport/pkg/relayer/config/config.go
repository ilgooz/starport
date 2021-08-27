package relayerconf

import (
	"fmt"
	"os"
	"reflect"

	"github.com/tendermint/starport/starport/pkg/confile"
)

const supportVersion = "2"

var configPath = os.ExpandEnv("$HOME/.starport/relayer/config.yml")

type Config struct {
	Version string  `yaml:"version"`
	Chains  []Chain `yaml:"chains,omitempty"`
	Paths   []Path  `yaml:"paths,omitempty"`
}

func (c Config) ChainByID(id string) (Chain, error) {
	for _, chain := range c.Chains {
		if chain.ID == id {
			return chain, nil
		}
	}
	return Chain{}, fmt.Errorf("chain %q cannot be found", id)
}

func (c Config) PathByID(id string) (Path, error) {
	for _, path := range c.Paths {
		if path.ID == id {
			return path, nil
		}
	}
	return Path{}, fmt.Errorf("path %q cannot be found", id)
}

type Chain struct {
	ID            string `yaml:"id"`
	Account       string `yaml:"account"`
	AddressPrefix string `yaml:"address_prefix"`
	RPCAddress    string `yaml:"rpc_address"`
	GasPrice      string `yaml:"gas_price,omitempty"`
	GasLimit      int64  `yaml:"gas_limit,omitempty"`
}

type Path struct {
	ID       string  `yaml:"id"`
	Ordering string  `yaml:"ordering,omitempty"`
	Src      PathEnd `yaml:"src"`
	Dst      PathEnd `yaml:"dst"`
}

type PathEnd struct {
	ChainID      string `yaml:"chain_id"`
	ConnectionID string `yaml:"connection_id,omitempty"`
	ChannelID    string `yaml:"channel_id,omitempty"`
	PortID       string `yaml:"port_id"`
	Version      string `yaml:"version,omitempty"`
	PacketHeight int64  `yaml:"packet_height,omitempty"`
	AckHeight    int64  `yaml:"ack_height,omitempty"`
}

func Get() (Config, error) {
	c := Config{}
	if err := confile.New(confile.DefaultYAMLEncodingCreator, configPath).Load(&c); err != nil {
		return c, err
	}
	if !reflect.DeepEqual(c, Config{}) && c.Version != supportVersion {
		return c, fmt.Errorf("your relayer setup is outdated. run 'rm %s' and configure relayer again", configPath)
	}
	return c, nil
}

func Save(c Config) error {
	c.Version = supportVersion
	return confile.New(confile.DefaultYAMLEncodingCreator, configPath).Save(c)
}
