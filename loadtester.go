package main

import (
    "fmt"
    "net/http"
    "time"
)

const url = "http://localhost:5050"
//Number of individual threads to make requests
const NUMTHREADS = 100

func main () {
    //Channel for threads to send their response times
    timeChannel := make (chan time.Duration)
    //spawn threads
    for i := 0; i < NUMTHREADS; i++ {
        go makeResponse(timeChannel)
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

func makeResponse (timeChannel chan time.Duration) {
    timeStart := time.Now()
    resp, err := http.Get(url)
    requestTime := time.Since(timeStart)

    if err != nil {
        fmt.Println("error with request", err)
        return
    }

    //Close http request after function is completed
    defer resp.Body.Close()

    if err != nil {
        fmt.Println("error with reading request")
        return
    }
    fmt.Println("Made request in ", requestTime)
    //send response time on channel
    timeChannel <- requestTime

    return
}