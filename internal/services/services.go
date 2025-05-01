package services

import (
	"em-task/cfg"
	"em-task/internal/clients"
	"em-task/internal/models"
	"em-task/internal/repository"
	"em-task/pkg/logging"
	"errors"
)

type EnrichServiceDeps struct {
	*cfg.Config
	*repository.EnrichedPersonsRepository
	*clients.AgeAPIClient
	*clients.GenderAPIClient
	*clients.NationalityAPIClient
	Logger *logging.Logger
}

type EnrichService struct {
	*cfg.Config
	*repository.EnrichedPersonsRepository
	*clients.AgeAPIClient
	*clients.GenderAPIClient
	*clients.NationalityAPIClient
	Logger *logging.Logger
}

func NewEnrichService(deps EnrichServiceDeps) *EnrichService {
	return &EnrichService{
		Config:                    deps.Config,
		EnrichedPersonsRepository: deps.EnrichedPersonsRepository,
		AgeAPIClient:              deps.AgeAPIClient,
		GenderAPIClient:           deps.GenderAPIClient,
		NationalityAPIClient:      deps.NationalityAPIClient,
		Logger:                    deps.Logger,
	}
}

func (s *EnrichService) EnrichPersonData(requestData repository.CreatePersonDTO) (*models.Person, error) {
	ageChan := make(chan uint, 1)
	genderChan := make(chan string, 1)
	nationalityChan := make(chan string, 1)
	errChan := make(chan error, 3)

	go func() {
		ageData, err := s.AgeAPIClient.GetAgeByName(requestData.Name)
		if err != nil {
			s.Logger.ErrorLogger.Errorf("EnrichPersonData: AgeAPI error: %v", err)
			errChan <- err
			return
		}
		ageChan <- ageData.Age
	}()

	go func() {
		genderData, err := s.GenderAPIClient.GetGenderByName(requestData.Name)
		if err != nil {
			s.Logger.ErrorLogger.Errorf("EnrichPersonData: GenderAPI error: %v", err)
			errChan <- err
			return
		}
		genderChan <- genderData.Gender
	}()

	go func() {
		nationalityData, err := s.NationalityAPIClient.GetNationalityByName(requestData.Name)
		if err != nil {
			s.Logger.ErrorLogger.Errorf("EnrichPersonData: NationalityAPI error: %v", err)
			errChan <- err
			return
		}

		nationalityChan <- GetMostProbableNationality(nationalityData)
	}()

	var (
		age         uint
		gender      string
		nationality string
	)

	for range 3 {
		select {
		case age = <-ageChan:
		case gender = <-genderChan:
		case nationality = <-nationalityChan:
		case err := <-errChan:
			return nil, err
		}
	}

	person := &models.Person{
		Name:        requestData.Name,
		Surname:     requestData.Surname,
		Patronymic:  requestData.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	createdPerson, err := s.EnrichedPersonsRepository.Create(person)
	if err != nil {
		s.Logger.ErrorLogger.Errorf("EnrichPersonData: Create error: %v", err)
		return nil, err
	}
	return createdPerson, nil
}

func (s *EnrichService) DeletePersonById(id int) error {
	exists, err := s.EnrichedPersonsRepository.ExistsById(id)
	if err != nil {
		s.Logger.ErrorLogger.Errorf("DeletePersonById: ExistsById error: %v", err)
		return err
	}
	if !exists {
		s.Logger.ErrorLogger.Errorf("DeletePersonById: person not found, id=%d", id)
		return errors.New("person not found")
	}

	if err := s.EnrichedPersonsRepository.DeleteByID(id); err != nil {
		s.Logger.ErrorLogger.Errorf("DeletePersonById: DeleteByID error: %v", err)
		return err
	}

	return nil
}

func (s *EnrichService) GetPeopleWithParams(limit int, offset int, filters repository.PersonParams) []models.Person {
	people := s.EnrichedPersonsRepository.GetPeople(limit, offset, filters)
	return people
}

func (s *EnrichService) EditPerson(id int, editedData repository.FullPersonDTO) (*models.Person, error) {
	exists, err := s.EnrichedPersonsRepository.ExistsById(id)
	if err != nil {
		s.Logger.ErrorLogger.Errorf("EditPerson: ExistsById error: %v", err)
		return nil, err
	}
	if !exists {
		s.Logger.ErrorLogger.Errorf("EditPerson: person not found, id=%d", id)
		return nil, errors.New("person not found")
	}

	editedPerson, err := s.EnrichedPersonsRepository.EditByID(id, &models.Person{
		Name:        editedData.Name,
		Surname:     editedData.Surname,
		Patronymic:  editedData.Patronymic,
		Age:         editedData.Age,
		Gender:      editedData.Gender,
		Nationality: editedData.Nationality,
	})
	if err != nil {
		s.Logger.ErrorLogger.Errorf("EditPerson: EditByID error: %v", err)
		return nil, err
	}

	return editedPerson, nil
}

func GetMostProbableNationality(nationalityData *clients.PersonNationalityData) string {
	var resultCountry clients.CountryData
	var maxProbability float64
	for _, country := range nationalityData.Country {
		if country.Probability >= maxProbability {
			maxProbability = country.Probability
			resultCountry = country
		}
	}

	return resultCountry.CointryId
}
