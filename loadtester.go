package main

import (
	"fmt"
	"os"
    "net/http"
    "time"
)

//Number of individual threads to make requests
const NUMTHREADS = 100
const DEFAULT_URL = "http://localhost:5050"

func main () {

	cmdArguments := os.Args[1:]
	for _, element := range cmdArguments {
		if element == "-c" {
			fmt.Println("we have found an option!")
		}
	}

	url := DEFAULT_URL
    //Channel for threads to send their response times
    timeChannel := make (chan time.Duration)
    //spawn threads
    for i := 0; i < NUMTHREADS; i++ {
        go makeResponse(timeChannel, url)
    }
    var averageTime time.Duration 
    //receive response times
    for i := 0; i < NUMTHREADS; i++ {
        select {
            case requestTime := <-timeChannel:
                averageTime += requestTime
        }
    }
    averageTime /= NUMTHREADS
    fmt.Println("Average time was", averageTime)

}

func makeResponse (timeChannel chan time.Duration, url string) {
    timeStart := time.Now()
    resp, err := http.Get(url)
    requestTime := time.Since(timeStart)
    //Close http request after function is completed
	defer resp.Body.Close()

    if err != nil {
        fmt.Println("error with request", err)
        return
    }

    fmt.Println("Made request in ", requestTime)
    //send response time on channel
    timeChannel <- requestTime
    return
}