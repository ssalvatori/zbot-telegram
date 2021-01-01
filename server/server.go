package server

import (
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Channel definition
type Channel struct {
	ID        int64
	Title     string
	AuthToken string
}

//Start http server in a given port
func Start(serverPort int, bot *tb.Bot, c interface{}) {
	log.Info(fmt.Sprintf("Starting http server at port: %d", serverPort))
	channels := []Channel{}
	mapstructure.Decode(c, &channels)

	http.HandleFunc("/messages", apiMessages(bot, channels))
	// http.HandleFunc("/modules", apiModules(bot, channels))
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", serverPort), nil)

	if err != nil {
		log.Fatal("ListenAndServe: " + err.Error())
	}

}

func apiMessages(bot *tb.Bot, channels []Channel) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug(fmt.Sprintf("GET params: %v", r.URL.Query()))

		authToken := r.URL.Query().Get("token")
		chatID := getChatID(authToken, channels)
		data := r.URL.Query().Get("data")

		if authToken != "" && chatID != 0 {

			if data != "" {
				var to = tb.Chat{}
				to.ID = chatID
				bot.Send(&to, data)
				w.Write([]byte(fmt.Sprintf("OK")))
			}

		}

		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf("Forbidden")))
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
