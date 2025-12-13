package repository

type HealthRepository interface {
	CheckDatabase() error
}
