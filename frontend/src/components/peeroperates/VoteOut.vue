<template>
    <el-form
        ref="ruleFormRef"
        :model="ruleForm"
        :rules="rules"
        label-width="120px"
        class="demo-ruleForm"
        status-icon
    >
        <el-form-item label="类型" prop="votetype">
            <el-select v-model="ruleForm.votetype" placeholder="投票类型">
                <el-option label="见证人投票" value="1" />
                <el-option label="社区节点投票" value="2" />
                <el-option label="轻节点押金" value="3" />
            </el-select>
        </el-form-item>
        <el-form-item label="地址" prop="address">
            <el-input v-model="ruleForm.address" type="text" />
        </el-form-item>
        <el-form-item label="见证人地址" prop="witness">
            <el-input v-model="ruleForm.witness" type="text" />
        </el-form-item>
        <el-form-item label="金额" prop="amount">
            <el-input v-model="ruleForm.amount" type="number" />
        </el-form-item>
        <el-form-item label="手续费" prop="gas">
            <el-input v-model="ruleForm.gas" type="number"/>
        </el-form-item>
        <el-form-item label="奖励分配比例" prop="rate">
            <el-slider v-model="ruleForm.rate" show-input />
        </el-form-item>
        <el-form-item label="Password" prop="pass">
            <el-input v-model="ruleForm.pass" type="password" show-password/>
        </el-form-item>

        <el-form-item>
            <el-button type="primary" @click="submitForm(ruleFormRef)">
                确认
            </el-button>
            <el-button @click="resetForm(ruleFormRef)">重置</el-button>
        </el-form-item>
    </el-form>
</template>


<script lang="ts" setup>
import { reactive, ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import {useStore} from "vuex";
import {ElMessage} from "element-plus";
import {VoteOut} from "../../../wailsjs/go/peer/Api";
const store = useStore()

interface RuleForm {
    votetype: number
    address:string
    witness:string
    pass: string
    amount: number
    gas: number
    rate:number
}

const ruleFormRef = ref<FormInstance>()
const ruleForm = reactive<RuleForm>({
    votetype:0,
    address:'',
    witness:'',
    pass: '',
    amount:100000,
    gas:0.1,
    rate:0,
})

const rules = reactive<FormRules<RuleForm>>({
    votetype: [
        {
            required: true,
            message: '类型不能为空',
            trigger: 'change',
        },
        { min: 1, max: 3, message: '类型不正确', trigger: 'change' },
    ],
})

const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate((valid, fields) => {
        VoteOut(store.getters.getPeerInfo.id,parseInt(ruleForm.votetype.toString()),ruleForm.address,parseFloat(ruleForm.amount.toString()),parseFloat(ruleForm.gas.toString()),parseFloat(ruleForm.rate.toString()),ruleForm.pass,'').then((result) => {
            ElMessage.success("成功！")
        }).catch(error => {
            // 处理错误
            ElMessage.error(error)
        });
    })
}

const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
}
</script>
