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
            <n-radio-group v-model:value="basicForm.selectedTrigger">
              <n-radio value="cycleExecution">{{ $t('循环执行') }}</n-radio>
              <n-radio value="elasticComment">{{ $t('弹幕评论') }}</n-radio>
              <n-radio value="gift">{{ $t('送礼物') }}</n-radio>
              <n-radio value="like">{{ $t('点赞') }}</n-radio>
              <n-radio value="joinRoom">{{ $t('进入直播间') }}</n-radio>
              <n-radio value="notice">{{ $t('警告提示') }}</n-radio>
            </n-radio-group>
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

          <!-- n-dynamic-input 用于显示和编辑话术文案 -->
          <n-dynamic-input
            v-model:value="basicForm.speechContents"
            :placeholder="$t('请输入话术文案')"
            :min="1"
            :max="22"
            :on-create="handleSpeechContentCreate"
            :on-remove="handleSpeechContentRemove"
            class="speech-content-container"
          >
            <template #create-button-default>
              <n-button type="primary" size="small">添加文案</n-button>
            </template>
            <template #default="{ value, index }">
              <div class="speech-content-item">
                <n-input
                  type="textarea"
                  :value="value"
                  :placeholder="$t('请输入话术文案')"
                  :autosize="{ minRows: 1, maxRows: 4 }"
                  @update:value="newValue => updateSpeechContent(index, newValue)"
                  @focus="() => activeInputIndex.value = index"
                  maxlength="2000"
                />
                <div class="word-count">{{ value.length }}/2000</div>
              </div>
            </template>
          </n-dynamic-input>
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
import { NModal, NInput, NInputNumber, NRadio, NRadioGroup, NTag, NButton, NDynamicInput } from 'naive-ui';
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
        selectedTrigger: 'elasticComment',
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
        speechContents: [''] // 初始化时至少有一个空输入框
      };
};

// 表单数据，优先从 localStorage 加载草稿
const basicForm = ref(loadDraft());

// watch initialData 逻辑：监听 props.initialData 变化，初始化表单数据
watch(() => props.initialData, (newData) => {
  if (newData && !localStorage.getItem('basicModuleDraft')) {
    const selectedTrigger = Object.entries(newData.triggers || {})
      .find(([key, value]) => value === true)?.[0] || 'elasticComment';

    // 修改 1：直接使用 script_content 数组，而不是 description
    // 说明：script_content 是后端返回的字符串数组（逗号分隔），直接赋值给 speechContents
    let speechContents = newData.script_content || [];
    // 确保至少有一个输入框
    if (speechContents.length === 0) {
      speechContents = [''];
    }
    // 限制最大输入框数量为 22（与 n-dynamic-input 的 max 属性一致）
    if (speechContents.length > 22) {
      speechContents = speechContents.slice(0, 22);
    }

    basicForm.value = {
      minTime: newData.intervalTimeStart || 0, // 修改 2：使用 intervalTimeStart 和 intervalTimeEnd
      maxTime: newData.intervalTimeEnd || 0,
      selectedTrigger,
      triggers: {
        cycleExecution: false,
        elasticComment: false,
        gift: false,
        like: false,
        joinRoom: false,
        notice: false,
        [selectedTrigger]: true
      },
      readStep: newData.readStep || 'random',
      modelRewrite: newData.isModelRewrite ? 'yes' : 'no', // 修改 3：映射 isModelRewrite
      rewriteFrequency: newData.rewriteFrequency || 0,
      speechContents
    };
  }
  console.log('BasicModuleDialog 初始化表单数据:', basicForm.value);
}, { immediate: true });

// 监听 selectedTrigger 变化，同步更新 triggers
watch(() => basicForm.value.selectedTrigger, (newTrigger) => {
  Object.keys(basicForm.value.triggers).forEach(key => {
    basicForm.value.triggers[key] = false;
  });
  if (newTrigger) {
    basicForm.value.triggers[newTrigger] = true;
  }
  console.log('触发条件更新:', basicForm.value.triggers);
}, { immediate: true });

