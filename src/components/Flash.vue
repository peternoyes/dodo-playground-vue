<template>
  <div>
    <p>
      Communication with Dodo requires a Serial Port connection and uses the WebSerial API
    </p>
    <p>
      Test the connection to report the system version of the connected Dodo. Flash to trasmit the game. Before proceeding with either option, ensure that the Dodo is on the home screen. 
    </p>
    <b-alert v-model="showStatus" :variant="statusVariant">{{ status }}</b-alert>
    <b-progress v-if="flashing" :value="progress" max="8192" class="mb-3"></b-progress>
    <b-button-group>
      <b-button @click="getVersion">Test Connection</b-button>
      <b-button @click="flashGame">Flash</b-button>      
    </b-button-group>
  </div>
</template>
<script>
export default {
  name: 'Flash',
  props: ['binary'],
  data() {
    return {
      showStatus: false,
      status: '',
      statusVariant: 'primary',
      progress: 0,
      flashing: false
    }
  },
  methods: {
    readLine: async function(reader) {
      // 2s timeout
      let timeout = setTimeout(async function() {           
        await reader.cancel()
      }, 2000)

      let text = ''      
      let failure = false
      while (true) { // eslint-disable-line
        let { value, done } = await reader.read()                
        if (value) {
          if (value.endsWith('\0')) {
            done = true
            value = value.slice(0, -1)
          }          
          text += value
        } else {          
          failure = true
          break
        }

        if (done) {          
          break
        }
      }

      clearTimeout(timeout)

      return {text, failure}
    },
    readChar: async function(reader, ms = 2000) {
      // 2s timeout
      let timeout = setTimeout(async function() {           
        await reader.cancel()
      }, ms)
      
      let text = ''
      let failure = false      
      let { value } = await reader.read()       
      if (value && value.length === 1) {        
        text = value
      } else {          
        failure = true        
      }

      clearTimeout(timeout)

      return {text, failure}
    },
    getVersion: async function() {
      this.status = 'Requesting Port...'
      this.statusVariant = 'primary'
      this.showStatus = true
      this.progress = 0
      this.flashing = false

      try
      {
        let port = await navigator.serial.requestPort();

        await port.open({ baudRate: 9600 })

        let decoder = new TextDecoderStream() // eslint-disable-line
        let inputDone = port.readable.pipeTo(decoder.writable)
        let inputStream = decoder.readable;
        let reader = inputStream.getReader()

        let encoder = new TextEncoderStream() // eslint-disable-line
        let outputDone = encoder.readable.pipeTo(port.writable)
        let outputStream = encoder.writable
        let writer = outputStream.getWriter()

        writer.write('V')
        writer.releaseLock()

        this.status = 'Connecting...'

        let {text, failure} = await this.readLine(reader)

        if (failure) {
          // reader is already cancelled in event of failure
          this.statusVariant = 'danger'
          this.status = 'Dodo is not responding'
        } else {
          await reader.cancel()
          // Success!
          this.statusVariant = 'success'
          this.status = `Connected to Dodo with version: ${text}`
        }

        reader.releaseLock()
        
        await inputDone.catch(() => {})

        await outputStream.getWriter().close()
        await outputDone

        await port.close()
      }
      catch (err) {
        this.statusVariant = 'danger'
        this.status = err.message
      }
    },
    flashGame: async function() {
      this.status = 'Requesting Port...'
      this.statusVariant = 'primary'
      this.showStatus = true
      this.progress = 0
      this.flashing = true

      try
      {
        let port = await navigator.serial.requestPort();

        await port.open({ baudRate: 9600 })

        let decoder = new TextDecoderStream() // eslint-disable-line
        let inputDone = port.readable.pipeTo(decoder.writable)
        let inputStream = decoder.readable;
        let reader = inputStream.getReader()

        let writer = port.writable.getWriter()
        let data = new Uint8Array(['G'.charCodeAt(0)])
        writer.write(data)        

        this.status = 'Connecting...'

        let {text, failure} = await this.readChar(reader)
        if (failure) {
          // reader is already cancelled in event of failure
          this.statusVariant = 'danger'
          this.status = 'Dodo is not responding'
        } else {
          if (text != 'A') {
            await reader.cancel()
            this.statusVariant = 'danger'
            this.status = 'Incorrect message received'
          } else {
            // Success!            
            this.status = 'Flashing...'

            // Write Fram
            for (let i = 0; i < 8192; i++) {
              data = new Uint8Array([this.binary[i]])
              await writer.write(data)
              this.progress = i
            }

            this.status = 'Loading...'

            // Read again
            let {text, failure} = await this.readChar(reader, 10000) // Give enough time to flash            
            if (failure) {
              // reader is already cancelled in event of failure
              this.statusVariant = 'danger'
              this.status = 'Dodo is not responding'
            } else {
              if (text != 'A') {
                await reader.cancel()
                this.statusVariant = 'danger'
                this.status = 'Incorrect message received'
              } else {                
                await reader.cancel()
                this.statusVariant = 'success'
                this.status = 'Successfully Flashed Game'
              }
            }            
          }
        }

        writer.releaseLock()
        reader.releaseLock()
        
        await inputDone.catch(() => {})

        await port.close()
      }
      catch (err) {
        this.statusVariant = 'danger'
        this.status = err.message        
      }
    }
  }
}
</script>

<style scoped>
</style>