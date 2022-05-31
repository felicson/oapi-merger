# oapi-merger

This tool helps to merge openapi 3 version specification to one file, based on references in files.

## Install
```bash
go install github.com/felicson/oapi-merger/cmd/oapi-merger@latest
```
After install step the `oapi-merger` binary should appear in your $GOPATH/bin directory.

## Usage
```bash
oapi-merger -wdir examples/api -spec openapi.yaml
```

Flags:  
    -wdir - directory with specification hierarchy  
    -spec - input root file with openapi spec (default: openapi.yaml)  
    -o    - where to output generated code, stdout is default