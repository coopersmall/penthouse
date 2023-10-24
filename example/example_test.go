// go:build unit

package example_test

import (
	"testing"

	. "github.com/coopersmall/penthouse"
	. "github.com/coopersmall/penthouse/example"
)

var (
	suite = Suite("Example Testing Suite")
)

var (
	db *DBMock
)

var _ = suite.BeforeAll(func() {
	db = NewDBMock()
})

var _ = suite.Describe("main test", func(ctx *Context) {
	var (
		result string
		err    error
		id     = 2

		itSucceeds = func(ctx *Context) {
			ctx.It("succeeds", func(assert Assert) {
				assert.Equal(result, "bob")
				assert.Equal(err, nil)
			})
		}

		itCallsDB = func(ctx *Context) {
			ctx.It("calls db", func(assert Assert) {
				assert.Equal(db.CallCount("Get"), 1)
				assert.Equal(db.CallParam("Get", 0, 0), id)
			})
		}
	)

	ctx.Before(func() {
		db.SetReturns("Get", "bob", nil)
	})

	ctx.JustBefore(func() {
		result, err = NewExample(db).GetCustomer(id)
	})

	ctx.After(func() {
		result = ""
		err = nil
		db.ClearReturns("Get")
	})

	itSucceeds(ctx)
	itCallsDB(ctx)

	ctx.Context("with different id", func(ctx *Context) {
		ctx.Before(func() {
			id = 3
		})

		itSucceeds(ctx)
		itCallsDB(ctx)

		ctx.Context("sub sub test", func(ctx *Context) {
			ctx.Before(func() {
				id = 2
			})

			itSucceeds(ctx)
			itCallsDB(ctx)
		})
	})

	ctx.Context("this is assert focused sub test", func(ctx *Context) {
		ctx.Before(func() {
			id = 7
		})

		itSucceeds(ctx)
		itCallsDB(ctx)

		ctx.It("changes the id", func(assert Assert) {
			// ...
			assert.Equal(id, 7)
			// ...
		})
	})
})

func Test(t *testing.T) {
	Run(t, suite)
}
