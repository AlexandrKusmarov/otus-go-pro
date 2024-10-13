package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// считать это как "максимум 0 ошибок", значит функция всегда будет возвращать `ErrErrorsLimitExceeded`;
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	taskCh := make(chan func() error) // Канал для отправки задач
	errCh := make(chan error, n)      // Буферизированный канал для ошибок
	var wg sync.WaitGroup
	var errCount int
	var mu sync.Mutex

	// Запуск n воркеров
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				task, ok := <-taskCh
				if !ok {
					return // Канал закрыт, завершаем воркер
				}
				if err := task(); err != nil {
					incrementErrCountAndGet(&mu, &errCount)
					mu.Lock()
					if errCount >= m {
						errCh <- ErrErrorsLimitExceeded
						mu.Unlock()
						return
					}
					mu.Unlock()
				}
			}
		}()
	}

	// Заполняем задачи
	go func() {
		defer close(taskCh)
		for _, task := range tasks {
			select {
			case taskCh <- task:
			case <-errCh:
				return // Прекращаем отправку задач, если лимит ошибок превышен
			}
		}
	}()

	// Ожидаем завершения всех воркеров
	wg.Wait()
	close(errCh)

	return checkErr(errCount, m)
}

func incrementErrCountAndGet(mu *sync.Mutex, errCount *int) {
	defer mu.Unlock()
	mu.Lock()
	*errCount++
}

func checkErr(errCount, m int) error {
	// Если было больше ошибок, возвращаем ошибку
	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
