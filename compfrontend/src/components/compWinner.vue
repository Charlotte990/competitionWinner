<template>
  <div>
    <h1>{{ msg }}</h1>
    <button v-if="!foundWinner" @click="findWinner()">Select Winner</button>
    <h2 v-else>{{ winnerName }}</h2>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'CompWinner',
  data () {
    return {
      msg: 'And the winner is...',
      foundWinner: false,
      winnerName: '',
      message: '',
      info: ''
    }
  },
  methods: {
    findWinner () {
      console.log('in func')
      axios
        .get('http://localhost:8080/entries?query={randomSelection{id,username}}')
        .then(response => (this.info = response.data.data.randomSelection))
        .catch(error => {
          this.message = 'api not working'
          console.log(error)
        })
      this.foundWinner = true
      this.winnerName = this.info.username
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
