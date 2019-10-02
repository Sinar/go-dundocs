package hansard_test

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"gopkg.in/yaml.v2"

	"github.com/Sinar/go-dundocs/internal/hansard"
)

var updateHansard = flag.Bool("updateHansard", false, "update Hansard .golden files")
var updateHansardPDF = flag.Bool("updateHansardPDF", false, "update Hansard .fixture PDFs")

func TestNewHansardQuestions(t *testing.T) {
	type args struct {
		fixtureLabel     string
		pdfPath          string
		hansardQuestions *[]hansard.HansardQuestion
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"case #1", args{"HDOC-Lisan-1-20", "raw/Lisan/SOALAN MULUT (1-20).pdf", &[]hansard.HansardQuestion{{"1", 1, 4}, {"2", 5, 5}}}, false},
		//{"case #2", args{"HDOC-Lisan-1-20", "raw/Lisan/SOALAN MULUT (1-20).pdf", nil}, false},
		{"case #3", args{"HDOC-BukanLisan-1-20", "raw/BukanLisan/1 - 20.pdf", &[]hansard.HansardQuestion{{"1", 1, 1}, {"2", 2, 2}, {"3", 3, 5}}}, false},
		//{"case #4", args{"HDOC-BukanLisan-1-20", "raw/BukanLisan/1 - 20.pdf", nil}, false},
		{"sad #1", args{"HDOC-Lisan-1-20", "raw/Lisan/SOALAN MULUT (1-20).pdfa", nil}, true},
		{"sad #2", args{"HDOC-BukanLisan-1-20", "raw/BukanLisan/1 - 20b.pdf", nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pdfDoc := samplePDFFromFixture(t, tt.args.fixtureLabel, tt.args.pdfPath)
			hansardQuestions := make([]hansard.HansardQuestion, 0, 20)
			if err := hansard.NewHansardQuestions(pdfDoc, &hansardQuestions); (err != nil) != tt.wantErr {
				t.Errorf("NewHansardQuestions() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Show diff if any ..
			if diff := cmp.Diff(tt.args.hansardQuestions, &hansardQuestions); diff != "" {
				t.Errorf("hansardQuestions mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestNewHansardDocument(t *testing.T) {
	type args struct {
		fixtureLabel string
		pdfPath      string
	}
	tests := []struct {
		name    string
		args    args
		want    *hansard.HansardDocument
		wantErr bool
	}{
		{"happy #1", args{"HDOC-Lisan-1-20", "raw/Lisan/SOALAN MULUT (1-20).pdf"}, nil, false},
		{"happy #2", args{"HDOC-BukanLisan-1-20", "raw/BukanLisan/1 - 20.pdf"}, nil, false},
		//{"happy #1", args{"", ""}, nil, false},
		//{"happy #1", args{"", ""}, nil, false},
		//{"happy #1", args{"", ""}, nil, false},
		//{"happy #1", args{"", ""}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hansardDocument := hansard.HansardDocument{StateAssemblySession: "testSessionName"}
			pdfDoc := loadPDFFromFixture(t, tt.args.fixtureLabel, tt.args.pdfPath)
			got := hansard.NewHansardDocumentContent(pdfDoc, &hansardDocument)
			//got, err := hansard.NewHansardDocument("sessionName", tt.args.pdfPath)
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("NewHansardDocument() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("NewHansardDocumentContent() mismatch (-want +got):\n%s", diff)
			}

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewHansardDocument() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

// Helper function to sample 5 pages from fixture
func samplePDFFromFixture(t *testing.T, fixtureLabel string, sourcePath string) *hansard.PDFDocument {
	// Mark as helper
	t.Helper()

	var pdfDoc *hansard.PDFDocument
	// Read from cache; if not exist; complain that need to update
	fixture := filepath.Join("testdata", "sample-"+fixtureLabel+".fixture")
	if *updateHansardPDF {
		// If run update; call the same function used by TOC to get the data
		pdfDoc, err := hansard.NewPDFDocument(sourcePath, &hansard.ExtractPDFOptions{NumPages: 5})
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
