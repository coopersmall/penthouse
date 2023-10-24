package gotesting

import (
	"fmt"
	"strings"
	"testing"
)

type testingSuite struct {
	setupSuite    func()
	teardownSuite func()
	beforeAll     func()
	afterAll      func()
	formatter     Formatter
	output        func(*testing.T) Output
	t             *testing.T
}

type opt func(*testingSuite)

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
	tests      []*t
	skip       bool
	focus      bool
	focused    map[string]*Context
}

type t struct {
	name  string
	runs  []func() error
	test  func(Assert)
	skip  bool
	focus bool
}

var Suite = newSuite

func newSuite(name string) *suite {
	return &suite{
		name:  name,
		tests: make([]*Context, 0),
		opts:  make([]opt, 0),
	}
}

func newContext(name string) *Context {
	return &Context{
		name:    name,
		before:  make([]func(), 0),
		after:   make([]func(), 0),
		focused: make(map[string]*Context),
		tests:   make([]*t, 0),
	}
}

func newT(name string) *t {
	return &t{
		name: name,
		runs: make([]func() error, 0),
		test: func(assert Assert) {},
	}
}

func newTestingSuite(t *testing.T) *testingSuite {
	return &testingSuite{
		t:         t,
		formatter: NewFormatter(),
		output:    NewOutput,
	}
}

func (t *Context) Before(fn func()) *Context {
	t.before = append(t.before, fn)
	return t
}

func (t *Context) JustBefore(fn func()) *Context {
	t.justBefore = append(t.justBefore, fn)
	return t
}

func (c *Context) It(name string, fn func(Assert)) *Context {
	tt := newT(fmt.Sprintf("%s/%s", c.name, name))
	tt.test = fn
	c.tests = append(c.tests, tt)
	return c
}

func (c *Context) XIt(name string, fn func(Assert)) *Context {
	tt := newT(fmt.Sprintf("%s/%s", c.name, name))
	tt.skip = true
	c.tests = append(c.tests, tt)
	return c
}

func (t *Context) After(fn func()) *Context {
	t.after = append(t.after, fn)
	return t
}

func (t *Context) Context(name string, fn func(*Context)) *Context {
	c := newContext(fmt.Sprintf("%s/%s", t.name, name))
	t.addChild(c)
	fn(c)
	return t
}

func (t *Context) XContext(name string, fn func(*Context)) *Context {
	c := newContext(fmt.Sprintf("%s/%s", t.name, name))
	c.skip = true
	t.addChild(c)
	fn(c)
	return t
}

func (t *Context) FContext(name string, fn func(*Context)) *Context {
	c := newContext(fmt.Sprintf("%s/%s", t.name, name))
	c.focus = true
	t.addChild(c)
	fn(c)
	return c
}

func (s *suite) Describe(name string, fn func(ctx *Context)) *suite {
	test := newContext(name)
	fn(test)
	s.tests = append(s.tests, test)

	return s
}

func (s *suite) BeforeAll(f func()) *suite {
	s.opts = append(s.opts, beforeAll(f))
	return s
}

func (s *suite) AfterAll(f func()) *suite {
	s.opts = append(s.opts, afterAll(f))
	return s
}

func (s *suite) SetupSuite(f func()) *suite {
	s.opts = append(s.opts, setupSuite(f))
	return s
}

func (s *suite) TeardownSuite(f func()) *suite {
	s.opts = append(s.opts, teardownSuite(f))
	return s
}

func beforeAll(before func()) opt {
	return func(s *testingSuite) {
		s.beforeAll = before
	}
}

func afterAll(after func()) opt {
	return func(s *testingSuite) {
		s.afterAll = after
	}
}

func setupSuite(setupSuite func()) opt {
	return func(s *testingSuite) {
		s.setupSuite = setupSuite
	}
}

func teardownSuite(teardownSuite func()) opt {
	return func(s *testingSuite) {
		s.teardownSuite = teardownSuite
	}
}

func Run(t *testing.T, s *suite) {
	ts := newTestingSuite(t)
	for _, opt := range s.opts {
		opt(ts)
	}

	focus := make(map[string]*Context)
	for i := range s.tests {
		has, ok := s.tests[i].focusContext()
		if !ok {
			continue
		}
		focus[has.name] = has
	}

	var (
		message string
		length  = 0
		tests   = make([]*Context, 0)
	)

	if len(focus) > 0 {
		for k := range focus {
			for _, t := range focus[k].focused {
				length += t.testLength()
				tests = append(tests, t)
			}
		}

		message = ts.formatter.Focus(fmt.Sprintf("%s : Focused %d tests", s.name, length))

	} else {
		for _, t := range s.tests {
			length += t.testLength()
		}

		message = ts.formatter.Title(fmt.Sprintf("%s : Running all %d tests", s.name, length))
		tests = s.tests
	}

	ts.output(t).Log(message)

	if ts.setupSuite != nil {
		ts.setupSuite()
	}

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

	if ts.teardownSuite != nil {
		ts.teardownSuite()
	}

	fmt.Println()
}

func test(c *Context) opt {
	return func(s *testingSuite) {
		switch {
		case c.skip:
			c.skipTests(s)

		case len(c.focused) > 0:
			for _, f := range c.focused {
				test(f)(s)
			}

		default:
			if s.beforeAll != nil {
				s.beforeAll()
			}

			for _, before := range c.before {
				before()
			}

			for _, before := range c.justBefore {
				before()
			}

			c.runTests(s)

			for _, after := range c.after {
				after()
			}

			if s.afterAll != nil {
				s.afterAll()
			}

		}

		for _, t := range c.children {
			test(t)(s)
		}

	}
}

func (c *Context) skipTests(suite *testingSuite) {
	suite.t.Run(c.name, func(t *testing.T) {
		l := c.testLength()

		var st strings.Builder
		for i := 0; i <= l; i++ {
			st.WriteString(suite.formatter.Skip())
		}

		suite.output(t).Skip(st.String())
		return
	})
}

func (c *Context) runTests(suite *testingSuite) {
	var (
		message string
		errs    []error
	)

	for _, tt := range c.tests {
		suite.t.Run(tt.name, func(t *testing.T) {
			tt.test(newAsserter(tt))

			for _, run := range tt.runs {
				if err := run(); err != nil {
					errs = append(errs, err)
				}
			}

			if len(errs) > 0 {
				message = suite.formatter.Failure(errs...)
			} else {
				message = suite.formatter.Success()
			}

			suite.output(t).Log(message)
		})

		message = ""
		errs = []error{}
	}
}

func (c *Context) addChild(child *Context) {
	child.before = append(c.before, child.before...)
	child.after = append(c.after, child.after...)
	child.justBefore = append(c.justBefore, child.justBefore...)
	c.children = append(c.children, child)
}

func (c *Context) testLength() int {
	length := 0
	for _, child := range c.children {
		length += child.testLength()
	}

	for i := 0; i < len(c.tests); i++ {
		length++
	}

	return length
}

func (c *Context) focusContext() (*Context, bool) {
	for _, t := range c.children {
		has, ok := t.focusContext()
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
