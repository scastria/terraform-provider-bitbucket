package client

import (
	"strconv"
	"strings"
)

const (
	RestrictionPath    = "/repositories/%s/%s/branch-restrictions"
	RestrictionPathGet = RestrictionPath + "/%s"
)

type Restriction struct {
	RepositoryId    string `json:"-"`
	Id              int    `json:"id,omitempty"`
	Kind            string `json:"kind,omitempty"`
	BranchMatchKind string `json:"branch_match_kind,omitempty"`
	BranchType      string `json:"branch_type,omitempty"`
	Pattern         string `json:"pattern,omitempty"`
	Value           int    `json:"value,omitempty"`
	UseExisting     bool   `json:"-"`
}
type RestrictionCollection struct {
	Values []Restriction `json:"values"`
}

func (r *Restriction) RestrictionEncodeId() string {
	return r.RepositoryId + IdSeparator + strconv.Itoa(r.Id)
}

func RestrictionDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
