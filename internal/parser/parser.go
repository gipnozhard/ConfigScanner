package parser

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Parser парсит конфигурации разных форматов
type Parser struct{}

// NewParser создает новый парсер
func NewParser() *Parser {
	return &Parser{}
}

// Parse парсит данные в map[string]interface{}
func (p *Parser) Parse(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}

	// Пробуем JSON
	err := json.Unmarshal(data, &result)
	if err == nil {
		return result, nil
	}

	// Пробуем YAML
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить: %v", err)
	}

	return result, nil
}

// DetectFormat определяет формат по содержимому
func (p *Parser) DetectFormat(data []byte) string {
	var test map[string]interface{}
	if json.Unmarshal(data, &test) == nil {
		return "json"
	}
	return "yaml"
}
