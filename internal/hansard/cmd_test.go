package hansard

import "testing"

func TestPlanAndSave(t *testing.T) {
	type args struct {
		conf Configuration
	}
	// This will be an actual files?
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PlanAndSave(tt.args.conf); (err != nil) != tt.wantErr {
				t.Errorf("PlanAndSave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadAndSplit(t *testing.T) {
	type args struct {
		conf Configuration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadAndSplit(tt.args.conf); (err != nil) != tt.wantErr {
				t.Errorf("LoadAndSplit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
