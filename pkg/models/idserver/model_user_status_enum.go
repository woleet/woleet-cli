/*
Woleet.ID Server

This is Woleet.ID Server API documentation.

API version: 1.2.8
Contact: contact@woleet.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package idserver

import (
	"encoding/json"
	"fmt"
)

// UserStatusEnum User status:<br> - `active`: the user is active: he can use his keys to sign<br> - `blocked`: the user is blocked: he cannot use his keys to sign 
type UserStatusEnum string

// List of UserStatusEnum
const (
	USERSTATUSENUM_ACTIVE UserStatusEnum = "active"
	USERSTATUSENUM_BLOCKED UserStatusEnum = "blocked"
)

// All allowed values of UserStatusEnum enum
var AllowedUserStatusEnumEnumValues = []UserStatusEnum{
	"active",
	"blocked",
}

func (v *UserStatusEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := UserStatusEnum(value)
	for _, existing := range AllowedUserStatusEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid UserStatusEnum", value)
}

// NewUserStatusEnumFromValue returns a pointer to a valid UserStatusEnum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewUserStatusEnumFromValue(v string) (*UserStatusEnum, error) {
	ev := UserStatusEnum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for UserStatusEnum: valid values are %v", v, AllowedUserStatusEnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v UserStatusEnum) IsValid() bool {
	for _, existing := range AllowedUserStatusEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to UserStatusEnum value
func (v UserStatusEnum) Ptr() *UserStatusEnum {
	return &v
}

type NullableUserStatusEnum struct {
	value *UserStatusEnum
	isSet bool
}

func (v NullableUserStatusEnum) Get() *UserStatusEnum {
	return v.value
}

func (v *NullableUserStatusEnum) Set(val *UserStatusEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableUserStatusEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableUserStatusEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserStatusEnum(val *UserStatusEnum) *NullableUserStatusEnum {
	return &NullableUserStatusEnum{value: val, isSet: true}
}

func (v NullableUserStatusEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserStatusEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

