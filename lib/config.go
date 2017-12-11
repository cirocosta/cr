package lib

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func ConfigFromFile(file string) (config Config, err error) {
	finfo, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			err = errors.Wrapf(err,
				"configuration file %s not found",
				file)
			return
		}

		err = errors.Wrapf(err,
			"unexpected error looking for config file %s",
			file)
		return
	}

	configContent, err := ioutil.ReadAll(finfo)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't properly read config file %s",
			file)
		return
	}

	err = yaml.Unmarshal(configContent, &config)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't properly parse yaml config file %s",
			file)
		return
	}

	return
}
