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
            <div class="card-placeholder">+</div>
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
                <p>{{ truncateDescription(module.scriptContent.join('\n') || '无描述') }}</p>
                <!-- 展示音频名称和路径（仅音频模块） -->
                <p v-if="module.moduleType === ModuleType.Audio && module.audioName">
                  音频名称：{{ module.audioName }}
                </p>
                <p v-if="module.moduleType === ModuleType.Audio && module.audioPath">
                  音频路径：{{ module.audioPath }}
                </p>
              </div>
              <div class="module-actions">
                <n-button type="primary" @click="editModule(module)">编辑</n-button>
                <n-button type="warning" @click="testModule(module)">测试</n-button>
                <n-button type="error" @click="deleteModule(module)">删除</n-button>
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { NCard, NButton, NScrollbar } from 'naive-ui';
import BasicModuleDialog from '@/views/ModuleConfig/components/BasicModuleDialog.vue';
import AudioModuleDialog from '@/views/ModuleConfig/components/AudioModuleDialog.vue';

// 后端 API 基础 URL
const API_BASE_URL = 'http://localhost:7073';

// 触发条件枚举，与后端保持一致
enum TriggerCondition {
  ExecuteLoop = "ExecuteLoop",        // 循环执行
  BarrageComment = "BarrageComment",  // 弹幕评论
  SendGift = "SendGift",              // 送礼物
  Like = "Like",                      // 点赞
  EnterLiveRoom = "EnterLiveRoom",    // 进入直播间
  WarningTip = "WarningTip",          // 警告提示
}

// 模块类型枚举
enum ModuleType {
  Base = "base",    // 基础模块
  Audio = "audio",  // 音频模块
}

// 读取步骤枚举
enum ReadStep {
  Random = "random",        // 随机读取
  Sequential = "sequential" // 顺序读取
}

// 定义模块接口，与后端数据结构保持一致
interface Module {
  id: number;                           // 模块ID
  moduleType: ModuleType;               // 模块类型（base/audio）
  orderNum: number;                     // 排序编号
  moduleName: string;                   // 模块名称
  intervalTimeStart: number;            // 间隔时间起始（秒）
  intervalTimeEnd: number;              // 间隔时间结束（秒）
  triggerConditions: TriggerCondition[];// 触发条件列表
  readStep: ReadStep;                   // 读取步骤（random/sequential）
  scriptContent: string[];              // 话术文案列表
  isModelRewrite: boolean;              // 是否启用大模型改写
  rewriteFrequency: number;             // 改写频率（秒）
  audioName?: string;                   // 音频文件名称（仅音频模块使用）
  audioPath?: string;                   // 音频文件路径（仅音频模块使用）
}

// 控制加载状态
const loading = ref(false);

// 控制基础模块弹窗的显示
const showBasicModule = ref(false);
const editingModule = ref<Module | null>(null);

// 控制音频交互模块弹窗的显示
const showAudioModule = ref(false);
const editingAudioModule = ref<Module | null>(null);

// 模块数据
const modules = ref<Module[]>([]);

// 防止重复保存
let isSaving = false;

