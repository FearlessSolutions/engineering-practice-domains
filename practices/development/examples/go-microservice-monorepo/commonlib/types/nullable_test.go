package types

import (
	"encoding/json"
	"github.com/jellydator/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type InnerStruct struct {
	ValueOne string `json:"valueOne"`
	ValueTwo int    `json:"valueTwo"`
}

type SampleDataStruct struct {
	NumberTest Nullable[int]         `json:"numberTest"`
	StringTest Nullable[string]      `json:"stringTest"`
	StructTest Nullable[InnerStruct] `json:"structTest"`
}

func TestNullable_UnmarshalFromNull(t *testing.T) {
	testPayload := `{"numberTest":null,"stringTest":null,"structTest":null}`
	var targetStruct SampleDataStruct
	unmarshalErr := json.Unmarshal([]byte(testPayload), &targetStruct)
	require.NoError(t, unmarshalErr)

	assert.False(t, targetStruct.NumberTest.IsPresent)
	assert.False(t, targetStruct.StringTest.IsPresent)
	assert.False(t, targetStruct.StructTest.IsPresent)
}

func TestNullable_UnmarshalFromEmpty(t *testing.T) {
	testPayload := `{}`
	var targetStruct SampleDataStruct
	unmarshalErr := json.Unmarshal([]byte(testPayload), &targetStruct)
	require.NoError(t, unmarshalErr)

	assert.False(t, targetStruct.NumberTest.IsPresent)
	assert.False(t, targetStruct.StringTest.IsPresent)
	assert.False(t, targetStruct.StructTest.IsPresent)
}

func TestNullable_UnmarshalFromGoodValue(t *testing.T) {
	testPayload := `
{
  "numberTest": 5,
  "stringTest": "hello",
  "structTest": {
    "valueOne": "goodbye",
    "valueTwo": 20
  }
}
`
	var targetStruct SampleDataStruct
	unmarshalErr := json.Unmarshal([]byte(testPayload), &targetStruct)
	require.NoError(t, unmarshalErr)

	expectedStruct := SampleDataStruct{
		NumberTest: Nullable[int]{
			Value:     5,
			IsPresent: true,
		},
		StringTest: Nullable[string]{
			Value:     "hello",
			IsPresent: true,
		},
		StructTest: Nullable[InnerStruct]{
			IsPresent: true,
			Value: InnerStruct{
				ValueOne: "goodbye",
				ValueTwo: 20,
			},
		},
	}
	assert.Equal(t, expectedStruct, targetStruct)
}

func TestNullable_MarshalsToNullValues(t *testing.T) {
	cases := []struct {
		Name  string
		Value any
	}{
		{
			Name:  "Pointer to struct",
			Value: &SampleDataStruct{},
		},
		// It turns out if you pass a struct by value and its MarshalJSON function has a pointer receiver,
		// json.Marshal will bypass the custom marshal implementation because the Go runtime can't acquire a
		// pointer to the struct passed by value
		{
			Name:  "Struct Value",
			Value: SampleDataStruct{},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			jsonData, marshalErr := json.Marshal(testCase.Value)
			require.NoError(t, marshalErr)

			expectedJson := `{"numberTest":null,"stringTest":null,"structTest":null}`
			assert.Equal(t, expectedJson, string(jsonData))
		})
	}
}

func TestNullable_MarshalsPresentValues(t *testing.T) {
	testStruct := SampleDataStruct{
		NumberTest: Nullable[int]{
			Value:     5,
			IsPresent: true,
		},
		StringTest: Nullable[string]{
			Value:     "hello",
			IsPresent: true,
		},
		StructTest: Nullable[InnerStruct]{
			Value: InnerStruct{
				ValueOne: "goodbye",
				ValueTwo: 10,
			},
			IsPresent: true,
		},
	}
	expectedJson := `{"numberTest":5,"stringTest":"hello","structTest":{"valueOne":"goodbye","valueTwo":10}}`
	cases := []struct {
		Name  string
		Value any
	}{
		{
			Name:  "Marshalling pointer to struct",
			Value: &testStruct,
		},
		{
			Name:  "Marshalling struct value",
			Value: testStruct,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			jsonData, marshalErr := json.Marshal(testCase.Value)
			require.NoError(t, marshalErr)

			assert.Equal(t, expectedJson, string(jsonData))
		})
	}
}

type ConditionalValidation struct {
	ShouldValidate bool          `json:"shouldValidate"`
	Value          Nullable[int] `json:"value"`
}

// The NullablePresent validation is generally useful if you want to require a nullable in specific cases

func (cv ConditionalValidation) Validate() error {
	var valueRules []validation.Rule

	if cv.ShouldValidate {
		valueRules = append(valueRules, NullablePresent[int]{})
	}

	return validation.ValidateStruct(&cv,
		validation.Field(&cv.Value, valueRules...),
	)
}

func TestNullable_IsPresentRule(t *testing.T) {
	shouldPassValidation := ConditionalValidation{
		ShouldValidate: false,
		Value:          Nullable[int]{},
	}
	shouldFailValidation := ConditionalValidation{
		ShouldValidate: true,
		Value:          Nullable[int]{},
	}

	shouldPassErr := shouldPassValidation.Validate()
	shouldFailErr := shouldFailValidation.Validate()

	assert.NoError(t, shouldPassErr)
	assert.Error(t, shouldFailErr)
}
