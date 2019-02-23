package tools

import (
	"github.com/satori/go.uuid"
	"strings"
)

func GetUUID() string {
	myUUID,err := uuid.NewV4()
	if err != nil {
		return "ERROR" 
	}
	return strings.Split(myUUID.String(), "-")[0]
}
