package services

import (
	"pos-app/interfaces"
	"pos-app/models"
)

type InventoryService struct {
	db     interfaces.Database
	logger interfaces.Logger
}

func NewInventoryService(db interfaces.Database, logger interfaces.Logger) *InventoryService {
	return &InventoryService{
		db:     db,
		logger: logger,
	}
}

func (is *InventoryService) AddItem(name string, price float64, stock int) error {
	item := models.Item{Name: name, Price: price, Stock: stock}
	err := is.db.SaveItem(item)
	if err == nil {
		is.logger.Log("Added item: " + name)
	}
	return err
}

func (is *InventoryService) GetItems() ([]models.Item, error) {
	return is.db.GetItems()
}
