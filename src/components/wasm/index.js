import './wasm_exec.js'

const funcs = {
  async init() {    
    const go = new global.Go()
    const wasm = await import('./main.wasm')
    await WebAssembly.instantiate(wasm.default, go.importObject).then((result) => {
      go.run(result.instance)      
    })
  },
  async run() {
    return global._dodo_simulator_run()
  }
}

export default funcs