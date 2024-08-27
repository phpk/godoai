<script setup lang="ts">
import { useClipboard } from "@vueuse/core";
const { copy } = useClipboard();
import { notifySuccess } from "@/util/msg";
const copiedText = "Copied to clipboard!";

// Props
const props = defineProps({
  // Text to copy to clipboard
  text: {
    type: String,
    default: "",
  },
});

// Copy Text
const copyText = () => {
  console.log(props.text);
  copy(props.text);
  notifySuccess(copiedText);
};
</script>

<template>
  <el-tooltip
    class="box-item"
    effect="dark"
    content="Copy"
    placement="top-start"
  >
    <el-button 
    icon="CopyDocument" @click="copyText" size="small" circle />
  </el-tooltip>
</template>

