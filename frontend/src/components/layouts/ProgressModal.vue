<template>
    <el-dialog v-model="dialogVisible" :title="title">
        <el-progress :percentage="percentage" :text-inside="true" :stroke-width="26" status="success" />
    </el-dialog>
</template>

<script setup lang="ts">
import { ref, defineProps } from 'vue';
import {  ElProgress, ElDialog } from 'element-plus';

// Define props for customizing the modal.
const props = defineProps({
    title: {
        type: String,
        default: 'Tips'
    },
});

const dialogVisible = ref(false);
const percentage = ref(0);
const title = ref(props.title);

// Opens the download progress modal.
const openModal = () => {
    dialogVisible.value = true;
};
const closeModel = () => {
    dialogVisible.value = false;
};
// Updates the progress of the download.
const updateProgress = (newPercentage: number) => {
    percentage.value = newPercentage;
};

defineExpose({ openModal, closeModel,updateProgress });
</script>

<style scoped>
.el-dialog__body {
    padding: 20px;
}

.el-progress {
    margin-bottom: 10px;
}
</style>