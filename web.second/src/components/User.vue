<template>
  <el-card>
    <el-row>
      <el-col>
        <el-input
          type="textarea"
          autosize
          placeholder="message"
          v-model="message">
        </el-input>
      </el-col>
      <el-col>
        <el-input
          type="textarea"
          autosize
          placeholder="policy"
          v-model="policy">
        </el-input>
      </el-col>
      <el-col>
        <el-button 
          type="success" 
          icon="el-icon-check" 
          circle
          @click="submit()"
        ></el-button>
      </el-col>
    </el-row>
  </el-card>
</template>

<script>
/* eslint-disable */
export default {
  name: 'User',
  props: {
    
  },
  data() {
    return {
      message: '',
      policy: '',
    }
  },
  methods:{
    submit(){
      console.log("submit()")
      console.log(this)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

  .el-col{
    text-align: center;
    margin-bottom: 10px;
  }
  
</style>
