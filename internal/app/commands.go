package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"strings"
	"transmission-telegram-bot/internal/pkg/rpc"
)

func start(id int64) tgbotapi.MessageConfig {
	s := `
	Hello! 👋
	Welcome to the Transmission Telegram Bot!
	For a list of commands you can use, respond
	with /command
`
	return tgbotapi.NewMessage(id, s)

}

func play(chatID int64, torrentID int, tc *rpc.TransmissionConnection) tgbotapi.MessageConfig {
	err := tc.StartTorrent(torrentID)
	if err != nil {
		log.Warn(err)
	}

	s := fmt.Sprintf("Started torrent ID %d.\n", torrentID)

	return tgbotapi.NewMessage(chatID, s)
}

func pause(chatID int64, torrentID int, tc *rpc.TransmissionConnection) tgbotapi.MessageConfig {
	err := tc.PauseTorrent(torrentID)
	if err != nil {
		log.Error(err)
	}

	s := fmt.Sprintf("Stopped torrent ID %d.\n", torrentID)

	return tgbotapi.NewMessage(chatID, s)
}

func list(chatID int64, tc *rpc.TransmissionConnection) tgbotapi.MessageConfig {

	torrents, err := tc.GetTorrentList(false)
	if err != nil {
		log.Warn(err)
	}

	s := strings.Builder{}

	if len(torrents) == 0 {
		s.WriteString("No torrents. Maybe you should add one with /add?")
	}

	for i, t := range torrents {
		s.WriteString(fmt.Sprintf("<b>ID:</b> %d\n", *t.ID))
		s.WriteString(fmt.Sprintf("<b>Name:</b> %s\n", *t.Name))
		s.WriteString(fmt.Sprintf("<b>Completion:</b> %.2f%%\n", 100*(*t.PercentDone)))
		s.WriteString(fmt.Sprintf("<b>State:</b> %s\n", t.Status.String()))
		if i < len(torrents) - 1 {
			s.WriteString("==========\n")
		}
	}

	msg := tgbotapi.NewMessage(chatID, s.String())
	msg.ParseMode = "HTML"
	return msg
}

func add(chatID int64, torrentFileURL string, tc *rpc.TransmissionConnection) tgbotapi.MessageConfig {
	s := strings.Builder{}
	torrent, err := tc.AddTorrent(torrentFileURL)
	if err != nil {
		log.Warn(err)
		s.WriteString(fmt.Sprintf("Failed to add torrent. %v", err))
	} else {
		s.WriteString(fmt.Sprintf("Added %s. See its status with the /list command.", *torrent.Name))
		log.Info(fmt.Sprintf("Added <code>%s</code>. See its status with the /list command.", *torrent.Name))
	}

	msg := tgbotapi.NewMessage(chatID, s.String())
	msg.ParseMode = "HTML"
	return msg

}

func remove(chatID int64, torrentID int, deleteData bool, tc *rpc.TransmissionConnection) tgbotapi.MessageConfig {
	s := strings.Builder{}

	err := tc.RemoveTorrent(torrentID, deleteData)
	if err != nil {
		s.WriteString(fmt.Sprintf("Failed to remove torrent. %v", err))
	} else {
		s.WriteString(fmt.Sprintf("Removed torrent."))
	}
	msg := tgbotapi.NewMessage(chatID, s.String())
	return msg
}