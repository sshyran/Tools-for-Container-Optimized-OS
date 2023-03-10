package cloudlogger

import (
	"io/ioutil"
	"testing"
	"time"

	"cloud.google.com/go/logging"
	"cos.googlesource.com/cos/tools.git/src/pkg/nodeprofiler/profiler"
	"cos.googlesource.com/cos/tools.git/src/pkg/nodeprofiler/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// fakeCPU is a struct that implements the Component interface.
type fakeCPU struct {
	CPUName string
	Metrics *profiler.USEMetrics
}

// CollectUtilization behavior with regards to type fakeCPU.
func (f *fakeCPU) CollectUtilization(outputs map[string]utils.ParsedOutput) error {
	f.Metrics.Utilization = float64(7) // an arbitrary number
	return nil
}

// CollectSaturation behavior with regards to type fakeCPU.
func (f *fakeCPU) CollectSaturation(outputs map[string]utils.ParsedOutput) error {
	f.Metrics.Saturation = true
	return nil
}

// CollectErrors behavior with regards to type fakeCPU.
func (f *fakeCPU) CollectErrors(outputs map[string]utils.ParsedOutput) error {
	return nil
}

// USEMetrics behavior with regards to type fakeCPU.
func (f *fakeCPU) USEMetrics() *profiler.USEMetrics {
	return f.Metrics
}

// Name behavior with regards to type CPU.
func (f *fakeCPU) Name() string {
	return f.CPUName
}

// AdditionalInformation returns unique information for fakeCPU
func (f *fakeCPU) AdditionalInformation() string {
	return ""
}

// fakeMemCap is a struct that implement the Component interface.
type fakeMemCap struct {
	MemCapName string
	Metrics    *profiler.USEMetrics
}

// CollectUtilization behavior with regards to type fakeMemCap.
func (f *fakeMemCap) CollectUtilization(outputs map[string]utils.ParsedOutput) error {
	f.Metrics.Utilization = float64(7) // any arbitrary number
	return nil
}

// CollectSaturation behavior with regards to type fakeMemCap.
func (f *fakeMemCap) CollectSaturation(outputs map[string]utils.ParsedOutput) error {
	f.Metrics.Saturation = true
	return nil
}

// CollectErrors behavior with regards to type fakeMemCap.
func (f *fakeMemCap) CollectErrors(outputs map[string]utils.ParsedOutput) error {
	return nil
}

// USEMetrics behavior with regards to type fakeMemCap.
func (f *fakeMemCap) USEMetrics() *profiler.USEMetrics {
	return f.Metrics
}

// Name behavior with regards to type fakeMemCap.
func (f *fakeMemCap) Name() string {
	return f.MemCapName
}

// AdditionalInformation behaviro with regards to fakeMemCap.
func (f *fakeMemCap) AdditionalInformation() string {
	return ""
}

// fakeStorageDevIO is a struct that implement the Component interface.
type fakeStorageDevIO struct {
	StorageDevIOName string
	Metrics          *profiler.USEMetrics
}

// CollectUtilization behavior with regards to type fakeStorageDevIO.
func (f *fakeStorageDevIO) CollectUtilization(outputs map[string]utils.ParsedOutput) error {
	f.Metrics.Utilization = float64(7) // any arbitrary number
	return nil
}

// CollectSaturation behavior with regards to type fakeStorageDevIO.
func (f *fakeStorageDevIO) CollectSaturation(outputs map[string]utils.ParsedOutput) error {
	f.Metrics.Saturation = true
	return nil
}

// CollectErrors behavior with regards to type fakeStorageDevIO.
func (f *fakeStorageDevIO) CollectErrors(outputs map[string]utils.ParsedOutput) error {
	return nil
}

// USEMetrics behavior with regards to type fakeStorageDevIO.
func (f *fakeStorageDevIO) USEMetrics() *profiler.USEMetrics {
	return f.Metrics
}

// Name behavior with regards to type fakeStorageDevIO.
func (f *fakeStorageDevIO) Name() string {
	return f.StorageDevIOName
}

func (f *fakeStorageDevIO) AdditionalInformation() string {
	return ""
}

// fakeTextLogger is a struct that implements the TextLogger interface.
type fakeTextLogger struct {
	logged string
}

// fakeStructuredLogger is a struct that implements the StructuredLogger
// interface.
type fakeStructuredLogger struct {
	buffer []logging.Entry
	logged []logging.Entry
}

// Printf behavior with regards to type fakeTextLogger.
func (f *fakeTextLogger) Printf(text string, a ...interface{}) {
	f.logged = text
}

// Log behavior with regards to type fakeStructuredLogger.
func (f *fakeStructuredLogger) Log(entry logging.Entry) {
	f.buffer = append(f.buffer, entry)
}

// Flush behavior with regards to type fakeStructuredLogger.
func (f *fakeStructuredLogger) Flush() error {
	f.logged = append(f.logged, f.buffer...)
	f.buffer = []logging.Entry{}
	return nil
}

// generateFakeProfilerOpts initializes profiler components and commands
// and returns them.
func generateFakeProfilerOpts() ([]profiler.Component, []profiler.Command) {
	fCPU := &fakeCPU{
		CPUName: "fakeCPU",
		Metrics: &profiler.USEMetrics{},
	}
	fMemCap := &fakeMemCap{
		MemCapName: "fakeMemCap",
		Metrics:    &profiler.USEMetrics{},
	}
	fStorageDevIO := &fakeStorageDevIO{
		StorageDevIOName: "fakeStorageDevIO",
		Metrics:          &profiler.USEMetrics{},
	}
	// populating fake components.
	components := []profiler.Component{fCPU, fMemCap, fStorageDevIO}
	// populating fake commands.
	cmds := []profiler.Command{}
	return components, cmds
}

// For every input, the logged string must be the expected output unless the
// input is an empty string. In that case, nothing gets logged.
func TestTableLogText(t *testing.T) {
	var tests = []struct {
		name       string
		input      string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "non-empty string log",
			input:      "Node Profiler",
			wantOutput: "Node Profiler",
			wantErr:    false,
		},
		{
			name:       "empty string log",
			input:      "",
			wantOutput: "",
			wantErr:    true,
		},
	}

	for _, test := range tests {
		var f *fakeTextLogger = &fakeTextLogger{}
		err := LogText(f, test.input)
		// err will not be nil if the user attempted to log an empty string.
		// ignoring the case in which the user logged empty string.
		if gotErr := err != nil; gotErr != test.wantErr {

			t.Errorf("LogText(%v, %v) = %v, wantErr %t", f, test.input, err, test.wantErr)
		}
		if diff := cmp.Diff(f.logged, test.wantOutput); diff != "" {
			t.Errorf("ran LogText(fakeTextLogger, %+v), but got mismatch between got and want (-got, +want): \n diff %s", test.input, diff)
		}
	}
}

