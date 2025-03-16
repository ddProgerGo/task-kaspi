package utils

import (
	"errors"
	"strconv"
	"time"
)

type IINInfo struct {
	Correct     bool   `json:"correct"`
	Sex         string `json:"sex"`
	DateOfBirth string `json:"date_of_birth"`
}

func ValidateIIN(iin string) (*IINInfo, error) {
	if len(iin) != 12 {
		return nil, errors.New("Неверная длина ИИН")
	}

	year, err := strconv.Atoi(iin[:2])
	if err != nil {
		return nil, errors.New("Неверный формат ИИН")
	}

	centuryIndicator, err := strconv.Atoi(string(iin[6]))
	if err != nil || centuryIndicator < 1 || centuryIndicator > 6 {
		return nil, errors.New("Неверный индикатор века")
	}

	century := 1800 + (centuryIndicator+1)/2*100
	fullYear := century + year
	month, _ := strconv.Atoi(iin[2:4])
	day, _ := strconv.Atoi(iin[4:6])

	dob := time.Date(fullYear, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("02.01.2006")

	sex := "male"
	if centuryIndicator%2 == 0 {
		sex = "female"
	}

	return &IINInfo{Correct: true, Sex: sex, DateOfBirth: dob}, nil
}
