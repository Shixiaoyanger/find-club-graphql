package dao

import (
	"find-club-graphql/config"
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
)

// mgo controller
type MgoDBCntlr struct {
	sess *mgo.Session
	db   *mgo.Database
}

var (
	DBNAME     = config.Conf.MongoDB.DBName
	globalSess *mgo.Session
	mongoURL   string
)

const (
	MongoCopyType  = "1"
	MongoCloneType = "2"
)

func init() {
	dbConf := config.Conf.MongoDB
	fmt.Println(dbConf)
	if dbConf.User != "" && dbConf.PW != "" {

		mongoURL = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", dbConf.User, dbConf.PW, dbConf.Host, dbConf.Port, dbConf.AdminDBName)
		fmt.Println(mongoURL, dbConf.Host)

	} else {
		mongoURL = fmt.Sprintf("mongodb://%s:%s", dbConf.Host, dbConf.Port)
		fmt.Println(mongoURL)

	}

	var err error
	globalSess, err = GetDBSession()
	if err != nil {
		fmt.Println("Failed to connect to MongoDB !")
		panic(err)
	}
	fmt.Println("Connect to MongoDB successfully!")
}

/****************************************** db session manage ****************************************/

// GetSession get the db session
func GetDBSession() (*mgo.Session, error) {
	globalMgoSession, err := mgo.DialWithTimeout(mongoURL, 10*time.Second)
	if err != nil {
		return nil, err
	}
	globalMgoSession.SetMode(mgo.Monotonic, true)
	//default is 4096
	globalMgoSession.SetPoolLimit(1000)
	return globalMgoSession, nil
}

func NewCloneMgoDBCntlr() *MgoDBCntlr {
	sess := globalSess.Clone()
	return &MgoDBCntlr{
		sess: sess,
		db:   sess.DB(DBNAME),
	}
}

func NewCopyMgoDBCntlr() *MgoDBCntlr {
	sess := globalSess.Copy()
	return &MgoDBCntlr{
		sess: sess,
		db:   sess.DB(DBNAME),
	}
}

func (m *MgoDBCntlr) Close() {
	m.sess.Close()
}

func (m *MgoDBCntlr) GetDB() *mgo.Database {
	return m.db
}

func (m *MgoDBCntlr) GetTable(tableName string) *mgo.Collection {
	return m.db.C(tableName)
}
