package kafka

type Message struct {
	Key, Value string
	Topic      string
	Partition  int32
	Offset     int64
}
