package helpers

import (
	"fmt"
	"imohamedsheta/gocrud/enums"
)

func LogError(err any) {
	fmt.Printf(enums.Red.Value()+"Error: %s\n"+enums.Reset.Value()+"\n", err)
}
