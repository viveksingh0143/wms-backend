package types

import (
	"encoding/json"
	"strings"
)

type InventoryStatus int

const (
	InventoryIn InventoryStatus = iota + 1
	InventoryOut
	InventoryRejected
	InventoryExpired
	InventoryApprovalWait
)

func (ps *InventoryStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(ps.String())
}

func (ps *InventoryStatus) UnmarshalJSON(data []byte) error {
	if err := ps.unmarshalFromInt(data, ps); err == nil {
		return nil
	}
	if err := ps.unmarshalFromString(data, ps); err != nil {
		return err
	}
	return nil
}

func (ps *InventoryStatus) unmarshalFromInt(data []byte, status *InventoryStatus) error {
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		*status = InventoryStatus(intVal)
		return nil
	}
	return json.Unmarshal(data, &intVal)
}

func (ps *InventoryStatus) unmarshalFromString(data []byte, status *InventoryStatus) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*status = ps.getStatusFromString(str)
	return nil
}

func (s *InventoryStatus) getStatusFromString(str string) InventoryStatus {
	switch strings.ToLower(str) {
	case "in":
		return InventoryIn
	case "out":
		return InventoryOut
	case "rejected":
		return InventoryRejected
	case "expired":
		return InventoryExpired
	case "waiting-approval":
		return InventoryApprovalWait
	default:
		return 0 // or some default value
	}
}

func (s *InventoryStatus) String() string {
	switch *s {
	case InventoryIn:
		return "IN"
	case InventoryOut:
		return "OUT"
	case InventoryRejected:
		return "REJECTED"
	case InventoryExpired:
		return "EXPIRED"
	case InventoryApprovalWait:
		return "WAITING-APPROVAL"
	default:
		return "UNKNOWN"
	}
}
