package biz

import (
	"fmt"
	"os"
)

func MakeDirs(dir ...string) (err error) {
	for _, s := range dir {
		if err := os.MkdirAll(s, 0755); err != nil {
			return fmt.Errorf("mkdir %s error: %s", dir, err.Error())
		}
	}
	return nil
}
