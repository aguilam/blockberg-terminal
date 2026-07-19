package main

import (
	"fmt"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

const initScript = `
	PRAGMA foreign_keys = ON;
	PRAGMA busy_timeout = 5000;

	CREATE TABLE IF NOT EXISTS barrel(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		x INTEGER NOT NULL,
		y INTEGER NOT NULL,
		z INTEGER NOT NULL,
		UNIQUE(x,y,z)
	);

	CREATE TABLE IF NOT EXISTS item_type(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		normalized_name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS seller(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		minecraft_uuid TEXT,
		username TEXT NOT NULL,
		normalized_username TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS item(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		normalized_name TEXT NOT NULL,
		type_id INTEGER REFERENCES item_type(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS barrel_item(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		item_id INTEGER REFERENCES item(id) ON DELETE CASCADE,
		barrel_id INTEGER REFERENCES barrel(id) ON DELETE CASCADE,
		seller_id INTEGER REFERENCES seller(id) ON DELETE SET NULL,
		quantity INTEGER NOT NULL,
		price INTEGER NOT NULL,
		benefit_ratio REAL NOT NULL,
		record_date TEXT DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS item_in_barrel(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		item_id INTEGER REFERENCES item(id) ON DELETE CASCADE,
		barrel_id INTEGER REFERENCES barrel(id) ON DELETE CASCADE,
		quantity INTEGER NOT NULL,
		record_date TEXT DEFAULT CURRENT_TIMESTAMP
	)
`

func InitDB(dbPath string) (*sqlite.Conn,error) {
	conn, err := sqlite.OpenConn(dbPath,sqlite.OpenReadWrite|sqlite.OpenCreate)
	if err != nil{
		return nil,fmt.Errorf("DB open error: %w",err)
	}
	err = sqlitex.ExecScript(conn,initScript)
	if err != nil {
		conn.Close() 
		return nil, fmt.Errorf("ошибка инициализации схемы: %w", err)
	}

	return conn, nil
}