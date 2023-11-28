package sshx_test

type promptMock struct {
}

func (p promptMock) Select(prompt string, defaultValue string, options []string) (int, error) {
	return 0, nil
}

func (p promptMock) Input(prompt string, defaultValue string) (string, error) {
	return "", nil
}
