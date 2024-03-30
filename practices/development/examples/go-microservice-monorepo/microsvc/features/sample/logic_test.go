package sample

import (
	"context"
	"errors"
	"example.com/sample/commonlib/logger"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zapcore"
	"testing"
)

type SampleLogicSuite struct {
	suite.Suite
	controller *gomock.Controller
}

func (suite *SampleLogicSuite) SetupSuite() {
	setupErr := logger.InitLogger(zapcore.InfoLevel, false)
	suite.Require().NoError(setupErr)
}

func (suite *SampleLogicSuite) SetupTest() {
	suite.controller = gomock.NewController(suite.T())
}

func (suite *SampleLogicSuite) TearDownTest() {
	suite.controller.Finish()
}

func TestSampleLogicSuite(t *testing.T) {
	suite.Run(t, new(SampleLogicSuite))
}

func (suite *SampleLogicSuite) TestGiveGreeting() {
	greetingReader := NewMockGreetingReader(suite.controller)

	greetingReader.EXPECT().RandomGreeting(gomock.Any()).Return("Hola", nil)

	greeting, greetingErr := CoreLogic{}.GiveGreeting(context.Background(), "Evan", greetingReader)
	suite.Assert().NoError(greetingErr)
	suite.Assert().Equal("Hola, Evan!", greeting)
}

func (suite *SampleLogicSuite) TestGiveGreetingFailsWhenReaderFails() {
	expectedError := errors.New("oops, something blew up")
	greetingReader := NewMockGreetingReader(suite.controller)

	greetingReader.EXPECT().RandomGreeting(gomock.Any()).Return("", expectedError)

	_, greetingErr := CoreLogic{}.GiveGreeting(context.Background(), "Evan", greetingReader)
	suite.Require().ErrorIs(greetingErr, expectedError)
}

func (suite *SampleLogicSuite) TestAddGreetingAddsGreeting() {
	greetingReader := NewMockGreetingReader(suite.controller)
	greetingWriter := NewMockGreetingWriter(suite.controller)

	listCall := greetingReader.EXPECT().
		List(gomock.Any()).
		Return([]string{"Hi", "Hello", "Howdy"}, nil)
	greetingWriter.EXPECT().
		AddGreeting(gomock.Any(), "G'day").
		After(listCall).
		Return(nil)

	addErr := CoreLogic{}.AddGreeting(context.Background(), "G'day", greetingReader, greetingWriter)
	suite.Assert().NoError(addErr)
}

func (suite *SampleLogicSuite) TestAddingDuplicateGreetingFails() {
	greetingReader := NewMockGreetingReader(suite.controller)
	greetingWriter := NewMockGreetingWriter(suite.controller)

	greetingReader.EXPECT().
		List(gomock.Any()).
		Return([]string{"Hi", "Hello", "Howdy"}, nil)

	addErr := CoreLogic{}.AddGreeting(context.Background(), "Hi", greetingReader, greetingWriter)
	suite.Assert().ErrorIs(addErr, ErrGreetingAlreadyExists)
}

func (suite *SampleLogicSuite) TestAddingGreetingFailsOnDbFail() {
	greetingReader := NewMockGreetingReader(suite.controller)
	greetingWriter := NewMockGreetingWriter(suite.controller)
	expectedErr := errors.New("oops")

	readerCall := greetingReader.EXPECT().
		List(gomock.Any()).
		Return([]string{"Hi", "Hello", "Howdy"}, nil)
	greetingWriter.EXPECT().
		AddGreeting(gomock.Any(), "G'day").
		After(readerCall).
		Return(expectedErr)

	addErr := CoreLogic{}.AddGreeting(context.Background(), "G'day", greetingReader, greetingWriter)
	suite.Assert().ErrorIs(addErr, expectedErr)
}
