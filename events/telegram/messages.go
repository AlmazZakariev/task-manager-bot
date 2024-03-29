package telegram

//TODO: нужно встроить константу вместо :% и даты
const msgHelp = `Я буду напоминать тебе о важных задачах! 
Для этого отправь сообщение в виде: 
"Текст задачи" :%ДД.ММ.ГГГГ чч:мм
ВАЖНО!!! Перед установкой даты символы " :%" обязательны

Пример корректного сообщения выглядит так:
Узнать, когда мне вышлют деньги за съёмку :%20.01.2024 12:00

Чтобы посмотреть список всех запланированных задач введите команду /show

Обрати внимание, что я нахожусь на стадии разработки,
поэтому пока не могу гарантировать конфиденциальность данных

Следите за обновлениями!`

// Чтобы удалить задачу ввведите команду /delete x
// Где вместо x - выступает номер этой задачи

const msgHello = "Привет! \n\n" + msgHelp

const (
	msgUnknownCommand = "Что-то пошло не так.👻 Проверьте список команд с помощью /help"
	msgNoSavedTasks   = "Кажется, у вас нет запланированных задач🤤"
	msgSaved          = "Задача добавлена!🤩"
	msgAlreadyExists  = "В списке уже есть точно такая же задача🤔\nСписок всех задач по команде /show"
	msgDeleted        = "Вы успешно избавились от этой задачи🗑"
	msgShowTask       = "Внимание!💡 Пришло время задачи:\n"

	msgDontReilised = "Это функция сейчас не поддерживается😞"
)
