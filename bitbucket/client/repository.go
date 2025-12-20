package client

const (
	RepositoryPath = "/repositories/%s/%s"
)

type Repository struct {
	Uuid        string  `json:"uuid,omitempty"`
	Project     Project `json:"project,omitempty"`
	Name        string  `json:"name,omitempty"`
	IsPrivate   bool    `json:"is_private"`
	UseExisting bool    `json:"-"`
}
