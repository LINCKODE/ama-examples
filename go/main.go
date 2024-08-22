package main

import (
	"fmt"
)

func main() {

	err := Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func Run() error {

	/*err := archiver.RunGRPCExample()
	if err != nil {
		return errors.Wrap(err, "running grpc example")
	}*/

	/*err := archiver.RunHTTPExample()
	if err != nil {
		return errors.Wrap(err, "running http example")
	}*/

	/*err := http.CreateAndBroadcastTransaction()
	if err != nil {
		return errors.Wrap(err, "running transaction example")
	}*/

	return nil
}
