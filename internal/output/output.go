package output

import (
	"fmt"
	"strings"

	"ConfigScanner/pkg/models"
)

// Formatter интерфейс форматирования вывода
type Formatter interface {
	Format(problems []models.Problem) string
}

// TextFormatter форматирует в человекочитаемый текст
type TextFormatter struct{}

// Format - реализует интерфейс Formatter для текстового вывода
// Преобразует срез проблем в отформатированную текстовую строку
// Использует strings.Builder для эффективной конкатенации строк
// Алгоритм работы:
// 1. Если проблем нет и ShowSuccess=true → выводим сообщение об успехе
// 2. Если проблем нет и ShowSuccess=false → возвращаем пустую строку
// 3. Если проблемы есть выводим нумерованный список с деталями
// Параметры:
// problems: срез найденных проблем
// Возвращает:
// string: отформатированный отчет
func (f TextFormatter) Format(problems []models.Problem) string {
	var builder strings.Builder

	if len(problems) == 0 {
		builder.WriteString("Проблем не найдено! Конфигурация безопасна.\n")
		return builder.String()
	}

	builder.WriteString("Найдены потенциальные проблемы:\n\n")

	// Проходим по всем найденным проблемам и выводим их
	// p - сама проблема (структура models.Problem)
	for _, p := range problems {
		// Выводим уровень опасности и название правила
		// Пример: "HIGH: plain-password. ..."
		builder.WriteString(fmt.Sprintf("%s: %s. ", p.LevelProblem, p.Rule))
		builder.WriteString(fmt.Sprintf("%s. ", p.ParseV))
		builder.WriteString(fmt.Sprintf("%s.", p.Recommendation))
		builder.WriteString("\n")
	}
	return builder.String()
}

// Print выводит проблемы в консоль
// Эта функция предоставляет удобный способ печати отчета, не требуя от вызывающего кода
// работать напрямую с fmt.Print
// Параметры: problems: срез найденных проблем, formatter: форматтер, определяющий вид вывода
// Пример использования:
//
//	output.Print(problems, output.TextFormatter{})
//
// Преимущества:
//   - Скрывает детали реализации (не нужно знать про fmt.Print)
//   - Легко заменить способ вывода (например, в файл или сеть)
//   - Единая точка входа для всех выводов
func Print(problems []models.Problem, formatter Formatter) {
	fmt.Print(formatter.Format(problems))
}
