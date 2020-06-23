package monkebase

import (
	"os"
	"testing"
)

var (
	CONNECTION string = os.Getenv("MONKEBASE_CONNECTION")
)

func Test_connect(test *testing.T) {
	connect("root:monkebase@tcp(35.196.32.228:3306)/imonke")
}
