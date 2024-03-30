package adapter

import (
	"context"
	"errors"
	"example.com/sample/commonlib/database"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"regexp"
	"testing"
)

type DatabaseGreetingReaderSuite struct {
	suite.Suite
	mockController *gomock.Controller
	mockConnection *database.MockConnection
	connectionCtx  context.Context
}

func TestDatabaseGreetingReaderSuite(t *testing.T) {
	suite.Run(t, new(DatabaseGreetingReaderSuite))
}

func (suite *DatabaseGreetingReaderSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockConnection = database.NewMockConnection(suite.mockController)
	suite.connectionCtx = database.CreateDerivativeMockContext(context.Background(), suite.mockConnection)
}

func (suite *DatabaseGreetingReaderSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *DatabaseGreetingReaderSuite) TestRandomGreetingRetrievesRandomEntry() {
	suite.mockConnection.EXPECT().
		Select(gomock.Any(), gomock.Any()).
		SetArg(0, []string{"A", "B", "C"}).
		Return(nil)

	greeting, greetingRetrieveErr := DatabaseGreetingReader{}.RandomGreeting(suite.connectionCtx)

	suite.Require().NoError(greetingRetrieveErr)
	// This is the best way I can think of to assert the greeting is either "A", "B", or "C"
	suite.Assert().Regexp(regexp.MustCompile("^[ABC]$"), greeting)
}

func (suite *DatabaseGreetingReaderSuite) TestRandomGreetingReturnsErrorOnDbError() {
	suite.mockConnection.EXPECT().
		Select(gomock.Any(), gomock.Any()).
		Return(errors.New("oops"))

	_, greetingRetrieveErr := DatabaseGreetingReader{}.RandomGreeting(suite.connectionCtx)
	suite.Require().Error(greetingRetrieveErr)
}
