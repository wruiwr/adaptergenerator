package gotestgen

import (
	"encoding/xml"
	"io/ioutil"
)

//go:generate testgen -dir=../gotestgen/xml/ -tests=system
func ParseXMLTestCase(file string, xmlTestCaseType interface{}) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, &xmlTestCaseType)
}
