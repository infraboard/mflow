// Code generated by github.com/infraboard/mcube/v2
// DO NOT EDIT

package pipeline

import (
	"bytes"
	"fmt"
	"strings"
)

// ParseTRIGGER_MODEFromString Parse TRIGGER_MODE from string
func ParseTRIGGER_MODEFromString(str string) (TRIGGER_MODE, error) {
	key := strings.Trim(string(str), `"`)
	v, ok := TRIGGER_MODE_value[strings.ToUpper(key)]
	if !ok {
		return 0, fmt.Errorf("unknown TRIGGER_MODE: %s", str)
	}

	return TRIGGER_MODE(v), nil
}

// Equal type compare
func (t TRIGGER_MODE) Equal(target TRIGGER_MODE) bool {
	return t == target
}

// IsIn todo
func (t TRIGGER_MODE) IsIn(targets ...TRIGGER_MODE) bool {
	for _, target := range targets {
		if t.Equal(target) {
			return true
		}
	}

	return false
}

// MarshalJSON todo
func (t TRIGGER_MODE) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(strings.ToUpper(t.String()))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON todo
func (t *TRIGGER_MODE) UnmarshalJSON(b []byte) error {
	ins, err := ParseTRIGGER_MODEFromString(string(b))
	if err != nil {
		return err
	}
	*t = ins
	return nil
}

// ParseEVENT_LEVELFromString Parse EVENT_LEVEL from string
func ParseEVENT_LEVELFromString(str string) (EVENT_LEVEL, error) {
	key := strings.Trim(string(str), `"`)
	v, ok := EVENT_LEVEL_value[strings.ToUpper(key)]
	if !ok {
		return 0, fmt.Errorf("unknown EVENT_LEVEL: %s", str)
	}

	return EVENT_LEVEL(v), nil
}

// Equal type compare
func (t EVENT_LEVEL) Equal(target EVENT_LEVEL) bool {
	return t == target
}

// IsIn todo
func (t EVENT_LEVEL) IsIn(targets ...EVENT_LEVEL) bool {
	for _, target := range targets {
		if t.Equal(target) {
			return true
		}
	}

	return false
}

// MarshalJSON todo
func (t EVENT_LEVEL) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(strings.ToUpper(t.String()))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON todo
func (t *EVENT_LEVEL) UnmarshalJSON(b []byte) error {
	ins, err := ParseEVENT_LEVELFromString(string(b))
	if err != nil {
		return err
	}
	*t = ins
	return nil
}

// ParseAUDIT_STAGEFromString Parse AUDIT_STAGE from string
func ParseAUDIT_STAGEFromString(str string) (AUDIT_STAGE, error) {
	key := strings.Trim(string(str), `"`)
	v, ok := AUDIT_STAGE_value[strings.ToUpper(key)]
	if !ok {
		return 0, fmt.Errorf("unknown AUDIT_STAGE: %s", str)
	}

	return AUDIT_STAGE(v), nil
}

// Equal type compare
func (t AUDIT_STAGE) Equal(target AUDIT_STAGE) bool {
	return t == target
}

// IsIn todo
func (t AUDIT_STAGE) IsIn(targets ...AUDIT_STAGE) bool {
	for _, target := range targets {
		if t.Equal(target) {
			return true
		}
	}

	return false
}

// MarshalJSON todo
func (t AUDIT_STAGE) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(strings.ToUpper(t.String()))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON todo
func (t *AUDIT_STAGE) UnmarshalJSON(b []byte) error {
	ins, err := ParseAUDIT_STAGEFromString(string(b))
	if err != nil {
		return err
	}
	*t = ins
	return nil
}
