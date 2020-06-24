package models

import (
	"fmt"
)

type ContentType struct {
	ID    uint   `gorm:"primary_key"`
	Model string `gorm:"size:255;"`
}

func (contentType ContentType) String() string {
	return fmt.Sprintf(contentType.Model)
}
