package compilerapi

import (
	"encoding/json"
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

type OutputWarning struct {
	Charno  int    `json:"charno"`
	Warning string `json:"warning"`
	Lineno  int    `json:"lineno"`
	File    string `json:"file"`
	Type    string `json:"type"`
	Line    string `json:"line"`
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

        // TODO support WHITESPACE_ONLY, SIMPLE_OPTIMIZATIONS, ADVANCED_OPTIMIZATIONS
	values.Set("compilation_level", "SIMPLE_OPTIMIZATIONS")

	values.Set("output_format", "json")
	values.Add("output_info", "compiled_code")
	values.Add("output_info", "statistics")
	values.Add("output_info", "warnings")
	values.Add("output_info", "errors")

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