func TestTableLogProfilerReport(t *testing.T) {
	// Retrieving testing data.
	inputFile1, inputFile2 := "testdata/testdata.txt", "testdata/testdata2.txt"
	inputFileData1, err := ioutil.ReadFile(inputFile1)
	if err != nil {
		t.Errorf("failed to open testing input file %v: %v\n", inputFile1, err)
	}
	inputFileData2, err := ioutil.ReadFile(inputFile2)
	if err != nil {
		t.Errorf("failed to open testing input file %v: %v\n", inputFile2, err)
	}
	// Retrieving profiler components and commands.
	components, cmds := generateFakeProfilerOpts()
	useReport := profiler.USEReport{
		Components: components,
	}
	var cInfos []componentInfo
	expected := struct {
		Metrics *profiler.USEMetrics
	}{
		Metrics: &profiler.USEMetrics{
			Timestamp:   time.Date(2021, time.July, 21, 9, 59, 30, 0, time.UTC),
			Interval:    time.Date(2021, time.July, 21, 10, 3, 0, 0, time.UTC).Sub(time.Date(2021, time.July, 21, 9, 59, 30, 0, time.UTC)),
			Utilization: 7,
			Saturation:  true,
			Errors:      0,
		},
	}
	for _, c := range useReport.Components {
		cInfos = append(cInfos, componentInfo{Name: c.Name(), Metrics: expected.Metrics, Additional: c.AdditionalInformation()})
	}

	var tests = []struct {
		name       string
		input      *LoggerOpts
		wantOutput []logging.Entry
		wantErr    bool
	}{
		{
			name: "valid logger options and non-empty json payload log.",
			input: &LoggerOpts{
				ProjID: "cos-interns-playground",
				ShCmds: []ShellCmdOpts{
					ShellCmdOpts{
						Command:     "bash testdata/testcmd.sh",
						CmdCount:    1,
						CmdInterval: 0 * time.Second,
						CmdTimeOut:  3 * time.Second,
					},
				},
				ProfilerCount:    1,
				ProfilerInterval: 0 * time.Second,
				Components:       components,
				ProfilerCmds:     cmds,
			},
			wantOutput: []logging.Entry{
				{
					Payload: struct {
						CommandName   string
						CommandOutput string
					}{
						CommandName:   "bash testdata/testcmd.sh",
						CommandOutput: string(inputFileData1),
					},
					Severity: logging.Debug,
				},

				{
					Payload: struct {
						Components []componentInfo
						Analysis   string
					}{
						Components: cInfos,
						Analysis:   useReport.Analysis,
					},
					Severity: logging.Debug,
				}},
			wantErr: false,
		},
		{
			name: "multiple commands executions and multiple profiler runs non-empty json payload log.",
			input: &LoggerOpts{
				ProjID: "cos-interns-playground",
				ShCmds: []ShellCmdOpts{
					ShellCmdOpts{
						Command:     "bash testdata/testcmd.sh",
						CmdCount:    1,
						CmdInterval: 0 * time.Second,
						CmdTimeOut:  3 * time.Second,
					}, ShellCmdOpts{
						Command:     "bash testdata/testhello.sh",
						CmdCount:    1,
						CmdInterval: 0 * time.Second,
						CmdTimeOut:  3 * time.Second,
					},
				},
				ProfilerCount:    2,
				ProfilerInterval: 0 * time.Second,
				Components:       components,
				ProfilerCmds:     cmds,
			},
			wantOutput: []logging.Entry{
				{
					Payload: struct {
						CommandName   string
						CommandOutput string
					}{
						CommandName:   "bash testdata/testcmd.sh",
						CommandOutput: string(inputFileData1),
					},
					Severity: logging.Debug,
				},

				{
					Payload: struct {
						CommandName   string
						CommandOutput string
					}{
						CommandName:   "bash testdata/testhello.sh",
						CommandOutput: string(inputFileData2),
					},
					Severity: logging.Debug,
				},

				{
					Payload: struct {
						Components []componentInfo
						Analysis   string
					}{
						Components: cInfos,
						Analysis:   useReport.Analysis,
					},
					Severity: logging.Debug,
				}, {
					Payload: struct {
						Components []componentInfo
						Analysis   string
					}{
						Components: cInfos,
						Analysis:   useReport.Analysis,
					},
					Severity: logging.Debug,
				}},
			wantErr: false,
		},
		{
			name: "invalid logger options payload log: empty command with CmdCount and/or CmdInterval.",
			input: &LoggerOpts{
				ProjID: "cos-interns-playground",
				ShCmds: []ShellCmdOpts{
					ShellCmdOpts{
						Command:     "",
						CmdCount:    1,
						CmdInterval: 0 * time.Second,
						CmdTimeOut:  3 * time.Second,
					},
				},
				ProfilerCount:    1,
				ProfilerInterval: 0 * time.Second,
				Components:       components,
				ProfilerCmds:     cmds,
			},
			wantOutput: nil,
			wantErr:    true,
		},
		{
			name: "invalid logger options payload log: inconsistent CmdCount and CmdInterval.",
			input: &LoggerOpts{
				ProjID: "cos-interns-playground",
				ShCmds: []ShellCmdOpts{
					ShellCmdOpts{
						Command:     "bash testdata/testcmd.sh",
						CmdCount:    0,
						CmdInterval: 3 * time.Second,
						CmdTimeOut:  3 * time.Second,
					},
				},
				ProfilerCount:    1,
				ProfilerInterval: 0 * time.Second,
				Components:       components,
				ProfilerCmds:     cmds,
			},
			wantOutput: nil,
			wantErr:    true,
		},
		{
			name: "invalid logger options payload log: inconsistent ProfilerCount and ProfilerInterval.",
			input: &LoggerOpts{
				ProjID: "cos-interns-playground",
				ShCmds: []ShellCmdOpts{
					ShellCmdOpts{
						Command:     "bash testdata/testcmd.sh",
						CmdCount:    0,
						CmdInterval: 3 * time.Second,
						CmdTimeOut:  3 * time.Second,
					},
				},
				ProfilerCount:    0,
				ProfilerInterval: 4 * time.Second,
				Components:       components,
				ProfilerCmds:     cmds,
			},
			wantOutput: nil,
			wantErr:    true,
		},
		{
			name: "invalid logger options payload log: no project ID/invalid project ID.",
			input: &LoggerOpts{
				ProjID: "",
				ShCmds: []ShellCmdOpts{
					ShellCmdOpts{
						Command:     "bash testdata/testcmd.sh",
						CmdCount:    0,
						CmdInterval: 3 * time.Second,
						CmdTimeOut:  3 * time.Second,
					},
				},
				ProfilerCount:    1,
				ProfilerInterval: 0 * time.Second,
				Components:       components,
				ProfilerCmds:     cmds,
			},
			wantOutput: nil,
			wantErr:    true,
		},
	}
	for _, test := range tests {
		var f *fakeStructuredLogger = &fakeStructuredLogger{}
		err := LogProfilerReport(f, test.input)
		if gotErr := err != nil; gotErr != test.wantErr {
			t.Errorf("LogProfilerReport(%v, %v) = %v, wantErr %t", f, test.input, err, test.wantErr)
		}

		if diff := cmp.Diff(f.logged, test.wantOutput, cmpopts.IgnoreFields(profiler.USEMetrics{}, "Timestamp", "Interval")); diff != "" {
			t.Errorf("ran LogProfilerReport(fakeStructuredLogger,%+v), but got mismatch between got and want (-got, +want): \n diff %s", test.input, diff)

		}
	}
}
