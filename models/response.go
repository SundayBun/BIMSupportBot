package models

type ResponseModel struct {
	Id       int    `json:"id,omitempty" db:"id"`
	Answer   string `json:"answer"  db:"b_answer"`
	Question string `json:"question" db:"b_question"`
}
