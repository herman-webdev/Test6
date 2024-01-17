package service

import (
	"awesomeProject/internal/apperror"
	"awesomeProject/internal/user/dto"
	"awesomeProject/internal/user/model"
	"awesomeProject/internal/user/storage"
	"awesomeProject/pkg/api/sort"
	"awesomeProject/pkg/logging"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type userService struct {
	repository storage.Repository
	logger     *logging.Logger
}

func NewUserService(repository storage.Repository, logger *logging.Logger) UserService {
	return &userService{
		repository: repository,
		logger:     logger,
	}
}

func (s *userService) GetAll(ctx context.Context, sortOptions sort.Options) ([]model.User, error) {
	options := storage.NewSortOptions(sortOptions.Field, sortOptions.Order)
	all, err := s.repository.FindAll(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get all authors due to error: %v", err)
	}
	return all, nil
}

func (s *userService) GetOne(ctx context.Context, uuid string) (model.User, error) {
	one, err := s.repository.FindOne(ctx, uuid)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get the user by UUID: %v", err)
	}

	if one.ID == "" {
		return model.User{}, fmt.Errorf("failed to get the user due to error: %v", err)
	}

	return one, nil
}

func (s *userService) UpdateOne(ctx context.Context, dto dto.UpdateUserDto, uuid string) error {
	updatedUser := &model.User{
		Name:       dto.Name,
		Surname:    dto.Surname,
		Patronymic: dto.Patronymic,
		Age:        dto.Age,
		Gender:     dto.Gender,
		CountryId:  dto.CountryId,
	}

	_, err := s.repository.FindOne(ctx, uuid)
	if err != nil {
		return fmt.Errorf("failed to get the user: %w", err)
	}

	if err := s.repository.Update(ctx, updatedUser, uuid); err != nil {
		return fmt.Errorf("failed to update the user: %w", err)
	}

	return nil
}

func (s *userService) DeleteOne(ctx context.Context, uuid string) error {
	_, err := s.repository.FindOne(ctx, uuid)
	if err != nil {
		return fmt.Errorf("failed to get the user: %w", err)
	}

	if err := s.repository.Delete(ctx, uuid); err != nil {
		return fmt.Errorf("failed to delete the user: %w", err)
	}

	return nil
}

func (s *userService) CreateUser(ctx context.Context, dto dto.CreateUserDto) error {
	newUser := &model.User{
		Name:       dto.Name,
		Surname:    dto.Surname,
		Patronymic: dto.Patronymic,
	}

	age, err := fetchAgeFromAgifyAPI(dto.Name)
	if err != nil {
		return apperror.InternalServerError("Failed to fetch age from Apify API", err.Error())
	}
	newUser.Age = age

	gender, err := fetchGenderFromGenderizeAPI(dto.Name)
	if err != nil {
		return apperror.InternalServerError("Failed to fetch gender from Gender API", err.Error())
	}
	newUser.Gender = gender

	countryID, err := fetchCountryFromNationalizeAPI(dto.Name)
	if err != nil {
		return apperror.InternalServerError("Failed to fetch country from Nationalize API", err.Error())
	}
	newUser.CountryId = countryID

	if err := s.repository.Create(ctx, newUser); err != nil {
		return apperror.InternalServerError("Failed to create user", err.Error())
	}

	return nil
}

func fetchAgeFromAgifyAPI(name string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var agifyResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&agifyResponse); err != nil {
		return 0, err
	}

	age, ok := agifyResponse["age"].(float64)
	if !ok {
		return 0, fmt.Errorf("age not found in Agify API response")
	}

	return int(age), nil
}

func fetchGenderFromGenderizeAPI(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var genderizeResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&genderizeResponse); err != nil {
		return "", err
	}

	gender, ok := genderizeResponse["gender"].(string)
	if !ok {
		return "", fmt.Errorf("gender not found in Genderize API response")
	}

	return gender, nil
}

func fetchCountryFromNationalizeAPI(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var nationalizeResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&nationalizeResponse); err != nil {
		return "", err
	}

	countries, ok := nationalizeResponse["country"].([]interface{})
	if !ok || len(countries) == 0 {
		return "", fmt.Errorf("country data not found in Nationalize API response")
	}

	highestProbability := 0.0
	var countryID string

	for _, country := range countries {
		countryData, ok := country.(map[string]interface{})
		if !ok {
			continue
		}

		probability, ok := countryData["probability"].(float64)
		if !ok {
			continue
		}

		if probability > highestProbability {
			highestProbability = probability
			countryID, ok = countryData["country_id"].(string)
			if !ok {
				continue
			}
		}
	}

	if countryID == "" {
		return "", fmt.Errorf("country not found in Nationalize API response")
	}

	return countryID, nil
}
