<template>
  <div class="module-config">
    <n-card :title="$t('模块配置')">
      <div class="top-section">
        <div class="module-cards">
          <!-- 基础模块 -->
          <div class="module-card">
            <div class="card-placeholder" @click="handleBasicModuleClick">+</div>
            <h3>{{ $t('基础模块') }}</h3>
            <p>{{ $t('说明：文案需要在当前音色影响的一环，通过设置全局性成文或文案，井共带你使用方式AI工作流对象') }}</p>
          </div>

          <!-- 工作流模块 -->
          <div class="module-card">
            <div class="card-placeholder">+</div>
            <h3>{{ $t('工作流对象') }}</h3>
            <p>{{ $t('说明：为coze，dify 等生成式AI工作流对象') }}</p>
          </div>
          <!-- 音文交互模块 -->
          <div class="module-card">
            <div class="card-placeholder" @click="handleAudioModuleClick">+</div>
            <h3>{{ $t('音文交互模块') }}</h3>
            <p>{{ $t('说明：设置影音成文或音档影响') }}</p>
          </div>
          <!-- 数字人模块 -->
          <div class="module-card">
            <div class="card-placeholder">+</div>
            <h3>{{ $t('数字人对象') }}</h3>
            <p>{{ $t('说明：数字人可以影响成文使用') }}</p>
          </div>
        </div>
      </div>

      <div class="content-section">
        <n-scrollbar style="max-height: calc(100vh - 400px)">
          <div class="module-grid">
            <div class="module-item" v-for="module in modules" :key="module.id">
              <div class="module-header">
                <h3>{{ module.name }}</h3>
                <p>{{ $t('周期时间') }}：{{ module.minTime }}~{{ module.maxTime }}s {{ $t('触发条件') }}：{{ getTriggerConditions(module) }}</p>
                <p>{{ module.description || module.audioFiles?.[0]?.name || '无描述' }}</p>
              </div>
              <div class="module-actions">
                <n-button type="primary" @click="editModule(module)">{{ $t('编辑') }}</n-button>
                <n-button type="warning" @click="testModule(module)">{{ $t('测试') }}</n-button>
                <n-button type="error" @click="deleteModule(module)">{{ $t('删除') }}</n-button>
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
    @submit="handleModuleSubmit"
    :initial-data="editingModule"
  />

  <!-- 音频交互模块配置弹窗 -->
  <AudioModuleDialog 
    :show="showAudioModule"
    @update:show="showAudioModule = $event"
    @submit="handleAudioModuleSubmit"
    :initial-data="editingAudioModule"
  />
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue';
import { NCard, NButton, NScrollbar } from 'naive-ui';
import { useI18n } from 'vue-i18n';
import BasicModuleDialog from './components/BasicModuleDialog.vue';
import AudioModuleDialog from './components/AudioModuleDialog.vue';
import { eventBUS } from '@/views/Home/utils/tools';

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
}

interface TopCard {
  id: number;
  name: string;
  description: string;
}

// 国际化
const { t: $t } = useI18n();

// 控制基础模块弹窗的显示
const showBasicModule = ref(false);
const editingModule = ref<Module | null>(null);

// 控制音频交互模块弹窗的显示
const showAudioModule = ref(false);
const editingAudioModule = ref<Module | null>(null);

// 模块数据
const modules = ref<Module[]>([]);

// 顶部卡片数据
const topCards = ref<TopCard[]>([]);

// 从 localStorage 加载初始模块数据
onMounted(() => {
  const savedModules = localStorage.getItem('modules');
  if (savedModules) {
    modules.value = JSON.parse(savedModules);
  } else {
    modules.value = [
      { id: 1, name: '控场话术模块', description: '相关数据...', minTime: 400, maxTime: 50, triggers: { elasticComment: true } },
      { id: 2, name: '弹幕话术模块', description: '相关数据...', minTime: 400, maxTime: 50, triggers: { elasticComment: true } },
      { id: 3, name: '收礼物话术模块', description: '相关数据...', minTime: 400, maxTime: 50, triggers: { gift: true } },
      { id: 4, name: '弹幕话术', description: '相关数据...', minTime: 400, maxTime: 50, triggers: { elasticComment: true } },
      { id: 5, name: '收礼乎音频', description: '相关数据...', minTime: 400, maxTime: 50, triggers: { gift: true } },
      { id: 6, name: '收礼物话术模块', description: '相关数据...', minTime: 400, maxTime: 50, triggers: { gift: true } },
      { id: 7, name: '弹幕话术', description: '相关数据...', minTime: 400, maxTime: 50, triggers: { elasticComment: true } },
      { id: 8, name: '收礼乎音频', description: '相关数据...', minTime: 400, maxTime: 50, triggers: { gift: true } }
    ];
  }

  topCards.value = [
    { id: 1, name: '基础模块', description: '说明：文案需要在当前音色影响的一环，通过设置全局性成文或文案，井共带你使用方式AI工作流对象' },
    { id: 2, name: '工作流对象', description: '说明：为coze，dify 等生成式AI工作流对象' },
    { id: 3, name: '音文交互模块', description: '说明：设置影音成文或音档影响' },
    { id: 4, name: '数字人对象', description: '说明：数字人可以影响成文使用' }
  ];
});

