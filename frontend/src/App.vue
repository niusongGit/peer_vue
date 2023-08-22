<script setup>
import PeerList from './components/PeerList.vue'
import Peer from './components/Peer.vue'
//import Test from './components/Test.vue'
</script>

<template>


<!--   <Test/>-->
    <PeerList/>
    <el-drawer v-model="$store.state.drawer" title="节点详情管理" :with-header="false" size="55%">
        <Peer/>
    </el-drawer>
</template>