<!-- BasicModuleDialog.vue -->
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
            <n-input-number v-model:value="basicForm.minTime" :min="0" style="width: 100px" />
            <span>{{ $t('秒到') }}</span>
            <n-input-number v-model:value="basicForm.maxTime" :min="0" style="width: 100px" />
            <span>{{ $t('秒') }}</span>
          </div>
        </div>
      </div>

      <!-- 触发条件 -->
      <div class="form-row">
        <div class="form-item">
          <span class="label">{{ $t('触发条件') }} *</span>
          <div class="trigger-conditions">
            <n-checkbox v-model:checked="basicForm.triggers.cycleExecution">{{ $t('循环执行') }}</n-checkbox>
            <n-checkbox v-model:checked="basicForm.triggers.elasticComment">{{ $t('弹幕评论') }}</n-checkbox>
            <n-checkbox v-model:checked="basicForm.triggers.gift">{{ $t('送礼物') }}</n-checkbox>
            <n-checkbox v-model:checked="basicForm.triggers.like">{{ $t('点赞') }}</n-checkbox>
            <n-checkbox v-model:checked="basicForm.triggers.joinRoom">{{ $t('进入直播间') }}</n-checkbox>
            <n-checkbox v-model:checked="basicForm.triggers.notice">{{ $t('警告提示') }}</n-checkbox>
          </div>
        </div>
      </div>

      <!-- 深度逻辑 -->
      <div class="form-row">
        <div class="form-item">
          <span class="label">{{ $t('读取步骤') }} *</span>
          <n-radio-group v-model:value="basicForm.readStep">
            <n-radio value="random">{{ $t('随机') }}</n-radio>
            <n-radio value="sequential">{{ $t('顺序') }}</n-radio>
          </n-radio-group>
        </div>
      </div>

      <!-- 大模型改写 -->
      <div class="form-row">
        <div class="form-item">
          <span class="label">{{ $t('大模型改写') }} *</span>
          <n-radio-group v-model:value="basicForm.modelRewrite">
            <n-radio value="yes">{{ $t('是') }}</n-radio>
            <n-radio value="no">{{ $t('否') }}</n-radio>
          </n-radio-group>
        </div>
      </div>

      <!-- 改写频率 -->
      <div class="form-row">
        <div class="form-item">
          <span class="label">{{ $t('改写频率') }} *</span>
          <div>
            <n-input-number v-model:value="basicForm.rewriteFrequency" :min="0" style="width: 100px" />
          </div>
        </div>
      </div>

      <!-- 文案 -->
      <div class="form-row">
        <div class="form-item">
          <span class="label">{{ $t('话术文案') }} *</span>
          
          <!-- 合并后的标签 -->
          <div class="speech-tags">
            <span class="label">{{ $t('变量支持') }}：</span>
            <n-tag class="speech-tag" @click="addSpeechTag('现在时间')">{{ $t('现在时间') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('现在日期')">{{ $t('现在日期') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('现在星期')">{{ $t('现在星期') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('下个整点')">{{ $t('下个整点') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('再过五分钟')">{{ $t('再过五分钟') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('再过十分钟')">{{ $t('再过十分钟') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('在线人数')">{{ $t('在线人数') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('直播间名称')">{{ $t('直播间名称') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('点赞用户名')">{{ $t('点赞用户名') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('送礼物用户名')">{{ $t('送礼物用户名') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('礼物名称')">{{ $t('礼物名称') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('发弹幕用户名')">{{ $t('发弹幕用户名') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('弹幕内容')">{{ $t('弹幕内容') }}</n-tag>
            <n-tag class="speech-tag" @click="addSpeechTag('直播间警告内容')">{{ $t('直播间警告内容') }}</n-tag>
          </div>
          <!-- 文案输入框 -->
          <n-input
            type="textarea"
            v-model:value="basicForm.speechContent"
            :placeholder="$t('请键入指令')"
            :autosize="{ minRows: 3, maxRows: 5 }"
            ref="speechInput"
          />
          <div class="word-count">{{ basicForm.speechContent.length }}/20000</div>
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
import { useI18n } from 'vue-i18n';
import { NModal, NInput, NInputNumber, NCheckbox, NRadio, NRadioGroup, NTag, NButton } from 'naive-ui';
import type { InputInst } from 'naive-ui';

const { t: $t } = useI18n();

const props = defineProps<{
  show: boolean;
  initialData: any;
}>();

const emits = defineEmits<{
  (e: 'update:show', value: boolean): void;
  (e: 'exit', value: any): void;
  (e: 'save', value: any): void;
}>();

// 从 localStorage 加载草稿
const loadDraft = () => {
  const draft = localStorage.getItem('basicModuleDraft');
  return draft
    ? JSON.parse(draft)
    : {
        minTime: 0,
        maxTime: 0,
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
        rewriteFrequency: 0,
        speechContent: ''
      };
};

// 表单数据，优先从 localStorage 加载草稿
const basicForm = ref(loadDraft());

// 监听 initialData 变化，填充表单（编辑模式）
watch(() => props.initialData, (newData) => {
  // 如果没有草稿（localStorage 中没有数据），则使用 initialData 初始化
  if (newData && !localStorage.getItem('basicModuleDraft')) {
    basicForm.value = {
      minTime: newData.minTime || 0,
      maxTime: newData.maxTime || 0,
      triggers: {
        cycleExecution: newData.triggers?.cycleExecution || false,
        elasticComment: newData.triggers?.elasticComment || false,
        gift: newData.triggers?.gift || false,
        like: newData.triggers?.like || false,
        joinRoom: newData.triggers?.joinRoom || false,
        notice: newData.triggers?.notice || false
      },
      readStep: newData.readStep || 'random',
      modelRewrite: newData.modelRewrite ? 'yes' : 'no',
      rewriteFrequency: newData.rewriteFrequency || 0,
      speechContent: newData.description || ''
    };
  }
  console.log('BasicModuleDialog 初始化表单数据:', basicForm.value);
}, { immediate: true });

// 监听表单变化，保存到 localStorage
watch(basicForm, (newFormData) => {
  localStorage.setItem('basicModuleDraft', JSON.stringify(newFormData));
  console.log('保存草稿到 localStorage (BasicModule):', newFormData);
}, { deep: true });

// 文案输入框的引用
const speechInput = ref<InputInst | null>(null);

// 点击标签插入到文案
const addSpeechTag = (tag: string) => {
  if (!speechInput.value) return;
  
  const textarea = speechInput.value.$el.querySelector('textarea');
  if (textarea) {
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const text = basicForm.value.speechContent;
    const tagWithBraces = `{${$t(tag)}}`;
    
    basicForm.value.speechContent = text.slice(0, start) + tagWithBraces + text.slice(end);
    
    setTimeout(() => {
      textarea.selectionStart = textarea.selectionEnd = start + tagWithBraces.length;
      textarea.focus();
    }, 0);
  }
};

// 退出并触发事件
const handleExit = () => {
  console.log('退出基础模块弹窗，当前表单数据:', basicForm.value);
  emits('exit', basicForm.value);
  emits('update:show', false);
};

// 保存并触发事件
const handleSave = () => {
  console.log('保存基础模块，当前表单数据:', basicForm.value);
  emits('save', basicForm.value);
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