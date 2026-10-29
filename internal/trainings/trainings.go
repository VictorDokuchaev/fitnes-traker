package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

const (
	dataSeparator  = ","
	dataPartsCount = 3 // количество шагов, тип тренировки и продолжительность.
)

// Поддерживаемые типы тренировок.
const (
	runningType = "Бег"
	walkingType = "Ходьба"
)

// Ошибки разбора и обработки данных о тренировке.
var (
	ErrInvalidFormat       = errors.New("неверный формат строки данных")
	ErrInvalidSteps        = errors.New("количество шагов должно быть больше нуля")
	ErrInvalidDuration     = errors.New("продолжительность должна быть больше нуля")
	ErrUnknownTrainingType = errors.New("неизвестный тип тренировки")
)

// Training содержит данные о тренировке пользователя.
type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

// Parse разбирает строку формата "3456,Ходьба,3h00m" и сохраняет данные в структуру.
// При ошибке поля структуры не изменяются.
func (t *Training) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, dataSeparator)
	if len(parts) != dataPartsCount {
		return fmt.Errorf("%w: %q", ErrInvalidFormat, datastring)
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("ошибка парсинга количества шагов: %w", err)
	}
	if steps <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidSteps, steps)
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return fmt.Errorf("ошибка парсинга продолжительности: %w", err)
	}
	if duration <= 0 {
		return fmt.Errorf("%w: %v", ErrInvalidDuration, duration)
	}

	t.Steps = steps
	t.TrainingType = parts[1]
	t.Duration = duration

	return nil
}

// ActionInfo возвращает строку с информацией о тренировке.
func (t Training) ActionInfo() (string, error) {
	distance := spentenergy.Distance(t.Steps, t.Height)
	meanSpeed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	var (
		calories float64
		err      error
	)

	switch t.TrainingType {
	case runningType:
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case walkingType:
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", ErrUnknownTrainingType
	}

	if err != nil {
		return "", fmt.Errorf("ошибка расчета калорий: %w", err)
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, meanSpeed, calories,
	), nil
}
