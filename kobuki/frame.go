package kobuki

// KobukiData is an interface for Packets in a Frame (e.g. Comand Packets or Feedback Packets)
type KobukiData interface {
	GetName() string
}
