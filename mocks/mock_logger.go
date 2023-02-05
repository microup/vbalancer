package mocks

type MockLogger struct {
}

func (m *MockLogger) Add(values ...interface{}) {

}

func (m *MockLogger) Close() error {
	return nil
}
