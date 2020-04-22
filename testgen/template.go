package main

var readerTemplateConcurrent = `package {{.PackageName}}

import	"encoding/xml"

// {{.TestName}}Test is an exported type for {{.TestName}} tests
type {{.TestName}}Test struct {
	XMLName	 xml.Name	` + "`" + `xml:"Test"` + "`" + `
	TestName string 	` + "`" + `xml:"Name,attr"` + "`" + `
	TestType string 	` + "`" + `xml:"Type,attr"` + "`" + `
	SystemParameters *SystemParameters	 ` + "`" + `xml:"SystemParameters"` + "`" + `
	TestCases []*{{.TestName}}TestCases	 ` + "`" + `xml:"TestCase"` + "`" + `
}

type {{.TestName}}TestCases struct {
	XMLName	 xml.Name	` + "`" + `xml:"TestCase"` + "`" + `
	CaseID string 	` + "`" + `xml:"ID,attr"` + "`" + `
	TestValues *{{.TestName}}TestValues ` + "`" + `xml:"InputValues"` + "`" + `
	TestOracles *{{.TestName}}TestOracles ` + "`" + `xml:"Oracles"` + "`" + `
}

type {{.TestName}}TestValues struct {
	XMLName	 xml.Name	` + "`" + `xml:"InputValues"` + "`" + `
	ConcurrentCalls []*Calls  ` + "`" + `xml:"Concurrent>Call"` + "`" + `
	SequentialCalls []*Calls ` + "`" + `xml:"Call"` + "`" + `
}

type Calls struct {
	XMLName	 xml.Name	` + "`" + `xml:"Call"` + "`" + `
	CallName string ` + "`" + `xml:"Name,attr"` + "`" + `
	TestValue *TestValue ` + "`" + `xml:"InputValue"` + "`" + `
}

type TestValue struct {
	XMLName	 xml.Name	` + "`" + `xml:"InputValue"` + "`" + `
	{{- range $inputName, $inputType := .InputTypes }}
    {{$inputName}} {{$inputType}} ` + "`" + `xml:"{{$inputName | stringFilter}}"` + "`" + `
    {{- end}}
}

type {{.TestName}}TestOracles struct {
	XMLName	 xml.Name	` + "`" + `xml:"Oracles"` + "`" + `
	{{- range $oracleName, $oracleType := .OracleTypes }}
    {{$oracleName}} {{$oracleType}} ` + "`" + `xml:"{{$oracleName | stringFilter}}"` + "`" + `
    {{- end}}
}

{{range $structName, $field := .InputTypeFields -}}
type {{$structName | stringFilter }} struct{
	{{- range $fieldName, $fieldType := $field }}
	{{$fieldName}} {{$fieldType}} ` + "`" + `xml:"{{$fieldName}}"` + "`" + `
	{{- end}}
}
{{- end}}

{{range $structName, $field := .OracleTypeFields -}}
type {{$structName | stringFilter }} struct{
	{{- range $fieldName, $fieldType := $field }}
	{{$fieldName}} {{$fieldType}} ` + "`" + `xml:"{{$fieldName}}"` + "`" + `
	{{- end}}
}
{{- end}}

{{ if eq .TestType "systemtest" -}}
type SystemParameters struct {
	XMLName	 xml.Name	` + "`" + `xml:"SystemParameters"` + "`" + `
	{{- range $fieldName, $fieldType := .SystemParameterTypes }}
	{{$fieldName}} {{$fieldType}} ` + "`" + `xml:"{{$fieldName | stringFilter}}"` + "`" + `
	{{- end}}
}
{{- end}}`

var readerTemplateDefault = `package {{.PackageName}}

import	"encoding/xml"

// {{.TestName}}Test is an exported type for {{.TestName}} tests
type {{.TestName}}Test struct {
	XMLName	 xml.Name	` + "`" + `xml:"Test"` + "`" + `
	TestName string 	` + "`" + `xml:"Name,attr"` + "`" + `
	TestType string 	` + "`" + `xml:"Type,attr"` + "`" + `
	SystemParameters *SystemParameters	 ` + "`" + `xml:"SystemParameters"` + "`" + `
	TestCases []*{{.TestName}}TestCases	 ` + "`" + `xml:"TestCase"` + "`" + `
}

type {{.TestName}}TestCases struct {
	XMLName	 xml.Name	` + "`" + `xml:"TestCase"` + "`" + `
	CaseID string 	` + "`" + `xml:"ID,attr"` + "`" + `
	TestValues *{{.TestName}}TestValues ` + "`" + `xml:"InputValues"` + "`" + `
	TestOracles *{{.TestName}}TestOracles ` + "`" + `xml:"Oracles"` + "`" + `
}

type {{.TestName}}TestValues struct {
	XMLName	 xml.Name	` + "`" + `xml:"InputValues"` + "`" + `
	{{- range $inputName, $inputType := .InputTypes }}
    {{$inputName}} {{$inputType}} ` + "`" + `xml:"{{$inputName | stringFilter}}"` + "`" + `
    {{- end}}
}

type {{.TestName}}TestOracles struct {
	XMLName	 xml.Name	` + "`" + `xml:"Oracles"` + "`" + `
	{{- range $oracleName, $oracleType := .OracleTypes }}
    {{$oracleName}} {{$oracleType}} ` + "`" + `xml:"{{$oracleName | stringFilter}}"` + "`" + `
    {{- end}}
}

{{range $structName, $field := .InputTypeFields -}}
type {{$structName | stringFilter }} struct{
	{{- range $fieldName, $fieldType := $field }}
	{{$fieldName}} {{$fieldType}} ` + "`" + `xml:"{{$fieldName}}"` + "`" + `
	{{- end}}
}
{{- end}}

{{range $structName, $field := .OracleTypeFields -}}
type {{$structName | stringFilter }} struct{
	{{- range $fieldName, $fieldType := $field }}
	{{$fieldName}} {{$fieldType}} ` + "`" + `xml:"{{$fieldName}}"` + "`" + `
	{{- end}}
}
{{- end}}

{{ if eq .TestType "systemtest" -}}
type SystemParameters struct {
	XMLName	 xml.Name	` + "`" + `xml:"SystemParameters"` + "`" + `
	{{- range $fieldName, $fieldType := .SystemParameterTypes }}
	{{$fieldName}} {{$fieldType}} ` + "`" + `xml:"{{$fieldName | stringFilter}}"` + "`" + `
	{{- end}}
}
{{- end}}`

