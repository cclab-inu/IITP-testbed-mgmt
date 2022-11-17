package types

// Pod Structure
type Pod struct {
	Namespace string `json:"namespace" bson:"namespace"`
	PodName   string `json:"pod_name" bson:"pod_name"`
	PodIP     string `json:"pod_ip" bson:"pod_ip"`
}
