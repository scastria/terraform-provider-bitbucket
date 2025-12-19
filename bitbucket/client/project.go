package client

const (
	ProjectPath = "/workspaces/%s/projects/%s"
)

type Project struct {
	Uuid string `json:"uuid,omitempty"`
	Key  string `json:"key,omitempty"`
}
