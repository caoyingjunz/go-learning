package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type Checker interface {
	Check() (warnings, errorList []error)
	Name() string
}

type FileAvailableCheck struct {
	Path  string
	Label string
}

func (fac FileAvailableCheck) Name() string {
	if len(fac.Label) != 0 {
		return fac.Label
	}
	return "TestDir"
}

func (fac FileAvailableCheck) Check() (warnings, errorList []error) {
	if _, err := os.Stat(fac.Path); err == nil {
		return nil, []error{errors.Errorf("%s already exists", fac.Path)}
	}
	return nil, nil
}

type TestCheck struct {
	Test string
}

func (t TestCheck) Name() string {
	return "Test"
}

func (t TestCheck) Check() (warnings, errorList []error) {
	// TODO
	return nil, nil
}

// RunChecks runs each check, displays it's warnings/errors, and once all
// are processed will exit if any errors occurred.
func RunChecks(checks []Checker) error {
	for _, c := range checks {
		name := c.Name()
		warnings, errs := c.Check()

		fmt.Println("name:", name)
		fmt.Println("warnings:", warnings)
		fmt.Println("errs:", errs)
	}

	return errors.Wrap(fmt.Errorf("wrap raw error"), "check failed")
}

// 接口实现检查 demo
func main() {
	checks := []Checker{
		FileAvailableCheck{Path: "/root/test1"},
		TestCheck{Test: "Test"},
	}

	err := RunChecks(checks)
	fmt.Println(err)
}
