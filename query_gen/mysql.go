/* WIP Under Construction */
package main

//import "fmt"
import "strings"
import "errors"

func init() {
	db_registry = append(db_registry,&Mysql_Adapter{Name:"mysql"})
}

type Mysql_Adapter struct
{
	Name string
	Stmts string
	Body string
}

func (adapter *Mysql_Adapter) get_name() string {
	return adapter.Name
}

func (adapter *Mysql_Adapter) simple_insert(name string, table string, columns string, fields string) error {
	if name == "" {
		return errors.New("You need a name for this statement")
	}
	if table == "" {
		return errors.New("You need a name for this table")
	}
	if len(columns) == 0 {
		return errors.New("No columns found for simple_insert")
	}
	if len(fields) == 0 {
		return errors.New("No input data found for simple_insert")
	}
	
	var querystr string = "INSERT INTO `" + table + "`("
	
	// Escape the column names, just in case we've used a reserved keyword
	for _, column := range _process_columns(columns) {
		if column.Type == "function" {
			querystr += column.Left + ","
		} else {
			querystr += "`" + column.Left + "`,"
		}
	}
	
	// Remove the trailing comma
	querystr = querystr[0:len(querystr) - 1]
	
	querystr += ") VALUES ("
	for _, field := range _process_fields(fields) {
		querystr += field.Name + ","
	}
	querystr = querystr[0:len(querystr) - 1]
	
	adapter.write_statement(name,querystr + ")")
	return nil
}

func (adapter *Mysql_Adapter) simple_replace(name string, table string, columns string, fields string) error {
	if name == "" {
		return errors.New("You need a name for this statement")
	}
	if table == "" {
		return errors.New("You need a name for this table")
	}
	if len(columns) == 0 {
		return errors.New("No columns found for simple_insert")
	}
	if len(fields) == 0 {
		return errors.New("No input data found for simple_insert")
	}
	
	var querystr string = "REPLACE INTO `" + table + "`("
	
	// Escape the column names, just in case we've used a reserved keyword
	for _, column := range _process_columns(columns) {
		if column.Type == "function" {
			querystr += column.Left + ","
		} else {
			querystr += "`" + column.Left + "`,"
		}
	}
	
	// Remove the trailing comma
	querystr = querystr[0:len(querystr) - 1]
	
	querystr += ") VALUES ("
	for _, field := range _process_fields(fields) {
		querystr += field.Name + ","
	}
	querystr = querystr[0:len(querystr) - 1]
	
	adapter.write_statement(name,querystr + ")")
	return nil
}

func (adapter *Mysql_Adapter) simple_update(name string, table string, set string, where string) error {
	if name == "" {
		return errors.New("You need a name for this statement")
	}
	if table == "" {
		return errors.New("You need a name for this table")
	}
	if set == "" {
		return errors.New("You need to set data in this update statement")
	}
	
	var querystr string = "UPDATE `" + table + "` SET "
	for _, item := range _process_set(set) {
		querystr += "`" + item.Column + "` ="
		for _, token := range item.Expr {
			switch(token.Type) {
				case "function","operator","number","substitute":
					querystr += " " + token.Contents + ""
				case "column":
					querystr += " `" + token.Contents + "`"
				case "string":
					querystr += " '" + token.Contents + "'"
			}
		}
		querystr += ","
	}
	
	// Remove the trailing comma
	querystr = querystr[0:len(querystr) - 1]
	
	if len(where) != 0 {
		querystr += " WHERE"
		for _, loc := range _process_where(where) {
			var left, right string
			
			if loc.LeftType == "column" {
				left = "`" + loc.LeftColumn + "`"
			} else {
				left = loc.LeftColumn
			}
			
			if loc.RightType == "column" {
				right = "`" + loc.RightColumn + "`"
			} else {
				right = loc.RightColumn
			}
			
			querystr += " " + left + " " + loc.Operator + " " + right + " AND "
		}
		querystr = querystr[0:len(querystr) - 4]
	}
	
	adapter.write_statement(name,querystr)
	return nil
}

