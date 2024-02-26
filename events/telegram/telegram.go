package telegram

import (
	"context"
	"errors"
	"fmt"
	"tasks-manager-bot/clients/telegram"
	"tasks-manager-bot/events"
	"tasks-manager-bot/lib/e"
	"tasks-manager-bot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserID   int
	UserName string
}

var (
	ErrUnknownEventType = errors.New("unknown message type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, U := range updates {
		res = append(res, event(U))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.TaskSending:
		return p.processTaskMessage(event)
	default:
		return e.WrapIfErr("can't process message", ErrUnknownEventType)
	}

}

// TODO: разобраться с context.Context
func (p *Processor) Execute() ([]events.Event, error) {
	tasks, err := p.storage.GetTasksToExecute(context.Background())
	if err != nil {
		return nil, e.Wrap("can't get task from storage", err)
	}
	if tasks == nil {
		return nil, nil
	}
	res := make([]events.Event, 0, len(*tasks))
	for _, task := range *tasks {
		res = append(res, eventFromTask(task))
	}
	return res, nil
}

func (p *Processor) processMessage(event events.Event) (err error) {
	defer func() { err = e.WrapIfErr("can't process message", err) }()
	meta, err := meta(event)
	if err != nil {
		return err
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.UserID, meta.UserName); err != nil {
		return err
	}
	return nil
}
func (p *Processor) processTaskMessage(event events.Event) (err error) {
	defer func() { err = e.WrapIfErr("can't process task sending message", err) }()
	meta, err := meta(event)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("%s\n%s", msgShowTask, event.Text)
	return p.tg.SendMessage(meta.ChatID, msg)
}
func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}
func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserID:   upd.Message.From.ID,
			UserName: upd.Message.From.UserName,
		}
	}
	return res
}
func eventFromTask(task storage.Task) events.Event {
	meta := Meta{
		UserID: task.UserId,
		ChatID: task.ChatId,
	}
	res := events.Event{
		Type: events.TaskSending,
		Text: task.Text,
		Meta: meta,
	}
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
