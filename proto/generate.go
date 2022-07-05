// Copyright (C) 2021 The datum-network Authors.

//+build ignore

package main

import (
"log"
"os"
"os/exec"
"path/filepath"
)

// go:generate go run scripts/protofmt.go .

// First generate extensions using standard proto compiler.
//go:generate protoc -I ../ -I . --gogofast_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor,paths=source_relative:ext ext.proto

// Then build our vanity compiler that uses the new extensions
//go:generate go build -o scripts/protoc-gen-gorosetta scripts/protoc_plugin.go

// Inception, go generate calls the script itself that then deals with generation.
// This is only done because go:generate does not support wildcards in paths.
//go:generate go run generate.go armada-common/carrier/api armada-common/carrier/db armada-common/carrier/netmsg/common armada-common/carrier/netmsg/consensus/twopc armada-common/carrier/netmsg/taskmng armada-common/carrier/p2p/v1  armada-common/carrier/rpc/debug/v1 armada-common/carrier/types armada-common/common/constant armada-common/datacenter/api armada-common/fighter/api/compute armada-common/fighter/api/data armada-common/fighter/types

// final, generate grpc gateway stub.
//go:generate protoc -I ../ -I ./armada-common -I . -I ../repos/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:./../pb ./armada-common/carrier/rpc/debug/v1/debug.proto
//go:generate protoc -I ../ -I ./armada-common -I . -I ../repos/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:./../pb ./armada-common/carrier/rpc/debug/v1/debug.proto

//go:generate protoc -I ../ -I ./armada-common -I . -I ../repos/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:./../pb ./armada-common/carrier/api/*.proto
//go:generate protoc -I ../ -I ./armada-common -I . -I ../repos/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:./../pb ./armada-common/carrier/api/*.proto


func main() {
	for _, path := range os.Args[1:] {
		matches, err := filepath.Glob(filepath.Join(path, "*proto"))
		if err != nil {
			log.Fatal(err)
		}
		log.Println(path, "returned:", matches)
		args := []string{
			"-I", "..",
			"-I", "./armada-common",
			"-I", ".",
			"-I", "../repos/grpc-gateway/third_party/googleapis",
			//"--plugin=protoc-gen-gorosetta=scripts/protoc-gen-gorosetta",
			//"--gorosetta_out=paths=source_relative:..",
			"--gogofast_out=plugins=grpc,paths=source_relative:../pb",
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
