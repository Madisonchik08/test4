package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: ожидалось 3 части, разделенные запятой, получено: %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать количество шагов в число: %w", err)
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать продолжительность: %w", err)
	}

	trainingType := parts[1]

	return steps, trainingType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64

	switch trainingType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:

		return "", fmt.Errorf("неизвестный тип тренировки: %s", trainingType)
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	avgSpeed := meanSpeed(steps, height, duration)

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		trainingType,
		duration.Hours(),
		dist,
		avgSpeed,
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные параметры: шаги, вес, рост и продолжительность должны быть положительными")
	}

	avgSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()
	calories := (weight * avgSpeed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные параметры: шаги, вес, рост и продолжительность должны быть положительными")
	}

	avgSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()
	baseCalories := (weight * avgSpeed * durationInMinutes) / minInH

	return baseCalories * walkingCaloriesCoefficient, nil
}
