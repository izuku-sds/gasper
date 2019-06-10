package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sdslabs/SWS/lib/utils"
)

var dbHost = `%`
var dbUser = "root"

type mysqlAgentServer struct{}

// CreateDB creates a database in the Mysql instance with the given database name, user and password
func CreateDB(database, username, password string) error {
	port := utils.ServiceConfig["mysql"].(map[string]interface{})["container_port"].(string)

	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%s)", port)
	connection := fmt.Sprintf("%s@%s/", dbUser, agentAddress)

	db, err := sql.Open("mysql", connection)

	if err != nil {
		return fmt.Errorf("Error while creating the database : %s", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + database)
	if err != nil {
		return fmt.Errorf("Error while creating the database : %s", err)
	}

	query := fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s'", username, dbHost, password)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("Error while creating the database : %s", err)
	}

	query = fmt.Sprintf("GRANT ALL ON %s.* TO '%s'@'%s'", database, username, dbHost)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("Error while creating the database : %s", err)
	}

	_, err = db.Exec("FLUSH PRIVILEGES")
	if err != nil {
		return fmt.Errorf("Error while flushing user priviliges : %s", err)
	}

	return nil
}

// DeleteDB deletes the database given by the database name and username
func DeleteDB(database, username string) error {
	port := utils.ServiceConfig["mysql"].(map[string]interface{})["container_port"].(string)

	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%s)", port)
	connection := fmt.Sprintf("%s@%s/", dbUser, agentAddress)

	db, err := sql.Open("mysql", connection)
	if err != nil {
		return fmt.Errorf("Error while connecting to database : %s", err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE DATABASE " + database)
	if err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err)
	}

	_, err = db.Exec(fmt.Sprintf("DROP USER '%s'@'%s'", username, dbHost))
	if err != nil {
		return fmt.Errorf("Error while deleting the user : %s", err)
	}
	return nil
}
