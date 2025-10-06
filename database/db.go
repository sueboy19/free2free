package database

import (
	"context"
	"gorm.io/gorm"
)

// DB interface to allow mocking and multiple implementations
type DB interface {
	WithContext(ctx interface{}) *ActualGormDB
	AutoMigrate(dst ...interface{}) error
	Where(query interface{}, args ...interface{}) *ActualGormDB
	First(dest interface{}, conds ...interface{}) error
	Create(value interface{}) error
	Delete(value interface{}, where ...interface{}) error
	Find(dest interface{}, conds ...interface{}) error
	Save(value interface{}) error
}

// ActualGormDB implements the DB interface using GORM
type ActualGormDB struct {
	Conn *gorm.DB
}

func (g *ActualGormDB) WithContext(ctx interface{}) *ActualGormDB {
	return &ActualGormDB{Conn: g.Conn.WithContext(ctx.(context.Context))}
}

func (g *ActualGormDB) AutoMigrate(dst ...interface{}) error {
	return g.Conn.AutoMigrate(dst...)
}

func (g *ActualGormDB) Where(query interface{}, args ...interface{}) *ActualGormDB {
	result := g.Conn.Where(query, args...)
	return &ActualGormDB{Conn: result}
}

func (g *ActualGormDB) First(dest interface{}, conds ...interface{}) error {
	return g.Conn.First(dest, conds...).Error
}

func (g *ActualGormDB) Create(value interface{}) error {
	return g.Conn.Create(value).Error
}

func (g *ActualGormDB) Delete(value interface{}, where ...interface{}) error {
	return g.Conn.Delete(value, where...).Error
}

func (g *ActualGormDB) Find(dest interface{}, conds ...interface{}) error {
	return g.Conn.Find(dest, conds...).Error
}

func (g *ActualGormDB) Save(value interface{}) error {
	return g.Conn.Save(value).Error
}

var GlobalDB *ActualGormDB

// For testing purposes - allows setting a mock database
func SetGlobalDB(db *ActualGormDB) {
	GlobalDB = db
}