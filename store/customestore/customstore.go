package customestore

import (
	"DobroBot/model"
	"DobroBot/store"
	"fmt"
)

type customstore struct {
	s map[int]model.User
}

func NewStore() store.Store {
	return &customstore{s: make(map[int]model.User, 10)}
}

func (s *customstore) Add(u model.User) error {
	if _, ok := s.s[u.Id]; !ok {
		s.s[u.Id] = u
	} else {
		return fmt.Errorf("user %v already exists", u.Id)
	}

	return nil
}

func (s *customstore) Get(id int) (model.User, error) {
	if k, ok := s.s[id]; ok {
		return k, nil
	} else {
		return model.User{}, fmt.Errorf("user %v doesnt exists", id)
	}
}

func (s *customstore) UpdateDontes(id int, count int) error {
	if k, ok := s.s[id]; ok {
		k.Donations += count
		s.s[id] = k
		return nil
	} else {
		return fmt.Errorf("user %v doesnt exists", id)
	}
}

func (s *customstore) GetAllWithDonates(count int) ([]int, error) {
	result := make([]int, 0, 10)

	for k, v := range s.s {
		if v.Donations >= count {
			result = append(result, k)
		}
	}

	return result, nil
}
