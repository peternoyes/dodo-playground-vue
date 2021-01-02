package main

import (
	"fmt"
	"syscall/js"
	"time"

	dodosim "github.com/peternoyes/dodo-sim"
)

var simulator *dodosim.SimulatorSync
var didStop bool
var keys string = ""
var c chan bool

func main() {
	fmt.Println("WASM: .main()")

	simulator = new(dodosim.SimulatorSync)

	js.Global().Set("_dodo_simulator_run", run())
	js.Global().Set("_dodo_simulator_stop", stop())
	js.Global().Set("_dodo_simulator_update_keys", updateKeys())

	c = make(chan bool)
	<-c
	fmt.Println("WASM: closing down")
}

func run() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		didStop = false

		fmt.Println("WASM: .run()")

		fmt.Println(args)

		ctx := args[0]
		wr := new(WebRenderer)
		wr.New(ctx)
		simulator.Renderer = wr

		speaker := new(WebSpeaker)
		speaker.New()
		simulator.Speaker = speaker

		firmware := jsToByteArray(args[1])
		fram := jsToByteArray(args[2])
		callback := args[3]

		simulator.CyclesPerFrame = func(cycles uint64) {
			if !didStop {
				callback.Invoke(uint32(cycles))
			}
		}

		simulator.SimulateSyncInit(firmware, fram, framFlusher)

		speaker.Enable()

		// Measure 2 seconds worth, to figure out delay
		start := time.Now()
		for i := 0; i < 40; i++ {
			simulator.PumpClock("")
		}
		elapsed := time.Since(start)
		fmt.Println(elapsed)

		// Calculate the delay necessary to achieve 20FPS
		delay := 1.0
		seconds := elapsed.Seconds()
		if seconds < 1.5 {
			delay = ((2.0 - seconds) / 40.0) * 1000.0
		}

		fmt.Println("Delay of :", delay)

		every(time.Millisecond*time.Duration(delay), func() bool {
			simulator.PumpClock(keys)

			if didStop {
				speaker.Disable()
				c <- true
			}

			return !didStop
		})

		return 0
	})
	return jsonFunc
}

func stop() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		didStop = true
		return 0
	})
	return jsonFunc
}

func updateKeys() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		keys = args[0].String()
		return 0
	})
	return jsonFunc
}

func jsToByteArray(arr js.Value) []byte {
	len := arr.Get("length").Int()
	buf := make([]byte, len)
	js.CopyBytesToGo(buf, arr)
	return buf
}

func every(duration time.Duration, fn func() bool) {
	time.AfterFunc(duration, func() {
		if !fn() {
			return
		}
		every(duration, fn)
	})
}

func framFlusher(f *dodosim.Fram) {
	/*
		if activeProject != "" {
			// Save to persistant storage or cookie
			data := f.Data[4:68]
			log.Printf("Flushing: ", data)
			localStorage.Call("setItem", activeProject, data)
		}*/
}
