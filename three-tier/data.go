package main

import(
	"time"
	"errors"
	"context"
	"strconv"
	"database/sql"
	
	"github.com/RileySun/Scaled/utils"
	
	//Caching
	gocache "github.com/patrickmn/go-cache"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/store/go_cache/v4"
)

//Struct
type DataService struct {
	DB *sql.DB
	
	//Caching
	ctx context.Context
	cacheManager *cache.Cache[[]byte]
}

//Create
func NewDataService() *DataService {
	//Load DB Creds
	creds := utils.LoadCredentials()
	
	//Create Object
	dataService := &DataService{
		DB:utils.NewDB(creds.Host, creds.Port, creds.User, creds.Pass, creds.Database),
		ctx:context.Background(),
	}
	
	//Load Caching using go-cache && gocache (what a confusing naming duo)
	cacheClient := gocache.New(5 * time.Minute, 10 * time.Minute)
	cacheStore := go_cache.NewGoCache(cacheClient)
	dataService.cacheManager = cache.New[[]byte](cacheStore)
	
	
	
	return dataService
}

//Public Methods
func (s *DataService) GetUserStatus(userID int) (string, error) {
	//Cache Check/Get
	value, cacheErr := s.getCache("User/Status-" + strconv.Itoa(userID))
	if cacheErr == nil {
		return value, nil
	}

	//Now Check DB
	var status string
	row := s.DB.QueryRow("SELECT `status` FROM Users WHERE `id` = ?;", userID)
	
	scanErr := row.Scan(&status)
	if scanErr != nil {
		return "", errors.New("Invalid User ID: " + strconv.Itoa(userID))
	}
	
	//Add to cache
	_ = s.setCache("User/Status-" + strconv.Itoa(userID), status)
	
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
	
	//Clear Cache for this user status
	_ = s.setCache("User/Status-" + strconv.Itoa(userID), newStatus)
	
	return nil
}

//Private Methods
func (s *DataService) getCache(key string) (string, error) {
	value, err := s.cacheManager.Get(s.ctx, key)
	if err != nil {
		return "", err
	}
	
	return string(value[:]), nil	
}

func (s *DataService) setCache(key, value string) error {
	err := s.cacheManager.Set(s.ctx, key, []byte(value))
	if err != nil {
		return err
	}
	return nil
}