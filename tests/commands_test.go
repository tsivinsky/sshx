package sshx_test

import (
	"reflect"
	"testing"

	"github.com/tsivinsky/sshx/config"
)

type questionAnswer map[int]string

type promptMock struct {
	qa                questionAnswer
	questionTracker   int
	MultiSelectAnswer []int
}

func (p *promptMock) Select(prompt string, defaultValue string, options []string) (int, error) {
	return 0, nil
}

func (p *promptMock) Input(prompt string, defaultValue string) (string, error) {
	p.questionTracker += 1
	return p.qa[p.questionTracker], nil
}

func (p *promptMock) MultiSelect(prompt string, defaultValues []string, options []string) ([]int, error) {
	return p.MultiSelectAnswer, nil
}

func TestAdd(t *testing.T) {
	t.Parallel()
	// test set up and tearDown
	testFile, testConf, err := setUp()
	defer tearDown(testFile)
	if err != nil {
		t.Errorf("error setting up test scenario: %v", err)
	}
	// populate test data
	testConf.Servers = []config.Server{{Name: "s1", User: "u1", Host: "h1"}}

	pm := &promptMock{questionTracker: 0, qa: questionAnswer{1: "s2", 2: "u2", 3: "h2"}}

	err = testConf.Add(pm)

	want := []config.Server{
		{Name: "s1", User: "u1", Host: "h1"},
		{Name: "s2", User: "u2", Host: "h2"}}

	if err != nil {
		t.Errorf("error adding config with mockPrompt")
	}
	if got := testConf.Servers; !reflect.DeepEqual(got, want) {
		t.Fatalf("not adding config correctly:\ngot: %v, want: %v", got, want)
	}

}

func TestRemove(t *testing.T) {
	t.Parallel()
	// test set up and tearDown
	testFile, testConf, err := setUp()
	defer tearDown(testFile)
	if err != nil {
		t.Errorf("error setting up test scenario: %v", err)
	}
	// populate test data
	testConf.Servers = []config.Server{
		{Name: "s1", User: "u1", Host: "h1"},
		{Name: "s2", User: "u2", Host: "h2"},
		{Name: "s3", User: "u3", Host: "h3"},
		{Name: "s4", User: "u4", Host: "h4"}}
	// populate mockPrompt with expected answers
	// remove s2 and s4
	mp := &promptMock{MultiSelectAnswer: []int{1, 3}}
	err = testConf.Remove(mp)
	if err != nil {
		t.Errorf("TestRemove: error running TestConf.Remove(mp).\nerror: %v", err)
	}
	want := []config.Server{
		{Name: "s1", User: "u1", Host: "h1"},
		{Name: "s3", User: "u3", Host: "h3"}}
	if got := testConf.Servers; !reflect.DeepEqual(got, want) {
		t.Fatalf("got: %v, want %v", got, want)
	}
}
