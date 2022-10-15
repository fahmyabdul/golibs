package databases

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // For Postgres Conn
)

// PostgreCore :
type PostgreCore struct {
	Host       string      `json:"host"`
	Port       int         `json:"port"`
	User       string      `json:"user"`
	Pass       string      `json:"pass"`
	DB         string      `json:"db"`
	Logger     *log.Logger // Optional
	Connection *sql.DB
}

// NewPostgre :
func NewPostgre(host string, port int, user, pass, db string, nLog *log.Logger) (*PostgreCore, error) {
	if nLog == nil {
		nLog = log.New(os.Stderr, "", log.LstdFlags)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, db)

	postgreConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = postgreConn.Ping()
	if err != nil {
		return nil, err
	}

	postgreConf := PostgreCore{
		Host:       host,
		Port:       port,
		User:       user,
		Pass:       pass,
		DB:         db,
		Logger:     nLog,
		Connection: postgreConn,
	}

	return &postgreConf, nil
}

// Close :
func (p *PostgreCore) Close() error {
	err := p.Connection.Close()
	if err != nil {
		return err
	}

	return nil
}
