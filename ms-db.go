package msAPI

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DBAccess interface {
	Query(sql string) (*sql.Rows, error)
	QuerySingle(sql string) (*sql.Row, error)
	Exec(sql string) (sql.Result, error)
	SaveChange()
	Close()
}

type dbaccess struct {
	conn *sql.DB
	tran *sql.Tx
}

func NewDBAccess(path string) (DBAccess, error) {
	db, err := sql.Open("mysql", path)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(10)
	db.SetMaxIdleConns(5)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &dbaccess{
		conn: db,
		tran: nil,
	}, nil
}

func GetDBDefaultPath() string {
	path := os.Getenv("dbaccessPath")
	if path == "" {
		path = "" //username:password@tcp(ip:port)/dbName?charset=utf8
	}
	return path
}

func (db *dbaccess) Query(sql string) (*sql.Rows, error) {
	stat, err := db.conn.Prepare(sql)
	if err != nil {
		return nil, err
	}

	rows, err := stat.Query()
	if err != nil {
		return nil, err
	}

	return rows, nil
}
func (db *dbaccess) QuerySingle(sql string) (*sql.Row, error) {
	stat, err := db.conn.Prepare(sql)
	if err != nil {
		return nil, err
	}

	row := stat.QueryRow()
	return row, nil
}
func (db *dbaccess) Exec(sql string) (sql.Result, error) {
	err := db.setTransition()
	if err != nil {
		return nil, err
	}

	//对SQL语句进行预处理
	stmt, err := db.conn.Prepare(sql)
	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec()
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (db *dbaccess) SaveChange() {
	if db.tran != nil {
		db.tran.Commit()
		db.tran = nil
	}
}
func (db *dbaccess) Close() {
	db.conn.Close()
}

func (db *dbaccess) setTransition() error {
	if db.tran == nil {
		tx, err := db.conn.Begin()
		if err != nil {
			return err
		}
		db.tran = tx
	}

	return nil
}
