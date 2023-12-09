package types

import (
	"encoding/json"
	"strings"
)

type Approval int

const (
	ApprovalYes Approval = iota + 1
	ApprovalNo
	ApprovalWait
)

func (s *Approval) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Approval) UnmarshalJSON(data []byte) error {
	if err := s.unmarshalFromInt(data, s); err == nil {
		return nil
	}
	if err := s.unmarshalFromString(data, s); err != nil {
		return err
	}
	return nil
}

func (s *Approval) unmarshalFromInt(data []byte, status *Approval) error {
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		*status = Approval(intVal)
		return nil
	}
	return json.Unmarshal(data, &intVal)
}

func (s *Approval) unmarshalFromString(data []byte, status *Approval) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*status = s.getApprovalFromString(str)
	return nil
}

func (s *Approval) getApprovalFromString(str string) Approval {
	switch strings.ToLower(str) {
	case "yes":
		return ApprovalYes
	case "no":
		return ApprovalNo
	case "wait":
		return ApprovalWait
	default:
		return 0 // or some default value
	}
}

func (s *Approval) String() string {
	switch *s {
	case ApprovalYes:
		return "YES"
	case ApprovalNo:
		return "NO"
	case ApprovalWait:
		return "WAIT"
	default:
		return "UNKNOWN"
	}
}
