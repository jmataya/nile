package nile

import "fmt"

type routeParam struct {
	Name  string
	Value interface{}
}

type routeMatch struct {
	PathTemplate string
	URI          string
	Params       []*routeParam
}

func (rm *routeMatch) AddPathTemplate(template string) {
	rm.PathTemplate = fmt.Sprintf("%s/%s", rm.PathTemplate, template)
}

func (rm *routeMatch) AddParam(name string, value interface{}) {
	rm.Params = append(rm.Params, &routeParam{Name: name, Value: value})
}
