package types

import (
	"encoding/json"
	"strings"
)

type ProcessStatus int

const (
	ProcessNotStarted ProcessStatus = iota + 1
	ProcessStarted
	ProcessClosed
	ProcessRejected
)

func (ps *ProcessStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(ps.String())
}

func (ps *ProcessStatus) UnmarshalJSON(data []byte) error {
	if err := ps.unmarshalFromInt(data, ps); err == nil {
		return nil
	}
	if err := ps.unmarshalFromString(data, ps); err != nil {
		return err
	}
	return nil
}

func (ps *ProcessStatus) unmarshalFromInt(data []byte, status *ProcessStatus) error {
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		*status = ProcessStatus(intVal)
		return nil
	}
	return json.Unmarshal(data, &intVal)
}

func (ps *ProcessStatus) unmarshalFromString(data []byte, status *ProcessStatus) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*status = ps.getStatusFromString(str)
	return nil
}

func (s *ProcessStatus) getStatusFromString(str string) ProcessStatus {
	switch strings.ToLower(str) {
	case "not started":
		return ProcessNotStarted
	case "started":
		return ProcessStarted
	case "closed":
		return ProcessClosed
	case "rejected":
		return ProcessRejected
	default:
		return 0 // or some default value
	}
}

func (s *ProcessStatus) String() string {
	switch *s {
	case ProcessNotStarted:
		return "NOT STARTED"
	case ProcessStarted:
		return "STARTED"
	case ProcessClosed:
		return "CLOSED"
	case ProcessRejected:
		return "REJECTED"
	default:
		return "UNKNOWN"
	}
}
