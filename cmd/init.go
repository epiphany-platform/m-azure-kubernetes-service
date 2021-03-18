package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"

	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	omitState bool

	name       string
	vmsRsaPath string

	rgName     string
	snName     string
	vnetName   string
	k8sVersion string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initializes module configuration file",
	Long:  `Initializes module configuration file (in ` + filepath.Join(defaultSharedDirectory, moduleShortName, configFileName) + `/). `,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("PreRun")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err).Msg("BindPFlags failed")
		}

		name = viper.GetString("name")
		vmsRsaPath = viper.GetString("vms_rsa")
		rgName = viper.GetString("rg_name")
		snName = viper.GetString("subnet_name")
		vnetName = viper.GetString("vnet_name")
		k8sVersion = viper.GetString("kubernetes_version")
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("init called")
		moduleDirectoryPath := filepath.Join(SharedDirectory, moduleShortName)
		configFilePath := filepath.Join(SharedDirectory, moduleShortName, configFileName)
		stateFilePath := filepath.Join(SharedDirectory, stateFileName)
		logger.Debug().Msg("ensure directories")
		err := ensureDirectory(moduleDirectoryPath)
		if err != nil {
			logger.Fatal().Err(err).Msg("ensureDirectory failed")
		}
		logger.Debug().Msg("load state file")
		state, err := loadState(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("loadState failed")
		}
		logger.Debug().Msg("load config file")
		config, err := loadConfig(configFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("loadConfig failed")
		}

		if !reflect.DeepEqual(state.AzKS, &st.AzKSState{}) && state.AzKS.Status != st.Initialized && state.AzKS.Status != st.Destroyed {
			logger.Fatal().Err(errors.New(string("unexpected state: " + state.AzKS.Status))).Msg("incorrect state")
		}

		logger.Debug().Msg("backup state file")
		err = backupFile(stateFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("backupFile failed")
		}
		logger.Debug().Msg("backup config file")
		err = backupFile(configFilePath)
		if err != nil {
			logger.Fatal().Err(err).Msg("backupFile failed")
		}

		config.GetParams().Name = to.StrPtr(name)
		config.GetParams().RsaPublicKeyPath = to.StrPtr(filepath.Join(SharedDirectory, fmt.Sprintf("%s.pub", vmsRsaPath)))
		config.GetParams().RgName = to.StrPtr(rgName)
		config.GetParams().SubnetName = to.StrPtr(snName)
		config.GetParams().VnetName = to.StrPtr(vnetName)
		config.GetParams().KubernetesVersion = to.StrPtr(k8sVersion)

		//initialize configuration using values from AzBIState
		if !omitState {
			if state.GetAzBIState().Status == st.Applied {
				if state.GetAzBIState().GetConfig().GetParams().GetNameV() != "" {
					config.GetParams().Name = to.StrPtr(state.GetAzBIState().GetConfig().GetParams().GetNameV())
					fmt.Println("Found and used 'name' parameter in existing AzBI configuration.")
				}
				if state.GetAzBIState().GetConfig().GetParams().GetRsaPublicKeyV() != "" {
					config.GetParams().RsaPublicKeyPath = to.StrPtr(state.GetAzBIState().GetConfig().GetParams().GetRsaPublicKeyV())
					fmt.Println("Found and used 'vms_rsa' parameter in existing AzBI configuration.")
				}
				if state.GetAzBIState().GetConfig().GetParams().GetLocationV() != "" {
					config.GetParams().Location = to.StrPtr(state.GetAzBIState().GetConfig().GetParams().GetLocationV())
					fmt.Println("Found and used 'location' parameter in existing AzBI configuration.")
				}
				if state.GetAzBIState().GetOutput().GetRgNameV() != "" {
					config.GetParams().RgName = to.StrPtr(state.GetAzBIState().GetOutput().GetRgNameV())
					fmt.Println("Found and used 'rg_name' parameter in existing AzBI output.")
				}
				if state.GetAzBIState().GetOutput().GetVnetNameV() != "" {
					config.GetParams().VnetName = to.StrPtr(state.GetAzBIState().GetOutput().GetVnetNameV())
					fmt.Println("Found and used 'vnet_name' parameter in existing AzBI output.")
				}
				for _, s := range state.GetAzBIState().GetConfig().GetParams().ExtractEmptySubnets() {
					if *s.Name == "azks" || *s.Name == "kubernetes" || *s.Name == "aks" {
						config.GetParams().SubnetName = s.Name
						fmt.Println("Found and used 'subnet.name' parameter in existing AzBI configuration.")
					}
				}
			}
		}

		state.AzKS.Status = st.Initialized

		logger.Debug().Msg("save config")
		err = saveConfig(configFilePath, config)
		if err != nil {
			logger.Fatal().Err(err).Msg("saveConfig failed")
		}
		logger.Debug().Msg("save state")
		err = saveState(stateFilePath, state)
		if err != nil {
			logger.Fatal().Err(err).Msg("saveState failed")
		}

		bytes, err := config.Marshal()
		if err != nil {
			logger.Fatal().Err(err).Msg("config.Marshal failed")
		}
		logger.Debug().Msg(string(bytes))
		fmt.Println("Initialized config: \n" + string(bytes))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&omitState, "omit_state", "o", false, "omit state values during initialization")

	initCmd.Flags().String("name", "epiphany", "prefix given to all resources created") //TODO rename to prefix
	initCmd.Flags().String("vms_rsa", "vms_rsa", "name of rsa keypair to be provided to machines")

	initCmd.Flags().String("rg_name", "epiphany-rg", "name of Azure Resource Group to be used")
	initCmd.Flags().String("subnet_name", "azks", "name of subnet to be used")
	initCmd.Flags().String("vnet_name", "epiphany-vnet", "name of vnet to be used")

	initCmd.Flags().String("kubernetes_version", "1.18.14", "version of kubernetes to be used")
}
