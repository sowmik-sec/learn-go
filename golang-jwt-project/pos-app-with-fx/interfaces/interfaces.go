package interfaces

import "pos-app-with-fx/models"

type Database interface {
	SaveItem(item models.Item) error
	GetItems() ([]models.Item, error)
}

type Logger interface {
	Log(msg string)
}
