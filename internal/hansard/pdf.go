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
	sourcePath string
}

const (
	MaxLineProcessed = 10
)

func NewPDFDocument(pdfPath string) (*PDFDocument, error) {
	f, r, err := pdf.Open(pdfPath)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, fmt.Errorf("PDFAccessErr: %w", err)
	}
	totalPage := r.NumPage()

	// Use it in another function ...
	//for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
	//	p := r.Page(pageIndex)
	//	if p.V.IsNull() {
	//		continue
	//	}
	//
	//	rows, _ := p.GetTextByRow()
	//	for _, row := range rows {
	//		println(">>>> row: ", row.Position)
	//		for _, word := range row.Content {
	//			fmt.Println(word.S)
	//		}
	//	}
	//}

	// Init it and fill it with the extracted info  earlier ..
	pdfDoc := PDFDocument{
		NumPages: totalPage,
	}

	return &pdfDoc, nil
}
