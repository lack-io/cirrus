package daemon

type Daemon interface {
	StartDaemon(root string)
	PauseDaemon()
}
