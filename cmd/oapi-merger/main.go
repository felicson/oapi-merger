package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

var (
	wdir string
	spec string
)

func init() {
	flag.StringVar(&wdir, "wdir", "./", "wdir=./")
	flag.StringVar(&spec, "spec", "openapi.yaml", "spec=./")
}

func main() {
	flag.Parse()

	if err := os.Chdir(wdir); err != nil {
		log.Fatal(err)
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, uri *url.URL) ([]byte, error) {
		return ioutil.ReadFile(uri.Path)
	}
	doc, err := loader.LoadFromFile(spec) // i.e "api/openapi.yaml")
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

	data, err := doc.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
