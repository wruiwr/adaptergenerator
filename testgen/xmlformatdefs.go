package main

import "encoding/xml"

type Global struct {
	XMLName         xml.Name         `xml:"Test"`
	TestName        string           `xml:"Name,attr"`
	TestType        string           `xml:"Type,attr"`
	TypeAssignments *TypeAssignments `xml:"TypeAssignments"`
}

type TypeAssignments struct {
	XMLName              xml.Name              `xml:"TypeAssignments"`
	SystemParameterTypes *SystemParameterTypes `xml:"SystemParameterTypes"`
	InputTypes           *[]InputType          `xml:"InputType"`
	OracleTypes          *[]OracleType         `xml:"OracleType"`
}

type SystemParameterTypes struct {
	XMLName xml.Name `xml:"SystemParameterTypes"`
	Fields  *[]Field `xml:"Field"`
}

type InputType struct {
	XMLName xml.Name `xml:"InputType"`
	Name    string   `xml:"Name,attr"`
	Type    string   `xml:"Type,attr"`
	Fields  *[]Field `xml:"Field"`
}

type OracleType struct {
	XMLName xml.Name `xml:"OracleType"`
	Name    string   `xml:"Name,attr"`
	Type    string   `xml:"Type,attr"`
	Fields  *[]Field `xml:"Field"`
}

type Field struct {
	XMLName xml.Name `xml:"Field"`
	Name    string   `xml:"Name,attr"`
	Type    string   `xml:"Type,attr"`
}
