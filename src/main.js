import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import 'github-markdown-css/github-markdown.css'

import Vue from 'vue'
import App from './App.vue'
import VueRouter from 'vue-router'
import VueCookies from 'vue-cookies'
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'
import VueShowdownPlugin from 'vue-showdown';
import axios from 'axios'
import VueAxios from 'vue-axios'
import router from './router'

Vue.use(VueRouter)
Vue.use(VueCookies)
Vue.use(BootstrapVue)
Vue.use(IconsPlugin)

axios.interceptors.request.use(
  (config) => {
    let token = VueCookies.get('Authorization')

    if (token) {
      config.headers['Authorization'] = token
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

Vue.use(VueAxios, axios)

Vue.use(VueShowdownPlugin, {
  // set default flavor of showdown
  flavor: 'github',
  // set default options of showdown (will override the flavor options)
  options: {
    emoji: false,
  },
});

Vue.config.productionTip = false

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
