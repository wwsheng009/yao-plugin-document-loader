package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"loader/loaders"
	"loader/textsplitter"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-hclog"
	jsoniter "github.com/json-iterator/go"
	"github.com/xuri/excelize/v2"
	"github.com/yaoapp/kun/grpc"
)

// PDF a simple pdf reader plugin
type DocumentLoader struct{ grpc.Plugin }

func (plugin *DocumentLoader) setLogFile() {
	var output io.Writer = os.Stdout
	//开启日志
	logroot := os.Getenv("GOU_TEST_PLG_LOG")
	if logroot == "" {
		logroot = "./logs"
	}
	if logroot != "" {
		logfile, err := os.OpenFile(path.Join(logroot, "docloader.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err == nil {
			output = logfile
		}
	}
	plugin.Plugin.SetLogger(output, grpc.Trace)
}

func getResponse(v interface{}, err error) (*grpc.Response, error) {

	if err != nil {
		bytes, err := jsoniter.Marshal(map[string]interface{}{"code": 400, "message": err.Error()})
		if err != nil {
			return nil, err
		}
		return &grpc.Response{Bytes: bytes, Type: "map"}, nil
	} else {
		bytes, err := jsoniter.Marshal(map[string]interface{}{"code": 200, "data": v})
		if err != nil {
			return nil, err
		}
		return &grpc.Response{Bytes: bytes, Type: "map"}, nil
	}

}

// Exec execute the plugin and return the result
func (doc *DocumentLoader) Exec(method string, args ...interface{}) (*grpc.Response, error) {
	doc.Logger.Log(hclog.Trace, "plugin method called", method)
	doc.Logger.Log(hclog.Trace, "args", args)

	if len(args) == 0 {
		return getResponse(nil, errors.New("missing file path"))
	}

	path, ok := args[0].(string)
	if !ok {
		return getResponse(nil, errors.New("invalid file path"))
	}

	ftype, err := getFileType(path)
	if err != nil {
		return getResponse(nil, err)
	}
	size := -1
	if len(args) > 1 {
		size, ok = args[1].(int)
	}

	switch strings.ToLower(method) {
	case "notation":
		if ftype == "DIR" {
			// Create a NotionDirectoryLoader instance
			loader := loaders.NewNotionDirectory(path)
			// Load documents from the test directory
			docs, err := loader.Load()
			if err != nil {
				return getResponse(nil, err)
			}
			bytes, err := jsoniter.Marshal(docs)
			if err != nil {
				return nil, err
			}
			return getResponse(bytes, err)
		} else {
			return getResponse(nil, fmt.Errorf("%s is not director", path))
		}
	case "text":
		f, err := os.Open(path)
		if err != nil {
			return getResponse(nil, err)
		}
		defer f.Close()
		finfo, err := f.Stat()
		if err != nil {
			return getResponse(nil, err)
		}
		var loader loaders.Loader
		var splitter textsplitter.TextSplitter
		switch ftype {
		case "WIZ":
			loader = loaders.NewWIZ(f, finfo.Size())
			splitter = textsplitter.NewRecursiveCharacter(textsplitter.WithChunkSize(size))
		case "DOCX":
			loader = loaders.NewDocx(f, finfo.Size())
			splitter = textsplitter.NewRecursiveCharacter(textsplitter.WithChunkSize(size))
		case "PPTX":
			loader = loaders.NewPPTX(f, finfo.Size())
			splitter = textsplitter.NewRecursiveCharacter(textsplitter.WithChunkSize(size))
		case "XLSX":
			if len(args) > 1 {
				password, ok := args[1].(string)
				if ok {
					loader = loaders.NewExcelx(f, excelize.Options{Password: password})
				}
			} else {
				loader = loaders.NewExcelx(f)
			}

			splitter = textsplitter.NewRecursiveCharacter(textsplitter.WithChunkSize(size))
		case "PDF":
			if len(args) > 1 {
				password, ok := args[1].(string)
				if ok {
					loader = loaders.NewPDF(f, finfo.Size(), loaders.PdfWithPassword(password))
				}
			} else {
				loader = loaders.NewPDF(f, finfo.Size())
			}
			splitter = textsplitter.NewRecursiveCharacter(textsplitter.WithChunkSize(size))
		case "MD":
			loader = loaders.NewText(f)
			splitter = textsplitter.NewMarkdownTextSplitter(textsplitter.WithChunkSize(size))
		case "HTML":
			loader = loaders.NewHTML(f)
			splitter = textsplitter.NewMarkdownTextSplitter(textsplitter.WithChunkSize(size))
		case "CSV":
			loader = loaders.NewCSV(f)
			splitter = textsplitter.NewMarkdownTextSplitter(textsplitter.WithChunkSize(size))
		case "TEXT":
			loader = loaders.NewText(f)
			splitter = textsplitter.NewRecursiveCharacter(textsplitter.WithChunkSize(size))
		}
		if loader == nil {
			return getResponse(nil, fmt.Errorf("%s not support:%s", ftype, path))
		}
		docs, err := loader.LoadAndSplit(context.Background(), splitter)
		return getResponse(docs, err)

	}
	return getResponse(nil, errors.New("invalid method"))

}

func main() {
	plugin := &DocumentLoader{}
	plugin.setLogFile()
	grpc.Serve(plugin)
}
