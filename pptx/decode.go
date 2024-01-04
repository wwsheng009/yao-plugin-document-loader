package pptx

import (
	"archive/zip"
	"encoding/xml"
	"io"
	"regexp"
)

func isSlide(file *zip.File) bool {
	matched, err := regexp.MatchString("slides/slide(\\d+).xml", file.Name)
	if err != nil {
		return false
	}
	return matched
}
func Read(r io.ReaderAt, size int64) ([][]string, error) {
	zipReader, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}

	texts := make([][]string, 0)
	for _, file := range zipReader.File {
		if isSlide(file) {
			slideTexts, err := getSingleSlide(file)
			if err != nil {
				return nil, err
			}
			texts = append(texts, slideTexts)
		}
	}
	return texts, nil

}

func getSingleSlide(s *zip.File) ([]string, error) {
	var slideTexts []string
	r, err := s.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	d := xml.NewDecoder(r)
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		t, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}

		if t.Name.Local == "t" {
			charDataTok, err := d.Token()
			if err != nil {
				return nil, err
			}
			charData, ok := charDataTok.(xml.CharData)
			if !ok {
				continue
			}
			slideTexts = append(slideTexts, string(charData))
		}
		if t.Name.Local == "p" {
			slideTexts = append(slideTexts, "\n")
		}
	}
	return slideTexts, nil
}
