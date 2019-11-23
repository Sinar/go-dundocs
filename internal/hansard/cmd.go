package hansard

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

func PlanAndSave(conf Configuration) error {
	// Pre-req; setup absolute paths?

	// Get pdfdoc
	pdfDoc, pderr := NewPDFDocument(conf.SourcePDFPath, conf.Options)
	if pderr != nil {
		return pderr
	}
	// get splitplan
	splitPlan := NewEmptySplitHansardDocumentPlan(conf.DataDir, conf.DataDir, conf.DUNSession)
	// Extract out the needed content
	serr := NewSplitHansardDocumentPlanContent(pdfDoc, splitPlan)
	if serr != nil {
		return serr
	}
	// persist the plan + split ..
	splitPlan.SavePlan()
	// get back here; is ok
	return nil
}

func LoadAndSplit(conf Configuration) error {
	// Pre-req; setup absolute paths?

	splitPlan := NewEmptySplitHansardDocumentPlan(conf.DataDir, conf.DataDir, conf.DUNSession)

	splitPlan.LoadPlan()
	splitPlan.ExecuteSplit(conf.SourcePDFPath, conf.WorkingDir)

	return nil
}
