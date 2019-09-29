package hansard_test

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/Sinar/go-dundocs/internal/hansard"
)

var updateHansard = flag.Bool("update", false, "update Hansard .golden files")
var updateHansardPDF = flag.Bool("updatePDF", false, "update Hansard .fixture PDFs")

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
		{"case #1", args{nil, nil}, false},
		{"case #2", args{nil, nil}, false},
		{"case #3", args{nil, nil}, false},
		{"case #4", args{nil, nil}, false},
		{"case #5", args{nil, nil}, false},
		{"case #6", args{nil, nil}, false},
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
		{"happy #1", args{"", ""}, nil, false},
		{"happy #1", args{"", ""}, nil, false},
		{"happy #1", args{"", ""}, nil, false},
		{"happy #1", args{"", ""}, nil, false},
		{"happy #1", args{"", ""}, nil, false},
		{"happy #1", args{"", ""}, nil, false},
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

// Helper function to load from fixture; safe/update as per necessary?
func loadPDFFromFixture(t *testing.T, fixtureLabel string, sourcePath string) *hansard.PDFDocument {
	// Mark as helper
	t.Helper()

	var pdfDoc *hansard.PDFDocument
	// Read from cache; if not exist; complain that need to update
	fixture := filepath.Join("testdata", fixtureLabel+".fixture")
	if *updateHansardPDF {
		// If run update; call the same function used by TOC to get the data
		pdfDoc, err := hansard.NewPDFDocument(sourcePath, nil)
		if err != nil {
			t.Fatalf("NewPDFDocument FAIL: %s", err.Error())
		}
		// Persist the data into the file
		w, werr := yaml.Marshal(pdfDoc)
		if werr != nil {
			t.Fatalf("Marshal FAIL: %s", werr.Error())
		}
		ioutil.WriteFile(fixture, w, 0644)
		return pdfDoc
	}
	// Normal path,read from fixture ,..
	want, rerr := ioutil.ReadFile(fixture)
	if rerr != nil {
		// Cannot proceed with one golden file update
		if os.IsNotExist(rerr) {
			t.Fatalf("Ensure run with -updateHansardPDF flag first time! ERR: %s", rerr.Error())
		}
		t.Fatalf("Unexpected error: %s", rerr.Error())
	}
	pdfDoc = &hansard.PDFDocument{}
	umerr := yaml.Unmarshal(want, pdfDoc)
	if umerr != nil {
		t.Fatalf("Unmarshal FAIL: %s", umerr.Error())
	}
	return pdfDoc
}
