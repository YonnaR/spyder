package main

import (
	"flag"
	"go_spider/core/spider"
	"log"
)

var (
	uri    string
	output string
)

/*
						CLI Initiation
	========================================================
	@uri is crawler start uri
	@out is output file (not required, default is data.json)
	========================================================
*/
func init() {

	flag.StringVar(&uri, "uri", "", "the scrapper start uri adress")
	flag.StringVar(&output, "out", "data.json", "name of the outpout file")

	flag.Parse()
	if uri == "" {
		log.Fatal("you need to set a valid start uri to use this client. use -h for example")
	}
}

/*
					Create process and start engine
==================================================================
	Check connection
	Create new processor for spidering
	Launch engine
==================================================================
*/
func main() {

	if !checkConnection() {
		log.Fatal("can't fetch website, check your internet connection")
	}

	proc := NewTripAdvisorProcessor(uri, output)

	spider.NewSpider(proc, "spyderf").
		AddUrl(proc.startURI, "html").
		//AddPipeline(pipeline.NewPipelineConsole()).            // Print result on screen
		//AddPipeline(pipeline.NewPipelineFile(proc.output)). // Print result in file
		OpenFileLog("./tmp"). // Error info or other useful info in spider will be logged in file of defalt path like "WD/log/log.2014-9-1".
		SetThreadnum(3).      // Crawl request by three Coroutines
		Run()
}
