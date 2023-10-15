package example_test

import (
	"testing"

	. "github.com/coopersmall/penthouse"
)

var (
	suite = Suite("Testing Suite").With()
)

var (
	num int
)

var _ = suite.BeforeAll(func() {
	num = 2
})

var _ = suite.Test("main test", func(c *Context) {
	c.Test(func(t *testing.T) {
		// ...
		// ...
		suite.Assert("this is a test", t, func(a Assert) {
			a.Equal(num, 2)
		})
	})

	c.Context("sub test 1", func(c *Context) {
		c.
			Before(func() {
				num = 3
			}).
			After(func() {
				num = 2
			}).
			JustBefore(func() {
				num = 5
			}).
			Test(func(t *testing.T) {
				// ...
				suite.Assert("this is a test", t, func(a Assert) {
					a.Equal(num, 5)
				})

				suite.Assert("this is another test", t, func(a Assert) {
					a.Equal(num, 5)
				})
				// ...
			}).
			FContext("sub sub test", func(c *Context) {
				c.
					Before(func() {
						num = 2
					}).
					JustBefore(func() {
						num = 5
					}).
					Test(func(t *testing.T) {
						// ...
						suite.Assert("this is a sub test", t, func(a Assert) {
							a.Equal(num, 5)
						})
						// ...
					})
			})

	})

	c.Context("this is a focused sub test", func(c *Context) {
		c.
			Before(func() {
				num = 2
			}).
			JustBefore(func() {
				num = 6
			}).
			Test(func(t *testing.T) {
				// ...
				// ...
				suite.Assert("this is a sub test", t, func(a Assert) {
					a.Equal(num, 6)
				})
			})
	})
})

func Test(t *testing.T) {
	Run(t, suite)
}
