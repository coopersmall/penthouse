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
	tests      []*test
	skip       bool
	focus      bool
	focused    map[string]*context
}

type test struct {
	name  string
	fn    func(t *testing.T)
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
		tests:   make([]*test, 0),
	}
}

func newRunner(name string) *test {
	return &test{
		name: name,
		fn:   func(t *testing.T) {},
	}
}

func newTestingSuite(t *testing.T) *testingSuite {
	return &testingSuite{
		t: t,
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

func (c *context) it(name string, fn func(t *testing.T)) *context {
	test := newRunner(fmt.Sprintf("%s/%s", c.name, name))
	test.fn = fn
	c.tests = append(c.tests, test)
	return c
}

func (c *context) xit(name string, fn func(t *testing.T)) *context {
	test := newRunner(fmt.Sprintf("%s/%s", c.name, name))
	test.skip = true
	c.tests = append(c.tests, test)
	return c
}

func (c *context) fit(name string, fn func(t *testing.T)) *context {
	test := newRunner(fmt.Sprintf("%s/%s", c.name, name))
	test.focus = true
	test.fn = fn
	c.tests = append(c.tests, test)
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

		message = focusTitle(fmt.Sprintf("%s : Focused %d tests", s.name, length))

	} else {
		for _, t := range s.tests {
			length += t.testLength()
		}

		message = title(fmt.Sprintf("%s : Running all %d tests", s.name, length))
		tests = s.tests
	}

	fmt.Println(message)

	if ts.setupSuite != nil {
		ts.setupSuite()
	}

	done := make(chan bool, len(tests))
	for _, t := range tests {
		go func(t *context) {
			runTest(t)(ts)
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

func runTest(c *context) opt {
	return func(s *testingSuite) {
		switch {
		case c.skip:
			c.skipTests(s)

		case len(c.focused) > 0:
			for _, f := range c.focused {
				runTest(f)(s)
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

			for _, test := range c.tests {
				s.t.Run(test.name, func(t *testing.T) {
					test.fn(t)
					if t.Failed() {
						fmt.Print(failure())
					} else {
						fmt.Print(success())
					}
				})
			}

			for _, after := range c.after {
				after()
			}

			if s.afterAll != nil {
				s.afterAll()
			}

		}

		for _, t := range c.children {
			runTest(t)(s)
		}

	}
}

func (c *context) skipTests(suite *testingSuite) {
	l := c.testLength()
	var st strings.Builder

	suite.t.Run(c.name, func(t *testing.T) {
		for i := 0; i <= l; i++ {
			st.WriteString(skip())
		}
	})
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
