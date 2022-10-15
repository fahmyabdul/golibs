package databases

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCore :
type MongoCore struct {
	Host    map[int]string `json:"host"`
	User    string         `json:"user"`
	Pass    string         `json:"pass"`
	DB      string         `json:"db"`
	Srv     bool           `json:"srv"`
	Cluster bool           `json:"cluster"`
	RsName  string         `json:"rs_name,omitempty"`
	Logger  *log.Logger    // Optional
}

// NewMongo :
func NewMongo(host map[int]string, user, pass, db string, srv bool, cluster bool, rsname string, nLog *log.Logger) (*mongo.Database, error) {
	if nLog == nil {
		nLog = log.New(os.Stderr, "", log.LstdFlags)
	}

	mongoConf := MongoCore{
		Host:    host,
		User:    user,
		Pass:    pass,
		DB:      db,
		Srv:     srv,
		Cluster: cluster,
		RsName:  rsname,
		Logger:  nLog,
	}

	if cluster {
		mongoConn, err := mongoConf.ClusterConnect()
		if err != nil {
			return nil, err
		}
		return mongoConn, nil
	}

	mongoConn, err := mongoConf.StandAloneConnect()
	if err != nil {
		return nil, err
	}
	return mongoConn, nil
}

// Non Cluster Connection Method
func (p *MongoCore) StandAloneConnect() (*mongo.Database, error) {
	p.Logger.Println("| Mongo | Connecting To Mongo StandAlone")

	mongoDBURI := "mongodb://"
	if p.Srv == true {
		mongoDBURI = "mongodb+srv://"
	}

	// sortedHostNumbers := make([]int, len(p.Host))
	// i := 0
	// for k := range p.Host {
	// 	sortedHostNumbers[i] = k
	// 	i++
	// }
	// sort.Ints(sortedHostNumbers)

	for khost, vhost := range p.Host {
		if p.User == "" || p.Pass == "" {
			mongoDBURI = fmt.Sprintf("%s%s/%s?retryWrites=true&w=majority&connect=direct", mongoDBURI, vhost, p.DB)
		} else {
			mongoDBURI = fmt.Sprintf("%s%s:%s@%s/%s?retryWrites=true&w=majority&connect=direct", mongoDBURI, p.User, p.Pass, vhost, p.DB)
		}

		client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBURI))
		if err != nil {
			continue
			// return nil, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = client.Connect(ctx)
		if err != nil {
			continue
			// return nil, err
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			continue
			// return nil, err
		}
		err = HostStatus(client.Database(p.DB), context.Background())
		if err != nil {
			p.Logger.Printf("| Mongo | Connecting To Mongo StandAlone | Host Number : %d | Error | %s\n", khost, err.Error())
			continue
		}

		p.Logger.Printf("| Mongo | Connecting To Mongo StandAlone | Host Number : %d | Success\n", khost)
		return client.Database(p.DB), nil
	}
	return nil, fmt.Errorf("Unable to connect to any configured hostname")
}

// Cluster or Replica Connection Method
func (p *MongoCore) ClusterConnect() (*mongo.Database, error) {
	p.Logger.Println("| Mongo | Connecting To Mongo Cluster")

	mongoDBURI := "mongodb://"
	if p.Srv == true {
		mongoDBURI = "mongodb+srv://"
	}

	if len(p.Host) == 1 {
		return nil, errors.New("A cluster connection cannot be made if only a host is specified")
	}
	var ClusterHost []string
	for _, v := range p.Host {
		ClusterHost = append(ClusterHost, v)
	}
	hostname := strings.Join(ClusterHost, ",")

	if p.User == "" || p.Pass == "" {
		mongoDBURI = fmt.Sprintf("%s%s/%s?retryWrites=true&w=majority&replicaSet=%s", mongoDBURI, hostname, p.DB, p.RsName)
	} else {
		mongoDBURI = fmt.Sprintf("%s%s:%s@%s/%s?retryWrites=true&w=majority&replicaSet=%s", mongoDBURI, p.User, p.Pass, hostname, p.DB, p.RsName)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	p.Logger.Printf("| Mongo | Connecting To Mongo Cluster | Success\n")
	return client.Database(p.DB), nil

}

func HostStatus(db *mongo.Database, ctx context.Context) error {
	collection := db.Collection("test")
	_, err := collection.InsertOne(ctx, bson.D{{"name", "unitest"}})
	if err == nil {
		err = collection.Drop(ctx)
		if err != nil {

			return err
		}
	}
	return err
}
