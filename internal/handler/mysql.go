package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type SConnect struct {
	Host     string
	Port     interface{}
	User     string
	Password string
	DBName   string
}

var connects = make(map[string]*sql.DB)

func respondDBError(c *gin.Context, message string, err error) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  "error",
		"message": message,
	})
}

func getDBFromRequest(c *gin.Context) (*sql.DB, bool) {
	db, err := connectMySQL(c)
	if err != nil {
		respondDBError(c, "failed to open mysql connection", err)
		return nil, false
	}
	return db, true
}

// TestConnect 测试 MySQL 连接
func TestConnect(c *gin.Context) {
	db, ok := getDBFromRequest(c)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		respondDBError(c, "failed to ping mysql", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// 链接MySQL 数据库
func connectMySQL(c *gin.Context) (*sql.DB, error) {
	connect := SConnect{
		Host:     c.DefaultQuery("host", ""),
		Port:     c.DefaultQuery("port", "3306"),
		User:     c.DefaultQuery("user", ""),
		Password: c.DefaultQuery("password", ""),
		DBName:   c.DefaultQuery("dbname", ""),
	}
	log.Println("connect", connect)
	dbPath := ""
	if connect.DBName != "" {
		dbPath = "/" + connect.DBName
	}
	dnsURI := fmt.Sprintf("%s:%s@tcp(%s:%s)%s?parseTime=true&charset=utf8mb4,utf8",
		connect.User,
		connect.Password,
		connect.Host,
		connect.Port,
		dbPath,
	)
	log.Println("buildDNSURI", dnsURI)
	if db, ok := connects[dnsURI]; ok {
		if err := db.Ping(); err == nil {
			return db, nil
		}
		_ = db.Close()
		delete(connects, dnsURI)
	}
	db, err := sql.Open("mysql", dnsURI)
	if err != nil {
		return nil, fmt.Errorf("open mysql connection: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}
	if connect.DBName != "" {
		connects[dnsURI] = db
	}
	return db, nil
}

// Version 获取 MySQL 版本
func Version(c *gin.Context) {
	db, ok := getDBFromRequest(c)
	if !ok {
		return
	}
	defer db.Close()

	var version string
	if err := db.QueryRow("SELECT VERSION()").Scan(&version); err != nil {
		respondDBError(c, "failed to query mysql version", err)
		return
	}

	log.Println("version", version)
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": version,
	})
}

// 全部数据库
func Databases(c *gin.Context) {
	db, ok := getDBFromRequest(c)
	if !ok {
		return
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		respondDBError(c, "failed to query databases", err)
		return
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			respondDBError(c, "failed to scan database name", err)
			return
		}
		databases = append(databases, dbName)
	}
	if err := rows.Err(); err != nil {
		respondDBError(c, "failed to iterate databases", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"databases": databases,
	})
}

// 全部表
func Tables(c *gin.Context) {
	if c.Query("dbname") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "missing dbname parameter",
		})
		return
	}
	db, ok := getDBFromRequest(c)
	if !ok {
		return
	}
	defer db.Close()

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		respondDBError(c, "failed to query tables", err)
		return
	}
	defer rows.Close()
	log.Println("rows", rows)
	tables := make([]string, 0)
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			respondDBError(c, "failed to scan table name", err)
			return
		}
		tables = append(tables, tableName)
	}
	if err := rows.Err(); err != nil {
		respondDBError(c, "failed to iterate tables", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"tables": tables,
	})
}