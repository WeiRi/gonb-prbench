package main

var stopWatchChan = make(chan bool, 1)

func SendToWatch() {
	stopWatchChan <- true
}

func CloseWatch() {
	close(stopWatchChan)
}
