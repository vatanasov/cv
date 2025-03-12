package extractor

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

var (
	ErrMissingFile = errors.New("missing file")
)

func ExtractXML(filename string) ([]byte, error) {
	var contents []byte
	file, err := os.Open(filename)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return contents, ErrMissingFile
		default:
			return contents, err
		}
	}

	defer file.Close()

	tempFile, err := os.CreateTemp("", "cv")
	if err != nil {
		return contents, err
	}

	defer tempFile.Close()

	cmd := exec.Command("pdfdetach", "-savefile", "attachment.xml", "-o", tempFile.Name(), filename)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return contents, err
	}

	contents, err = io.ReadAll(tempFile)
	if err != nil {
		return contents, err
	}

	return contents, nil
}
