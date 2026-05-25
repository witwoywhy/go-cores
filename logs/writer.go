package logs

import (
	"encoding/json"
	"log"

	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/pubsub"
)

type stdoutWriter struct{}

func (w *stdoutWriter) Write(fields map[string]any, l logger.Logger) {
	b, err := json.Marshal(fields)
	if err == nil {
		log.Println(string(b))
	}
}

type producerWriter struct {
	producer pubsub.Producer
}

func (w *producerWriter) Write(fields map[string]any, l logger.Logger) {
	w.producer.Produce(fields[apps.TraceID].(string), fields, l)
}