var mainTestTemplateConcurrent = `package {{.PackageName}}

import (
	"flag"
	"os"
	"testing"
	"reflect"
	"sync"
)

type tester struct {
	t  *testing.T
}

var (
    {{- range $name := .Tests }}
	{{$name}}Test  {{$name | stringTitle}}Test
	{{- end}}
)

func TestMain(m *testing.M) {
	// Flag definitions.
    {{- range $name := .Tests }}
	var {{$name}}TCsDir = flag.String(
		"{{$name}}TCsDir",
		"",
		"path to the file for {{$name | stringTitle}} tests",
	)
	{{- end}}

    // Parse and validate flags.
	flag.Parse()

	// Load test cases from XML files:
	{{- range $name := .Tests }}
	ParseXMLTestCase(*{{$name}}TCsDir, &{{$name}}Test)
	{{- end}}
    
	// Run tests/benchmarks.
	res := m.Run()
	os.Exit(res)
}

{{range $name := .Tests}}  
func Test{{$name | stringTitle}}(t *testing.T){

{{if eq $name "system"}}
	tester := &tester{t}
	t.Logf("test name=%v", {{$name}}Test.TestName)
    
    for _, cs := range {{$name}}Test.TestCases{
		t.Logf("test case ID=%v", cs.CaseID)
		// concurrent executions
		if cs.TestValues.ConcurrentCalls == nil {
			t.Log("no concurrent executions")
		}else {
			var wg sync.WaitGroup
			wg.Add(len(cs.TestValues.ConcurrentCalls))
			for _, concurrentCall :=  range cs.TestValues.ConcurrentCalls {
				t.Logf("concurrent call name=%v",concurrentCall.CallName)
				go func(testValue *TestValue) {

					// TODO: invoke funcExecutor with call names and test values to call other functions, based on the SUT
					funcExecutor(t, tester, concurrentCall.CallName, testValue)

					wg.Done()
				}(concurrentCall.TestValue)
			}
			wg.Wait()
		}
		// sequential executions
		if cs.TestValues.SequentialCalls == nil {
			t.Log("no sequential executions")
		}else{
			for _, sequentailCall := range cs.TestValues.SequentialCalls {
				t.Logf("sequentail call name=%v",sequentailCall.CallName)

				// TODO: invoke funcExecutor with call names and test values to call other functions, based on the SUT
				funcExecutor(t, tester, sequentailCall.CallName, sequentailCall.TestValue)
			}
		}

		//TODO: compare test results against with test oracles, based on the SUT
	}
{{else}}
	for _, cs := range {{$name}}Test.TestCases{
	
	}
{{end}}
}
{{end}}

// funcExecutor is used to executed different functions.
func funcExecutor(t *testing.T, tester interface{}, funcName string, params ...interface{}) {
	t.Helper()
	inputArgs := make([]reflect.Value, len(params))
	for i, param := range params {
		inputArgs[i] = reflect.ValueOf(param)
	}
	fn := reflect.ValueOf(tester).MethodByName(funcName)
	if !fn.IsValid() {
		t.Errorf("method '%s' not found", funcName)
	}
	fn.Call(inputArgs)
}
`

var mainTestTemplateDefault = `package {{.PackageName}}

import (
	"flag"
	"os"
	"testing"
)

var (
    {{- range $name := .Tests }}
	{{$name}}Test  {{$name | stringTitle}}Test
	{{- end}}
)

func TestMain(m *testing.M) {
	// Flag definitions.
    {{- range $name := .Tests }}
	var {{$name}}TCsDir = flag.String(
		"{{$name}}TCsDir",
		"",
		"path to the file for {{$name | stringTitle}} tests",
	)
	{{- end}}

    // Parse and validate flags.
	flag.Parse()

	// Load test cases from XML files:
	{{- range $name := .Tests }}
	ParseXMLTestCase(*{{$name}}TCsDir, &{{$name}}Test)
	{{- end}}
    
	// Run tests/benchmarks.
	res := m.Run()
	os.Exit(res)
}

{{range $name := .Tests }}  
func Test{{$name | stringTitle}}(t *testing.T){
	
	for _, cs := range {{$name}}Test.TestCases{
	
	}
}
{{end}}
`
