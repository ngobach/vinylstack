package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	parse("Fav~Y3NuX3BsYXlsaXN0fjE0NzI=")
	return
	//
	var outfile string
	flag.StringVar(&outfile, "out", "./out.json", "Output file")
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("A playlist url must be provided!")
		os.Exit(1)
	}
}
