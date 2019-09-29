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
	hansardDoc := HansardDocument{StateAssemblySession: sessionName}
	// Load the PDFDoc; should we check length?
	pdfDoc, err := NewPDFDocument(pdfPath, nil)
	if err != nil {
		return nil, err
	}
	// Detect HansardType from the first page ... and fill it up ..
	if len(pdfDoc.Pages) < 1 {
		return nil, fmt.Errorf("Needs to have at least one page to be considered a valid Hansard!!!", len(pdfDoc.Pages))
	}
	hansardDoc.HansardType = detectHansardType(pdfDoc.Pages[0])
	hansardQuestions := make([]HansardQuestion, 0, 20)
	hqerr := NewHansardQuestions(pdfDoc, &hansardQuestions)
	if hqerr != nil {
		return nil, hqerr
	}
	// Fill up the questions ...
	hansardDoc.HansardQuestions = hansardQuestions

	return &hansardDoc, nil
}

func detectHansardType(firstPage PDFPage) HansardType {

	// default value ..
	return HANSARD_SPOKEN
}

func NewHansardQuestions(pdfDoc *PDFDocument, hansardQuestions *[]HansardQuestion) error {
	// Iterate through all pages
	for _, r := range pdfDoc.Pages {
		// Init a new hansardQuestion struct
		hansardQuestion := HansardQuestion{PageNumStart: r.PageNo}
		for _, rowContent := range r.PDFTxtSameLines {
			if isStartOfQuestionSection(rowContent) {
				foundQuestionNum, exerr := extractQuestionNum(rowContent)
				fmt.Println("FOUND Question %s in page %d", foundQuestionNum, r.PageNo)
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
			// Update end page num as we go along ..
			hansardQuestion.PageNumEnd = r.PageNo
		}
	}
	return nil
}

func isStartOfQuestionSection(rowContent string) bool {

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
