<template>
  <el-footer>
    <el-row>
      <el-col :span='24'>
        @ Fudan BlockChain Lab.
      </el-col>
    </el-row>
  </el-footer>
</template>

<script>
export default {
  name: 'Footer',
  props: {

  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .el-col{
    text-align: center;
    font-size: 80%;
    color: grey;
  }
</style>
