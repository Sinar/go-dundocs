package hansard_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"

	"gopkg.in/yaml.v2"

	"github.com/Sinar/go-dundocs/internal/hansard"
)

var updateSplitterGolden = flag.Bool("updateSplitterGolden", false, "update SplitHansardDocPlan .golden files")

//var updateSplitterFixture = flag.Bool("updateSplitterFixture", false, "update SplitHansardDocPlan .fixture plans")

var testDir = ""

func TestMain(m *testing.M) {
	log.Println("Do stuff BEFORE the tests!")
	// Dynamic value; if in internal, chdir into hansard .. looks dirty :(
	// **** START UGLY *******
	//_, filename, _, _ := runtime.Caller(0)
	//wddir := path.Join(path.Dir(filename), "..")
	wddir, _ := os.Getwd()
	_, curDir := filepath.Split(wddir)
	//fmt.Println("CUR_DIR:", curDir)
	if curDir == "internal" {
		cderr := os.Chdir(filepath.Join(wddir, "hansard"))
		if cderr != nil {
			panic(cderr)
		}
	}
	testDir, _ = os.Getwd()
	// **** END UGLY *******
	exitVal := m.Run()
	log.Println("Do stuff AFTER the tests!")
	os.Exit(exitVal)
}

func TestNewSplitHansardDocumentPlan(t *testing.T) {
	type args struct {
		fixtureLabel string
	}
	tests := []struct {
		name    string
		args    args
		want    *hansard.SplitHansardDocumentPlan
		wantErr bool
	}{
		{"happy #1", args{"HDOC-Lisan-1-20"}, nil, false},
		{"happy #2", args{"HDOC-BukanLisan-1-20"}, nil, false},
		{"sad #1", args{"Bad-HDOC-Lisan-1-20"}, nil, true},
		{"sad #2", args{"Bad-HDOC-BukanLisan-1-20"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prep temp dir for usage ..
			// Prepare TempDir for working with it
			dir, err := ioutil.TempDir("", "dundocs-splitter")
			if err != nil {
				log.Fatal(err)
			}
			// Comment out below if need to see the output in dir
			defer os.RemoveAll(dir)
			// DEBUG
			//log.Println("Dir is ", dir)
			// Load sample PDFs and check against the plan passed in
			// use the sample version so test cases not too long .. WILL be loaded from testdata fixture
			// Prefix with  plan-sample as it is a PLAN made from SAMPLE of the Fixture in testdata already
			// should FAIL if not found!
			pdfDoc := samplePDFFromFixture(t, tt.args.fixtureLabel, "")
			// Init
			confDUNSession := "confDUNSession"
			got := hansard.SplitHansardDocumentPlan{
				HansardDocument: hansard.HansardDocument{
					StateAssemblySession: confDUNSession,
					HansardQuestions:     []hansard.HansardQuestion{},
				},
			}
			// QUESTION: What to do if failed load??
			// Sav eplan  fixture here? or just check structure
			//if got := hansard.NewSplitHansardDocumentPlanContent(pdfDoc, &splitPlan); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewSplitHansardDocumentPlan() = %v, want %v", got, tt.want)
			//}
			serr := hansard.NewSplitHansardDocumentPlanContent(pdfDoc, &got)
			if (serr != nil) != tt.wantErr {
				t.Errorf("NewSplitHansardDocumentPlanContent() error = %v, wantErr %v", serr, tt.wantErr)
				return
			}
			// want load from goldenImage
			// Save the PLan for use by LoadPLan
			goldenLabel := "plan-sample-" + tt.args.fixtureLabel
			want := hansard.SplitHansardDocumentPlan{
				HansardDocument: hansard.HansardDocument{
					StateAssemblySession: confDUNSession,
					HansardQuestions:     []hansard.HansardQuestion{},
				},
			}
			if *updateSplitterGolden {
				loadPlanFromGolden(t, goldenLabel, &got.HansardDocument)
				want.HansardDocument = got.HansardDocument
			} else {
				// load from cache ..
				var wantHansardDocument hansard.HansardDocument
				loadPlanFromGolden(t, goldenLabel, &wantHansardDocument)
				want.HansardDocument = wantHansardDocument
			}
			// Finish updating want .. need to ignore unexprted ..
			if diff := cmp.Diff(want, got, cmpopts.IgnoreUnexported(hansard.SplitHansardDocumentPlan{})); diff != "" {
				t.Errorf("NewSplitHansardDocumentPlanContent() mismatch (-want +got):\n%s", diff)
				// DEBUG diff using alternative method
				//litter.Dump(tt.want)
				//litter.Dump(got)
			}
		})
	}
}

