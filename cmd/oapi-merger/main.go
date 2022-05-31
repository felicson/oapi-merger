package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

var (
	flagWorkDir    string
	flagSpecFile   string
	flagOutputFile string
)

func main() {
	flag.StringVar(&flagWorkDir, "wdir", "./", "Working directory")
	flag.StringVar(&flagSpecFile, "spec", "openapi.yaml", "Entry for specification description, openapi.yaml is default")
	flag.StringVar(&flagOutputFile, "o", "", "Where to output generated code, stdout is default")

	flag.Parse()

	if err := os.Chdir(flagWorkDir); err != nil {
		log.Fatal(err)
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, uri *url.URL) ([]byte, error) {
		return os.ReadFile(uri.Path)
	}
	doc, err := loader.LoadFromFile(flagSpecFile) // i.e "api/openapi.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// validate before merging
	if err = doc.Validate(loader.Context); err != nil {
		log.Fatal(err)
	}

	doc.InternalizeRefs(context.Background(), nil)

	// validate after merging
	if err = doc.Validate(loader.Context); err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", " ")
	if err := enc.Encode(doc); err != nil {
		log.Fatal(err)
	}
	if flagOutputFile != "" {
		if err := os.WriteFile(flagOutputFile, buf.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}
	} else {
		println(buf.String())
	}
}
