package uuid

import uuid "github.com/satori/go.uuid"

func Get() string {
	return uuid.NewV4().String()
}