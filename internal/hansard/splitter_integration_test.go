package hansard_test

import (
	"reflect"
	"testing"

	"github.com/Sinar/go-dundocs/internal/hansard"
)

func TestNewSplitHansardDocumentPlan(t *testing.T) {
	tests := []struct {
		name string
		want *hansard.SplitHansardDocumentPlan
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hansard.NewSplitHansardDocumentPlan(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSplitHansardDocumentPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitHansardDocumentPlan_SavePlan(t *testing.T) {
	type fields struct {
		workingDir      string
		planDir         string
		hansardDocument hansard.HansardDocument
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &hansard.SplitHansardDocumentPlan{
				WorkingDir:      tt.fields.workingDir,
				PlanDir:         tt.fields.planDir,
				HansardDocument: tt.fields.hansardDocument,
			}
			if err := s.SavePlan(); (err != nil) != tt.wantErr {
				t.Errorf("SavePlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSplitHansardDocumentPlan_LoadPlan(t *testing.T) {
	type fields struct {
		workingDir      string
		planDir         string
		hansardDocument hansard.HansardDocument
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &hansard.SplitHansardDocumentPlan{
				WorkingDir:      tt.fields.workingDir,
				PlanDir:         tt.fields.planDir,
				HansardDocument: tt.fields.hansardDocument,
			}
			if err := s.LoadPlan(); (err != nil) != tt.wantErr {
				t.Errorf("LoadPlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSplitHansardDocumentPlan_ExecuteSplit(t *testing.T) {
	type fields struct {
		workingDir      string
		planDir         string
		hansardDocument hansard.HansardDocument
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &hansard.SplitHansardDocumentPlan{
				WorkingDir:      tt.fields.workingDir,
				PlanDir:         tt.fields.planDir,
				HansardDocument: tt.fields.hansardDocument,
			}
			if err := s.ExecuteSplit(); (err != nil) != tt.wantErr {
				t.Errorf("ExecuteSplit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
