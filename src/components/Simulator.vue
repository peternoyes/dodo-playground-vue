<template>
  <div 
    tabindex="0"
    v-on:keydown.a="A = true"
    v-on:keyup.a="A = false"
    v-on:keydown.b="B = true"
    v-on:keyup.b="B = false"
    v-on:keydown.up="Up = true"
    v-on:keyup.up="Up = false"
    v-on:keydown.down="Down = true"
    v-on:keyup.down="Down = false"
    v-on:keydown.left="Left = true"
    v-on:keyup.left="Left = false"
    v-on:keydown.right="Right = true"
    v-on:keyup.right="Right = false"
    >
    <center>
      <img src="../assets/dodo_large.png" width=764 usemap="#game_pad">
      <map name="game_pad">
        <area shape="rect" coords="509,374,573,439" href="#" @mousedown="A = true" @mouseup="A = false" />
        <area shape="rect" coords="604,374,668,439" href="#" @mousedown="B = true" @mouseup="B = false" />
        <area shape="poly" coords="137,278,184,323,137,369,92,323" href="#" @mousedown="Up = true" @mouseup="Up = fase" />
        <area shape="poly" coords="137,393,184,439,137,484,92,439" href="#" @mousedown="Down = true" @mouseup="Down = false" />
        <area shape="poly" coords="80,336,125,380,80,427,35,380" href="#" @mousedown="Left = true" @mouseup="Left = false" />
        <area shape="poly" coords="194,336,241,380,194,427,150,380" href="#" @mousedown="Right = true" @mouseup="Right = false" />
      </map>
      <canvas
        id="sim-canvas"
        style="position:absolute; top:104px; left:271px; z-index:1"
        width="256" 
        height="128"></canvas>
    </center>
    <p>Frames Per Cycle: {{ cycles }}</p>
  </div>
</template>

<script>
import dodo from './wasm/index.js'

export default {
  name: 'Simulator',
  props: ['firmware', 'binary', 'id', 'stopped'],
  data() {
    return {
      cycles: 0,
      A: false,
      B: false,
      Up: false,
      Down: false,
      Left: false,
      Right: false
    }
  },
  computed: {
    keyState: function() {
      let keys = ""
      if (this.A) {
        keys += "A"
      }
      if (this.B) {
        keys += "B"
      }
      if (this.Up) {
        keys += "U"
      }
      if (this.Down) {
        keys += "D"
      }
      if (this.Left) {
        keys += "L"
      }
      if (this.Right) {
        keys += "R"
      }
      return keys
    }
  },
  watch: {
    keyState: async function(val) {
      await dodo.updateKeys(val)
    },
    stopped: async function(val) {
      if (val) {        
        await dodo.stop()
      }
    }
  },
  async mounted() {
    let canvas = document.getElementById('sim-canvas')
    let ctx = canvas.getContext('2d')

    var buffer
    switch (this.firmware) {
        case '1.0.1': 
            buffer = await import('../assets/firmware_1.0.1.bin')
            break
        case '1.1.0_beta2':
            buffer = await import('../assets/firmware_1.1.0_beta2.bin')
            break
    }

    var firmwareBytes = new Uint8Array(buffer.default)

    let binary = this.binary
    let id = this.id

    // Load any persistent data from localStorage
    let data = localStorage.getItem(id)
    if (data) {
      let fram = JSON.parse(`[${data}]`)
      for (let i = 0; i < 64; i++) {
        binary[i+4] = fram[i]
      }
    }

    await dodo.init()

    // Callback is called whenver persistant data needs to be written to localStorage
    await dodo.setFramFlusher(function(data) {
      localStorage.setItem(id, data)
    })

    await dodo.run(ctx, firmwareBytes, binary, (cycles) => {
      this.cycles = cycles
    })
  },
  async unmounted() {      
      await dodo.stop()
  }
}
</script>

<style scoped>
</style>
