<template>
  <div class="module-config">
    <n-card title="模块配置">
      <div class="top-section">
        <div class="module-cards">
          <!-- 基础模块 -->
          <div class="module-card">
            <div class="card-placeholder" @click="handleBasicModuleClick">+</div>
            <h3>基础模块</h3>
            <p>说明：文案配置是ai直播关键的一环，通过设定自动化生成对应文案，并支持基础使用方式和AI工作流对接</p>
          </div>

          <!-- 工作流模块 -->
          <div class="module-card">
            <div class="card-placeholder">+</div>
            <h3>工作流对象</h3>
            <p>说明：支持coze，dify等主流工作流对接</p>
          </div>

          <!-- 音文交互模块 -->
          <div class="module-card">
            <div class="card-placeholder" @click="handleAudioModuleClick">+</div>
            <h3>音频文件模块</h3>
            <p>说明：设定音频触发对接程序</p>
          </div>

          <!-- 数字人模块 -->
          <div class="module-card">
            <div class="card-placeholder" @click="handleNumberPersonModuleClick">+</div>
            <h3>数字人模块</h3>
            <p>说明：数字人对口型对接说明</p>
          </div>
        </div>
      </div>

      <div class="content-section">
        <n-scrollbar style="max-height: calc(100vh - 400px)">
          <div v-if="loading" class="loading-state">
            <p>加载中...</p>
          </div>
          <div v-else-if="modules.length === 0" class="empty-state">
            <p>暂无模块，请点击上方卡片添加新模块。</p>
          </div>
          <div v-else class="module-grid">
            <div class="module-item" v-for="module in modules" :key="module.id">
              <div class="module-header">
                <h3>{{ module.moduleName }}</h3>
                <p>周期时间：{{ module.intervalTimeStart }}~{{ module.intervalTimeEnd }}s 触发条件：{{ getTriggerConditions(module) }}</p>
                <div v-if="module.scriptContent && module.scriptContent.length > 0" class="script-content-box">
                  <p v-for="(script, index) in module.scriptContent" :key="index">
                    {{ script }}
                  </p>
                </div>
                <p v-else>无描述</p>
                <p v-if="module.moduleType === ModuleType.Audio && module.audioName">
                  音频名称：{{ module.audioName }}
                </p>
                <p v-if="module.moduleType === ModuleType.Audio && module.audioPath">
                  音频路径：{{ module.audioPath }}
                </p>
              </div>
              <div class="module-actions">
                <n-button type="primary" @click="editModule(module)">编辑</n-button>
                <!-- <n-button type="error" @click="deleteModule(module)">删除</n-button> -->
                <n-popconfirm
                  @positive-click="deleteModule(module)"
                  @negative-click="cancelDelete"
                  positive-text="确认"
                  negative-text="取消"
                >
                  <template #trigger>
                    <n-button type="error">删除</n-button>
                  </template>
                  确定要删除模块 "{{ module.moduleName }}" 吗？
                </n-popconfirm>
              </div>
            </div>
          </div>
        </n-scrollbar>
      </div>
    </n-card>

    <!-- 基础模块配置弹窗 -->
    <BasicModuleDialog
      :show="showBasicModule"
      @update:show="showBasicModule = $event"
      @exit="handleBasicModuleExit"
      @save="handleBasicModuleSave"
      :initial-data="editingModule"
    />

    <!-- 音频交互模块配置弹窗 -->
    <AudioModuleDialog
      :show="showAudioModule"
      @update:show="showAudioModule = $event"
      @exit="handleAudioModuleExit"
      @save="handleAudioModuleSave"
      :initial-data="editingAudioModule"
    />

    <!-- 数字人模块配置弹窗 -->
    <NumberPersonDialog
      :show="showNumberPersonModule"
      @update:show="showNumberPersonModule = $event"
      @exit="handleNumberPersonModuleExit"
      @save="handleNumberPersonModuleSave"
      :initial-data="editingNumberPersonModule"
    />
    
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { NCard, NButton, NScrollbar, NPopconfirm } from 'naive-ui';
import BasicModuleDialog from '@/views/ModuleConfig/components/BasicModuleDialog.vue';
import AudioModuleDialog from '@/views/ModuleConfig/components/AudioModuleDialog.vue';
import NumberPersonDialog from '@/views/dataListView/components/NumberPersonDialog.vue';

// 后端 API 基础 URL
const API_BASE_URL = 'http://localhost:7072';

