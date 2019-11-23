package hansard

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	papi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
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
	// DEBUG
	//fmt.Println("INSIDE: XSPLIT; splitting ", absoluteSrcPDF)
	// Prepare the splitout folder scratch space + final destination
	scratchDir := filepath.Join(absoluteSplitOutput, "scratch")
	createDirIfNotExist(scratchDir)
	// Split the document into the  final destination
	pserr := prepareSplitAPI(absoluteSrcPDF, scratchDir)
	if pserr != nil {
		return pserr
	}
	// Traverse the HansardDoc
	// Traverse  each question in the HansardDoc
	for _, hansardQuestion := range s.HansardDocument.HansardQuestions {
		// DEBUG!
		//spew.Dump(hansardQuestion)
		// fmt.Sprintf("%s-soalan-%s.pdf", label, hq.QuestionNum)
		finalFileName := "wassup!!" // Is this to be derived? template?
		// Can also prepare the needed scratch spaces?
		// Derive  the needed basename ** TODO ***
		//srcBasename := "SRC_BASENAME"
		//fmt.Println("Basename derived from: ", absoluteSrcPDF, " is ", srcBasename)
		// ** TODO *** Prepare the location of merged directory too ..
		//// Ensure the merged directory is there ..
		//createDirIfNotExist(filepath.Join(absoluteSplitOutput, srcBasename))
		////finalMergedPDFPath := fmt.Sprintf("%s/splitout/%s-soalan-%s-%s.pdf", currentWorkingDir, label, hansardType, hq.QuestionNum)
		//finalMergedPDFPath := filepath.Join(absoluteSplitOutput, srcBasename, fmt.Sprintf("%s-soalan-%s.pdf", label, hq.QuestionNum))

		// DO the actuak split ..
		fmt.Println("Split Question: ", hansardQuestion.QuestionNum)
		// Write the splitoutput  to this final location
		fmt.Println("Save  split file into Dir: "+absoluteSplitOutput+" w/ fileName: ", finalFileName)
		ssqerr := splitSingleQuestion(
			s.HansardDocument.StateAssemblySession,
			absoluteSrcPDF, absoluteSplitOutput, finalFileName,
			hansardQuestion)
		if ssqerr != nil {
			// Need to determine if recoverable??
			// What is recoverable; for now; give  up
			//panic(ssqerr)
			return ssqerr
		}
		// Below is the nicer API I think ..
		// ssqerr := splitSingleQuestion(absoluteScratchDir, absoluteSplitOutputFMT, hansardQuestion)
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

// Use the latest API v2
func prepareSplitAPI(absoluteSrcPDF, scratchDir string) error {
	// Relax validation  --> https://github.com/hhrutter/pdfcpu/issues/80
	conf := pdfcpu.NewDefaultConfiguration()
	// DEBUG
	fmt.Println("In prepareSplitAPI!!")
	fmt.Println("Split ", absoluteSrcPDF, " to singles in ", scratchDir)
	sperr := papi.SplitFile(absoluteSrcPDF, scratchDir, 1, conf)
	if sperr != nil {
		return sperr
	}
	// DEBUG
	//q.Q(o)

	return nil
}

func splitSingleQuestion(label, absoluteSrcPDF, absoluteSplitOutput, finalFileName string, hq HansardQuestion) error {
	fmt.Println("LABEL: ", label)
	// Derive  the needed basename
	//srcBasename := "SRC_BASENAME"
	_, srcBasename, _ := normalizeToAbsolutePath(absoluteSrcPDF)

	fmt.Println("Basename derived from: ", absoluteSrcPDF, " is ", srcBasename)

	// Pre-reqs are done; now can start the split itself ..
	var pagesToMerge []string

	for i := hq.PageNumStart; i <= hq.PageNumEnd; i++ {
		//sourcePDFPath := fmt.Sprintf("%s/raw/splitout/%s/%s/pages/%s_%d.pdf", currentWorkingDir, hansardType, sessionName, sessionName, i)
		sourcePDFPath := filepath.Join(absoluteSplitOutput, "scratch", fmt.Sprintf("%s_%d.pdf", srcBasename, i))
		pagesToMerge = append(pagesToMerge, sourcePDFPath)
	}

	// Ensure the merged directory is there ..
	createDirIfNotExist(filepath.Join(absoluteSplitOutput, srcBasename))
	//finalMergedPDFPath := fmt.Sprintf("%s/splitout/%s-soalan-%s-%s.pdf", currentWorkingDir, label, hansardType, hq.QuestionNum)
	finalMergedPDFPath := filepath.Join(absoluteSplitOutput, srcBasename, fmt.Sprintf("%s-soalan-%s.pdf", label, hq.QuestionNum))
	fmt.Println(">>>=========== Merged file at: ", finalMergedPDFPath, " ==============<<<<<<")

	// Relax validation  --> https://github.com/hhrutter/pdfcpu/issues/80
	// Real-life data are pretty broken ..
	conf := pdfcpu.NewDefaultConfiguration()
	// Not needed
	//conf.ValidationMode = pdfcpu.ValidationRelaxed
	merr := papi.MergeFile(pagesToMerge, finalMergedPDFPath, conf)
	if merr != nil {
		return merr
	}
	return nil
}

// Backup  command to run if API having  issues
func prepareSplit() error {
	panic(fmt.Errorf("NOT IMPLEMENTED!!! FUNC: %s", "prepareSplit"))
	return nil
}

// Normalize to Absolute Path
func normalizeToAbsolutePath(relativePath string) (absoluteDir string, baseName string, extension string) {
	// Scenarios:
	// 	#1 Handle multiple '.' in filename
	//	#2 Windows path
	//	#3 No extensions
	//	#4 UTF-8 filenames?
	absolutePath, aerr := filepath.Abs(relativePath)
	if aerr != nil {
		panic(aerr)
	}
	//absoluteDir = filepath.Dir(absolutePath)
	//baseName = filepath.Base(absolutePath)
	absoluteDir, baseName = filepath.Split(absolutePath)
	filepath.SplitList(absolutePath)
	extension = filepath.Ext(absolutePath)
	if strings.ToLower(extension) == ".pdf" {
		baseName = strings.TrimSuffix(baseName, extension)
		// Order matters!
		extension = "pdf"
	} else {
		// Make ti explicit if not matcging PDF doc!
		extension = ""
	}
	// Any '.' transform to '_'
	baseName = strings.ReplaceAll(baseName, ".", "_")
	return absoluteDir, baseName, extension
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

// Create the needed directory if it does not exist; candidate for nonstdlib
// https://siongui.github.io/2017/03/28/go-create-directory-if-not-exist/
func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
