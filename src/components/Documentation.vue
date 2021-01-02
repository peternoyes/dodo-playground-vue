<template>
  <div>
    <VueShowdown :markdown="apiContent" class="markdown-body"></VueShowdown>
  </div>
</template>

<script>
import api from '../assets/api.md'

export default {
  name: 'Documentation',
  data() {
    return {
      apiContent: '',
    }
  },
  mounted() {
    this.apiContent = api
  },
}
</script>

<style scoped>
</style>
