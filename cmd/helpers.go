package cmd

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/epiphany-platform/e-structures/utils/load"

	azks "github.com/epiphany-platform/e-structures/azks/v0"
	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/to"
)

func ensureDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func backupFile(path string) error {
	logger.Debug().Msgf("backupFile(%s)", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	} else {
		backupPath := path + ".backup"

		input, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(backupPath, input, 0644)
		if err != nil {
			return err
		}
		return nil
	}
}

func checkAndLoad(stateFilePath string, configFilePath string) (*azks.Config, *st.State, error) {
	logger.Debug().Msgf("checkAndLoad(%s, %s)", stateFilePath, configFilePath)
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("state file does not exist, please run init first")
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("config file does not exist, please run init first")
	}

	state, err := load.State(stateFilePath)
	if err != nil {
		return nil, nil, err
	}

	config, err := load.AzKSConfig(configFilePath)
	if err != nil {
		return nil, nil, err
	}

	return config, state, nil
}

func produceOutput(m map[string]interface{}) *azks.Output {
	logger.Debug().Msgf("Received output map: %#v", m)

	return &azks.Output{
		KubeConfig: to.StrPtr(m["kubeconfig"].(string)),
	}
}
