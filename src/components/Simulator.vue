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
    :style="wrapperStyle"
    >
    <center>
      <div v-responsive="['hidden-all', 'lg', 'xl']" v-view="largeHandler">
        <img src="../assets/dodo_large.png" width=764 usemap="#game_pad">      
        <map name="game_pad">
          <area shape="rect" coords="509,374,573,439" href="#" @mousedown="A = true" @mouseup="A = false" />
          <area shape="rect" coords="604,374,668,439" href="#" @mousedown="B = true" @mouseup="B = false" />
          <area shape="poly" coords="137,278,184,323,137,369,92,323" href="#" @mousedown="Up = true" @mouseup="Up = false" />
          <area shape="poly" coords="137,393,184,439,137,484,92,439" href="#" @mousedown="Down = true" @mouseup="Down = false" />
          <area shape="poly" coords="80,336,125,380,80,427,35,380" href="#" @mousedown="Left = true" @mouseup="Left = false" />
          <area shape="poly" coords="194,336,241,380,194,427,150,380" href="#" @mousedown="Right = true" @mouseup="Right = false" />
        </map>
      </div>
      <canvas
        id="sim-canvas"
        :style="canvasStyle"
        width="256" 
        height="128"></canvas>
      <div v-responsive="['hidden-lg', 'hidden-xl']" style="position: relative; width: 256px; height: 150px;">
        <button class="gameButton gameButton" style="top: 60px; left: 140px;"
          @mousedown="A = true" @mouseup="A = false"
          @touchstart="A = true" @touchend="A = false">A</button>
        <button class="gameButton gameButton" style="top: 60px; left: 200px;"
          @mousedown="B = true" @mouseup="B = false"
          @touchstart="B = true" @touchend="B = false">B</button>
        <button class="gameButton gameButtonRotated" style="top: 25px; left: 40px;"
          @mousedown="Up = true" @mouseup="Up = false"
          @touchstart="Up = true" @touchend="Up = false" />
        <button class="gameButton gameButtonRotated" style="top: 60px; left: 5px;"
          @mousedown="Left = true" @mouseup="Left = false"
          @touchstart="Left = true" @touchend="Left = false" />
        <button class="gameButton gameButtonRotated" style="top: 60px; left: 75px;"
          @mousedown="Right = true" @mouseup="Right = false"
          @touchstart="Right = true" @touchend="Right = false" />
        <button class="gameButton gameButtonRotated" style="top: 95px; left: 40px;"
          @mousedown="Down = true" @mouseup="Down = false"
          @touchstart="Down = true" @touchend="Down = false" />
      </div>        
    </center>
    <div>
      <p>Frames Per Cycle: {{ cycles }}</p>
    </div>
  </div>
</template>

<script>
import dodo from './wasm/index.js'

export default {
  name: 'Simulator',
  props: ['firmware', 'binary', 'id', 'stopped'],
  data() {
    return {
      isLarge: true,
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
    canvasStyle: function() {
      if (this.isLarge) {
        return 'position:absolute; top:104px; left:271px; z-index:1'
      } else {
        return 'position:relative; top:10px; z-index:1'
      }
    },
    wrapperStyle: function() {
      if (this.isLarge) {
        return ''
      } else {
        return 'background-color: #097600'
      }
    },
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
  },
  methods: {
    largeHandler(e) {
      if (e.type === 'enter') {
        this.isLarge = true
      } else if (e.type === 'exit') {
        this.isLarge = false
      }
    }
  }
}
</script>

<style scoped>
  .gameButton {
    -webkit-appearance: none !important;
    -moz-appearance: none !important;
    appearance: none !important;
    -webkit-touch-callout:none !important;
    -webkit-user-select:none !important;
    -khtml-user-select:none !important;
    -moz-user-select:none !important;
    -ms-user-select:none !important;
    user-select:none !important;
    -webkit-tap-highlight-color:rgba(0,0,0,0) !important;    
    position: absolute;
    height:44px;
    width:44px;
    border: 1px solid #ff0000;
    color: #ffffff;
    background-color: #ff0000;
  }
  .gameButtonRotated {
    transform: rotate(45deg);
  }
</style>