func TestSplitHansardDocumentPlan_SavePlan(t *testing.T) {
	type args struct {
		fixtureLabel string
	}
	tests := []struct {
		name            string
		args            args
		dataDir         string
		expectedPlanDir string
		wantErr         bool
		expectPlanFile  bool
	}{
		{"default data folder", args{"HDOC-Lisan-1-20"}, "", "data/HDOC-Lisan-1-20", false, true},
		{"custom data folder", args{"HDOC-BukanLisan-1-20"}, "custom/datadir", "custom/datadir/HDOC-BukanLisan-1-20", false, true},
		{"absolute custom data folder", args{"HDOC-BukanLisan-1-20"}, "/tmp/datadir", "/tmp/datadir/HDOC-BukanLisan-1-20", false, true}, // Will do this when have absolutePath helper function
		{"bad plan not saved", args{"Bad-HDOC-Lisan-1-20"}, "", "", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare TempDir for working with it
			dir, err := ioutil.TempDir("", "dundocs-splitter")
			if err != nil {
				log.Fatal(err)
			}
			// Comment out below if need to see the output in dir
			defer os.RemoveAll(dir)
			// DEBUG
			//log.Println("Dir is ", dir)
			// want load from goldenImage
			// Save the PLan for use by LoadPLan
			goldenLabel := "plan-sample-" + tt.args.fixtureLabel
			absoluteDataDir := hansard.GetAbsoluteDataDir(dir, tt.dataDir)
			absolutePlanDir := absoluteDataDir + "/" + tt.args.fixtureLabel
			s := hansard.NewEmptySplitHansardDocumentPlan(absoluteDataDir, absolutePlanDir, "confDUNSession")
			// DEBUG
			//spew.Dump(s)
			// load from Golden copy; assume past tests has been run ..
			loadPlanFromGolden(t, goldenLabel, &s.HansardDocument)
			// DEBUG
			//fmt.Println(litter.Sdump(s))
			// Run the save, catch invalid plans
			serr := s.SavePlan()
			if (serr != nil) != tt.wantErr {
				t.Errorf("SavePlan() error = %v, wantErr %v", serr, tt.wantErr)
				return
			}
			// If OK, do the rest of the checks
			if serr == nil {
				if serr := s.SavePlan(); (serr != nil) != tt.wantErr {
					t.Errorf("SavePlan() error = %v, wantErr %v", serr, tt.wantErr)
					return
				}
				// Check planFile actually created and exists ..
				if _, err := os.Stat(absolutePlanDir); os.IsNotExist(err) {
					t.Errorf("PLANDIR_MISSING: %s", err)
					return
				}
				// Check split file in known place
				if _, err := os.Stat(absolutePlanDir + "/split.yml"); os.IsNotExist(err) {
					// DEBUG
					//fmt.Println("EXPECT:" + expectedPlanDir + "/split.yml")
					// path/to/whatever does not exist
					t.Errorf("SPLITPLAN_MISSING: %s", err)
					return
				}
				// If plan file exist; check its content; is it necessary?
				contentYAML, rerr := ioutil.ReadFile(absolutePlanDir + "/split.yml")
				if rerr != nil {
					t.Fatalf("ERR: %s", rerr.Error())
				}
				var gotPlan hansard.HansardDocument
				umerr := yaml.Unmarshal(contentYAML, &gotPlan)
				if umerr != nil {
					panic(umerr)
				}
				if diff := cmp.Diff(s.HansardDocument, gotPlan); diff != "" {
					t.Errorf("Plan mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestSplitHansardDocumentPlan_LoadPlan(t *testing.T) {
	type args struct {
		fixtureLabel string
	}
	tests := []struct {
		name     string
		args     args
		planFile string
		wantErr  bool
	}{
		// NOTE: Not too nice; still tied to physical path that are not in testdata
		{"happy #1", args{"HDOC-Lisan-1-20"}, "testdata/plan-sample-HDOC-Lisan-1-20.golden", false},
		{"happy #2", args{"HDOC-BukanLisan-1-20"}, "testdata/plan-sample-HDOC-BukanLisan-1-20.golden", false},
		{"sad #1", args{"Bad-HDOC-Lisan-1-20"}, "testdata/plan-sample-Bad-HDOC-Lisan-1-20.golden", false},
		{"sad #2", args{"Bad-HDOC-BukanLisan-1-20"}, "testdata/plan-sample-Bad-HDOC-BukanLisan-1-20.golden", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup where the plan is
			absoluteDataDir, dderr := filepath.Abs(".")
			if dderr != nil {
				panic(dderr)
			}
			absolutePlanFile, plerr := filepath.Abs(tt.planFile)
			if plerr != nil {
				panic(plerr)
			}
			// DEBUG
			fmt.Println("DATA_PATH: ", absoluteDataDir, " PLAN_PATH: ", absolutePlanFile)
			s := hansard.NewEmptySplitHansardDocumentPlan(absoluteDataDir, absolutePlanFile, "dun15sesi3")
			// Load it into the struct
			if err := s.LoadPlan(); (err != nil) != tt.wantErr {
				t.Errorf("LoadPlan() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := &s.HansardDocument
			// Now compare against the Plan related to the PDF Fixture (from Hansard Integration) we test against
			// which will be persisted to Golden plans
			// We pass "" for pdfPath as it will load from testdata; failing it should FAIL LOUDLY
			want := loadPlanSampleFixture(t, tt.args.fixtureLabel, "")

			// Show diff if any ..
			// Compare against the expected  output? Which is one of the previous HansardDoc files ..
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("LoadPlan mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSplitHansardDocumentPlan_ExecuteSplit(t *testing.T) {
	type fields struct {
		srcPDFPath      string
		planDir         string
		hansardDocument hansard.HansardDocument
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		wantedFiles []string
	}{
		{"no plan #1", fields{"testdata/BukanLisan_41_60_36_39.pdf", "testdata/happy1-plan-execute", hansard.HansardDocument{}}, true, []string{}},
		{"happy #1", fields{"testdata/BukanLisan_41_60_36_39.pdf", "testdata/happy1-plan-execute", hansard.HansardDocument{
			StateAssemblySession: "dun15sesi3",
			HansardType:          hansard.HANSARD_WRITTEN,
			HansardQuestions: []hansard.HansardQuestion{
				{"58", 1, 1},
				{"59", 2, 3},
				{"60", 4, 4},
			},
		}}, false, []string{
			"BukanLisan_41_60_36_39/dun15sesi3-soalan-58.pdf",
			"BukanLisan_41_60_36_39/dun15sesi3-soalan-59.pdf",
			"BukanLisan_41_60_36_39/dun15sesi3-soalan-60.pdf",
		}},
		{"happy #2", fields{"testdata/Lisan_Mulut_261_272.pdf", "testdata/happy2-plan-execute", hansard.HansardDocument{
			StateAssemblySession: "dun15sesi3",
			HansardType:          hansard.HANSARD_SPOKEN,
			HansardQuestions: []hansard.HansardQuestion{
				{"269", 1, 1},
				{"270", 2, 4},
				{"271", 5, 5},
			},
		}}, false, []string{
			"Lisan_Mulut_261_272/dun15sesi3-soalan-269.pdf",
			"Lisan_Mulut_261_272/dun15sesi3-soalan-270.pdf",
			"Lisan_Mulut_261_272/dun15sesi3-soalan-271.pdf",
		}},
	}
	// Prepare the PDF test cases we will use in the tests ..
	// Same as used in  hansard_integration_tests; we will run them  here to set it up
	// to be copied and  used again and again?

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load up the SplitPlan aka HansardDocument
			// Prepare TempDir for working with it
			dir, err := ioutil.TempDir("", "dundocs-splitter")
			if err != nil {
				log.Fatal(err)
			}
			// Comment out below if need to see the output in dir
			defer os.RemoveAll(dir)
			// DEBUG
			log.Println("Dir is ", dir)
			// Calculate the absolute versions
			absoluteDataDir, dderr := filepath.Abs(dir)
			if dderr != nil {
				panic(dderr)
			}
			absolutePlanDir, plerr := filepath.Abs(tt.fields.planDir)
			if plerr != nil {
				panic(plerr)
			}
			absoluteSrcPDF, sperr := filepath.Abs(tt.fields.srcPDFPath)
			if sperr != nil {
				panic(sperr)
			}
			// DEBUG
			//fmt.Println("DATA_PATH: ", absoluteDataDir, " PLAN_PATH: ", absolutePlanDir)
			s := hansard.NewEmptySplitHansardDocumentPlan(absoluteDataDir, absolutePlanDir, "")
			// Plan loaded; assume it is extracted from split.yml
			s.HansardDocument = tt.fields.hansardDocument
			// DEBUG
			//spew.Dump(s)
			// Default splitout location; will be created if needed ..
			absoluteSplitOutput := dir + "/splitout"
			// Execute the actual split; with one test PDF?
			exerr := s.ExecuteSplit(absoluteSrcPDF, absoluteSplitOutput)
			if (exerr != nil) != tt.wantErr {
				t.Errorf("ExecuteSplit() error = %v, wantErr %v", exerr, tt.wantErr)
			}

			if exerr != nil {
				// Expected; error, just leave
				return
			}

			for _, fileName := range tt.wantedFiles {
				// DEBUG
				//fmt.Println("Check split file at: ", absoluteSplitOutput+"/"+fileName)
				_, rerr := ioutil.ReadFile(filepath.Join(absoluteSplitOutput, fileName))
				if rerr != nil {
					t.Errorf("cannot read: %s", rerr)
				}
				// Open the output in the WorkDir and see if it has the expected number of files
				// Pattern in dataDir ==> base_filename_question_<nn>
				// Should we check content? probably no need .. what about filename; is it important; no for now ..
			}
		})
	}
}

// Helper functions for tests

// Loader Plan for use in Load phase (Full)
func loadPlanFromFixture(t *testing.T, fixtureLabel string, pdfPath string) *hansard.HansardDocument {
	// TODO: Maybe if needed in the future
	t.Helper()
	return nil
}

// Loader Plan for use in Load phase (Sample)
func loadPlanSampleFixture(t *testing.T, fixtureLabel string, pdfPath string) *hansard.HansardDocument {
	t.Helper()

	// Init ..
	//var hansardDoc *hansard.HansardDocument
	//pdfDoc := loadPDFFromFixture(t, fixtureLabel, pdfPath)
	// use the sample version so test cases not too long .. WILL be loaded from testdata fixture
	pdfDoc := samplePDFFromFixture(t, fixtureLabel, pdfPath)
	// Above  WILL FAIL for first case if TestNewHansardQuestions not run yet   when add new scenario

	// Should StateAssemblySession be so deep here?
	hansardDoc := &hansard.HansardDocument{StateAssemblySession: "confDUNSession"}
	hansard.NewHansardDocumentContent(pdfDoc, hansardDoc)
	// Read from cache; if not exist; complain that need to update
	//goldenPath := filepath.Join("testdata", fixtureLabel+".golden")
	//fmt.Println("GOLDEN FILE: ", goldenPath)

	// For comparison, if empty; fill it with zero out value
	if len(hansardDoc.HansardQuestions) == 0 {
		hansardDoc.HansardQuestions = []hansard.HansardQuestion{}
	}
	// For comparison; nil HansardQuestions; need to fill it!!
	// Alternative way to do it ..
	//if hansardDoc.HansardQuestions == nil {
	//	hansardDoc.HansardQuestions = []hansard.HansardQuestion{}
	//}

	// If it did NOT fatakl; we just proceed  ..
	return hansardDoc
}

// Loader Plan for use in ExecuteSplit
func loadPlanFromGolden(t *testing.T, goldenLabel string, hansardDoc *hansard.HansardDocument) {
	// Assumption: There is the equivalent fixture from hansard_integration_test
	t.Helper()

	// Read from cache; if not exist; complain that need to update
	goldenPath := filepath.Join("testdata", goldenLabel+".golden")
	// DEBUG
	//fmt.Println("GOLDEN FILE: ", goldenPath)

	// Case when update  flag is passed; update and eye-ball the changes ,,
	if *updateSplitterGolden {
		// Persist the golden file once ..
		w, werr := yaml.Marshal(hansardDoc)
		if werr != nil {
			t.Fatalf("Marshal FAIL: %s", werr.Error())
		}
		ioutil.WriteFile(goldenPath, w, 0644)
		// Done; get out!
		return
	}

	// Normal path,read from fixture ,..
	golden, rerr := ioutil.ReadFile(goldenPath)
	if rerr != nil {
		if os.IsNotExist(rerr) {
			t.Fatalf("Ensure run with flag -updateSplitterGolden first time! ERR: %s", rerr.Error())
		}
		t.Fatalf("Unexpected error: %s", rerr.Error())
	}
	umerr := yaml.Unmarshal(golden, hansardDoc)
	if umerr != nil {
		t.Fatalf("Unmarshal FAIL: %s", umerr.Error())
	}

	return
}
