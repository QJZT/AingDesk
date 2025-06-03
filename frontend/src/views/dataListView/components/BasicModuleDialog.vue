<template>
  <n-modal
    :show="show"
    @update:show="value => emits('update:show', value)"
    preset="dialog"
    style="width: 830px"
    :mask-closable="false"
    :show-icon="false"
  >
    <template #header>
      <div style="display: flex; align-items: center; width: 100%; gap: 8px;">
        <span style="font-size: 16px; font-weight: bold; color: #18a058;">模块备注</span>
        <n-input
          v-model:value="basicForm.moduleName"
          placeholder="请输入模块名称"
          style="width: 120px; font-size: 16px; font-weight: bold;"
          @blur="handleModuleNameBlur"
        />
      </div>
    </template>
    <div class="basic-module-form">
      <!-- 触发条件（第一排） -->
      <div class="form-row">
        <div class="form-item">
          <span class="label">{{ $t('触发条件') }} *</span>
          <div class="trigger-conditions">
            <n-radio-group v-model:value="basicForm.selectedTrigger">
              <n-radio value="sceneLoop">{{ $t('控场循环') }}</n-radio>
              <n-radio value="intervalLoop">{{ $t('间隔循环') }}</n-radio>
              <n-radio value="barrageComment">{{ $t('弹幕评论') }}</n-radio>
              <n-radio value="sendGift">{{ $t('送礼物') }}</n-radio>
              <n-radio value="like">{{ $t('点赞') }}</n-radio>
              <n-radio value="enterLiveRoom">{{ $t('进入直播间') }}</n-radio>
              <n-radio value="shareRoom">{{ $t('分享直播间') }}</n-radio>
              <n-radio value="followRoom">{{ $t('关注直播间') }}</n-radio>
            </n-radio-group>
          </div>
        </div>
      </div>

      <!-- 时间间隔（第二排） -->
      <div class="form-row" v-if="basicForm.selectedTrigger !== 'sceneLoop'">
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

      <!-- 读取步骤 -->
      <div class="form-row" v-if="isReadStepSupported">
        <div class="form-item">
          <span class="label">{{ $t('读取步骤') }} *</span>
          <n-radio-group v-model:value="basicForm.readStep">
            <n-radio value="random">{{ $t('随机') }}</n-radio>
            <n-radio value="sequential">{{ $t('顺序') }}</n-radio>
          </n-radio-group>
        </div>
      </div>

      <!-- 文案 -->
      <div class="form-row">
        <div class="form-item">
          <!-- 动态标题：弹幕评论时显示“提示词”，其他显示“话术文案” -->
          <span class="label">{{ basicForm.selectedTrigger === 'barrageComment' ? $t('提示词') : $t('话术文案') }} *</span>

          <!-- 动态显示支持的变量标签 -->
          <div class="speech-tags">
            <span class="label">{{ $t('变量支持') }}：</span>
            <!-- {
              "下个整点": "2",
              "再过五分钟(时分)": "2点01分",
              "现在日期(月日)": "5月9日",
              "现在时间(时分)": "1点56分",
              "现在时间(时分秒)": "1点56分18秒"
            } -->
            <n-tag v-if="isVariableSupported('现在时间(时分)')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('现在时间(时分)')">{{ $t('现在时间(时分)') }}</n-tag>
            <n-tag v-if="isVariableSupported('现在时间(时分秒)')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('现在时间(时分秒)')">{{ $t('现在时间(时分秒)') }}</n-tag>
            <n-tag v-if="isVariableSupported('现在日期(月日)')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('现在日期(月日)')">{{ $t('现在日期(月日)') }}</n-tag>
            <n-tag v-if="isVariableSupported('下个整点')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('下个整点')">{{ $t('下个整点') }}</n-tag>
            <n-tag v-if="isVariableSupported('再过五分钟(时分)')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('再过五分钟(时分)')">{{ $t('再过五分钟(时分)') }}</n-tag>
            <n-tag v-if="isVariableSupported('在线人数')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('在线人数')">{{ $t('在线人数') }}</n-tag>
            <n-tag v-if="isVariableSupported('直播间名称')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('直播间名称')">{{ $t('直播间名称') }}</n-tag>
            <n-tag v-if="isVariableSupported('点赞用户名')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('点赞用户名')">{{ $t('点赞用户名') }}</n-tag>
            <n-tag v-if="isVariableSupported('送礼用户名')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('送礼用户名')">{{ $t('送礼用户名') }}</n-tag>
            <n-tag v-if="isVariableSupported('礼物名称')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('礼物名称')">{{ $t('礼物名称') }}</n-tag>
            <n-tag v-if="isVariableSupported('弹幕内容')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('弹幕内容')">{{ $t('弹幕内容') }}</n-tag>
            <!-- <n-tag v-if="isVariableSupported('直播间警告内容')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('直播间警告内容')">{{ $t('直播间警告内容') }}</n-tag> -->
            <n-tag v-if="isVariableSupported('弹幕用户')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('弹幕用户')">{{ $t('弹幕用户') }}</n-tag>
            <n-tag v-if="isVariableSupported('礼物数量')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('礼物数量')">{{ $t('礼物数量') }}</n-tag>
            <n-tag v-if="isVariableSupported('进入直播间用户名')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('进入直播间用户名')">{{ $t('进入直播间用户名') }}</n-tag>
            <n-tag v-if="isVariableSupported('分享直播间用户名')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('分享直播间用户名')">{{ $t('分享直播间用户名') }}</n-tag>
            <n-tag v-if="isVariableSupported('关注直播用户名')" class="speech-tag" @mousedown.prevent="prepareAddSpeechTag('关注直播用户名')">{{ $t('关注直播用户名') }}</n-tag>
          </div>

          <!-- 弹幕评论使用单个大输入框 -->
          <div v-if="basicForm.selectedTrigger === 'barrageComment'" class="speech-content-container">
            <div class="speech-content-item">
              <n-input
                type="textarea"
                v-model:value="singleSpeechContent"
                :placeholder="$t('请输入提示词')"
                :autosize="{ minRows: 4, maxRows: 8 }"
                @focus="handleFocus(-1)"
                @blur="handleBlur(-1)"
                maxlength="2000"
              />
              <div class="word-count">{{ singleSpeechContent.length }}/2000</div>
            </div>
          </div>

          <!-- 其他触发条件使用动态输入框 -->
          <div v-else class="speech-content-container">
            <n-dynamic-input
              v-model:value="basicForm.speechContents"
              :placeholder="$t('请输入话术文案')"
              :min="1"
              :max="2000"
              :on-create="handleSpeechContentCreate"
              :on-remove="handleSpeechContentRemove"
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
                    @focus="handleFocus(index)"
                    @blur="handleBlur(index)"
                    maxlength="2000"
                  />
                  <div class="word-count">{{ value.length }}/2000</div>
                </div>
              </template>
            </n-dynamic-input>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部按钮 -->
    <template #action>
      <div style="text-align: right">
        <n-button type="warning" @click="handleExit">退出</n-button>
        <n-button type="primary" @click="handleSave" style="margin-left: 12px">保存</n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch, computed, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';
