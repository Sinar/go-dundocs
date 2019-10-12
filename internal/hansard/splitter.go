package hansard

type SplitHansardDocumentPlan struct {
	WorkingDir      string
	PlanDir         string
	HansardDocument HansardDocument
}

func NewSplitHansardDocumentPlan() *SplitHansardDocumentPlan {
	splitPlan := SplitHansardDocumentPlan{}

	return &splitPlan
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
