<template>
  <n-modal
    :show="show"
    @update:show="value => emits('update:show', value)"
    preset="dialog"
    title="数字人模块"
    style="width: 800px"
    :mask-closable="false"
  >
    <div class="number-person-module-form">
      <!-- 数字人功能按钮 -->
      <div class="form-row">
        <div class="form-item">
          <span class="label">数字人功能</span>
          <div class="button-group">
            <n-button type="primary" @click="handleAddDigitalHuman" :loading="isAddingDigitalHuman">
              新增数字人
            </n-button>
            <n-button type="success" @click="handleStartDigitalHuman" :loading="isStartingDigitalHuman" style="margin-left: 12px">
              开启数字人
            </n-button>
          </div>
          <div v-if="operationMessage" class="operation-message" :class="{ 'success': operationSuccess, 'error': !operationSuccess }">
            {{ operationMessage }}
          </div>
        </div>
      </div>

      <!-- 数字人配置 -->
      <!-- <div class="form-row">
        <div class="form-item">
          <span class="label">数字人名称 *</span>
          <n-input 
            v-model:value="numberPersonForm.moduleName" 
            placeholder="请输入数字人名称"
            @blur="handleModuleNameBlur"
          />
        </div>
      </div>

      数字人描述
      <div class="form-row">
        <div class="form-item">
          <span class="label">数字人描述</span>
          <n-input 
            v-model:value="numberPersonForm.description" 
            type="textarea" 
            placeholder="请输入数字人描述"
          />
        </div> -->
      <!-- </div> -->

      <!-- 添加网页嵌入框架 -->
      <div v-if="showWebFrame" class="web-frame-container">
        <div class="frame-header">
          <span>数字人Web应用</span>
          <n-button size="small" @click="closeWebFrame">关闭</n-button>
        </div>
        <iframe :src="webAppUrl" class="web-frame"></iframe>
      </div>
    </div>
    
    <!-- 底部按钮 -->
    <template #action>
      <div style="text-align: right">
        <n-button type="warning" @click="handleExit">退出</n-button>
       
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { NModal, NInputNumber, NButton, NRadioGroup, NRadio, NInput } from 'naive-ui';

// 引入Electron的ipcRenderer用于与主进程通信
const { ipcRenderer } = window.require('electron');

// 触发条件枚举
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
  const draft = localStorage.getItem('numberPersonModuleDraft');
  return draft
    ? JSON.parse(draft)
    : {
        minTime: 30,
        maxTime: 50,
        selectedTrigger: 'SceneLoop',
        moduleName: '数字人模块',
        description: '',
      };
};

// 表单数据
const numberPersonForm = ref(loadDraft());

// 操作状态
const isAddingDigitalHuman = ref(false);
const isStartingDigitalHuman = ref(false);
const operationMessage = ref('');
const operationSuccess = ref(true);

// 监听弹窗显示状态
watch(() => props.show, (newShow) => {
  if (newShow) {
    if (props.initialData) {
      // 编辑模式：使用初始数据
      numberPersonForm.value = {
        minTime: props.initialData.intervalTimeStart || 30,
        maxTime: props.initialData.intervalTimeEnd || 50,
        selectedTrigger: props.initialData.triggerConditions?.[0] || 'SceneLoop',
        moduleName: props.initialData.moduleName || '数字人模块',
        description: props.initialData.description || '',
      };
    } else {
      // 新建模式：重置表单
      numberPersonForm.value = loadDraft();
    }
    // 清除草稿
    localStorage.removeItem('numberPersonModuleDraft');
    // 重置操作消息
    operationMessage.value = '';
  } else {
    // 弹窗关闭时重置所有状态
    isAddingDigitalHuman.value = false;
    isStartingDigitalHuman.value = false;
    operationMessage.value = '';
    operationSuccess.value = true;
  }
}, { immediate: true });

// 初始化或编辑模式
watch(() => props.initialData, (newData) => {
  if (newData && !localStorage.getItem('numberPersonModuleDraft')) {
    numberPersonForm.value = {
      minTime: newData.intervalTimeStart || 30,
      maxTime: newData.intervalTimeEnd || 50,
      selectedTrigger: newData.triggerConditions?.[0] || 'SceneLoop',
      moduleName: newData.moduleName || '数字人模块',
      description: newData.description || '',
    };
  }
  console.log('NumberPersonDialog 初始化表单数据:', numberPersonForm.value);
}, { immediate: true });

// 监听表单变化，保存草稿
watch(numberPersonForm, (newFormData) => {
  localStorage.setItem('numberPersonModuleDraft', JSON.stringify(newFormData));
  console.log('保存草稿到 localStorage (NumberPersonModule):', newFormData);
}, { deep: true });