import { NModal, NInput, NInputNumber, NRadio, NRadioGroup, NTag, NButton, NDynamicInput } from 'naive-ui';
import type { Ref } from 'vue';

const { t: $t } = useI18n();

const props = defineProps<{
  show: boolean;
  initialData: any; // 父组件传递的初始数据
}>();

const emits = defineEmits<{
  (e: 'update:show', value: boolean): void;
  (e: 'exit', value: any): void;
  (e: 'save', value: any): void;
}>();

// 初始化表单数据
const basicForm = ref({
  minTime: 0,                    // 最小时间间隔（秒）
  maxTime: 0,                    // 最大时间间隔（秒）
  selectedTrigger: 'sceneLoop', // 修改默认值
  triggers: {                    // 触发条件的状态映射
    seneLoop: true,           // 控场循环
    intervalLoop: false,         // 间隔循环
    barrageComment: false,       // 弹幕评论
    sendGift: false,            // 送礼物
    like: false,                // 点赞
    enterLiveRoom: false,       // 进入直播间
    shareRoom: false,       // 分享直播间
    followRoom: false,        // 关注主播
  },
  readStep: 'random',           // 读取步骤：random-随机，sequential-顺序
  moduleName: '基础模块',        // 模块名称
  speechContents: [''],         // 话术文案内容数组
});

const singleSpeechContent = ref(basicForm.value.speechContents[0] || '');

const isReadStepSupported = computed(() => {
  const supportedTriggers = ['sceneLoop', 'intervalLoop', 'sendGift', 'like', 'enterLiveRoom', 'shareRoom', 'followRoom'];
  return supportedTriggers.includes(basicForm.value.selectedTrigger);
});

