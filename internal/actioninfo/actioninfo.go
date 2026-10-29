package actioninfo

import (
	"fmt"
	"log"
)

// DataParser описывает активность, которая умеет разбирать строку с данными
// и формировать по ним информационное сообщение.
type DataParser interface {
	Parse(datastring string) error
	ActionInfo() (string, error)
}

// Info разбирает каждую строку набора данных и выводит информацию об активности.
// Ошибки разбора и формирования информации логируются, обработка продолжается.
func Info(dataset []string, dp DataParser) {
	for _, data := range dataset {
		if err := dp.Parse(data); err != nil {
			log.Printf("ошибка парсинга данных %q: %v", data, err)
			continue
		}

		info, err := dp.ActionInfo()
		if err != nil {
			log.Printf("ошибка формирования информации для данных %q: %v", data, err)
			continue
		}

		fmt.Println(info)
	}
}
