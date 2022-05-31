// Copyright (C) 2021 The datum-network Authors.

// +build ignore

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var(
	debug          = os.Getenv("BUILDDEBUG") != ""
	goCmd          string
)

func main() {
	log.SetFlags(0)
	goCmd = "go"
	proto()
}

func proto() {

	// go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	// go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	// go get -u github.com/golang/protobuf/protoc-gen-go
	pv := protobufVersion()
	repo := "https://github.com/gogo/protobuf.git"
	path := filepath.Join("repos", "protobuf")
	apipath := filepath.Join("repos", "grpc-gateway")

	runPrint(goCmd, "get", fmt.Sprintf("github.com/gogo/protobuf/protoc-gen-gogofast@%v", pv))
	os.MkdirAll("repos", 0755)

	if _, err := os.Stat(path); err != nil {
		runPrint("git", "clone", repo, path)
	} else {
		runPrintInDir(path, "git", "fetch")
	}
	if _, err := os.Stat(apipath); err != nil {
		runPrint("git", "clone", "-b", "v1.16.0", "https://github.com/grpc-ecosystem/grpc-gateway.git", apipath)
	}

	runPrintInDir(path, "git", "checkout", pv)

	runPrint(goCmd, "generate", "proto/generate.go")
}

func protobufVersion() string {
	bs, err := runError(goCmd, "list", "-f", "{{.Version}}", "-m", "github.com/gogo/protobuf")
	if err != nil {
		log.Fatal("Getting protobuf version:", err)
	}
	return string(bs)
}

func runError(cmd string, args ...string) ([]byte, error) {
	if debug {
		t0 := time.Now()
		log.Println("runError:", cmd, strings.Join(args, " "))
		defer func() {
			log.Println("... in", time.Since(t0))
		}()
	}
	ecmd := exec.Command(cmd, args...)
	bs, err := ecmd.CombinedOutput()
	return bytes.TrimSpace(bs), err
}

func runPrint(cmd string, args ...string) {
	runPrintInDir(".", cmd, args...)
}

func runPrintInDir(dir string, cmd string, args ...string) {
	if debug {
		t0 := time.Now()
		log.Println("runPrint:", cmd, strings.Join(args, " "))
		defer func() {
			log.Println("... in", time.Since(t0))
		}()
	}
	ecmd := exec.Command(cmd, args...)
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr
	ecmd.Dir = dir
	err := ecmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}