// 从后端获取模块数据
async function fetchModules() {
  console.log('开始加载模块数据...');
  loading.value = true;
  try {
    console.log('发起请求:', `${API_BASE_URL}/base-modules`);
    const response = await fetch(`${API_BASE_URL}/base-modules`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    console.log('响应状态:', response.status);
    
    const rawResponse = await response.text();
    console.log('原始响应:', rawResponse);
    
    if (!rawResponse) {
      throw new Error('后端返回空响应');
    }
    
    let backendData;
    try {
      backendData = JSON.parse(rawResponse);
      console.log('解析后的数据:', backendData);
    } catch (parseError) {
      console.error('JSON 解析失败:', parseError);
      throw new Error(`后端返回数据无法解析为 JSON: ${rawResponse}`);
    }
    
    if (!response.ok) {
      throw new Error(`获取模块列表失败，状态码: ${response.status}，错误信息: ${backendData.error || '未知错误'}`);
    }
    
    if (!Array.isArray(backendData)) {
      throw new Error('后端返回数据格式错误，期望数组');
    }
    
    const mappedData = backendData
      .map((item) => mapBackendToFrontend(item))
      .filter((module): module is Module => module !== null && module.id !== undefined);
    console.log('映射后的数据:', mappedData);
    
    modules.value = mappedData;
  } catch (error) {
    console.error('加载模块失败:', error);
    modules.value = [];
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

  return {
    id: backendModule.id || 0,
    moduleType: (backendModule.module_type || 'base') as ModuleType,
    orderNum: backendModule.order_num || 0,
    moduleName: backendModule.module_name || '未知模块',
    intervalTimeStart: backendModule.interval_time_start || 0,
    intervalTimeEnd: backendModule.interval_time_end || 0,
    triggerConditions: Array.isArray(backendModule.trigger_conditions)
      ? backendModule.trigger_conditions as TriggerCondition[]
      : [],
    readStep: (backendModule.read_step || 'random') as ReadStep,
    scriptContent: Array.isArray(backendModule.script_content)
      ? backendModule.script_content
      : [],
    isModelRewrite: backendModule.is_model_rewrite || false,
    rewriteFrequency: backendModule.rewrite_frequency || 0,
    audioName: backendModule.audio_name || '',
    audioPath: backendModule.audio_path || '',
  };
}

// 将前端数据映射到后端格式
function mapFrontendToBackend(frontendModule: Module): any {
  const maxOrderNum = modules.value.length > 0
    ? Math.max(...modules.value.map(m => m.orderNum))
    : 0;

  return {
    id: frontendModule.id,
    module_type: frontendModule.moduleType,
    order_num: frontendModule.id === 0 ? maxOrderNum + 1 : frontendModule.orderNum,
    module_name: frontendModule.moduleName,
    interval_time_start: frontendModule.intervalTimeStart,
    interval_time_end: frontendModule.intervalTimeEnd,
    trigger_conditions: frontendModule.triggerConditions,
    read_step: frontendModule.readStep,
    script_content: frontendModule.scriptContent,
    is_model_rewrite: frontendModule.isModelRewrite,
    rewrite_frequency: frontendModule.rewriteFrequency,
    audio_name: frontendModule.audioName || '',
    audio_path: frontendModule.audioPath || '',
  };
}

// 保存模块到后端
async function saveModule(module: Module, isEdit: boolean) {
  if (isSaving) {
    console.log('正在保存中，请勿重复提交');
    return;
  }
  if (module.intervalTimeStart >= module.intervalTimeEnd) {
    console.error('保存失败：intervalTimeStart 必须小于 intervalTimeEnd');
    return;
  }

  try {
    isSaving = true;
    const method = isEdit ? 'PUT' : 'POST';
    const url = isEdit
      ? `${API_BASE_URL}/base-modules/${module.id}`
      : `${API_BASE_URL}/base-modules`;
    const response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(mapFrontendToBackend(module)),
    });
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(`${isEdit ? '更新' : '新增'}模块失败，错误信息: ${errorData.error}`);
    }
    console.log(isEdit ? '模块更新成功' : '模块新增成功');
    await fetchModules();
  } catch (error) {
    console.error('保存模块失败:', error);
  } finally {
    isSaving = false;
  }
}

