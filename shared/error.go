package shared

import (
	"log"
	"os"
)

func CheckError(err error) int {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return 0
}
