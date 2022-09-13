package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/macaron.v1"
)

var pwd string

func init() {
	var err error
	pwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

func main() {
	m := macaron.Classic()
	m.Get("/static/*", staticHandler)
	log.Fatal(http.ListenAndServe(":8081", m))
}

func staticHandler(ctx *macaron.Context) (status int, out string) {
	path := filepath.Join(pwd, "static", filepath.Clean(ctx.Params("*")))
	log.Println(path)
	file, err := os.Open(path)
	if err != nil {
		return http.StatusNotFound, ""
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return http.StatusInternalServerError, ""
	}
	if fileInfo.IsDir() {
		return http.StatusForbidden, ""
	}
	var sb strings.Builder
	_, err = io.Copy(&sb, file)
	if err != nil {
		return http.StatusInternalServerError, ""
	}
	return http.StatusOK, sb.String()
}
