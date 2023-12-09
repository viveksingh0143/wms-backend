package types

import (
	"encoding/json"
	"strings"
)

type FillStatus int

const (
	FillStatusEmpty FillStatus = iota + 1
	FillStatusFull
)

func (s *FillStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *FillStatus) UnmarshalJSON(data []byte) error {
	if err := s.unmarshalFromInt(data, s); err == nil {
		return nil
	}
	if err := s.unmarshalFromString(data, s); err != nil {
		return err
	}
	return nil
}

func (s *FillStatus) unmarshalFromInt(data []byte, status *FillStatus) error {
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		*status = FillStatus(intVal)
		return nil
	}
	return json.Unmarshal(data, &intVal)
}

func (s *FillStatus) unmarshalFromString(data []byte, status *FillStatus) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*status = s.getFillStatusFromString(str)
	return nil
}

func (s *FillStatus) getFillStatusFromString(str string) FillStatus {
	switch strings.ToLower(str) {
	case "empty":
		return FillStatusEmpty
	case "full":
		return FillStatusFull
	default:
		return 0 // or some default value
	}
}

func (s *FillStatus) String() string {
	switch *s {
	case FillStatusEmpty:
		return "EMPTY"
	case FillStatusFull:
		return "FULL"
	default:
		return "UNKNOWN"
	}
}
