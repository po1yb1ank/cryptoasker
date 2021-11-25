package db

import (
	"encoding/json"
	"time"

	// register driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const (
	configPort     = "db.port"
	configUser     = "db.user"
	configPassword = "db.password"
	configDBName   = "db.name"
)

type Client struct {
	port     string
	user     string
	password string
	dbname   string
	conn     *sqlx.DB
}

func NewClient() *Client {
	return &Client{
		port:     viper.GetString(configPort),
		user:     viper.GetString(configUser),
		password: viper.GetString(configPassword),
		dbname:   viper.GetString(configDBName),
	}
}

func (c *Client) Connect() error {
	conn, err := sqlx.Connect("mysql", c.user+":"+c.password+"@tcp(localhost"+c.port+")/"+c.dbname)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}
func (c *Client) CreateSchema() *Client {
	if c.conn != nil {
		c.conn.MustExec("CREATE TABLE `cryptopairs` (`id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY," +
			"`rawjson` json NOT NULL, `timestamp` timestamp NOT NULL)")
	}
	return c
}
func (c *Client) InsertRawJSON(rawjson json.RawMessage) error {
	_, err := c.conn.Exec("INSERT INTO cryptopairs (rawjson, timestamp) VALUES(?,?)",rawjson, timestamp())
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) GetLastRawJSON() (rawjson json.RawMessage, err error) {
	rows, err := c.conn.Query("SELECT rawjson FROM cryptopairs ORDER BY timestamp DESC LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		if err = rows.Scan(&rawjson); err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rawjson, err
}
func timestamp() string {
	return time.Now().UTC().Format("2006-01-02 03:04:05")
}
