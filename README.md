# golang-ecal
A [high-level interface](https://github.com/Blutkoete/golang-ecal/tree/master/ecal) from [GO](https://golang.org/) to [eCAL](https://github.com/continental/ecal) based on a [low-level interface](https://github.com/Blutkoete/golang-ecal/tree/master/ecalc) generated via [SWIG](http://swig.org/).


## Status
This is the initial version providing an abstraction of the SWIG-generated low-level interface for initialization, publishers and subscribers.

*[ecal](https://github.com/Blutkoete/golang-ecal/tree/master/ecal)*: This is the high-level interface. It allows initialization, data publishing and reception.

*[ecalc](https://github.com/Blutkoete/golang-ecal/tree/master/ecal)*: This is the pure SWIG-generated low-level interface.

As the full low-level interface is accessible, you can do whatever the eCAL C interface allows you to do. More GO-like approaches like channels are only available via the high-level interface and thus are currently limited to publishers and subscribers. The publisher and subscriber interface are complete excluding event callbacks. Services and other functionality are currently only available via the low-level ecalc interface.

## Usage
GO is about simplicity, so the high-level interface initializes a lot of settings with defaults if you do not call the initialization functions yourself.

Sending:

    package main
    
    import (
        "log"
        "time"
    
        "github.com/Blutkoete/golang-ecal/ecal"
    )
    
    func main() {
    	var pub ecal.PublisherIf
	    var pubChannel chan<- ecal.Message
        var err error
        pub, pubChannel, err = ecal.PublisherCreate("Hello", "base:std::string", "", true)
        if err != nil {
            log.Fatal(err)
        }
        defer pub.Destroy()
    
        go func() {
            count := 1
            for ecal.Ok() {
                message := ecal.Message{Content: []byte(fmt.Sprintf("Hello World from Go (%d)", count)),
                    Timestamp: -1}
                count += 1
                select {
                case pubChannel <- message:
                    log.Printf("Sent \"%s\"\n", message.Content)
                case <-time.After(time.Second):
                }
                <-time.After(250 * time.Millisecond)
            }
    
            pub.Stop()
        }()
    
        for !pub.IsStopped() {
            <-time.After(time.Second)
        }
    
        <-time.After(time.Second)
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
        var sub ecal.SubscriberIf
        var subChannel <-chan ecal.Message
        var err error
        sub, subChannel, err = ecal.SubscriberCreate("Hello", "base:std::string", "", true, 1024)
        if err != nil {
            log.Fatal(err)
        }
        defer sub.Destroy()
    
        go func() {
            for ecal.Ok() {
                select {
                case message := <-subChannel:
                    log.Printf("Received \"%s\"\n", message.Content)
                case <-time.After(time.Second):
                }
            }
    
            sub.Stop()
        }()
    
        for !sub.IsStopped() {
            <-time.After(time.Second)
        }
    
    }
