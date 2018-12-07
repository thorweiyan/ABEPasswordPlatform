import Vue from 'vue';
import axios from 'axios';
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import App from './App.vue';

Vue.config.productionTip = false;

Vue.use(ElementUI);

/* eslint-disable */

new Vue({
  render: h => h(App),
  data(){
    return{
      name:"app name",
    }
  },
  methods:{
    async getinfo(){
      console.log("hello app");
      console.log(this.name);
      await axios.post('/user', {
        firstName: 'Fred',
        lastName: 'Flintstone'
      })
      .then(function (response) {
        console.log(response);
      })
      .catch(function (error) {
        console.log(error);
      });
    },
  },

  created() {},
  beforeCreate() {},
  beforeMount() {},
  mounted() {
    this.getinfo();
  },

}).$mount('#app');
