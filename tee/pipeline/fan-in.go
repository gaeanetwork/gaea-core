package pipeline

import "sync"

// FanIn multiplex or combine multiple streams into one stream
func FanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})
	multiple := func(c <-chan interface{}) {
		defer wg.Done()

		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiple(c)
	}

	// Wait all read operations over
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}
