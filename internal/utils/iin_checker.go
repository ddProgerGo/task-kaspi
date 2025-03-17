package utils

import (
	"errors"
	"strconv"
	"time"
)

const typeDate = "02.01.2006"

type IINInfo struct {
	Correct     bool   `json:"correct"`
	Sex         string `json:"sex"`
	DateOfBirth string `json:"date_of_birth"`
}

func ValidateIIN(iin string) (*IINInfo, error) {
	if len(iin) != 12 {
		return nil, errors.New("ИИН должен содержать 12 цифр")
	}

	for _, r := range iin {
		if r < '0' || r > '9' {
			return nil, errors.New("ИИН должен содержать только цифры")
		}
	}

	year, _ := strconv.Atoi(iin[0:2])
	month, _ := strconv.Atoi(iin[2:4])
	day, _ := strconv.Atoi(iin[4:6])

	centuryGender, _ := strconv.Atoi(string(iin[6]))
	var fullYear int
	var gender string

	switch centuryGender {
	case 1, 2:
		fullYear = 1800 + year
	case 3, 4:
		fullYear = 1900 + year
	case 5, 6:
		fullYear = 2000 + year
	default:
		return nil, errors.New("Некорректный 7-й символ в ИИН")
	}

	if centuryGender%2 == 1 {
		gender = "male"
	} else {
		gender = "female"
	}

	dateOfBirth := time.Date(fullYear, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if dateOfBirth.Month() != time.Month(month) || dateOfBirth.Day() != day {
		return nil, errors.New("Некорректная дата рождения в ИИН")
	}

	if !isValidChecksum(iin) {
		return nil, errors.New("Некорректная контрольная цифра ИИН")
	}

	return &IINInfo{
		true,
		gender,
		dateOfBirth.Format(typeDate),
	}, nil
}

func isValidChecksum(iin string) bool {
	weights1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	weights2 := []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 1, 2}

	sum := 0
	for i := 0; i < 11; i++ {
		digit, _ := strconv.Atoi(string(iin[i]))
		sum += digit * weights1[i]
	}

	controlDigit := sum % 11
	if controlDigit == 10 {
		sum = 0
		for i := 0; i < 11; i++ {
			digit, _ := strconv.Atoi(string(iin[i]))
			sum += digit * weights2[i]
		}
		controlDigit = sum % 11
	}

	expectedDigit, _ := strconv.Atoi(string(iin[11]))
	return controlDigit < 10 && controlDigit == expectedDigit
}
