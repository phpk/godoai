<script lang="ts" setup>
import { useCreationStore } from "@/stores/creation";
const creationStore = useCreationStore();
const emit = defineEmits(["setBaseFn", "setTargetFn", "newFn"]);
const changeRoom = async (id:any) => {
  const data: any = await creationStore.getActicle(id);
  if (!data) return;
  creationStore.activeId = id;
  creationStore.baseTitle = data.title;
  creationStore.currentCate = data.cate;
  emit("setBaseFn", data.base);
  emit("setTargetFn", data.content);
  creationStore.showLeft = false;
};
const deleteRoom = async (id:any) => {
  await creationStore.delArticle(id);
  if (creationStore.activeId < 1 || creationStore.activeId == id) {
    //newCreate();
    emit("newFn", id);
  } else {
    changeRoom(creationStore.activeId);
  }
  creationStore.showLeft = false;
};
</script>
<template>
  <div
    class="list-item"
    v-for="(item, index) in creationStore.articleList"
    :key="index"
  >
    <div class="list-title" @click="changeRoom(item.id)">
      {{ item.title }}
    </div>
    <el-button
      class="delete-btn"
      icon="Delete"
      size="small"
      @click.stop="deleteRoom(item.id)"
      circle
    ></el-button>
  </div>

  <el-row
    justify="center"
    style="margin-top: 15px"
    v-if="creationStore.page.total > creationStore.page.size"
  >
    <el-pagination
      background
      layout="prev, pager, next"
      v-model:current-page="creationStore.page.current"
      v-model:page-size="creationStore.page.size"
      :total="creationStore.page.total"
      @current-change="(val:any) => creationStore.pageClick(val)"
    />
  </el-row>
</template>

<style scoped>
@import "@/assets/left.scss";
</style>