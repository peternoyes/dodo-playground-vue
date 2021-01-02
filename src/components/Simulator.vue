<template>
  <div>
    <canvas id="sim-canvas" width="256" height="128" style="border:1px solid #000000;"></canvas>
  </div>
</template>

<script>
import dodo from './wasm/index.js'

export default {
  name: 'Simulator',
  data() {
    return {      
    }
  },
  async mounted() {
    let canvas = document.getElementById('sim-canvas')
    let ctx = canvas.getContext('2d')    

    await dodo.init()
    await dodo.run(ctx)
  },
}
</script>

<style scoped>
</style>
