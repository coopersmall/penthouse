// go:build unit

package example_test

import (
	"testing"

	. "github.com/coopersmall/penthouse/suite"
)

var suite = Suite("Example Testing Suite")

var _ = Describe("main test", func() {
	var id = 2

	Before(func() {
		id = 3
	})

	It("succeeds", func(assert Assert) {
		assert.Equal(id, 3)
		// ...
	})

	Context("sub test", func() {
		Before(func() {
			id = 6
		})

		It("succeeds", func(assert Assert) {
			assert.Equal(id, 6)
		})

		Context("sub sub test", func() {
			JustBefore(func() {
				id = 2
			})

			It("succeeds", func(assert Assert) {
				assert.Equal(id, 2)
			})
		})
	})

	Context("sub test 2", func() {
		It("succeeds", func(assert Assert) {
			assert.Equal(id, 3)
		})
	})

	Context("this is assert focused sub test", func() {
		Before(func() {
			id = 7
		})

		It("changes the id", func(assert Assert) {
			// ...
			assert.Equal(id, 7)
			// ...
		})
	})
})

func Test(t *testing.T) {
	Run(t)
}
