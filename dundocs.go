package dundocs

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/Sinar/go-dundocs/internal/hansard"
)

type DUNDocs struct {
	// DUN Session Label
	DUNSession string

	Conf Configuration

	Options *ExtractPDFOptions
}

// CommandMode specifies the operation being executed.
type CommandMode int

// The available commands.
const (
	PLAN CommandMode = iota
	SPLIT
	RESET
)

type ExtractPDFOptions struct {
	StartPage int
	NumPages  int
}

type Configuration struct {
	// Source PDF can be anywhere; maybe make it a Reader to be read direct from S3?
	SourcePDFPath string

	// ./raw + ./data folders are assumed to be relative to this dir
	WorkingDir string

	// Data directory name; can be relative or absolute?
	DataDir string
}

func NewDUNDocs() *DUNDocs {
	dunDocs := DUNDocs{}
	return &dunDocs
}
func (dd *DUNDocs) Plan() {
	log.Println("In Plan ..")
	//pdfPath := dd.Conf.SourcePDFPath
	//conf := Configuration{}
	//splitPlan := hansard.NewSplitHansardDocumentPlan(
	//	dd.Conf.SourcePDFPath, dd.Conf.WorkingDir,
	//	dd.Conf.DataDir, dd.DUNSession,
	//	dd.Options.StartPage, dd.Options.NumPages)
	//// Perissit it
	//splitPlan.SavePlan()

	var options *hansard.ExtractPDFOptions
	// Fill in the options if needed ..
	if dd.Options != nil {
		options = &hansard.ExtractPDFOptions{
			StartPage: dd.Options.StartPage,
			NumPages:  dd.Options.NumPages,
		}
	}
	// Prepare config
	conf := hansard.Configuration{
		DUNSession:    dd.DUNSession,
		WorkingDir:    dd.Conf.WorkingDir,
		DataDir:       dd.Conf.DataDir,
		SourcePDFPath: dd.Conf.SourcePDFPath,
		Options:       options,
	}
	// DEBUG
	//spew.Dump(c)
	//err := hansard.PlanAndSave(conf)
	//if err != nil {
	//	panic(err)
	//}
	// Above  no need; as equivalent to below; can be r efectored out I guess ..
	splitPlan := hansard.NewSplitHansardDocumentPlan(conf.SourcePDFPath, conf.WorkingDir, conf.DataDir, conf.DUNSession, conf.Options)
	// Try to persist the plan
	sperr := splitPlan.SavePlan()
	if sperr != nil {
		panic(sperr)
	}
	// if see flag; then call the following; not executed by default ..
	//hansard.LoadAndSplit()
}

func (dd *DUNDocs) Split() {
	log.Println("In Split ..")
	// Setup; possibly into a helper function? Seen it repeated a few places
	workingDir := "."
	if dd.Conf.WorkingDir != "" {
		workingDir = dd.Conf.WorkingDir
	}
	dataDir := ""
	if dd.Conf.DataDir != "" {
		dataDir = dd.Conf.DataDir
	}
	// Needs absolute Path for Source  .. absoluteSrcPDF, sperr := filepath.Abs(tt.fields.srcPDFPath)
	sourcePDFPath, sperr := filepath.Abs(dd.Conf.SourcePDFPath)
	if sperr != nil {
		panic(sperr)
	}
	absoluteDataDir := hansard.GetAbsoluteDataDir(workingDir, dataDir)
	// Extract out filename as folder for split.yml plan
	// https://stackoverflow.com/questions/13027912/trim-strings-suffix-or-extension
	basePDFPath := filepath.Base(sourcePDFPath)
	// Plan file name is hardcoded split.yml
	absolutePlanFile := absoluteDataDir + fmt.Sprintf("/%s/split.yml", strings.TrimSuffix(basePDFPath, filepath.Ext(basePDFPath)))
	// Maybe  above to be moved  into NewEmptySplitHansardDocumentPlan?
	splitPlan := hansard.NewEmptySplitHansardDocumentPlan(
		absoluteDataDir, absolutePlanFile, dd.DUNSession)
	// Load plan
	lderr := splitPlan.LoadPlan()
	if lderr != nil {
		panic(lderr)
	}
	// DEBUG
	//spew.Dump(splitPlan.HansardDocument)
	// Setup  + execute; use the default splitOutDir
	absoluteSplitOutput := hansard.GetAbsoluteSplitOutDir(workingDir, "")
	exerr := splitPlan.ExecuteSplit(sourcePDFPath, absoluteSplitOutput)
	if exerr != nil {
		panic(exerr)
	}
	// Below also looks like a recurring pattern; restructure it  somehow? Not needed?
	//var options *hansard.ExtractPDFOptions
	//// Fill in the options if needed ..
	//if dd.Options != nil {
	//	options = &hansard.ExtractPDFOptions{
	//		StartPage: dd.Options.StartPage,
	//		NumPages:  dd.Options.NumPages,
	//	}
	//}
	//// Prepare config
	//conf := hansard.Configuration{
	//	DUNSession:    dd.DUNSession,
	//	WorkingDir:    dd.Conf.WorkingDir,
	//	DataDir:       dd.Conf.DataDir,
	//	SourcePDFPath: dd.Conf.SourcePDFPath,
	//	Options:       options,
	//}
	//hansard.LoadAndSplit(conf)
}

func (dd *DUNDocs) Reset() {
	log.Println("In Reset ...")
	// Clean up plan
	// Clean up split pages folder
	// Clean up merged pages location
}
