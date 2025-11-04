package services

import (
	"fmt"
	"pos-app-with-fx/interfaces"
	"pos-app-with-fx/models"
)

type InMemoryDB struct {
	items []models.Item
}

func NewInMemoryDB() interfaces.Database {
	return &InMemoryDB{}
}

func (db *InMemoryDB) SaveItem(item models.Item) error {
	db.items = append(db.items, item)
	return nil
}

func (db *InMemoryDB) GetItems() ([]models.Item, error) {
	return db.items, nil
}

type ConsoleLogger struct{}

func NewConsoleLogger() interfaces.Logger {
	return &ConsoleLogger{}
}

func (cl *ConsoleLogger) Log(msg string) {
	fmt.Println(msg)
}
