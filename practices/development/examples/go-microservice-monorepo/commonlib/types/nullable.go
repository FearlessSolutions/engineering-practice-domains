package types

import (
	"encoding/json"
	"errors"
)

// Nullable is a generic type that allows us to represent values that are or are not present. It can be used in JSON DTOs,
// as it has a custom implementation for json.Marshal and json.Unmarshal to handle null values.
type Nullable[T any] struct {
	// Value is the underlying nullable value. It can be read if IsPresent is true.
	Value T
	// IsPresent states the presence of the value. If it is true, the value is present.
	IsPresent bool
}

// UnmarshalJSON implements json.Unmarshaler for Nullable
//
//goland:noinspection GoMixedReceiverTypes
func (nullable *Nullable[T]) UnmarshalJSON(bytes []byte) error {
	var zeroValue T
	if string(bytes) != "null" {
		if unmarshalErr := json.Unmarshal(bytes, &zeroValue); unmarshalErr != nil {
			return unmarshalErr
		}

		nullable.Value = zeroValue
		nullable.IsPresent = true
	}

	return nil
}

// MarshalJSON implements json.Marshaler for Nullable, and intentionally uses a value receiver so that this method is
// invoked during serialization even if this nullable in question is not a pointer. If we were to use a pointer receiver
// here, the inner "value" and "isPresent" fields would get serialized if this Nullable is passed to json.Marshal as
// a value.
//
//goland:noinspection GoMixedReceiverTypes
func (nullable Nullable[T]) MarshalJSON() ([]byte, error) {
	if !nullable.IsPresent {
		return []byte("null"), nil
	}

	return json.Marshal(nullable.Value)
}

// NullablePresent is a validation rule satisfying validation.Rule. It fails validation on a nullable field if
// the nullable value is not present. This rule is generally useful in cases where a nullable field should
// be required conditionally, based on other fields present in the data structure
type NullablePresent[NestedValue any] struct{}

// Validate implements validation.Rule for NullablePresent
func (rule NullablePresent[NestedValue]) Validate(value interface{}) error {
	var nullableValue *Nullable[NestedValue]
	switch realValue := value.(type) {
	case Nullable[NestedValue]:
		nullableValue = &realValue
	case *Nullable[NestedValue]:
		nullableValue = realValue
	default:
		panic("the dtos.NullablePresent validation rule only works with Nullable[T] types, and the NestedValue generic must match the type T, please use is.Required instead or make the validation type match that of the nullable")
	}

	if !nullableValue.IsPresent {
		return errors.New("nullable value was not present, but should have been")
	}

	return nil
}
