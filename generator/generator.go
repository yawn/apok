package generator

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"time"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	"github.com/yawn/apok"
)

const (
	actionsPkg  = "actions"
	mainPkg     = "github.com/yawn/apok"
	servicesPkg = "services"
)

var now = time.Now()

type Generator struct {
	directory string
}

func New(directory string) (*Generator, error) {

	for _, pkg := range []string{
		actionsPkg,
		servicesPkg,
	} {

		dir := path.Join(directory, pkg)

		if err := os.MkdirAll(dir, 0777); err != nil {
			return nil, errors.Wrapf(err, "failed to create directory %q", dir)
		}

	}

	return &Generator{
		directory: directory,
	}, nil

}

func (g *Generator) Actions(service apok.Service, actions []apok.Action) error {

	var (
		all  jen.Dict = make(map[jen.Code]jen.Code, len(actions))
		file          = header(actionsPkg)
	)

	for _, action := range actions {

		all[jen.Lit(action.Name)] = jen.Op("&").Id(action.Var())

		_struct(file.Var().Id(fmt.Sprintf("%s_%s", service.Var(), action.Var())).Op("="), "Action", action)

		file.Line()

	}

	// add all actions under the service var name
	file.Var().Id(service.Var()).Op("=").Map(jen.String()).Op("*").Id("Action").Values(all)

	name := path.Join(g.directory, actionsPkg, fmt.Sprintf("%s.go", service.Filename()))

	fh, err := os.Create(name)

	if err != nil {
		return errors.Wrapf(err, "failed to create file %q", name)
	}

	if err := file.Render(fh); err != nil {
		return errors.Wrapf(err, "failed to render code in file %q", name)
	}

	fh.Close()

	return nil

}

func (g *Generator) Services(services []apok.Service) error {

	var all jen.Dict = make(map[jen.Code]jen.Code, len(services))

	for _, service := range services {

		all[jen.Lit(service.ActionPrefix)] = jen.Op("&").Id(service.Var())

		file := header(servicesPkg)

		_struct(file.Var().Id(service.Var()).Op("="), "Service", service)

		name := path.Join(g.directory, servicesPkg, fmt.Sprintf("%s.go", service.Filename()))

		fh, err := os.Create(name)

		if err != nil {
			return errors.Wrapf(err, "failed to create file %q", name)
		}

		if err := file.Render(fh); err != nil {
			return errors.Wrapf(err, "failed to render code in file %q", name)
		}

		fh.Close()

	}

	name := path.Join(g.directory, servicesPkg, "all.go")

	fh, err := os.Create(name)

	file := header(servicesPkg)

	file.Var().Id("All").Op("=").Map(jen.String()).Op("*").Qual(mainPkg, "Service").Values(all)

	if err != nil {
		return errors.Wrapf(err, "failed to create all-service file")
	}

	if err := file.Render(fh); err != nil {
		return errors.Wrapf(err, "failed to render code in all-service")
	}

	fh.Close()

	return nil

}

func header(pkg string) *jen.File {

	// TODO: embded the JSON at the end of the file

	file := jen.NewFile(pkg)
	file.HeaderComment(fmt.Sprintf("This file was generated on %s - do not edit.", now.Format(time.RFC3339)))
	file.Line()

	return file

}

// plan here is to become pretty dnamic when it comes to struct creation
func _struct(ptr *jen.Statement, name string, s interface{}) *jen.Statement {

	if ptr == nil {
		ptr = jen.Empty()
	}

	var (
		fields = make(map[jen.Code]jen.Code)
		typ    = reflect.TypeOf(s)
	)

	for i := 0; i < typ.NumField(); i++ {

		var (
			key = typ.Field(i).Name
			val = reflect.ValueOf(s).Field(i).Interface()
		)

		switch t := val.(type) {

		// TODO: don't use Values(), they don't look nice

		case []string:

			var blocks []jen.Code

			for _, e := range t {
				blocks = append(blocks, jen.Lit(e))
			}

			if len(blocks) > 0 {
				fields[jen.Id(key)] = jen.Index().String().Values(blocks...)
			}

		case []*apok.Action:

		case []apok.ContextKey:

			var blocks []jen.Code

			for _, e := range t {
				blocks = append(blocks, _struct(nil, "ContextKey", e))
			}

			if len(blocks) > 0 {
				fields[jen.Id(key)] = jen.Index().Qual(mainPkg, "ContextKey").Values(blocks...)
			}

		case []apok.Resource:

			var blocks []jen.Code

			for _, e := range t {
				blocks = append(blocks, _struct(nil, "Resource", e))
			}

			if len(blocks) > 0 {
				fields[jen.Id(key)] = jen.Index().Qual(mainPkg, "Resource").Values(blocks...)
			}

		default:
			fields[jen.Id(key)] = jen.Lit(val)
		}

	}

	ptr.Qual(mainPkg, name).Values(jen.Dict(fields))

	return ptr

}
