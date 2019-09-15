package hansard

import (
	"reflect"
	"testing"
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
