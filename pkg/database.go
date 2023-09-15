package pkg

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Container struct {
	DB *sqlx.DB
}

var DBCon *Container

func (c *Container) CreateMessage(msg Message) error {
	_, err := c.DB.Exec("INSERT INTO messages (topic, data) VALUES (?, ?)", msg.Topic, msg.Body)
	return err
}

func (c *Container) CreateMessages(msgs []Message) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO messages (topic, data) VALUES (?, ?)")
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		if _, err := stmt.Exec(msg.Topic, msg.Body); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func InitDB() {
	conn, err := sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		log.Panicln("failed to connect -> ", err)
	}

	if err := conn.Ping(); err != nil {
		log.Panicln("failed to ping -> ", err)
	}

	// Run migration
	if _, err := conn.Exec(getSchema()); err != nil {
		log.Panicln("failed to migrate -> ", err)
	}

	DBCon = &Container{DB: conn}
}

// Fake migration
func getSchema() string {
	return `
        CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            topic TEXT NOT NULL,
            data JSON NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );
    `
}
