<template>
  <div>
    <h2 class="subtitle">
      IG engagement rate calculator
    </h2>

    <form class="form" v-on:submit.prevent>
      <div class="form-group">
        <label class="form-label" for="username">Username</label>
        <input  type="text" v-model="username" placeholder="thetwobohemians" class="form-control" id="username">
      </div>
      <button
          type="submit"
          @click="onSubmit"
          class="btn"
      >SUBMIT</button>
    </form>
    <h3 class="subtitle">
      Engagement rate: {{rate}} %
    </h3>
  </div>
</template>

<script>
import axios from "axios";
  export default {
    data() {
      return {
        rate: '',
        username: ''
      }
    },
    methods: {
      async onSubmit() {
          await axios.get(`http://localhost:3000/cal-engagement-rate/${this.username}`)
              .then(response => {
                this.rate = response.data * 100
              })
              .catch(e => {
                console.log(e)
              })
        }
    }
  }
</script>

<style scoped>

</style>