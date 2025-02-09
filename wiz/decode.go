package wiz

import (
	"archive/zip"
	"io"
	"loader/utils"
)

func Read(r io.ReaderAt, size int64) (string, error) {
	zipReader, err := zip.NewReader(r, size)
	if err != nil {
		return "", err
	}
	var document string
	for _, file := range zipReader.File {
		if file.Name == "index.html" {
			f, err := file.Open()
			if err != nil {
				return "", err
			}
			document, err = utils.GetHtmlText(f)
			f.Close()
			if err != nil {
				return "", err
			}
			break
		}
	}
	return document, nil
}
