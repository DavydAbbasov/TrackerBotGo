package events
//Worker pool — управляемый параллелизм:
// не больше M тяжёлых задач одновременно.
import "context"

//пул воркеров (горутины‑потребители).
//Пул — это M горутин, каждая делает бесконечный
//цикл чтения из канала и обработки:

type WorkerPool interface {
	Start(ctx context.Context, bus Bus, workers int, handle func(Event)) // запускает M воркеров
	// (можно вернуть функцию Stop или просто полагаться на ctx.Done + закрытие bus)
}
