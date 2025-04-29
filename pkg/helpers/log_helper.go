package helpers

import (
	"fmt"
	"imohamedsheta/gocrud/pkg/enums"
)

func LogConsoleError(err any) {
	fmt.Printf(enums.Red.Value()+"Error: %s\n"+enums.Reset.Value()+"\n", err)
}
