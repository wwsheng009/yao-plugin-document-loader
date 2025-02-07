package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// isPlainTextFile checks if the file is a plain text file
func isPlainTextFile(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the first 512 bytes of the file
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return false, err
	}

	// Check if bytes are in printable range
	for _, b := range buf {
		if (b < 32 || b > 126) && b != 10 && b != 13 {
			return false, nil
		}
	}
	return true, nil
}
func isTextFile(filename string) (bool, error) {
    file, err := os.Open(filename)
    if err != nil {
        return false, err
    }
    defer file.Close()

    reader := bufio.NewReader(file)

    // Read up to 1024 bytes or until EOF, whichever comes first
    buf, err := reader.Peek(1024)
    if err != nil && err != bufio.ErrBufferFull {
        return false, err
    }

    // Check for non-printable characters
    nonPrintableCount := 0
    for _, ch := range buf {
        if !unicode.IsPrint(rune(ch)) && !unicode.IsSpace(rune(ch)) {
            nonPrintableCount++
            if nonPrintableCount > 10 { // Arbitrary threshold for non-text files
                return false, nil
            }
        }
    }

    // If we've made it here with few non-printable chars, it's likely text
    return true, nil
}


// getFileType returns the file type based on the file extension
func getFileType(fileName string) (string, error) {
	fileType := "Unknown"
	info, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return fileType, err
			// fmt.Println("File or directory does not exist.")
		} else {
			return fileType, err
			// Handle other potential errors.
			// fmt.Println("Error:", err)
		}
	}
	if info.IsDir() {
		return "DIR", nil
		// fmt.Println(fileName, "is a directory.")
	}

	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {

	case ".docx":
		fileType = "DOCX"
	case ".xlsx":
		fileType = "XLSX"
	case ".pptx":
		fileType = "PPTX"
	case ".pdf":
		fileType = "PDF"
	case ".md", ".mdx":
		fileType = "MD"
	case ".html":
		fileType = "HTML"
	case ".csv":
		fileType = "CSV"
	case ".ziw":
		fileType = "WIZ"
	case ".txt", ".text", ".log":
		fileType = "TEXT"
	// Add more cases as needed
	default:
		istext, err := isPlainTextFile(fileName)
		if istext && err == nil {
			fileType = "TEXT"
		}
	}
	return fileType, nil
}