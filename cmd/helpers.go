package cmd

import (
	"errors"
	azks "github.com/epiphany-platform/e-structures/azks/v0"
	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/to"
	"io/ioutil"
	"os"
)

func ensureDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func loadState(path string) (*st.State, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return st.NewState(), nil
	} else {
		state := &st.State{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = state.Unmarshal(bytes)
		if err != nil {
			return nil, err
		}
		if state.AzKS == nil {
			state.AzKS = &st.AzKSState{}
		}
		return state, nil
	}
}

func saveState(path string, state *st.State) error {
	bytes, err := state.Marshal()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadConfig(path string) (*azks.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return azks.NewConfig(), nil
	} else {
		config := &azks.Config{}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = config.Unmarshal(bytes)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
}

func saveConfig(path string, config *azks.Config) error {
	bytes, err := config.Marshal()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func backupFile(path string) error {
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
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("state file does not exist, please run init first")
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("config file does not exist, please run init first")
	}

	state, err := loadState(stateFilePath)
	if err != nil {
		return nil, nil, err
	}

	config, err := loadConfig(configFilePath)
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
