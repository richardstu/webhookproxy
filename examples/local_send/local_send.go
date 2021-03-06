package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/grokify/gotilla/fmt/fmtutil"

	log "github.com/Sirupsen/logrus"
	cc "github.com/commonchat/commonchat-go"
	"github.com/grokify/webhookproxy/src/adapters"
	"github.com/grokify/webhookproxy/src/config"
	"github.com/grokify/webhookproxy/src/util"
	"github.com/valyala/fasthttp"

	"github.com/grokify/webhookproxy/src/handlers/appsignal"
	"github.com/grokify/webhookproxy/src/handlers/apteligent"
	"github.com/grokify/webhookproxy/src/handlers/circleci"
	"github.com/grokify/webhookproxy/src/handlers/codeship"
	"github.com/grokify/webhookproxy/src/handlers/confluence"
	"github.com/grokify/webhookproxy/src/handlers/datadog"
	"github.com/grokify/webhookproxy/src/handlers/deskdotcom"
	"github.com/grokify/webhookproxy/src/handlers/enchant"
	"github.com/grokify/webhookproxy/src/handlers/gosquared"
	"github.com/grokify/webhookproxy/src/handlers/gosquared2"
	"github.com/grokify/webhookproxy/src/handlers/heroku"
	"github.com/grokify/webhookproxy/src/handlers/librato"
	"github.com/grokify/webhookproxy/src/handlers/magnumci"
	"github.com/grokify/webhookproxy/src/handlers/marketo"
	"github.com/grokify/webhookproxy/src/handlers/opsgenie"
	"github.com/grokify/webhookproxy/src/handlers/papertrail"
	"github.com/grokify/webhookproxy/src/handlers/pingdom"
	"github.com/grokify/webhookproxy/src/handlers/raygun"
	"github.com/grokify/webhookproxy/src/handlers/runscope"
	"github.com/grokify/webhookproxy/src/handlers/semaphore"
	"github.com/grokify/webhookproxy/src/handlers/statuspage"
	"github.com/grokify/webhookproxy/src/handlers/travisci"
	"github.com/grokify/webhookproxy/src/handlers/userlike"
	"github.com/grokify/webhookproxy/src/handlers/victorops"
)

const (
	GLIP_WEBHOOK_ENV  = "GLIP_WEBHOOK"
	SLACK_WEBHOOK_ENV = "SLACK_WEBHOOK"
)

type Sender struct {
	Adapter adapters.Adapter
}

func (sender *Sender) SendCcMessage(ccMsg cc.Message, err error) {
	if err != nil {
		panic(fmt.Sprintf("Bad Test Message: %v\n", err))
	}
	req, resp, err := sender.Adapter.SendMessage(ccMsg)
	fmt.Printf("RESPONSE_STATUS_CODE [%v]\n", resp.StatusCode())
	if err != nil {
		fmt.Printf("ERROR [%v]\n", err)
	}
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}

func main() {
	log.SetLevel(log.DebugLevel)

	guidPointer := flag.String("guid", "", "Glip webhook GUID or URL")
	examplePointer := flag.String("example", "", "Example message type")
	adapterType := flag.String("adapter", "", "Adapter")

	flag.Parse()
	webhookURLOrUID := strings.TrimSpace(*guidPointer)
	example := strings.ToLower(strings.TrimSpace(*examplePointer))

	fmt.Printf("LENGUID[%v]\n", len(webhookURLOrUID))
	fmt.Printf("GUID [%v]\n", webhookURLOrUID)
	fmt.Printf("EXAMPLE [%v]\n", example)

	if len(example) < 1 {
		panic("Usage: send_example.go -hook=<GUID> -adapter=glip -example=raygun")
	}

	cfg := config.Configuration{
		IconBaseURL:    "https://grokify.github.io/webhookproxy/icons/",
		LogrusLogLevel: log.DebugLevel}

	sender := Sender{}
	if *adapterType == "glip" {
		if len(webhookURLOrUID) < 1 {
			webhookURLOrUID = os.Getenv(GLIP_WEBHOOK_ENV)
			fmt.Printf("GLIP_GUID_ENV [%v]\n", webhookURLOrUID)
		}
		adapter, err := adapters.NewGlipAdapter(webhookURLOrUID)
		if err != nil {
			panic("Incorrect Webhook GUID or URL")
		}
		sender.Adapter = adapter
	} else if *adapterType == "slack" {
		if len(webhookURLOrUID) < 1 {
			webhookURLOrUID = os.Getenv(SLACK_WEBHOOK_ENV)
			fmt.Printf("SLACK_GUID_ENV [%v]\n", webhookURLOrUID)
		}
		adapter, err := adapters.NewSlackAdapter(webhookURLOrUID)
		if err != nil {
			panic("Incorrect Webhook GUID or URL")
		}
		sender.Adapter = adapter
	} else {
		panic("Invalid Adapter")
	}

	exampleData, err := util.NewExampleData()
	if err != nil {
		panic(fmt.Sprintf("Invalid Example Data: %v\n", err))
	}
	fmtutil.PrintJSON(exampleData)

	switch example {
	case "appsignal":
		source := exampleData.Data[appsignal.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(appsignal.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "apteligent":
		source := exampleData.Data[apteligent.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(apteligent.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "circleci":
		sender.SendCcMessage(circleci.ExampleMessage(cfg, exampleData))
	case "codeship":
		sender.SendCcMessage(codeship.ExampleMessage(cfg, exampleData))
	case "confluence":
		source := exampleData.Data[confluence.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(confluence.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "datadog":
		sender.SendCcMessage(datadog.ExampleMessage(cfg, exampleData))
	case "deskdotcom":
		source := exampleData.Data[deskdotcom.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(deskdotcom.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "enchant":
		sender.SendCcMessage(enchant.ExampleMessage(cfg, exampleData))
	case "gosquared":
		source := exampleData.Data[gosquared.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(gosquared.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "gosquared2":
		source := exampleData.Data[gosquared.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(gosquared2.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "heroku":
		sender.SendCcMessage(heroku.ExampleMessage(cfg, exampleData))
	case "librato":
		source := exampleData.Data[librato.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(librato.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "magnumci":
		sender.SendCcMessage(magnumci.ExampleMessage(cfg, exampleData))
	case "marketo":
		source := exampleData.Data[marketo.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(marketo.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "opsgenie":
		source := exampleData.Data[opsgenie.HandlerKey]
		for i, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(opsgenie.ExampleMessage(cfg, exampleData, eventSlug))
			if i == 8 {
				time.Sleep(2000 * time.Millisecond)
			}
		}
	case "papertrail":
		source := exampleData.Data[papertrail.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(papertrail.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "pingdom":
		source := exampleData.Data[pingdom.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(pingdom.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "raygun":
		sender.SendCcMessage(raygun.ExampleMessage(cfg, exampleData))
	case "runscope":
		sender.SendCcMessage(runscope.ExampleMessage(cfg, exampleData))
	case "semaphore":
		source := exampleData.Data[semaphore.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(semaphore.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "statuspage":
		source := exampleData.Data[statuspage.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(statuspage.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "travisci":
		sender.SendCcMessage(travisci.ExampleMessage(cfg, exampleData))
	case "userlike":
		source := exampleData.Data[userlike.HandlerKey]
		for _, eventSlug := range source.EventSlugs {
			sender.SendCcMessage(userlike.ExampleMessage(cfg, exampleData, eventSlug))
		}
	case "victorops":
		sender.SendCcMessage(victorops.ExampleMessage(cfg, exampleData))
	default:
		panic(fmt.Sprintf("Unknown webhook source %v\n", example))
	}
}
