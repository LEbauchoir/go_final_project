package database

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type DbHelper struct {
	Db *sql.DB
}

func InitDb() (*DbHelper, error) {
	appPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dbFile := filepath.Join(appPath, "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			install = true
		} else {
			return nil, fmt.Errorf("can't check db file: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}
	dbHelper := &DbHelper{Db: db}
	if install {
		if err := dbHelper.createTables(); err != nil {
			return nil, err
		}
	}
	// Проверка наличия таблицы tasks
	if err := dbHelper.checkTableExists("tasks"); err != nil {
		return nil, err
	}
	return dbHelper, nil
}
func (d *DbHelper) createTables() error {
	queryCreate := `
		CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			title TEXT NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat VARCHAR(128) NOT NULL DEFAULT ""
		);
		CREATE INDEX IF NOT EXISTS date_index ON scheduler (date);
	`
	_, err := d.Db.Exec(queryCreate)
	return err
}
func (d *DbHelper) checkTableExists(tableName string) error {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	var name string
	err := d.Db.QueryRow(query, tableName).Scan(&name)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Таблица %s не существует", tableName)
			return d.createTables()
		}
		return err
	}
	log.Printf("Таблица %s существует", tableName)
	return nil
}