// 触发条件枚举，与后端保持一致
enum TriggerCondition {
  SceneLoop = "SceneLoop",         // 控场循环
  IntervalLoop = "IntervalLoop",   // 间隔循环：按照指定间隔时间循环执行
  BarrageComment = "BarrageComment", // 弹幕评论：是否自动发送弹幕评论
  SendGift = "SendGift",           // 送礼物：是否自动发送礼物
  Like = "Like",                   // 点赞：是否自动点赞
  EnterLiveRoom = "EnterLiveRoom", // 进入直播间：是否自动进入直播间
  ShareRoom = "ShareRoom",         // 分享直播间
  FollowRoom = "FollowRoom",       // 关注直播间
}

// 模块类型枚举
enum ModuleType {
  Base = "base",
  Audio = "audio",
  NumberPerson = "numberPerson",
}

// 读取步骤枚举
enum ReadStep {
  Random = "random",
  Sequential = "sequential"
}

// 处理取消删除的操作
function cancelDelete() {
  console.log('删除操作已取消');
  window.$message?.info('已取消删除操作');
}

// 定义模块接口，与后端数据结构保持一致
interface Module {
  id: number;
  moduleType: ModuleType;
  orderNum: number;
  moduleName: string;
  intervalTimeStart: number;
  intervalTimeEnd: number;
  triggerConditions: TriggerCondition[];
  readStep: ReadStep;
  scriptContent: string[];
  isModelRewrite: boolean;
  rewriteFrequency: number;
  audioName?: string;
  audioPath?: string;
}

// 控制加载状态
const loading = ref(false);

// 控制基础模块弹窗的显示
const showBasicModule = ref(false);
const editingModule = ref<Module | null>(null);

// 控制音频交互模块弹窗的显示
const showAudioModule = ref(false);
const editingAudioModule = ref<Module | null>(null);

// 控制数字人模块弹窗的显示
const showNumberPersonModule = ref(false);
const editingNumberPersonModule = ref<Module | null>(null);

// 模块数据
const modules = ref<Module[]>([]);

// 防止重复保存
let isSaving = false;

