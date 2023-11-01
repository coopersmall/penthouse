package suite

import (
	"fmt"
	"testing"
)

var (
	currentSuite   *suite
	currentContext *context
	currentTest    *test
)

func Suite(name string) *suite {
	currentSuite = newSuite(name)
	return currentSuite
}

func (s *suite) BeforeAll(fn func()) *suite {
	s.opts = append(s.opts, func(ts *testingSuite) {
		ts.beforeAll = fn
	})
	return s
}

func (s *suite) AfterAll(fn func()) *suite {
	s.opts = append(s.opts, func(ts *testingSuite) {
		ts.afterAll = fn
	})
	return s
}

func (s *suite) SetupSuite(fn func()) *suite {
	s.opts = append(s.opts, func(s *testingSuite) {
		s.setupSuite = fn
	})
	return s
}

func (s *suite) TeardownSuite(fn func()) *suite {
	s.opts = append(s.opts, func(s *testingSuite) {
		s.teardownSuite = fn
	})
	return s
}

func Describe(name string, fn func()) *suite {
	if currentSuite == nil {
		panic("Describe must be called with an active suite")
	}

	currentContext = newContext(name)
	currentContext.parent = nil
	fn()

	currentSuite.contexts = append(currentSuite.contexts, currentContext)
	return currentSuite
}

func FDescribe(name string, fn func()) *suite {
	if currentSuite == nil {
		panic("FDescribe must be called with an active suite")
	}

	currentContext = newContext(name)
	currentContext.focus = true
	fn()

	currentSuite.contexts = append(currentSuite.contexts, currentContext)
	return currentSuite
}

func XDescribe(name string, fn func()) *suite {
	if currentSuite == nil {
		panic("XDescribe must be called with an active suite")
	}

	currentContext = newContext(name)
	currentContext.skip = true
	fn()

	currentSuite.contexts = append(currentSuite.contexts, currentContext)
	return currentSuite
}

func Before(fn func()) {
	if currentContext == nil {
		panic("Before must be called with an active context")
	}

	currentContext.before = append(currentContext.before, fn)
}

func JustBefore(fn func()) {
	if currentContext == nil {
		panic("JustBefore must be called with an active context")
	}

	currentContext.justBefore = append(currentContext.justBefore, fn)
}

func After(fn func()) {
	if currentContext == nil {
		panic("After must be called with an active context")
	}

	currentContext.after = append(currentContext.after, fn)
}

func Context(name string, fn func()) {
	if currentContext == nil {
		panic("Context must be called inside a Describe")
	}

	ctx := newContext(fmt.Sprintf("%s/%s", currentContext.name, name))
	ctx.parent = currentContext
	currentContext.children = append(currentContext.children, ctx)

	currentContext = ctx
	fn()
	currentContext = currentContext.parent
}

func XContext(name string, fn func()) {
	if currentContext == nil {
		panic("XContext must be called inside a Describe")
	}

	ctx := newContext(fmt.Sprintf("%s/%s", currentContext.name, name))
	ctx.skip = true
	ctx.parent = currentContext
	currentContext.children = append(currentContext.children, ctx)

	currentContext = ctx
	fn()
	currentContext = currentContext.parent
}

func FContext(name string, fn func()) {
	if currentContext == nil {
		panic("FContext must be called inside a Describe")
	}

	ctx := newContext(fmt.Sprintf("%s/%s", currentContext.name, name))
	ctx.focus = true
	ctx.parent = currentContext
	currentContext.children = append(currentContext.children, ctx)

	currentContext = ctx
	fn()
	currentContext = currentContext.parent
}

func It(name string, fn func(t *testing.T)) {
	if currentContext == nil {
		panic("It must be called inside a Describe or Context")
	}

	if currentSuite.focused {
		return
	}

	test := newTest(fmt.Sprintf("%s/%s", currentContext.name, name))
	test.focus = currentContext.focus
	test.fn = fn
	test.context = currentContext
	currentContext.tests = append(currentContext.tests, test)
	currentSuite.tests = append(currentSuite.tests, test)
}

func XIt(name string, fn func(t *testing.T)) {
	if currentContext == nil {
		panic("XIt must be called inside a Describe or Context")
	}

	if currentSuite.focused {
		return
	}

	test := newTest(fmt.Sprintf("%s/%s", currentContext.name, name))
	test.skip = true
	test.fn = fn
	test.context = currentContext
	currentContext.tests = append(currentContext.tests, test)
	currentSuite.tests = append(currentSuite.tests, test)
}

func FIt(name string, fn func(t *testing.T)) {
	if currentContext == nil {
		panic("FIt must be called inside a Describe or Context")
	}

	if !currentSuite.focused {
		currentSuite.focused = true
		currentSuite.tests = []*test{}
	}

	test := newTest(fmt.Sprintf("%s/%s", currentContext.name, name))
	test.focus = true
	test.fn = fn
	test.context = currentContext
	currentContext.tests = append(currentContext.tests, test)
	currentSuite.tests = append(currentSuite.tests, test)
}

func Run(t *testing.T) {
	run(t, currentSuite)
}
