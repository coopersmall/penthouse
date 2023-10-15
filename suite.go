package gotesting

import (
	"fmt"
	"testing"
)

var Suite = newSuite

type suite struct {
	name  string
	tests []*Context
	opts  []opt
}

type Context struct {
	name       string
	before     []func()
	justBefore []func()
	after      []func()
	children   []*Context
	runner     Runner
	skip       bool
	focus      bool
	focused    map[string]*Context
}

type Runner struct {
	name  string
	runs  []func(t *testing.T)
	skip  bool
	focus bool
}

func (t *Context) Before(fn func()) *Context {
	t.before = append(t.before, fn)
	return t
}

func (t *Context) JustBefore(fn func()) *Context {
	t.justBefore = append(t.justBefore, fn)
	return t
}

func (t *Context) Test(fn func(t *testing.T)) *Context {
	t.runner.runs = append(t.runner.runs, fn)
	return t
}

func (t *Context) XTest(fn func(t *testing.T)) *Context {
	t.runner.runs = append(t.runner.runs, fn)
	t.runner.skip = true
	return t
}

func (t *Context) After(fn func()) *Context {
	t.after = append(t.after, fn)
	return t
}

func (t *Context) Context(name string, fn func(*Context)) *Context {
	c := newContext(fmt.Sprintf("%s/%s", t.name, name))
	fn(c)
	addChild(t, c)
	return t
}

func (t *Context) XContext(name string, fn func(*Context)) *Context {
	c := newContext(fmt.Sprintf("%s/%s", t.name, name))
	c.skip = true
	fn(c)
	addChild(t, c)
	return t
}

func (t *Context) FContext(name string, fn func(*Context)) *Context {
	c := newContext(fmt.Sprintf("%s/%s", t.name, name))
	c.focus = true
	fn(c)
	addChild(t, c)
	return c
}

func newSuite(name string) *suite {
	return &suite{
		name:  name,
		tests: make([]*Context, 0),
		opts:  make([]opt, 0),
	}
}

func (s *suite) Assert(name string, t *testing.T, fn func(a Assert)) {
	a := newAsserter(name, t)
	fn(a)
}

func (s *suite) Test(name string, fn func(ctx *Context)) *suite {
	test := newContext(name)
	fn(test)
	s.tests = append(s.tests, test)

	return s
}

func (s *suite) Before(name string, f func()) *suite {
	s.opts = append(s.opts, before(name, f))
	return s
}

func (s *suite) BeforeAll(f func()) *suite {
	s.opts = append(s.opts, beforeAll(f))
	return s
}

func (s *suite) With(opts ...opt) *suite {
	s.opts = append(s.opts, opts...)
	return s
}

func Run(t *testing.T, s *suite) {
	ts := newTestingSuite(t)
	for _, opt := range s.opts {
		opt(ts)
	}

	fmt.Println("Running suite:", s.name)

	for i := range s.tests {
		focusTest(s.tests[i])
	}

	done := make(chan bool, len(s.tests))
	for _, t := range s.tests {
		go func(t *Context) {
			test(t)(ts)
			done <- true
		}(t)
	}

	for i := 0; i < len(s.tests); i++ {
		<-done
	}

}

type opt func(*testingSuite)

func test(c *Context) opt {
	return func(s *testingSuite) {
		if s.beforeAll != nil {
			s.beforeAll()
		}

		for _, before := range c.before {
			before()
		}

		for _, before := range c.justBefore {
			before()
		}

		testChildren := func() {
			for _, t := range c.children {
				test(t)(s)
			}
		}

		switch {
		case c.skip:
			l := len(getAllChildren(c))
			for i := 0; i <= l; i++ {
				fmt.Print(Yellow("•"))
			}
		case len(c.focused) > 0:
			for _, f := range c.focused {
				test(f)(s)
			}
		case c.runner.runs == nil:
			testChildren()
		default:
			for _, run := range c.runner.runs {
				run(s.T())
			}
			testChildren()
		}

		for _, after := range c.after {
			after()
		}

		if s.afterAll != nil {
			s.afterAll()
		}
	}
}

func newContext(name string) *Context {
	return &Context{
		name:    name,
		before:  make([]func(), 0),
		after:   make([]func(), 0),
		focused: make(map[string]*Context),
		runner: Runner{
			runs: make([]func(t *testing.T), 0),
		},
	}
}

func addChild(parent *Context, child *Context) {
	child.before = append(parent.before, child.before...)
	child.after = append(parent.after, child.after...)
	child.justBefore = append(parent.justBefore, child.justBefore...)
	parent.children = append(parent.children, child)
}

func getAllChildren(c *Context) []*Context {
	children := make([]*Context, 0)
	for _, child := range c.children {
		children = append(children, child)
		children = append(children, getAllChildren(child)...)
	}
	return children
}

func focusTest(c *Context) (*Context, bool) {
	for _, t := range c.children {
		has, ok := focusTest(t)
		if !ok {
			continue
		}

		c.focused[has.name] = has
	}

	if c.focus {
		return c, true
	}

	return nil, false
}

func SetupSuite(setupSuite func()) opt {
	return func(s *testingSuite) {
		s.setupSuite = setupSuite
	}
}

func TeardownSuite(teardownSuite func()) opt {
	return func(s *testingSuite) {
		s.teardownSuite = teardownSuite
	}
}

func before(test string, before func()) opt {
	return func(s *testingSuite) {
		s.before[test] = before
	}
}

func beforeAll(before func()) opt {
	return func(s *testingSuite) {
		s.beforeAll = before
	}
}

func after(test string, after func()) opt {
	return func(s *testingSuite) {
		s.after[test] = after
	}
}

func afterAll(after func()) opt {
	return func(s *testingSuite) {
		s.afterAll = after
	}
}

type testingSuite struct {
	setupSuite    func()
	teardownSuite func()
	beforeAll     func()
	afterAll      func()
	before        map[string]func()
	after         map[string]func()
	t             *testing.T
}

func newTestingSuite(t *testing.T) *testingSuite {
	return &testingSuite{
		before: make(map[string]func()),
		after:  make(map[string]func()),
		t:      t,
	}
}

func (s *testingSuite) T() *testing.T {
	return s.t
}
