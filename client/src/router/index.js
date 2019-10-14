import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
import demo from '@/components/demo'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'HelloWorld',
      component: HelloWorld
    },
    {
      path: '/u',
      name: 'demo',
      component: demo
    }
  ],
  mode: "history" //把Router的mode修改为history模式,VueRouter默认的模式为HASH模式
})
