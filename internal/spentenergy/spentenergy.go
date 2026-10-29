package spentenergy

import (
	"errors"
	"fmt"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

// Ошибки валидации входных параметров расчета калорий.
var (
	ErrInvalidSteps    = errors.New("количество шагов должно быть больше нуля")
	ErrInvalidWeight   = errors.New("вес должен быть больше нуля")
	ErrInvalidHeight   = errors.New("рост должен быть больше нуля")
	ErrInvalidDuration = errors.New("продолжительность должна быть больше нуля")
)

// WalkingSpentCalories возвращает количество калорий, потраченных при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateParams(steps, weight, height, duration); err != nil {
		return 0, err
	}

	return spentCalories(steps, weight, height, duration) * walkingCaloriesCoefficient, nil
}

// RunningSpentCalories возвращает количество калорий, потраченных при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateParams(steps, weight, height, duration); err != nil {
		return 0, err
	}

	return spentCalories(steps, weight, height, duration), nil
}

func spentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	meanSpeed := MeanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	return weight * meanSpeed * durationInMinutes / minInH
}

func validateParams(steps int, weight, height float64, duration time.Duration) error {
	switch {
	case steps <= 0:
		return fmt.Errorf("%w: %d", ErrInvalidSteps, steps)
	case weight <= 0:
		return fmt.Errorf("%w: %.2f", ErrInvalidWeight, weight)
	case height <= 0:
		return fmt.Errorf("%w: %.2f", ErrInvalidHeight, height)
	case duration <= 0:
		return fmt.Errorf("%w: %v", ErrInvalidDuration, duration)
	}

	return nil
}

// MeanSpeed возвращает среднюю скорость в км/ч.
// Для некорректного количества шагов или продолжительности возвращает 0.
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || duration <= 0 {
		return 0
	}

	return Distance(steps, height) / duration.Hours()
}

// Distance возвращает дистанцию в километрах по количеству шагов и росту пользователя в метрах.
func Distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	return float64(steps) * stepLength / mInKm
}
