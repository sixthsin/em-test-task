package repository

import (
	"em-task/internal/models"
	"em-task/pkg/db"
	"em-task/pkg/logging"
)

type EnrichedPersonsRepository struct {
	Db     *db.Db
	Logger *logging.Logger
}

func NewRepository(db *db.Db, logger *logging.Logger) *EnrichedPersonsRepository {
	return &EnrichedPersonsRepository{
		Db:     db,
		Logger: logger,
	}
}

func (r *EnrichedPersonsRepository) Create(person *models.Person) (*models.Person, error) {
	result := r.Db.DB.Create(person)
	if result.Error != nil {
		r.Logger.ErrorLogger.Errorf("Create: %v", result.Error)
		return nil, result.Error
	}
	return person, nil
}

func (r *EnrichedPersonsRepository) ExistsById(id int) (bool, error) {
	result := r.Db.DB.First(&models.Person{}, id)
	if result.Error != nil {
		r.Logger.ErrorLogger.Errorf("ExistsById: %v", result.Error)
		return false, result.Error
	}

	return true, nil
}

func (r *EnrichedPersonsRepository) DeleteByID(id int) error {
	result := r.Db.DB.Delete(&models.Person{}, id)
	if result.Error != nil {
		r.Logger.ErrorLogger.Errorf("DeleteByID: %v", result.Error)
		return result.Error
	}

	return nil
}

func (r *EnrichedPersonsRepository) GetPeople(limit, offset int, filters PersonParams) []models.Person {
	var people []models.Person
	query := r.Db.DB.Model(&models.Person{})

	if filters.Name != "" {
		query = query.Where("name LIKE ?", "%"+filters.Name+"%")
	}
	if filters.Surname != "" {
		query = query.Where("surname LIKE ?", "%"+filters.Surname+"%")
	}
	if filters.Patronymic != "" {
		query = query.Where("patronymic LIKE ?", "%"+filters.Patronymic+"%")
	}
	if filters.Age > 0 {
		query = query.Where("age = ?", filters.Age)
	}
	if filters.AgeFrom > 0 && filters.AgeTo > 0 {
		query = query.Where("age BETWEEN ? AND ?", filters.AgeFrom, filters.AgeTo)
	} else if filters.AgeFrom > 0 {
		query = query.Where("age >= ?", filters.AgeFrom)
	} else if filters.AgeTo > 0 {
		query = query.Where("age <= ?", filters.AgeTo)
	}
	if filters.Gender != "" {
		query = query.Where("gender = ?", filters.Gender)
	}
	if filters.Nationality != "" {
		query = query.Where("nationality = ?", filters.Nationality)
	}

	query.Limit(limit).Offset(offset).Find(&people)

	return people
}

func (r *EnrichedPersonsRepository) EditByID(id int, editedData *models.Person) (*models.Person, error) {
	var updatedPerson models.Person
	if err := r.Db.DB.Model(&updatedPerson).
		Where("id = ?", id).
		Updates(editedData).
		First(&updatedPerson).Error; err != nil {
		r.Logger.ErrorLogger.Errorf("EditByID: %v", err)
		return nil, err
	}

	return &updatedPerson, nil
}
