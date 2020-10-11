# gotestgen

### This tool can generate go code for the test adapter by using a template.

# gotestgen

### This tool can generate go code for the test adapter by using templates.

## Templates
Currently there are 4 templates:
1. for generating reader: <br />
    a. readerTemplateConcurrent <br />
    b. readerTemplateDefault
2. for generating tester: <br />
    c. mainTestTemplateConcurrent <br />
    d. mainTestTemplateDefault

These templates can be chosen based on the "Type" tag in the xml file of each test. <br />
For example: Type="concurrentsystemtest" can lead to templates a and c to be chosen.       

## How to use
1.in reader.go, users need to provide the source directory of generated xml files 
and the test names (in lower case). <br />
For example: <br />
    a. //go:generate testgen -dir=../gotestgen/xml/ -tests=system <br />
    b. //go:generate testgen -dir=../gotestgen/xml/ -tests=system,prepareQF,acceptQF

2.the generated xml file should provide some basic information about the test,
which includes "MethodCalls", "SystemParameterTypes", "InputType", "OracleType" (all with names) and "SystemParameters". <br />
For example: 
```xml
    <MethodCalls>
        <CallName>SendVote</CallName>
    </MethodCalls>
    <TypeAssignments>
        <SystemParameterTypes>
            <Field Name="NumberOfWorker" Type="int"></Field>
        </SystemParameterTypes>
        <InputType Name="VoteInput" Type="VoteInput">
            <Field Name="WorkerID" Type="WorkerID"></Field>
            <Field Name="VoteValue" Type="VoteEnum"></Field>
        </InputType>
        <OracleType Name="DecisionOracle" Type="DecisionSlice">
            <Field Name="WorkerID" Type="WorkerID"></Field>
            <Field Name="DecisionValue" Type="DecisionEnum"></Field>
        </OracleType>
        <OracleType Name="FinalDecision" Type="DecisionEnum"></OracleType>
    </TypeAssignments>
    <SystemParameters>
        <NumberOfWorker>3</NumberOfWorker>
    </SystemParameters>
```    

For any method names provided by "MethodCalls", the tester of the generated test adapter has an function executor
"funcExecutor" to execute these method calls to interact with SUT:
``` 
// funcExecutor is used to executed different method calls to SUT.
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
``` 
Therefore, based on call names of "MethodCalls", an interface will be generated in the tester of the generated test
adapter, for example:
```go   
...
type systemTest interface {
	Start()
	SendVote(v *TestValue)
	OutputChecker(cs *SystemTestCases)
}
...
...
// TODO: implement start method to setup and start SUT
func (t tester) Start() {}
// TODO: implement method call to SUT
func (t tester) SendVote(v *TestValue) {}
// TODO: implement outputChecker method
func (t tester) OutputChecker(cs *SystemTestCases) {} 
```   
The user needs to implement these methods to interact with the SUT. <br />
3.if the test involves concurrent executions, then the xml format should give the information of 
"concurrent" and "sequential" execution call names. <br />
for example, the test case below involves two concurrent calls ("SendVote") executed concurrently 
and followed by another "SendVote" call sequentially:
```xml
 <TestCase ID="1">
        <InputValues>
            <Concurrent>
                <Call Name="SendVote">
                    <InputValue>
                        <Vote>
                            <WorkerID>1</WorkerID>
                            <VoteValue>0</VoteValue>
                        </Vote>
                    </InputValue>
                </Call>
                <Call Name="SendVote">
                    <InputValue>
                        <Vote>
                            <WorkerID>2</WorkerID>
                            <VoteValue>1</VoteValue>
                        </Vote>
                    </InputValue>
                </Call>
            </Concurrent>
            <Call Name="SendVote">
                <InputValue>
                    <Vote>
                        <WorkerID>3</WorkerID>
                        <VoteValue>0</VoteValue>
                    </Vote>
                </InputValue>
            </Call>
        </InputValues>
        <Oracles>
            <Decision>
                <WorkerID>1</WorkerID>
                <DecisionValue>1</DecisionValue>
            </Decision>
            <Decision>
                <WorkerID>2</WorkerID>
                <DecisionValue>1</DecisionValue>
            </Decision>
            <FinalDecision>1</FinalDecision>
        </Oracles>
    </TestCase>
```   
## Examples
All current examples to generate reader and tester are in xml folder. <br />
testsystemtpc.xml is an example xml format for generating a test adapter 
that can perform concurrent executions of the system by invoking different calls provided in the xml file. <br />
testsystem_paxos.xml and testunit_paxos.xml are the example of xml format for generating a test adapter for testing non concurrent executions.

## TO RUN
In terminal, under gotestgen, run: make <br /> 
(note: users need to choose where they want to install the tool by updating Makefile)
