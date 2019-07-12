<template>
  <el-container>

    <el-header>
      <Header 
      info='ABE'
      />
    </el-header>

    <el-main>
      <User/>
    </el-main>

    <el-footer>
      <Footer/>
    </el-footer>

  </el-container>
</template>

<script>
import Header from './components/Header.vue';
import User from './components/User.vue';
import Footer from './components/Footer.vue';
/* eslint-disable */
export default {
  name: 'app',
  data(){
    return{
      name:"container",
    }
  },
  components: {
    Header,
    Footer,
    User,
  },
  methods:{},
  created(){},
  mounted(){
    console.log(this.name)
  },
}
</script>

<style>
/* #app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
} */
</style>
