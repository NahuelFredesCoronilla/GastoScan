package expense

import "fmt"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (service *Service) List() ([]Expense, error) {
	return service.repository.FindAll()
}

func (service *Service) Get(id uint) (Expense, error) {
	expense, err := service.repository.FindByID(id)
	if err != nil {
		return Expense{}, fmt.Errorf("Error al obtener el gasto: %w", err)
	}

	return expense, nil
}

func (service *Service) Create(expense Expense) (Expense, error) {
	if expense.Name == "" {
		return Expense{}, fmt.Errorf("El nombre del gasto no puede estar vacío")
	}

	return service.repository.Create(expense)
}
