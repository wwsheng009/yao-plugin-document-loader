package docx

import (
	"archive/zip"
	"encoding/xml"
	"io"
)

type Document struct {
	Body Body `xml:"body"`
}

type Body struct {
	Paragraphs []Paragraph `xml:"p"`
}

type Paragraph struct {
	Texts []Text `xml:"r>t"`
}

type Text struct {
	Content string `xml:",chardata"`
}

func Read(r io.ReaderAt, size int64) ([]Paragraph, error) {
	// Open the .docx file
	zipReader, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	// Find the document.xml file and read its content
	var documentXML []byte
	for _, file := range zipReader.File {
		if file.Name == "word/document.xml" {
			f, err := file.Open()
			if err != nil {
				return nil, err
			}

			documentXML, err = io.ReadAll(f)
			f.Close()
			if err != nil {
				return nil, err
			}
			break
		}
	}
	// Unmarshal the XML content into a Document struct
	var doc Document
	err = xml.Unmarshal(documentXML, &doc)
	if err != nil {
		return nil, err
	}
	return doc.Body.Paragraphs, nil
}
