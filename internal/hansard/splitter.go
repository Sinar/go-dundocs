package hansard

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Configuration of a Context from outside-in ..
type Configuration struct {
	// DUN Session Label
	DUNSession string

	// ./raw + ./data folders are assumed to be relative to this dir
	WorkingDir string

	// Data directory name; can be relative or absolute?
	DataDir string

	// Source PDF can be anywhere; maybe make it a Reader to be read direct from S3?
	SourcePDFPath string

	// Options?
	Options *ExtractPDFOptions
}

type SplitHansardDocumentPlan struct {
	dataDir         string
	PlanDir         string
	HansardDocument HansardDocument
}

func GetAbsoluteDataDir(workingDir, dataDir string) string {
	// If absolute already not needed to do anything ..
	if filepath.IsAbs(dataDir) {
		return dataDir
	}
	// OK, now apply the rules
	var absoluteDataDir string
	// If no dataDir; do a default
	if dataDir == "" {
		absoluteDataDir = workingDir + "/data"
	} else {
		absoluteDataDir = workingDir + fmt.Sprintf("/%s", dataDir)
	}
	// DEBUG
	//if absoluteDataDir == "" {
	//	panic(fmt.Errorf("BEFORE: %s AFTER: %s", dataDir, absoluteDataDir))
	//}
	return absoluteDataDir
}

func NewEmptySplitHansardDocumentPlan(absoluteDataDir, absolutePlanDir, sessionName string) *SplitHansardDocumentPlan {
	// Assume: sourcePDFFilename stripped off; validation here??
	// Assume: dataDir and PlanDir must become absolute before passing it back? Validate?
	if !(filepath.IsAbs(absoluteDataDir) && filepath.IsAbs(absolutePlanDir)) {
		panic(fmt.Errorf("DATA: %s + PLAN: %s MUST BE ABSOLUTE!", absoluteDataDir, absolutePlanDir))
	}
	// If absolute dataDir; just take it  as is, no use for workingDir
	//absoluteDataDir := GetAbsoluteDataDir(workingDir, dataDir)
	// Extract out filename as folder for split.yml plan
	// https://stackoverflow.com/questions/13027912/trim-strings-suffix-or-extension
	//basePDFPath := filepath.Base(sourcePDFPath)
	//planDir := absoluteDataDir + fmt.Sprintf("/%s", strings.TrimSuffix(basePDFPath, filepath.Ext(basePDFPath)))
	//// DEBUG
	//fmt.Println("PLAN_PATH: ", planDir)
	// Do abs conversion here?? for PlanDir only? Is it needed; is relative good enough? Maybe ..
	//<TODO>??
	// Assemble the pieces here ..
	splitPlan := SplitHansardDocumentPlan{
		dataDir: absoluteDataDir,
		PlanDir: absolutePlanDir,
		HansardDocument: HansardDocument{
			StateAssemblySession: sessionName,
			HansardQuestions:     []HansardQuestion{},
		},
	}
	return &splitPlan
}

func NewSplitHansardDocumentPlan(conf Configuration) *SplitHansardDocumentPlan {
	// If we need to customize any options; put  it above ..
	// Get PDF content
	pdfDoc, nperr := NewPDFDocument(conf.SourcePDFPath, conf.Options)
	if nperr != nil {
		panic(nperr)
	}
	// Once have the  content  to be  processed; pass it all on after using the config  to adjust things ..
	// TODO: Make the  dataDir into absolute  before  instantiate
	// It will later need to  be append with Type (string version) + filename .. + split.yml
	// TODO: Filename needs to be  extracted out; and need to handle those cases with '.' in filename
	//splitPlan := SplitHansardDocumentPlan{
	//	dataDir: conf.WorkingDir,
	//	PlanDir: planDir,
	//	HansardDocument: HansardDocument{
	//		StateAssemblySession: conf.DUNSession,
	//	},
	//}
	absoluteDataDir := GetAbsoluteDataDir(conf.WorkingDir, conf.DataDir)
	// Extract out filename as folder for split.yml plan
	// https://stackoverflow.com/questions/13027912/trim-strings-suffix-or-extension
	basePDFPath := filepath.Base(conf.SourcePDFPath)
	absolutePlanDir := absoluteDataDir + fmt.Sprintf("/%s", strings.TrimSuffix(basePDFPath, filepath.Ext(basePDFPath)))
	splitPlan := NewEmptySplitHansardDocumentPlan(absoluteDataDir, absolutePlanDir, conf.DUNSession)
	// Fill in the needed  plan
	err := NewSplitHansardDocumentPlanContent(pdfDoc, splitPlan)
	if err != nil {
		panic(err)
	}
	// DEBUG
	//spew.Dump(splitPlan)
	return splitPlan
}

