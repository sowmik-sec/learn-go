package services

import (
	"fmt"
	"pos-app/models"
)

type InMemoryDB struct {
	items []models.Item
}

func (db *InMemoryDB) SaveItem(item models.Item) error {
	db.items = append(db.items, item)
	return nil
}

func (db *InMemoryDB) GetItems() ([]models.Item, error) {
	return db.items, nil
}

type ConsoleLogger struct{}

func (cl *ConsoleLogger) Log(msg string) {
	fmt.Println(msg)
}
