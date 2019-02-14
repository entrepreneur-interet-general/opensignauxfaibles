package testtools

import "dbmongo/lib/engine"

// DiscardOutputChannel
func DiscardOutputChannel(outputChannel chan engine.Tuple) {
	for range outputChannel {
	}
}

func DiscardEventChannel(eventChannel chan engine.Event) {
	for range eventChannel {
	}
}

func EventChannelToSlice(eventChannel chan engine.Event) []engine.Event {
	var events []engine.Event
	for e := range eventChannel {
		events = append(events, e)
	}
	return events
}

func OutputChannelToSlice(outputChannel chan engine.Tuple) []engine.Tuple {
	var data []engine.Tuple
	for d := range outputChannel {
		data = append(data, d)
	}
	return data
}
