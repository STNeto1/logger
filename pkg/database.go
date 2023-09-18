package pkg

import (
	"log"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Container struct {
	DB *sqlx.DB
}

var DBCon *Container

func (c *Container) CreateMessage(msg Message) error {
	sql, args := sqlbuilder.NewInsertBuilder().InsertInto("messages").Cols("topic", "data").Values(msg.Topic, msg.Body).Build()

	_, err := c.DB.Exec(sql, args...)
	return err
}

func (c *Container) CreateMessages(msgs []Message) error {
	builder := sqlbuilder.NewInsertBuilder().InsertInto("messages").Cols("topic", "data")

	for _, msg := range msgs {
		builder.Values(msg.Topic, msg.Body)
	}

	sql, args := builder.Build()

	_, err := c.DB.Exec(sql, args...)
	return err
}

func (c *Container) GetMessages() ([]Message, error) {
	var msgs []Message

	sql, args := sqlbuilder.NewSelectBuilder().Select("*").From("messages").Build()

	if err := c.DB.Select(&msgs, sql, args...); err != nil {
		return nil, err
	}

	return msgs, nil
}

func (c *Container) GetTopics() ([]TopicMetadata, error) {
	var data []TopicMetadata

	sql, args := sqlbuilder.NewSelectBuilder().Select("topic", "count(topic)").From("messages").GroupBy("topic").Build()

	if err := c.DB.Select(&data, sql, args...); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Container) GetMessagesByTopic(topic string) ([]Message, error) {
	var msgs []Message

	sb := sqlbuilder.NewSelectBuilder().Select("*").From("messages")
	sql, args := sb.Where(sb.Equal("topic", topic)).Build()

	log.Println(sql, args)

	if err := c.DB.Select(&msgs, sql, args...); err != nil {
		return nil, err
	}

	return msgs, nil
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
	return sqlbuilder.NewCreateTableBuilder().
		CreateTable("messages").IfNotExists().
		Define("id", "INTEGER PRIMARY KEY AUTOINCREMENT").
		Define("topic", "TEXT NOT NULL").Define("data", "JSON NULL").
		Define("created_at", "DATETIME DEFAULT CURRENT_TIMESTAMP").
		String()
}
