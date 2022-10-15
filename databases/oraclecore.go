package databases

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/godror/godror"
)

// OracleCore:
type OracleCore struct {
	Host    map[int]string `json:"host"`
	User    string         `json:"user"`
	Pass    string         `json:"pass"`
	DB      string         `json:"db"`
	Timeout string         `json:"timeout"`
	Logger  *log.Logger    // Optional
}

// NewOracle :
func NewOracle(host map[int]string, user, pass, db, timeout string, nLog *log.Logger) (*sql.DB, error) {
	if nLog == nil {
		nLog = log.New(os.Stderr, "", log.LstdFlags)
	}

	oracleCore := OracleCore{
		Host:    host,
		User:    user,
		Pass:    pass,
		DB:      db,
		Timeout: timeout,
		Logger:  nLog,
	}

	return oracleCore.Connect()
}

func (p *OracleCore) Connect() (*sql.DB, error) {
	for keyHost, currentHost := range p.Host {
		connectionString := fmt.Sprintf(`user="%s" password="%s" connectString="%s/%s?connect_timeout=%s"`, p.User, p.Pass, currentHost, p.DB, p.Timeout)

		db, err := sql.Open("godror", connectionString)
		if err != nil {
			p.Logger.Printf("| Oracle | Connecting To Host Number : %d | Error | %s\n", keyHost, err.Error())
			continue
		}

		p.Logger.Printf("| Oracle | Connecting To Host Number : %d | Success\n", keyHost)

		return db, nil
	}

	return nil, fmt.Errorf("Unable to connect to any configured host")
}