// 删除模块
async function deleteModuleFromBackend(module: Module) {
  try {
    const response = await fetch(`${API_BASE_URL}/base-modules/${module.id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(`删除模块失败，错误信息: ${errorData.error}`);
    }
    console.log('模块删除成功');
    await fetchModules();
  } catch (error) {
    console.error('删除模块失败:', error);
  }
}

// 组件挂载时加载数据
onMounted(async () => {
  console.log('组件挂载，开始加载后端数据');
  await fetchModules();
});

// 点击“基础模块”卡片时打开弹窗
function handleBasicModuleClick() {
  console.log('打开基础模块弹窗');
  editingModule.value = null;
  showBasicModule.value = true;
}

// 点击“音频交互模块”卡片时打开弹窗
function handleAudioModuleClick() {
  console.log('打开音频交互模块弹窗');
  editingAudioModule.value = null;
  showAudioModule.value = true;
}

// 处理基础模块弹窗退出时的提示
function handleBasicModuleExit(formData: any) {
  console.log('基础模块退出，当前表单数据:', formData);
}

// 处理基础模块弹窗保存时的逻辑
function handleBasicModuleSave(formData: any) {
  console.log('基础模块保存，当前表单数据:', formData);
  const module: Module = editingModule.value
    ? {
        ...editingModule.value,
        moduleName: formData.moduleName || editingModule.value.moduleName,
        intervalTimeStart: formData.minTime,
        intervalTimeEnd: formData.maxTime,
        triggerConditions: formData.triggers
          ? Object.entries(formData.triggers)
              .filter(([_, value]) => value)
              .map(([key]) => {
                switch (key) {
                  case 'cycleExecution':
                    return TriggerCondition.ExecuteLoop;
                  case 'elasticComment':
                    return TriggerCondition.BarrageComment;
                  case 'gift':
                    return TriggerCondition.SendGift;
                  case 'like':
                    return TriggerCondition.Like;
                  case 'joinRoom':
                    return TriggerCondition.EnterLiveRoom;
                  case 'notice':
                    return TriggerCondition.WarningTip;
                  default:
                    return null;
                }
              })
              .filter((condition): condition is TriggerCondition => condition !== null)
          : editingModule.value.triggerConditions,
        readStep: (formData.readStep || editingModule.value.readStep) as ReadStep,
        scriptContent: formData.speechContents || editingModule.value.scriptContent,
        isModelRewrite: formData.modelRewrite === 'yes',
        rewriteFrequency: formData.rewriteFrequency,
      }
    : {
        id: 0,
        moduleType: ModuleType.Base,
        orderNum: 0,
        moduleName: formData.moduleName || '基础模块',
        intervalTimeStart: formData.minTime,
        intervalTimeEnd: formData.maxTime,
        triggerConditions: formData.triggers
          ? Object.entries(formData.triggers)
              .filter(([_, value]) => value)
              .map(([key]) => {
                switch (key) {
                  case 'cycleExecution':
                    return TriggerCondition.ExecuteLoop;
                  case 'elasticComment':
                    return TriggerCondition.BarrageComment;
                  case 'gift':
                    return TriggerCondition.SendGift;
                  case 'like':
                    return TriggerCondition.Like;
                  case 'joinRoom':
                    return TriggerCondition.EnterLiveRoom;
                  case 'notice':
                    return TriggerCondition.WarningTip;
                  default:
                    return null;
                }
              })
              .filter((condition): condition is TriggerCondition => condition !== null)
          : [],
        readStep: (formData.readStep || 'random') as ReadStep,
        scriptContent: formData.speechContents || [],
        isModelRewrite: formData.modelRewrite === 'yes',
        rewriteFrequency: formData.rewriteFrequency || 0,
      };
  saveModule(module, !!editingModule.value);
}

// 处理音频交互模块弹窗退出时的提示
function handleAudioModuleExit(formData: any) {
  console.log('音频模块退出，当前表单数据:', formData);
}

// 处理音频交互模块弹窗保存时的逻辑
function handleAudioModuleSave(formData: any) {
  console.log('音频模块保存，当前表单数据:', formData);
  const module: Module = editingAudioModule.value
    ? {
        ...editingAudioModule.value,
        moduleName: formData.moduleName || editingAudioModule.value.moduleName,
        intervalTimeStart: formData.minTime,
        intervalTimeEnd: formData.maxTime,
        triggerConditions: formData.triggers
          ? Object.entries(formData.triggers)
              .filter(([_, value]) => value)
              .map(([key]) => {
                switch (key) {
                  case 'cycleExecution':
                    return TriggerCondition.ExecuteLoop;
                  case 'elasticComment':
                    return TriggerCondition.BarrageComment;
                  case 'gift':
                    return TriggerCondition.SendGift;
                  case 'like':
                    return TriggerCondition.Like;
                  case 'joinRoom':
                    return TriggerCondition.EnterLiveRoom;
                  case 'notice':
                    return TriggerCondition.WarningTip;
                  default:
                    return null;
                }
              })
              .filter((condition): condition is TriggerCondition => condition !== null)
          : editingAudioModule.value.triggerConditions,
        readStep: (formData.readStep || editingAudioModule.value.readStep) as ReadStep,
        scriptContent: formData.speechContents || editingAudioModule.value.scriptContent,
        audioName: formData.audioName || editingAudioModule.value.audioName,
        audioPath: formData.audioPath || editingAudioModule.value.audioPath,
      }
    : {
        id: 0,
        moduleType: ModuleType.Audio,
        orderNum: 0,
        moduleName: formData.moduleName || '音频模块',
        intervalTimeStart: formData.minTime,
        intervalTimeEnd: formData.maxTime,
        triggerConditions: formData.triggers
          ? Object.entries(formData.triggers)
              .filter(([_, value]) => value)
              .map(([key]) => {
                switch (key) {
                  case 'cycleExecution':
                    return TriggerCondition.ExecuteLoop;
                  case 'elasticComment':
                    return TriggerCondition.BarrageComment;
                  case 'gift':
                    return TriggerCondition.SendGift;
                  case 'like':
                    return TriggerCondition.Like;
                  case 'joinRoom':
                    return TriggerCondition.EnterLiveRoom;
                  case 'notice':
                    return TriggerCondition.WarningTip;
                  default:
                    return null;
                }
              })
              .filter((condition): condition is TriggerCondition => condition !== null)
          : [],
        readStep: (formData.readStep || 'random') as ReadStep,
        scriptContent: formData.speechContents || [],
        isModelRewrite: false,
        rewriteFrequency: 0,
        audioName: formData.audioName,
        audioPath: formData.audioPath,
      };
  saveModule(module, !!editingAudioModule.value);
}

// 编辑模块
function editModule(module: Module) {
  console.log('编辑模块:', module);
  if (module.moduleType === ModuleType.Audio) {
    editingAudioModule.value = { ...module };
    showAudioModule.value = true;
  } else {
    editingModule.value = { ...module };
    showBasicModule.value = true;
  }
}

// 测试模块
function testModule(module: Module) {
  console.log('测试模块:', module);
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
      case TriggerCondition.ExecuteLoop:
        return '循环执行';
      case TriggerCondition.BarrageComment:
        return '弹幕评论';
      case TriggerCondition.SendGift:
        return '送礼物';
      case TriggerCondition.Like:
        return '点赞';
      case TriggerCondition.EnterLiveRoom:
        return '进入直播间';
      case TriggerCondition.WarningTip:
        return '管理警告';
      default:
        return '';
    }
  }).filter(name => name);
  return conditionNames.length > 0 ? conditionNames.join('、') : '无限制';
}

// 截断描述以防止过长
function truncateDescription(description: string, maxLength: number = 50): string {
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
  }

  .module-actions {
    display: flex;
    gap: 10px;
  }
}
</style>