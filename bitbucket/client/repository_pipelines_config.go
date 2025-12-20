package client

const (
	RepositoryPipelinesConfigPath = "/repositories/%s/%s/pipelines_config"
)

type RepositoryPipelinesConfig struct {
	RepositoryId string `json:"-"`
	Enabled      bool   `json:"enabled"`
}
