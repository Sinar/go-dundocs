package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/Sinar/go-dundocs"

	"github.com/turbinelabs/cli"
	"github.com/turbinelabs/cli/command"
)

func main() {
	fmt.Println("Welcome to go-dundocs!! Another project from SinarProject!!!")
	// run the Main function, which calls os.Exit with the appropriate exit status
	mkSubCmdCLI().Main()
}

type planRunner struct {
	sessionDUNName string
	sourcePDFPath  string
}

func CmdPlan() *command.Cmd {
	runner := &planRunner{}
	cmd := command.Cmd{
		Name:        "plan",
		Summary:     "Plan based on source PDF",
		Usage:       "--session=duntest-sesi42 --source=raw/BukanLisan/abc.pdf",
		Description: "Analyze the source PDF, plan the split which is stored into the data folder by default",
		Runner:      runner,
	}

	cmd.Flags.StringVar(&runner.sourcePDFPath, "source", "", "[Required] Source PDF to Analyze")
	cmd.Flags.StringVar(&runner.sessionDUNName, "session", "", "[Required] State Assembly Session Name to label the split by Questions PDFs")

	return &cmd
}

func (f *planRunner) Run(cmd *command.Cmd, args []string) command.CmdErr {
	//  Pre-req checks
	var errMessage string
	if f.sessionDUNName == "" || f.sourcePDFPath == "" {
		errMessage = `
**************** ERROR *************************
================================================
Ensure DUN Session Name and Source PDFs entered!
=================================================
`
		return command.CmdErr{
			Cmd:     nil,
			Code:    command.CmdErrCodeBadInput,
			Message: errMessage,
		}
	}
	// Get started with the Plan Operation
	fmt.Println("Pre-reqs passed; now running Plan!!")
	conf := dundocs.Configuration{
		SourcePDFPath: f.sourcePDFPath,
		//WorkingDir:    "",
		//DataDir:       "",
	}
	dd := dundocs.DUNDocs{
		DUNSession: f.sessionDUNName,
		Conf:       conf,
		Options:    nil,
	}
	// If need to debug ..
	if globalFlags.verbose {
		spew.Dump(dd)
	}
	// Execute the Plan
	dd.Plan()

	return command.NoError()
}

// The typical pattern is to provide a public CmdXYZ() func for each
// sub-command you wish to provide. This function should initialize the
// command.Cmd, the command.Runner, and flags.

// The private command.Runner implementation should contain any state needed
// to execute the command. The values should be initialized via flags declared
// in the CmdXYZ() function.
type splitRunner struct {
	sourcePDFPath string
}

func CmdSplit() *command.Cmd {
	// typically the command.Runner is initialized only with internally-defined
	// state; all necessary external state should be provided via flags. One can
	// inline the initializaton of the command.Runner in the command.Cmd
	// initialization if no flags are necessary, but it's often convenient to
	// have a typed reference
	runner := &splitRunner{}

	cmd := &command.Cmd{
		Name:        "split",
		Summary:     "Split source PDF following plan",
		Usage:       "--source=raw/Lisan/def.pdf",
		Description: "Load the planned Split Instruction and Execute Split on Source PDF",
		Runner:      runner,
	}

	cmd.Flags.StringVar(&runner.sourcePDFPath, "source", "", "[Required] Source PDF to Analyze")

	return cmd
}

// Run does the actual work, based on state provided by flags, and the
// args remaining after the flags have been parsed.
func (f *splitRunner) Run(cmd *command.Cmd, args []string) command.CmdErr {
	// argument validation should occur at the top of the function, and
	// errors should be reported via the cmd.BadInput or cmd.BadInputf methods
	//if len(args) < 1 {
	//	return cmd.BadInput("missing \"string\" argument.")
	//}
	//str := args[0]
	//if globalFlags.verbose {
	//	fmt.Printf("Splitting \"%s\"\n", str)
	//}
	//split := strings.Split(str, f.delim)
	//for i, term := range split {
	//	if globalFlags.verbose {
	//		fmt.Printf("[%d] ", i)
	//	}
	//	fmt.Println(term)
	//}

	//  Pre-req checks
	var errMessage string
	if f.sourcePDFPath == "" {
		errMessage = `
***** ERROR ******************
==============================
Ensure Source PDFs entered!
==============================
`
		return command.CmdErr{
			Cmd:     nil,
			Code:    command.CmdErrCodeBadInput,
			Message: errMessage,
		}
	}
	// Get started with the Split Operation
	fmt.Println("Pre-reqs passed; now running Split!!")
	conf := dundocs.Configuration{
		SourcePDFPath: f.sourcePDFPath,
		//WorkingDir:    "",
		//DataDir:       "",
	}
	dd := dundocs.DUNDocs{
		//DUNSession: f.sessionDUNName,
		Conf:    conf,
		Options: nil,
	}
	// For debugging purpose
	if globalFlags.verbose {
		spew.Dump(dd)
	}
	// Execute the Split action ..
	dd.Split()
	// In this case, there was no error. Errors should be returned via the
	// cmd.Error or cmd.Errorf methods.
	return command.NoError()
}

// while not mandatory, keeping globally-configured flags in a single struct
// makes it obvious where they came from at access time.
type globalFlagsT struct {
	verbose bool
}

var globalFlags = globalFlagsT{}

func mkSubCmdCLI() cli.CLI {
	// make a new CLI passing the description and version and one or more sub commands
	c := cli.NewWithSubCmds(
		"CLI to process PDFs from Selangor State Assembly",
		"1.0.0",
		CmdPlan(),
		CmdSplit(),
	)

	// Global flags can be used to modify global state
	c.Flags().BoolVar(&globalFlags.verbose, "verbose", false, "Produce verbose output")

	return c
}