// 处理新增数字人按钮点击
const handleAddDigitalHuman = async () => {
  try {
    isAddingDigitalHuman.value = true;
    operationMessage.value = '正在启动新增数字人程序...';
    operationSuccess.value = true;
    
    // 调用Electron主进程执行批处理文件
    const result = await ipcRenderer.invoke('execute-bat-file', 'd:\\androidWork\\AingDesk\\build\\extraResources\\exe\\QimuDHweb\\01提取视频检查点.bat');
    
    console.log('新增数字人程序启动结果:', result);
    operationMessage.value = '新增数字人程序已启动';
    operationSuccess.value = true;
  } catch (error) {
    console.error('启动新增数字人程序失败:', error);
    operationMessage.value = `启动新增数字人程序失败: ${error.message || '未知错误'}`;
    operationSuccess.value = false;
  } finally {
    isAddingDigitalHuman.value = false;
  }
};

// 添加网页框架相关状态
const showWebFrame = ref(false);
const webAppUrl = ref('http://localhost:7860');

// 处理开启数字人按钮点击
const handleStartDigitalHuman = async () => {
  try {
    isStartingDigitalHuman.value = true;
    operationMessage.value = '正在启动数字人程序...';
    operationSuccess.value = true;
    
    // 执行批处理文件
    const result = await ipcRenderer.invoke('execute-bat-file', 'd:\\androidWork\\AingDesk\\build\\extraResources\\exe\\QimuDHweb\\02启动webapp.bat');
    
    console.log('数字人程序启动结果:', result);
    operationMessage.value = '数字人程序已启动 ✅';
    operationSuccess.value = true;
    
    // 10秒后清除提示消息
    setTimeout(() => {
      operationMessage.value = '';
    }, 10000);
    
  } catch (error) {
    console.error('启动数字人程序失败:', error);
    operationMessage.value = `启动数字人程序失败: ${error.message || '未知错误'}`;
    operationSuccess.value = false;
  } finally {
    isStartingDigitalHuman.value = false;
  }
};

// 关闭网页框架
const closeWebFrame = () => {
  showWebFrame.value = false;
};

// 退出
const handleExit = () => {
  console.log('退出数字人模块弹窗，当前表单数据:', numberPersonForm.value);
  // 清除草稿
  localStorage.removeItem('numberPersonModuleDraft');
  emits('exit', numberPersonForm.value);
  emits('update:show', false);
};

// 保存
const handleSave = () => {
  if (!numberPersonForm.value.moduleName || !numberPersonForm.value.moduleName.trim()) {
    window.$message?.error('保存失败：数字人名称不能为空');
    return;
  }

  const formData = {
    id: props.initialData?.id || 0,
    moduleType: 'numberPerson',
    moduleName: numberPersonForm.value.moduleName,
    intervalTimeStart: 0,
    intervalTimeEnd: 0,
    triggerConditions: [TriggerCondition.SceneLoop],
    readStep: props.initialData?.readStep || 'random',
    scriptContent: props.initialData?.scriptContent || [],
    isModelRewrite: props.initialData?.isModelRewrite || false,
    rewriteFrequency: props.initialData?.rewriteFrequency || 0,
    description: numberPersonForm.value.description || '',
  };

  console.log('保存数字人模块，当前表单数据:', formData);
  emits('save', formData);
  emits('update:show', false);
  
  // 清除草稿和重置表单
  localStorage.removeItem('numberPersonModuleDraft');
  numberPersonForm.value = loadDraft();
};

const handleModuleNameBlur = () => {
  if (!numberPersonForm.value.moduleName || numberPersonForm.value.moduleName.trim() === '') {
    numberPersonForm.value.moduleName = '数字人模块';
  }
  localStorage.setItem('numberPersonModuleDraft', JSON.stringify(numberPersonForm.value));
};
</script>

<style scoped lang="scss">
@use "@/assets/base";

.number-person-module-form {
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
    
    .button-group {
      margin-top: 10px;
    }
    
    .operation-message {
      margin-top: 10px;
      padding: 8px;
      border-radius: 4px;
      
      &.success {
        background-color: #f0f9eb;
        color: #67c23a;
      }
      
      &.error {
        background-color: #fef0f0;
        color: #f56c6c;
      }
    }
  }
}

// 添加网页框架样式
.web-frame-container {
  margin-top: 20px;
  border: 1px solid #eee;
  border-radius: 4px;
  overflow: hidden;
  
  .frame-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background-color: #f5f5f5;
    border-bottom: 1px solid #eee;
  }
  
  .web-frame {
    width: 100%;
    height: 500px;
    border: none;
  }
}
</style>