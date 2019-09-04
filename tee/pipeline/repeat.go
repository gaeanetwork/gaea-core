package pipeline

// RepeatFn repeat fn count to chan
func RepeatFn(done chan interface{}, fn func() interface{}, count int) <-chan interface{} {
	resultStream := make(chan interface{})
	go func() {
		defer close(resultStream)
		for index := 0; index < count; index++ {
			select {
			case <-done:
				return
			case resultStream <- fn():
			}
		}
	}()

	return resultStream
}
