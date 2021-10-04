package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"youtube-stream-notifier-linebot/controller"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func LineCallbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := controller.Bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeLeave {
			log.Println("Left group", event.Source.GroupID)
			// Delete Group in DB
			DeleteGroup(event.Source.GroupID)
		}
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				switch {
				case event.Source.GroupID != "":
					//In the group
					if strings.EqualFold(message.Text, "Farewell") {
						if _, err = controller.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Chao")).Do(); err != nil {
							log.Fatal(err)
						}
						// Leave the group
						controller.Bot.LeaveGroup(event.Source.GroupID).Do()

					}
				default:
					if _, err = controller.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("What do you mean?")).Do(); err != nil {
						log.Print(err)
					}
					log.Println(message.Text)
				}
			}
		}
		if event.Type == linebot.EventTypeJoin {
			// If join into a Group
			if event.Source.GroupID != "" {
				if groupRes, err := controller.Bot.GetGroupSummary(event.Source.GroupID).Do(); err == nil {
					retString := fmt.Sprintf("Thanks for inviting me to: %s\nI'll notify you when my YouTube stream starts\nIf you want me to leave, just say 'Farewell'", groupRes.GroupName)
					if _, err = controller.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(retString)).Do(); err != nil {
						log.Fatal(err)
					}
					log.Println("Joined group:", event.Source.GroupID)

					// Insert new Group to DB
					InsertGroup(event.Source.GroupID, groupRes.GroupName)

				} else {
					//GetGroupSummary fail/.
					log.Printf("GetGroupSummary:%x", err)
				}
			}
		}
	}
}

func PubsubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// handle callback from pubsubhubub
	switch r.Method {
	case "GET":
		log.Println(r.Method, r.URL)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, r.URL.Query()["hub.challenge"][0])
		return
	case "POST":
		log.Println(r.Method, r.URL)
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(nil)
		}
		title, video_link := ParseXML(data)
		log.Println("title:", title, "link:", video_link)

		// broadcast message to every user
		retString := fmt.Sprintf("My new YouTute stream is on-live\nTitle: %s\nLink: %s\n", title, video_link)
		if _, err = controller.Bot.BroadcastMessage(linebot.NewTextMessage(retString)).Do(); err != nil {
			log.Fatal(err)
		}

		// broadcast message to every group it joins
		for _, group_obj := range ListGroup() {
			group_id := group_obj.group_id
			if _, err = controller.Bot.PushMessage(group_id, linebot.NewTextMessage(retString)).Do(); err != nil {
				log.Print(err)
			}
		}
		return
	}
}

func Subscribe() {
	// construct message body to POST
	YT_CHANNEL_ID := os.Getenv("YT_CHANNEL_ID")
	HUB_URL := "https://pubsubhubbub.appspot.com/subscribe"
	TOPIC_URL := "https://www.youtube.com/xml/feeds/videos.xml?channel_id=" + YT_CHANNEL_ID
	CALLBACK_URL := "https://youtube-stream-notifier-linebot.herokuapp.com/subscribe/"

	body := url.Values{}
	body.Set("hub.callback", CALLBACK_URL)
	body.Set("hub.topic", TOPIC_URL)
	body.Set("hub.verify", "async")
	body.Set("hub.mode", "subscribe")

	r, err := http.NewRequest("POST", HUB_URL, strings.NewReader(body.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// build http connection and send request
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}

	// check the response status
	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		log.Fatal(res.Status)
	} else {
		log.Println(res.Status)
	}

}

func Hello(w http.ResponseWriter, r *http.Request) {
	msg := "Hello World!"
	fmt.Fprintf(w, "%s", msg)
}
