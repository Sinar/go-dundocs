package hansard_test

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/sanity-io/litter"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"

	"gopkg.in/yaml.v2"

	"github.com/Sinar/go-dundocs/internal/hansard"
)

var updateSplitterGolden = flag.Bool("updateSplitterGolden", false, "update SplitHansardDocPlan .golden files")

//var updateSplitterFixture = flag.Bool("updateSplitterFixture", false, "update SplitHansardDocPlan .fixture plans")

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
			if serr := s.SavePlan(); (serr != nil) != tt.wantErr {
				t.Errorf("SavePlan() error = %v, wantErr %v", serr, tt.wantErr)
				return
			}
			// Check planDir actually created and exists ..
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
			// If plan file exist; check its string equivalent
			content, rerr := ioutil.ReadFile(absolutePlanDir + "/split.yml")
			if rerr != nil {
				t.Fatalf("ERR: %s", rerr.Error())
			}
			if diff := cmp.Diff(litter.Sdump(s), litter.Sdump(content)); diff != "" {
				t.Errorf("Plan mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func TestSplitHansardDocumentPlan_LoadPlan(t *testing.T) {
	type args struct {
		fixtureLabel string
	}
	tests := []struct {
		name    string
		args    args
		planDir string
		wantErr bool
	}{
		// NOTE: Not too nice; still tied to physical path that are not in testdata
		{"happy #1", args{"HDOC-Lisan-1-20"}, "testdata/split-case-1.yaml", false},
		{"happy #2", args{"HDOC-BukanLisan-1-20"}, "testdata/split-case-2.yaml", false},
		{"sad #1", args{"Bad-HDOC-Lisan-1-20"}, "testdata/split-case-3.yaml", false},
		{"sad #2", args{"Bad-HDOC-BukanLisan-1-20"}, "testdata/split-case-4.yaml", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup where the plan is
			s := &hansard.SplitHansardDocumentPlan{
				PlanDir: tt.planDir,
			}
			// Load it into the struct
			if err := s.LoadPlan(); (err != nil) != tt.wantErr {
				t.Errorf("LoadPlan() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := s.HansardDocument
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
	type args struct {
		planLabel string // Plan Fixture .. not needed, we can make it small enough
	}
	type fields struct {
		planDir         string
		hansardDocument hansard.HansardDocument
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"no plan #1", fields{"", hansard.HansardDocument{
			StateAssemblySession: "",
		}}, true},
		{"happy #1", fields{"", hansard.HansardDocument{
			StateAssemblySession: "",
			HansardType:          0,
			HansardQuestions:     nil,
		}}, false},
		{"happy #2", fields{"", hansard.HansardDocument{
			HansardType:      0,
			HansardQuestions: nil,
		}}, false},
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
			//defer os.RemoveAll(dir)
			log.Println("Dir is ", dir)

			s := &hansard.SplitHansardDocumentPlan{
				PlanDir:         tt.fields.planDir,
				HansardDocument: tt.fields.hansardDocument,
			}
			// Prep first; load actual plan from the test cases ..
			lerr := s.LoadPlan()
			if lerr != nil {
				// Should NOT happen!
				t.Fatal(lerr)
			}
			// Execute the actual split ...
			if err := s.ExecuteSplit(); (err != nil) != tt.wantErr {
				t.Errorf("ExecuteSplit() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Open the output in the WorkDir and see if it has the expected number of files
			// Should we check content? probably no need .. what about filename; is it important; no for now ..
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
	hansardDoc := &hansard.HansardDocument{StateAssemblySession: "testGoldenPlan"}
	hansard.NewHansardDocumentContent(pdfDoc, hansardDoc)
	// Read from cache; if not exist; complain that need to update
	//goldenPath := filepath.Join("testdata", fixtureLabel+".golden")
	//fmt.Println("GOLDEN FILE: ", goldenPath)

	// Read it from planPath; which in this case is just the Golden Plan ..
	loadPlanFromGolden(t, fixtureLabel, hansardDoc)
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
