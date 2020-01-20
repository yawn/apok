package apok

import (
	"strings"

	"github.com/iancoleman/strcase"
)

type Service struct {
	ActionPrefix string `json:"actionPrefix"`
	ArnFormat    string `json:"arnFormat"`
	HasResource  bool   `json:"hasResource"`
	Name         string `json:"name"`
}

func (s Service) Filename() string {

	name := strings.ReplaceAll(s.ActionPrefix, "-", "_")

	return name

}

func (s Service) Var() string {

	name := strings.ReplaceAll(s.ActionPrefix, "-", " ")
	name = strings.ReplaceAll(s.ActionPrefix, "_", " ")
	name = strcase.ToCamel(name)

	return name

}
