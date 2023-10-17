package gotesting

import (
	"fmt"
	"strings"
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
	tests      []*T
	skip       bool
	focus      bool
	focused    map[string]*Context
}

type T struct {
	name  string
	runs  []func() error
	skip  bool
	focus bool
}

func (t *T) Assert(name string, fn func(a Assert)) {
	a := newAsserter(name, t)
	fn(a)
}

func (t *Context) Before(fn func()) *Context {
	t.before = append(t.before, fn)
	return t
}

func (t *Context) JustBefore(fn func()) *Context {
	t.justBefore = append(t.justBefore, fn)
	return t
}

func (t *Context) Test(fn func(t *T)) *Context {
	test := &T{
		name: t.name,
		runs: make([]func() error, 0),
	}
	fn(test)
	t.tests = append(t.tests, test)
	return t
}

func (t *Context) XTest(fn func(t *testing.T)) *Context {
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

	runTests := func(tests []*Context) {
		done := make(chan bool, len(tests))
		for _, t := range tests {
			go func(t *Context) {
				test(t)(ts)
				done <- true
			}(t)
		}

		for i := 0; i < len(tests); i++ {
			<-done
		}
	}

	focus := make(map[string]*Context)
	for i := range s.tests {
		has, ok := focusTest(s.tests[i])
		if !ok {
			continue
		}
		focus[has.name] = has
	}

	if len(focus) > 0 {
		for k := range focus {
			plural := "s"
			if len(focus[k].focused) == 1 {
				plural = ""
			}
			length := 0
			for _, t := range focus[k].focused {
				length += len(getAllRuns(t))
			}
			message := fmt.Sprintf(Orange("Focused %s")+": Running %d test"+plural, k, length)
			fmt.Println(strings.Repeat("-", len(message)))
			fmt.Println(message)
			fmt.Println(strings.Repeat("-", len(message)))

			tests := make([]*Context, 0)
			for _, t := range focus[k].focused {
				tests = append(tests, t)
			}
			runTests(tests)
		}
	} else {
		message := fmt.Sprintf(Cyan("%s")+": Running all tests", s.name)
		fmt.Println(strings.Repeat("-", len(message)))
		fmt.Println(message)
		fmt.Println(strings.Repeat("-", len(message)))

		runTests(s.tests)
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
		default:
			for _, t := range c.tests {
				for _, before := range c.before {
					before()
				}

				for _, before := range c.justBefore {
					before()
				}

				for _, run := range t.runs {
					s.t.Run(t.name, func(t *testing.T) {
						if err := run(); err != nil {
							fmt.Print(Red("•"))
							t.Error(err)
							return
						}
						fmt.Print(Green("•"))

					})
				}

				for _, after := range c.after {
					after()
				}
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
		tests:   make([]*T, 0),
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

func getAllRuns(c *Context) []func() error {
	runs := make([]func() error, 0)
	for _, child := range c.children {
		runs = append(runs, getAllRuns(child)...)
	}

	for _, t := range c.tests {
		for _, run := range t.runs {
			runs = append(runs, run)
		}
	}

	return runs
}

func focusTest(c *Context) (*Context, bool) {
	for _, t := range c.children {
		has, ok := focusTest(t)
		if ok {
			c.focused[has.name] = has
			continue
		}
	}

	if c.focus {
		return c, true
	}

	if len(c.focused) > 0 {
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
