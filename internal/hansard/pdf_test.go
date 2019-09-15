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
		{"missing file, xerrors", args{""}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPDFDocument(tt.args.pdfPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPDFDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.NumPages, tt.want.NumPages) {
				t.Errorf("NewPDFDocument() got = %v, want %v", got, tt.want.NumPages)
			}
		})
	}
}
