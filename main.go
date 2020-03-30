package main

import (
	"flag"
	"fmt"
	"github.com/thanbaiks/vinylstack/core"
	"github.com/thanbaiks/vinylstack/downloader"
	"github.com/thanbaiks/vinylstack/downloader/csn"
	"github.com/thanbaiks/vinylstack/exporter"
	"github.com/thanbaiks/vinylstack/patcher"
	"time"
)

type FactoryWithFlag struct {
	downloader.Factory
	value *string
}

func MustSuccess(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	startTime := time.Now()
	registerFactories()
	var list []FactoryWithFlag
	for _, f := range downloader.Factories() {
		value := flag.String(f.CommandPrefix(), "", f.CommandHelp())
		list = append(list, FactoryWithFlag{
			f,
			value,
		})
		fmt.Printf("++ %s registered\n", f.Name())
	}
	flag.Parse()
	for _, item := range list {
		if len(*item.value) > 0 {
			fmt.Printf("++ [Processing][%s]: %s\n", item.Name(), *item.value)
			d := item.Create(*item.value)
			err := d.Download(&core.DefaultStore)
			if err != nil {
				panic(err)
			}
			fmt.Printf("++ [Processed ][%s]: %s\n", item.Name(), *item.value)
			core.DefaultStore.Dump()
		}
	}
	pl := patcher.LoadPatchList("vinylstack.json")
	if pl != nil {
		pl.Patch(&core.DefaultStore)
	}
	exptr := exporter.NewExporter("_dist_")
	MustSuccess(exptr.Prepare())
	MustSuccess(exptr.DownloadAndExport(&core.DefaultStore))
	fmt.Println("Job finished after", time.Since(startTime))
}

func registerFactories() {
	downloader.RegisterFactory(csn.ChiaSeNhacFactory{})
}
