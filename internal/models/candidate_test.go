package models

import (
	"github.com/gkampitakis/go-snaps/snaps"
	"io"
	"os"
	"testing"
)

func TestFromXML(t *testing.T) {
	testXmlFile, err := os.Open("../testdata/test.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer testXmlFile.Close()
	testXml, err := io.ReadAll(testXmlFile)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		xmlContents []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{xmlContents: testXml}, false},
		{"success", args{xmlContents: []byte{}}, true},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := FromXML(tt.args.xmlContents)
				if (err != nil) != tt.wantErr {
					t.Errorf("fromXML() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				snaps.MatchSnapshot(t, got)
			},
		)
	}
}
