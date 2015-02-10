package tracks

// Provider implementation for unit tests

import "strconv"

type providerMock struct {
	prefix  string
	current int
}

func (mock *providerMock) Provide() TrackId {
	defer func() { mock.current++ }()
	next := mock.prefix + strconv.Itoa(mock.current)
	return TrackId(next)
}

func (mock *providerMock) Prefix() string {
	return mock.prefix
}
