package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
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
	DB() (*sql.DB, error)
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

func (d *dbImpl) DB() (*sql.DB, error) {
	return d.conn.DB()
}

var DBInstance DB

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DBInstance = &dbImpl{conn: gormDB}

	// AutoMigrate all models - add your models here
	// DBInstance.AutoMigrate(&User{}, &Activity{}, &Location{}, &Match{}, &Review{}, etc.)

	return nil
}

func GetDB() DB {
	if DBInstance == nil {
		panic("Database not initialized. Call InitDB() first.")
	}
	return DBInstance
}
