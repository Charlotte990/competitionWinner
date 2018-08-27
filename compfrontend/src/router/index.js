import Vue from 'vue'
import Router from 'vue-router'
import CompWinner from '@/components/CompWinner'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'CompWinner',
      component: CompWinner
    }
  ],
  mode: 'history'
})
