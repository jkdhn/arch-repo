package main

import "os"

func main() {
	if deployCmd.Execute() != nil {
		os.Exit(1)
	}
}
