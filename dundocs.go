package dundocs

import (
	"log"

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

	c := hansard.Configuration{
		DUNSession:    "",
		WorkingDir:    "",
		DataDir:       "",
		SourcePDFPath: "",
		Options:       nil,
	}
	hansard.PlanAndSave(c)

	// if see flag; then call the following; not executed by default ..
	//hansard.LoadAndSplit()
}

func (dd *DUNDocs) Split() {
	log.Println("In Split ..")
	// Load plan
	splitPlan := hansard.NewEmptySplitHansardDocumentPlan(
		dd.Conf.DataDir, "", dd.DUNSession)
	splitPlan.LoadPlan()
	splitPlan.ExecuteSplit(dd.Conf.SourcePDFPath, dd.Conf.DataDir)

	c := hansard.Configuration{
		DUNSession:    "",
		WorkingDir:    "",
		DataDir:       "",
		SourcePDFPath: "",
		Options:       nil,
	}
	hansard.LoadAndSplit(c)
}

func (dd *DUNDocs) Reset() {
	log.Println("In Reset ...")
	// Clean up plan
	// Clean up split pages folder
	// Clean up merged pages location
}
