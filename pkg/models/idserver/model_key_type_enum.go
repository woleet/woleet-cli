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

// KeyTypeEnum The type of the key (only `bip39` is supported).<br> - `bip39`: Bitcoin BIP39 
type KeyTypeEnum string

// List of KeyTypeEnum
const (
	KEYTYPEENUM_BIP39 KeyTypeEnum = "bip39"
)

// All allowed values of KeyTypeEnum enum
var AllowedKeyTypeEnumEnumValues = []KeyTypeEnum{
	"bip39",
}

func (v *KeyTypeEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := KeyTypeEnum(value)
	for _, existing := range AllowedKeyTypeEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid KeyTypeEnum", value)
}

// NewKeyTypeEnumFromValue returns a pointer to a valid KeyTypeEnum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewKeyTypeEnumFromValue(v string) (*KeyTypeEnum, error) {
	ev := KeyTypeEnum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for KeyTypeEnum: valid values are %v", v, AllowedKeyTypeEnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v KeyTypeEnum) IsValid() bool {
	for _, existing := range AllowedKeyTypeEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to KeyTypeEnum value
func (v KeyTypeEnum) Ptr() *KeyTypeEnum {
	return &v
}

type NullableKeyTypeEnum struct {
	value *KeyTypeEnum
	isSet bool
}

func (v NullableKeyTypeEnum) Get() *KeyTypeEnum {
	return v.value
}

func (v *NullableKeyTypeEnum) Set(val *KeyTypeEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableKeyTypeEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableKeyTypeEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKeyTypeEnum(val *KeyTypeEnum) *NullableKeyTypeEnum {
	return &NullableKeyTypeEnum{value: val, isSet: true}
}

func (v NullableKeyTypeEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKeyTypeEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

