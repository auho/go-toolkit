package client

import "time"

var (
	ContentTypePlain = "text/plain"
	ContentTypeJson  = "application/json"

	notifyStateShutdown = "shutDown"
	notifyStateContinue = "continue"
)

func isNotifyStateShutDown(s string) bool {
	return s == notifyStateShutdown
}

func isNotifyStateContinue(s string) bool {
	return s == notifyStateContinue
}

// backoff 按照时间避退
// 按照起始时间 b，每次增加 2 倍，直到最大时间 mb
func backoff(b time.Duration, mb time.Duration) time.Duration {
	if b >= mb {
		return mb
	}

	time.Sleep(b)

	return b * 2
}
