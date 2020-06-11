# golang-ecal
A [high-level interface](https://github.com/Blutkoete/golang-ecal/tree/master/ecal) from [GO](https://golang.org/) to [eCAL](https://github.com/continental/ecal) based on a [low-level interface](https://github.com/Blutkoete/golang-ecal/tree/master/ecalc) generated via [SWIG](http://swig.org/).


## Status
This is the initial proof-of-concept version providing an abstraction of the SWIG-generated low-level interface for initialization, publishers and subscribers.

*[ecal](https://github.com/Blutkoete/golang-ecal/tree/master/ecal)*: This is the high-level interface. It allows initialization, data publishing and reception.

*[ecalc](https://github.com/Blutkoete/golang-ecal/tree/master/ecal)*: This is the pure SWIG-generated low-level interface.

As the full low-level interface is accessible, you can do whatever the eCAL C interface allows you to do. More GO-like approaches like channels are only available via the high-level interface and thus are currently limited to publishers and subscribers. The publisher and subscriber interface is also not complete, but as you can retrieve the handle from a high-level publisher, you could use the low-level interface to set whatever other settings you want.

## Usage

Sending:

    package main
    
    import (
        "log"
        "os"
    
        "github.com/Blutkoete/golang-ecal/ecal"
    )
    
    func main() {
    
        err := ecal.Initialize(os.Args, "golang_ecal_sample", ecal.InitDefault)
        if err != nil {
            log.Fatal(err)
        }
        defer ecal.Finalize(ecal.InitAll)
        log.Println("Initialized successfully.")
    
        var publisher ecal.PublisherIf
        publisher, err = ecal.PublisherCreate("foo", "base:std::string", "", true)
        if err != nil {
            log.Fatal(err)
        }
        defer publisher.Destroy()
    
        publisherInput := publisher.GetChannel()
        message := ecal.Message { Content:   []byte("HELLO WORLD FROM GO"),
                                  Timestamp: -1 }
        for ecal.Ok() {
            log.Println("sending", len(message.Content), "bytes")
            publisherInput <- message
            ecal.ProcessSleepMS(250)
        }
        
    }

Receiving:

    package main
    
    import (
        "log"
        "os"
        "time"
    
        "github.com/Blutkoete/golang-ecal/ecal"
    )
    
    func main() {
    
        err := ecal.Initialize(os.Args, "golang_ecal_sample", ecal.InitDefault)
        if err != nil {
            log.Fatal(err)
        }
        defer ecal.Finalize(ecal.InitAll)
        log.Println("Initialized successfully.")
    
        var subscriber ecal.SubscriberIf
        subscriber, err = ecal.SubscriberCreate("Hello", "base:std::string", "", true, 1024)
        if err != nil {
            log.Fatal(err)
        }
        defer subscriber.Destroy()
    
        receiveChannel := subscriber.GetChannel()
        timeout := time.NewTimer(1 * time.Second)
    
        for ecal.Ok() {
            log.Println("Waiting for message ...")
            select {
            case message := <- receiveChannel:
                timeout.Reset(1 * time.Second)
                log.Printf("Received a message with %d bytes: \"%s\"", len(message.Content), string(message.Content))
            case <-timeout.C:
                timeout.Reset(1 * time.Second)
                continue
            }
        }
    
    }
