//go:generate go run ../mock/generate example.go

package example

type DB interface {
	Get(id int) (string, error)
	Set(id int, value string) error
}

type Example struct {
	db DB
}

func NewExample(db DB) *Example {
	return &Example{
		db: db,
	}
}

func (e *Example) GetCustomer(id int) (string, error) {
	return e.db.Get(id)
}

func (e *Example) SetCustomer(id int, value string) error {
	return e.db.Set(id, value)
}
