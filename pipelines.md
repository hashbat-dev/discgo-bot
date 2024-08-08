# Constructing Pipelines

Channels are uniquely suited to constructing pipelines in Go because they fulfill all our basic requirements. They can receive and emit values, they can safely be used concurrently, they can be ranged over, and they are reified by the language. Looking at using channels:

```golang
generator := func(done <-chan interface{}, integers ...int) <-chan int {
   intStream := make(chan int) 
   go func() {
    defer close(intStream)
    for _, i := range integers {
        select {
            case <- done:
                return
            case intStream <-i:
        }
    }
   }()
   return intstream
}

multiply:= func(
    done <-chan interface{},
    intStream <-chan int,
    multiplier int,
) <-chan int{
    multipliedStream := make(chan int)
    go func() {
        defer close(multipliedStream)
        for i := range intStream {
            select {
                case <-done:
                    return
                case multipliedStream <- i* multiplier:
            }
        }
    }()
    return multipliedStream
}

add := func(
    done <-chan interface{},
    intStream <-chan int,
    additive int,
) <-chan int {
   addedStream := make(chan int) 
   go func() {
    defer close(addedStream)
    for i := range intStream {
        select {
            case <-done:
                return
            case addedStream <- i+additive:
        }
    }
   }()
   return addedStream
}

done := make(chan interface{})
defer close(done)
// the done channel is passed to each as this communicates our signal to kill the processing, and it bubbles up through the pipeline to each one
intStream := generator (done, 1, 2, 3, 4)
pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

// Notice how the <-chan int intStream is passed to multiply, then multiply passes the same data type to add, and add passes back to multiply



for v := range pipeline {
    fmt.Println(v)
}
```

The above outputs
6
10
14
18
To see in an example:
https://go.dev/play/p/BURKnwtvdks

On first glance this might seem like a lot of code to achieve something simple, but the benefits become apparent once we start looking at concurrent structures.

## Fan-Out, Fan-in
Some stages in processing can be computationally expensive, a good example is the functionality we've written for handling image manipulation, particulary for processing gifs when we might want to reverse, or mirror them.

This is a good candidate for the fan-out fan-in pattern. First though, we have to question whether it meets two criteria, 

- Does it rely on values that the stage had calculated before?
- Does it take a long time to run?

Order dependence - frames are ordered in a gif, and rows of pixels have a particular place in the vertical order. These pieces of information however, can be captured when spinning up processes, as an index.
Speed - we know from practise that these things can be computationally expensive, taking up to 30-50 seconds at peak for larger images to apply transformations.

Fanning out, fortunately, is extremely straight-forward to write:
```golang
   // using the earlier functions... 
   numMultipliers := runtime.NumCPU()
   intStream := generator (done, 1, 2, 3, 4)
   multipliers := make([]<-chan int, numMultipliers)
   for i := 0; i < numMultipliers; i++ {
       multipliers[i] = multiply(done, intStream, i)
   }
```

Here it's fairly plan that what we do is iterate over a series of values that we need to use, and then spawn channels with goRoutines that handle our processing. To actually do something with this (get the processed values back) we need to then write the fan-in portion..

```golang
fanIn := func(
    done <- chan interface{},
    channels ...<-chan interface {},
) <-chan interface{} {
    var wg sync.WaitGroup
    multiplexedStream := make(chan interface{}) 

    multiplex := func(c <-chan interface{}){
        defer wg.Done()
        for i := range c{
            select {
                case <-done:
                    return
                case multiplexedStream <-i:
            }
        }
    }

    // Select from all the channels
    wg.Add(len(channels))
    for _, c := range channels {
       go multiplex(c)
    }

    // Wait for all the reads to complete
    go func() {
        wg.Wait()
        close(multiplexedStream)
    }()
    return multiplexedStream
}
```

The value of multiplexedStream as type `<-chan interface{}` can effectively then be consumed the same way prior pipelines were. If we needed to do something with these to re-organise data (such as in the case of image frames and row indices) we can do this at the point at which we consume this.
