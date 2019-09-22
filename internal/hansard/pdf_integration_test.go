package hansard_test

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Sinar/go-dundocs/internal/hansard"
	"github.com/sanity-io/litter"
	"gopkg.in/yaml.v2"
)

var update = flag.Bool("update", false, "update .fixture files")

func TestNewPDFDocument(t *testing.T) {
	type args struct {
		pdfPath string
		options *hansard.ExtractPDFOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"missing file, xerrors", args{"testdata/bogus.pdf", nil}, true},
		{"happy path #1, Lisan", args{"/Users/leow/GOMOD/go-dundocs/raw/Lisan/SOALAN MULUT (261-272).pdf", &hansard.ExtractPDFOptions{NumPages: 19}}, false},
		{"happy path #2, BukanLisan", args{"/Users/leow/GOMOD/go-dundocs/raw/BukanLisan/41 - 60.pdf", nil}, false},
		{"table files, lampiran", args{"/Users/leow/GOMOD/go-dundocs/raw/Lampiran/LAMPIRAN B JAWAPAN SOALAN NO.213.pdf", nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hansard.NewPDFDocument(tt.args.pdfPath, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPDFDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Special case of error; no structure; is nil
			if err != nil {
				if got != nil {
					t.Errorf("NewPDFDocument() should be nil; but is NOT!")
				}
				return
			}
			// Load from fixture if no update flag passed
			fixture := filepath.Join("testdata", tt.name+".fixture")
			if *update {
				// Write struct into YAML
				litter.Dump(got)

				wb, merr := yaml.Marshal(got)
				if merr != nil {
					t.Fatalf("Marsahl FAIL: %s", merr.Error())
				}
				// Write YAML in to it ..
				ioutil.WriteFile(fixture, wb, 0644)
			}
			// Load into  struct from YAML
			b, rerr := ioutil.ReadFile(fixture)
			if rerr != nil {
				// Cannot proceed with one golden file update
				if os.IsNotExist(rerr) {
					t.Fatalf("Ensure run with -update flag first time! ERR: %s", rerr.Error())
				}
				t.Fatalf("Unexpected error: %s", rerr.Error())
			}
			var want *hansard.PDFDocument
			umerr := yaml.Unmarshal(b, &want)
			if umerr != nil {
				t.Fatalf("Unmarsahl FAIL: %s", umerr.Error())
			}

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("TestNewPDFDocument() mismatch (-want +got):\n%s", diff)
			}

			//if !reflect.DeepEqual(got, want) {
			//	t.Errorf("NewPDFDocument() got = %v, want %v", got, want)
			//}
		})
	}
}
