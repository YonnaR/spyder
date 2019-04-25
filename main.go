package main

import (
	"flag"
	"go_spider/core/pipeline"
	"go_spider/core/spider"
	"log"
)

var (
	uri    string
	output string
	thread int
)

func init() {
	flag.StringVar(&uri, "uri", "", "the scrapper start uri adress")
	flag.StringVar(&output, "out", "data.json", "name of the outpout file")
	flag.IntVar(&thread, "thread", 3, "thread authorized")

	flag.Parse()
	if uri == "" {
		log.Fatal("you need to set a valid start uri for scrapping. use -h for example")
	}
}

func main() {

	proc := NewTripAdvisorProcessor(uri, output, thread)
	if !checkConnection() {
		log.Fatal("can't fetch website, check your internet connection")
	}
	spider.NewSpider(proc, "scrappeyrf").
		AddUrl(proc.startURI, "html").
		//AddPipeline(pipeline.NewPipelineConsole()).            // Print result on screen
		AddPipeline(pipeline.NewPipelineFile(proc.output)). // Print result in file
		OpenFileLog("./tmp").                               // Error info or other useful info in spider will be logged in file of defalt path like "WD/log/log.2014-9-1".
		SetThreadnum(3).                                    // Crawl request by three Coroutines
		Run()
}
