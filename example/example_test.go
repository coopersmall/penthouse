// go:build unit

package example_test

import (
	"testing"

	. "github.com/coopersmall/penthouse"
)

var suite = Suite("Example Testing Suite")

var _ = Describe("main test", func() {
	var id = 2

	Before(func() {
		id = 3
	})

	It("succeeds", func(t *testing.T) {
		if id != 3 {
			t.Errorf("id should be 3, but got %d", id)
		}
		// ...
	})

	Context("sub test", func() {
		Before(func() {
			id = 6
		})

		It("succeeds", func(t *testing.T) {
			if id != 6 {
				t.Errorf("id should be 6, but got %d", id)
			}
			// ...
		})

		Context("sub sub test", func() {
			JustBefore(func() {
				id = 2
			})

			It("succeeds", func(t *testing.T) {
				if id != 2 {
					t.Errorf("id should be 2, but got %d", id)
				}
				// ...
			})
		})
	})

	Context("sub test 2", func() {
		It("succeeds", func(t *testing.T) {
			if id != 3 {
				t.Errorf("id should be 3, but got %d", id)
			}
			// ...
		})
	})

	Context("this is assert focused sub test", func() {
		Before(func() {
			id = 7
		})

		It("changes the id", func(t *testing.T) {
			// ...
			if id != 7 {
				t.Errorf("id should be 7, but got %d", id)
			}
			// ...
		})
	})
})

func Test(t *testing.T) {
	Run(t)
}
