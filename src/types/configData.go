package types

type ConfigCiliumHubble struct {
	HubbleURL  string `json:"hubble_url,omitempty" bson:"hubble_url,omitempty"`
	HubblePort string `json:"hubble_port,omitempty" bson:"hubble_port,omitempty"`
}
