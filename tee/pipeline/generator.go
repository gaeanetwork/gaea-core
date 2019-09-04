package pipeline

// Generator generator from target value
func Generator(done chan interface{}, target ...interface{}) <-chan interface{} {
	outputStream := make(chan interface{})
	go func() {
		defer close(outputStream)
		for _, value := range target {
			select {
			case <-done:
				return
			case outputStream <- value:
			}

		}
	}()

	return outputStream
}
