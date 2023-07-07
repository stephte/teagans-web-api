package dtos

import (
	"fmt"
)


func genQueryStr(path string, limit int, page int, sort string) string {
	return fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, limit, page, sort)
}
