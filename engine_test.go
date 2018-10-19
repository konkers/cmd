package cmd

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewEngine(t *testing.T) {
	e := NewEngine()
	if e == nil {
		t.Errorf("NewEngine() returned nil")
	}
}

func TestAddHandler(t *testing.T) {
	e := NewEngine()
	h := func(ctx interface{}, args []string) error { return nil }

	err := e.AddCommand("test", "help", h, 0)
	if err != nil {
		t.Errorf("AddCommand() failed: %v", err)
	}

	err = e.AddCommand("test", "help", h, 0)
	if err == nil {
		t.Errorf("Duplicate AddCommand() succeeded")
	}
}

func TestRemoveHandler(t *testing.T) {
	e := NewEngine()
	h := func(ctx interface{}, args []string) error { return nil }

	err := e.RemoveCommand("test")
	if err == nil {
		t.Errorf("Expected error when removing non-existant command")
	}

	err = e.AddCommand("test", "help", h, 0)
	if err != nil {
		t.Errorf("AddCommand() failed: %v", err)
	}

	err = e.RemoveCommand("test")
	if err != nil {
		t.Errorf("RemoveCommand() failed: %v", err)
	}

	err = e.ExecString(nil, 10, "test")
	if err == nil {
		t.Errorf("Expected error when executing non-existant command")
	}
}

func TestExecHandler(t *testing.T) {
	var handlerArgs []string
	var err error
	h := func(ctx interface{}, args []string) error {
		handlerArgs = args
		if len(args) >= 1 && args[0] == "bad" {
			return fmt.Errorf("bad")
		} else {
			return nil
		}
	}

	e := NewEngine()
	e.AddCommand("test", "help", h, 0)

	type testdata struct {
		ctx           interface{}
		userLevel     int
		args          []string
		errorExpected bool
		argsExpected  bool
	}
	tests := []testdata{
		{nil, 0, []string{}, true, false},               // 0
		{nil, 0, []string{"nope"}, true, false},         // 1
		{nil, -1, []string{"test"}, true, false},        // 2
		{nil, 0, []string{"test"}, false, true},         // 3
		{nil, 0, []string{"test", "good"}, false, true}, // 4
		{nil, 0, []string{"test", "bad"}, true, true},   // 5
	}

	for i, test := range tests {
		handlerArgs = nil
		err = e.Exec(test.ctx, test.userLevel, test.args)
		if test.errorExpected && err == nil {
			t.Errorf("%d: Expected error but none was returned.", i)
		} else if !test.errorExpected && err != nil {
			t.Errorf("%d: Expected no error but one was returned.", i)
		}

		if test.argsExpected {
			if !reflect.DeepEqual(handlerArgs, test.args[1:]) {
				t.Errorf("%d: handlerArgs(%#v) != exepected(%#v)", i, handlerArgs, test.args[1:])
			}
		} else {
			if handlerArgs != nil {
				t.Errorf("%d: expected handlerArgs be nil, instead they are %#v",
					i, handlerArgs)
			}
		}
	}
}

func TestExecString(t *testing.T) {
	var handlerArgs []string
	var err error
	h := func(ctx interface{}, args []string) error {
		handlerArgs = args
		if len(args) >= 1 && args[0] == "bad" {
			return fmt.Errorf("bad")
		} else {
			return nil
		}
	}

	e := NewEngine()
	e.AddCommand("test", "help", h, 0)

	type testdata struct {
		ctx           interface{}
		userLevel     int
		args          string
		errorExpected bool
		expectedArgs  []string
	}
	tests := []testdata{
		{nil, 0, "\"", true, nil},                      // 0
		{nil, 0, "", true, nil},                        // 1
		{nil, 0, "nope", true, nil},                    // 2
		{nil, -1, "test", true, nil},                   // 3
		{nil, 0, "test", false, []string{}},            // 4
		{nil, 0, "test good", false, []string{"good"}}, // 5
		{nil, 0, "test bad", true, []string{"bad"}},    // 6
	}

	for i, test := range tests {
		handlerArgs = nil
		err = e.ExecString(test.ctx, test.userLevel, test.args)
		if test.errorExpected && err == nil {
			t.Errorf("%d: Expected error but none was returned.", i)
		} else if !test.errorExpected && err != nil {
			t.Errorf("%d: Expected no error but one was returned.", i)
		}

		if test.expectedArgs != nil {
			if !reflect.DeepEqual(handlerArgs, test.expectedArgs) {
				t.Errorf("%d: handlerArgs(%#v) != exepected(%#v)",
					i, handlerArgs, test.expectedArgs)
			}
		} else {
			if handlerArgs != nil {
				t.Errorf("%d: expected handlerArgs be nil, instead they are %#v",
					i, handlerArgs)
			}
		}
	}
}
