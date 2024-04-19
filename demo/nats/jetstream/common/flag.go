package common

import "flag"

var (
	FlagIsProducer = flag.Bool("producer", false, "Set to true to run as a producer")
	FlagNatsAddr   = flag.String("nats", "nats://nats1:4222", "NATS server address")
	FlagNumProcess = flag.Int("num", 10, "Number of concurrent processes")
)
