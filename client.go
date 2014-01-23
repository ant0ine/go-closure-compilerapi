package compilerapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	GoogleEndpointUrl = "http://closure-compiler.appspot.com/compile"
)

type Client struct {
}

type OutputError struct {
	Charno int    `json:"charno"`
	Error  string `json:"error"`
	Lineno int    `json:"lineno"`
	File   string `json:"file"`
	Type   string `json:"type"`
	Line   string `json:"line"`
}

func (e *OutputError) AsLogline() string {
	return fmt.Sprintf("\033[36;1m[%d, %d]\033[31m error: \033[0m%s\n\t%s\n",
		e.Lineno,
		e.Charno,
		e.Error,
		e.Line,
	)
}

type OutputWarning struct {
	Charno  int    `json:"charno"`
	Warning string `json:"warning"`
	Lineno  int    `json:"lineno"`
	File    string `json:"file"`
	Type    string `json:"type"`
	Line    string `json:"line"`
}

func (w *OutputWarning) AsLogline() string {
	return fmt.Sprintf("\033[36;1m[%d, %d]\033[33m warning: \033[0m%s\n\t%s\n",
		w.Lineno,
		w.Charno,
		w.Warning,
		w.Line,
	)
}

type OutputServerError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type OutputStatistics struct {
	OriginalSize   int `json:"originalSize"`
	CompressedSize int `json:"compressedSize"`
	CompileTime    int `json:"compileTime"`
}

type Output struct {
	CompiledCode string             `json:"compiledCode"`
	Errors       []OutputError      `json:"errors"`
	Warnings     []OutputWarning    `json:"warnings"`
	ServerErrors *OutputServerError `json:"serverErrors"`
	Statistics   OutputStatistics   `json:"statistics"`
}

func (client *Client) buildRequest(jsCode []byte) *http.Request {

	values := url.Values{}
	values.Set("js_code", string(jsCode[:]))
	values.Set("output_format", "json")
	values.Add("output_info", "compiled_code")
	values.Add("output_info", "statistics")
	values.Add("output_info", "warnings")
	values.Add("output_info", "errors")

	// TODO support ECMASCRIPT3, ECMASCRIPT5, ECMASCRIPT5_STRICT
	values.Set("language", "ECMASCRIPT5_STRICT")

	// TODO support WHITESPACE_ONLY, SIMPLE_OPTIMIZATIONS, ADVANCED_OPTIMIZATIONS
	values.Set("compilation_level", "SIMPLE_OPTIMIZATIONS")

	req, err := http.NewRequest(
		"POST",
		GoogleEndpointUrl,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		log.Fatalf(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func (client *Client) Compile(jsCode []byte) *Output {

	httpClient := http.Client{}

	req := client.buildRequest(jsCode)
	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf(err.Error())
	}

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}

	output := Output{}
	err = json.Unmarshal(content, &output)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return &output
}
