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
	ctx.Test(func(t *T) {
		// ...
		t.Assert("some test", func(assert Assert) {
			// ...
			assert.Equal(num, 2)
			// ...
		})
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
			Test(func(t *T) {
				// ...
				t.Assert("some test", func(assert Assert) {
					assert.Equal(num, 5)
				})
				// ...
			}).
			Test(func(t *T) {
				// ...
				t.Assert("some test", func(assert Assert) {
					assert.Equal(num, 5)
				})
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
					Test(func(t *T) {
						// ...
						t.Assert("some test", func(assert Assert) {
							assert.Equal(num, 6)
						})
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
			Test(func(t *T) {
				t.Assert("some test", func(assert Assert) {
					assert.Equal(num, 6)
				})
			})
	})
})

func Test(t *testing.T) {
	Run(t, suite)
}
