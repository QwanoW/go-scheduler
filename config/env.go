package config

import (
	"errors"
	"log"
	"os"
	"strconv"
)

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func GetenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}

//func GetenvInt(key string) (int, error) {
//	s, err := GetenvStr(key)
//	if err != nil {
//		return 0, err
//	}
//	v, err := strconv.Atoi(s)
//	if err != nil {
//		return 0, err
//	}
//	return v, nil
//}

func GetenvInt64(key string) int64 {
	s, err := GetenvStr(key)
	if err != nil {
		log.Fatal(err)
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return int64(v)
}

//func GetenvBool(key string) (bool, error) {
//	s, err := GetenvStr(key)
//	if err != nil {
//		return false, err
//	}
//	v, err := strconv.ParseBool(s)
//	if err != nil {
//		return false, err
//	}
//	return v, nil
//}
