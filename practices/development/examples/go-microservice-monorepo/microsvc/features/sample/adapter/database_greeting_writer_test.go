package adapter

import (
	"context"
	"errors"
	"example.com/sample/commonlib/database"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
)

type DatabaseGreetingWriterSuite struct {
	suite.Suite
	mockController *gomock.Controller
	mockConnection *database.MockConnection
	connContext    context.Context
}

func TestDatabaseGreetingWriterSuite(t *testing.T) {
	suite.Run(t, new(DatabaseGreetingWriterSuite))
}

func (suite *DatabaseGreetingWriterSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockConnection = database.NewMockConnection(suite.mockController)
	suite.connContext = database.CreateDerivativeMockContext(context.Background(), suite.mockConnection)
}

func (suite *DatabaseGreetingWriterSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *DatabaseGreetingWriterSuite) TestAddGreeting() {
	suite.mockConnection.EXPECT().
		Exec(gomock.Any(), "G'day").
		Return(nil, nil)

	addErr := DatabaseGreetingWriter{}.AddGreeting(suite.connContext, "G'day")
	suite.Assert().NoError(addErr)
}

func (suite *DatabaseGreetingWriterSuite) TestAddGreetingFailsOnDbFail() {
	expectedErr := errors.New("oops")
	suite.mockConnection.EXPECT().Exec(gomock.Any(), "G'day").Return(nil, expectedErr)

	addErr := DatabaseGreetingWriter{}.AddGreeting(suite.connContext, "G'day")
	suite.Assert().ErrorIs(addErr, expectedErr)
}
