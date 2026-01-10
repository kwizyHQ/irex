package main

import (
	"flag"
	"fmt"
)

func main() {
	analyze := flag.Bool("a", false, "Run Analyze")
	profile := flag.Bool("p", false, "Run Profile")
	flag.Parse()

	if *analyze {
		fmt.Println("Running Analyze...")
		Analyze()
	} else if *profile {
		fmt.Println("Running Profile...")
		Profile()
	} else {
		fmt.Println("No action specified. Use -a for analyze or -p for profile.")
	}
}
