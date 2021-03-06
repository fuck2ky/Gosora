/*
*
* Gosora MySQL Interface
* Copyright Azareal 2017 - 2020
*
 */
package install

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Azareal/Gosora/query_gen"
	_ "github.com/go-sql-driver/mysql"
)

//var dbCollation string = "utf8mb4_general_ci"

func init() {
	adapters["mysql"] = &MysqlInstaller{dbHost: ""}
}

type MysqlInstaller struct {
	db         *sql.DB
	dbHost     string
	dbUsername string
	dbPassword string
	dbName     string
	dbPort     string
}

func (ins *MysqlInstaller) SetConfig(dbHost string, dbUsername string, dbPassword string, dbName string, dbPort string) {
	ins.dbHost = dbHost
	ins.dbUsername = dbUsername
	ins.dbPassword = dbPassword
	ins.dbName = dbName
	ins.dbPort = dbPort
}

func (ins *MysqlInstaller) Name() string {
	return "mysql"
}

func (ins *MysqlInstaller) DefaultPort() string {
	return "3306"
}

func (ins *MysqlInstaller) dbExists(dbName string) (bool, error) {
	var waste string
	err := ins.db.QueryRow("SHOW DATABASES LIKE '" + dbName + "'").Scan(&waste)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	} else if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func (ins *MysqlInstaller) InitDatabase() (err error) {
	_dbPassword := ins.dbPassword
	if _dbPassword != "" {
		_dbPassword = ":" + _dbPassword
	}
	db, err := sql.Open("mysql", ins.dbUsername+_dbPassword+"@tcp("+ins.dbHost+":"+ins.dbPort+")/")
	if err != nil {
		return err
	}

	// Make sure that the connection is alive..
	err = db.Ping()
	if err != nil {
		return err
	}
	fmt.Println("Successfully connected to the database")

	ins.db = db
	ok, err := ins.dbExists(ins.dbName)
	if err != nil {
		return err
	}

	if !ok {
		fmt.Println("Unable to find the database. Attempting to create it")
		_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + ins.dbName)
		if err != nil {
			return err
		}
		fmt.Println("The database was successfully created")
	}

	/*fmt.Println("Switching to database ", ins.dbName)
	_, err = db.Exec("USE " + ins.dbName)
	if err != nil {
		return err
	}*/
	db.Close()

	db, err = sql.Open("mysql", ins.dbUsername+_dbPassword+"@tcp("+ins.dbHost+":"+ins.dbPort+")/" + ins.dbName)
	if err != nil {
		return err
	}

	// Make sure that the connection is alive..
	err = db.Ping()
	if err != nil {
		return err
	}
	fmt.Println("Successfully connected to the database")

	// Ready the query builder
	ins.db = db
	qgen.Builder.SetConn(db)
	return qgen.Builder.SetAdapter("mysql")
}

func(ins *MysqlInstaller) createTable(f os.FileInfo) error {
	table := strings.TrimPrefix(f.Name(), "query_")
	ext := filepath.Ext(table)
	if ext != ".sql" {
		return nil
	}
	table = strings.TrimSuffix(table, ext)
	
	// ? - This is mainly here for tests, although it might allow the installer to overwrite a production database, so we might want to proceed with caution
	q := "DROP TABLE IF EXISTS `" + table + "`;"
	_, err := ins.db.Exec(q)
	if err != nil {
		fmt.Println("Failed query:", q)
		fmt.Println("e:",err)
		return err
	}

	data, err := ioutil.ReadFile("./schema/mysql/" + f.Name())
	if err != nil {
		return err
	}
	data = bytes.TrimSpace(data)

	_, err = ins.db.Exec(string(data))
	if err != nil {
		fmt.Println("Failed query:", string(data))
		fmt.Println("e:",err)
		return err
	}
	fmt.Printf("Created table '%s'\n", table)

	return nil
}

func (ins *MysqlInstaller) TableDefs() (err error) {
	fmt.Println("Creating the tables")
	files, err := ioutil.ReadDir("./schema/mysql/")
	if err != nil {
		return err
	}

	// TODO: Can we reduce the amount of boilerplate here?
	after := []string{"activity_stream_matches"}
	c1 := make(chan os.FileInfo)
	c2 := make(chan os.FileInfo)
	e := make(chan error)
	var wg sync.WaitGroup
	r := func(c chan os.FileInfo) {
		wg.Add(1)
		for f := range c {
			err := ins.createTable(f)
			if err != nil {
				e <- err
			}
		}
		wg.Done()
	}
	go r(c1)
	go r(c2)

	var a []os.FileInfo
Outer:
	for i, f := range files {
		if !strings.HasPrefix(f.Name(), "query_") {
			continue
		}
		table := strings.TrimPrefix(f.Name(), "query_")
		ext := filepath.Ext(table)
		if ext != ".sql" {
			continue
		}
		table = strings.TrimSuffix(table, ext)
		for _, tbl := range after {
			if tbl == table {
				a = append(a, f)
				continue Outer
			}
		}
		if i%2 == 0 {
			c1 <- f
		} else {
			c2 <- f
		}
	}
	close(c1)
	close(c2)
	wg.Wait()
	close(e)
	
	var first error
	for err := range e {
		if first == nil {
			first = err
		}
	}
	if first != nil {
		return first
	}

	for _, f := range a {
		if !strings.HasPrefix(f.Name(), "query_") {
			continue
		}
		err := ins.createTable(f)
		if err != nil {
			return err
		}
	}

	return nil
}

// ? - Moved this here since it was breaking the installer, we need to add this at some point
/* TODO: Implement the html-attribute setting type before deploying this */
/*INSERT INTO settings(`name`,`content`,`type`) VALUES ('meta_desc','','html-attribute');*/

func (ins *MysqlInstaller) InitialData() error {
	fmt.Println("Seeding the tables")
	data, err := ioutil.ReadFile("./schema/mysql/inserts.sql")
	if err != nil {
		return err
	}
	data = bytes.TrimSpace(data)

	statements := bytes.Split(data, []byte(";"))
	for key, sBytes := range statements {
		statement := string(sBytes)
		if statement == "" {
			continue
		}
		statement += ";"

		fmt.Println("Executing query #" + strconv.Itoa(key) + " " + statement)
		_, err = ins.db.Exec(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ins *MysqlInstaller) CreateAdmin() error {
	return createAdmin()
}

func (ins *MysqlInstaller) DBHost() string {
	return ins.dbHost
}

func (ins *MysqlInstaller) DBUsername() string {
	return ins.dbUsername
}

func (ins *MysqlInstaller) DBPassword() string {
	return ins.dbPassword
}

func (ins *MysqlInstaller) DBName() string {
	return ins.dbName
}

func (ins *MysqlInstaller) DBPort() string {
	return ins.dbPort
}
