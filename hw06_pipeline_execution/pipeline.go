package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(forwarder(in, done))
	}
	return in
}

// forwarder - передает данные между стейджами, учитывая сигнал остановки.
func forwarder(in In, done In) Out {
	lint := false
	out := make(Bi)
	go func() {
		defer func() {
			close(out) // Закрываем выходной канал
			// Читаем остаточные данные из канала in, чтобы избежать блокировок
			for range in {
				// lint
				if lint {
					break
				}
			}
		}()

		for {
			select {
			case <-done:
				// Если пришел сигнал завершения, выходим из горутины
				return
			case v, ok := <-in:
				if !ok {
					// Если входной канал закрыт, завершаем работу
					return
				}
				select {
				case <-done:
					// Если пришел сигнал завершения, выходим из горутины
					return

				// Отправляем результат в выходной канал, если нет сигнала завершения
				case out <- v:
					break
				}
			}
		}
	}()

	return out
}
