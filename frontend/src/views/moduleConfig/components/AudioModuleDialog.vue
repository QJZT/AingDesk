<template>
  <n-modal
    :show="show"
    @update:show="value => emits('update:show', value)"
    preset="dialog"
    title="音频文件模块"
    style="width: 800px"
    :mask-closable="false"
  >
    <div class="audio-module-form">
      <!-- 触发条件 -->
      <div class="form-row">
        <div class="form-item">
          <span class="label">触发条件 *</span>
          <div class="trigger-conditions">
            <n-radio-group v-model:value="audioForm.selectedTrigger">
              <n-radio value="SceneLoop">控场循环</n-radio>
              <n-radio value="IntervalLoop">间隔循环</n-radio>
              <n-radio value="BarrageComment">弹幕评论</n-radio>
              <n-radio value="SendGift">送礼物</n-radio>
              <n-radio value="Like">点赞</n-radio>
              <n-radio value="EnterLiveRoom">进入直播间</n-radio>
              <n-radio value="ShareRoom">分享直播间</n-radio>
              <n-radio value="FollowRoom">关注直播间</n-radio>
            </n-radio-group>
          </div>
        </div>
      </div>

      <!-- 时间间隔 -->
      <div class="form-row" v-if="audioForm.selectedTrigger !== 'SceneLoop'">
        <div class="form-item">
          <span class="label">时间间隔（秒）*</span>
          <div class="time-range">
            <n-input-number
              v-model:value="audioForm.minTime"
              :min="0"
              style="width: 100px"
              show-button
              placeholder="最小值"
            />
            <span style="margin: 0 10px;">到</span>
            <n-input-number
              v-model:value="audioForm.maxTime"
              :min="0"
              style="width: 100px"
              show-button
              placeholder="最大值"
            />
            <span style="margin-left: 20px;">执行一次</span>
          </div>
        </div>
      </div>

      <!-- 选择上传音频文件 -->
      <div class="form-row">
        <div class="form-item">
          <span class="label">选择上传音频文件 *</span>
          <n-upload
            :file-list="audioForm.audioFiles"
            accept="audio/*"
            :max="1"
            list-type="text"
            @change="handleAudioChange"
            :disabled="uploading"
          >
            <n-upload-dragger>
              <div style="padding: 20px; text-align: center;">
                <p v-if="uploading">上传中，请稍候...</p>
                <p v-else-if="audioForm.audioName">已上传：{{ audioForm.audioName }}</p>
                <p v-else>点击或拖拽文件到此区域上传</p>
              </div>
            </n-upload-dragger>
          </n-upload>
        </div>
      </div>
    </div>

    <!-- 底部按钮 -->
    <template #action>
      <div style="text-align: right">
        <n-button type="warning" @click="handleExit">退出</n-button>
        <n-button type="error" @click="handleSave" style="margin-left: 12px" :disabled="uploading || !audioForm.selectedTrigger || !audioForm.audioPath">保存</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { NModal, NInputNumber, NButton, NUpload, NUploadDragger, NRadioGroup, NRadio } from 'naive-ui';
import type { Ref } from 'vue';

// 后端 API 基础 URL
const API_BASE_URL = 'http://localhost:7072'; // 确保与后端端口一致

const props = defineProps<{
  show: boolean;
  initialData: any;
}>();

const emits = defineEmits<{
  (e: 'update:show', value: boolean): void;
  (e: 'exit', value: any): void;
  (e: 'save', value: any): void;
}>();

// 加载草稿
const loadDraft = () => {
  const draft = localStorage.getItem('audioModuleDraft');
  return draft
    ? JSON.parse(draft)
    : {
        minTime: 30,
        maxTime: 50,
        selectedTrigger: 'SceneLoop',
        audioFiles: [],
        audioName: '',
        audioPath: '',
      };
};

// 表单数据
const audioForm = ref(loadDraft());

// 上传状态
const uploading = ref(false);

// 监听弹窗显示状态
watch(() => props.show, (newShow) => {
  if (newShow) {
    if (props.initialData) {
      // 编辑模式：使用初始数据
      audioForm.value = {
        minTime: props.initialData.intervalTimeStart || 30,
        maxTime: props.initialData.intervalTimeEnd || 50,
        selectedTrigger: props.initialData.triggerConditions?.[0] || 'SceneLoop',
        audioFiles: [],
        audioName: props.initialData.audioName || '',
        audioPath: props.initialData.audioPath || '',
      };
    } else {
      // 新建模式：重置表单
      audioForm.value = {
        minTime: 30,
        maxTime: 50,
        selectedTrigger: 'SceneLoop',
        audioFiles: [],
        audioName: '',
        audioPath: '',
      };
    }
    // 清除草稿
    localStorage.removeItem('audioModuleDraft');
  }
}, { immediate: true });

// 移除原有的 props.initialData 监听器
// watch(() => props.initialData, ...) 删除这段代码

// 初始化或编辑模式
watch(() => props.initialData, (newData) => {
  if (newData && !localStorage.getItem('audioModuleDraft')) {
    audioForm.value = {
      minTime: newData.intervalTimeStart || 30,
      maxTime: newData.intervalTimeEnd || 50,
      selectedTrigger: newData.triggerConditions?.[0] || 'SceneLoop',
      audioFiles: [],
      audioName: newData.audioName || '',
      audioPath: newData.audioPath || '',
    };
  }
  console.log('AudioModuleDialog 初始化表单数据:', audioForm.value);
}, { immediate: true });

// 监听表单变化，保存草稿
watch(audioForm, (newFormData) => {
  const draftData = {
    minTime: newFormData.minTime,
    maxTime: newFormData.maxTime,
    selectedTrigger: newFormData.selectedTrigger,
    audioName: newFormData.audioName,
    audioPath: newFormData.audioPath,
  };
  localStorage.setItem('audioModuleDraft', JSON.stringify(draftData));
  console.log('保存草稿到 localStorage (AudioModule):', draftData);
}, { deep: true });

