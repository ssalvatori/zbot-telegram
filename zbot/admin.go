package zbot

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/ssalvatori/zbot-telegram/user"
	zbotMiddleware "github.com/ssalvatori/zbot-telegram/zbot/middleware"
	tele "gopkg.in/telebot.v4"
)

// Admin represents a bot administrator with their Telegram user ID and username.
type Admin struct {
	ID       int64
	Username string
}

// IsAdmin checks if the given Telegram user ID is an admin.
// TODO: b.cfg.Admin should be replaced by b.admins
func (b *Bot) IsAdmin(userID int64) bool {
	slog.Debug("List of admins", "admins", b.cfg.Admins)
	slog.Debug("Checking admin status for user", "user_id", userID)
	return userID != 0 && slices.Contains(b.cfg.Admins, userID)
}

// setupAdministrators verifies the configured admin IDs and sends them a notification that the bot is online.
func (b *Bot) setupAdministrators(tb *tele.Bot) []Admin {
	slog.Info("Setting up administrators", "admin_ids", b.cfg.Admins)

	adminList := []Admin{}

	for _, adminID := range b.cfg.Admins {
		chat, err := tb.ChatByID(int64(adminID))
		if err != nil {
			slog.Error("get chat info error", "err", err, "admin_id", adminID)
			continue
		}
		slog.Debug("Retrieved admin chat info", "admin_id", adminID, "chat", chat)
		if _, err := tb.Send(chat, fmt.Sprintf("Hey %s, I'm online now", chat.Username)); err != nil {
			slog.Error("send admin online notification error", "err", err, "admin_id", adminID)
		}
		adminList = append(adminList, Admin{ID: adminID, Username: chat.Username})

	}

	return adminList

}

// registerAdminCommands register administrator's commands
func (b *Bot) registerAdminCommands(tb *tele.Bot) {
	// adminGroup := tb.Group()
	// adminGroup.Use(zbotMiddleware.AdminOnly(b.IsAdmin))

	adminCommands := []tele.Command{}

	for _, admin := range b.admins {

		// for cmd := b.adminCommandsList.List.Front(); cmd != nil; cmd = cmd.Next() {
		for _, cmd := range b.adminCmds {

			slog.Debug("Registering command", "cmd", cmd.CmdString, "for admin_username", admin.Username, "admin_id", admin.ID)

			tb.Handle(cmd.CmdString, func(c tele.Context) error {
				//TODO: Run the Command.ProcessText
				admin := user.BuildUser(c.Sender(), b.db)
				chatID := int(c.Chat().ID)
				chatName := ""

				msg, err := cmd.ProcessCommand(c.Text(), admin, chatID, chatName, true)
				if err != nil {
					slog.Error("Error processing admin command", "err", err)
					return nil
				}
				return c.Send(msg)
			}, zbotMiddleware.AdminOnly(b.IsAdmin))

			adminCommands = append(adminCommands, tele.Command{
				Text:        cmd.CmdString,
				Description: cmd.ShowDescription(),
			})
		}

		slog.Debug("Registering help command for admin", "admin_id", admin.ID)

		tb.Handle("/help", func(c tele.Context) error {
			helpMessage := "Available commands:\n"
			for _, cmd := range b.adminCmds {
				msg := c.Message()
				slog.Debug("Help command received", "cmd", cmd.CmdString, "from_user", msg.Sender.Username, "chat_id", msg.Chat.ID)
				helpMessage += fmt.Sprintf("%s - %s\n", cmd.CmdString, cmd.ShowDescription())
			}
			return c.Send(helpMessage)
		}, zbotMiddleware.AdminOnly(b.IsAdmin))

		adminCommands = append(adminCommands, tele.Command{
			Text:        "/help",
			Description: "Show bots commands help",
		})

		// Commands should be set for each admin individually only one time
		scope := tele.CommandScope{
			Type:   tele.CommandScopeChat,
			ChatID: int64(admin.ID), // in a private chat, chat ID == user ID
		}

		err := tb.SetCommands(adminCommands, scope)
		if err != nil {
			slog.Error("set admin commands error", "err", err)
		}
	}

}