const isVariableSupported = computed(() => (variable: string) => {
  const variableMap = {
    '现在时间(时分)': ['sceneLoop', 'intervalLoop', 'barrageComment', 'enterLiveRoom'],
    '现在时间(时分秒)': ['sceneLoop', 'intervalLoop', 'barrageComment', 'enterLiveRoom'],
    '现在日期(月日)': ['sceneLoop', 'intervalLoop', 'barrageComment', 'enterLiveRoom'],
    '下个整点': ['sceneLoop', 'intervalLoop', 'barrageComment', 'enterLiveRoom'],
    '再过五分钟(时分)': ['sceneLoop', 'intervalLoop', 'barrageComment', 'enterLiveRoom'],
    '在线人数': ['sceneLoop', 'intervalLoop'],
    '直播间名称': ['sceneLoop', 'intervalLoop'],
    '点赞用户名': ['like'],
    '送礼用户名': ['sendGift'],
    '礼物名称': ['sendGift'],
    '弹幕内容': ['barrageComment'],
    // '直播间警告内容': ['sceneLoop', 'intervalLoop'],
    '弹幕用户': ['barrageComment'],
    '礼物数量': ['sendGift'],
    '进入直播间用户名': ['enterLiveRoom'],
    '分享直播间用户名': ['shareRoom'],
    '关注直播用户名': ['followRoom'],
  };
  return variableMap[variable as keyof typeof variableMap]?.includes(basicForm.value.selectedTrigger) || false;
});

// 监听 props.initialData，用于编辑模式
watch(() => props.initialData, (newData) => {
  if (newData) {
    // 编辑模式：使用 initialData 初始化表单
    basicForm.value.moduleName = newData.moduleName || '基础模块';
    basicForm.value.minTime = newData.intervalTimeStart || 0;
    basicForm.value.maxTime = newData.intervalTimeEnd || 0;
    
    // 修改触发条件的映射逻辑
    let selectedTrigger = 'sceneLoop';
    if (newData.triggerConditions && newData.triggerConditions.length > 0) {
      switch (newData.triggerConditions[0]) {
        case 'SceneLoop':
          selectedTrigger = 'sceneLoop';
          break;
        case 'IntervalLoop':
          selectedTrigger = 'intervalLoop';
          break;
        case 'BarrageComment':
          selectedTrigger = 'barrageComment';
          break;
        case 'SendGift':
          selectedTrigger = 'sendGift';
          break;
        case 'Like':
          selectedTrigger = 'like';
          break;
        case 'EnterLiveRoom':
          selectedTrigger = 'enterLiveRoom';
          break;
        case 'ShareRoom':
          selectedTrigger = 'shareRoom';
          break;
        case 'FollowRoom':
          selectedTrigger = 'followRoom';
          break;
      }
    }
    basicForm.value.selectedTrigger = selectedTrigger;
    
    // 更新触发条件的选中状态
    Object.keys(basicForm.value.triggers).forEach(key => {
      basicForm.value.triggers[key] = key === selectedTrigger;
    });
    basicForm.value.readStep = newData.readStep || 'random';
    basicForm.value.speechContents = newData.scriptContent || [''];
    singleSpeechContent.value = newData.scriptContent[0] || '';
  } else {
    // 添加模式：重置表单
    basicForm.value = {
      minTime: 0,
      maxTime: 0,
      selectedTrigger: 'sceneLoop',
      triggers: {
        sceneLoop: true,            // 修改为与前面定义一致的值
        intervalLoop: false,
        barrageComment: false,
        sendGift: false,
        like: false,
        enterLiveRoom: false,
        shareRoom: false,           // 修改为与前面定义一致的值
        followRoom: false,          // 修改为与前面定义一致的值
      },
      readStep: 'random',
      moduleName: '基础模块', // 默认值设置为“基础模块”
      speechContents: [''],
    };
    singleSpeechContent.value = '';
  }
}, { immediate: true });

// 监听 selectedTrigger 的变化
watch(() => basicForm.value.selectedTrigger, (newTrigger) => {
  Object.keys(basicForm.value.triggers).forEach(key => {
    basicForm.value.triggers[key] = key === newTrigger;
  });
  if (newTrigger === 'barrageComment') {
    singleSpeechContent.value = basicForm.value.speechContents[0] || '';
  } else if (basicForm.value.selectedTrigger === 'barrageComment') {
    basicForm.value.speechContents = singleSpeechContent.value ? [singleSpeechContent.value] : [''];
  }
}, { immediate: true });

