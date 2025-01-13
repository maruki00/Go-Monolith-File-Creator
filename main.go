package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
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
	"factory":     "factories",
	"formatter":   "formatters",
	"middleware":  "middlewares",
	"model":       "models",
	"repository":  "repositories",
	"request":     "requests",
	"service":     "services",
	"contract":    "contracts",
}

var projectFolders = []string{
	"internal",
	"configs",
	"cmd",
	"migrations",
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

func (obj *Framework) GetPath(module string) string {
	root, _ := os.Getwd()
	path := root + "/internal/"
	os.MkdirAll(path, 0754)
	os.Chdir(path)
	paths := strings.Split(module, ".")
	path += strings.Join(paths, "/")
	fmt.Println("this is the path for : ", module, path)
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
	path := obj.GetPath(obj.Module)
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

// here
func (obj *Framework) MakeOperation(implementIface bool) {

	op, ok := folders[strings.ToLower(obj.Operation)]
	if !ok {
		panic("Operation Not Supported")
	}
	path := obj.GetPath(obj.Module)
	name := strings.ToLower(obj.Name)

	filePath := path + "/" + op + "/" + name + obj.firstUppercase(obj.Operation) + ".go"
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		panic("file already exists")
	}
	pkg := obj.GetPackage(op)
	class := obj.firstUppercase(name) + obj.firstUppercase(obj.Operation)

	file, _ := os.Create(filePath)
	defer file.Close()
	_, _ = file.WriteString(fmt.Sprintf("package %s\n\n\ntype %s struct {\n\n\t//implementation goes here\n\n}\n\n", pkg, class))
	if !implementIface {
		return
	}
	impls := obj.ImplemementInterface(class, "")

	fmt.Println("implementation : ", impls)

}

func (obj *Framework) firstUppercase(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}

func usage() {
	H := ` 
	./fm -make [option] -name [name] -module [module]
	options:
		-controller
		-dto
		-enum
		-error
		-eventhandler
		-event
		-factory
		-middleware
		-model
		-repository
		-request
		-service
		-contract
		-init
	examples:
		- ./fm -make module -name v1.module1
		- ./fm -make controller -name controller1 -module v1.module
		- ./fm -make init -name controller1 -module v1.module
	`
	fmt.Println(H)
}

func (obj *Framework) InitProject() {

	for _, folder := range projectFolders {
		//os.Mkdir(path.Join(pathProject, folder), 0754)
		os.MkdirAll(folder, 0754)
	}
}

type Hello interface {
	Sleep() int
}

func (obj *Framework) ImplemementInterface(pathIface string, class string) string {
	path := obj.GetPath(class)
	fmt.Println(path)
	return path
}

func GetInterfaceBody(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic("file doesnt eexists")
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic("error reading the file")
	}

	functions := make([]string, 0)
	regex := `(\w+)\s*\(([^)]*)\)\s*(\w+)?\s*(error)?\s*$`
	r := regexp.MustCompile(regex)
	matches := r.FindAllStringSubmatch(string(data), -1)
	for _, match := range matches {
		functions = append(functions, strings.Join(match, " "))
	}
	return functions
}

func main() {

	var op string
	var name string
	var module string
	var use *string
	var iface *string
	var initP *string

	use = flag.String("show", "", "--show usage")
	initP = flag.String("init", "", "--init usage")
	flag.StringVar(&op, "make", "", "-make controller")
	flag.StringVar(&name, "name", "", "-name helloworld")
	flag.StringVar(&module, "module", "", "-module v1.module1")
	iface = flag.String("interface", "", "-interface v1.module1.interfaces.iface1")
	flag.Parse()

	if *use == "usage" {
		usage()
		return
	}
	fm := NewFramework(op, name, module)
	if *initP != "" {
		fm.InitProject()
		return
	}

	if *use != "usage" && (name == "" || op == "") {
		usage()
		return
	}

	if op == "module" {
		fm.MakeModule()
	} else {
		if *iface != "" {
			fm.MakeOperation(true)
		} else {
			fm.MakeOperation(true)
		}

	}
}

// fm --make   controller --name helloworld --module help
// fm --remove controller --name helloworld
