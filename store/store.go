package store

import "DobroBot/model"

type Store interface {
	Add(u model.User) error
	Get(id int) (model.User, error)
	UpdateDontes(id int, count int) error
	GetAllWithDonates(count int) ([]int, error)
}