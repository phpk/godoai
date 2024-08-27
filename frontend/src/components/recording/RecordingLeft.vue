<script lang="ts" setup>
import {useRecordingStore} from "@/stores/recording.ts"
const recordingStore = useRecordingStore();
</script>
<template>
  <div
    class="list-item"
    v-for="(item, index) in recordingStore.recordingList"
    :key="index"
  >
    <div class="list-title" @click="recordingStore.changeRecording(item.id)">
      {{ item.name }}
    </div>
    <el-button
      class="delete-btn"
      icon="Delete"
      size="small"
      @click.stop="recordingStore.deleteRecording(item.id)"
      circle
    ></el-button>
  </div>

  <el-row
    justify="center"
    style="margin-top: 15px"
    v-if="recordingStore.page.total > recordingStore.page.size"
  >
    <el-pagination
      background
      layout="prev, pager, next"
      v-model:current-page="recordingStore.page.current"
      v-model:page-size="recordingStore.page.size"
      :total="recordingStore.page.total"
      @current-change="(val:any) => recordingStore.pageClick(val)"
    />
  </el-row>
</template>

<style scoped>
@import "@/assets/left.scss";
</style>