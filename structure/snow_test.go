package structure

import (
	"fmt"
	"testing"
)

func TestSnow(t *testing.T) {
	snow := new(Snowflake)
	res := snow.NextVal()
	fmt.Println("xxx")
	fmt.Println(res)

}
