# New, better version that is actually maintained
I created this years ago and never updated it. There is a new, better and maintained version by [DownerCase](https://github.com/DownerCase) available [here](https://github.com/DownerCase/ecal-go). 

# golang-ecal
A [high-level interface](https://github.com/Blutkoete/golang-ecal/tree/master/ecal) from [GO](https://golang.org/) to [eCAL](https://github.com/continental/ecal) based on a [low-level interface](https://github.com/Blutkoete/golang-ecal/tree/master/ecalc) generated via [SWIG](http://swig.org/).


## Status
This is the initial version providing an abstraction of the SWIG-generated low-level interface for initialization, publishers and subscribers.

*[ecal](https://github.com/Blutkoete/golang-ecal/tree/master/ecal)*: This is the high-level interface. It allows initialization, data publishing and reception.

*[ecalc](https://github.com/Blutkoete/golang-ecal/tree/master/ecal)*: This is the pure SWIG-generated low-level interface.

As the full low-level interface is accessible, you can do whatever the eCAL C interface allows you to do. More GO-like approaches like channels are only available via the high-level interface and thus are currently limited to publishers and subscribers. The publisher and subscriber interface are complete excluding event callbacks. Services and other functionality are currently only available via the low-level ecalc interface.

## Usage
GO is about simplicity, so the high-level interface initializes a lot of settings with defaults if you do not call the initialization functions yourself.

### Simple Publisher
This publisher corresponds to the *ecal_sample_minimal_snd* example coming with eCAL. It is able to communicate with *ecal_sample_minimal_rec* or with the "Simple Subscriber" example. You can also call
    
    $ golang-ecal minimal_snd
    
to have a simple publisher running.

    package main
    
    import (
        "log"
        "time"
    
        "github.com/Blutkoete/golang-ecal/ecal"
    )
    
    func main() {
	    pub, pubChannel, err := ecal.PublisherCreate("Hello", "base:std::string", "", true)
	    if err != nil {
            log.Fatal(err)
	    }
	    defer pub.Destroy()

	    count := 0
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
    }

### Simple Subscriber
This subscriber corresponds to the *ecal_sample_minimal_rec* example coming with eCAL. It is able to communicate with *ecal_sample_minimal_snd* or with the "Simple Publisher" example. You can also call
    
    $ golang-ecal minimal_rec
    
to have a simple subscriber running.

    package main
    
    import (
        "log"
        "time"
    
        "github.com/Blutkoete/golang-ecal/ecal"
    )
    
    func main() {
	    sub, subChannel, err := ecal.SubscriberCreate("Hello", "base:std::string", "", true, 1024)
	    if err != nil {
            log.Fatal(err)
	    }
	    defer sub.Destroy()

	    for ecal.Ok() {
            select {
            case message := <-subChannel:
                log.Printf("Received \"%s\"\n", message.Content)
            case <-time.After(time.Second):
            }
        }
    }

### Other examples
*golang-ecal_sample* implements two more examples: Sending and receiving [protobuf](https://developers.google.com/protocol-buffers/) messages via eCAL.

    $ golang-ecal person_snd
    
corresponds to *ecal_sample_person_snd* and

    $ golang-ecal person_rec
    
to *ecal_sample_person_rec*. For reading eCAL protobuf messages, it basically boils down to very few lines for unmarshalling the []byte message content:

    ...
    case message := <-subChannel:
    person := &pbexample.Person{}
    proto.Unmarshal(message.Content, person)
    log.Printf("Received \"%s\"\n", person)
    ...
    
See the source of [golang-ecal_sample](https://github.com/Blutkoete/golang-ecal/blob/master/golang-ecal_sample.go) for more details.
    
