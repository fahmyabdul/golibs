package databases

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// PostgreCoreGORM :
type PostgreCoreGORM struct {
	Host       string      `json:"host"`
	Port       int         `json:"port"`
	User       string      `json:"user"`
	Pass       string      `json:"pass"`
	DB         string      `json:"db"`
	Logger     *log.Logger // Optional
	Connection *gorm.DB
}

// NewPostgreGORM :
func NewPostgreGORM(host string, port int, user, pass, db string, nLog *log.Logger) (*PostgreCoreGORM, error) {
	if nLog == nil {
		nLog = log.New(os.Stderr, "", log.LstdFlags)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, db)
	gormConnPostgre, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	gormConf := PostgreCoreGORM{
		Host:       host,
		Port:       port,
		User:       user,
		DB:         db,
		Connection: gormConnPostgre,
		Logger:     nLog,
	}

	return &gormConf, nil
}

// Close :
func (a *PostgreCoreGORM) Close() error {
	err := a.Connection.Close()

	return err
}

// Conditions :
type Conditions struct {
	Operator string      `json:"operator"`
	Value1   interface{} `json:"value_1"`
	Value2   interface{} `json:"value_2"`
}

// Filtering :
func Filtering(db *gorm.DB, conditions map[string]Conditions) *gorm.DB {
	db.Error = nil

	for key, value := range conditions {
		switch strings.ToLower(value.Operator) {
		case "lt", "lte", "gt", "gte", "eq", "ne":
			db = normalOperator(db, key, value)
		case "like":
			// key ILIKE val_1
			db = likeOperator(db, key, value)
		case "rng":
			// key BETWEEN val_1 AND val_2
			db = rangeOperator(db, key, value)
		case "rne":
			// key > val_1 AND key < val_2
			db = rangeNotEqualOperator(db, key, value)
		case "nir":
			// key NOT BETWEEN val_1 AND val_2
			db = notInRangeOperator(db, key, value)
		default:
			db.Error = fmt.Errorf("unrecognized operator")
		}
	}

	return db
}

// normalOperator :
func normalOperator(db *gorm.DB, key string, condition Conditions) *gorm.DB {
	operator := strings.ToLower(condition.Operator)
	switch operator {
	case "lt":
		operator = "<"
	case "lte":
		operator = "<="
	case "gt":
		operator = ">"
	case "gte":
		operator = ">="
	case "eq":
		operator = "="
	case "ne":
		operator = "!="
	}
	statement := fmt.Sprintf("%s %s ?", key, operator)
	db = db.Where(statement, condition.Value1)

	return db
}

// likeOperator :
func likeOperator(db *gorm.DB, key string, condition Conditions) *gorm.DB {
	val, ok := condition.Value1.(string)
	if !ok {
		db.Error = fmt.Errorf("value_1 of %s must be of string type", key)
		return db
	}
	statement := fmt.Sprintf("%s ILIKE ?", key)
	db = db.Where(statement, "%"+val+"%")

	return db
}

// rangeOperator :
func rangeOperator(db *gorm.DB, key string, condition Conditions) *gorm.DB {
	_, ok := condition.Value1.(float64)
	if !ok {
		db.Error = fmt.Errorf("value_1 of %s must be of integer or float type", key)
		return db
	}
	_, ok = condition.Value2.(float64)
	if !ok {
		db.Error = fmt.Errorf("value_2 of %s must be of integer or float type", key)
		return db
	}
	statement := fmt.Sprintf("%s BETWEEN ? AND ?", key)
	db = db.Where(statement, condition.Value1, condition.Value2)

	return db
}

// rangeNotEqualOperator :
func rangeNotEqualOperator(db *gorm.DB, key string, condition Conditions) *gorm.DB {
	_, ok := condition.Value1.(float64)
	if !ok {
		db.Error = fmt.Errorf("value_1 of %s must be of integer or float type", key)
		return db
	}
	_, ok = condition.Value2.(float64)
	if !ok {
		db.Error = fmt.Errorf("value_2 of %s must be of integer or float type", key)
		return db
	}
	statement := fmt.Sprintf("%s > ? AND %s < ?", key, key)
	db = db.Where(statement, condition.Value1, condition.Value2)

	return db
}

// notInRangeOperator :
func notInRangeOperator(db *gorm.DB, key string, condition Conditions) *gorm.DB {
	_, ok := condition.Value1.(float64)
	if !ok {
		db.Error = fmt.Errorf("value_1 of %s must be of integer or float type", key)
		return db
	}
	_, ok = condition.Value2.(float64)
	if !ok {
		db.Error = fmt.Errorf("value_2 of %s must be of integer or float type", key)
		return db
	}
	statement := fmt.Sprintf("%s NOT BETWEEN ? AND ?", key)
	db = db.Where(statement, condition.Value1, condition.Value2)

	return db
}
