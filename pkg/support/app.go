package support

import (
	"fmt"
	"imohamedsheta/gocrud/enums"
)

// showHowToUse shows how to use the application
func PrintHowToUseApp() {
	fmt.Println("\n" + enums.Blue.Value() + "📌 Usage:" + enums.Reset.Value())
	fmt.Println("  " + enums.Green.Value() + "▶ To start the server:    go run . serve" + enums.Reset.Value())
	fmt.Println("  " + enums.Green.Value() + "▶ To run CLI commands:    go run . <command>" + enums.Reset.Value() + "\n")
}