// 监听表单变化，保存到 localStorage
watch(basicForm, (newFormData) => {
  localStorage.setItem('basicModuleDraft', JSON.stringify(newFormData));
  console.log('保存草稿到 localStorage (BasicModule):', newFormData);
}, { deep: true });

// n-dynamic-input 相关方法
const handleSpeechContentCreate = () => {
  return '';
};

const handleSpeechContentRemove = (index: number) => {
  console.log(`移除文案输入框: 索引 ${index}`);
};

const updateSpeechContent = (index: number, newValue: string) => {
  basicForm.value.speechContents[index] = newValue;
};

// addSpeechTag 逻辑
const activeInputIndex = ref<number | null>(null);

const addSpeechTag = (tag: string) => {
  if (activeInputIndex.value === null) {
    console.warn('未选择任何输入框，无法插入标签');
    return;
  }

  const index = activeInputIndex.value;
  const currentValue = basicForm.value.speechContents[index];
  const tagWithBraces = `{${$t(tag)}}`;
  basicForm.value.speechContents[index] = currentValue + tagWithBraces;

  setTimeout(() => {
    const textarea = document.querySelectorAll('.speech-content-item textarea')[index] as HTMLTextAreaElement;
    if (textarea) {
      textarea.focus();
      textarea.selectionStart = textarea.selectionEnd = currentValue.length + tagWithBraces.length;
    }
  }, 0);
};

// 保存逻辑
const handleSave = () => {
  if (basicForm.value.minTime >= basicForm.value.maxTime) {
    console.error('保存失败：最小时间必须小于最大时间');
    return;
  }

  const nonEmptyContents = basicForm.value.speechContents.filter(content => content.trim());
  if (nonEmptyContents.length === 0) {
    console.error('保存失败：至少需要输入一条非空话术文案');
    return;
  }

  // 修改 4：构造符合后端 Module 接口的数据
  // 说明：将 speechContents 直接作为 scriptContent 数组，不用 join('\n')
  const formData = {
    id: props.initialData?.id || 0, // 如果是编辑模式，传递 id
    moduleType: 'base', // 基础模块
    orderNum: props.initialData?.orderNum || 0,
    moduleName: props.initialData?.moduleName || '基础模块',
    intervalTimeStart: basicForm.value.minTime,
    intervalTimeEnd: basicForm.value.maxTime,
    triggerConditions: Object.entries(basicForm.value.triggers)
      .filter(([_, value]) => value)
      .map(([key]) => ({ type: key })), // 转换为 TriggerCondition 数组
    readStep: basicForm.value.readStep,
    scriptContent: basicForm.value.speechContents.filter(content => content.trim()), // 过滤空内容
    isModelRewrite: basicForm.value.modelRewrite === 'yes',
    rewriteFrequency: basicForm.value.rewriteFrequency
  };

  console.log('保存基础模块，当前表单数据:', formData);
  emits('save', formData);
  emits('update:show', false);
};

// 退出逻辑
const handleExit = () => {
  console.log('退出基础模块弹窗，当前表单数据:', basicForm.value);
  emits('exit', basicForm.value);
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

    .speech-content-container {
      display: flex;
      flex-direction: column;
      gap: 10px;

      :deep(.n-dynamic-input-item) {
        display: flex;
        align-items: center;
        gap: 5px;
      }

      .speech-content-item {
        display: flex;
        flex-direction: column;
        width: 740px;
        margin-right: 0;

        .n-input {
          textarea {
            padding: 5px 8px;
            font-size: 14px;
            line-height: 1.5;
            border-radius: 4px;
          }
        }

        .word-count {
          text-align: right;
          color: #999;
          font-size: 12px;
          margin-top: 2px;
        }
      }

      :deep(.n-button--small) {
        padding: 2px 6px;
        font-size: 12px;
      }

      :deep(.n-dynamic-input-create-button) {
        margin-top: 10px;
        .n-button {
          padding: 2px 8px;
          font-size: 12px;
        }
      }
    }
  }
}
</style>