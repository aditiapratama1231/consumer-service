package category

type CategoryService interface {
	CreateCategory() error
	UpdateCategory() error
	DeleteCategory() error
}