// 从后端获取模块数据（查）
async function fetchModules() {
  console.log('开始加载模块数据...');
  loading.value = true;
  try {
    const response = await fetch(`${API_BASE_URL}/base-modules`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    if (!response.ok) {
      throw new Error(`获取模块列表失败，状态码: ${response.status}`);
    }
    const backendData = await response.json();
    console.log('后端返回数据:', backendData);

    const mappedData = backendData
      .map((item) => mapBackendToFrontend(item))
      .filter((module): module is Module => module !== null && module.id !== undefined);
    console.log('映射后的数据:', mappedData);

    // 清空并重新赋值，确保响应式更新
    modules.value = [];
    modules.value = [...mappedData];
  } catch (error) {
    console.error('加载模块失败:', error);
    modules.value = [];
    window.$message?.error('加载模块失败：' + error.message);
  } finally {
    loading.value = false;
    console.log('加载完成，当前模块数据:', modules.value);
  }
}

// 将后端数据映射到前端格式
function mapBackendToFrontend(backendModule: any): Module | null {
  if (!backendModule || typeof backendModule !== 'object') {
    console.warn('无效的后端数据:', backendModule);
    return null;
  }

  const triggerConditions = Array.isArray(backendModule.trigger_conditions)
    ? backendModule.trigger_conditions.map((condition: string) => {
        switch (condition) {
          case 'SceneLoop':
            return TriggerCondition.SceneLoop;
          case 'IntervalLoop':
            return TriggerCondition.IntervalLoop;
          case 'BarrageComment':
            return TriggerCondition.BarrageComment;
          case 'SendGift':
            return TriggerCondition.SendGift;
          case 'Like':
            return TriggerCondition.Like;
          case 'EnterLiveRoom':
            return TriggerCondition.EnterLiveRoom;
          case 'ShareRoom':
            return TriggerCondition.ShareRoom;
          case 'FollowRoom':
            return TriggerCondition.FollowRoom;
          default:
            console.warn('未知的触发条件:', condition);
            return condition as TriggerCondition;
        }
      })
    : [];

  let scriptContent: string[] = [];
  if (backendModule.script_content) {
    if (Array.isArray(backendModule.script_content)) {
      scriptContent = backendModule.script_content;
    } else if (typeof backendModule.script_content === 'string') {
      scriptContent = backendModule.script_content.split(',').map((s: string) => s.trim());
    }
  }

  return {
    id: backendModule.id || 0,
    moduleType: (backendModule.module_type || 'base') as ModuleType,
    orderNum: backendModule.order_num || 0,
    moduleName: backendModule.module_name || '未知模块',
    intervalTimeStart: backendModule.interval_time_start || 0,
    intervalTimeEnd: backendModule.interval_time_end || 0,
    triggerConditions,
    readStep: (backendModule.read_step || 'random') as ReadStep,
    scriptContent,
    isModelRewrite: backendModule.is_model_rewrite || false,
    rewriteFrequency: backendModule.rewrite_frequency || 0,
    audioName: backendModule.audio_name || '',
    audioPath: backendModule.audio_path || '',
  };
}

// 将前端数据映射到后端格式
function mapFrontendToBackend(frontendModule: Module): any {
  const maxOrderNum = modules.value.length > 0
    ? Math.max(...modules.value.map((m) => m.orderNum))
    : 0;

  const triggerConditions = frontendModule.triggerConditions.map((condition) => {
    switch (condition) {
      case TriggerCondition.SceneLoop:
        return 'SceneLoop';
      case TriggerCondition.IntervalLoop:
        return 'IntervalLoop';
      case TriggerCondition.BarrageComment:
        return 'BarrageComment';
      case TriggerCondition.SendGift:
        return 'SendGift';
      case TriggerCondition.Like:
        return 'Like';
      case TriggerCondition.EnterLiveRoom:
        return 'EnterLiveRoom';
      case TriggerCondition.ShareRoom:
        return 'ShareRoom';
      case TriggerCondition.FollowRoom:
        return 'FollowRoom';
      default:
        console.warn('未映射的触发条件:', condition);
        return condition;
    }
  });

  const backendData = {
    id: frontendModule.id,
    module_type: frontendModule.moduleType,
    order_num: frontendModule.id === 0 ? maxOrderNum + 1 : frontendModule.orderNum,
    module_name: frontendModule.moduleName,
    interval_time_start: frontendModule.intervalTimeStart,
    interval_time_end: frontendModule.intervalTimeEnd,
    trigger_conditions: triggerConditions,
    read_step: frontendModule.readStep,
    script_content: frontendModule.scriptContent,
    is_model_rewrite: frontendModule.isModelRewrite || false,
    rewrite_frequency: frontendModule.rewriteFrequency || 0,
    audio_name: frontendModule.audioName || '',
    audio_path: frontendModule.audioPath || '',
  };
  console.log('转换后的后端数据:', backendData);
  return backendData;
}

// 保存模块到后端（增/改）
async function saveModule(module: Module, isEdit: boolean) {
  if (isSaving) {
    console.log('正在保存中，请勿重复提交');
    return;
  }

  try {
    isSaving = true;
    const method = isEdit ? 'PUT' : 'POST';
    const url = isEdit
      ? `${API_BASE_URL}/base-modules/${module.id}`
      : `${API_BASE_URL}/base-modules`;

    const backendData = mapFrontendToBackend(module);
    console.log('发送到后端的最终数据:', backendData);

    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(backendData),
    });

    const responseData = await response.json();

    if (!response.ok) {
      throw new Error(responseData.error || `保存失败，状态码：${response.status}`);
    }

    console.log('后端保存响应:', responseData);
    await fetchModules(); // 刷新数据
    window.$message?.success(isEdit ? '更新成功' : '添加成功');
  } catch (error) {
    console.error('保存失败:', error);
    window.$message?.error('保存失败：' + (error.message || '未知错误'));
  } finally {
    isSaving = false;
  }
}

