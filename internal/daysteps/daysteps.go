package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// TODO: реализовать функцию
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных: ожидалось 'шаги,продолжительность', получено '%s'", data)
	}
	steps, err := strconv.Atoi(parts[0])
	if err != nil {

		return 0, 0, fmt.Errorf("не удалось преобразовать шаги в число: %w", err)
	}

	//Количество шагов должно быть больше 0.
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным числом, получено: %d", steps)
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		// При ошибке вернуть нули и ошибку.
		return 0, 0, fmt.Errorf("не удалось преобразовать продолжительность: %w", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть положительной, получено: %v", duration)
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {

		log.Printf("Ошибка при обработке данных: %v\n", err)
		return ""
	}

	distanceMeters := float64(steps) * stepLength

	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("Ошибка при расчете калорий: %v\n", err)
		return ""
	}

	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		calories,
	)

	return result
}
