package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	g "github.com/selabhvl/gotestgen"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"reflect"
)

type Options struct {
	PackageName          string
	ReaderOutputFileName string
	TesterOutputFileName string
	Tests                []string
	TestName             string
	TestType             string
	SystemParameterTypes map[string]string
	InputTypes           map[string]string
	OracleTypes          map[string]string
	InputTypeFields      map[string]map[string]string
	OracleTypeFields     map[string]map[string]string
}


func main() {

	options, dir := parseFlags()
	readerGen(options, dir)

	testerGen(options)
}

func parseFlags() (*Options, *string) {
	dir := flag.String("dir", ".", "path to system test file")
	tests := flag.String("tests", "", "comma separated test names")

	flag.Parse()
	// get test names
	testNames := strings.Split(*tests, ",")

	// get current working dir
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// get package name of working dir
	_, packageName := filepath.Split(workDir)

	return &Options{
		PackageName: packageName,
		Tests:       testNames,
	}, dir
}

func readerGen(options *Options, dir *string) {

	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".xml" {
			fmt.Printf("process XML file: %q\n", path)
			// TODO process XML file

			var global Global

			g.ParseXMLTestCase(path, &global)

			options = getGlobalXMLInformation(options,&global)

			tmpl := createTemplate(options.ReaderOutputFileName, &readerTemplateConcurrent /*&readerTemplate*/, template.FuncMap{"stringFilter": stringFilter})
			generator(tmpl, options, options.ReaderOutputFileName)

		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
	}
}

func testerGen(options *Options) {

	options.TesterOutputFileName = "main_test.go"
	tmpl := createTemplate(options.TesterOutputFileName, &mainTestTemplateConcurrent, template.FuncMap{"stringTitle": strings.Title})
	generator(tmpl, options, options.TesterOutputFileName)
}

func generator(tmpl *template.Template, options *Options, outputFileName string) {

	createDirectoryIfNotExist(outputFileName)
	f, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(f)
	defer func() {
		writer.Flush()
		f.Close()
	}()

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, *options)
	if err != nil {
		panic(err)
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	writer.Write(p)
}

// createTemplate create a template
func createTemplate(testName string, tmplName *string, funcmap map[string]interface{}) *template.Template {

	tmpl, err := template.New(testName).Funcs(funcmap).Parse(*tmplName)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func createDirectoryIfNotExist(file string) {
	directory := filepath.Dir(file)
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		fmt.Println("directory is not exist", directory)
		os.MkdirAll(directory, 0777)
	}
}

// getGlobalXMLInformation gets global information from xml files
func getGlobalXMLInformation(options *Options, global *Global) *Options {

	options.ReaderOutputFileName = "reader" + global.TestName + ".go"
	options.TestName = global.TestName
	options.TestType = global.TestType

	//configurationSetup := global.ReaderSetup.ConfigurationSetup
	//inputArguments := global.ReaderSetup.InputArguments
	//outputResults := global.ReaderSetup.OutputResults

	// get system parameter types for xml to options
	if global.TypeAssignments.SystemParameterTypes.Fields != nil {
		options.SystemParameterTypes = make(map[string]string)

		for _, field := range *global.TypeAssignments.SystemParameterTypes.Fields {
			options.SystemParameterTypes[field.Name] = field.Type
		}
	}

	// get input types from xml to options
	options.InputTypes = make(map[string]string)
	for _, input := range *global.TypeAssignments.InputTypes {
		options.InputTypes[input.Name] = input.Type

		if input.Fields != nil {
			field := make(map[string]string)
			options.InputTypeFields = make(map[string]map[string]string)

			for _, f := range *input.Fields {
				field[f.Name] = f.Type
			}

			options.InputTypeFields[input.Name] = field
		}
	}

	// get oracle types from xml to options
	options.OracleTypes = make(map[string]string)
	for _, oracle := range *global.TypeAssignments.OracleTypes {
		options.OracleTypes[oracle.Name] = oracle.Type

		if oracle.Fields != nil {
			field := make(map[string]string)
			options.OracleTypeFields = make(map[string]map[string]string)

			for _, f := range *oracle.Fields {
				field[f.Name] = f.Type
			}

			options.OracleTypeFields[oracle.Name] = field
		}
	}

	// TODO: the generic function is a little complex to understand or use to replace the code above?...
	/*
		options.InputTypes = make(map[string]string)
		options.OracleTypes = make(map[string]string)
		options.InputTypeFields = make(map[string]map[string]string)
		options.OracleTypeFields = make(map[string]map[string]string)
		options.InputTypes, options.InputTypeFields = getInputOracleTypesInfo(options.InputTypes, options.InputTypeFields, global.TypeAssignments.InputTypes)
		options.OracleTypes, options.OracleTypeFields = getInputOracleTypesInfo(options.OracleTypes, options.OracleTypeFields, global.TypeAssignments.OracleTypes)
	*/
	return options
}

func stringFilter(str string) string {

	if str[len(str)-1:] == "s" {
		return str[:len(str)-1]
	}

	return str
}


// getInputOutputInfo is a generic function to get input/output information
func getInputOracleTypesInfo(inputOracles map[string]string,
	inputOracleFields map[string]map[string]string, typesInfo interface{}) (map[string]string, map[string]map[string]string) {

	value := reflect.ValueOf(typesInfo).Elem()

	for i := 0; i < value.Len(); i++ {
		inputOracles[value.Index(i).FieldByName("Name").String()] = value.Index(i).FieldByName("Type").String()

		if !value.Index(i).FieldByName("Fields").IsNil() {
			field := make(map[string]string)

			testField := value.Index(i).FieldByName("Fields").Elem()

			for j := 0; j < testField.Len(); j++ {
				field[testField.Index(j).FieldByName("Name").String()] = testField.Index(j).FieldByName("Type").String()
			}
			inputOracleFields[value.Index(i).FieldByName("Name").String()] = field
		} else {
			fmt.Println("TestOutputFields is nil")
		}
	}
	fmt.Println("inputOutputs:", inputOracles)
	fmt.Println("inputOutputFields:", inputOracleFields)
	return inputOracles, inputOracleFields
}