package dundocs_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/Sinar/go-dundocs"
)

func TestDUNDocs_Plan(t *testing.T) {
	type fields struct {
		Conf dundocs.Configuration
	}
	// Test normal plan
	// Test plan + split in one; flags being passed in
	// SImulate from UI too?
	tests := []struct {
		name   string
		fields fields
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dd := &dundocs.DUNDocs{
				Conf: tt.fields.Conf,
			}
			// DEBUG
			spew.Dump(dd)

		})
	}
}

func TestDUNDocs_Split(t *testing.T) {
	type fields struct {
		Conf dundocs.Configuration
	}
	// Test normal split
	// SImulate from UI too?
	tests := []struct {
		name   string
		fields fields
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dd := &dundocs.DUNDocs{
				Conf: tt.fields.Conf,
			}
			// DEBUG
			spew.Dump(dd)
		})
	}
}

func TestDUNDocs_Reset(t *testing.T) {
	type fields struct {
		Conf dundocs.Configuration
	}
	tests := []struct {
		name   string
		fields fields
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dd := &dundocs.DUNDocs{
				Conf: tt.fields.Conf,
			}
			// DEBUG
			spew.Dump(dd)
		})
	}
}
