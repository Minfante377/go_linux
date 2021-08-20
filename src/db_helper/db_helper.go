package db_helper

import (
	"database/sql"
	"fmt"
	"logger_helper"

	_ "github.com/mattn/go-sqlite3"
)


func InitDb(db_name string, table_name string) int {
	logger_helper.LogInfo(fmt.Sprintf("Initializing db %s...", db_name))
	db, err := sql.Open("sqlite3", db_name)
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Failed to init db %s",
										   err.Error()))
		return -1
	}

	defer db.Close()
	logger_helper.LogInfo(fmt.Sprintf("Creating user table %s", table_name))
	var query string = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
								      uid INTEGER PRIMARY KEY AUTOINCREMENT,
									  username VARCHAR(64) NULL,
									  user_id VARCHAR(64) NULL
								  );`, table_name)
	stmt, err := db.Prepare(query)
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Failed to create user table %s",
										   err.Error()))
		return -1
	}

	_, err = stmt.Exec()
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Failed to create user table %s",
										   err.Error()))
		return -1
	}

	logger_helper.LogInfo("Database initialized succesfully")
	return 0
}


func AddUser(db_name string, table_name string,
			 username string, user_id int) int {
	logger_helper.LogInfo(fmt.Sprintf("Creating user %s...", username))
	db, err := sql.Open("sqlite3", db_name)
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Failed to connect to db %s",
										   err.Error()))
		return -1
	}

	defer db.Close()
	var query string = fmt.Sprintf(
		`INSERT INTO %s (username, user_id) values(?, ?)`,
		table_name)
	stmt, err := db.Prepare(query)
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Failed to add user %s",
										   err.Error()))
		return -1
	}

	_, err = stmt.Exec(username, user_id)
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Failed to add user %s",
										   err.Error()))
		return -1
	}


	logger_helper.LogInfo("User created")
	return 0
}


func DeleteUser(db_name string, table_name string, user_id int) int {
	logger_helper.LogInfo(fmt.Sprintf("Deleting user %d...", user_id))
	db, err := sql.Open("sqlite3", db_name)
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Failed to connect to db %s",
										   err.Error()))
		return -1
	}

	defer db.Close()
	var query string = fmt.Sprintf(
		`DELETE FROM %s WHERE user_id=?`,
		table_name)
	stmt, err := db.Prepare(query)
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Failed to delete user %s",
										   err.Error()))
		return -1
	}

	_, err = stmt.Exec(user_id)
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Failed to delete user %s",
										   err.Error()))
		return -1
	}

	logger_helper.LogInfo("User Deleted")
	return 0
}


func GetUsers(db_name string, table_name string) ([]string, []int) {
	db, err := sql.Open("sqlite3", db_name)
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Failed to connect to db %s",
										   err.Error()))
		return nil, nil
	}

	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", table_name))
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Failed to query users %s",
										   err.Error()))
		return nil, nil
	}

	var username, uid string
	var user_id int
	var usernames []string
	var user_ids []int
	for rows.Next() {
		err = rows.Scan(&uid, &username, &user_id)
		if err != nil {
			logger_helper.LogError(fmt.Sprintf("Failed to query user %s",
											   err.Error()))
			return nil, nil
		}
		usernames = append(usernames, username)
		user_ids = append(user_ids, user_id)
	}

	return usernames, user_ids
}
