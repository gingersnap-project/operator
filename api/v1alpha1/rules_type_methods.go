package v1alpha1

import (
	"fmt"
)

func (x DBType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", DBType_name[int32(x)])), nil
}

func (x *DBType) UnmarshalJSON(b []byte) error {
	*x = DBType(DBType_value[string(b[1:len(b)-1])])
	return nil
}

func (x KeyFormat) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", KeyFormat_name[int32(x)])), nil
}

func (x *KeyFormat) UnmarshalJSON(b []byte) error {
	*x = KeyFormat(KeyFormat_value[string(b[1:len(b)-1])])
	return nil
}
