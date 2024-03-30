package testhelper

import (
	"encoding/json"
	"io"

	"github.com/stretchr/testify/suite"
)

// UnmarshalBody converts a raw response from a controller into a concrete type (unmarshalTarget). unmarshalTarget
// should be a pointer to a data structure so this function can write data into it. This function accepts the test
// suite as an argument, so it can auto-fail the test in the event unmarshalling or closing the body fails
func UnmarshalBody(ste *suite.Suite, responseBody io.ReadCloser, unmarshalTarget any) {
	decoder := json.NewDecoder(responseBody)
	defer func() {
		_ = responseBody.Close()
	}()

	parseErr := decoder.Decode(unmarshalTarget)
	ste.Require().NoError(parseErr)
}
