package hansard

import (
	"testing"

	"github.com/ledongthuc/pdf"
)

func Test_extractContentPreviewByRow(t *testing.T) {
	type args struct {
		row            *pdf.Row
		contentPreview *[]string
		i              int
	}
	// Get a few page from common PDF ..
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := extractContentPreviewByRow(tt.args.row, tt.args.contentPreview, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("extractContentPreviewByRow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_extractPDFPageContent(t *testing.T) {
	type args struct {
		pdfPage *PDFPage
		p       pdf.Page
	}
	// We have a test PDF file we will use in all the tests
	// or maybe a few other alternatives ..
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := extractPDFPageContent(tt.args.pdfPage, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("extractPDFPageContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
