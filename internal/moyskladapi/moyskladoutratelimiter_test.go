package moyskladapi_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"mstorefgo/internal/moyskladapi" // замените на актуальный путь
)

// Базовый тест работы лимитера
func TestRatelimiter_Basic(t *testing.T) {
	limit := 1
	interval := time.Second
	r := moyskladapi.NewRatelimiter(limit, interval)
	defer r.Stop()

	start := time.Now()
	r.Wait() // не должен блокироваться
	r.Wait() // должен ждать ~1 секунду
	duration := time.Since(start)

	if duration < interval {
		t.Fatalf("Ожидалась задержка не менее 1 секунды, получено: %v", duration)
	}
}

// Тест остановки лимитера
func TestRatelimiter_Stop(t *testing.T) {
	r := moyskladapi.NewRatelimiter(1, time.Second)
	r.Stop()

	// Проверяем, что канал закрыт
	_, ok := <-r.Chan()
	if ok {
		t.Fatal("Канал должен быть закрыт после Stop()")
	}
}

// Нагрузочный тест с выводом статистики
func TestRatelimiter_Workload(t *testing.T) {
	const (
		limit        = 40               // 100 запросов
		interval     = 3 * time.Second  // в секунду
		testDuration = 30 * time.Second // общее время теста
		reportPeriod = 3 * time.Second  // период отчетности
		workers      = 10               // количество горутин
	)

	r := moyskladapi.NewRatelimiter(limit, interval)
	defer r.Stop()

	var (
		requestCount uint64 // счетчик запросов
		ctx, cancel  = context.WithTimeout(context.Background(), testDuration)
	)
	defer cancel()

	// Запускаем рабочих
	for i := 0; i < workers; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					r.Wait()
					atomic.AddUint64(&requestCount, 1)
				}
			}
		}()
	}

	// Сбор статистики
	var (
		prevCount    uint64
		start        = time.Now()
		statsTicker  = time.NewTicker(reportPeriod)
		maxPerPeriod = uint64(limit * int(reportPeriod/interval))
	)
	defer statsTicker.Stop()

	for {
		select {
		case <-statsTicker.C:
			total := atomic.LoadUint64(&requestCount)
			current := total - prevCount
			prevCount = total

			t.Logf("[%v] Запросов за период: %d (макс. допуск: %d ± 10%%)",
				time.Since(start).Round(time.Second),
				current,
				maxPerPeriod,
			)

			// Проверка что не превышаем лимит более чем на 10%
			if float64(current) > float64(maxPerPeriod)*1.1 {
				t.Errorf("Превышен лимит: %d > %d", current, maxPerPeriod)
			}

		case <-ctx.Done():
			t.Logf("Всего обработано запросов: %d", atomic.LoadUint64(&requestCount))
			return
		}
	}
}
