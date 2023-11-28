package sshx_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/tsivinsky/sshx/config"
)

type questionAnswer map[int]string
type promptMock struct {
	qa              questionAnswer
	questionTracker int
}

func (p *promptMock) Select(prompt string, defaultValue string, options []string) (int, error) {
	return 0, nil
}

func (p *promptMock) Input(prompt string, defaultValue string) (string, error) {
	p.questionTracker += 1
	return p.qa[p.questionTracker], nil
}

func (p *promptMock) MultiSelect(prompt string, defaultValues []string, options []string) ([]int, error) {
	return []int{}, nil
}

func TestAdd(t *testing.T) {
	t.Parallel()
	testFile, err := os.CreateTemp("", "*_test_config.json")

	if err != nil {
		t.Errorf("error creating tempFile for testing purposes")
	}
	defer os.Remove(testFile.Name())

	testConf, err := config.NewConfig(config.WithFile(testFile.Name()))
	if err != nil {
		t.Errorf("test add: error creating testConf")
	}
	testConf.Servers = []config.Server{{Name: "s1", User: "u1", Host: "h1"}}

	pm := &promptMock{questionTracker: 0, qa: questionAnswer{1: "s2", 2: "u2", 3: "h2"}}

	want := []config.Server{
		{Name: "s1", User: "u1", Host: "h1"},
		{Name: "s2", User: "u2", Host: "h2"}}

	err = testConf.Add(pm)
	if err != nil {
		t.Errorf("error adding config with mockPrompt")
	}
	if got := testConf.Servers; !reflect.DeepEqual(got, want) {
		t.Fatalf("not adding config correctly:\ngot: %v, want: %v", got, want)
	}

}
