package expense

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (repository *Repository) FindAll() ([]Expense, error) {
	var expenses []Expense
	err := repository.db.Find(&expenses).Error
	return expenses, err
}

func (repository *Repository) FindByID(id uint) (Expense, error) {
	var expense Expense
	err := repository.db.First(&expense, id).Error
	return expense, err
}

func (repository *Repository) Create(expense Expense) (Expense, error) {
	err := repository.db.Create(&expense).Error
	return expense, err
}
