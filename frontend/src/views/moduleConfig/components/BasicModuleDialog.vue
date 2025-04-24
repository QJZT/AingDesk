<template>
    <n-modal 
      :show="show"
      @update:show="value => emits('update:show', value)"
      preset="dialog" 
      :title="$t('基础模块')"
      style="width: 800px"
      :mask-closable="false"
    >
      <div class="basic-module-form">
        <!-- 时间间隔 -->
        <div class="form-row">
          <div class="form-item">
            <span class="label">{{ $t('时间间隔') }} *</span>
            <div class="time-range">
              <span>{{ $t('抖音值') }}:</span>
              <n-input-number v-model:value="basicForm.minTime" :min="0" style="width: 100px" />
              <span>{{ $t('秒') }}</span>
              <span style="margin-left: 20px;">{{ $t('到达') }}:</span>
              <n-input-number v-model:value="basicForm.maxTime" :min="0" style="width: 100px" />
            </div>
          </div>
        </div>
  
        <!-- 触发条件 -->
        <div class="form-row">
          <div class="form-item">
            <span class="label">{{ $t('触发条件') }} *</span>
            <div class="trigger-conditions">
              <n-checkbox v-model:checked="basicForm.triggers.cycleExecution">{{ $t('催件执行') }}</n-checkbox>
              <n-checkbox v-model:checked="basicForm.triggers.elasticComment">{{ $t('弹幕野怪') }}</n-checkbox>
              <n-checkbox v-model:checked="basicForm.triggers.gift">{{ $t('送礼物') }}</n-checkbox>
              <n-checkbox v-model:checked="basicForm.triggers.like">{{ $t('点赞') }}</n-checkbox>
              <n-checkbox v-model:checked="basicForm.triggers.joinRoom">{{ $t('进入直播间') }}</n-checkbox>
              <n-checkbox v-model:checked="basicForm.triggers.notice">{{ $t('管理提示') }}</n-checkbox>
            </div>
          </div>
        </div>
  
        <!-- 深度逻辑 -->
        <div class="form-row">
          <div class="form-item">
            <span class="label">{{ $t('深度逻辑') }} *</span>
            <n-radio-group v-model:value="basicForm.readStep">
              <n-radio value="random">{{ $t('是') }}</n-radio>
              <n-radio value="none">{{ $t('否') }}</n-radio>
            </n-radio-group>
          </div>
        </div>
  
        <!-- 大模型或者 -->
        <div class="form-row">
          <div class="form-item">
            <span class="label">{{ $t('大模型或者') }} *</span>
            <n-radio-group v-model:value="basicForm.modelRewrite">
              <n-radio value="yes">{{ $t('是') }}</n-radio>
              <n-radio value="no">{{ $t('否') }}</n-radio>
            </n-radio-group>
          </div>
        </div>
  
        <!-- 改写概率 -->
        <div class="form-row">
          <div class="form-item">
            <span class="label">{{ $t('改写概率') }} *</span>
            <div>
              <span>{{ $t('抖音值概率') }}:</span>
              <n-input-number v-model:value="basicForm.rewriteFrequency" :min="0" style="width: 100px" />
            </div>
          </div>
        </div>
  
        <!-- 文案 -->
        <div class="form-row">
          <div class="form-item">
            <span class="label">{{ $t('文案') }} *</span>
            <!-- 第一组标签 -->
            <div class="speech-tags">
              <n-tag
                v-for="tag in speechTags"
                :key="tag"
                class="speech-tag"
                @click="addSpeechTag(tag)"
              >
                {{ $t(tag) }}
              </n-tag>
            </div>
            <!-- 第二组标签 -->
            <div class="speech-tags">
              <n-tag
                v-for="tag in additionalTags"
                :key="tag"
                class="speech-tag"
                @click="addSpeechTag(tag)"
              >
                {{ $t(tag) }}
              </n-tag>
            </div>
            <!-- 文案输入框 -->
            <n-input
              type="textarea"
              v-model:value="basicForm.speechContent"
              :placeholder="$t('请键入指令')"
              :autosize="{ minRows: 3, maxRows: 5 }"
            />
            <div class="word-count">{{ basicForm.speechContent.length }}/20000</div>
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
  import { NModal, NInput, NInputNumber, NCheckbox, NRadio, NRadioGroup, NTag, NButton } from 'naive-ui';
  
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
  const basicForm = ref({
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
    readStep: 'random',
    modelRewrite: 'yes',
    rewriteFrequency: 180,
    speechContent: ''
  });
  
  // 监听 initialData 变化，填充表单（编辑模式）
  watch(() => props.initialData, (newData) => {
    if (newData) {
      basicForm.value = {
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
        readStep: newData.readStep || 'random',
        modelRewrite: newData.modelRewrite || 'yes',
        rewriteFrequency: newData.rewriteFrequency || 180,
        speechContent: newData.description || ''
      };
    }
  }, { immediate: true });
  
  // 第一组标签
  const speechTags = ref([
    '现在时间',
    '现在日期',
    '现在星期',
    '下次结果',
    '聊天去分析',
    '聊天去分析',
    '在线人数',
    '直播间名称'
  ]);
  
  // 第二组标签
  const additionalTags = ref([
    '场景通用名',
    '选择通用用户名',
    '北极星名',
    '进入直播间用户名',
    '美導購通用名',
    '3种高频语',
    '直播間頻率高語'
  ]);
  
  // 点击标签插入到文案
  const addSpeechTag = (tag: string) => {
    basicForm.value.speechContent += ` ${$t(tag)}`;
  };
  
  // 关闭弹窗
  const handleClose = () => {
    emits('update:show', false);
  };
  
  // 提交表单
  const handleSubmit = () => {
    emits('submit', basicForm.value);
    emits('update:show', false);
  };
  </script>
  
  <style scoped lang="scss">
  @use "@/assets/base";
  
  .basic-module-form {
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
  
      .speech-tags {
        margin-bottom: 10px;
        display: flex;
        gap: 10px;
        flex-wrap: wrap;
  
        .speech-tag {
          cursor: pointer;
          border-radius: 15px;
          background-color: #f0f0f5;
          &:hover {
            background-color: #e0e0e5;
          }
        }
      }
  
      .word-count {
        text-align: right;
        color: #999;
        margin-top: 5px;
      }
    }
  }
  </style>