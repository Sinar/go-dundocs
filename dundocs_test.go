package dundocs

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestDUNDocs_Plan(t *testing.T) {
	type fields struct {
		DUNSession string
		Conf       Configuration
		Options    *ExtractPDFOptions
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dd := &DUNDocs{
				DUNSession: tt.fields.DUNSession,
				Conf:       tt.fields.Conf,
				Options:    tt.fields.Options,
			}
			// Execute the plan
			dd.Plan()
			spew.Dump(dd)
		})
	}
}

func TestDUNDocs_Reset(t *testing.T) {
	type fields struct {
		DUNSession string
		Conf       Configuration
		Options    *ExtractPDFOptions
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dd := &DUNDocs{
				DUNSession: tt.fields.DUNSession,
				Conf:       tt.fields.Conf,
				Options:    tt.fields.Options,
			}
			spew.Dump(dd)

		})
	}
}

func TestDUNDocs_Split(t *testing.T) {
	type fields struct {
		DUNSession string
		Conf       Configuration
		Options    *ExtractPDFOptions
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dd := &DUNDocs{
				DUNSession: tt.fields.DUNSession,
				Conf:       tt.fields.Conf,
				Options:    tt.fields.Options,
			}
			spew.Dump(dd)

		})
	}
}

func TestNewDUNDocs(t *testing.T) {
	tests := []struct {
		name string
		want *DUNDocs
	}{
		// TODO: Add test cases.
		// Scemario: Options, 2 pages out of 5?
		// Default source
		// Default data
		// Default plan
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDUNDocs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDUNDocs() = %v, want %v", got, tt.want)
			}
		})
	}
}
