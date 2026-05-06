package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func respondDBError(c *gin.Context, message string, err error) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  "error",
		"message": message,
	})
}
func getConnectFromRequest(c *gin.Context) SConnect {
	return SConnect{
		Host:     c.Query("host"),
		Port:     c.Query("port"),
		User:     c.Query("user"),
		Password: c.Query("password"),
		DBName:   c.Query("dbname"),
	}
}

// TestConnect 测试 MySQL 连接
func TestConnect(c *gin.Context) {
	err := TestConnectService(getConnectFromRequest(c))
	if err != nil {
		respondDBError(c, "failed to connect to mysql", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// Version 获取 MySQL 版本
func Version(c *gin.Context) {
	version, err := VersionService(getConnectFromRequest(c))
	if err != nil {
		respondDBError(c, "failed to get mysql version", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": version,
	})
}

// 全部数据库
func Databases(c *gin.Context) {
	databases, err := DatabasesService(getConnectFromRequest(c))
	if err != nil {
		respondDBError(c, "failed to get databases", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"databases": databases,
	})
}

// 全部表
func Tables(c *gin.Context) {
	tables, err := TablesService(getConnectFromRequest(c))
	if err != nil {
		respondDBError(c, "failed to get tables", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"tables": tables,
	})
}

// 返回列信息
func Columns(c *gin.Context) {
	columnDetails, err := ColumnsService(getConnectFromRequest(c), c.Query("tablename"))
	if err != nil {
		respondDBError(c, "failed to get column details", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"columns": columnDetails,
	})
}

// 列表查询 分页 返回表结构
func List(c *gin.Context) {
	tablename := c.Query("tablename")
	if tablename == "" {
		respondDBError(c, "tablename is required", nil)
		return
	}
	page := c.Query("page")
	pageSize := c.Query("pageSize")
	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")
	if page == "" {
		page = "1"
	}
	if pageSize == "" {
		pageSize = "100"
	}
	if sortBy == "" {
		sortBy = "id"
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	data, err := ListService(getConnectFromRequest(c), tablename, pageInt, pageSizeInt, sortBy, sortOrder)
	if err != nil {
		respondDBError(c, "failed to list table data", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   data,
	})
}