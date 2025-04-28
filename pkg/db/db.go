package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() error {
	var err error
	DB, err = sql.Open("sqlite3", "./factory.db")
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	err = CreateTables()
	if err != nil {
		return err
	}

	err = EnsureDefaultProductionLines()
	if err != nil {
		return err
	}

	return nil
}

func CreateTables() error {
	// Таблица продуктов
	_, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            gtin TEXT NOT NULL,
            label_data TEXT
        )
    `)
	if err != nil {
		return err
	}

	// Таблица заданий (обновленная)
	_, err = DB.Exec(`
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        product_id INTEGER NOT NULL,
        line_id INTEGER NOT NULL DEFAULT 1,
        date TEXT NOT NULL,
        batch_number TEXT NOT NULL,
        status TEXT NOT NULL DEFAULT 'новое',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (product_id) REFERENCES products(id),
        FOREIGN KEY (line_id) REFERENCES production_lines(id)
    )
`)
	if err != nil {
		return err
	}

	// Таблица производственных линий (новая)
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS production_lines (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL
        )
    `)
	if err != nil {
		return err
	}

	// Таблица истории изменений заданий (новая)
	_, err = DB.Exec(`
    CREATE TABLE IF NOT EXISTS task_history (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        task_id INTEGER NOT NULL,
        old_status TEXT NOT NULL,
        new_status TEXT NOT NULL,
        changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (task_id) REFERENCES tasks(id)
    )
`)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS mark_codes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL,
    code TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'new',
    used_at TIMESTAMP,
    file_position INTEGER,  -- Новое поле для хранения позиции в файле
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);
        
        -- Создаем индекс для ускорения поиска по коду
        CREATE INDEX IF NOT EXISTS idx_mark_codes_code ON mark_codes(code);
        
        -- Создаем индекс для поиска по task_id
        CREATE INDEX IF NOT EXISTS idx_mark_codes_task_id ON mark_codes(task_id);
        
        -- Создаем индекс для поиска по статусу
        CREATE INDEX IF NOT EXISTS idx_mark_codes_status ON mark_codes(status);
    `)
	if err != nil {
		return err
	}

	return nil
}
