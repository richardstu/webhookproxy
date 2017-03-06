package semaphore

import (
	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhook-proxy-go/src/util"
)

func ExampleMessage(data util.ExampleData, eventSlug string) (cc.Message, error) {
	bytes, err := data.ExampleMessageBytes(HandlerKey, eventSlug)
	if err != nil {
		return cc.Message{}, err
	}
	return Normalize(bytes)
}