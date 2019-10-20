package hansard

import (
	"github.com/davecgh/go-spew/spew"
)

// Configuration of a Context from outside-in ..
type Configuration struct {
	// DUN Session Label
	DUNSession string

	// ./raw + ./data folders are assumed to be relative to this dir
	WorkingDir string

	// Source PDF can be anywhere; maybe make it a Reader to be read direct from S3?
	SourcePDFPath string

	// Options?
	Options *ExtractPDFOptions
}

type SplitHansardDocumentPlan struct {
	workingDir      string
	PlanDir         string
	HansardDocument HansardDocument
}

func NewSplitHansardDocumentPlan(conf Configuration) *SplitHansardDocumentPlan {
	// If we need to customize any options; put  it above ..
	// Get PDF content
	pdfDoc, nperr := NewPDFDocument(conf.SourcePDFPath, conf.Options)
	if nperr != nil {
		panic(nperr)
	}
	// Once have the  content  to be  processed; pass it all on after using the config  to adjust things ..
	// If no PLanDir; do a default
	planDir := conf.WorkingDir + "/data"
	// TODO: Make the  workingDir into absolute  before  instantiate
	// It will later need to  be append with Type (string version) + filename .. + split.yml
	// TODO: Filename needs to be  extracted out; and need to handle those cases with '.' in filename
	splitPlan := SplitHansardDocumentPlan{
		workingDir: conf.WorkingDir,
		PlanDir:    planDir,
		HansardDocument: HansardDocument{
			StateAssemblySession: conf.DUNSession,
		},
	}
	// Fill in the needed  plan
	err := NewSplitHansardDocumentPlanContent(pdfDoc, &splitPlan)
	if err != nil {
		panic(err)
	}
	// DEBUG
	spew.Dump(splitPlan)
	return &splitPlan
}

func NewSplitHansardDocumentPlanContent(pdfDoc *PDFDocument, splitPlan *SplitHansardDocumentPlan) error {
	err := NewHansardDocumentContent(pdfDoc, &splitPlan.HansardDocument)
	if err != nil {
		return err
	}
	return nil
}

func (s *SplitHansardDocumentPlan) SavePlan() error {
	// Persist HansardDoc into storage; whatever it may be ..
	return nil
}

func (s *SplitHansardDocumentPlan) LoadPlan() error {
	// Load into the struct HansardDoc from the persistent storage ..
	return nil
}

func (s *SplitHansardDocumentPlan) ExecuteSplit() error {
	// Traverse the HansardDoc
	// Traverse  each question in the HansardDoc

	return nil
}

//  INternal helper methods
//  Use  whatever is available; start with the  PDF filename itself
// if being passed via command; not overwrite; only when  empty
func (s *SplitHansardDocumentPlan) detectSessionName(sourcePDFFileName string) string {
	return sourcePDFFileName
}

// Helper functions
func loadSplitHansardDocPlan(splitPlanPath string) *HansardDocument {
	splitHansardDocPlan := HansardDocument{}

	return &splitHansardDocPlan
}

// Peg the API on the v0.1.25 version; no  support for v0.2 yet :(
// OK we'll try the latest API first ..
func prepareSplitAPI() error {
	return nil
}

// Backup  command to run if API having  issues
func prepareSplit() error {
	return nil
}

// Normalize to Absolute Path
func normalizeToAbsolutePath(relativePath string) (absolutePath string, baseName string, extension string) {
	// Scenarios:
	// 	#1 Handle multiple '.' in filename
	//	#2 Windows path
	//	#3 No extensions
	//	#4 UTF-8 filenames?
	return absolutePath, baseName, extension
}

// For testing, ensure can get unique  TemoDir
func normalizeTempDirAbsolutePath(relativePath string) (absolutePath string, baseName string, extension string) {
	// Scenarios:
	// 	#1 Handle multiple '.' in filename
	//	#2 Windows path
	//	#3 No extensions
	//	#4 UTF-8 filenames?
	return absolutePath, baseName, extension
}
