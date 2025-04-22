package config

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type Config struct {
	HttpPort        int    `json:"http_port" toml:"http_port"`
	RpcPort         int    `json:"rpc_port" toml:"rpc_port"`
	ExecuteRpc      string `json:"execute_rpc" toml:"execute_rpc"`
	BeaconRpc       string `json:"beacon_rpc" toml:"beacon_rpc"`
	HonestBeaconRpc string `json:"honest_beacon_rpc" toml:"honest_beacon_rpc"`
	DbConnect       string `json:"dbconnect" toml:"dbconnect"`
	SwagHost        string `json:"swag_host" toml:"swag_host"`
	RewardFile      string `json:"reward_file" toml:"reward_file"`
}

var _cfg *Config = nil

func ParseConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("get config failed", "err", err)
		panic(err)
	}
	err = toml.Unmarshal(data, &_cfg)
	// err = json.Unmarshal(data, &_cfg)
	if err != nil {
		log.Error("unmarshal config failed", "err", err)
		panic(err)
	}
	return _cfg, nil
}

func GetConfig() *Config {
	return _cfg
}

var (
	DefaultCors    = []string{"*"} // Default cors domain for the apis
	DefaultVhosts  = []string{"*"} // Default virtual hosts for the apis
	DefaultOrigins = []string{"*"} // Default origins for the apis
	DefaultPrefix  = ""            // Default prefix for the apis
	DefaultModules = []string{}    // enable all module.
	//DefaultModules = []string{"time", "block", "attest"}
)

const (
	APIBatchItemLimit         = 2000
	APIBatchResponseSizeLimit = 250 * 1000 * 1000
)

func GetSafeEpochEndInterval() int64 {
	return 3
}
