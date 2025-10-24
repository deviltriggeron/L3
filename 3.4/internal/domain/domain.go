package domain

type ServerConfig struct {
	Host string
	Port string
}

type MinioCfg struct {
	Endpoint        string
	Bucket          string
	BucketProcessed string
	AccessKey       string
	SecretKey       string
	SSL             bool
}

type ConfigBroker struct {
	Broker  string
	GroupID string
	Topic   string
}
