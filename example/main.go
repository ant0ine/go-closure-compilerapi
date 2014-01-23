package main

import (
        "github.com/ant0ine/go-closure-compilerapi"
        "log"
)

func main() {
        client := compilerapi.Client{}
        output := client.Compile([]byte("var i = 0 // test"))
	log.Printf("%+v", output)
}
