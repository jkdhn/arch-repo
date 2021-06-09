package main

import "os"

func main() {
	if serverCmd.Execute() != nil {
		os.Exit(1)
	}
}
