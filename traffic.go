package core

import (
	"encoding/json"
	"time"
)

type Traffic struct {
	Date     time.Time `json:"date"`
	DataType string    `json:"dataType"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Value    int64     `json:"value"`
}

func (t *Traffic) UnmarshalJSON(p []byte) error {
	var aux struct {
		Date     string `json:"date"`
		DataType string `json:"dataType,string"`
		Name     string `json:"name,string"`
		Type     string `json:"type,string"`
		Value    int64  `json:"value,int"`
	}

	err := json.Unmarshal(p, &aux)
	if err != nil {
		return err
	}

	foo, err := time.Parse(time.RFC3339, aux.Date)
	if err != nil {
		return err
	}

	t.Name = aux.Name
	t.Date = foo
	return nil
}
