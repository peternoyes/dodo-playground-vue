package main

import (
	"fmt"
	"syscall/js"

	dodosim "github.com/peternoyes/dodo-sim"
)

var simulator *dodosim.SimulatorSync

func main() {
	fmt.Println("WASM: .main()")

	simulator = new(dodosim.SimulatorSync)

	js.Global().Set("_dodo_simulator_run", run())

	<-make(chan bool)
}

func run() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("WASM: .run()")
		return "This is a test"
	})
	return jsonFunc
}
