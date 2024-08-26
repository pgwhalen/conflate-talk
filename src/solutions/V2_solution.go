package solutions

import (
	"time"
)

func ConflateV2[T any](retryInterval time.Duration) (chan<- T, <-chan T) {
	outCh := make(chan T)
	inCh := make(chan T)
	go func() {
		var lastMsg T
		var retryTimer *time.Timer
		var retryCh <-chan time.Time
		for {
			select {
			case lastMsg = <-inCh:
				if retryTimer != nil {
					retryTimer.Stop()
					retryCh = nil
				}
			case <-retryCh:
			}

			select {
			case outCh <- lastMsg:
			default:
				retryTimer = time.NewTimer(retryInterval)
				retryCh = retryTimer.C
			}
		}
	}()
	return inCh, outCh
}
