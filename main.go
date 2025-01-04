package main

import (
	"flag"
	"os"
	"strings"
)

// fm -make controller -name sdf -modue m1
// fm -make module -name m2
// fm -rm

type Framework struct {
	Operation string
	Name      string
	Module    string
}

func NewFramework(op, name, module string) *Framework {
	return &Framework{
		Name:      name,
		Operation: op,
		Module:    module,
	}
}

func (obj *Framework) MakeModule() {
	path := "internal"
	os.Chdir(path)
	name_space := ""
	paths := strings.Split(obj.Module, ".")
	for _, p := range paths {
		name_space += p + "_"
		path += "/" + p
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			os.Mkdir(path, 0755)
		}
		os.Chdir(path)
	}

}

func main() {

	ops := map[string]string{
		"controller":  "controllers",
		"dto":         "dtos",
		"enum":        "enums",
		"error":       "errors",
		"eventhandle": "eventhandles",
		"event":       "events",
		"factorie":    "factories",
		"formatter":   "formatters",
		"middleware":  "middlewares",
		"model":       "models",
		"repositorie": "repositories",
		"request":     "requests",
		"service":     "services",
	}

	var op string
	var name string
	var module string

	flag.StringVar(&op, "op", "", "-make")
}

// fm --make   controller --name helloworld --module help
// fm --remove controller --name helloworld
