package sshx_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tsivinsky/sshx/config"
)

// dict to simulate prompt input reply
// index is #question
// the index is tracked on questionTracker
type questionAnswer map[int]string

type promptMock struct {
	qa                 questionAnswer
	questionTracker    int
	multiSelectAnswer  []int
	simpleSelectAnswer int
}

func (p *promptMock) Select(prompt string, defaultValue string, options []string) (int, error) {
	return p.simpleSelectAnswer, nil
}

func (p *promptMock) Input(prompt string, defaultValue string) (string, error) {
	p.questionTracker += 1
	return p.qa[p.questionTracker], nil
}

func (p *promptMock) MultiSelect(prompt string, defaultValues []string, options []string) ([]int, error) {
	return p.multiSelectAnswer, nil
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
	mp := &promptMock{multiSelectAnswer: []int{1, 3}}
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

func TestUpdate(t *testing.T) {
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
	pm := &promptMock{
		questionTracker:    0,
		qa:                 questionAnswer{1: "s2B", 2: "u2B", 3: "h2B"},
		simpleSelectAnswer: 1}
	// execute tested function
	err = testConf.Update(pm)
	if err != nil {
		t.Errorf("TestUpdate - error executing testConf.Update: %v", err)
	}
	want := []config.Server{
		{Name: "s1", User: "u1", Host: "h1"},
		{Name: "s2B", User: "u2B", Host: "h2B"},
		{Name: "s3", User: "u3", Host: "h3"},
		{Name: "s4", User: "u4", Host: "h4"}}
	if got := testConf.Servers; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestList(t *testing.T) {
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
	got := new(bytes.Buffer)
	testConf.List(got)

	want := "s1: u1@h1\ns2: u2@h2\ns3: u3@h3\ns4: u4@h4\n"
	if err != nil {
		t.Errorf("TestList: error marshalling testConf.Servers")
	}
	if got.String() != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}
