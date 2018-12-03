<template>
  <el-header>
    <el-row>
      <el-col :span='24'>
        {{info}}
      </el-col>
    </el-row>
  </el-header>
</template>

<script>
export default {
  name: 'Header',
  props: {
    info: String
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .el-col{
    text-align: center;
    font-size: 300%;
  }
</style>
