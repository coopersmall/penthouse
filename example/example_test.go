package example_test

import (
	"fmt"
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
	ctx.Test("top level test", func(assert Assert) {
		// ...
		assert.Equal(num, 2)
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
				fmt.Print(num)
			}).
			Test("some test", func(assert Assert) {
				// ...
				fmt.Print(num)
				assert.Equal(num, 5)
				// ...
			}).
			Test("some other test", func(assert Assert) {
				// ...
				assert.Equal(num, 5)
			}).
			Context("sub sub test", func(ctx *Context) {
				ctx.
					Before(func() {
						num = 2
					}).
					JustBefore(func() {
						num = 6
					}).
					XTest("some test in the subtest", func(assert Assert) {
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
			Test("tis a focused subtest", func(assert Assert) {
				// ...
				assert.Equal(num, 6)
				// ...
			})
	})
})

func Test(t *testing.T) {
	Run(t, suite)
}
