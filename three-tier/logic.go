package main

import(
	"errors"
	"strconv"
	
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//Struct
type LogicService struct {
	dataService *DataService
}

//Create
func NewLogicService(ds *DataService) *LogicService {
	logic := &LogicService{
		dataService:ds,
	}
	
	return logic
}

//Methods
func (s *LogicService) GetStatus(userID string) (string, error) {
	//Validate userID
	userInt, userErr := strconv.Atoi(userID)
	if userErr != nil {
		return "", errors.New("Invalid User ID: " + userID)
	}
	
	//Get Status from DB
	status, err := s.dataService.GetUserStatus(userInt)
	if err != nil {
		return status, err
	}
	
	//Extra Logic (in this case upset statuses will be redacted for privacy)
	switch status {
		case "Upset":
			status = "Private"
		case "Sad":
			status = "Private"
		case "Angry":
			status = "Private"
	}
	
	//Return
	return status, nil
}

func (s *LogicService) SetStatus(userID, newStatus string) error {
	//Validate userID
	userInt, userErr := strconv.Atoi(userID)
	if userErr != nil {
		return errors.New("Invalid User ID: " + userID)
	}
	
	//Extra Logic (in this case capitalizing first letter of word and lower case everyting else)
	caser := cases.Title(language.English)
	fixedStatus := caser.String(newStatus)
	
	//Set Status in DB
	err := s.dataService.SetUserStatus(userInt, fixedStatus)
	if err != nil {
		return err
	}
	
	//Return
	return nil
}