<script lang="ts" setup>
import { useRecordingStore } from "@/stores/recording";
import { watchEffect, ref, toRaw } from "vue";
import { t } from "@/i18n/index";
import { notifyError } from "@/util/msg";
const recordingStore = useRecordingStore();
const formData = ref({
  name: recordingStore.title,
  time: new Date(),
  local: "",
  desc: "",
});
watchEffect(async () => {
  if (recordingStore.currentId > 0) {
    const data = await recordingStore.getRecording(recordingStore.currentId);
    formData.value = {
      name: data.name,
      time: data.time,
      local: data.local,
      desc: data.desc,
    };
  } else {
    formData.value = {
      name: recordingStore.title,
      time: new Date(),
      local: "",
      desc: "",
    };
  }
});
const emit = defineEmits(["saveFn"]);
async function save() {
  const saveData = toRaw(formData.value);
  if (saveData.name == "") {
    notifyError(t('recording.inputName'));
    return;
  }
  emit("saveFn", saveData);
}
</script>
<template>
  <el-form
    ref="form"
    :model="formData"
    label-width="100px"
    style="margin-top: 15px"
  >
  <el-form-item :label="t('recording.name')">
      <el-input
        v-model="formData.name"
        prefix-icon="House"
        clearable
        :placeholder="t('recording.meetingName')"
      ></el-input>
    </el-form-item>
    <el-form-item :label="t('recording.location')">
      <el-input
        v-model="formData.local"
        prefix-icon="HomeFilled"
        clearable
        :placeholder="t('recording.meetingLocation')"
      ></el-input>
    </el-form-item>
    <el-form-item :label="t('recording.time')">
      <el-date-picker
          v-model="formData.time"
          type="datetime"
          :placeholder="t('recording.meetingTime')"
        />
    </el-form-item>
    <el-form-item :label="t('recording.remarks')">
      <el-input
        type="textarea"
        v-model="formData.desc"
        rows="5"
        clearable
        :placeholder="t('recording.meetingRemarks')"
      ></el-input>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="save">{{
        t("common.confim")
      }}</el-button>
    </el-form-item>
  </el-form>
  
</template>

<style lang="scss">
.dp__input {
  height: 57px;
}
</style>
