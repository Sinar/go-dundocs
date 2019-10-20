package hansard

import (
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
				workingDir:      tt.fields.workingDir,
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
			if err := prepareSplitAPI(); (err != nil) != tt.wantErr {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare TemoDir working ..
			gotAbsolutePath, gotBaseName, gotExtension := normalizeToAbsolutePath(tt.args.relativePath)
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
