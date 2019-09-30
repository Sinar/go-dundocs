package hansard

import (
	"testing"
)

func Test_extractQuestionNum(t *testing.T) {
	type args struct {
		rowContent string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractQuestionNum(tt.args.rowContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractQuestionNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractQuestionNum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isStartOfQuestionSection(t *testing.T) {
	type args struct {
		rowContent string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isStartOfQuestionSection(tt.args.rowContent); got != tt.want {
				t.Errorf("isStartOfQuestionSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_detectHansardType(t *testing.T) {
	type args struct {
		firstPage PDFPage
	}
	tests := []struct {
		name string
		args args
		want HansardType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := detectHansardType(tt.args.firstPage); got != tt.want {
				t.Errorf("detectHansardType() = %v, want %v", got, tt.want)
			}
		})
	}
}
