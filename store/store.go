package store

import (
	"github.com/arvnd79/ie-proj/model"
	"github.com/jinzhu/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// GetUserByUserName retrieves user from database based on it's ID
// this method loads user's URLs and Requests lists
// returns error if user was not found
func (s *Store) GetUserByUserName(username string) (*model.User, error) {
	usr := &model.User{}
	err := s.db.Preload("Urls").Preload("Urls.Requests").First(usr, model.User{Username: username}).Error
	if err != nil {
		return nil, err
	}
	return usr, nil
}

// AddUser add's a user to the database
func (s *Store) AddUser(user *model.User) error {
	return s.db.Create(user).Error
}

// AddURL add's a url to the database
func (s *Store) AddURL(url *model.URL) error {
	return s.db.Create(url).Error
}

func (s *Store) GetAllURLs() ([]model.URL, error) {
	var urls []model.URL
	err := s.db.Model(&model.URL{}).Find(&urls).Error
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// GetURLById retrieves a URL from database based on it's ID
// returns error if an URL was not fount
func (s *Store) GetURLById(id uint) (*model.URL, error) {
	url := &model.URL{}
	err := s.db.Preload("Requests").First(url, id).Error
	if err != nil {
		return nil, err
	}
	return url, nil
}

// GetURLByUser retrieves urls for this user
// returns error if nothing was found
func (s *Store) GetURLsByUser(userID uint) ([]model.URL, error) {
	var urls []model.URL
	err := s.db.Model(&model.URL{}).Where("user_id == ?", userID).Find(&urls).Error
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// UpdateURL updates a URL to it's new value
func (s *Store) UpdateURL(url *model.URL) error {
	return s.db.Model(url).Update(url).Error
}

func (s *Store) FetchAlerts(userID uint) ([]model.URL, error) {
	var urls []model.URL
	err := s.db.Model(&model.URL{}).Where("user_id == ? and failed_times >= threshold", userID).Find(urls).Error
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// IncrementFailed increments failed_times of a URL
func (s *Store) IncrementFailed(url *model.URL) error {
	url.FailedTimes += 1
	return s.UpdateURL(url)
}

// AddRequest adds a request to database
func (s *Store) AddRequest(req *model.Request) error {
	return s.db.Create(req).Error
}
