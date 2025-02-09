package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"mime"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

type Context struct {
	IndentLevel int      // 当前缩进级别
	InPre       bool     // 是否在pre/code标签内
	InTable     bool     // 是否在表格中
	TableBuffer []string // 表格数据缓存
}

var (
	blockElements = map[string]bool{
		"p": true, "div": true, "h1": true, "h2": true, "h3": true,
		"h4": true, "h5": true, "h6": true, "ul": true, "ol": true,
		"li": true, "br": true, "hr": true, "pre": true, "code": true,
		"blockquote": true,
		"table":      true, "tr": true, "td": true, "th": true,
	}
	// 需要保留属性的标签及对应属性
	attrWhitelist = map[string][]string{
		"a":   {"href"},
		"img": {"src", "alt"},
	}
)

func GetHtmlText(r io.Reader) (string, error) {

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
	// 将所有的<br>标签替换为换行符
	doc.Find("br").Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithHtml("\n")
	})
	// 移除不需要的元素
	doc.Find("script, style, comment()").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	var buf strings.Builder
	body := doc.Find("body")
	if body.Length() == 0 {
		body = doc.Contents()
	}

	// 递归处理节点
	// walkNode(body.Nodes[0], &buf, false)
	walkNode(body.Nodes[0], &buf, Context{
		IndentLevel: 0,
		InPre:       false,
		InTable:     false,
	})
	text := buf.String()
	re := regexp.MustCompile(`\s*\n\s*`)
	text = re.ReplaceAllString(text, "\n")
	text = strings.ReplaceAll(text, " \n", "\n")
	return text, nil
}

func walkNode(n *html.Node, buf *strings.Builder, ctx Context) {
	switch n.Type {
	case html.TextNode:
		processText(n.Data, buf, ctx)

	case html.ElementNode:
		// tagName := n.Data
		newCtx := processElement(n, buf, ctx)

		// 处理子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walkNode(c, buf, newCtx)
		}

		// 后处理
		postProcessElement(n, buf, newCtx)

	case html.CommentNode:
		// 忽略注释

	case html.DoctypeNode:
		// 忽略文档类型

	case html.RawNode:
		// 处理CDATA节点
		if n.Data == "![CDATA[" {
			buf.WriteString(strings.TrimSpace(n.FirstChild.Data))
		}
	}
}

func processText(text string, buf *strings.Builder, ctx Context) {
	if ctx.InPre {
		buf.WriteString(text)
		return
	}

	// 普通文本处理
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return
	}

	if ctx.InTable {
		// 表格内容暂存
		ctx.TableBuffer = append(ctx.TableBuffer, trimmed)
	} else {
		writeIndent(buf, ctx.IndentLevel)
		buf.WriteString(trimmed + " ")
	}
}

func processElement(n *html.Node, buf *strings.Builder, ctx Context) Context {
	tagName := n.Data
	newCtx := ctx

	switch tagName {
	case "pre", "code":
		newCtx.InPre = true
		buf.WriteString("\n")

	case "blockquote":
		newCtx.IndentLevel += 2
		writeIndent(buf, ctx.IndentLevel)
		buf.WriteString("> ")

	case "ul", "ol":
		newCtx.IndentLevel++

	case "li":
		writeIndent(buf, ctx.IndentLevel-1)
		buf.WriteString("* ")

	case "table":
		newCtx.InTable = true
		newCtx.TableBuffer = []string{}

	case "tr":
		newCtx.TableBuffer = []string{}

	case "td", "th":
		// 表格单元格处理

	case "br":
		if !ctx.InPre {
			buf.WriteString("\n")
		}

	case "a", "img":
		// 处理白名单属性
		if attrs := extractAttributes(n, tagName); attrs != "" {
			buf.WriteString(attrs)
		}
	}

	return newCtx
}

func postProcessElement(n *html.Node, buf *strings.Builder, ctx Context) {
	tagName := n.Data

	switch tagName {
	case "pre", "code":
		buf.WriteString("\n")
		ctx.InPre = false

	case "blockquote":
		ctx.IndentLevel -= 2
		buf.WriteString("\n")

	case "ul", "ol":
		ctx.IndentLevel--
		buf.WriteString("\n")

	case "td", "th":
		// 单元格处理
		if ctx.InTable {
			ctx.TableBuffer = append(ctx.TableBuffer, "|")
		}

	case "tr":
		if ctx.InTable && len(ctx.TableBuffer) > 0 {
			writeIndent(buf, ctx.IndentLevel)
			buf.WriteString(strings.Join(ctx.TableBuffer, " "))
			buf.WriteString("\n")
		}

	case "table":
		buf.WriteString("\n")
		ctx.InTable = false
	}

	if blockElements[tagName] && !ctx.InPre {
		buf.WriteString("\n")
	}
}

func writeIndent(buf *strings.Builder, level int) {
	buf.WriteString(strings.Repeat("  ", level))
}

func extractAttributes(n *html.Node, tag string) string {
	var attrs []string
	for _, attr := range n.Attr {
		if contains(attrWhitelist[tag], attr.Key) {
			attrs = append(attrs, fmt.Sprintf("[%s=%q]", attr.Key, attr.Val))
		}
	}
	return strings.Join(attrs, "")
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
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
