package main

import(
	"errors"
	"strconv"
	"database/sql"
	
	"github.com/RileySun/Scaled/utils"
)

//Struct
type DataService struct {
	DB *sql.DB
}

//Create
func NewDataService() *DataService {
	creds := utils.LoadCredentials()
	
	dataService := &DataService{
		DB:utils.NewDB(creds.Host, creds.Port, creds.User, creds.Pass, creds.Database),
	}
	
	return dataService
}

//Methods
func (s *DataService) GetUserStatus(userID int) (string, error) {
	var status string
	row := s.DB.QueryRow("SELECT `status` FROM Users WHERE `id` = ?;", userID)
	
	scanErr := row.Scan(&status)
	if scanErr != nil {
		return "", errors.New("Invalid User ID: " + strconv.Itoa(userID))
	}
	
	return status, nil
}

func (s *DataService) SetUserStatus(userID int, newStatus string) error {
	
	statement, prepErr := s.DB.Prepare("UPDATE Users SET `status`=? WHERE `id` = ?;")
	if prepErr != nil {
        return errors.New("Internal Service Error")
    }
    
	result, execErr := statement.Exec(newStatus, userID)
	rowsChanged, _ := result.RowsAffected() //Rows Affected always returns a nil error, why golang why (https://cs.opensource.google/go/go/+/master:src/database/sql/driver/driver.go;l=534?q=RowsAffected()&ss=go%2Fgo).
	if execErr != nil || rowsChanged == 0 {
		return errors.New("Invalid User ID: " + strconv.Itoa(userID))
	} //Double Check (sometimes there are no errors but no rows are affected)
	
	return nil
}