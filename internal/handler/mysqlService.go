package handler

import (
	"database/sql"
	"fmt"
	"log"

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

// 链接MySQL 数据库
func connectMySQL(connect SConnect) (*sql.DB, error) {
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

// 测试连接 MySQL 数据库
func TestConnectService(connect SConnect) error {
	db, err := connectMySQL(connect)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping mysql: %w", err)
	}
	return nil
}

// 获取 MySQL 版本
func VersionService(connect SConnect) (string, error) {
	db, err := connectMySQL(connect)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var version string
	if err := db.QueryRow("SELECT VERSION()").Scan(&version); err != nil {
		return "", fmt.Errorf("failed to query mysql version: %w", err)
	}
	return version, nil
}

// 获取 MySQL 全部数据库
func DatabasesService(connect SConnect) ([]string, error) {
	db, err := connectMySQL(connect)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("failed to query databases: %w", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed to scan database name: %w", err)
		}
		databases = append(databases, dbName)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate databases: %w", err)
	}

	return databases, nil
}

// 获取 MySQL 全部表
func TablesService(connect SConnect) ([]string, error) {
	db, err := connectMySQL(connect)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tables: %w", err)
	}

	return tables, nil
}

// 获取 MySQL 列信息
func ColumnsService(connect SConnect, tablename string) ([]map[string]interface{}, error) {
	db, err := connectMySQL(connect)
	if err != nil {
		return nil, err
	}

	columns, err := db.Query("DESCRIBE " + tablename)
	if err != nil {
		return nil, fmt.Errorf("failed to describe table: %w", err)
	}
	defer columns.Close()

	var columnDetails []map[string]interface{}
	for columns.Next() {
		var field, colType string
		var null, key, defaultValue, extra sql.NullString
		if err := columns.Scan(&field, &colType, &null, &key, &defaultValue, &extra); err != nil {
			return nil, fmt.Errorf("failed to scan column details: %w", err)
		}

		defaultOutput := interface{}(nil)
		if defaultValue.Valid {
			defaultOutput = defaultValue.String
		}

		columnDetails = append(columnDetails, map[string]interface{}{
			"field":   field,
			"type":    colType,
			"null":    null.String,
			"key":     key.String,
			"default": defaultOutput,
			"extra":   extra.String,
		})
	}
	if err := columns.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate column details: %w", err)
	}

	return columnDetails, nil
}

// 获取 MySQL 列表
func ListService(connect SConnect, tablename string, page int, pageSize int, sortBy string, sortOrder string) ([]map[string]interface{}, error) {
	db, err := connectMySQL(connect)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	rows, err := db.Query("SELECT * FROM " + tablename + " ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanTargets := make([]interface{}, len(columns))
		for i := range values {
			scanTargets[i] = &values[i]
		}

		if err := rows.Scan(scanTargets...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result := make(map[string]interface{}, len(columns))
		for i, col := range columns {
			if raw, ok := values[i].([]byte); ok {
				result[col] = string(raw)
				continue
			}
			result[col] = values[i]
		}

		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return results, nil
}
