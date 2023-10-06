package types

import (
	"encoding/json"
	"strings"
)

type Status int

const (
	StatusActive Status = iota + 1
	StatusInactive
	StatusBanned
)

func (s *Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Status) UnmarshalJSON(data []byte) error {
	if err := unmarshalFromInt(data, s); err == nil {
		return nil
	}
	if err := unmarshalFromString(data, s); err != nil {
		return err
	}
	return nil
}

func unmarshalFromInt(data []byte, s *Status) error {
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		*s = Status(intVal)
		return nil
	}
	return json.Unmarshal(data, &intVal)
}

func unmarshalFromString(data []byte, s *Status) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*s = getStatusFromString(str)
	return nil
}

func getStatusFromString(str string) Status {
	switch strings.ToLower(str) {
	case "active":
		return StatusActive
	case "inactive":
		return StatusInactive
	case "banned":
		return StatusBanned
	default:
		return 0 // or some default value
	}
}

func (s *Status) String() string {
	switch *s {
	case StatusActive:
		return "ACTIVE"
	case StatusInactive:
		return "INACTIVE"
	case StatusBanned:
		return "BANNED"
	default:
		return "UNKNOWN"
	}
}
