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
// TestConnect 测试 MySQL 连接
func TestConnect(c *gin.Context){
	db, err := connectMySQL(c)
	log.Println(db)
	log.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to open mysql connection",
		})
		return
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to ping mysql",
		})
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
	dnsURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&charset=utf8mb4,utf8",
		connect.User,
		connect.Password,
		connect.Host,
		connect.Port,
	)
	if connect.DBName != "" {
		dnsURI += connect.DBName
	}
	log.Println("buildDNSURI", dnsURI)
	db, err := sql.Open("mysql", dnsURI)
	if err != nil {
		return nil, fmt.Errorf("open mysql connection: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}
	return db, nil
}
// Version 获取 MySQL 版本
func Version(c *gin.Context) {
	db, err := connectMySQL(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to open mysql connection",
		})
		return
	}
	defer db.Close()

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to query mysql version",
		})
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
	db, err := connectMySQL(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to open mysql connection",
		})
		return
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to query databases",
		})
		return
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to scan database name",
			})
			return
		}
		databases = append(databases, dbName)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"databases": databases,
	})
}