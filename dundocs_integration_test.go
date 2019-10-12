package dundocs_test

import (
	"testing"

	"github.com/Sinar/go-dundocs/internal/hansard"

	"github.com/davecgh/go-spew/spew"

	"github.com/Sinar/go-dundocs"
)

func TestDUNDocs_Plan(t *testing.T) {
	type fields struct {
		Conf dundocs.Configuration
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"", fields{dundocs.Configuration{
			"", hansard.HANSARD_INVALID,
			"", "",
			dundocs.PLAN},
		}},
		{"", fields{dundocs.Configuration{
			"", hansard.HANSARD_INVALID,
			"", "",
			dundocs.PLAN},
		}},
		{"", fields{dundocs.Configuration{
			"", hansard.HANSARD_INVALID,
			"", "",
			dundocs.PLAN},
		}},
		{"", fields{dundocs.Configuration{
			"", hansard.HANSARD_INVALID,
			"", "",
			dundocs.PLAN},
		}},
		{"", fields{dundocs.Configuration{
			"", hansard.HANSARD_INVALID,
			"", "",
			dundocs.PLAN},
		}},
	}
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
	tests := []struct {
		name   string
		fields fields
	}{
		{"", fields{dundocs.Configuration{
			"", hansard.HANSARD_INVALID,
			"", "",
			dundocs.SPLIT}},
		},
		{"", fields{dundocs.Configuration{
			"", hansard.HANSARD_INVALID,
			"", "",
			dundocs.SPLIT}},
		},
		{"", fields{dundocs.Configuration{
			"", hansard.HANSARD_INVALID,
			"", "",
			dundocs.SPLIT}},
		},
		{"", fields{dundocs.Configuration{
			"", hansard.HANSARD_INVALID,
			"", "",
			dundocs.SPLIT}},
		},
		{"", fields{dundocs.Configuration{
			"", hansard.HANSARD_INVALID,
			"", "",
			dundocs.SPLIT}},
		},
	}
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
	}{
		{"", fields{dundocs.Configuration{
			"", hansard.HANSARD_INVALID,
			"", "",
			dundocs.RESET}},
		},
	}
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
