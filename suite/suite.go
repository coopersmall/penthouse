package suite

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
	output        Output
	t             *testing.T
}

type opt func(*testingSuite)

type suite struct {
	name  string
	tests []*context
	opts  []opt
}

type context struct {
	name       string
	before     []func()
	justBefore []func()
	after      []func()
	children   []*context
	tests      []*runnter
	skip       bool
	focus      bool
	focused    map[string]*context
}

type runnter struct {
	name  string
	runs  []func() error
	test  func(Assert)
	skip  bool
	focus bool
}

func newSuite(name string) *suite {
	return &suite{
		name:  name,
		tests: make([]*context, 0),
		opts:  make([]opt, 0),
	}
}

func newContext(name string) *context {
	return &context{
		name:    name,
		before:  make([]func(), 0),
		after:   make([]func(), 0),
		focused: make(map[string]*context),
		tests:   make([]*runnter, 0),
	}
}

func newRunner(name string) *runnter {
	return &runnter{
		name: name,
		runs: make([]func() error, 0),
		test: func(assert Assert) {},
	}
}

func newTestingSuite(t *testing.T) *testingSuite {
	return &testingSuite{
		t:         t,
		formatter: NewFormatter(),
		output:    NewOutput(),
	}
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

func (t *context) beforeEach(fn func()) *context {
	t.before = append(t.before, fn)
	return t
}

func (t *context) justBeforeEach(fn func()) *context {
	t.justBefore = append(t.justBefore, fn)
	return t
}

func (c *context) it(name string, fn func(Assert)) *context {
	runner := newRunner(fmt.Sprintf("%s/%s", c.name, name))
	runner.test = fn
	c.tests = append(c.tests, runner)
	return c
}

func (c *context) xit(name string, fn func(Assert)) *context {
	runner := newRunner(fmt.Sprintf("%s/%s", c.name, name))
	runner.skip = true
	c.tests = append(c.tests, runner)
	return c
}

func (c *context) fit(name string, fn func(Assert)) *context {
	runner := newRunner(fmt.Sprintf("%s/%s", c.name, name))
	runner.focus = true
	runner.test = fn
	c.tests = append(c.tests, runner)
	return c
}

func (t *context) afterEach(fn func()) *context {
	t.after = append(t.after, fn)
	return t
}

func run(t *testing.T, s *suite) {
	ts := newTestingSuite(t)
	for _, opt := range s.opts {
		opt(ts)
	}

	focus := make(map[string]*context)
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
		tests   = make([]*context, 0)
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

	ts.output.Log(message, t)

	if ts.setupSuite != nil {
		ts.setupSuite()
	}

	done := make(chan bool, len(tests))
	for _, t := range tests {
		go func(t *context) {
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

func test(c *context) opt {
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

func (c *context) skipTests(suite *testingSuite) {
	suite.t.Run(c.name, func(t *testing.T) {
		l := c.testLength()

		var st strings.Builder
		for i := 0; i <= l; i++ {
			st.WriteString(suite.formatter.Skip())
		}

		suite.output.Skip(st.String(), t)
		return
	})
}

func (c *context) runTests(suite *testingSuite) {
	var (
		message string
		errs    []error
	)

	for _, runner := range c.tests {
		suite.t.Run(runner.name, func(t *testing.T) {
			runner.test(newAsserter(runner))

			for _, run := range runner.runs {
				if err := run(); err != nil {
					errs = append(errs, err)
				}
			}

			if len(errs) > 0 {
				message = suite.formatter.Failure(errs...)
				suite.output.Error(message, t)
			} else {
				message = suite.formatter.Success()
				suite.output.Log(message, t)
			}
		})

		message = ""
		errs = []error{}
	}
}

func (c *context) addChild(child *context) {
	child.before = append(c.before, child.before...)
	child.after = append(c.after, child.after...)
	child.justBefore = append(c.justBefore, child.justBefore...)
	c.children = append(c.children, child)
}

func (c *context) testLength() int {
	length := 0
	for _, child := range c.children {
		length += child.testLength()
	}

	for i := 0; i < len(c.tests); i++ {
		length++
	}

	return length
}

func (c *context) focusContext() (*context, bool) {
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