// 监听 singleSpeechContent 的变化
watch(singleSpeechContent, (newValue) => {
  if (basicForm.value.selectedTrigger === 'barrageComment') {
    basicForm.value.speechContents = [newValue];
  }
});

enum TriggerCondition {
  SceneLoop = "SceneLoop",         // 控场循环
  IntervalLoop = "IntervalLoop",   // 间隔循环
  BarrageComment = "BarrageComment", // 弹幕评论
  SendGift = "SendGift",           // 送礼物
  Like = "Like",                   // 点赞
  EnterLiveRoom = "EnterLiveRoom", // 进入直播间
  ShareRoom = "ShareRoom",         // 分享直播间
  FollowRoom = "FollowRoom",       // 关注直播间
}

const truncateText = (text: string, maxLength: number = 50): string => {
  if (text.length > maxLength) {
    return text.slice(0, maxLength) + '...';
  }
  return text;
};

const handleSpeechContentCreate = () => '';
const handleSpeechContentRemove = (index: number) => console.log(`移除文案输入框: 索引 ${index}`);
const updateSpeechContent = (index: number, newValue: string) => {
  basicForm.value.speechContents[index] = newValue;
};

const activeInputIndex = ref<number | null>(null);
const activeInputRef = ref<HTMLTextAreaElement | null>(null);
const lastSelection = ref<{ start: number; end: number } | null>(null);
let pendingTag: string | null = null;

const handleFocus = (index: number) => {
  activeInputIndex.value = index;
  nextTick(() => {
    const textarea = index === -1
      ? document.querySelector('.speech-content-item textarea') as HTMLTextAreaElement
      : document.querySelectorAll('.speech-content-item textarea')[index] as HTMLTextAreaElement;
    if (textarea) {
      activeInputRef.value = textarea;
      lastSelection.value = { start: textarea.selectionStart, end: textarea.selectionEnd };
    }
  });
};

const handleBlur = (index: number) => {
  const textarea = index === -1
    ? document.querySelector('.speech-content-item textarea') as HTMLTextAreaElement
    : document.querySelectorAll('.speech-content-item textarea')[index] as HTMLTextAreaElement;
  if (textarea) {
    lastSelection.value = { start: textarea.selectionStart, end: textarea.selectionEnd };
  }
};

const prepareAddSpeechTag = (tag: string) => {
  pendingTag = tag;
  const textarea = activeInputRef.value;
  if (textarea) {
    lastSelection.value = { start: textarea.selectionStart, end: textarea.selectionEnd };
  }
  setTimeout(() => {
    addSpeechTag(tag);
    pendingTag = null;
  }, 0);
};

const addSpeechTag = (tag: string) => {
  let textarea = activeInputRef.value;
  if (!textarea && activeInputIndex.value !== null) {
    textarea = activeInputIndex.value === -1
      ? document.querySelector('.speech-content-item textarea') as HTMLTextAreaElement
      : document.querySelectorAll('.speech-content-item textarea')[activeInputIndex.value] as HTMLTextAreaElement;
    if (textarea) {
      textarea.focus();
      activeInputRef.value = textarea;
    }
  }

  if (!textarea) {
    console.warn('未找到聚焦的输入框，无法插入标签');
    return;
  }

  let valueRef: Ref<string>;
  let index: number | null = activeInputIndex.value;

  if (basicForm.value.selectedTrigger === 'barrageComment') {
    valueRef = singleSpeechContent;
    index = -1;
  } else {
    if (index === null) {
      console.warn('未选择动态输入框，无法插入标签');
      return;
    }
    valueRef = ref(basicForm.value.speechContents[index]);
  }

  const start = lastSelection.value?.start || 0;
  const end = lastSelection.value?.end || 0;
  const currentValue = valueRef.value || '';
  const tagWithBraces = `{${$t(tag)}}`;
  const newValue = currentValue.substring(0, start) + tagWithBraces + currentValue.substring(end);

  if (index === -1) {
    singleSpeechContent.value = newValue;
  } else {
    basicForm.value.speechContents[index] = newValue;
  }

  nextTick(() => {
    textarea.focus();
    const newPosition = start + tagWithBraces.length;
    textarea.selectionStart = newPosition;
    textarea.selectionEnd = newPosition;
    lastSelection.value = { start: newPosition, end: newPosition };
  });
};