func NewSplitHansardDocumentPlanContent(pdfDoc *PDFDocument, splitPlan *SplitHansardDocumentPlan) error {
	err := NewHansardDocumentContent(pdfDoc, &splitPlan.HansardDocument)
	if err != nil {
		return err
	}
	// Some basic validtion check; need to be valid type
	if splitPlan.HansardDocument.HansardType == HANSARD_INVALID {
		return errors.New("INVALID PLAN!!")
	}
	if len(splitPlan.HansardDocument.HansardQuestions) == 0 {
		return errors.New("EMPTY PLAN!!")
	}
	return nil
}

func (s *SplitHansardDocumentPlan) SavePlan() error {
	// Persist HansardDoc into storage; whatever it may be ..
	// we know where the plan is to be saved ..
	// Pre-req plan checks first .. is this a dupe?
	if s.HansardDocument.HansardType == HANSARD_INVALID {
		return errors.New("INVALID PLAN!!")
	}
	if len(s.HansardDocument.HansardQuestions) == 0 {
		return errors.New("EMPTY PLAN!!")
	}
	// Make dir if needed ..
	mkerr := os.MkdirAll(s.PlanDir, 0744)
	if mkerr != nil {
		// If no permission; die!
		if os.IsPermission(mkerr) {
			panic(mkerr)
		}
	}
	// OK, pre-reqs done
	b, merr := yaml.Marshal(s.HansardDocument)
	if merr != nil {
		return merr
	}
	// Write into the pre-defined filename ..
	werr := ioutil.WriteFile(s.PlanDir+"/split.yml", b, 0644)
	if werr != nil {
		return werr
	}
	// All okie ..
	return nil
}

func (s *SplitHansardDocumentPlan) LoadPlan() error {
	// Load into the struct HansardDoc from the persistent storage ..
	hansardDoc := HansardDocument{}
	b, rerr := ioutil.ReadFile(s.PlanDir)
	if rerr != nil {
		return rerr
	}
	umerr := yaml.Unmarshal(b, &hansardDoc)
	if umerr != nil {
		return umerr
	}
	// attach plan
	s.HansardDocument = hansardDoc
	// All OK!
	return nil
}

func (s *SplitHansardDocumentPlan) ExecuteSplit(absoluteSrcPDF, absoluteSplitOutput string) error {
	// Assume: absoluteSrcPDF must be absolute before passing it back? Validate?
	if !(filepath.IsAbs(absoluteSrcPDF)) {
		panic(fmt.Errorf("PDF: %s MUST BE ABSOLUTE!", absoluteSrcPDF))
	}

	if len(s.HansardDocument.HansardQuestions) == 0 {
		return fmt.Errorf("Empty %s", absoluteSrcPDF)
	}
	fmt.Println("INSIDE: XSPLIT; splitting ", absoluteSrcPDF)
	// Traverse the HansardDoc
	// Traverse  each question in the HansardDoc
	for _, hansardQuestion := range s.HansardDocument.HansardQuestions {
		// DEBUG!
		//spew.Dump(hansardQuestion)
		finalFileName := "wassup!!"
		// DO the actuak split ..
		fmt.Println("Split Question: ", hansardQuestion.QuestionNum)
		// Write the splitoutput  to this final location
		fmt.Println("Save  split file into Dir: "+absoluteSplitOutput+" w/ fileName: ", finalFileName)
	}

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
