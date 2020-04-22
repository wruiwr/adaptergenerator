package gotestgen

import "encoding/xml"

// SystemTest is an exported type for System tests
type SystemTest struct {
	XMLName          xml.Name           `xml:"Test"`
	TestName         string             `xml:"Name,attr"`
	TestType         string             `xml:"Type,attr"`
	SystemParameters *SystemParameters  `xml:"SystemParameters"`
	TestCases        []*SystemTestCases `xml:"TestCase"`
}

type SystemTestCases struct {
	XMLName     xml.Name           `xml:"TestCase"`
	CaseID      string             `xml:"ID,attr"`
	TestValues  *SystemTestValues  `xml:"InputValues"`
	TestOracles *SystemTestOracles `xml:"Oracles"`
}

type SystemTestValues struct {
	XMLName         xml.Name `xml:"InputValues"`
	ConcurrentCalls []*Calls `xml:"Concurrent>Call"`
	SequentialCalls []*Calls `xml:"Call"`
}

type Calls struct {
	XMLName   xml.Name   `xml:"Call"`
	CallName  string     `xml:"Name,attr"`
	TestValue *TestValue `xml:"InputValue"`
}

type TestValue struct {
	XMLName xml.Name `xml:"InputValue"`
	Vote    Vote     `xml:"Vote"`
}

type SystemTestOracles struct {
	XMLName       xml.Name    `xml:"Oracles"`
	Decisions     []*Decision `xml:"Decision"`
	FinalDecision int         `xml:"FinalDecision"`
}

type Vote struct {
	VoteValue int `xml:"VoteValue"`
	WorkerID  int `xml:"WorkerID"`
}

type Decision struct {
	DecisionValue int `xml:"DecisionValue"`
	WorkerID      int `xml:"WorkerID"`
}

type SystemParameters struct {
	XMLName        xml.Name `xml:"SystemParameters"`
	NumberOfWorker int      `xml:"NumberOfWorker"`
}
