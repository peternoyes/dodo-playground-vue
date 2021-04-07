<template>
  <div id="app">
    <b-container fluid class="d-flex flex-column flex-grow-1 vh-100 px-0 overflow-hidden">    
      <b-navbar toggleable="lg" type="dark" variant="dark" class="flex-shrink-0">
        <b-navbar-brand>Dodo Playground</b-navbar-brand>
        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
        <b-collapse id="nav-collapse" is-nav>
          <b-navbar-nav>
            <b-nav-item-dropdown :text="displayLanguage">
              <b-dropdown-item v-for="option in languageOptions"
                :key="option.value"
                :value="option.value"
                @click="language = option.value">{{ option.text }}</b-dropdown-item>
            </b-nav-item-dropdown>            
            <b-nav-item-dropdown :text="firmware">
              <b-dropdown-item v-for="version in firmwareVerisons"
                :key="version"
                :value="version"
                @click="firmware = version">{{ version }}</b-dropdown-item>
            </b-nav-item-dropdown>
            <b-nav-form>            
              <b-button pill variant="info" class="ml-2" :disabled="!showEditor" @click="run"><span><b-icon icon="play" /> Run</span></b-button>
              <b-button pill v-if="hasToken" variant="info" class="ml-2" :disabled="!showEditor" ><span><b-icon icon="cloud-check" /> Save</span></b-button>
              <b-button pill variant="success" class="ml-2" :disabled="!showEditor"  @click="download"><span><b-icon icon="download" /> fram.bin</span></b-button>
              <b-button pill v-if="canFlash" variant="warning" class="ml-2" :disabled="!showEditor" @click="flash" ><span><b-icon icon="lightning" /> Flash</span></b-button>
            </b-nav-form>          
          </b-navbar-nav>
          <b-navbar-nav class="ml-auto">
            <b-nav-form>
              <b-button v-if="!hasToken" variant="dark" @click="login"><span><b-icon icon="github" /> Sign in with GitHub</span></b-button>
              <b-button v-if="hasToken" variant="info" @click="logout">Sign Out</b-button>
            </b-nav-form>
          </b-navbar-nav>
        </b-collapse>
      </b-navbar>
      <b-row class="flex-grow-1 overflow-hidden" no-gutters>
        <b-col cols="8" class="mh-100 overflow-auto py-2">          
            <editor v-if="showEditor" v-model="content" @init="editorInit" lang="c_cpp" theme="chrome" width="100%" height="100%"></editor>            
            <div v-if="!showEditor" class="h-100 w-100 text-center" style="display: table;">              
              <span style="display: table-cell; vertical-align: middle;">Please load or create a new project</span>
            </div>
        </b-col>
        <b-col class="mh-100 overflow-auto">
          <div class="border-left px-2 py-2">    
            <Documentation v-if="!hasToken" />
          </div>
          <b-tabs v-if="hasToken" class="border-left py-2 px-2" content-class="mt-3">
            <b-tab title="Projects" active><Projects v-on:project-loaded="loadProject" v-bind:user="user" /></b-tab>
            <b-tab title="Documentation"><Documentation /></b-tab>
          </b-tabs>
        </b-col>        
      </b-row>         
      <b-row :class="`text-white px-2 flex-shrink-0 ${status.background}`">
        <b-col cols="12" class="py-1">
          {{ status.message }}
        </b-col>
      </b-row>
    </b-container>
    <b-modal id="bv-modal-simulator"
      size="lg"
      @hide="stop"
      header-bg-variant="dark"
      header-text-variant="light"
      centered hide-footer title="Simulator">
      <Simulator v-bind:firmware="firmware" v-bind:binary="binary" v-bind:stopped="stopped" />
    </b-modal>
    <b-modal id="bv-modal-flash"
      size="lg"
      header-bg-variant="dark"
      header-text-variant="light"
      centered hide-footer title="Flash">
        <p>
          After choosing the appropriate COM port, ensure that Dodo is on the home screen before selecting 'Start' which will immediately flash Dodo. After successfully flashing Dodo, it will take a few additional seconds to load the game.
        </p>
        <b-button @click="openPort">Open Port</b-button>
    </b-modal>
  </div>
</template>

<script>

import Documentation from '../components/Documentation.vue'
import Projects from '../components/Projects.vue'
import Simulator from '../components/Simulator.vue'
import DefaultCode from '../assets/default.c'
import helpers from '../helpers.js'

