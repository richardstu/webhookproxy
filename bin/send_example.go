package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	//log "github.com/Sirupsen/logrus"
	"github.com/grokify/glip-go-webhook"
	"github.com/grokify/glip-webhook-proxy-go/src/config"
	"github.com/grokify/glip-webhook-proxy-go/src/handlers/appsignal"
	"github.com/grokify/glip-webhook-proxy-go/src/handlers/confluence"
	"github.com/grokify/glip-webhook-proxy-go/src/handlers/enchant"
	"github.com/grokify/glip-webhook-proxy-go/src/handlers/heroku"
	"github.com/grokify/glip-webhook-proxy-go/src/handlers/magnumci"
	"github.com/grokify/glip-webhook-proxy-go/src/handlers/raygun"
	"github.com/grokify/glip-webhook-proxy-go/src/handlers/semaphoreci"
	"github.com/grokify/glip-webhook-proxy-go/src/handlers/travisci"
	"github.com/grokify/glip-webhook-proxy-go/src/handlers/userlike"
	"github.com/grokify/glip-webhook-proxy-go/src/util"
)

const (
	GLIP_WEBHOOK_ENV = "GLIP_WEBHOOK"
)

func main() {
	guidPointer := flag.String("guid", "", "Glip webhook GUID or URL")
	examplePointer := flag.String("example", "", "Example message type")
	flag.Parse()
	guid := strings.TrimSpace(*guidPointer)
	example := strings.ToLower(strings.TrimSpace(*examplePointer))

	fmt.Printf("LENGUID[%v]\n", len(guid))
	if len(guid) < 1 {
		guid = os.Getenv(GLIP_WEBHOOK_ENV)
		fmt.Printf("HERE [%v]\n", guid)
	}

	//glip, _ := glipwebhook.NewGlipWebhookClient(guid)
	fmt.Printf("GUID [%v]\n", guid)
	fmt.Printf("EXAMPLE [%v]\n", example)

	if len(example) < 1 {
		panic("Usage: send_example.go -hook=<GUID> -example=raygun")
	}

	glipClient, err := glipwebhook.NewGlipWebhookClient(guid)
	if err != nil {
		panic("Incorrect Webhook GUID or URL")
	}

	config.GLIP_ACTIVITY_INCLUDE_INTEGRATION_NAME = true

	switch example {
	case "appsignal":
		SendAppsignal(glipClient, guid)
	case "confluence":
		SendConfluence(glipClient, guid)
	case "enchant":
		glipMsg, err := enchant.ExampleMessageGlip()
		if err != nil {
			panic("Bad Test Message")
		}
		util.SendGlipWebhook(glipClient, guid, glipMsg)
	case "heroku":
		glipMsg, err := heroku.ExampleMessageGlip()
		if err != nil {
			panic("Bad Test Message")
		}
		util.SendGlipWebhook(glipClient, guid, glipMsg)
	case "magnumci":
		glipMsg, err := magnumci.ExampleMessageGlip()
		if err != nil {
			panic(fmt.Sprintf("Bad Test Message [%v]", err))
		}
		util.SendGlipWebhook(glipClient, guid, glipMsg)
	case "raygun":
		glipMsg, err := raygun.ExampleMessageGlip()
		if err != nil {
			panic("Bad Test Message")
		}
		util.SendGlipWebhook(glipClient, guid, glipMsg)
	case "semaphoreci":
		SendSemaphoreci(glipClient, guid)
	case "travisci":
		glipMsg, err := travisci.ExampleMessageGlip()
		if err != nil {
			panic("Bad Test Message")
		}
		util.SendGlipWebhook(glipClient, guid, glipMsg)
	case "userlike":
		SendUserlike(glipClient, guid)
	default:
		panic(fmt.Sprintf("Unknown webhook source %v\n", example))
	}
}

func SendAppsignal(glipClient glipwebhook.GlipWebhookClient, guid string) {
	glipMsg, err := appsignal.ExampleMarkerMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
	glipMsg, err = appsignal.ExampleExceptionMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
	glipMsg, err = appsignal.ExamplePerformanceMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
}

func SendConfluence(glipClient glipwebhook.GlipWebhookClient, guid string) {
	glipMsg, err := confluence.ExamplePageCreatedMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
	glipMsg, err = confluence.ExampleCommentCreatedMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
}

func SendSemaphoreci(glipClient glipwebhook.GlipWebhookClient, guid string) {
	glipMsg, err := semaphoreci.ExampleBuildMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
	glipMsg, err = semaphoreci.ExampleDeployMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
}

func SendUserlike(glipClient glipwebhook.GlipWebhookClient, guid string) {
	glipMsg, err := userlike.ExampleOfflineMessageReceiveMessageGlip()
	if err != nil {
		panic(fmt.Sprintf("Bad Test Message [%v]", err))
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)

	glipMsg, err = userlike.ExampleChatWidgetConfigMessageGlip()
	if err != nil {
		panic("Bad Test Message")
	}
	util.SendGlipWebhook(glipClient, guid, glipMsg)
	//return
	if 1 == 1 {
		for i, event := range userlike.ChatMetaEvents {
			//continue
			fmt.Printf("%v %v\n", i, event)
			glipMsg, err := userlike.ExampleUserlikeChatMetaMessageGlip(event)
			if err != nil {
				panic(fmt.Sprintf("Bad Test Message: %v", err))
			}
			util.SendGlipWebhook(glipClient, guid, glipMsg)
		}
	}
	for i, event := range userlike.OperatorEvents {
		fmt.Printf("%v %v\n", i, event)
		glipMsg, err := userlike.ExampleUserlikeOperatorMessageGlip(event)
		if err != nil {
			panic(fmt.Sprintf("Bad Test Message: %v", err))
		}
		util.SendGlipWebhook(glipClient, guid, glipMsg)
	}
}
