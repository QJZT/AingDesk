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
                <h3>{{ module.name }}</h3>
                <p>周期时间：{{ module.minTime }}~{{ module.maxTime }}s 触发条件：{{ getTriggerConditions(module) }}</p>
                <p>{{ truncateDescription(module.description || '无描述') }}</p>
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
  </div>

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
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { NCard, NButton, NScrollbar } from 'naive-ui';
import BasicModuleDialog from '@/views/ModuleConfig/components/BasicModuleDialog.vue';
import AudioModuleDialog from '@/views/ModuleConfig/components/AudioModuleDialog.vue';
import { eventBUS } from '@/views/Home/utils/tools';

// 后端 API 基础 URL
const API_BASE_URL = 'http://localhost:7077';

// 定义模块的接口
interface Module {
    id: number;
    name: string;
    description?: string;
    minTime: number;
    maxTime: number;
    triggers: {
        cycleExecution?: boolean;
        elasticComment?: boolean;
        gift?: boolean;
        like?: boolean;
        joinRoom?: boolean;
        notice?: boolean;
    };
    audioFiles?: File[];
    readStep?: 'random' | 'sequential';
    modelRewrite?: boolean;
    rewriteFrequency?: number;
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
        const response = await fetch(`${API_BASE_URL}/base-modules`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        console.log('API 响应状态:', response.status);
        const rawResponse = await response.text();
        console.log('原始响应内容:', rawResponse);
        console.log('原始响应长度:', rawResponse.length);
        console.log('原始响应前100个字符:', rawResponse.substring(0, 100));
        console.log('原始响应后100个字符:', rawResponse.substring(rawResponse.length - 100));
        try {
            JSON.parse(rawResponse);
            console.log('JSON 解析成功');
        } catch (parseError) {
            console.error('JSON 解析失败:', parseError);
            throw parseError;
        }
        if (!response.ok) {
            let errorData;
            try {
                errorData = JSON.parse(rawResponse);
            } catch (e) {
                throw new Error(`获取模块列表失败，状态码: ${response.status}，原始响应: ${rawResponse}`);
            }
            throw new Error(`获取模块列表失败，状态码: ${response.status}，错误信息: ${errorData.error}`);
        }
        const backendData = JSON.parse(rawResponse);
        console.log('后端返回数据:', backendData);
        if (!Array.isArray(backendData)) {
            throw new Error('后端返回数据格式错误，期望数组');
        }
        const mappedData = backendData.map((item) => {
            try {
                return mapBackendToFrontend(item);
            } catch (mapError) {
                console.error('映射数据时出错:', mapError, '原始数据:', item);
                return null;
            }
        }).filter(module => module && module.id);
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
// 将后端数据映射到前端格式的函数
// 参数 backendModule: 后端返回的单个模块数据（可能是任意类型，需要做类型检查）
// 返回值: 符合 Module 接口的对象，或者 null（如果数据无效）
function mapBackendToFrontend(backendModule: any): Module | null {
    if (!backendModule || typeof backendModule !== 'object') {
        console.warn('无效的后端数据:', backendModule);
        return null;
    }

    const triggers: Module['triggers'] = {
        cycleExecution: Array.isArray(backendModule.trigger_conditions) && backendModule.trigger_conditions.includes('ExecuteLoop') || false,
        elasticComment: Array.isArray(backendModule.trigger_conditions) && backendModule.trigger_conditions.includes('BarrageComment') || false,
        gift: Array.isArray(backendModule.trigger_conditions) && backendModule.trigger_conditions.includes('SendGift') || false,
        like: Array.isArray(backendModule.trigger_conditions) && backendModule.trigger_conditions.includes('Like') || false,
        joinRoom: Array.isArray(backendModule.trigger_conditions) && backendModule.trigger_conditions.includes('EnterLiveRoom') || false,
        notice: Array.isArray(backendModule.trigger_conditions) && backendModule.trigger_conditions.includes('WarningTip') || false
    };

    return {
        id: backendModule.id || 0, // 修正为小写 id
        name: backendModule.module_name || '未知模块', // 修正为 module_name
        description: backendModule.script_content || '', // 修正为 script_content
        minTime: backendModule.interval_time_start || 0, // 修正为 interval_time_start
        maxTime: backendModule.interval_time_end || 0, // 修正为 interval_time_end
        triggers,
        readStep: backendModule.read_step || 'random', // 修正为 read_step
        modelRewrite: backendModule.is_model_rewrite || false, // 修正为 is_model_rewrite
        rewriteFrequency: backendModule.rewrite_frequency || 0 // 修正为 rewrite_frequency
    };
}

// 将前端数据映射到后端格式
function mapFrontendToBackend(frontendModule: Module): any {
  const triggerConditions: string[] = [];
  if (frontendModule.triggers.cycleExecution) triggerConditions.push('ExecuteLoop');
  if (frontendModule.triggers.elasticComment) triggerConditions.push('BarrageComment');
  if (frontendModule.triggers.gift) triggerConditions.push('SendGift');
  if (frontendModule.triggers.like) triggerConditions.push('Like');
  if (frontendModule.triggers.joinRoom) triggerConditions.push('EnterLiveRoom');
  if (frontendModule.triggers.notice) triggerConditions.push('WarningTip');

  // 计算 OrderNum：如果是新增模块，设置为当前最大 OrderNum + 1
  const maxOrderNum = modules.value.length > 0 
    ? Math.max(...modules.value.map(m => m.id)) // 近似替代 OrderNum
    : 0;

  return {
    ID: frontendModule.id,
    OrderNum: frontendModule.id === 0 ? maxOrderNum + 1 : frontendModule.id,
    ModuleName: frontendModule.name,
    IntervalTimeStart: frontendModule.minTime,
    IntervalTimeEnd: frontendModule.maxTime,
    TriggerConditions: triggerConditions,
    ReadStep: frontendModule.readStep || 'random',
    ScriptContent: frontendModule.description || '',
    IsModelRewrite: frontendModule.modelRewrite || false,
    RewriteFrequency: frontendModule.rewriteFrequency || 0
  };
}

// 保存模块到后端
async function saveModule(module: Module, isEdit: boolean) {
    if (isSaving) {
        console.log('正在保存中，请勿重复提交');
        return;
    }
    if (module.minTime >= module.maxTime) {
        console.error('保存失败：minTime 必须小于 maxTime');
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
            body: JSON.stringify(mapFrontendToBackend(module))
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
      method: 'DELETE'
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
        ? { ...editingModule.value, ...formData, modelRewrite: formData.modelRewrite === 'yes' }
        : {
            id: 0,
            name: '基础模块',
            description: formData.speechContent || '用户输入的文案',
            minTime: formData.minTime,
            maxTime: formData.maxTime,
            triggers: formData.triggers,
            readStep: formData.readStep || 'random',
            modelRewrite: formData.modelRewrite === 'yes',
            rewriteFrequency: formData.rewriteFrequency
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
    ? { ...editingAudioModule.value, ...formData }
    : {
        id: 0,
        name: '音频模块',
        description: formData.speechContent || '',
        minTime: formData.minTime,
        maxTime: formData.maxTime,
        triggers: formData.triggers,
        readStep: formData.readStep || 'random',
        audioFiles: formData.audioFiles
      };
  saveModule(module, !!editingAudioModule.value);
}

// 编辑模块
function editModule(module: Module) {
  if (module.audioFiles) {
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
  const triggers = module.triggers || {};
  const conditions: string[] = [];
  if (triggers.cycleExecution) conditions.push('循环执行');
  if (triggers.elasticComment) conditions.push('弹幕评论');
  if (triggers.gift) conditions.push('送礼物');
  if (triggers.like) conditions.push('点赞');
  if (triggers.joinRoom) conditions.push('进入直播间');
  if (triggers.notice) conditions.push('管理警告');
  return conditions.length > 0 ? conditions.join('、') : '无限制';
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