func (adapter *Mysql_Adapter) simple_select(name string, table string, columns string, where string, orderby string/*, offset int, maxCount int*/) error {
	if name == "" {
		return errors.New("You need a name for this statement")
	}
	if table == "" {
		return errors.New("You need a name for this table")
	}
	if len(columns) == 0 {
		return errors.New("No columns found for simple_select")
	}
	
	// Slice up the user friendly strings into something easier to process
	var colslice []string = strings.Split(strings.TrimSpace(columns),",")
	
	var querystr string = "SELECT "
	
	// Escape the column names, just in case we've used a reserved keyword
	for _, column := range colslice {
		querystr += "`" + strings.TrimSpace(column) + "`,"
	}
	
	// Remove the trailing comma
	querystr = querystr[0:len(querystr) - 1]
	
	querystr += " FROM `" + table + "`"
	if len(where) != 0 {
		querystr += " WHERE"
		for _, loc := range _process_where(where) {
			var left, right string
			
			if loc.LeftType == "column" {
				left = "`" + loc.LeftColumn + "`"
			} else {
				left = loc.LeftColumn
			}
			
			if loc.RightType == "column" {
				right = "`" + loc.RightColumn + "`"
			} else {
				right = loc.RightColumn
			}
			
			querystr += " " + left + " " + loc.Operator + " " + right + " AND "
		}
		querystr = querystr[0:len(querystr) - 4]
	}
	
	if len(orderby) != 0 {
		querystr += " ORDER BY "
		for _, column := range _process_orderby(orderby) {
			querystr += column.Column + " " + strings.ToUpper(column.Order) + ","
		}
		querystr = querystr[0:len(querystr) - 1]
	}
	
	adapter.write_statement(name,strings.TrimSpace(querystr))
	return nil
}

func (adapter *Mysql_Adapter) simple_left_join(name string, table1 string, table2 string, columns string, joiners string, where string, orderby string/*, offset int, maxCount int*/) error {
	if name == "" {
		return errors.New("You need a name for this statement")
	}
	if table1 == "" {
		return errors.New("You need a name for the left table")
	}
	if table2 == "" {
		return errors.New("You need a name for the right table")
	}
	if len(columns) == 0 {
		return errors.New("No columns found for simple_left_join")
	}
	if len(joiners) == 0 {
		return errors.New("No joiners found for simple_left_join")
	}
	
	var querystr string = "SELECT "
	
	for _, column := range _process_columns(columns) {
		var source, alias string
		
		// Escape the column names, just in case we've used a reserved keyword
		if column.Table != "" {
			source = "`" + column.Table + "`.`" + column.Left + "`"
		} else if column.Type == "function" {
			source = column.Left
		} else {
			source = "`" + column.Left + "`"
		}
		
		if column.Alias != "" {
			alias = " AS `" + column.Alias + "`"
		}
		querystr += source + alias + ","
	}
	
	// Remove the trailing comma
	querystr = querystr[0:len(querystr) - 1]
	
	querystr += " FROM `" + table1 + "` LEFT JOIN `" + table2 + "` ON "
	for _, joiner := range _process_joiner(joiners) {
		querystr += "`" + joiner.LeftTable + "`.`" + joiner.LeftColumn + "`=`" + joiner.RightTable + "`.`" + joiner.RightColumn + "` AND "
	}
	// Remove the trailing AND
	querystr = querystr[0:len(querystr) - 4]
	
	if len(where) != 0 {
		querystr += " WHERE"
		for _, loc := range _process_where(where) {
			var left, right string
			
			if loc.LeftTable != "" {
				left = "`" + loc.LeftTable + "`.`" + loc.LeftColumn + "`"
			} else if loc.LeftType == "column" {
				left = "`" + loc.LeftColumn + "`"
			} else {
				left = loc.LeftColumn
			}
			
			if loc.RightTable != "" {
				right = "`" + loc.RightTable + "`.`" + loc.RightColumn + "`"
			} else if loc.RightType == "column" {
				right = "`" + loc.RightColumn + "`"
			} else {
				right = loc.RightColumn
			}
			
			querystr += " " + left + " " + loc.Operator + " " + right + " AND "
		}
		querystr = querystr[0:len(querystr) - 4]
	}
	
	if len(orderby) != 0 {
		querystr += " ORDER BY "
		for _, column := range _process_orderby(orderby) {
			querystr += column.Column + " " + strings.ToUpper(column.Order) + ","
		}
		querystr = querystr[0:len(querystr) - 1]
	}
	
	adapter.write_statement(name,strings.TrimSpace(querystr))
	return nil
}

func (adapter *Mysql_Adapter) write() error {
	out := `// Code generated by. DO NOT EDIT.
/* This file was generated by Gosora's Query Generator. The thing above is to tell GH this file is generated. */
// +build !pgsql !sqlite !mssql
package main

import "log"
import "database/sql"

` + adapter.Stmts + `
func gen_mysql() (err error) {
	if debug {
		log.Print("Building the generated statements")
	}
` + adapter.Body + `
	return nil
}
`
	return write_file("./gen_mysql.go", out)
}

// Internal method, not exposed in the interface
func (adapter *Mysql_Adapter) write_statement(name string, querystr string ) {
	adapter.Stmts += "var " + name + "_stmt *sql.Stmt\n"
	
	adapter.Body += `	
	log.Print("Preparing ` + name + ` statement.")
	` + name + `_stmt, err = db.Prepare("` + querystr + `")
	if err != nil {
		return err
	}
	`
}
