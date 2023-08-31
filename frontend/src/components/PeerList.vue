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
                <el-text class="mx-1" :type="scope.row.status?'primary':'info'">{{scope.row.current_block}}</el-text>/<el-text class="mx-1" :type="scope.row.status?'success':'info'" tag="b">{{scope.row.highest_block}}</el-text>
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



    <el-dialog v-model="dialogFormVisible" title="节点配置" destroy-on-close center draggable>
        <el-form  :model="peerConfForm" ref="peerConfFormRef" :rules="rules">

            <el-form-item label="分组" :label-width="formLabelWidth"  prop="group">
                <el-select v-model="peerConfForm.group" placeholder="Please select a group">
                    <el-option
                        v-for="(item) in groupfilters"
                        :key="item.value"
                        :label="item.value"
                        :value="item.value"
                    />
                </el-select>
            </el-form-item>
            <el-form-item label="节点名" :label-width="formLabelWidth" prop="name">
                <el-input v-model="peerConfForm.name" autocomplete="off" />
            </el-form-item>
            <el-form-item label="IP" :label-width="formLabelWidth"  prop="ip">
                <el-input v-model="peerConfForm.ip" autocomplete="off" />
            </el-form-item>
            <el-form-item label="Port" :label-width="formLabelWidth" prop="port">
                <el-input v-model="peerConfForm.port" autocomplete="off" type="number"/>
            </el-form-item>
            <el-form-item label="RPC用户名" :label-width="formLabelWidth" prop="username">
                <el-input v-model="peerConfForm.username" autocomplete="off" />
            </el-form-item>
            <el-form-item label="RPC密码" :label-width="formLabelWidth" prop="password">
                <el-input v-model="peerConfForm.password" autocomplete="off" />
            </el-form-item>
        </el-form>
        <template #footer>
      <span class="dialog-footer">
        <el-button @click="closeFormDialog">取消</el-button>
        <el-button type="primary" @click="peerConfFormSubmit">
          保存
        </el-button>
      </span>
        </template>
    </el-dialog>




    <el-dialog v-model="groupDialogFormVisible" width="40%" title="添加分组" destroy-on-close  center draggable>
        <el-form :model="groupForm">
            <el-form-item label="分组名" :label-width="formLabelWidth" prop="group">
                <el-input v-model="groupForm.group" autocomplete="off" />
            </el-form-item>
        </el-form>
        <template #footer>
      <span class="dialog-footer">
        <el-button @click="groupDialogFormVisible = false">取消</el-button>
        <el-button type="primary" @click="groupFormSubmit">
          保存
        </el-button>
      </span>
        </template>
    </el-dialog>


    
    <el-dialog
        v-model="tagDialogVisible"
        title="分组管理"
        destroy-on-close
        center
        draggable
    >
            <div>
                <el-tag
                    :key="index"
                    v-for="(tag, index) in dynamicTags"
                    closable
                    :disable-transitions="false"
                    @close="showPopover(tag)"
                >
                    <span v-if="!tag.editable" @click="tagHandleEdit(index)">{{ tag.value }}</span>
                    <el-input
                        v-else
                        v-model="tag.value"
                        size="small"
                        ref="editInput"
                        @keyup.enter.native="$event.target.blur()"
                        @blur="tagHandleInputConfirm(index)"
                    ></el-input>

                    <el-popover :visible="tag.popover" placement="bottom" :width="160"  trigger="manual">
                        <p>删除该分组和该分组下所有节点?</p>
                        <div style="text-align: right; margin: 0">
                            <el-button size="small" text @click="tag.popover = false">否</el-button>
                            <el-button size="small" type="primary" @click="tagHandleClose(tag)"
                            >是</el-button>
                        </div>
                        <template #reference><span></span></template>
                    </el-popover>

                </el-tag>
                <!--
                    *  @keyup.enter.native="$event.target.blur()" 表示使用keyup（回车键）事件触发blur（失去焦点）事件，为了防止回车后同时触发两个事件而调用两次newTagHandleInputConfirm方法
                  -->

                <el-input
                    class="input-new-tag"
                    v-if="newTagInputVisible"
                    v-model="newTagInputValue"
                    ref="saveTagInput"
                    size="small"
                    @keyup.enter.native="$event.target.blur()"
                    @blur="newTagHandleInputConfirm"
                ></el-input>
                <el-button v-else class="button-new-tag" size="small" @click="showNewTagInput">+ 新增分组</el-button>
            </div>
    </el-dialog>
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
            rules : {
                group: [
                    { required: true, message: '请输入组名', trigger: 'blur' },
                ],
                name: [
                    { required: true, message: '请输入节点名', trigger: 'blur' },
                ],
                ip: [
                    { required: true, message: '请输入ip', trigger: 'blur' },
                ],
                port: [
                    { required: true, message: '请输入端口', trigger: 'blur' },
                    {pattern:  /^[1-9]\d*$/,"message": "只能输入大于0的正整数", trigger: 'blur'}
                ],
                username: [
                    { required: true, message: '请输入RPC用户名', trigger: 'blur' }
                ],
                password: [
                    { required: true, message: '请输入RPC密码', trigger: 'blur' },
                ]
            },
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
            newTagInputValue: ''
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
        loadPeerList(){
            GetPeerList().then((result) => {
                if (result != null){
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
        loadGroups(){
            GetPeerGroup().then((result) => {
                const groups = []
                for (var item in result) {
                    groups.push({
                                text: result[item],value: result[item]
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
        switchChange(item){
            SetStatus(item.id,item.status).then((result) => {
                console.log('id :'+item.id+' status:'+item.status)
            }).catch(error => {
                // 处理错误
                ElMessage.error(error)
            });
        },
        showFormDialog(item) {
            this.dialogFormVisible = true;
            if (item) {
                this.editingItem = true;
                this.peerConfForm = item;
            } else {
                this.editingItem = false;
                this.peerConfForm.id='';
                this.peerConfForm.group='';
                this.peerConfForm.name='';
                this.peerConfForm.ip='';
                this.peerConfForm.port=0;
                this.peerConfForm.username='';
                this.peerConfForm.password='';
                this.peerConfForm.status=false;
                this.peerConfForm.highest_block=0;
                this.peerConfForm.current_block=0;
                this.peerConfForm.snapshot_height=0;
                this.peerConfForm.total_balance=0;
                this.peerConfForm.types=null;
                this.peerConfForm.used_time=0;
                this.peerConfForm.updated_at=0;
                this.peerConfForm.error='';
                this.peerConfForm.default_address=null;
                this.peerConfForm.addresses=[];
            }
        },
        closeFormDialog() {
            this.dialogFormVisible = false;
        },
        peerConfFormSubmit() {
            const formRef = unref(this.$refs.peerConfFormRef);
            if (!formRef) return
            //通过ref的值触发验证
            formRef.validate((valid) => {
                if (valid) {
                    const newItem = {
                        id: this.peerConfForm.id,
                        group: this.peerConfForm.group,
                        name: this.peerConfForm.name,
                        ip: this.peerConfForm.ip,
                        port: parseInt(this.peerConfForm.port.toString()),
                        username: this.peerConfForm.username,
                        password: this.peerConfForm.password,
                        // status: this.peerConfForm.status,
                        // highest_block: this.peerConfForm.highest_block,
                        // current_block: this.peerConfForm.current_block,
                        // balance: this.peerConfForm.balance,
                        // types: this.peerConfForm.types,
                        // delay: this.peerConfForm.delay,
                        // is_del: this.peerConfForm.is_del,
                    };
                    if (this.editingItem) {
                        EditPeer(newItem).then((result) => {
                            this.peerConfForm = result;
                            this.loadPeerList();
                            this.closeFormDialog();
                        }).catch(error => {
                            // 处理错误
                            ElMessage.error(error)
                        });
                    } else {
                        AddPeer(newItem).then((result) => {
                            //this.$store.state.peerList.push(result)
                            this.loadPeerList();
                            this.closeFormDialog();
                    }).catch(error => {
                            // 处理错误
                            ElMessage.error(error)
                        });
                    }
                } else {
                    console.log("未通过");
                }
            });
        },
        deleteRow(item) {
            const index = this.$store.state.peerList.findIndex((data) => data.id === item.id);
            if (index !== -1) {
                this.$store.state.peerList.splice(index, 1);
                DelPeer(item.id).then((result) => {
                    console.log('删除id :'+item.id+' result:'+result)
                }).catch(error => {
                    // 处理错误
                    ElMessage.error(error)
                });
            }
        },
        filterTag(value,row,column){
            this.selectedGroups = column.filteredValue //记录选中的项
            return true //让筛选不起作用，因为在loadPeerList()中已经做了分组筛选
        },
        getTagType(type){
            if (type == 1) {
                return 'success';
            } else if (type == 2) {
                return 'danger';
            } else if (type == 3) {
                return 'warning';
            } else {
                return 'info';
            }
        },
        getTagText(type){
            if (type == 1) {
                return '见证人';
            } else if (type == 2) {
                return '社区节点';
            } else if (type == 3) {
                return '轻节点';
            } else {
                return '普通节点';
            }
        },
        getUsedTimeType(row){
            if (row.status=== false || row.used_time === 0) { //关闭状态置灰
                return 'info';
            } else if (0 < row.used_time && row.used_time <= 50) {
                return 'success';
            } else if (50 <row.used_time && row.used_time <= 200) {
                return 'warning';
            } else {
                return 'danger';
            }
        },
        showDrawer(row){
            const index = this.$store.state.peerList.findIndex((data) => data.id === row.id);
            if (index !== -1) {
                this.$store.state.peerInfoIndex = index //抽屉中应该展示的节点的Index
                this.$store.state.selectedAddress = row.default_address?row.default_address.address:'' //默认选中的地址
                this.$store.state.drawer = true //打开抽屉
                this.$store.commit('setAddressesOptions')//打开抽屉时加载级联地址数据
            }

        },
        filterChange(filters) { //将当前筛选所选中的分组记录下来
            this.selectedGroups= filters['groups']
        },
        groupFormSubmit() { //弹窗功能添加分组
            const group = this.groupForm.group.trim()

            if (group === '') {
                ElMessage.error("请填写分组名！")
                return false
            }
            AddGroup(group).then((result) => {
                this.groupfilters.push({text: group,value:group})
                this.groupDialogFormVisible = false
                this.loadGroups()
                this.groupForm.group = ''
            }).catch(error => {
                // 处理错误
                ElMessage.error(error)
                return false
            });

        },
        showTagDialog(){ //所有分组tag
            const tags = [];
            this.groupfilters.forEach(function (item, index) {
                tags.push({ value: item.value, editable: false,tmp:item.value,popover:false });
            });
            this.dynamicTags = tags
            this.tagDialogVisible=true
        },
        showPopover(tag) { //显示删除分组提示
            tag.popover = true;
        },
        tagHandleClose(tag) { //删除分组
            const index = this.dynamicTags.findIndex(t => t.value === tag.value);
            if (index !== -1) {
                tag.popover = false;
                //删除
                DelGroup(tag.value.trim()).then((result) => {
                    this.loadGroups()
                }).catch(error => {
                    // 处理错误
                    ElMessage.error(error)
                    return false
                });
                this.dynamicTags.splice(index, 1);
            }
        },

        showNewTagInput() { //显示分组添加修改框
            this.newTagInputVisible = true;
            this.$nextTick(() => {
                this.$refs.saveTagInput.focus();
            });
        },

        tagHandleEdit(index) { //显示分组修改框
            this.dynamicTags[index].editable = true;
            this.dynamicTags[index].tmp = this.dynamicTags[index].value;
            this.$nextTick(() => {
                const input = this.$refs.editInput[index];
                if (input) {
                    input.$el.querySelector('input').focus();
                }
            });
        },

        tagHandleInputConfirm(index) {//修改分组
            const tag = this.dynamicTags[index];
            tag.editable = false;
            if (tag.value.trim()) { //有值
                if (tag.value.trim() !== tag.tmp.trim()) { //有修改
                    console.log(tag)
                    //修改
                    EditGroup(tag.value.trim(),tag.tmp.trim()).then((result) => {
                        this.loadGroups()
                        tag.tmp = tag.value.trim()
                    }).catch(error => {
                        // 处理错误
                        ElMessage.error(error)
                        return false
                    });
                }
            } else {
                tag.value = tag.tmp.trim() //当为空时不修改
            }
        },
        newTagHandleInputConfirm() {//添加分组tag
            const newGorup = this.newTagInputValue.trim();
            if (newGorup) {
                AddGroup(newGorup).then((result) => {
                    this.dynamicTags.push({ value: newGorup, editable: false,tmp:newGorup,popover:false});
                    this.loadGroups()
                }).catch(error => {
                    // 处理错误
                    ElMessage.error(error)
                    return false
                });
            }
            this.newTagInputValue = '';
            this.newTagInputVisible = false;
        }
    },
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