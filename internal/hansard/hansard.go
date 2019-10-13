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

type ErrorQuestionsHasInvalid struct {
	badQuestionsCount int
}

func (e *ErrorQuestionsHasInvalid) Error() string {
	return fmt.Sprintf("Has %d bad Questions!", e.badQuestionsCount)
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

	// Look out for pertanyaan pattern .. is duplicated
	// TODO: Should sync it up somehow; having it embedded so deep here might
	// not be too good ..
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
	// Setup regexp once
	re := regexp.MustCompile(`(?i)^.*?(\d+).*bertanya\s+kepada.*$`)
	// TODO: Might have to break up line cases; what other special characters will appear here?
	// It fails with the content; but will it appear in real life? ==> "\n\n\n 50 bertanya kepada yab menteri besar Azmin ALI "
	sm := re.FindStringSubmatch(rowContent)
	if sm != nil {
		// DEBUG:
		//fmt.Println("FOUND NUM: ", sm[1])
		return sm[1], nil
	}
	return "", nil
}

func NewHansardQuestions(pdfDoc *PDFDocument, hansardQuestions *[]HansardQuestion) error {
	if pdfDoc == nil {
		return fmt.Errorf("pdfDoc is nil!")
	}
	var hansardQuestion *HansardQuestion
	var badQuestionsCount int
	// Iterate through all pages
	for _, r := range pdfDoc.Pages {
		if isStartOfQuestionSection(r) {
			// Special case: first round ..
			if hansardQuestion != nil {
				// Before  append; let's check previous and flag if got bad  question
				if hansardQuestion.QuestionNum == "0" {
					badQuestionsCount++
				}
				//  Otherwise, attach as per needed ..
				*hansardQuestions = append(*hansardQuestions, *hansardQuestion)
			}
			// Init a new hansardQuestion struct
			hansardQuestion = &HansardQuestion{QuestionNum: "0", PageNumStart: r.PageNo}
			for _, rowContent := range r.PDFTxtSameLines {
				foundQuestionNum, exerr := extractQuestionNum(rowContent)
				// DEBUG ..
				//fmt.Println(fmt.Sprintf("FOUND Question %s in page %d", foundQuestionNum, r.PageNo))
				if exerr != nil {
					return exerr
				}
				if foundQuestionNum != "" {
					// fill it in, the foundQuestionNum; need to strip?
					hansardQuestion.QuestionNum = strings.TrimSpace(foundQuestionNum)
					break
				}
			}
		}
		// Put some protection checks ..
		if hansardQuestion != nil {
			// Update end page num as we go along ..
			hansardQuestion.PageNumEnd = r.PageNo
		}
	}
	// Special case: last line; code below probably can be refactored!
	if hansardQuestion != nil {
		// Before  append; let's check previous and flag if got bad  question
		if hansardQuestion.QuestionNum == "0" {
			badQuestionsCount++
		}
		//  Otherwise, attach as per needed ..
		*hansardQuestions = append(*hansardQuestions, *hansardQuestion)
	}
	// If have badQuestionsCount; flag it; NOT fatal; but to be handled by caller
	if badQuestionsCount > 0 {
		return fmt.Errorf("NewHansardQuestions FAIL: %w", &ErrorQuestionsHasInvalid{badQuestionsCount})
	}
	// Reached here; all OK and peachy!
	return nil
}
