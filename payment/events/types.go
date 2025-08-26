package events

import "time"

//(структура) модель события.
//Event struct — «пакет», который хендлер положит в очередь,
// а воркеры потом обработают.
type Event struct {
	Type       string    // evtType
	ID         string    // evtID
	Raw        []byte    // исходный body []byte (на будущее для повторного парсинга)
	ReceivedAt time.Time // время получения
}
//Зачем Raw: воркер потом может распарсить
//«по‑взрослому» в разные модели, не нагружая хендлер.