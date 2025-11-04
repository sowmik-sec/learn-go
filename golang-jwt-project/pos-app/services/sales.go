package services

import (
	"errors"
	"fmt"
	"pos-app/interfaces"
	"pos-app/models"
)

type SalesService struct {
	inventory *InventoryService
	logger    interfaces.Logger
}

func NewSalesService(inventory *InventoryService, logger interfaces.Logger) *SalesService {
	return &SalesService{inventory: inventory, logger: logger}
}

func (ss *SalesService) Checkout(itemsWantToSell []models.ItemSell) (float64, error) {
	total := 0.0
	itemsInStore, err := ss.inventory.GetItems()

	if err != nil {
		return 0.0, err
	}

	for _, itemSell := range itemsWantToSell {
		found := false
		for i, itemStore := range itemsInStore {
			if itemSell.Name == itemStore.Name {
				if itemSell.Qty > itemStore.Stock {
					ss.logger.Log("Insufficient stock for " + itemSell.Name)
					return 0, errors.New("insufficient stock for " + itemSell.Name)
				}
				itemsInStore[i].Stock -= itemSell.Qty
				ss.logger.Log("Sold: " + itemSell.Name + " Qty: " + fmt.Sprintf("%d", itemSell.Qty))
				total += itemStore.Price * float64(itemSell.Qty)
				found = true
				break
			}
		}
		if !found {
			ss.logger.Log("Item not found: " + itemSell.Name)
			return 0, errors.New("item not found: " + itemSell.Name)
		}
	}
	return total, nil
}
