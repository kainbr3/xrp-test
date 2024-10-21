package repositories

import (
	m "go.mongodb.org/mongo-driver/mongo"
)

type IDocument interface {
	GetCollection(repo *Repository) *m.Collection
	GetName() string
}
