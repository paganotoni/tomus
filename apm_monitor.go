package tomus

//APMMonitor is an interface for each of our APM monitors
type APMMonitor interface {

	//Track is useful for background operations or functions that need to be measured
	Track(name string, fn func() error) error

	//Listen is the method that will start listening for buffalo events to create/finish spans.
	Listen() error
}
