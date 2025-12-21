package client

import "strings"

const (
	EnvironmentPath       = "/repositories/%s/%s/environments"
	EnvironmentPathGet    = EnvironmentPath + "/%s"
	EnvironmentPathUpdate = EnvironmentPathGet + "/changes"
)

type Environment struct {
	RepositoryId string                  `json:"-"`
	Uuid         string                  `json:"uuid,omitempty"`
	Name         string                  `json:"name,omitempty"`
	Type         EnvironmentType         `json:"environment_type,omitempty"`
	Restrictions EnvironmentRestrictions `json:"restrictions,omitempty"`
	UseExisting  bool                    `json:"-"`
}
type EnvironmentType struct {
	Name string `json:"name,omitempty"`
}
type EnvironmentRestrictions struct {
	AdminOnly bool `json:"admin_only"`
}
type EnvironmentChanges struct {
	Change EnvironmentChange `json:"change,omitempty"`
}
type EnvironmentChange struct {
	Name         string                  `json:"name,omitempty"`
	Restrictions EnvironmentRestrictions `json:"restrictions,omitempty"`
}
type EnvironmentCollection struct {
	Values []Environment `json:"values"`
}

func (e *Environment) EnvironmentEncodeId() string {
	return e.RepositoryId + IdSeparator + e.Uuid
}

func EnvironmentDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
