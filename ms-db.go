package msAPI

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type IDBAccess interface {
	Query(sql string) (*sql.Rows, error)
	QuerySingle(sql string) (*sql.Row, error)
	Exec(sql string) (sql.Result, error)
	SaveChange()
	Close()
}

type dbaccess struct {
	connStr string
	db      *sql.DB
	tran    *sql.Tx
}

func DBAccess(path string) IDBAccess {
	return &dbaccess{
		connStr: path,
		db:      nil,
		tran:    nil,
	}
}

func DefaultConnStr() string {
	path := os.Getenv("dbaccessPath")
	if path == "" {
		path = "" //username:password@tcp(ip:port)/dbName?charset=utf8
	}
	return path
}
func ConnStr(uname, pwd, ip, port, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", uname, pwd, ip, port, dbName)
}

func (dba *dbaccess) Query(sql string) (*sql.Rows, error) {
	err := dba.connect()
	if err != nil {
		return nil, err
	}

	stat, err := dba.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stat.Close()

	rows, err := stat.Query()
	if err != nil {
		return nil, err
	}

	return rows, nil
}
func (dba *dbaccess) QuerySingle(sql string) (*sql.Row, error) {
	err := dba.connect()
	if err != nil {
		return nil, err
	}

	stat, err := dba.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stat.Close()

	row := stat.QueryRow()
	return row, nil
}
func (dba *dbaccess) Exec(sql string) (sql.Result, error) {
	err := dba.connect()
	if err != nil {
		return nil, err
	}

	err = dba.setTrans()
	if err != nil {
		return nil, err
	}

	//对SQL语句进行预处理
	stmt, err := dba.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (dba *dbaccess) SaveChange() {
	if dba.tran != nil {
		dba.tran.Commit()
		dba.tran = nil
	}
}
func (dba *dbaccess) Close() {
	if dba.db != nil {
		dba.db.Close()
		dba.db = nil
	}
}

func (dba *dbaccess) connect() error {
	if dba.db != nil {
		return nil
	}

	db, err := sql.Open("mysql", dba.connStr)
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(10)
	db.SetMaxIdleConns(5)
	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}

	dba.db = db
	return nil
}
func (dba *dbaccess) setTrans() error {
	if dba.tran == nil {
		tx, err := dba.db.Begin()
		if err != nil {
			return err
		}
		dba.tran = tx
	}

	return nil
}
