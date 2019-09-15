package hansard

import (
	"reflect"
	"testing"

	"github.com/ledongthuc/pdf"
)

func TestNewPDFDocument(t *testing.T) {
	type args struct {
		pdfPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *PDFDocument
		wantErr bool
	}{
		{"missing file, xerrors", args{"testdata/bogus.pdf"}, nil, true},
		{"happy path, Lisan", args{"/Users/leow/GOMOD/go-dundocs/raw/Lisan/SOALAN MULUT (261-272).pdf"}, &PDFDocument{NumPages: 18}, false},
		{"happy path, BukanLisan", args{"/Users/leow/GOMOD/go-dundocs/raw/BukanLisan/41 - 60.pdf"}, &PDFDocument{NumPages: 39}, false},
		{"table files, lampiran", args{"/Users/leow/GOMOD/go-dundocs/raw/Lampiran/LAMPIRAN B JAWAPAN SOALAN NO.213.pdf"}, &PDFDocument{NumPages: 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPDFDocument(tt.args.pdfPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPDFDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				if !reflect.DeepEqual(got.NumPages, tt.want.NumPages) {
					t.Errorf("NewPDFDocument() got = %v, want %v", got, tt.want.NumPages)
				}
			}
		})
	}
}

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
