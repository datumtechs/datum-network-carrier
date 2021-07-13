// Copyright (C) 2021 The RosettaNet Authors.

//+build ignore

package main

import (
"log"
"os"
"os/exec"
"path/filepath"
)

//go:generate go run scripts/protofmt.go .

// First generate extensions using standard proto compiler.
//go:generate protoc -I ../ -I . --gogofast_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor,paths=source_relative:ext ext.proto

// Then build our vanity compiler that uses the new extensions
//go:generate go build -o scripts/protoc-gen-gorosetta scripts/protoc_plugin.go

// Inception, go generate calls the script itself that then deals with generation.
// This is only done because go:generate does not support wildcards in paths.
//go:generate go run generate.go lib/db lib/types lib/api lib/center/api lib/p2p/v1 lib/rpc/v1 lib/consensus/twopc lib/fighter/computesvc lib/fighter/datasvc lib/fighter/common

// final, generate grpc gateway stub.
//go:generate protoc -I ../ -I . -I ../repos/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:./.. ./lib/rpc/v1/debug.proto
//go:generate protoc -I ../ -I . -I ../repos/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:./.. ./lib/rpc/v1/debug.proto

// go:generate protoc -I ../ -I . -I ../repos/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:./.. ./lib/api/*.proto
// go:generate protoc -I ../ -I . -I ../repos/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:./.. ./lib/api/*.proto

func main() {
	for _, path := range os.Args[1:] {
		matches, err := filepath.Glob(filepath.Join(path, "*proto"))
		if err != nil {
			log.Fatal(err)
		}
		log.Println(path, "returned:", matches)
		args := []string{
			"-I", "..",
			"-I", ".",
			"-I", "../repos/grpc-gateway/third_party/googleapis",
			//"--plugin=protoc-gen-gorosetta=scripts/protoc-gen-gorosetta",
			//"--gorosetta_out=paths=source_relative:..",
			"--gogofast_out=plugins=grpc,paths=source_relative:..",
		}
		args = append(args, matches...)
		cmd := exec.Command("protoc", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal("Failed generating", path)
		}
	}
}
