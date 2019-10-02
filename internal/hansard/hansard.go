package hansard

import (
	"fmt"
	"strings"
)

type HansardType int

const (
	HANSARD_SPOKEN HansardType = iota
	HANSARD_WRITTEN
	HANSARD_DEBATE
)

type HansardQuestion struct {
	QuestionNum  string
	PageNumStart int
	PageNumEnd   int
}

type HansardDocument struct {
	StateAssemblySession string
	HansardType          HansardType
	HansardQuestions     []HansardQuestion
}

func NewHansardDocument(sessionName string, pdfPath string) (*HansardDocument, error) {
	// Load the PDFDoc; should we check length?
	pdfDoc, err := NewPDFDocument(pdfPath, nil)
	if err != nil {
		return nil, err
	}
	// Process the  HansardDoc ..
	hansardDoc := HansardDocument{StateAssemblySession: sessionName}
	cerr := NewHansardDocumentContent(pdfDoc, &hansardDoc)
	if cerr != nil {
		return nil, cerr
	}
	// Any post processing??
	return &hansardDoc, nil
}

func detectHansardType(firstPage PDFPage) (HansardType, error) {
	// default value ..

	// Is this fatal?
	return -1, fmt.Errorf("could not detect a type")
}

func NewHansardDocumentContent(pdfDoc *PDFDocument, hansardDoc *HansardDocument) error {
	// validation checks
	if len(pdfDoc.Pages) < 1 {
		return fmt.Errorf("Needs at least one page for valid Hansard!!! Found: %d", len(pdfDoc.Pages))
	}
	// Detect HansardType from the first page ... and fill it up ..
	hansardType, derr := detectHansardType(pdfDoc.Pages[0])
	if derr != nil {
		return derr
	}
	hansardDoc.HansardType = hansardType
	// Extract out Questions metadata for all pages ..
	hansardQuestions := make([]HansardQuestion, 0, 20)
	hqerr := NewHansardQuestions(pdfDoc, &hansardQuestions)
	if hqerr != nil {
		return hqerr
	}
	// Fill up the questions ...
	hansardDoc.HansardQuestions = hansardQuestions
	// All OK?
	return nil
}

func isStartOfQuestionSection(rowContents []string) bool {

	// Look out for pertanyaan pattern

	// Look for question number pattern

	// Look for topic

	// Look for who ask for it ..

	return false
}

func extractQuestionNum(rowContent string) (string, error) {
	foundQuestionNum := ""
	// pattern to look out for is
	// <digit>.  Bertanya kepada ...
	return foundQuestionNum, nil
}

func NewHansardQuestions(pdfDoc *PDFDocument, hansardQuestions *[]HansardQuestion) error {
	if pdfDoc == nil {
		return fmt.Errorf("pdfDoc is nil!")
	}
	// Iterate through all pages
	for _, r := range pdfDoc.Pages {
		// Init a new hansardQuestion struct
		hansardQuestion := HansardQuestion{PageNumStart: r.PageNo}
		if isStartOfQuestionSection(r.PDFTxtSameLines) {
			for _, rowContent := range r.PDFTxtSameLines {
				foundQuestionNum, exerr := extractQuestionNum(rowContent)
				// DEBUG ..
				fmt.Println(fmt.Sprintf("FOUND Question %s in page %d", foundQuestionNum, r.PageNo))
				if exerr != nil {
					return exerr
				}
				if foundQuestionNum != "" {
					// fill it in, the foundQuestionNum; need to strip?
					hansardQuestion.QuestionNum = strings.TrimSpace(foundQuestionNum)
					*hansardQuestions = append(*hansardQuestions, hansardQuestion)
					break
				}
			}
		}
		// Update end page num as we go along ..
		hansardQuestion.PageNumEnd = r.PageNo
	}
	return nil
}
