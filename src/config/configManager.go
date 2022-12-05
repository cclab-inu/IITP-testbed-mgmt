package config

import (
	types "github.com/cclab.inu/testbed-mgmt/src/types"
	"github.com/spf13/viper"
)

// ====================== //
// == Global Variables == //
// ====================== //

var CurrentCfg types.Configuration

// =========================== //
// == Configuration Loading == //
// =========================== //

func LoadConfigCluster() types.ConfigCluster {
	cfgCluster := types.ConfigCluster{}

	cfgCluster.Master = viper.GetString("cluster.master")
	cfgCluster.Worker1 = viper.GetString("cluster.worker1")
	cfgCluster.Worker2 = viper.GetString("cluster.worker2")

	return cfgCluster
}

func LoadConfigCiliumHubble() types.ConfigCiliumHubble {
	cfgHubble := types.ConfigCiliumHubble{}

	cfgHubble.HubbleURL = viper.GetString("cilium-hubble.url")
	cfgHubble.HubblePort = viper.GetString("cilium-hubble.port")

	return cfgHubble
}

func LoadConfigKubeArmor() types.ConfigKubeArmorRelay {
	cfgKubeArmor := types.ConfigKubeArmorRelay{}

	cfgKubeArmor.KubeArmorRelayURL = viper.GetString("kubearmor.url")
	cfgKubeArmor.KubeArmorRelayPort = viper.GetString("kubearmor.port")

	return cfgKubeArmor
}

func LoadConfigFromFile() {
	CurrentCfg = types.Configuration{}

	// default
	CurrentCfg.Home = viper.GetString("project.home")
	CurrentCfg.User = viper.GetString("project.user")

	// load cluster related config
	CurrentCfg.ConfigCluster = LoadConfigCluster()

	// load cilium hubble relay
	CurrentCfg.ConfigCiliumHubble = LoadConfigCiliumHubble()

	// load kubearmor relay config
	CurrentCfg.ConfigKubeArmorRelay = LoadConfigKubeArmor()
}

// ============================ //
// == Get Configuration Info == //
// ============================ //

func GetCurrentCfg() types.Configuration {
	return CurrentCfg
}

func GetCfgCluster() types.ConfigCluster {
	return CurrentCfg.ConfigCluster
}

func GetCfgCiliumHubble() types.ConfigCiliumHubble {
	return CurrentCfg.ConfigCiliumHubble
}

func GetCfgKubeArmor() types.ConfigKubeArmorRelay {
	return CurrentCfg.ConfigKubeArmorRelay
}
