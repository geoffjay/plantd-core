package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config interface{}

// LoadConfig reads in a configuration file from a set of locations and
// deserializes it into a Config instance.
func LoadConfig(name string, c interface{}) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	envPrefix := fmt.Sprintf("PLANTD_%s", strings.ToUpper(name))
	envConfig := fmt.Sprintf("%s_CONFIG", envPrefix)

	config := viper.New()

	file := os.Getenv(envConfig)
	if file == "" {
		config.SetConfigName(name)
		config.AddConfigPath(".")
		config.AddConfigPath(fmt.Sprintf("%s/.config/plantd", home))
		config.AddConfigPath("/etc/plantd")
	} else {
		var extension string
		base := filepath.Base(file)
		if strings.HasSuffix(base, "yaml") ||
			strings.HasSuffix(base, "json") ||
			strings.HasSuffix(base, "toml") {
			// strip the file type for viper
			parts := strings.Split(filepath.Base(file), ".")
			base = strings.Join(parts[:len(parts)-1], ".")
			extension = parts[len(parts)-1]
		} else {
			return errors.New("configuration does not support that extension type")
		}
		config.SetConfigName(base)
		config.SetConfigType(extension)
		config.SetConfigFile(file)
		config.AddConfigPath(filepath.Dir(file))
	}

	err = config.ReadInConfig()
	if err != nil {
		return err
	}

	config.SetEnvPrefix(envPrefix)
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()

	err = config.Unmarshal(&c)
	if err != nil {
		return err
	}

	return nil
}
