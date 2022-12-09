package config

import (
	types "github.com/cclab.inu/testbed-mgmt/src/types"
	"github.com/spf13/viper"
)

// ====================== //
// == Global Variables == //
// ====================== //

var CurrentCfg types.Configuration

func init() {
	viper.AddConfigPath("./conf")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	LoadConfigFromFile()
}

// ======================== //
// == Load Configuration == //
// ======================== //

func LoadConfigCluster() types.ConfigCluster {
	cfgCluster := types.ConfigCluster{}

	cfgCluster.Master = viper.GetString("cluster.master")
	cfgCluster.Worker1 = types.ConfigWorker{
		IP:    viper.GetString("cluster.worker1.ip"),
		SSHID: viper.GetString("cluster.worker1.id"),
		SSHPW: viper.GetString("cluster.worker1.pw"),
	}
	cfgCluster.Worker2 = types.ConfigWorker{
		IP:    viper.GetString("cluster.worker2.ip"),
		SSHID: viper.GetString("cluster.worker2.id"),
		SSHPW: viper.GetString("cluster.worker2.pw"),
	}

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

	cfgKubeArmor.KubeArmorURL = viper.GetString("kubearmor.url")
	cfgKubeArmor.KubeArmorPort = viper.GetString("kubearmor.port")

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
