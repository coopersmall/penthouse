package suite

import (
	"fmt"
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
	name    string
	opts    []opt
	focused bool
	skip    bool
	tests   []*test
}

type context struct {
	name       string
	before     []func()
	justBefore []func()
	after      []func()
	focused    bool
	skip       bool
	parent     *context
	children   []*context
}

type test struct {
	name    string
	fn      func(t *testing.T)
	skip    bool
	focused bool
	context *context
}

func newSuite(name string) *suite {
	return &suite{
		name: name,
		opts: make([]opt, 0),
	}
}

func newContext(name string) *context {
	return &context{
		name:   name,
		before: make([]func(), 0),
		after:  make([]func(), 0),
	}
}

func newTest(name string) *test {
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

func run(t *testing.T, s *suite) {
	fmt.Print(title(s))

	ts := newTestingSuite(t)
	for _, opt := range s.opts {
		opt(ts)
	}

	if ts.setupSuite != nil {
		ts.setupSuite()
	}

	for _, test := range s.tests {
		if s.skip || test.skip {
			skipTest(test)(ts)
			continue
		}
		runTest(test)(ts)
	}

	if ts.teardownSuite != nil {
		ts.teardownSuite()
	}

	fmt.Println()
}

func runTest(test *test) opt {
	return func(s *testingSuite) {
		if s.beforeAll != nil {
			s.beforeAll()
		}

		test.context.runBefore()
		test.context.runJustBefore()

		s.t.Run(test.name, func(t *testing.T) {
			test.fn(t)

			if t.Failed() {
				fmt.Print(failure())
			} else {
				fmt.Print(success())
			}
		})

		test.context.runAfter()

		if s.afterAll != nil {
			s.afterAll()
		}
	}
}

func skipTest(test *test) opt {
	return func(s *testingSuite) {
		s.t.Run(test.name, func(t *testing.T) {
			t.SkipNow()
		})
		fmt.Print(skip())
	}
}

func (c *context) runBefore() {
	if parent := c.parent; parent != nil {
		parent.runBefore()
	}

	for _, before := range c.before {
		before()
	}
}

func (c *context) runJustBefore() {
	if parent := c.parent; parent != nil {
		parent.runJustBefore()
	}

	for _, justBefore := range c.justBefore {
		justBefore()
	}
}

func (c *context) runAfter() {
	if parent := c.parent; parent != nil {
		parent.runAfter()
	}

	for _, after := range c.after {
		after()
	}
}
