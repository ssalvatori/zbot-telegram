package middleware

import (
	"log/slog"

	tele "gopkg.in/telebot.v4"
)

// AdminOnly restricts a handler to senders for which isAdmin returns true.
// Updates without a sender or Telegram user ID are dropped (fail closed).
func AdminOnly(isAdmin func(userID int64) bool) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {

			sender := c.Sender()
			if sender == nil {
				slog.Warn("admin check: no sender in update")
				return nil
			}

			if sender.ID == 0 {
				slog.Warn("admin check: sender has no Telegram ID", "user", sender.FirstName)
				return nil
			}

			if !isAdmin(sender.ID) {
				slog.Warn("admin check: non-admin user rejected", "user_id", sender.ID)
				return nil
			}

			return next(c)
		}
	}
}
