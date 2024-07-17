package database

import (
	"database/sql"
	_ "embed"
	"log"
	"os"
	"path/filepath"
)

type DbHelper struct {
	Db *sql.DB
}

var createTableSQL = `
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    due_date DATE,
    completed BOOLEAN NOT NULL CHECK (completed IN (0, 1))
    date DATE
);
`

var createIndexSQL = `
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks (due_date);
`

var createSchedulerTableSQL = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    due_date DATE,
    completed BOOLEAN NOT NULL CHECK (completed IN (0, 1))
);
`

func InitDb() (*DbHelper, error) {
	appPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dbFile := filepath.Join(appPath, "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if os.IsNotExist(err) {
		install = true
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
	_, err := d.Db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	_, err = d.Db.Exec(createIndexSQL)
	if err != nil {
		return err
	}
	_, err = d.Db.Exec(createSchedulerTableSQL)
	if err != nil {
		return err
	}
	_, err = d.Db.Exec("INSERT INTO scheduler (id, title, description, due_date, completed) SELECT id, title, description, due_date, completed FROM tasks")
	return err
}
func (d *DbHelper) checkTableExists(tableName string) error {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	var name string
	err := d.Db.QueryRow(query, tableName).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Таблица %s не существует", tableName)
			return d.createTables()
		}
		return err
	}
	log.Printf("Таблица %s существует", tableName)
	return nil
}