// 每次 modules 更新时保存到 localStorage
function saveModules() {
  localStorage.setItem('modules', JSON.stringify(modules.value));
}

// 点击“基础模块”卡片时打开弹窗
function handleBasicModuleClick() {
  editingModule.value = null;
  showBasicModule.value = true;
}

// 点击“音频交互模块”卡片时打开弹窗
function handleAudioModuleClick() {
  editingAudioModule.value = null;
  showAudioModule.value = true;
}

// 处理基础模块弹窗提交的数据
function handleModuleSubmit(formData: any) {
  if (editingModule.value) {
    const index = modules.value.findIndex(m => m.id === editingModule.value.id);
    if (index !== -1) {
      modules.value[index] = {
        ...modules.value[index],
        name: '基础模块（已编辑）',
        description: formData.speechContent || '用户输入的文案',
        minTime: formData.minTime,
        maxTime: formData.maxTime,
        triggers: formData.triggers
      };
      eventBUS.$emit('moduleUpdated', modules.value[index]);
    }
  } else {
    const newModule: Module = {
      id: modules.value.length + 1,
      name: '新基础模块',
      description: formData.speechContent || '用户输入的文案',
      minTime: formData.minTime,
      maxTime: formData.maxTime,
      triggers: formData.triggers
    };
    modules.value.push(newModule);
    eventBUS.$emit('moduleAdded', newModule);
  }
  saveModules();
}

// 处理音频交互模块弹窗提交的数据
function handleAudioModuleSubmit(formData: any) {
  if (editingAudioModule.value) {
    const index = modules.value.findIndex(m => m.id === editingAudioModule.value.id);
    if (index !== -1) {
      modules.value[index] = {
        ...modules.value[index],
        name: '音频模块（已编辑）',
        minTime: formData.minTime,
        maxTime: formData.maxTime,
        triggers: formData.triggers,
        audioFiles: formData.audioFiles
      };
      eventBUS.$emit('moduleUpdated', modules.value[index]);
    }
  } else {
    const newModule: Module = {
      id: modules.value.length + 1,
      name: '新音频模块',
      minTime: formData.minTime,
      maxTime: formData.maxTime,
      triggers: formData.triggers,
      audioFiles: formData.audioFiles
    };
    modules.value.push(newModule);
    eventBUS.$emit('moduleAdded', newModule);
  }
  saveModules();
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
  modules.value = modules.value.filter(m => m.id !== module.id);
  saveModules();
  eventBUS.$emit('moduleDeleted', module);
  console.log('删除模块:', module);
}

// 获取触发条件文本
function getTriggerConditions(module: Module) {
  const triggers = module.triggers || {};
  const conditions: string[] = [];
  if (triggers.cycleExecution) conditions.push($t('新兵执行'));
  if (triggers.elasticComment) conditions.push($t('弹幕野怪'));
  if (triggers.gift) conditions.push($t('送礼物'));
  if (triggers.like) conditions.push($t('点赞'));
  if (triggers.joinRoom) conditions.push($t('进入直播间'));
  if (triggers.notice) conditions.push($t('管理提示'));
  return conditions.length > 0 ? conditions.join('、') : $t('无限制');
}

// 监听事件
eventBUS.$on('moduleAdded', (module: Module) => {
  nextTick(() => {
    console.log('模块已添加:', module);
  });
});

eventBUS.$on('moduleUpdated', (module: Module) => {
  nextTick(() => {
    console.log('模块已更新:', module);
  });
});

eventBUS.$on('moduleDeleted', (module: Module) => {
  nextTick(() => {
    console.log('模块已删除:', module);
  });
});

eventBUS.$on('executeAudioOnce', (formData: any) => {
  console.log('执行音频一次:', formData);
});
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