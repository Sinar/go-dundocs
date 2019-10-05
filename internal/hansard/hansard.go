package hansard

import (
	"fmt"
	"regexp"
	"strings"
)

type HansardType int

const (
	HANSARD_INVALID HansardType = iota
	HANSARD_SPOKEN
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
	for _, rowContent := range firstPage.PDFTxtSameLines {
		normalizedContent := strings.ToLower(rowContent)
		// Look  out for pertanyaan
		hasQuestion, err := regexp.MatchString("pertanyaan", normalizedContent)
		if err != nil {
			return HANSARD_INVALID, err
		}
		if hasQuestion {
			// Has potential; do the further checks ..
			// Look out for mulut
			hasSpokenHansardType, serr := regexp.MatchString("mulut", normalizedContent)
			if serr != nil {
				return HANSARD_INVALID, serr
			}
			// If found match; get out IMMEDIATELY!
			if hasSpokenHansardType {
				return HANSARD_SPOKEN, nil
			}
			// Look out for tulis
			hasWrittenHansardType, werr := regexp.MatchString("tulis", normalizedContent)
			if werr != nil {
				return HANSARD_INVALID, werr
			}
			if hasWrittenHansardType {
				return HANSARD_WRITTEN, nil
			}
		}
	}
	// If get here without a match, no type FOUND! INVALID default ..
	return HANSARD_INVALID, nil
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

func isStartOfQuestionSection(currentPage PDFPage) bool {

	// Look out for pertanyaan pattern
	hansardType, err := detectHansardType(currentPage)
	if err != nil {
		panic(err)
	}
	if hansardType == HANSARD_INVALID {
		return false
	}
	// More sophisticated checks later ?? If ever ..
	// At this point we know if it is SPOKEN, WRITTEn etc?
	// Look for question number pattern

	// Look for topic

	// Look for who ask for it ..

	return true
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
		if isStartOfQuestionSection(r) {
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
