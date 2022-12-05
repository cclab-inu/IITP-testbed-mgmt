package types

type ConfigCluster struct {
	Master  string `json:"master,omitempty" bson:"master,omitempty"`
	Worker1 string `json:"worker1,omitempty" bson:"worker1,omitempty"`
	Worker2 string `json:"worker2,omitempty" bson:"worker2,omitempty"`
}

type ConfigKubeArmorRelay struct {
	KubeArmorURL  string `json:"kubearmor_url,omitempty" bson:"kubearmor_url,omitempty"`
	KubeArmorPort string `json:"kubearmor_port,omitempty" bson:"kubearmor_port,omitempty"`
}

type ConfigCiliumHubble struct {
	HubbleURL  string `json:"hubble_url,omitempty" bson:"hubble_url,omitempty"`
	HubblePort string `json:"hubble_port,omitempty" bson:"hubble_port,omitempty"`
}

type Configuration struct {
	Home string `json:"home,omitempty" bson:"home,omitempty"`
	User string `json:"user,omitempty" bson:"user,omitempty"`

	ConfigCluster        ConfigCluster        `json:"config_cluster,omitempty" bson:"config_cluster,omitempty"`
	ConfigCiliumHubble   ConfigCiliumHubble   `json:"config_cilium_hubble,omitempty" bson:"config_cilium_hubble,omitempty"`
	ConfigKubeArmorRelay ConfigKubeArmorRelay `json:"config_kubearmor_relay,omitempty" bson:"config_kubearmor_relay,omitempty"`
}
