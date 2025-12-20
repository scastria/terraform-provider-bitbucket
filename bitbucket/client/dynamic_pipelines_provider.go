package client

const (
	DynamicPipelinesProviderPath = "/repositories/%s/%s/pipelines-config/dynamic-pipelines-provider"
)

type DynamicPipelinesProvider struct {
	RepositoryId string `json:"-"`
	AppAri       string `json:"appAri"`
}