// 删除模块（删）
async function deleteModuleFromBackend(module: Module) {
  try {
    const response = await fetch(`${API_BASE_URL}/base-modules/${module.id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || '未知错误');
    }
    console.log('模块删除成功');
    await fetchModules();
    window.$message?.success('模块删除成功');
  } catch (error) {
    console.error('删除模块失败:', error);
    window.$message?.error('删除模块失败：' + error.message);
  }
}

// 组件挂载时加载数据
onMounted(async () => {
  console.log('组件挂载，开始加载后端数据');
  await fetchModules();
});

// 点击“基础模块”卡片时打开弹窗（始终以添加模式）
function handleBasicModuleClick() {
  console.log('打开基础模块弹窗');
  editingModule.value = null; // 始终以添加模式打开，填充默认数据
  showBasicModule.value = true;
}

// 点击“音频交互模块”卡片时打开弹窗
function handleAudioModuleClick() {
  console.log('打开音频交互模块弹窗');
  editingAudioModule.value = null; // 新增时设为 null
  showAudioModule.value = true;
}

// 处理基础模块弹窗退出时的提示
function handleBasicModuleExit(formData: any) {
  console.log('基础模块退出，当前表单数据:', formData);
}

// 处理基础模块弹窗保存时的逻辑
function handleBasicModuleSave(formData: any) {
  console.log('基础模块保存，当前表单数据:', formData);
  const module: Module = {
    id: formData.id || 0,
    moduleType: ModuleType.Base,
    orderNum: 0, // 新增时设为 0，由后端生成
    moduleName: formData.moduleName || '基础模块',
    intervalTimeStart: formData.intervalTimeStart,
    intervalTimeEnd: formData.intervalTimeEnd,
    triggerConditions: formData.triggerConditions.map((condition: string) => {
      switch (condition) {
        case 'SceneLoop':
          return TriggerCondition.SceneLoop;
        case 'IntervalLoop':
          return TriggerCondition.IntervalLoop;
        case 'BarrageComment':
          return TriggerCondition.BarrageComment;
        case 'SendGift':
          return TriggerCondition.SendGift;
        case 'Like':
          return TriggerCondition.Like;
        case 'EnterLiveRoom':
          return TriggerCondition.EnterLiveRoom;
        case 'ShareRoom':
          return TriggerCondition.ShareRoom;
        case 'FollowRoom':
          return TriggerCondition.FollowRoom;
        default:
          console.warn('未知触发条件:', condition);
          return condition as TriggerCondition;
      }
    }),
    readStep: (formData.readStep || 'random') as ReadStep,
    scriptContent: formData.scriptContent || [],
    isModelRewrite: formData.isModelRewrite || false,
    rewriteFrequency: formData.rewriteFrequency || 0,
  };
  console.log('构造的 Module 对象:', module);
  saveModule(module, !!formData.id);
}

// 处理音频交互模块弹窗退出时的提示
function handleAudioModuleExit(formData: any) {
  console.log('音频模块退出，当前表单数据:', formData);
}

// 处理音频交互模块弹窗保存时的逻辑
function handleAudioModuleSave(formData: any) {
  console.log('音频模块保存，当前表单数据:', formData);
  const module: Module = {
    id: formData.id || 0,
    moduleType: ModuleType.Audio,
    orderNum: 0, // 新增时设为 0，由后端生成
    moduleName: formData.moduleName || '音频模块',
    intervalTimeStart: formData.intervalTimeStart,
    intervalTimeEnd: formData.intervalTimeEnd,
    triggerConditions: formData.triggerConditions.map((condition: string) => {
      switch (condition) {
        case 'SceneLoop':
          return TriggerCondition.SceneLoop;
        case 'IntervalLoop':
          return TriggerCondition.IntervalLoop;
        case 'BarrageComment':
          return TriggerCondition.BarrageComment;
        case 'SendGift':
          return TriggerCondition.SendGift;
        case 'Like':
          return TriggerCondition.Like;
        case 'EnterLiveRoom':
          return TriggerCondition.EnterLiveRoom;
        case 'ShareRoom':
          return TriggerCondition.ShareRoom;
        case 'FollowRoom':
          return TriggerCondition.FollowRoom;
        default:
          console.warn('未知触发条件:', condition);
          return condition as TriggerCondition;
      }
    }),
    readStep: (formData.readStep || 'random') as ReadStep,
    scriptContent: formData.scriptContent || [],
    isModelRewrite: formData.isModelRewrite || false,
    rewriteFrequency: formData.rewriteFrequency || 0,
    audioName: formData.audioName || '',
    audioPath: formData.audioPath || '',
  };
  console.log('构造的 Module 对象:', module);
  saveModule(module, !!formData.id);
}

// 点击"数字人模块"卡片时打开弹窗
function handleNumberPersonModuleClick() {
  console.log('打开数字人模块弹窗');
  editingNumberPersonModule.value = null; // 新增时设为 null
  showNumberPersonModule.value = true;
}

// 处理数字人模块弹窗退出时的提示
function handleNumberPersonModuleExit(formData: any) {
  console.log('数字人模块退出，当前表单数据:', formData);
}

// 处理数字人模块弹窗保存时的逻辑
function handleNumberPersonModuleSave(formData: any) {
  console.log('数字人模块保存，当前表单数据:', formData);
  const module: Module = {
    id: formData.id || 0,
    moduleType: ModuleType.NumberPerson,
    orderNum: 0, // 新增时设为 0，由后端生成
    moduleName: formData.moduleName || '数字人模块',
    intervalTimeStart: formData.intervalTimeStart,
    intervalTimeEnd: formData.intervalTimeEnd,
    triggerConditions: formData.triggerConditions.map((condition: string) => {
      switch (condition) {
        case 'SceneLoop':
          return TriggerCondition.SceneLoop;
        case 'IntervalLoop':
          return TriggerCondition.IntervalLoop;
        case 'BarrageComment':
          return TriggerCondition.BarrageComment;
        case 'SendGift':
          return TriggerCondition.SendGift;
        case 'Like':
          return TriggerCondition.Like;
        case 'EnterLiveRoom':
          return TriggerCondition.EnterLiveRoom;
        case 'ShareRoom':
          return TriggerCondition.ShareRoom;
        case 'FollowRoom':
          return TriggerCondition.FollowRoom;
        default:
          console.warn('未知触发条件:', condition);
          return condition as TriggerCondition;
      }
    }),
    readStep: (formData.readStep || 'random') as ReadStep,
    scriptContent: formData.scriptContent || [],
    isModelRewrite: formData.isModelRewrite || false,
    rewriteFrequency: formData.rewriteFrequency || 0,
  };
  console.log('构造的 Module 对象:', module);
  saveModule(module, !!formData.id);
}

// 编辑模块
function editModule(module: Module) {
  console.log('编辑模块:', module);
  if (module.moduleType === ModuleType.Audio) {
    editingAudioModule.value = { ...module };
    showAudioModule.value = true;
  } else if (module.moduleType === ModuleType.NumberPerson) {
    editingNumberPersonModule.value = { ...module };
    showNumberPersonModule.value = true;
  } else {
    editingModule.value = { ...module };
    showBasicModule.value = true;
  }
}

// 删除模块
function deleteModule(module: Module) {
  deleteModuleFromBackend(module);
  console.log('删除模块:', module);
}

// 获取触发条件文本
function getTriggerConditions(module: Module) {
  const conditions = module.triggerConditions || [];
  const conditionNames: string[] = conditions.map((condition) => {
    switch (condition) {
      case TriggerCondition.SceneLoop:
        return '控场循环';
      case TriggerCondition.IntervalLoop:
        return '间隔循环';
      case TriggerCondition.BarrageComment:
        return '弹幕评论';
      case TriggerCondition.SendGift:
        return '送礼物';
      case TriggerCondition.Like:
        return '点赞';
      case TriggerCondition.EnterLiveRoom:
        return '进入直播间';
      case TriggerCondition.ShareRoom:
        return '分享直播间';
      case TriggerCondition.FollowRoom:
        return '关注直播间';
      default:
        return condition;
    }
  }).filter((name) => name);
  return conditionNames.length > 0 ? conditionNames.join('、') : '无限制';
}

// 截断描述以防止过长（保留但未使用）
function truncateDescription(description: string, maxLength: number = 30): string {
  if (description.length > maxLength) {
    return description.slice(0, maxLength) + '...';
  }
  return description;
}
</script>

<style scoped lang="scss">
@use "@/assets/base";

.module-config {
  padding: 20px;

  .top-section {
    margin-bottom: 20px;
  }

  .module-cards {
    display: flex;
    justify-content: space-between;
    gap: 30px;
  }

  .module-card {
    flex: 1;
    background: linear-gradient(135deg, #3b82f6, #8b5cf6);
    color: white;
    padding: 20px;
    border-radius: 8px;
    text-align: center;
    position: relative;

    .card-placeholder {
      font-size: 24px;
      position: absolute;
      top: 10px;
      right: 10px;
      cursor: pointer;
      width: 30px;
      height: 30px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: rgba(255, 255, 255, 0.2);
      border-radius: 50%;
      &:hover {
        background: rgba(255, 255, 255, 0.3);
      }
    }

    h3 {
      margin: 10px 0;
      font-size: 16px;
    }

    p {
      margin: 0;
      font-size: 12px;
    }
  }

  .content-section {
    width: 100%;
  }

  .loading-state {
    text-align: center;
    padding: 20px;
    color: #666;
  }

  .empty-state {
    text-align: center;
    padding: 20px;
    color: #666;
  }

  .module-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
  }

  .module-item {
    border: 1px solid #f0f0f0;
    padding: 20px;
    border-radius: 8px;
    min-height: 180px;
    box-sizing: border-box;
  }

  .module-header {
    margin-bottom: 10px;

    h3 {
      margin: 0;
      font-size: 16px;
    }

    p {
      margin: 5px 0;
      color: #666;
      font-size: 12px;
    }

    .script-content-box {
      max-height: 80px;
      overflow: hidden;
    }
  }

  .module-actions {
    display: flex;
    gap: 10px;
  }
}
</style>