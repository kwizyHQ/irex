package main

import (
	"flag"
	"fmt"
)

func main() {
	analyze := flag.Bool("a", false, "Run Analyze")
	profile := flag.Bool("p", false, "Run Profile")
	watcher := flag.Bool("w", false, "Run Watcher")
	proxy := flag.String("x", "", "Run Proxy Command")

	flag.Parse()

	if *analyze {
		fmt.Println("Running Analyze...")
		Analyze()
	} else if *profile {
		fmt.Println("Running Profile...")
		Profile()
	} else if *watcher {
		fmt.Println("Running Watcher...")
		WatchAndBuild()
	} else if *proxy != "" {
		fmt.Println("Running Proxy Command...")
		RunProxyCommand(*proxy, flag.Args())
	}
}
