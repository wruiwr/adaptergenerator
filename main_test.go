package gotestgen

import (
	"flag"
	"os"
	"reflect"
	"sync"
	"testing"
)

type tester struct {
	t *testing.T
}

var (
	systemTest SystemTest
)

func TestMain(m *testing.M) {
	// Flag definitions.
	var systemTCsDir = flag.String(
		"systemTCsDir",
		"",
		"path to the file for System tests",
	)

	// Parse and validate flags.
	flag.Parse()

	// Load test cases from XML files:
	ParseXMLTestCase(*systemTCsDir, &systemTest)

	// Run tests/benchmarks.
	res := m.Run()
	os.Exit(res)
}

func TestSystem(t *testing.T) {

	tester := &tester{t}
	t.Logf("test name=%v", systemTest.TestName)

	for _, cs := range systemTest.TestCases {
		t.Logf("test case ID=%v", cs.CaseID)
		// concurrent executions
		if cs.TestValues.ConcurrentCalls == nil {
			t.Log("no concurrent executions")
		} else {
			var wg sync.WaitGroup
			wg.Add(len(cs.TestValues.ConcurrentCalls))
			for _, concurrentCall := range cs.TestValues.ConcurrentCalls {
				t.Logf("concurrent call name=%v", concurrentCall.CallName)
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
		} else {
			for _, sequentailCall := range cs.TestValues.SequentialCalls {
				t.Logf("sequentail call name=%v", sequentailCall.CallName)

				// TODO: invoke funcExecutor with call names and test values to call other functions, based on the SUT
				funcExecutor(t, tester, sequentailCall.CallName, sequentailCall.TestValue)
			}
		}

		//TODO: compare test results against with test oracles, based on the SUT
	}

}

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
