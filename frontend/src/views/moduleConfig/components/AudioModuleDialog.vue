<template>
    <n-modal 
      :show="show"
      @update:show="value => emits('update:show', value)"
      preset="dialog" 
      :title="$t('音频交互模块')"
      style="width: 800px"
      :mask-closable="false"
    >
      <div class="audio-module-form">
        <!-- 触发条件 -->
        <div class="form-row">
          <div class="form-item">
            <span class="label">{{ $t('触发条件') }} *</span>
            <div class="trigger-conditions">
              <n-checkbox v-model:checked="audioForm.triggers.cycleExecution">{{ $t('新兵执行') }}</n-checkbox>
              <n-checkbox v-model:checked="audioForm.triggers.elasticComment">{{ $t('弹幕野怪') }}</n-checkbox>
              <n-checkbox v-model:checked="audioForm.triggers.gift">{{ $t('送礼物') }}</n-checkbox>
              <n-checkbox v-model:checked="audioForm.triggers.like">{{ $t('点赞') }}</n-checkbox>
              <n-checkbox v-model:checked="audioForm.triggers.joinRoom">{{ $t('进入直播间') }}</n-checkbox>
              <n-checkbox v-model:checked="audioForm.triggers.notice">{{ $t('管理提示') }}</n-checkbox>
            </div>
          </div>
        </div>
  
        <!-- 时间间隔 -->
        <div class="form-row">
          <div class="form-item">
            <span class="label">{{ $t('时间间隔') }} *</span>
            <div class="time-range">
              <span>{{ $t('抖音值') }}:</span>
              <n-input-number v-model:value="audioForm.minTime" :min="0" style="width: 100px" />
              <span>{{ $t('秒') }}</span>
              <span style="margin-left: 20px;">{{ $t('抖音值上限') }}:</span>
              <n-input-number v-model:value="audioForm.maxTime" :min="0" style="width: 100px" />
              <span>{{ $t('秒') }}</span>
              <n-button type="primary" style="margin-left: 20px;" @click="executeOnce">{{ $t('执行一次') }}</n-button>
            </div>
          </div>
        </div>
  
        <!-- 选择上传音频文件 -->
        <div class="form-row">
          <div class="form-item">
            <span class="label">{{ $t('选择上传音频文件') }} *</span>
            <n-upload
              v-model:file-list="audioForm.audioFiles"
              accept="audio/*"
              :max="1"
              list-type="text"
              @change="handleAudioChange"
            >
              <n-upload-dragger>
                <div style="padding: 20px; text-align: center;">
                  <p>{{ $t('点击或拖拽文件到此区域上传') }}</p>
                </div>
              </n-upload-dragger>
            </n-upload>
          </div>
        </div>
      </div>
      
      <!-- 底部按钮 -->
      <template #footer>
        <div style="text-align: right">
          <n-button @click="handleClose">{{ $t('取消') }}</n-button>
          <n-button type="primary" @click="handleSubmit" style="margin-left: 12px">{{ $t('确定') }}</n-button>
        </div>
      </template>
    </n-modal>
  </template>
  
  <script setup lang="ts">
  import { ref, watch } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { NModal, NInputNumber, NCheckbox, NButton, NUpload, NUploadDragger } from 'naive-ui';
  import { eventBUS } from '@/views/Home/utils/tools';
  
  const { t: $t } = useI18n();
  
  const props = defineProps<{
    show: boolean;
    initialData: any;
  }>();
  
  const emits = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'submit', value: any): void;
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
  
  // 执行一次（模拟功能）
  const executeOnce = () => {
    console.log('执行一次:', audioForm.value);
    eventBUS.$emit('executeAudioOnce', audioForm.value);
  };
  
  // 关闭弹窗
  const handleClose = () => {
    emits('update:show', false);
  };
  
  // 提交表单
  const handleSubmit = () => {
    emits('submit', audioForm.value);
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
  </style>