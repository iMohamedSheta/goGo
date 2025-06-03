package helpers

import (
	"fmt"

	"github.com/iMohamedSheta/xapp/pkg/enums"
)

func LogConsoleError(err any) {
	fmt.Printf(enums.Red.Value()+"Error: %s\n"+enums.Reset.Value()+"\n", err)
}
