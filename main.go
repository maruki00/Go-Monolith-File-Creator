package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// fm -make controller -name sdf -modue m1
// fm -make module -name m2
// fm -rm
var folders = map[string]string{
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

func (obj *Framework) GetPath() string {
	root, _ := os.Getwd()
	path := root + "/internal/"
	os.Chdir(path)
	paths := strings.Split(obj.Module, ".")
	path += strings.Join(paths, "/")
	return path
}

func (obj *Framework) MakeModuleDir() {
	path, _ := os.Getwd()
	os.Chdir(path)
	paths := strings.Split(obj.Module, ".")
	for _, p := range paths {
		path += "/" + p
		fmt.Println(path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			_ = os.Mkdir(path, 0755)

		}
		os.Chdir(path)
	}
}

func (obj *Framework) MakeModule() {
	path := obj.GetPath()
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		panic("module already exists")
	}
	obj.MakeModuleDir()
	for _, folder := range folders {
		_ = os.Mkdir(fmt.Sprintf("%s/%s", path, folder), 0755)
		_, _ = os.Create(fmt.Sprintf("%s/%s/.gitkeep", path, folder))
	}
}

func (obj *Framework) GetPackage(op string) string {
	paths := strings.Split(obj.Module, ".")
	paths = append(paths, op)
	return strings.Join(paths, "_")
}

func (obj *Framework) MakeOperation() {

	op, ok := folders[strings.ToLower(obj.Operation)]
	if !ok {
		panic("Operation Not Supported")
	}
	path := obj.GetPath()
	name := strings.ToLower(obj.Name)

	filePath := path + "/" + op + "/" + name + obj.firstUppercase(obj.Operation) + ".go"
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		panic("file already exists")
	}
	pkg := obj.GetPackage(op)
	class := obj.firstUppercase(name) + obj.firstUppercase(obj.Operation)

	file, _ := os.Create(filePath)
	defer file.Close()
	_, _ = file.WriteString(fmt.Sprintf("package %s\n\n\ntype %s struct {\n\n\t//implementation goes here\n\n}", pkg, class))

}

func (obj *Framework) firstUppercase(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}

func main() {

	var op string
	var name string
	var module string

	flag.StringVar(&op, "make", "", "-make controller")
	flag.StringVar(&name, "name", "", "-name helloworld")
	flag.StringVar(&module, "module", "", "-module v1.module1")
	flag.Parse()

	fm := NewFramework(op, name, module)
	if op == "module" {
		fm.MakeModule()
	} else {
		fm.MakeOperation()
	}
}

// fm --make   controller --name helloworld --module help
// fm --remove controller --name helloworld
