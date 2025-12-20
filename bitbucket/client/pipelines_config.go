package client

const (
	PipelinesConfigPath = "/repositories/%s/%s/pipelines_config"
)

type PipelinesConfig struct {
	RepositoryId string `json:"-"`
	Enabled      bool   `json:"enabled"`
}
