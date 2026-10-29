package daysteps

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
	dataPartsCount = 2 // количество шагов и продолжительность прогулки.
)

// Ошибки разбора строки с данными о прогулке.
var (
	ErrInvalidFormat   = errors.New("неверный формат строки данных")
	ErrInvalidSteps    = errors.New("количество шагов должно быть больше нуля")
	ErrInvalidDuration = errors.New("продолжительность должна быть больше нуля")
)

// DaySteps содержит данные о дневной прогулке пользователя.
type DaySteps struct {
	Steps    int
	Duration time.Duration
	personaldata.Personal
}

// Parse разбирает строку формата "678,0h50m" и сохраняет данные в структуру.
// При ошибке поля структуры не изменяются.
func (ds *DaySteps) Parse(datastring string) (err error) {
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

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return fmt.Errorf("ошибка парсинга продолжительности: %w", err)
	}
	if duration <= 0 {
		return fmt.Errorf("%w: %v", ErrInvalidDuration, duration)
	}

	ds.Steps = steps
	ds.Duration = duration

	return nil
}

// ActionInfo возвращает строку с информацией о прогулке.
func (ds DaySteps) ActionInfo() (string, error) {
	distance := spentenergy.Distance(ds.Steps, ds.Height)

	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", fmt.Errorf("ошибка расчета калорий: %w", err)
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, distance, calories,
	), nil
}
