package dundocs_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/Sinar/go-dundocs"
)

func TestDUNDocs_Plan(t *testing.T) {
	type fields struct {
		DUNSession string
		Conf       dundocs.Configuration
		Options    *dundocs.ExtractPDFOptions
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"BukanLisan", fields{
			"duntest-sesi01",
			dundocs.Configuration{
				SourcePDFPath: "./raw/BukanLisan/41 - 60.pdf",
				//WorkingDir:    ".",
				//DataDir:       "./data",
			},
			nil},
		},
		{"Lisan", fields{
			"duntest-sesi01",
			dundocs.Configuration{
				SourcePDFPath: "./raw/Lisan/SOALAN MULUT (261-272).pdf",
				WorkingDir:    "/tmp",
				DataDir:       "DUNDOC",
			},
			nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dd := &dundocs.DUNDocs{
				DUNSession: tt.fields.DUNSession,
				Conf:       tt.fields.Conf,
				Options:    tt.fields.Options,
			}
			dd.Plan()
			// TODO: Should it check content; or only location; or if it does not
			// 	go up in flames is considered OK?
		})
	}
}

func TestDUNDocs_Split(t *testing.T) {
	type fields struct {
		DUNSession string
		Conf       dundocs.Configuration
		Options    *dundocs.ExtractPDFOptions
	}
	// Test normal split
	// SImulate from UI too?
	tests := []struct {
		name   string
		fields fields
	}{
		{"BukanLisan", fields{
			"duntest-sesi01",
			dundocs.Configuration{
				SourcePDFPath: "./raw/BukanLisan/41 - 60.pdf",
				//WorkingDir:    ".",
				//DataDir:       "./data",
			},
			nil},
		},
		{"Lisan", fields{
			"duntest-sesi01",
			dundocs.Configuration{
				SourcePDFPath: "./raw/Lisan/SOALAN MULUT (261-272).pdf",
				WorkingDir:    "/tmp",
				DataDir:       "DUNDOC",
			},
			nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dd := &dundocs.DUNDocs{
				DUNSession: tt.fields.DUNSession,
				Conf:       tt.fields.Conf,
				Options:    tt.fields.Options,
			}
			// DEBUG
			//spew.Dump(dd)
			dd.Split()
			// Where to check where the output is placed
			// default will be <workingDir>/splitout ?
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
