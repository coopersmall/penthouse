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

var _ = suite.Test("main test", func(ctx *Context) {
	ctx.Before(func() {
		num = 3
	})

	ctx.Test("name", func(assert Assert) {
		// ...
		assert.Equal(num, 3)
		// ...
	})

	ctx.Context("sub test 1", func(ctx *Context) {
		ctx.
			Before(func() {
				num = 3
			}).
			After(func() {
				num = 2
			}).
			JustBefore(func() {
				num = 5
			}).
			Test("name", func(assert Assert) {
				// ...
				assert.Equal(num, 5)
				// ...
			}).
			Test("some other name", func(assert Assert) {
				// ...
				assert.Equal(num, 5)
				// ...
			}).
			Context("sub sub test", func(ctx *Context) {
				ctx.
					Before(func() {
						num = 2
					}).
					JustBefore(func() {
						num = 6
					}).
					Test("", func(assert Assert) {
						// ...
						assert.Equal(num, 6)
						// ...
					})
			})

	})

	ctx.Context("this is assert focused sub test", func(ctx *Context) {
		ctx.
			Before(func() {
				num = 2
			}).
			JustBefore(func() {
				num = 6
			}).
			Test("", func(assert Assert) {
				// ...
				assert.Equal(num, 6)
				// ...
			})
	})
})

func Test(t *testing.T) {
	Run(t, suite)
}
