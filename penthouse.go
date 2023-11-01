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
	fn()

	currentSuite.tests = append(currentSuite.tests, currentContext)
	return currentSuite
}

func FDescribe(name string, fn func()) *suite {
	if currentSuite == nil {
		panic("FDescribe must be called with an active suite")
	}

	currentContext = newContext(name)
	currentContext.focus = true
	fn()

	currentSuite.tests = append(currentSuite.tests, currentContext)
	return currentSuite
}

func XDescribe(name string, fn func()) *suite {
	if currentSuite == nil {
		panic("XDescribe must be called with an active suite")
	}

	currentContext = newContext(name)
	currentContext.skip = true
	fn()

	currentSuite.tests = append(currentSuite.tests, currentContext)
	return currentSuite
}

func Before(fn func()) {
	if currentContext == nil {
		panic("Before must be called with an active context")
	}

	currentContext.beforeEach(fn)
}

func JustBefore(fn func()) {
	if currentContext == nil {
		panic("JustBefore must be called with an active context")
	}

	currentContext.justBeforeEach(fn)
}

func After(fn func()) {
	if currentContext == nil {
		panic("After must be called with an active context")
	}

	currentContext.afterEach(fn)
}

func Context(name string, fn func()) {
	if currentContext == nil {
		panic("Context must be called inside a Describe")
	}

	ctx := newContext(fmt.Sprintf("%s/%s", currentContext.name, name))
	oldCtx := currentContext
	oldCtx.addChild(ctx)
	currentContext = ctx
	fn()
	currentContext = oldCtx
}

func XContext(name string, fn func()) {
	if currentContext == nil {
		panic("XContext must be called inside a Describe")
	}

	ctx := newContext(fmt.Sprintf("%s/%s", currentContext.name, name))
	ctx.skip = true
	oldCtx := currentContext
	oldCtx.addChild(ctx)
	currentContext = ctx
	fn()
	currentContext = oldCtx
}

func FContext(name string, fn func()) {
	if currentContext == nil {
		panic("FContext must be called inside a Describe")
	}

	ctx := newContext(fmt.Sprintf("%s/%s", currentContext.name, name))
	ctx.focus = true
	oldCtx := currentContext
	oldCtx.addChild(ctx)
	currentContext = ctx
	fn()
	currentContext = oldCtx
}

func It(name string, fn func(t *testing.T)) {
	if currentContext == nil {
		panic("It must be called inside a Describe or Context")
	}

	currentContext.it(name, fn)
}

func XIt(name string, fn func(t *testing.T)) {
	if currentContext == nil {
		panic("XIt must be called inside a Describe or Context")
	}

	currentContext.xit(name, fn)
}

func FIt(name string, fn func(t *testing.T)) {
	if currentContext == nil {
		panic("FIt must be called inside a Describe or Context")
	}

	currentContext.fit(name, fn)
}

func Run(t *testing.T) {
	run(t, currentSuite)
}
