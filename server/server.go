package server

import (
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"log/slog"
	"os"
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
	slog.Info("Starting http server", "port", serverPort)
	channels := []Channel{}
	err := mapstructure.Decode(c, &channels)
	if err != nil {
		slog.Error("decode channels error", "err", err)
		os.Exit(1)
	}

	http.HandleFunc("/messages", apiMessages(bot, channels))
	// http.HandleFunc("/modules", apiModules(bot, channels))
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", serverPort), nil)

	if err != nil {
		slog.Error("ListenAndServe failed", "err", err)
		os.Exit(1)
	}

}

func apiMessages(bot *tele.Bot, channels []Channel) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("GET params", "params", r.URL.Query())

		authToken := r.URL.Query().Get("token")
		chatID := getChatID(authToken, channels)
		data := r.URL.Query().Get("data")

		if authToken != "" && chatID != 0 {

			if data != "" {
				var to = tele.Chat{}
				to.ID = chatID
				_, err := bot.Send(&to, data)
				if err != nil {
					slog.Error("could not send message", "err", err)
				}
				_, err = w.Write([]byte("OK"))
				if err != nil {
					slog.Error("write error", "err", err)
				}
			}

		}

		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("Forbidden"))

		if err != nil {
			slog.Error("write forbidden error", "err", err)
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
