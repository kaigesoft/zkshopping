package config

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"github.com/go-sql-driver/mysql"
	"github.com/cihub/seelog"
	"time"
)

var (
	dbMutex    sync.RWMutex
	db         = map[string]*sql.DB{}
	dbAddress  = map[string]string{}
	dbName     = map[string]string{}
	dbUsername = map[string]string{}
	dbPassword = map[string]string{}
)

func init() {
	mysql.RegisterDial("failover", failoverDial)
}

func initDatabase(config Config) {
	switch config.String("mode") {
	case "dev":
		dbAddress["db"] = "127.0.0.1:3306"
		dbUsername["db"] = "root"
		dbPassword["db"] = "123456"
		dbName["db"] = "zkshopping"
	case "release":
		dbAddress["db"] = "xxxx:xxx"
		dbUsername["db"] = "xxx"
		dbPassword["db"] = "xxx"
		dbName["db"] = "zkshopping"
	default:
		dbAddress["db"] = "127.0.0.1:3306"
		dbUsername["db"] = "root"
		dbPassword["db"] = "123456"
		dbName["db"] = "zkshopping"
	}
}

func failoverDial(addr string) (conn net.Conn, err error) {
	addrs := strings.Split(addr, ",")
	for _, v := range addrs {
		nd := net.Dialer{Timeout: 3 * time.Second}
		conn, err = nd.Dial("tcp", v)
		if err == nil {
			if strings.HasSuffix(v, ":6090") {
				return &secConn{conn}, err
			}
			return conn, err
		}
		seelog.Error("dial tcp ", v, " failed:", err)
	}
	return
}

type secConn struct {
	net.Conn
}

func (s *secConn) Read(b []byte) (n int, err error) {
	n, err = s.Conn.Read(b)
	if n > 0 {
		for i := 0; i < n; i++ {
			b[i] = b[i] ^ 0x66
		}
	}
	return
}

func (s *secConn) Write(b []byte) (n int, err error) {
	cb := make([]byte, len(b))
	if len(b) > 0 {
		for i := 0; i < len(b); i++ {
			cb[i] = b[i] ^ 0x66
		}
	}
	n, err = s.Conn.Write(cb)
	return
}

func GetDBConnect(dbname string) func() *sql.DB {
	return func() *sql.DB {
		dbMutex.RLock()
		conn := db[dbname]
		if conn != nil {
			dbMutex.RUnlock()
			return conn
		}
		dbMutex.RUnlock()

		dbMutex.Lock()
		defer dbMutex.Unlock()
		conn = db[dbname]
		if conn != nil {
			return conn
		}

		connStr := fmt.Sprintf("%s:%s@failover(%s)/%s?parseTime=true&loc=Local&charset=utf8&timeout=5s&readTimeout=30s&writeTimeout=30s", dbUsername[dbname], dbPassword[dbname], dbAddress[dbname], dbName[dbname])
		var err error
		conn, err = sql.Open("mysql", connStr)
		if err != nil {
			log.Fatal("Connect to database failed: ", connStr)
			log.Panicln(err)
		}
		conn.SetMaxIdleConns(10)
		db[dbname] = conn
		return conn
	}
}
