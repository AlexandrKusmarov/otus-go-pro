package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// считать это как "максимум 0 ошибок", значит функция всегда будет возвращать `ErrErrorsLimitExceeded`;
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var wg sync.WaitGroup
	var errCount int64
	taskCh := make(chan Task, len(tasks))
	errCh := make(chan error, m)

	// Инициализируем канал с задачами, наполнив его списком переданных задач
	go func() {
		defer close(taskCh)
		for _, task := range tasks {
			taskCh <- task
		}
	}()

	// Создаем n параллельных горутин для выполнения задач
	// Внутрь каждой горутины передаем задачи из канала задач,
	// где так же выполняется проверка на переполнение допустимого количества ошибок
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if task == nil {
					continue
				}
				if err := task(); err != nil {
					// Селект нужен для потокобезопасной работы с каналом, блокирующая операция.
					select {
					// Если не удалось записать в канал - будет выполнен блок default. Это защита от deadlock.
					case errCh <- err:
					default:
					}
					if atomic.AddInt64(&errCount, 1) >= int64(m) {
						return
					}
				}
			}
		}()
	}

	// Ожидаем завершения всех горутин
	wg.Wait()

	// Проверяем, было ли превышено количество ошибок
	if atomic.LoadInt64(&errCount) >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
