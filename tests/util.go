package tests

import "os"

var (
	gopath = os.Getenv("GOPATH")
)

// Help fulfillment the test.

func initDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	return nil
}
