package models

type LevelProblem string

const (
	Low    LevelProblem = "LOW"
	Medium LevelProblem = "MEDIUM"
	High   LevelProblem = "HIGH"
)

type Problem struct {
	Path           string       `json:"path"`           //принимаем путь к файлу
	ParseV         string       `json:"explanation"`    //распарсим конфиг
	Rule           string       `json:"rule"`           //проверка по заданному набору правил
	LevelProblem   LevelProblem `json:"level_problem"`  //уровень найденных проблем
	Recommendation string       `json:"recommendation"` //краткое объяснение и рекомендации
}
