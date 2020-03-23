package mocks

type Manager struct {
	NotifyFunc func(string) error
	CloseFunc  func() error
}

func (m *Manager) Notify(s string) error {
	return m.NotifyFunc(s)
}

func (m *Manager) Close() error {
	return m.CloseFunc()
}
