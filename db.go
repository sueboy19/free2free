package main

import (
	"context"

	"gorm.io/gorm"
)

type DB interface {
	AutoMigrate(dst ...interface{}) error
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	WithContext(ctx context.Context) *gorm.DB
	Preload(query string, args ...interface{}) *gorm.DB
	// 補齊常用方法
	Model(value interface{}) *gorm.DB
	Update(column string, value interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Joins(query string, args ...interface{}) *gorm.DB
	Raw(sql string, values ...interface{}) *gorm.DB
	Order(value interface{}) *gorm.DB
}

type dbImpl struct {
	conn *gorm.DB
}

func (d *dbImpl) AutoMigrate(dst ...interface{}) error {
	return d.conn.AutoMigrate(dst...)
}
func (d *dbImpl) Create(value interface{}) *gorm.DB {
	return d.conn.Create(value)
}
func (d *dbImpl) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return d.conn.First(dest, conds...)
}
func (d *dbImpl) Where(query interface{}, args ...interface{}) *gorm.DB {
	return d.conn.Where(query, args...)
}
func (d *dbImpl) Save(value interface{}) *gorm.DB {
	return d.conn.Save(value)
}
func (d *dbImpl) WithContext(ctx context.Context) *gorm.DB {
	return d.conn.WithContext(ctx)
}
func (d *dbImpl) Preload(query string, args ...interface{}) *gorm.DB {
	return d.conn.Preload(query, args...)
}
func (d *dbImpl) Model(value interface{}) *gorm.DB { return d.conn.Model(value) }
func (d *dbImpl) Update(column string, value interface{}) *gorm.DB {
	return d.conn.Model(nil).Update(column, value)
}
func (d *dbImpl) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return d.conn.Find(dest, conds...)
}
func (d *dbImpl) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return d.conn.Delete(value, conds...)
}
func (d *dbImpl) Joins(query string, args ...interface{}) *gorm.DB {
	return d.conn.Joins(query, args...)
}
func (d *dbImpl) Raw(sql string, values ...interface{}) *gorm.DB {
	return d.conn.Raw(sql, values...)
}
func (d *dbImpl) Order(value interface{}) *gorm.DB {
	return d.conn.Order(value)
}
