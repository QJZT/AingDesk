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
                      <n-checkbox v-model:checked="audioForm.triggers.cycleExecution">循环执行</n-checkbox>
                      <n-checkbox v-model:checked="audioForm.triggers.elasticComment">弹幕评论</n-checkbox>
                      <n-checkbox v-model:checked="audioForm.triggers.gift">送礼物</n-checkbox>
                      <n-checkbox v-model:checked="audioForm.triggers.like">点赞</n-checkbox>
                      <n-checkbox v-model:checked="audioForm.triggers.joinRoom">进入直播间</n-checkbox>
                      <n-checkbox v-model:checked="audioForm.triggers.notice">管理警告</n-checkbox>
                  </div>
              </div>
          </div>

          <!-- 时间间隔 -->
          <div class="form-row">
              <div class="form-item">
                  <span class="label">时间间隔 *</span>
                  <div class="time-range">
                      <n-input-number 
                          v-model:value="audioForm.minTime" 
                          :min="0" 
                          style="width: 100px" 
                          show-button
                      />
                      <span style="margin: 0 10px;">到</span>
                      <n-input-number 
                          v-model:value="audioForm.maxTime" 
                          :min="0" 
                          style="width: 100px" 
                          show-button
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
                      v-model:file-list="audioForm.audioFiles"
                      accept="audio/*"
                      :max="1"
                      list-type="text"
                      @change="handleAudioChange"
                  >
                      <n-upload-dragger>
                          <div style="padding: 20px; text-align: center;">
                              <p>点击或拖拽文件到此区域上传</p>
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
              <n-button type="error" @click="handleSave" style="margin-left: 12px">保存</n-button>
          </div>
      </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { NModal, NInputNumber, NCheckbox, NButton, NUpload, NUploadDragger } from 'naive-ui';

const props = defineProps<{
  show: boolean;
  initialData: any;
}>();

const emits = defineEmits<{
  (e: 'update:show', value: boolean): void;
  (e: 'exit', value: any): void;
  (e: 'save', value: any): void;
}>();

// 表单数据
const audioForm = ref({
  minTime: 30,
  maxTime: 50,
  triggers: {
      cycleExecution: false,
      elasticComment: true,
      gift: false,
      like: false,
      joinRoom: false,
      notice: false
  },
  audioFiles: [] as File[]
});

// 监听 initialData 变化，填充表单（编辑模式）
watch(() => props.initialData, (newData) => {
  if (newData) {
      audioForm.value = {
          minTime: newData.minTime || 30,
          maxTime: newData.maxTime || 50,
          triggers: {
              cycleExecution: newData.triggers?.cycleExecution || false,
              elasticComment: newData.triggers?.elasticComment || false,
              gift: newData.triggers?.gift || false,
              like: newData.triggers?.like || false,
              joinRoom: newData.triggers?.joinRoom || false,
              notice: newData.triggers?.notice || false
          },
          audioFiles: newData.audioFiles || []
      };
  }
}, { immediate: true });

// 处理音频文件上传
const handleAudioChange = (options: { fileList: File[] }) => {
  audioForm.value.audioFiles = options.fileList;
};

// 退出并输出提示
const handleExit = () => {
  console.log('退出音频模块弹窗，当前表单数据:', audioForm.value);
  emits('exit', audioForm.value);
  emits('update:show', false);
};

// 保存并输出提示
const handleSave = () => {
  console.log('保存音频模块，当前表单数据:', audioForm.value);
  emits('save', audioForm.value);
  emits('update:show', false);
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
</style>