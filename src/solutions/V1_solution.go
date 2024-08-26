package solutions

func ConflateV1[T any]() (chan<- T, <-chan T) {
	outCh := make(chan T)
	inCh := make(chan T)
	go func() {
		for {
			lastMsg := <-inCh

			select {
			case outCh <- lastMsg:
			default:
			}
		}
	}()
	return inCh, outCh
}