// 处理音频文件上传
const handleAudioChange = async (options: { fileList: any[] }) => {
  const file = options.fileList[0]?.file;
  if (!file) {
    audioForm.value.audioFiles = [];
    audioForm.value.audioName = '';
    audioForm.value.audioPath = '';
    return;
  }

  uploading.value = true;
  const formData = new FormData();
  formData.append('file', file);

  try {
    const response = await fetch(`${API_BASE_URL}/upload-audio`, {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || '上传失败');
    }

    const result = await response.json();
    console.log('文件上传成功，返回数据:', result);

    if (!result.path) {
      throw new Error('上传成功但未返回有效的音频路径');
    }

    audioForm.value.audioName = file.name;
    audioForm.value.audioPath = result.path;
    audioForm.value.audioFiles = options.fileList;
    window.$message?.success('文件上传成功');
  } catch (error) {
    console.error('文件上传失败:', error);
    window.$message?.error('文件上传失败：' + error.message);
    audioForm.value.audioFiles = [];
    audioForm.value.audioName = '';
    audioForm.value.audioPath = '';
  } finally {
    uploading.value = false;
  }
};

// 退出
const handleExit = () => {
  console.log('退出音频模块弹窗，当前表单数据:', audioForm.value);
  // 清除草稿
  localStorage.removeItem('audioModuleDraft');
  emits('exit', audioForm.value);
  emits('update:show', false);
};

// 保存
enum TriggerCondition {
  SceneLoop = "SceneLoop",
  IntervalLoop = "IntervalLoop",
  BarrageComment = "BarrageComment",
  SendGift = "SendGift",
  Like = "Like",
  EnterLiveRoom = "EnterLiveRoom",
  ShareRoom = "ShareRoom",
  FollowRoom = "FollowRoom",
}

const handleSave = () => {
  if (!audioForm.value.selectedTrigger) {
    window.$message?.error('保存失败：请先选择触发条件');
    return;
  }

  if (audioForm.value.selectedTrigger !== 'SceneLoop') {
    if (audioForm.value.minTime === undefined || audioForm.value.maxTime === undefined) {
      window.$message?.error('保存失败：时间间隔不能为空');
      return;
    }
    if (audioForm.value.minTime >= audioForm.value.maxTime) {
      window.$message?.error('保存失败：最小时间必须小于最大时间');
      return;
    }
  }

  if (!audioForm.value.audioPath || !audioForm.value.audioPath.trim()) {
    window.$message?.error('保存失败：请先上传有效的音频文件');
    return;
  }

  console.log('保存前 audioPath:', audioForm.value.audioPath);

  let intervalTimeStart = audioForm.value.minTime;
  let intervalTimeEnd = audioForm.value.maxTime;
  if (audioForm.value.selectedTrigger === 'SceneLoop') {
    intervalTimeStart = 0;
    intervalTimeEnd = 0;
  }

  let triggerCondition;
  switch (audioForm.value.selectedTrigger) {
    case 'SceneLoop':
      triggerCondition = TriggerCondition.SceneLoop;
      break;
    case 'IntervalLoop':
      triggerCondition = TriggerCondition.IntervalLoop;
      break;
    case 'BarrageComment':
      triggerCondition = TriggerCondition.BarrageComment;
      break;
    case 'SendGift':
      triggerCondition = TriggerCondition.SendGift;
      break;
    case 'Like':
      triggerCondition = TriggerCondition.Like;
      break;
    case 'EnterLiveRoom':
      triggerCondition = TriggerCondition.EnterLiveRoom;
      break;
    case 'ShareRoom':
      triggerCondition = TriggerCondition.ShareRoom;
      break;
    case 'FollowRoom':
      triggerCondition = TriggerCondition.FollowRoom;
      break;
    default:
      triggerCondition = TriggerCondition.SceneLoop;
  }

  const formData = {
    id: props.initialData?.id || 0,
    moduleType: 'audio',
    moduleName: props.initialData?.moduleName || '音频模块',  // 保留原有模块名称
    intervalTimeStart,
    intervalTimeEnd,
    triggerConditions: [triggerCondition],
    readStep: props.initialData?.readStep || '',
    scriptContent: props.initialData?.scriptContent || [],
    isModelRewrite: props.initialData?.isModelRewrite || false,
    rewriteFrequency: props.initialData?.rewriteFrequency || 0,
    audioName: audioForm.value.audioName,
    audioPath: audioForm.value.audioPath,
  };

  console.log('保存音频模块，当前表单数据:', formData);
  emits('save', formData);
  emits('update:show', false);
  
  // 清除草稿和重置表单
  localStorage.removeItem('audioModuleDraft');
  audioForm.value = loadDraft();
};
</script>

<style scoped lang="scss">
@use "@/assets/base";

.audio-module-form {
  padding: 20px;
  max-height: 70vh;
  overflow-y: auto;

  .form-row {
    margin-bottom: 20px;
  }

  .form-item {
    .label {
      display: block;
      margin-bottom: 8px;
      font-weight: bold;
    }

    .time-range {
      display: flex;
      align-items: center;
      gap: 10px;
    }

    .trigger-conditions {
      display: flex;
      gap: 20px;
      flex-wrap: wrap;
    }
  }

  .n-upload-dragger {
    border: 2px dashed #d9d9d9;
    border-radius: 4px;
    padding: 20px;
    text-align: center;
    cursor: pointer;
    transition: border-color 0.3s;

    &:hover {
      border-color: #409eff;
    }
  }
}
</style>