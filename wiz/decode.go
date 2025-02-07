package wiz

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"mime"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html/charset"
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
			document, err = getHtmlText(f)
			f.Close()
			if err != nil {
				return "", err
			}
			break
		}
	}
	return document, nil
}

func getHtmlText(r io.ReadCloser) (string, error) {

	// rd := bytes.NewReader([]byte("xxxxxxxxx"))
	// buffer2 := make([]byte, 512)
	// n2, err := rd.Read(buffer2)
	// if err != nil {
	// 	return "", err
	// }
	// print(n2)

	buffer := make([]byte, 512)
	n, err := r.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}
	charsetLabel := params["charset"] // This is the actual charset (e.g., "utf-8")

	// Handle character encoding
	utf8Reader, err := charset.NewReaderLabel(charsetLabel, io.MultiReader(bytes.NewReader(buffer[:n]), r))
	if err != nil {
		return "", err
	}
	utf8Reader, err = trimBOM(utf8Reader)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		return "", err
	}

	var sel *goquery.Selection
	if doc.Has("body") != nil {
		sel = doc.Find("body").Contents()
	} else {
		sel = doc.Contents()
	}
	sanitized := bluemonday.UGCPolicy().Sanitize(sel.Text())
	pagecontent := strings.TrimSpace(sanitized)
	
	re := regexp.MustCompile(`\s*\n\s*`)
	pagecontent = re.ReplaceAllString(pagecontent, "\n")
	println(pagecontent)

	return pagecontent, nil
}

// trimBOM trims the byte order mark from the beginning of an io.Reader if present.
func trimBOM(reader io.Reader) (io.Reader, error) {
	bom := []byte("\ufeff")
	buffer := make([]byte, len(bom))

	// Read the first few bytes
	n, err := reader.Read(buffer)
	if err != nil {
		return nil, err
	}

	// Check if BOM is present
	if n == len(bom) && bytes.Equal(buffer, bom) {
		// BOM is present, return a new reader without the BOM
		return io.MultiReader(bytes.NewReader(buffer[n:]), reader), nil
	}

	// BOM is not present, return the original reader
	return io.MultiReader(bytes.NewReader(buffer[:n]), reader), nil
}

// UTF16ToUTF8 converts UTF-16 encoded data from an io.Reader to UTF-8.
func UTF16ToUTF8(reader io.Reader, byteOrder binary.ByteOrder) ([]byte, error) {
	var utf8Data []byte
	bufReader := bufio.NewReader(reader)
	eof := false

	for !eof {
		var rune uint16
		err := binary.Read(bufReader, byteOrder, &rune)
		if err == io.EOF {
			eof = true
			err = nil
		} else if err != nil {
			return nil, err
		}

		if !eof {
			utf8Data = append(utf8Data, encodeRune(rune)...)
		}
	}

	return utf8Data, nil
}

// encodeRune converts a single UTF-16 rune to its UTF-8 representation.
func encodeRune(r uint16) []byte {
	runeValue := utf16.Decode([]uint16{r})
	buf := make([]byte, 3)
	n := utf8.EncodeRune(buf, runeValue[0])
	return buf[:n]
}