export default {
  name: 'App',
  components: {    
    Documentation,
    Projects,
    Simulator,
    editor: require('vue2-ace-editor')    
  },
  data() {
    return {
      status: {
        message: "Ready. Click 'Run' to try your game in the simulator.",
        background: 'bg-success'
      },
      user: {},
      title: '',
      content: '',
      binary: null,
      stopped: false,
      language: 'c',
      token: '',
      languageOptions: [
        { value: 'c', text: 'C'},
        { value: 'assembly', text: 'Assembly'}      
      ],
      firmware: '1.0.1',
      firmwareVerisons: [
        '1.0.1',
        '1.1.0_beta2'
      ]
    }
  },
  computed: {
    displayLanguage: function() {
      return this.languageOptions.find(o => o.value == this.language).text
    },
    hasToken: function() {
      return this.token && this.token.length !== 0
    },
    showEditor: function() {
      return !this.hasToken || this.title && this.title.length !== 0
    },
    hasProject: function() {
      return this.title && this.title.length !== 0
    },
    canFlash: function() {
      return 'serial' in navigator;
    }
  },
  async mounted() {    
    this.token = this.$cookies.get('Authorization')
    let code = this.$route.query.code
    if (code) {
      this.axios.get(`${process.env.VUE_APP_API_URL}/code/${code}`)
        .then((response) => {
          this.content = response.data
          this.language = response.headers['x-language']        
        })
        .catch((reason) => {          
          // todo change background color
          this.setStatus(reason.message, 'bg-error')          
        })

    } else {
      if (!this.hasToken) { 
        this.content = DefaultCode
        this.language = 'c'
      } else {
        this.language = 'c'
      }
    }

    if (this.hasToken) {
      this.axios.get(`${process.env.VUE_APP_API_URL}/user`)
        .then((response) => {
          this.user = response.data          
        })
        .catch((reason) => {
          console.log(reason)
          // How to handle error here?
        })
    } else {
      this.user = {}
    }
  },  
  methods: {
    editorInit: function () {
      require('brace/ext/language_tools') //language extension prerequsite...
      require('brace/mode/c_cpp')                
      require('brace/mode/less')
      require('brace/theme/chrome')
      require('brace/snippets/javascript') //snippet
    },
    compileCode: async function() {      
      let config = {
        headers: {
          'Content-Type': 'application/text',
          'X-Language': this.language,
          'X-Version': this.firmware          
        }
      }

      return await this.axios.post(`${process.env.VUE_APP_API_URL}/build`, this.content, config)
    },
    setStatus: function(message, background) {
      this.status.message = message
      this.status.background = background
    },
    run: async function() {
      // TODO Save      
      this.setStatus('Compiling...', 'bg-info')
     
      try {
        let response = await this.compileCode()
        this.binary = helpers.b64toUint8Array(response.data.binary)

        this.setStatus('Loading Simulator...', 'bg-success')

        // TODO: load from local storage for Fram (What does this mean??)
        
        // Start Simulator
        this.stopped = false
        this.$bvModal.show('bv-modal-simulator')
      } catch (error) {
        this.setStatus(error.response.data, 'bg-danger')
      }
    },
    stop: async function() {
      this.stopped = true
      this.setStatus('Simulator Stopped.', 'bg-success')
    },
    download: async function() {
      this.setStatus('Compiling...', 'bg-info')
      
      try {
        let response = await this.compileCode()              
        
        let blob = helpers.b64toBlob(response.data.binary, "application/octet-string")

        let link = document.createElement("a")        
        link.href = window.URL.createObjectURL(blob)        
        link.download = "fram.bin"
        link.click()

        this.setStatus('Successfully created fram.bin', 'bg-success')
      } catch (error) {
        this.setStatus(error.response.data, 'bg-danger')
      }
    },
    flash: async function() {
      let response = await this.compileCode()
      this.binary = helpers.b64toUint8Array(response.data.binary)

      this.setStatus('Flashing...', 'bg-success')

      this.$bvModal.show('bv-modal-flash')
    },
    openPort: async function() {
      let port = await navigator.serial.requestPort();

      await port.open({ baudrate: 9600 })

    },
    login: function() {
      console.log('Login...')
      window.location.href = `${process.env.VUE_APP_API_URL}/login`
    },
    logout: function() {
      console.log('Logout...')
      window.location.href = `${process.env.VUE_APP_API_URL}/logout`
    },
    loadProject: function(project) {
      this.title = project.title   
      this.content = project.source
      this.language = project.language
      this.firmware = project.version            
    }
  },
}
</script>

<style>
#app {
  font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: left;
}
</style>
