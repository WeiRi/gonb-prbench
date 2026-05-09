// Production stub for VictoriaMetrics app/vmalert/notifier/alertmanager.go (PR #8258).
// Pre-PR: send creates alertsToSend := alerts[:0] sharing backing array;
// concurrent send() calls with same alerts slice race on element slots.
package notifier

type Alert struct {
	Name  string
	Value float64
}

func applyRelabelingIfNeeded(a Alert) []string { return []string{a.Name} }

// send mirrors the BUG-state with aliased backing array.
func send(alerts []Alert) []Alert {
	alertsToSend := alerts[:0] // RACE: aliased backing array
	for _, a := range alerts {
		lbls := applyRelabelingIfNeeded(a)
		if lbls == nil {
			continue
		}
		alertsToSend = append(alertsToSend, a) // RACE: write to shared backing slot
	}
	return alertsToSend
}
