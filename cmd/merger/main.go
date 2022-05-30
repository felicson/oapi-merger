package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
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
	doc, err := loader.LoadFromFile(filepath.Join(spec)) // i.e "api/openapi.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if err = doc.Validate(loader.Context); err != nil {
		log.Fatal(err)
	}

	doc.InternalizeRefs(context.Background(), func(ref string) string {
		if idx := strings.Index(ref, "#"); idx > 1 {
			ref = ref[idx+2:]
			return ref[strings.LastIndex(ref, "/")+2:]
		}
		return strings.TrimRight(filepath.Base(ref), filepath.Ext(ref))
	})

	if err = doc.Validate(loader.Context); err != nil {
		log.Fatal(err)
	}

	data, err := doc.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
