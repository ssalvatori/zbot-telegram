package server

import (
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

//Channel definition
type Channel struct {
	ID        int64
	Title     string
	AuthToken string
}

//Start http server in a given port
func Start(serverPort int, bot *tele.Bot, c interface{}) {
	log.Info(fmt.Sprintf("Starting http server at port: %d", serverPort))
	channels := []Channel{}
	err := mapstructure.Decode(c, &channels)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/messages", apiMessages(bot, channels))
	// http.HandleFunc("/modules", apiModules(bot, channels))
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", serverPort), nil)

	if err != nil {
		log.Fatal("ListenAndServe: " + err.Error())
	}

}

func apiMessages(bot *tele.Bot, channels []Channel) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug(fmt.Sprintf("GET params: %v", r.URL.Query()))

		authToken := r.URL.Query().Get("token")
		chatID := getChatID(authToken, channels)
		data := r.URL.Query().Get("data")

		if authToken != "" && chatID != 0 {

			if data != "" {
				var to = tele.Chat{}
				to.ID = chatID
				_, err := bot.Send(&to, data)
				if err != nil {
					log.Error("Could not set the message")
					log.Error(err)
				}
				_, err = w.Write([]byte("OK"))
				if err != nil {
					log.Error(err)
				}
			}

		}

		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("Forbidden"))

		if err != nil {
			log.Error(err)
		}
	}
}

//getChatId return the chat_id associated with that token
func getChatID(token string, channels []Channel) int64 {

	for i := range channels {
		if channels[i].AuthToken == token {
			return channels[i].ID
		}
	}

	return 0
}
