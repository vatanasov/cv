package extractor

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

func TestExtract(t *testing.T) {
	resultFile, err := os.Open("../testdata/test.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer resultFile.Close()
	result, err := io.ReadAll(resultFile)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr error
	}{
		{"missing file", args{"invalid_file"}, nil, ErrMissingFile},
		{"works", args{"../testdata/test_with_attachment.pdf"}, result, nil},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := ExtractXML(tt.args.filename)
				if err != nil {
					switch {
					case errors.Is(err, tt.wantErr):
						return
					default:
						t.Errorf("ExtractXML() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
				}
				if !bytes.Equal(got, tt.want) {
					t.Errorf("ExtractXML() got = %q, want %q", got, tt.want)
				}
			},
		)
	}
}
