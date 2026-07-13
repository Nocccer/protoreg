package tests_test

import (
	"testing"

	"github.com/nocccer/protoreg/tests"
	"github.com/stretchr/testify/suite"
)

func TestGeneratorTags(t *testing.T) {
	suite.Run(t, new(GeneratorTagsTestSuite))
}

type GeneratorTagsTestSuite struct {
	suite.Suite
}

func (s *GeneratorTagsTestSuite) TestCustomTagKey() {
	custom := tests.CustomStructKey{}

	intf, ok := any(&custom).(interface {
		Marshal() ([]uint16, error)
		Unmarshal([]uint16) error
	})
	s.True(ok)
	s.NotNil(intf)
}

func (s *GeneratorTagsTestSuite) TestOnlyMarshaler() {
	onlyMarshaler := tests.OnlyMarshaler{}

	intf1, ok := any(&onlyMarshaler).(interface {
		Marshal() ([]uint16, error)
	})
	s.True(ok)
	s.NotNil(intf1)

	intf2, ok := any(&onlyMarshaler).(interface {
		Unmarshal([]uint16) error
	})
	s.False(ok)
	s.Nil(intf2)
}

func (s *GeneratorTagsTestSuite) TestOnlyUnmarshaler() {
	onlyUnmarshaler := tests.OnlyUnmarshaler{}

	intf1, ok := any(&onlyUnmarshaler).(interface {
		Marshal() ([]uint16, error)
	})
	s.False(ok)
	s.Nil(intf1)

	intf2, ok := any(&onlyUnmarshaler).(interface {
		Unmarshal([]uint16) error
	})
	s.True(ok)
	s.NotNil(intf2)
}

func (s *GeneratorTagsTestSuite) TestCustomFuncNames() {
	custom := tests.CustomFuncNames{}

	intf1, ok := any(&custom).(interface {
		Marshal() ([]uint16, error)
		Unmarshal([]uint16) error
	})
	s.False(ok)
	s.Nil(intf1)

	intf2, ok := any(&custom).(interface {
		CustomMarshal() ([]uint16, error)
		CustomUnmarshal([]uint16) error
	})
	s.True(ok)
	s.NotNil(intf2)
}
