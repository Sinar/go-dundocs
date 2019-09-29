package hansard_test

import (
	"reflect"
	"testing"

	"github.com/Sinar/go-dundocs/internal/hansard"
)

func TestNewHansardQuestions(t *testing.T) {
	type args struct {
		pdfDoc           *hansard.PDFDocument
		hansardQuestions *[]hansard.HansardQuestion
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := hansard.NewHansardQuestions(tt.args.pdfDoc, tt.args.hansardQuestions); (err != nil) != tt.wantErr {
				t.Errorf("NewHansardQuestions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewHansardDocument(t *testing.T) {
	type args struct {
		sessionName string
		pdfPath     string
	}
	tests := []struct {
		name    string
		args    args
		want    *hansard.HansardDocument
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hansard.NewHansardDocument(tt.args.sessionName, tt.args.pdfPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHansardDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHansardDocument() got = %v, want %v", got, tt.want)
			}
		})
	}
}
