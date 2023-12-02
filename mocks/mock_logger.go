package mocks

type MockLogger struct {
}

func (m *MockLogger) Add(...interface{}) {

}

func (m *MockLogger) Close() error {
	return nil
}
