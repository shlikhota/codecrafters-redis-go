package main

import (
	"errors"
	"strconv"
	"time"
)

type storage struct {
	storage map[string]storageValue
}

type storageValue struct {
	val        string
	expiration *time.Time
}

func NewStorage() *storage {
	return &storage{
		storage: make(map[string]storageValue),
	}
}

func (s *storage) Get(key string) (string, error) {
	val, ok := s.storage[key]
	if ok && (val.expiration == nil || val.expiration.After(time.Now())) {
		return val.val, nil
	}
	return "", errors.New("key doesn't exist")
}

func (s *storage) Set(key string, value string, args []string) bool {
	var exp *time.Time
	for i, arg := range args {
		if arg == "PX" {
			expireInMs, err := strconv.Atoi(args[i+1])
			if err == nil {
				e := time.Now().Add(time.Duration(expireInMs) * time.Millisecond)
				exp = &e
			}
		}
	}
	s.storage[key] = storageValue{value, exp}
	return true
}
