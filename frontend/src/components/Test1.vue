<template>
    <el-button-group class="ml-4">
        <el-button type="primary" @click="this.groupDialogFormVisible = true"><el-icon><Grid /></el-icon>添加分组</el-button>
        <el-button type="primary" @click="showFormDialog(null)"><el-icon><Place /></el-icon>添加节点</el-button>
        <el-button type="primary" @click="showTagDialog()"><el-icon><Management /></el-icon>管理分组</el-button>
    </el-button-group>
    <el-table ref="tableRef" :data="$store.state.peerList" style="width: 100%"  @filter-change="filterChange">
        <el-table-column prop="status" :min-width="30">
            <template  #default="scope">
                <el-switch size="small" v-model="scope.row.status" :active-value="true" :inactive-value="false" @change="switchChange(scope.row)"/>
            </template>
        </el-table-column>
        <el-table-column prop="group" label="分组" :filters="groupfilters" :filter-method="filterTag" filter-placement="bottom-end" :min-width="40" column-key="groups"/>
        <el-table-column prop="name" label="节点名" :min-width="40"/>
        <el-table-column prop="highest_block" label="高度">
            <template #default="scope">
                <el-text class="mx-1" :type="getCurrentBlockType(scope.row)">{{scope.row.current_block}}</el-text>/<el-text class="mx-1" :type="scope.row.status?'success':'info'" tag="b">{{scope.row.highest_block}}</el-text>
            </template>
        </el-table-column>
        <el-table-column prop="total_balance" label="余额"/>
        <el-table-column prop="types" label="角色" :min-width="40">
            <template #default="scope">
                <el-tag v-for="(num, type) in scope.row.types" :type="getTagType(scope.row.status?type:-1)" size="small">{{ getTagText(type) }}{{num>1?num:''}}</el-tag>
            </template>
        </el-table-column>
        <el-table-column prop="used_time" label="延迟" :min-width="30">
            <template #default="scope">
                <el-text class="mx-1" :type="getUsedTimeType(scope.row)">{{scope.row.used_time}}ms</el-text>
            </template>
        </el-table-column>
        <el-table-column fixed="right" label="Operations" :min-width="55">
            <template #default="scope">
                <el-button link type="primary" @click="showDrawer(scope.row)" ><el-icon><Expand /></el-icon></el-button>
                <el-button link type="primary" @click.prevent="showFormDialog(scope.row)"><el-icon><Edit /></el-icon></el-button>
                <el-popconfirm title="是否删除该节点?" confirm-button-text="是" cancel-button-text="否" @confirm="deleteRow(scope.row)">
                    <template #reference>
                        <el-button link type="primary"><el-icon><Delete /></el-icon></el-button>
                    </template>
                </el-popconfirm>

            </template>
        </el-table-column>
    </el-table>

</template>
<script>
import {unref} from 'vue'
import { ElMessage } from 'element-plus'
import {AddGroup,AddPeer,DelPeer,DelGroup,EditGroup, EditPeer, GetPeerList,GetPeerGroup, SetStatus} from '../../wailsjs/go/peer/Api'
export default {
    data() {
        return {
            timer:0,
            peerConfFormRef:null,
            tableRef:null,
            dialogFormVisible: false,
            groupDialogFormVisible: false,
            tagDialogVisible: false,
            peerConfForm: {
                id: '',
                group: '',
                name: '',
                ip: '',
                port: 0,
                username: '',
                password: '',
                status: false,
                highest_block: 0,
                current_block: 0,
                snapshot_height: 0,
                total_balance: 0,
                types: null,
                used_time: 0,
                updated_at:0,
                error:"",
                default_address: null,
                addresses:[],
            },
            editingItem: false,//节点配置弹窗是添加还是编辑
            formLabelWidth:'100px',
            groupfilters:[ //筛选功能的分组列表
            ],
            selectedGroups:[],
            groupForm:{
                group:''
            },
            dynamicTags: [
                // { value: 'Tag 1', editable: false,tmp:'',popover:false  }, value分组名，editable编辑框是否显示，tmp用于保存修改之前的分组名在修改时作对比，popover控制删除提示框是否显示
                // { value: 'Tag 2', editable: false,tmp:'',popover:false  },
                // { value: 'Tag 3', editable: false ,tmp:'',popover:false }
            ],
            newTagInputVisible: false,
            newTagInputValue: '',
            peerCurrentBocks:{} //用来记录区块高度是否在增长
        };
    },
    created: function () {
        this.loadPeerList();
        this.loadGroups()
    },
    mounted() {
        this.timer = setInterval(() => {
            this.loadPeerList(); //500毫秒更新一次节点列表数据
        }, 500)
    },
    beforeUnmount() {
        clearInterval(this.timer)
    },
    methods: {
        loadPeerList() {
            GetPeerList().then((result) => {
                if (result != null) {
                    if (this.selectedGroups.length > 0) {
                        // 过滤出 selectedGroups 中存在的 group
                        const selected = result.filter(item => this.selectedGroups.includes(item.group));
                        // 将相同 group 的节点对象放在一起
                        const grouped = selected.reduce((obj, item) => {
                            obj[item.group] = [...(obj[item.group] || []), item];
                            return obj;
                        }, {});
                        // 将 grouped 的值放在数组前面,实现把同一组的数据放在一起
                        this.$store.state.peerList = [].concat(...Object.values(grouped), ...result.filter(item => !this.selectedGroups.includes(item.group)));
                    } else {
                        this.$store.state.peerList = result;
                    }
                }
            }).catch(error => {
                // 处理错误
                console.log(error)
            });
        },
        loadGroups() {
            GetPeerGroup().then((result) => {
                const groups = []
                for (var item in result) {
                    groups.push({
                        text: result[item], value: result[item]
                    })
                    // const index = this.groupfilters.findIndex((data) => data.text === item);
                    // if (index === -1) {
                    //     this.groupfilters.push({
                    //         text: result[item],value: result[item]
                    //     })
                    // }
                }
                this.groupfilters = groups
            }).catch(error => {
                // 处理错误
                console.log(error)
            });
        },
        getCurrentBlockType(row){
            if (!row.status){
                return 'info'
            }
            if (this.peerCurrentBocks[row.id] == null){
                this.peerCurrentBocks[row.id] = [row.current_block]
            }else {
                const timex = new Date().getTime()
                if (this.peerCurrentBocks[row.id].length>10){
                    if (timex%10 == 0){
                        this.peerCurrentBocks[row.id].push(row.current_block)
                        this.peerCurrentBocks[row.id].splice(0,1)
                    }
                    const sum = this.peerCurrentBocks[row.id].reduce((total,num)=>total+num)
                    if ((sum/this.peerCurrentBocks[row.id].length == row.current_block)){
                        return 'danger'
                    }
                }else {
                    this.peerCurrentBocks[row.id].push(row.current_block)
                }
            }
            return 'primary'
        }
    }
};
</script>

<style scoped>
.el-tag + .el-tag {
    margin-left: 10px;
}
.button-new-tag {
    margin-left: 10px;
    height: 32px;
    line-height: 30px;
    padding-top: 0;
    padding-bottom: 0;
}
.input-new-tag {
    width: 90px;
    margin-left: 10px;
    vertical-align: bottom;
}
</style>
上面的代码请在保证逻辑不变的情况下优化getCurrentBlockType方法