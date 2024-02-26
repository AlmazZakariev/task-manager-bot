package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"tasks-manager-bot/lib/e"
	"tasks-manager-bot/storage"
	"time"
	"unicode/utf8"
)

const (
	DltCmd   = "/delete"
	ChowCmd  = "/show"
	StartCmd = "/start"
	HelpCmd  = "/help"
)
const (
	DateStartString = " :%"
	DateLayout      = "2.01.2006 15:04"
)

// TODO: разобраться с context.Context
func (p *Processor) doCmd(text string, chatID int, userID int, userName string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, userName)

	isAddCmd, err := isAddCmd(text)
	if err != nil {
		//TODO: теряется ошибка из функции isAddCmd
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
	if isAddCmd {
		return p.saveTask(chatID, text, userID, userName)
	}
	switch text {
	case DltCmd:
		return p.DeleteTask(chatID, userID)
	case StartCmd:
		return p.SendHello(chatID)
	case HelpCmd:
		return p.SendHelp(chatID)
	case ChowCmd:
		return p.ShowTasks(chatID, userID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}
func (p *Processor) saveTask(chatID int, messageText string, userID int, userName string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()
	date, taskText, err := parseMessage(messageText)
	task := &storage.Task{
		//TODO: в какой-то момент нужно получать количество задач текущего пользователя для формирования id, оно будет не уникальным и представлять
		//пока вырезал поле ID из-за отсутствия смысла

		UserId:   chatID,
		Date:     date,
		Text:     taskText,
		ChatId:   chatID,
		UserName: userName,
	}
	isExists, err := p.storage.IsExists(context.Background(), task)
	if err != nil {
		return nil
	}

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}
	if err := p.storage.Save(context.Background(), task); err != nil {
		return err
	}
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) SendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) SendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) ShowTasks(chatID int, userID int) error {
	tasks, err := p.storage.GetAllTasks(context.Background(), userID)
	if err != nil {
		return e.Wrap("Can't show tasks", err)
	}
	msg := "Список ваших задач\n"

	for i, task := range *tasks {
		msg += fmt.Sprintf("%d. %s: %v\n", i+1, task.Text, task.Date.Format(DateLayout))
	}
	return p.tg.SendMessage(chatID, msg)
}
func (p *Processor) DeleteTask(chatID int, userID int) error {
	return p.tg.SendMessage(chatID, msgDontReilised)
}

func parseMessage(msg string) (date time.Time, text string, err error) {
	defer func() { err = e.WrapIfErr("can't parse message", err) }()

	text, err = msgTaskText(msg)
	if err != nil {
		//TODO: почему не могу вернуть время nil?????
		return time.Now(), "", err
	}
	date, err = msgTaskDate(msg)
	if err != nil {
		//TODO: почему не могу вернуть время nil?????
		return time.Now(), "", err
	}
	return date, text, err

}

func msgTaskDate(msg string) (date time.Time, err error) {
	defer func() { err = e.WrapIfErr("can't get task date from message", err) }()
	dateFirtsIndex, err := dateStartPos(msg)
	if err != nil {
		//TODO: почему не могу вернуть время nil?????
		return time.Now(), err
	}

	dateString := strings.TrimSpace(string([]rune(msg)[dateFirtsIndex:]))
	date, err = time.ParseInLocation(DateLayout, dateString, time.Local)

	if err != nil {
		return time.Now(), err
	}

	return date, err
}

func msgTaskText(msg string) (text string, err error) {
	defer func() { err = e.WrapIfErr("can't get task text from message", err) }()

	textLastIndex, err := textEndedPos(msg)
	if err != nil {
		return "", err
	}
	text = string([]rune(msg)[:textLastIndex])

	return text, err
}
func isAddCmd(text string) (AddCmd bool, err error) {
	defer func() { err = e.WrapIfErr("can't define is it AddCmd", err) }()
	pattern := `.+` + DateStartString + `\s*\d\d.\d\d.\d\d\d\d\s*\d\d:\d\d\s*`
	matched, err := regexp.MatchString(pattern, text)
	if err != nil {

		return false, err

	}
	if !matched {

		return false, nil
	}
	return true, err
}

// DSS mean DateStartString
func dateStartPos(text string) (pos int, err error) {

	iDSS := strings.Index(text, DateStartString)
	if iDSS == -1 {
		return 0, errors.New("can't find date start position")
	}
	return utf8.RuneCountInString(text[:iDSS]) + dssLen(), nil

}
func textEndedPos(text string) (pos int, err error) {
	DSS, err := dateStartPos(text)
	if err != nil {
		return 0, e.WrapIfErr("can't find text ended position", err)
	}
	return DSS - dssLen(), nil
}
func dssLen() int {
	return len([]rune(DateStartString))
}

//DSS mean DateStartString
