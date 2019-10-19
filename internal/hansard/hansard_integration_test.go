package hansard_test

import (
	"errors"
	"flag"
	"fmt"
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
		{"case #1", args{
			"HDOC-Lisan-1-20", "raw/Lisan/SOALAN MULUT (1-20).pdf",
			&[]hansard.HansardQuestion{
				{"1", 1, 4},
				{"2", 5, 5},
			},
		}, false},
		{"case #3", args{
			"HDOC-BukanLisan-1-20", "raw/BukanLisan/1 - 20.pdf",
			&[]hansard.HansardQuestion{
				{"1", 1, 1},
				{"2", 2, 2},
				{"3", 3, 5},
			},
		}, false},
		{"sad #1", args{
			"Bad-HDOC-Lisan-1-20", "raw/Lisan/SOALAN MULUT (1-20).pdfa",
			&[]hansard.HansardQuestion{
				{"8", 18, 18}, // This simulates scanned pages messed up ordering
				{"9", 19, 11}, // This is to simulate messing the ordering in PDFs
				{"5", 12, 13},
				{"6", 14, 15},
				{"10", 20, 21},
				{"10", 22, 23},
				{"15", 30, 36},
				{"19", 37, 37},
				{"0", 38, 39},
			},
		}, true},
		// We expect QuestionNum of "0" for pages with some marker but could NOT recognize question number!
		{"sad #2", args{
			"Bad-HDOC-BukanLisan-1-20", "raw/BukanLisan/1 - 20b.pdf",
			&[]hansard.HansardQuestion{
				{"0", 11, 12},
				{"0", 29, 30},
				{"0", 32, 33},
			},
		}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pdfDoc := samplePDFFromFixture(t, tt.args.fixtureLabel, tt.args.pdfPath)
			hansardQuestions := make([]hansard.HansardQuestion, 0, 20)
			// Run function  ..
			err := hansard.NewHansardQuestions(pdfDoc, &hansardQuestions)
			// Check if expectation is fulfilled
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHansardQuestions() error = %v, wantErr %v", err, tt.wantErr)
			}
			//  For errors; check out Error Type to see if it is  recoverable
			if err != nil {
				//  Old way is
				// if err == hansard.ErrorQuestionsHasInvalid ..
				// Below does not work; apparently is not expression -> hansard.ErrorQuestionsHasInvalid. Why?
				//if errors.Is(err, hansard.ErrorQuestionsHasInvalid) {
				//	t.Errorf("ERR: %v", err)
				//}
				// Below would be the old way to unwrap the embedded error and check the type
				errQInvalid, ok := errors.Unwrap(err).(*hansard.ErrorQuestionsHasInvalid)
				if ok {
					fmt.Println("RECOVERABLE_OLD: ", errQInvalid.Error())
				} else {
					// Is more serious error?
				}
				// Below is the equivalent of above;  but better as it unwraps all and try to match ..
				var invalidQError *hansard.ErrorQuestionsHasInvalid
				if errors.As(err, &invalidQError) {
					fmt.Println("RECOVERABLE_NEW: ", invalidQError.Error())
				} else {
					// Is more serioous error; cannot be recovered!
					panic(err)
					//fmt.Println(err)
				}
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
		want    hansard.HansardDocument
		wantErr bool
	}{
		{"happy #1", args{"HDOC-Lisan-1-20", "raw/Lisan/SOALAN MULUT (1-20).pdf"}, hansard.HansardDocument{
			StateAssemblySession: "testSessionName",
			HansardType:          hansard.HANSARD_SPOKEN,
			HansardQuestions: []hansard.HansardQuestion{
				{"1", 1, 4},
				{"2", 5, 8},
				{"3", 9, 9},
				{"4", 10, 11},
				{"5", 12, 13},
				{"6", 14, 15},
				{"7", 16, 17},
				{"8", 18, 18},
				{"9", 19, 19},
				{"10", 20, 21},
				{"11", 22, 23},
				{"12", 24, 25},
				{"13", 26, 28},
				{"14", 29, 29},
				{"15", 30, 31},
				{"16", 32, 32},
				{"17", 33, 34},
				{"18", 35, 36},
				{"19", 37, 37},
				{"20", 38, 39},
			},
		}, false},
		{"happy #2", args{"HDOC-BukanLisan-1-20", "raw/BukanLisan/1 - 20.pdf"}, hansard.HansardDocument{
			StateAssemblySession: "testSessionName",
			HansardType:          hansard.HANSARD_WRITTEN,
			HansardQuestions: []hansard.HansardQuestion{
				{"1", 1, 1},
				{"2", 2, 2},
				{"3", 3, 5},
				{"4", 6, 6},
				{"5", 7, 7},
				{"6", 8, 9},
				{"7", 10, 10},
				{"8", 11, 17},
				{"9", 18, 18},
				{"10", 19, 19},
				{"11", 20, 20},
				{"12", 21, 21},
				{"13", 22, 24},
				{"14", 25, 26},
				{"15", 27, 28},
				{"16", 29, 30},
				{"17", 31, 31},
				{"18", 32, 34},
				{"19", 35, 36},
				{"20", 37, 37},
			},
		}, false},
		//{"happy #1", args{"", ""}, nil, false},
		//{"happy #1", args{"", ""}, nil, false},
		//{"happy #1", args{"", ""}, nil, false},
		//{"happy #1", args{"", ""}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load Fixture for test
			pdfDoc := loadPDFFromFixture(t, tt.args.fixtureLabel, tt.args.pdfPath)
			// DEBUG
			//spew.Dump(pdfDoc)
			got := hansard.HansardDocument{StateAssemblySession: "testSessionName"}
			err := hansard.NewHansardDocumentContent(pdfDoc, &got)
			// DEBUG
			//spew.Dump(got)
			// Catch bad pdf raw data .. no cases yet ..
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHansardDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("NewHansardDocumentContent() mismatch (-want +got):\n%s", diff)
				// DEBUG diff using alternative method
				//litter.Dump(tt.want)
				//litter.Dump(got)
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
