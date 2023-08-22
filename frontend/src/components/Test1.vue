vue
<template>
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
                @keyup.enter.native="tagHandleInputConfirm(index)"
                @blur="tagHandleInputConfirm(index)"
            ></el-input>

            <el-popover :visible="tag.popover" placement="bottom" :width="160"  trigger="manual">
                <p>删除该分组和该分组下所有节点?</p>
                <div style="text-align: right; margin: 0">
                    <el-button size="small" text @click="tag.popover = false">否</el-button>
                    <el-button size="small" type="primary" @click="tagHandleClose(tag)"
                    >是</el-button
                    >
                </div>
                <template #reference><span></span></template>
            </el-popover>

        </el-tag>
        <el-input
            class="input-new-tag"
            v-if="newTagInputVisible"
            v-model="newTagInputValue"
            ref="saveTagInput"
            size="small"
            @keyup.enter.native="newTagHandleInputConfirm"
            @blur="newTagHandleInputConfirm"
        ></el-input>
        <el-button v-else class="button-new-tag" size="small" @click="showNewTagInput">+ New Tag</el-button>
    </div>
</template>

<script>
export default {
    data() {
        return {
            dynamicTags: [
                { value: 'Tag 1', editable: false,popover:false },
                { value: 'Tag 2', editable: false,popover:false },
                { value: 'Tag 3', editable: false,popover:false }
            ],
            newTagInputVisible: false,
            popoverRef: null,
            newTagInputValue: ''
        };
    },
    methods: {
        showPopover(tag) {
            tag.popover = true;
        },
        popoverShow() {
            this.$nextTick(() => {
                const popover = this.$refs.popover;
                if (popover) {
                    popover.$refs.reference.$el.querySelector('.el-tag__close').click();
                }
            });
        },
        tagHandleClose(tag) {
            const index = this.dynamicTags.findIndex((t) => t.value === tag.value);
            if (index !== -1) {
                this.dynamicTags.splice(index, 1);
            }
            tag.popover = false;
        },

        showNewTagInput() {
            this.newTagInputVisible = true;
            this.$nextTick(() => {
                this.$refs.saveTagInput.focus();
            });
        },

        tagHandleEdit(index) {
            this.dynamicTags[index].editable = true;
            console.log("tagHandleEdit:",this.dynamicTags[index])
            this.$nextTick(() => {
                const input = this.$refs.editInput[index];
                if (input) {
                    input.$el.querySelector('input').focus();
                }
            });
        },

        tagHandleInputConfirm(index) {
            const tag = this.dynamicTags[index];
            console.log("tagHandleInputConfirm:",tag)
            if (tag.value.trim()) {
                tag.editable = false;
            } else {
                this.dynamicTags.splice(index, 1);
            }
        },

        newTagHandleInputConfirm() {
            const newTagInputValue = this.newTagInputValue.trim();
            if (newTagInputValue) {
                this.dynamicTags.push({ value: newTagInputValue, editable: false });
            }
            this.newTagInputValue = '';
            this.newTagInputVisible = false;
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
上面代码中，el-popover弹出时显示的位置不在触发它显示的close按钮所在的下方，请修改代码实现