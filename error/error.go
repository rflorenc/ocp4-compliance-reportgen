package error

import "log"

// CheckError wraps log.Fatal(err) and panic(err)
func CheckError(msg string, err error) {
	if err != nil {
		log.Fatalf("%v :%v", msg, err)
		panic(err)
	}
}
