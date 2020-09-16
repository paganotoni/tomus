package tomus

// APMMonitor is an interface for the shared monitors in this package or
// subpackages, APM monitors listen to Buffalo events and report that
// to specific platforms like Newrelic or Datadog.
type APMMonitor interface {
	// Listen would make the APM monitor to starts to listen for Buffalo events.
	Listen() error

	// Track is useful for things like BackgroundTransactions
	// that cannot be tracked with events.
	Track(string, func() error) error
}
