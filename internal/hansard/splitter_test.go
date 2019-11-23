package hansard

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSplitHansardDocumentPlan_detectSessionName(t *testing.T) {
	type fields struct {
		workingDir      string
		planDir         string
		HansardDocument HansardDocument
	}
	type args struct {
		sourcePDFFileName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SplitHansardDocumentPlan{
				dataDir:         tt.fields.workingDir,
				PlanDir:         tt.fields.planDir,
				HansardDocument: tt.fields.HansardDocument,
			}
			if got := s.detectSessionName(tt.args.sourcePDFFileName); got != tt.want {
				t.Errorf("detectSessionName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadSplitHansardDocPlan(t *testing.T) {
	type args struct {
		splitPlanPath string
	}
	tests := []struct {
		name string
		args args
		want *HansardDocument
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadSplitHansardDocPlan(tt.args.splitPlanPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadSplitHansardDocPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepareSplit(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := prepareSplit(); (err != nil) != tt.wantErr {
				t.Errorf("prepareSplit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_prepareSplitAPI(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := prepareSplitAPI("", ""); (err != nil) != tt.wantErr {
				t.Errorf("prepareSplitAPI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_normalizeToAbsolutePath(t *testing.T) {
	type args struct {
		relativePath string
	}
	tests := []struct {
		name             string
		args             args
		wantAbsolutePath string
		wantBaseName     string
		wantExtension    string
	}{
		{"happy #1", args{"testdata/abc.pdf"}, "", "abc", "pdf"},
		{"caps PDF #2", args{"testdata/ABC.Def.PDF"}, "", "ABC_Def", "pdf"},
		{"weird #3", args{"/tmp/abc.def.123.pdf"}, "/tmp/", "abc_def_123", "pdf"},
		{"no pdf ext #4", args{"/tmp/abc.def.123"}, "/tmp/", "abc_def_123", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare TemoDir working ..
			// Prepare TempDir for working with it
			dir, err := ioutil.TempDir("", "dundocs-splitter")
			if err != nil {
				log.Fatal(err)
			}
			// Comment out below if need to see the output in dir
			defer os.RemoveAll(dir)
			log.Println("Dir is ", dir)
			// Go to dir to execute
			cerr := os.Chdir(dir)
			if cerr != nil {
				panic(cerr)
			}
			gotAbsolutePath, gotBaseName, gotExtension := normalizeToAbsolutePath(tt.args.relativePath)
			// If tt.wantAbsolutePath is empty; means  make it  to be tempDir
			if tt.wantAbsolutePath == "" {
				// Only if relative is actually relative
				if !filepath.IsAbs(tt.args.relativePath) {
					b, _ := filepath.Split(tt.args.relativePath)
					// NOTE: A bot of a hack, need a more elegant way? Maybe no need to check too ..
					tt.wantAbsolutePath = filepath.Join("/private", dir, b) + "/"
				}
			}
			if gotAbsolutePath != tt.wantAbsolutePath {
				t.Errorf("normalizeToAbsolutePath() gotAbsolutePath = %v, want %v", gotAbsolutePath, tt.wantAbsolutePath)
			}
			if gotBaseName != tt.wantBaseName {
				t.Errorf("normalizeToAbsolutePath() gotBaseName = %v, want %v", gotBaseName, tt.wantBaseName)
			}
			if gotExtension != tt.wantExtension {
				t.Errorf("normalizeToAbsolutePath() gotExtension = %v, want %v", gotExtension, tt.wantExtension)
			}
		})
	}
}

func Test_normalizeTempDirAbsolutePath(t *testing.T) {
	type args struct {
		relativePath string
	}
	tests := []struct {
		name             string
		args             args
		wantAbsolutePath string
		wantBaseName     string
		wantExtension    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAbsolutePath, gotBaseName, gotExtension := normalizeTempDirAbsolutePath(tt.args.relativePath)
			if gotAbsolutePath != tt.wantAbsolutePath {
				t.Errorf("normalizeTempDirAbsolutePath() gotAbsolutePath = %v, want %v", gotAbsolutePath, tt.wantAbsolutePath)
			}
			if gotBaseName != tt.wantBaseName {
				t.Errorf("normalizeTempDirAbsolutePath() gotBaseName = %v, want %v", gotBaseName, tt.wantBaseName)
			}
			if gotExtension != tt.wantExtension {
				t.Errorf("normalizeTempDirAbsolutePath() gotExtension = %v, want %v", gotExtension, tt.wantExtension)
			}
		})
	}
}
