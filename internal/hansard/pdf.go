package hansard

import (
	"fmt"

	"github.com/ledongthuc/pdf"
)

type PDFPage struct {
	PageNo           int
	PDFPlainText     string
	PDFTxtSameLines  []string // combined content with same line .. proxy for changes
	PDFTxtSameStyles []string // combined content with same style .. proxy for changes
}

type PDFDocument struct {
	NumPages   int
	Pages      []PDFPage
	SourcePath string
}

type ExtractPDFOptions struct {
	StartPage int
	NumPages  int
}

const (
	MaxLineProcessed = 10
)

func NewPDFDocument(pdfPath string, options *ExtractPDFOptions) (*PDFDocument, error) {
	f, r, err := pdf.Open(pdfPath)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, fmt.Errorf("PDFAccessErr: %w", err)
	}

	startPage := 1
	totalPage := r.NumPage()
	// DEBUG
	//totalPage = 3
	if options != nil {
		if options.StartPage > 1 {
			startPage = options.StartPage
		}
		if options.NumPages > 0 {
			totalPage = options.NumPages
		}
	}
	var pdfPages []PDFPage
	// Init it and fill it with the extracted info  earlier ..
	pdfDoc := PDFDocument{
		NumPages:   totalPage,
		Pages:      pdfPages,
		SourcePath: pdfPath,
	}

	for pageIndex := startPage; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		// New Page to be Processed ..
		newPageProcessed := PDFPage{
			PageNo:           pageIndex,
			PDFTxtSameLines:  []string{},
			PDFTxtSameStyles: []string{},
		}
		exerr := extractPDFPageContent(&newPageProcessed, p)
		if exerr != nil {
			panic(exerr)
		}
		// If OK, append them ..
		pdfDoc.Pages = append(pdfDoc.Pages, newPageProcessed)
	}

	return &pdfDoc, nil
}

func extractPDFPageContent(pdfPage *PDFPage, p pdf.Page) error {

	rows, _ := p.GetTextByRow()
	contentPreview := make([]string, 0, 20)
	for i, row := range rows {
		// DEBUG
		//fmt.Println("LINE: ", i)
		exerr := extractContentPreviewByRow(row, &contentPreview)
		if exerr != nil {
			// HOWTO multi error?
			//return &pdfDoc, fmt.Errorf(": %w", exerr)
			panic(exerr)
		}
		if i > MaxLineProcessed {
			//fmt.Println("Processed enough!! break!!")
			break
		}
	}
	pdfPage.PDFTxtSameLines = contentPreview

	return nil
}

func extractContentPreviewByRow(row *pdf.Row, contentPreview *[]string) error {
	// Use it in another function ...
	var currentContentLine string
	//		println(">>>> row: ", row.Position)
	for _, word := range row.Content {
		// DEBUG
		//fmt.Println(word.S)
		// Is trabnsformation needed?
		currentContentLine += word.S
	}
	//	}

	//fmt.Println(currentContentLine)
	*contentPreview = append(*contentPreview, currentContentLine)
	return nil
}
