package utilmock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type UtilMock struct {
	mock.Mock
}

func (m *UtilMock) DateToString(t time.Time) string {
	args := m.Called(t)
	return args.String(0)
}

func (m *UtilMock) StringToTime(s string) time.Time {
	args := m.Called(s)
	return args.Get(0).(time.Time)
}

func (m *UtilMock) InTimeSpanStart(start, end, check time.Time) bool {
	args := m.Called(start, end, check)
	return args.Bool(0)
}

func (m *UtilMock) InTimeSpanEnd(start, end, check time.Time) bool {
	args := m.Called(start, end, check)
	return args.Bool(0)
}