const handleSave = () => {
  // 校验必填字段
  if (!basicForm.value.moduleName.trim()) {
    window.$message?.error('保存失败：模块名称不能为空');
    return;
  }


  let intervalTimeStart = basicForm.value.minTime;
  let intervalTimeEnd = basicForm.value.maxTime;
  if (basicForm.value.selectedTrigger === 'sceneLoop') {
    intervalTimeStart = 0;
    intervalTimeEnd = 1; 
  } else if (basicForm.value.minTime >= basicForm.value.maxTime) {
    window.$message?.error('保存失败：最小时间必须小于最大时间');
    return;
  }

  const nonEmptyContents = basicForm.value.speechContents.filter(content => content.trim());
  if (nonEmptyContents.length === 0) {
    window.$message?.error('保存失败：至少需要输入一条非空话术文案');
    return;
  }

  let triggerCondition;
  switch (basicForm.value.selectedTrigger) {
    case 'sceneLoop':
      triggerCondition = TriggerCondition.SceneLoop;
      break;
    case 'intervalLoop':
      triggerCondition = TriggerCondition.IntervalLoop;
      break;
    case 'barrageComment':
      triggerCondition = TriggerCondition.BarrageComment;
      break;
    case 'sendGift':
      triggerCondition = TriggerCondition.SendGift;
      break;
    case 'like':
      triggerCondition = TriggerCondition.Like;
      break;
    case 'enterLiveRoom':
      triggerCondition = TriggerCondition.EnterLiveRoom;
      break;
    case 'shareRoom':
      triggerCondition = TriggerCondition.ShareRoom;
      break;
    case 'followRoom':
      triggerCondition = TriggerCondition.FollowRoom;
      break;
    default:
      triggerCondition = TriggerCondition.SceneLoop;
  }

  const formData = {
    id: props.initialData?.id || 0,
    moduleType: 'base',
    moduleName: basicForm.value.moduleName,
    intervalTimeStart: intervalTimeStart,
    intervalTimeEnd: intervalTimeEnd,
    triggerConditions: [triggerCondition],
    readStep: basicForm.value.readStep,
    scriptContent: nonEmptyContents,
    isModelRewrite: false,
    rewriteFrequency: 0
  };

  console.log('保存基础模块，当前表单数据:', formData);
  emits('save', formData);
  emits('update:show', false);
  
  // 清除本地存储的草稿
  localStorage.removeItem('basicModuleDraft');
  
  // 重置表单数据为默认值
  basicForm.value = {
    minTime: 0,
    maxTime: 0,
    selectedTrigger: 'sceneLoop',
    triggers: {
      sceneLoop: true,            // 修改为与前面定义一致的值
      intervalLoop: false,
      barrageComment: false,
      sendGift: false,
      like: false,
      enterLiveRoom: false,
      shareRoom: false,           // 修改为与前面定义一致的值
      followRoom: false,          // 修改为与前面定义一致的值
    },
    readStep: 'random',
    moduleName: '基础模块',
    speechContents: [''],
  };
  singleSpeechContent.value = '';
};

const handleExit = () => {
  console.log('退出基础模块弹窗，当前表单数据:', basicForm.value);
  emits('exit', basicForm.value);
  emits('update:show', false);
};

const handleModuleNameBlur = () => {
  if (!basicForm.value.moduleName || basicForm.value.moduleName.trim() === '') {
    basicForm.value.moduleName = '基础模块';
  }
  localStorage.setItem('basicModuleDraft', JSON.stringify(basicForm.value));
};
</script>

<style scoped lang="scss">
@use "@/assets/base";

.basic-module-form {
  padding: 20px;
  max-height: 70vh;
  overflow-y: auto;
  box-sizing: border-box;
  overflow-x: hidden;

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
      gap: 10px;
      flex-wrap: wrap;
      justify-content: flex-start;

      :deep(.n-radio) {
        flex: 1 0 150px;
        max-width: 180px;
        box-sizing: border-box;
      }
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
        width: 100%;
        max-width: 760px;
        margin-right: 0;

        .n-input {
          textarea {
            padding: 5px 8px;
            font-size: 14px;
            line-height: 1.5;
            border-radius: 4px;
            width: 100%;
            box-sizing: border-box;
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

  /* 强制隐藏图标，保留其他样式 */
  :deep(.n-dialog__icon) {
    display: none !important;
  }
}
</style>