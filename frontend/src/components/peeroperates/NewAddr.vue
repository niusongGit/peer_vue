<template>
    <el-form
            ref="ruleFormRef"
            :model="ruleForm"
            :rules="rules"
            label-width="120px"
            class="demo-ruleForm"
            status-icon
    >
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
import {GetNewAddress} from '../../../wailsjs/go/peer/Api'
import {ElMessage} from "element-plus";
import {useStore} from "vuex";
const store = useStore()
interface RuleForm {
    pass: string

}

const ruleFormRef = ref<FormInstance>()
const ruleForm = reactive<RuleForm>({
    pass: '',
})

const rules = reactive<FormRules<RuleForm>>({
    pass: [
        {
            required: true,
            message: '密码不能为空',
            trigger: 'change',
        },
        { min: 6, max: 14, message: '密码长度不正确', trigger: 'blur' },
    ],
})

const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate((valid, fields) => {
        if (valid) {

            GetNewAddress(store.getters.getPeerInfo.id,ruleForm.pass).then((result) => {
                ElMessage.success("成功！")
            }).catch(error => {
                // 处理错误
                ElMessage.error(error)
            });

            console.log('submit!')
        } else {
            console.log('error submit!', fields)
        }
    })
}

const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
}
</script>
