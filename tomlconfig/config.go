package tomlconfig

import (
	"github.com/BurntSushi/toml"
)

func ParseTomlConfig(filepath string, v interface{}) error {
	_, err := toml.DecodeFile(filepath, v)
	return err
}
