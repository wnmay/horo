package db

import (
	"github.com/wnmay/horo/services/payment-service/internal/domain"
	"github.com/wnmay/horo/services/payment-service/internal/ports/outbound"
	"gorm.io/gorm"
)

type personModel struct {
	ID   string `gorm:"primaryKey;type:uuid"`
	Name string `gorm:"not null;index"`
}

func (personModel) TableName() string { return "people" }

type GormPersonRepository struct{ db *gorm.DB }

var _ outbound.PersonRepository = (*GormPersonRepository)(nil)

func NewGormPersonRepository(db *gorm.DB) *GormPersonRepository {
	//only in dev phase
	_ = db.AutoMigrate(&personModel{})
	return &GormPersonRepository{db: db}
}

func (r *GormPersonRepository) Save(p domain.Person) error {
	return r.db.Create(&personModel{ID: p.ID, Name: p.Name}).Error
}

func (r *GormPersonRepository) GetAll() ([]domain.Person, error) {
	var rows []personModel
	if err := r.db.Order("name ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Person, 0, len(rows))
	for _, m := range rows {
		out = append(out, domain.Person{ID: m.ID, Name: m.Name})
	}
	return out, nil
}