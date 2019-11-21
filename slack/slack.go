package slack

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/nlopes/slack"
)

/*
   NOTE: command_arg_1 and command_arg_2 represent optional parameteras that you define
   in the Slack API UI
*/
const helpMessage = "type in '@Go-Syntax-Helper <package> <functionName>'"

/*
   CreateSlackClient sets up the slack RTM (real-timemessaging) client library,
   initiating the socket connection and returning the client.
   DO NOT EDIT THIS FUNCTION. This is a fully complete implementation.
*/
func CreateSlackClient(apiKey string) *slack.RTM {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection() // goroutine!
	return rtm
}

/*
   RespondToEvents waits for messages on the Slack client's incomingEvents channel,
   and sends a response when it detects the bot has been tagged in a message with @<botTag>.

   EDIT THIS FUNCTION IN THE SPACE INDICATED ONLY!
*/
func RespondToEvents(slackClient *slack.RTM) {
	for msg := range slackClient.IncomingEvents {
		fmt.Println("Event Received: ", msg.Type)
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			botTagString := fmt.Sprintf("<@%s> ", slackClient.GetInfo().User.ID)
			if !strings.Contains(ev.Msg.Text, botTagString) {
				continue
			}
			message := strings.Replace(ev.Msg.Text, botTagString, "", -1)

			// TODO: Make your bot do more than respond to a help command. See notes below.
			// Make changes below this line and add additional funcs to support your bot's functionality.
			// sendHelp is provided as a simple example. Your team may want to call a free external API
			// in a function called sendResponse that you'd create below the definition of sendHelp,
			// and call in this context to ensure execution when the bot receives an event.

			// START SLACKBOT CUSTOM CODE
			// ===============================================================
			sendResponse(slackClient, message, ev.Channel)
			sendHelp(slackClient, message, ev.Channel)
			sendExample(slackClient, message, ev.Channel)
			// ===============================================================
			// END SLACKBOT CUSTOM CODE
		default:

		}
	}
}

// sendHelp is a working help message, for reference.
func sendHelp(slackClient *slack.RTM, message, slackChannel string) {
	if strings.ToLower(message) != "help" {
		return
	}
	slackClient.SendMessage(slackClient.NewOutgoingMessage(helpMessage, slackChannel))
}

// sendResponse is NOT unimplemented --- write code in the function body to complete!

func sendResponse(slackClient *slack.RTM, message, slackChannel string) {
	// command := strings.ToLower(message)
	println("[RECEIVED] sendResponse:", message)

	splitMessage := strings.Split(message, " ")

	// if len(splitMessage) > 2 {
	// 	slackClient.SendMessage(slackClient.NewOutgoingMessage("Uh oh - too many arguments! Type in `@go-syntax-helper` for proper formating and help.", slackChannel))
	// 	return
	// }

	if len(splitMessage) == 2 {
		outputString := fmt.Sprintf("https://golang.org/pkg/%s/#%s", splitMessage[0], splitMessage[1])
		slackClient.SendMessage(slackClient.NewOutgoingMessage(outputString, slackChannel))
	} else if len(splitMessage) == 1 && message != "help"{
		outputString := fmt.Sprintf("https://golang.org/pkg/%s", message)
		slackClient.SendMessage(slackClient.NewOutgoingMessage(outputString, slackChannel))
	}

	// fmt.Println("Split message:", splitMessage)


	// START SLACKBOT CUSTOM CODE
	// ===============================================================
	// TODO:
	//      1. Implement sendResponse for one or more of your custom Slackbot commands.
	//         You could call an external API here, or create your own string response. Anything goes!
	//      2. STRETCH: Write a goroutine that calls an external API based on the data received in this function.
	// ===============================================================
	// END SLACKBOT CUSTOM CODE
}

func sendExample(slackClient *slack.RTM, message, slackChannel string) {
	splitMessage := strings.Split(message, " ")

	c := colly.NewCollector(
			// colly.AllowedDomains("golang.org"),
	)
	fmt.Println("Has gotten c")
	// fmt.Println(c)

	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	fmt.Println("In onhtml" )
	// 	fmt.Println(e.Text)
	// })

	c.OnHTML("h2[id=ToUpper]+pre", func(e *colly.HTMLElement) {
		// id := e.Attr("id")
		// Print link
		// fmt.Printf("Link found: %q -> %s\n", e.Text, id)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		// c.Visit(e.Request.AbsoluteURL(link))
		fmt.Println(e.Text)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	if len(splitMessage) == 3 && splitMessage[0] == "example" {
		// fmt.Println("Printing if statement")
		link :=fmt.Sprintf("https://golang.org/pkg/%s/", splitMessage[1])
		c.Visit(link)
	}
}
