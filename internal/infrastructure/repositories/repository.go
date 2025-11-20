package repositories

type Repository[E any, ID any] interface {
	Create(entity *E) error
	FindByID(id ID) (*E, error)
	FindAll() ([]*E, error)
	Update(entity *E) error
	DeleteByID(id ID) error
}
