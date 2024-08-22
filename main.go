package main

import (
	"awesomeProject/archiver"
	"fmt"
	"github.com/pkg/errors"
)

const httpAddress = "https://testapi.qubic.org"
const grpcAddress = "213.170.135.5:8003"

func main() {

	err := Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func Run() error {

	err := archiver.RunGRPCExample()
	if err != nil {
		return errors.Wrap(err, "running grpc example")
	}

	return nil
}
