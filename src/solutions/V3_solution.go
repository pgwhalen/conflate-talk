package solutions

import "time"

type Conflater[C any] interface {
	ConflateWith(latest C) C
	ZeroValue() C
}

func ConflateV3[T Conflater[T]](retryInterval time.Duration) (chan<- T, <-chan T) {
	outCh := make(chan T)
	inCh := make(chan T)
	go func() {
		var conflatedMessage T
		var retryTimer *time.Timer
		var retryCh <-chan time.Time
		for {
			select {
			case lastMsg := <-inCh:
				conflatedMessage = conflatedMessage.ConflateWith(lastMsg)
				if retryTimer != nil {
					retryTimer.Stop()
					retryCh = nil
				}
			case <-retryCh:
			}

			select {
			case outCh <- conflatedMessage:
				conflatedMessage = T.ZeroValue(conflatedMessage)
			default:
				retryTimer = time.NewTimer(retryInterval)
				retryCh = retryTimer.C
			}
		}
	}()
	return inCh, outCh
}
