<template>
  <div class="live-console" style="" :style="{backgroundColor: themeThinkBg}" >
    <!-- 设置区域 -->
     <div style="padding: 12px 12px 12px 12px;">
      <n-card title="基础设置" style="" :segmented="{content: true,footer:true}"
                    header-style="padding:10px;font-size:14px"
                    footer-style="padding:10px" content-style="padding:10px;height:100%">
                    <div class="settings-bar">
      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">时区设置</n-text>
        <n-select 
          v-model:value="selectedTimezone"
          :options="timezoneOptions"
          placeholder="选择时区"
          style="width: 250px;"
          @update:value="handleLanguageChange"
        />
      </div>
      
      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">语种设置 {{ selectedLanguageLabel }}</n-text>
        <n-select 
          v-model:value="selectedLanguage"
          :options="languageOptions"
          placeholder="选择语种"
          style="width: 150px;"
          @update:value="(value, option) => {
            selectedLanguageLabel = option.label
            handleLanguageChange(value)
          }"
        />
      </div>
      
      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">音频驱动</n-text>
        <n-select 
          v-model:value="selectedAudioDriver"
          :options="audioDeviceOptions"
          placeholder="选择驱动"
          style="width: 180px;"
          @update:value="handleLanguageChange"
        />
      </div>

      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">语音模型</n-text>
        <n-select 
          v-model:value="selectedSpeechModel"
          :options="speedModelOptions"
          placeholder="选择驱动"
          style="width: 180px;"
          @update:value="initializeSpeechModel"
        />
      </div>
    </div>
    <div class="settings-bar">
      <div class="setting-item">
       
       <n-text depth="3" style="margin-right: 8px;">Ai提示词:</n-text>
       <div>
        <n-button-group>
          <n-button ghost @click="showPromptDialog = true">
            改写提示词
          </n-button>
          <n-button ghost>
            弹幕提示词
          </n-button>
        </n-button-group>
       </div>
     </div>
      <div class="setting-item">
       
        <n-text depth="3" style="margin-right: 8px;">知识库:</n-text>
        <div>
          <n-popover trigger="click" style="">
                        <template #trigger>
                            <n-button :type="activeKnowledgeForChat.length ? 'primary' : 'default'" ghost
                                style="height: 34px;" icon-placement="left" :focusable="false">
                                <template #icon>
                                    <i class="i-tdesign:folder"></i>
                                </template>
                                {{ $t("知识库") }}
                            </n-button>
                        </template>
                        <KnowledgeChoosePanel />
      </n-popover>
        </div>
      </div>

      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">模型选择</n-text>
        <n-select
          v-model:value="selectedModel"
          :options="Object.entries(showModelList).map(([category, items]) => ({
            label: category,
            value: category,
            disabled: true,  // 添加disabled属性禁止选择
            children: items.map(item => ({
              label: item.title,
              value: item.model,
              supplierName: item.supplierName
            }))
          }))"
          placeholder="请选择模型"
          style="width: 400px;"
          clearable
         @update:value="(value, option) => handleSelectedModelChange(value, option)"
        />
      </div>
      
       <!-- 新增直播链接输入和启动按钮 -->
       <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">确认启动</n-text>
        <n-input-group>
          <n-button 
            v-if="!start"
            type="primary" 
            @click="registerModules"
            :loading="loading"
            style="width: 120px;"
          >
            <template #icon>
              <i class="i-tdesign:play-circle"></i>
            </template>
            启动系统
          </n-button>

          <n-button 
            v-if="start"
            type="error" 
            @click="start = false"
            :loading="loading"
            style="width: 120px;"
          >
            <template #icon>
              <i class="i-tdesign:stop-circle"></i>
            </template>
            停止系统
          </n-button>
        </n-input-group>
      </div>
    </div>
    </n-card>   
     </div>
                 
   
    <!-- 三列布局 -->
    <!-- height: calc(100vh - 460px) -->
    <!-- :style="{backgroundColor: themeThinkBg}" -->
    <div class="three-column-layout">
      <n-card title="模块" style="" :segmented="{content: true,footer:true}"
                    header-style="padding:10px;font-size:14px"
                    footer-style="padding:10px" content-style="padding:10px;height:100%">
        <n-infinite-scroll style="height: 480px;overflow-y: auto;background: var(--n-color-embedded);
      border-radius: var(--n-border-radius);"  :distance="10">
          <div class="">
            <div class="control-block" v-for="(block, index) in modules" :key="index" style="">
              <n-divider v-if="index != 0"  style="font-size: 15px;height: 1px;margin: 12px 0;"></n-divider>
            <div class="block-header">
              <h4>{{ block.module_name }}</h4>
            <!-- <n-switch v-model:value="block.isActive"  size="small" @update:value="handleVolumeChange(block)"/> -->
            <n-switch v-model:value="block.retAi"  size="small" @update:value="handleVolumeChange(block)">
              <template #checked>
                AI改写
              </template>
              <template #unchecked>
                AI改写
              </template>
            </n-switch>

            <n-switch v-model:value="block.isActive"  size="small" @update:value="handleVolumeChange(block)">
              <template #checked>
                开启
              </template>
              <template #unchecked>
                关闭
              </template>
            </n-switch>

            </div>
            <div style="height: 5px;"></div>
            <div style="display: flex; align-items: center; gap: 8px; width: 100%;">
              <div style="width: 30%;font-size: 12px;">
              音量:{{block.volume}}%
            </div>
            <n-slider
                v-model:value="block.volume"
                :min="0"
                :max="100"
                :step="1"
                style="width: 70%;"
                tooltip
                placement="bottom"
                :format-tooltip="(value) => `音量: ${value}%`"
                @update:value="handleVolumeChange(block)"
              />
            </div>
            <div style="height: 5px;"></div>
            <div style="display: flex; align-items: center; gap: 8px; width: 100%;">
              <div style="width: 30%;font-size: 12px">
                音速:{{block.speed}}x
              </div>
              <n-slider
                v-model:value="block.speed"
                :min="0.01"
                :max="2"
                :step="0.01"
                style="width: 70%;"
                tooltip
                placement="bottom"
                :format-tooltip="(value) => `音速: ${value}x`"
                 @update:value="handleVolumeChange(block)"
              />
            </div> 
            <div style="height: 5px;"></div>
            <div style="display: flex; align-items: center; gap: 8px; width: 100%;">
              <div style="width: 30%;font-size: 12px;">
                选择人声
            </div>
            <n-select 
                v-model:value="block.selectedNameId"
                :options="humanVoiceOptions"
                style="width: 70%;"
                placeholder="选择人声"
                clearable
                size="small"
                @update:value="handleVolumeChange(block)"
              />
            </div>
          </div>  
        </div>
        </n-infinite-scroll>
      </n-card>

      <!-- 观众互动 -->
      <!-- <n-card title="虚拟列表-滚动定位" style="height:100%" :segmented="{content: true,footer:true}"
                    header-style="padding:10px;font-size:14px"
                    footer-style="padding:10px" content-style="padding:10px;height:100%"> -->

      <n-card title="观众互动" style="" :segmented="{content: true,footer:true}"
                    header-style="padding:10px;font-size:14px"
                    footer-style="padding:10px" content-style="padding:10px;height:100%">
        <!-- <template #header-extra>
          <n-switch 
            v-model:value="autoReadMode" 
          >
            <template #checked>自动朗读</template>
            <template #unchecked>关闭朗读</template>
          </n-switch>
        </template> -->
        <div class="setting-item" >
        <n-text depth="3" style="margin-right: 8px;">直播弹幕接入</n-text>
        <n-input 
          v-model:value="streamUrl"
          placeholder="输入直播间链接"
        >
          <template #suffix>
            <n-button 
              type="info" 
              @click="toggleLive"
              :loading="loading"
              text
            >
            接入
            </n-button>
            <div style="width: 5px;"></div>
            <n-button 

              type="warning" 
              @click="outToggleLive"
              :loading="loading"
              text
            >
            断开
            </n-button>
          </template>
        </n-input>
      </div>
      <n-divider style="height: 1px;margin: 12px 0;"></n-divider>
            <div style="">
                <n-checkbox-group v-model:value="selectedFilters">
                  <div style="margin-bottom: 16px;">
              <n-space align="center" size="small" style="flex-wrap: wrap;">
                <n-text depth="3" style="font-size: 14px; margin-right: 2px;">筛选查看:</n-text>
                <n-button  strong secondary
                  v-for="filter in filterOptions"
                  :key="filter.value"
                  :type="selectedFilters.includes(filter.value) ? 'primary' : 'tertiary'" 
                  size="small" 
                  round
                  checkable
                  :checked="selectedFilters.includes(filter.value)"
                  @click="toggleFilter(filter.value)"
                  style="margin-right: 0px;"
                >
                  {{ filter.label }}
                </n-button>
              </n-space>
            </div>
            </n-checkbox-group>
          <n-divider style="height: 1px;margin: 12px 0;">
          </n-divider>
        </div>
        <n-infinite-scroll style="height: 180px; overflow-y: auto;" :style="{backgroundColor: themeThinkBg}"  :distance="10">
          <div >
            <div 
            v-for="(msg, index) in messages" 
            :key="index" 
            class="message-item"
            style="white-space: nowrap; overflow: hidden; text-overflow: ellipsis; display: flex; align-items: center; gap: 4px;"
          >
          <!-- getMessageColor -->
          <span class="content" style="line-height: 1;" :style="{
            color: getMessageColor(msg.data_types),
            whiteSpace: 'nowrap',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            display: 'flex',
            alignItems: 'center',
            gap: '4px',
            filter: 'brightness(1.2) contrast(1.2)'
          }">{{ msg.content }}</span>
        </div>
          </div>
        </n-infinite-scroll>
        <n-divider title-placement="left" style="font-size: 15px;height: 1px;margin: 12px 0;">
          </n-divider>
        <n-infinite-scroll style="height: 80px; overflow-y: auto;" :style="{backgroundColor: themeThinkBg}"  :distance="10">
          <div >
            <div 
             v-for="(data, index) in playList"
            :key="index" 
            class="message-item"
            style="white-space: nowrap; overflow: hidden; text-overflow: ellipsis; display: flex; align-items: center; gap: 4px;"
          >
          <!-- getMessageColor -->
          <span class="content" style="line-height: 1;" :style="{
            // color: getMessageColor(msg.data_types),
            whiteSpace: 'nowrap',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            display: 'flex',
            alignItems: 'center',
            gap: '4px',
            filter: 'brightness(1.2) contrast(1.2)'
          }">{{ data.content }}</span>
        </div>
          </div>
        </n-infinite-scroll>

        <!-- <div class="setting-item">
        <n-input 
          v-model:value="EnterLiveRoomUserName"
          placeholder="输入测试用户名"
          style="width: 200px;"
        /> -->
      <!-- </div> -->
        <!-- <div class="message-input">
          <n-input-group>
            <n-button text @click="toggleInputMode" class="mode-switch">
                <template #icon>
                  <i :class="isVoiceMode ? 'i-tdesign:keyboard w-20 h-20' : 'i-tdesign:microphone-1 w-20 h-20'"></i>
                </template>
            </n-button>
            
            <n-input 
              v-if="!isVoiceMode"
              v-model:value="newMessage" 
              placeholder="发送消息..."
              @keyup.enter="sendAnnouncement"
            />
            
            <n-button 
              v-else
              type="primary" 
              round
              class="voice-button"
            >
                <template #icon>
                  <i class="i-tdesign:microphone-1 w-20 h-20"></i>
                </template>
              <span>按住说话</span>
            </n-button>
            
            <n-button 
              type="primary" 
              @click="sendAnnouncement"
              :disabled="isVoiceMode"
            >
              <template #icon>
                  <i class="i-tdesign:send w-20 h-20"></i>
              </template>
            </n-button>
          </n-input-group>
        </div> -->
      </n-card>

      <!-- 直播数据 -->
      <!-- <n-card title="音频数据" style="height:100%" :segmented="{content: true,footer:true}"
                    header-style="padding:10px;font-size:14px"
                    footer-style="padding:10px" content-style="padding:10px;height:100%"> -->
        <!-- 设置间隔时间 -->
        <!-- <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">间隔时间(ms)</n-text>
        <n-input-number 
          v-model:value="intervalTime"
          :min="0"
          :max="500000"
          :step="100"
          style="width: 100px"
        />
      </div> -->
      <!-- 分割线 -->
      <!-- <n-divider  style="font-size: 15px;height: 1px;margin: 12px 0;"></n-divider>
        
        <n-infinite-scroll style="height: 280px; overflow-y: auto;" :style="{backgroundColor: themeThinkBg}" :distance="10">
          <div class="audio-blocks"> -->
            <!-- usePlayList -->
          <!-- <div class="audio-block" v-for="(data, index) in playList" :key="index" style="">
            <div class="audio-info">
              {{ data.content }} -->
              <!-- <div style="display: flex; justify-content: space-between; align-items: center;">
                <h4>{{ audio.title }}</h4>
                <n-button 
                  type="error" 
                  text
                  @click=""
                >
                  <template #icon>
                    <i class="i-tdesign:delete w-15 h-15"></i>
                  </template>
                </n-button>
              </div>
              <p>{{ audio.description }}</p> -->
            <!-- </div> -->
            <!-- <audio controls :src="audio.url"></audio> -->
          <!-- </div>
        </div>
        </n-infinite-scroll>   
      </n-card> -->
        <!-- <div class="stats-grid">
          <n-statistic label="观看人数" :value="viewerCount" />
          <n-statistic label="点赞数" :value="likeCount" />
          <n-statistic label="评论数" :value="commentCount" />
          <n-statistic label="直播时长" :value="duration" />
        </div>
      </n-card> -->
    </div>

    <!-- 互动消息区域 -->
    <!-- display: grid;
    // grid-template-rows: 50px 22px 2fr 22px 1fr 140px;
    grid-template-rows: 50px 50px 50px 50px 50px 25px 22px 1fr 35px  1fr 170px; -->

    <n-modal v-model:show="showPromptDialog">
      <n-card style="width:800px" title="编辑AI提示词" :bordered="false">
        <n-space vertical>
          <n-space align="center" size="small">
            <n-text>变量</n-text>
            <n-button size="small" @click="promptText += '{语种}'">语种</n-button>
          </n-space>
          <n-input
            v-model:value="promptText"
            type="textarea"
            placeholder="请输入AI提示词..."
            :autosize="{ minRows: 10 }"
          />
        </n-space>
        <template #footer>
          <n-space justify="end">
            <n-button @click="showPromptDialog = false">取消</n-button>
            <n-button type="primary" @click="handlePromptSave">保存</n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { 
  NCard, 
  NButton, 
  NTag,
  NInput,
  NInputGroup,
  NStatistic,
  NScrollbar
} from 'naive-ui'
import { io } from "socket.io-client"
import { log } from 'mermaid/dist/logger.js'
import KnowledgeChoosePanel from "@/views/KnowleadgeStore/components/KnowledgeChoosePanel.vue";
import { getKnowledgeStoreData } from "../KnowleadgeStore/store";
import { getHeaderStoreData } from '../Header/store';
import { getSoftSettingsStoreData } from '@/views/SoftSettings/store';
import { useMessage } from "naive-ui";
const {  showModelList } = getHeaderStoreData()
const {
    themeColors,
    themeMode
} = getSoftSettingsStoreData()
const themeThinkBg = computed(() => {
    if (themeMode.value == "light") {
        return themeColors.value.thinkWrapperLight
    } else {
        {
            return themeColors.value.thinlWrapperDark
        }
    }
})
const autoReadMode = ref(false)
const selectedModel = ref('') // 当前选中的模型
const isLive = ref(false)
const loading = ref(false)
const streamUrl = ref('')
const newMessage = ref('')
const isVoiceMode = ref(false)

const model_api = ref("")
const parameters_api = ref("")
const supplierName_api = ref("")
const handleSelectedModelChange = async (value, option) => {
  if (value == "") {
    model_api.value = ""
    parameters_api.value = ""
    return
  }
  supplierName_api.value = option.supplierName
  if (option.supplierName == "ollama") {
        [model_api.value, parameters_api.value] = option.value.split(":")
    } else {
        model_api.value = option.value
        parameters_api.value = ""
    }
}


const {
    activeKnowledge,
    activeKnowledgeForChat,
} = getKnowledgeStoreData()
function copyStreamUrl() {
  navigator.clipboard.writeText(streamUrl.value)
}

function sendAnnouncement() {
  if (newMessage.value.trim()) {
    messages.value.push({
      user: '系统公告',
      content: newMessage.value
    })
    newMessage.value = ''
  }
}

function toggleInputMode() {
  isVoiceMode.value = !isVoiceMode.value
}


const showPromptDialog = ref(false)
const promptText = ref('')//循环改写提示词

const handlePromptSave = () => {
  // 保存提示词逻辑
  showPromptDialog.value = false
  // 发送提示词到后端
  setModulesKv('prompt_rewrite', {"data":promptText.value})
}

const getPromptRewrite =async  () => {
  let  kv =await getModulesKv("prompt_rewrite")
  if (kv) {
    promptText.value = kv.data
  }
} 


// 新增状态
const selectedTimezone = ref('')
const selectedLanguage = ref('')
const selectedLanguageLabel = ref('')
const selectedAudioDriver = ref('')
const selectedSpeechModel = ref('')
const handleLanguageChange = async (value) => {
  try {
    const response = await fetch('http://127.0.0.1:7073/set_config', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ 
        timeZone: selectedTimezone.value,
        language:selectedLanguage.value,
        soundCard:selectedAudioDriver.value
       })
    })
    const data = await response.json()
    console.log('语言设置成功:', data)
  } catch (error) {
    console.error('语言设置失败:', error)
  }
}

const initializeSpeechModel = async () => {
  try {
   const response = await fetch('http://127.0.0.1:7073/initialize', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        model:selectedSpeechModel.value
       })
    })
    const data = await response.json()
    console.log('语音模型初始化成功:', data)
   }
   catch (error) {
    console.error('语音模型初始化失败:', error) 
  }
}


const timezoneOptions = [
  { label: '国际日期变更线以西', value: 'Etc/GMT+12' },
  { label: '中途岛', value: 'Pacific/Midway' },
  { label: '夏威夷标准时间', value: 'Pacific/Honolulu' },
  { label: '阿拉斯加标准时间', value: 'America/Anchorage' },
  { label: '太平洋标准时间（美国西海岸）', value: 'America/Los_Angeles' },
  { label: '山地标准时间（美国）', value: 'America/Denver' },
  { label: '中部标准时间（美国）', value: 'America/Chicago' },
  { label: '东部标准时间（美国东海岸）', value: 'America/New_York' },
  { label: '智利标准时间', value: 'America/Santiago' },
  { label: '巴西标准时间', value: 'America/Sao_Paulo' },
  { label: '南乔治亚岛', value: 'Atlantic/South_Georgia' },
  { label: '亚速尔群岛', value: 'Atlantic/Azores' },
  { label: '协调世界时（伦敦）', value: 'UTC' },
  { label: '中欧时间（巴黎）', value: 'Europe/Paris' },
  { label: '东欧时间（雅典）', value: 'Europe/Athens' },
  { label: '莫斯科时间', value: 'Europe/Moscow' },
  { label: '迪拜时间', value: 'Asia/Dubai' },
  { label: '巴基斯坦标准时间', value: 'Asia/Karachi' },
  { label: '印度标准时间', value: 'Asia/Kolkata' },
  { label: '孟加拉国时间', value: 'Asia/Dhaka' },
  { label: '泰国时间', value: 'Asia/Bangkok' },
  { label: '中国标准时间', value: 'Asia/Shanghai' },
  { label: '日本标准时间', value: 'Asia/Tokyo' },
  { label: '澳大利亚东部标准时间', value: 'Australia/Sydney' },
  { label: '所罗门群岛', value: 'Pacific/Guadalcanal' },
  { label: '新西兰标准时间', value: 'Pacific/Auckland' }
]

const languageOptions = [
{ label: '英语', value: 'en' },
  { label: '西班牙语', value: 'es' },
  { label: '法语', value: 'fr' },
  { label: '德语', value: 'de' },
  { label: '意大利语', value: 'it' },
  { label: '葡萄牙语', value: 'pt' },
  { label: '波兰语', value: 'pl' },
  { label: '土耳其语', value: 'tr' },
  { label: '俄语', value: 'ru' },
  { label: '荷兰语', value: 'nl' },
  { label: '捷克语', value: 'cs' },
  { label: '阿拉伯语', value: 'ar' },
  { label: '简体中文', value: 'zh-cn' },
  { label: '日语', value: 'ja' },
  { label: '匈牙利语', value: 'hu' },
  { label: '韩语', value: 'ko' }
]

const audioDriverOptions = [
  { label: '默认驱动', value: 'default' },
  { label: 'ASIO', value: 'asio' },
  { label: 'WASAPI', value: 'wasapi' }
]

const socket = ref<any>(null) 
let reconnectAttempts = 0
const MAX_RECONNECT_ATTEMPTS = Infinity // 设置为无限重连
const modules = ref([])

onUnmounted(() => {
  if (socket.value) {
    socket.value.disconnect()
    socket.value = null
  }
})

const messages = ref<Array<{content: string,data_types: string , id?: string}>>([])
  const filterOptions = ref([
  { label: '弹幕消息', value: 'CommentEvent', checked: true },
  { label: '礼物消息', value: 'GiftEvent', checked: true },
  { label: '进入直播间', value: 'JoinEvent', checked: true },
  { label: '分享', value: 'ShareEvent', checked: true },
  { label: '关注', value: 'FollowEvent', checked: true },
  { label: '点赞', value: 'LikeEvent', checked: true },
])
const selectedFilters = ref(['CommentEvent', 'GiftEvent','JoinEvent','ShareEvent',"FollowEvent","LikeEvent"]) // 默认选中的过滤项

function addMessage(message: {content: string , data_types: string , id?: string}) {
  const newMsg = {
    ...message,
    data_types: message.data_types,
    id: Date.now().toString() // 为每条消息添加唯一ID
  }
  messages.value = [...messages.value.slice(-99), newMsg]
}
const getMessageColor = (type) => {
  switch(type) {
    case 'CommentEvent': return '#4A5568'
    case 'GiftEvent': return '#2E7D32'   
    case 'JoinEvent': return '#1565C0'  
    case 'ShareEvent': return '#D84315'   
    case 'FollowEvent': return '#6A1B9A' 
    case 'LikeEvent': return '#C62828'  
    default: return '#546E7A'       
  }
}

const handleVolumeChange = (block) => {
  console.log("block----------------------------");
  console.log(block);
  setModulesKv(String(block.id),block)
}

const mames = ref([])
const getMames =async  () => {
  const response = await fetch('http://127.0.0.1:7072/names', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        mames.value = await response.json();
}

const getModulesKv =async  (id) => {
  const response = await fetch('http://127.0.0.1:7072/get_kv', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ key: String(id) }),
        });
        const data = await response.json();
        return data.kv;
}

// 新增：音频设备下拉选项
const audioDeviceOptions = ref([])

const getAudioDevices = async () => {
  try {
    const response = await fetch('http://127.0.0.1:7073/get_sound_cards', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })
    const data = await response.json()
    audioDeviceOptions.value = data.map(item => ({
      label: item,
      value: item
    }))
    // 检查是否有"CABLE Input"，有则设为默认值
    const cable = audioDeviceOptions.value.find(opt => opt.label === 'CABLE Output (VB-Audio Virtual Cable)')
    if (cable) {
      selectedAudioDriver.value = cable.value
    }
  } catch (error) {
    console.error('获取音频设备失败:', error)
  }
}

const speedModelOptions = ref([])

const getSpeedModels = async () => {
  try {
    const response = await fetch('http://127.0.0.1:7073/get_model_files', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      } 
    })
    const data = await response.json()
    speedModelOptions.value = (data.files || []).map(item => ({
      label: item,
      value: item
    }))
    
  }  catch (error) {
    console.error('获取语音模型失败:', error)
  }
}

const setModulesKv =async  (id,kv) => {
  const response = await fetch('http://127.0.0.1:7072/set_kv', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ key: String(id) ,kv:kv}),
        });
}

const fetchModules = async () => {
  try {
    console.log('正在获取模块数据...')
    let  newList   =[]
    const response = await fetch('http://127.0.0.1:7072/base-modules')
    let list = await response.json()
    console.log("lis23132t:",list);
    for (const  [index, module] of list.entries()) {
      console.log("module.id:",module.id);
      let  kv =await getModulesKv(module.id)
      console.log("kv:",kv);
      module.volume = kv.volume || 50
      module.isActive= kv.isActive || false
      module.retAi= kv.retAi || false
      module.speed = kv.speed || 1.0
      module.selectedNameId = kv.selectedNameId || ""
      setModulesKv(module.id,module)
      newList.push(module)
    }
    modules.value=newList
  } catch (error) {
    console.error('获取模块数据失败:', error)
  }
}

onMounted(() => {
  fetchModules()
  getMames()
  getPromptRewrite()
  getAudioDevices()
  getHumanVoiceFiles()
  getSpeedModels()
  // setInterval(fetchModules, 25000)
  return
  if (socket.value && socket.value.connected) return
    //  { 
    //       's': '渲染完整数据',
    //       'u': 相关用户名,
    //       'c': 弹幕内容,
    //       'g': 礼物名称,
    //       'gn': 礼物数量,
    //  }
  socket.value = io('ws://127.0.0.1:7073', {
    reconnection: true,
    reconnectionAttempts: MAX_RECONNECT_ATTEMPTS,
    reconnectionDelay: 5,
    randomizationFactor: 5
  })
    socket.value.on('connect', () => {
      reconnectAttempts = 0
      console.log('Socket.IO连接成功')
    })
    socket.value.on('reconnect_attempt', (attemptNumber) => {
    reconnectAttempts = attemptNumber
    console.log(`正在尝试第${attemptNumber}次重连...`)
    })
    socket.value.on('connect_error', (err) => {
      console.error('Socket.IO连接错误:', err)
    })
    socket.value.on('ConnectEvent', function(message) {
      if (!selectedFilters.value.includes("ConnectEvent")){
        return;
      }
              console.log("ConnectEvent:",message);
              addMessage({content: message.s ,data_types: "ConnectEvent"})
    });
    socket.value.on('CommentEvent', function(message) {//弹幕消息
      EnterBarrageUserName.value = message.u //弹幕用户名
      EnterBarrageContent.value = message.c //弹幕内容
      if (!selectedFilters.value.includes("CommentEvent")){
        return;
      }
              console.log("CommentEvent:",message);
              addMessage({content: message.s,data_types: "CommentEvent"})
    });
    socket.value.on('GiftEvent', function(message) { //礼物消息
      EnterGiftUserName.value = message.u //送礼物的用户名
      EnterGiftGoodsName.value = message.g // 礼物名称
      EnterGiftNum.value = message.gn //礼物数量
      if (!selectedFilters.value.includes("GiftEvent")){
        return;
      }
              console.log("GiftEvent:",message);
              addMessage({content: message.s,data_types: "GiftEvent"})
    });
    socket.value.on('JoinEvent', function(message) {// 进入直播间
      EnterLiveRoomUserName.value = message.u
      if (!selectedFilters.value.includes("JoinEvent")){
        return;
      }
              console.log("JoinEvent:",message);
              addMessage({content: message.s,data_types: "JoinEvent"})
    });
    socket.value.on('ShareEvent', function(message) { //分享
      if (!selectedFilters.value.includes("ShareEvent")){
        return;
      }
              console.log("ShareEvent:",message);
              addMessage({content: message.s,data_types: "ShareEvent"})
    });
    socket.value.on('FollowEvent', function(message) { //关注
      if (!selectedFilters.value.includes("FollowEvent")){
        return;
      }
              console.log("FollowEvent:",message);
              addMessage({content: message.s,data_types: "FollowEvent"})
    });
    socket.value.on('LikeEvent', function(message) { //点赞
      EnterSupportRoomUserName.value = message.u
      if (!selectedFilters.value.includes("LikeEvent")){
        return;
      }
              console.log("LikeEvent:",message);
              addMessage({content: message.s,data_types: "LikeEvent"})
    });
    // socket.value.on('LiveEndEvent', function(message) { //直播结束
    //           console.log("LiveEndEvent:",message);
    //           addMessage({content: message})
    // });
    // socket.value.on('DisconnectEvent', function(message) { //断开连接
    //           console.log("DisconnectEvent:",message);
    //           addMessage({content: message})
    // });
})

const toggleFilter = (value) => {
  const index = selectedFilters.value.indexOf(value)
  if (index > -1) {
    selectedFilters.value.splice(index, 1)
  } else {
    selectedFilters.value.push(value)
  }
}
function toggleLive() {
  loading.value = true
  setTimeout(() => {
    isLive.value = !isLive.value
    loading.value = false
  }, 1000)
  socket.value.emit('sync_tk', { id: streamUrl.value });
}

function outToggleLive() {
  loading.value = true
  setTimeout(() => {
    isLive.value = !isLive.value
    loading.value = false
  }, 1000)
  socket.value.emit('sync_tk', { id: "" });
}



interface PlayItem {
    content: any;
    filename: string;
    play_mode: string;
  }
const playList = ref<PlayItem[]>([])//播放列表
const uesPlayList = ref<PlayItem[]>([])//已播放列表
const start = ref(false)//启动开关

//注册模块，
const registerModules = () => {
    //循环模块列表 干活
    start.value = true
    for (const module of modules.value) {
        // api\control-go\model\base_module.go
//     TriggerExecuteLoop    = "ExecuteLoop"    // 执行循环：是否自动循环执行任务
//     TriggerIntervalLoop   = "IntervalLoop"   // 间隔循环：按照指定间隔时间循环执行
//     TriggerBarrageComment = "BarrageComment" // 弹幕评论：是否自动发送弹幕评论
//     TriggerSendGift       = "SendGift"       // 送礼物：是否自动发送礼物
//     TriggerLike           = "Like"           // 点赞：是否自动点赞
//     TriggerEnterLiveRoom  = "EnterLiveRoom"  // 进入直播间：是否自动进入直播间
//     TriggerWarningTip     = "WarningTip"     // 警告提示：是否显示警告提示
// )
        //循环读取
        if (module.trigger_conditions.includes("ExecuteLoop")) { // 执行循环：是否自动循环执行任务
            ExecuteLoop(module)
        }
        if (module.trigger_conditions.includes("IntervalLoop")) { //间隔循环：按照指定间隔时间循环执行
            ExecuteLoop(module)
        }
        //弹幕评论 BarrageComment
        if (module.trigger_conditions.includes("BarrageComment")) {// 弹幕
            BarrageComment(module)
        }
        //送礼物 SendGift
        if (module.trigger_conditions.includes("SendGift")) {
          SendGift(module)
        }

        //点赞 Like
        if (module.trigger_conditions.includes("Like")) {
          EnterLiveRoom(module)
        }
        //进入直播间 EnterLiveRoom
        if (module.trigger_conditions.includes("EnterLiveRoom")) {
          EnterLiveRoom(module)
        }
    }
    playListConsumption()
}
const intervalTime = ref(1000) // 默认1000毫秒 1秒=1000毫秒
//消费playList
const playListConsumption= async () => {
    do {
        if (playList.value.length > 0) { //队列消费
          await new Promise(resolve => setTimeout(resolve, intervalTime.value))   // 等待1000毫秒
          console.log("kai:",playList.value);
            const item = playList.value.shift() //出队
            console.log("jie:",playList.value);
            if (item) {
                uesPlayList.value.push(item) //已播放列表
                if (uesPlayList.value.length > 100) {
                    uesPlayList.value.shift() //删除第一条
                }
                let data = await play_task_voice_api(item.filename, item.play_mode); //播放
                // await new Promise(resolve => setTimeout(resolve, data.duration * 1000)) //等待
            }
        }
        await new Promise(resolve => setTimeout(resolve, 100))   
    } while (start.value);
}

//循环模块处理
const ExecuteLoop= async (module) => {
    let index = 0 //当前播放索引
    const script_content_len = module.script_content.length //脚本长度
    do {
        if (!module.isActive) { 
          await new Promise(resolve => setTimeout(resolve, 1000))
          continue
        }

        if (module.read_step == "random") { //随机
            if (index == 0) { 
                module.script_content = shuffleArray(module.script_content)//将数组 打乱
            }
        }
        if (module.interval_time_start != 0 && module.interval_time_end!= 0) { //间隔时间
            const randomTime = Math.floor(Math.random() * (module.interval_time_end - module.interval_time_start + 1)) + module.interval_time_start; //随机时间
            await new Promise(resolve => setTimeout(resolve, randomTime * 1000)) //等待
        }

        let content = module.script_content[index] 
        while (playList.value.length >= 2) { //等待队列 消费完探入
            await new Promise(resolve => setTimeout(resolve, 1000))
        }
        // 是否改写
        if (module.retAi) { //是否改写
          if (model_api.value == "") { //是否改写
            console.log("请选择模型");
            await new Promise(resolve => setTimeout(resolve, 4000))
            continue
          }
          let prompt =  ReplaceText(promptText.value) //提示词 赋值变量
          content = ReplaceText(content) //内容赋值变量
          const  apidata =  await DisposableSendApi(
            model_api.value,
            parameters_api.value,
            content, //文本
            prompt, //提示词
            supplierName_api.value, //供应商名称
          )
          console.log("apidata:",apidata);
          content = apidata
        }

        const newFileName =   crypto.randomUUID()+ "_" + Date.now() + '.wav' //生成文件名
        let ok = await generate_wav_api(
            content, //文本
            selectedLanguage.value || "en",  //语种
            newFileName, //生成文件名
            module.selectedNameId, //音色文件名
            module.speed, //生成速度
            module.volume / 100 //生成音量
          ) //生成音量
        if (!ok) { //生成失败
            await new Promise(resolve => setTimeout(resolve, 1000))
            continue
        }
        playList.value.push({ //入队
            content: content, //内容
            filename: newFileName, //文件名
            play_mode: "serial", //播放模式
        })    
        if (index == script_content_len - 1) { //循环
            index = 0
        } else {
            index = index + 1
        }
    } while (start.value);
}

// 进入直播间 模块
const EnterLiveRoomUserName = ref("")  //进入直播间用户名

const EnterLiveRoom= async (module) => {
    let index = 0 //当前播放索引
    const script_content_len = module.script_content.length //脚本长度
    do {
        if (!module.isActive) { 
          await new Promise(resolve => setTimeout(resolve, 1000))
          continue
        }
        if (module.read_step == "random") { //随机
            if (index == 0) { 
                module.script_content = shuffleArray(module.script_content)//将数组 打乱
            }
        }
        if (module.interval_time_start != 0 && module.interval_time_end!= 0) { //间隔时间
            const randomTime = Math.floor(Math.random() * (module.interval_time_end - module.interval_time_start + 1)) + module.interval_time_start; //随机时间
            await new Promise(resolve => setTimeout(resolve, randomTime * 1000)) //等待
            console.log("4");
        }
        while (EnterLiveRoomUserName.value == "") { //等待有用户进入直播间
            await new Promise(resolve => setTimeout(resolve, 1000))
        }
        let content = module.script_content[index]  //欢迎{用户名称}的到来，您你的到来让直播间蓬革
        const newFileName =   crypto.randomUUID()+ "_" + Date.now() + '.wav' //生成文件名
        let ok = await generate_wav_api(
            content, //文本
            selectedLanguage.value || "en",  //语种
            newFileName, //生成文件名
            module.selectedNameId, //音色文件名
            module.speed, //生成速度
            module.volume / 100 //生成音量
          ) //生成音量
        if (!ok) { //生成失败
            await new Promise(resolve => setTimeout(resolve, 1000))
            continue
        }
        // 插入到头部
        playList.value.unshift({ //入队
            content: content, //内容
            filename: newFileName, //文件名
            play_mode: "top", //播放模式
        })    
        EnterLiveRoomUserName.value = "" //清空
        if (index == script_content_len - 1) { //循环
            index = 0
        } else {
            index = index + 1
        }
    } while (start.value);
}

// 点赞 Like
const EnterSupportRoomUserName = ref("")  //点赞直播间 用户名
const Like= async (module) => {
    let index = 0 //当前播放索引
    const script_content_len = module.script_content.length //脚本长度
    do {
        if (!module.isActive) { 
          await new Promise(resolve => setTimeout(resolve, 1000))
          continue
        }
        if (module.read_step == "random") { //随机
            if (index == 0) { 
                module.script_content = shuffleArray(module.script_content)//将数组 打乱
            }
        }
        if (module.interval_time_start != 0 && module.interval_time_end!= 0) { //间隔时间
            const randomTime = Math.floor(Math.random() * (module.interval_time_end - module.interval_time_start + 1)) + module.interval_time_start; //随机时间
            await new Promise(resolve => setTimeout(resolve, randomTime * 1000)) //等待
        }
        while (EnterSupportRoomUserName.value == "") { //等待有用户进入直播间
            await new Promise(resolve => setTimeout(resolve, 1000))
        }
        let content = module.script_content[index] 
        const newFileName =   crypto.randomUUID()+ "_" + Date.now() + '.wav' //生成文件名
        let ok = await generate_wav_api(
            content, //文本
            selectedLanguage.value || "en",  //语种
            newFileName, //生成文件名
            module.selectedNameId, //音色文件名
            module.speed, //生成速度
            module.volume / 100 //生成音量
          ) //生成音量
        if (!ok) { //生成失败
            await new Promise(resolve => setTimeout(resolve, 1000))
            continue
        }
        // 插入到头部
        playList.value.unshift({ //入队
            content: content, //内容
            filename: newFileName, //文件名
            play_mode: "serial", //播放模式
        })    
        EnterSupportRoomUserName.value = "" //清空
        if (index == script_content_len - 1) { //循环
            index = 0
        } else {
            index = index + 1
        }
    } while (start.value);
}

//送礼物数据
const EnterGiftUserName = ref("")  //送礼物 用户名
const EnterGiftNum = ref("")  //送礼物数量
const EnterGiftGoodsName = ref("")  //送的礼物名字
const SendGift= async (module) => {
    let index = 0 //当前播放索引
    const script_content_len = module.script_content.length //脚本长度
    do {
        if (!module.isActive) { 
          await new Promise(resolve => setTimeout(resolve, 1000))
          continue
        }
        if (module.read_step == "random") { //随机
            if (index == 0) { 
                module.script_content = shuffleArray(module.script_content)//将数组 打乱
            }
        }
        if (module.interval_time_start != 0 && module.interval_time_end!= 0) { //间隔时间
            const randomTime = Math.floor(Math.random() * (module.interval_time_end - module.interval_time_start + 1)) + module.interval_time_start; //随机时间
            await new Promise(resolve => setTimeout(resolve, randomTime * 1000)) //等待
        }
        while (EnterGiftUserName.value == "") { //等待有用户送礼物
            await new Promise(resolve => setTimeout(resolve, 1000))
        }
        let content = module.script_content[index] 
        const newFileName =   crypto.randomUUID()+ "_" + Date.now() + '.wav' //生成文件名
        let ok = await generate_wav_api(
            content, //文本
            selectedLanguage.value || "en",  //语种
            newFileName, //生成文件名
            module.selectedNameId, //音色文件名
            module.speed, //生成速度
            module.volume / 100 //生成音量
          ) //生成音量
        if (!ok) { //生成失败
            await new Promise(resolve => setTimeout(resolve, 1000))
            continue
        }
        // 插入到头部
        playList.value.unshift({ //入队
            content: content, //内容
            filename: newFileName, //文件名
            play_mode: "serial", //播放模式
        })    
        EnterGiftUserName.value = "" //清空
        EnterGiftNum.value = "" //清空
        EnterGiftGoodsName.value = "" //清空
        if (index == script_content_len - 1) { //循环
            index = 0
        } else {
            index = index + 1
        }
    } while (start.value);
}
//弹幕评论数据
const EnterBarrageUserName = ref("")  //弹幕评论 用户名
const EnterBarrageContent = ref("")  //弹幕评论 内容
const BarrageComment= async (module) => {
    // let index = 0 //当前播放索引
    // const script_content_len = module.script_content.length //脚本长度
    do {
        if (!module.isActive) { 
          await new Promise(resolve => setTimeout(resolve, 2000))
          continue
        }
        //弹幕间隔时间
        if (module.interval_time_start != 0 && module.interval_time_end!= 0) { //间隔时间
            const randomTime = Math.floor(Math.random() * (module.interval_time_end - module.interval_time_start + 1)) + module.interval_time_start; //随机时间
            await new Promise(resolve => setTimeout(resolve, randomTime * 1000)) //等待
        }
        while (EnterBarrageContent.value == "") { //弹幕内容
            await new Promise(resolve => setTimeout(resolve, 1000))
        }
        if (model_api.value == "") { //是否改写
            console.log("请选择模型");
            await new Promise(resolve => setTimeout(resolve, 4000))
            continue
        }
        let content = "" // 
        let prompt =  ReplaceText(module.script_content[0]) //提示词
          let ai_user_content =  ReplaceText(EnterBarrageContent.value) //弹幕内容作为ai输入
            const  apidata =  await DisposableSendApi(
              model_api.value,
              parameters_api.value,
              ai_user_content, //文本
              prompt, //提示词
              supplierName_api.value, //供应商名称
            )
            console.log("apidata:",apidata);
            content = apidata
        const newFileName =   crypto.randomUUID()+ "_" + Date.now() + '.wav' //生成文件名
        let ok = await generate_wav_api(
            content, //文本
            selectedLanguage.value || "en",  //语种
            newFileName, //生成文件名
            module.selectedNameId, //音色文件名
            module.speed, //生成速度
            module.volume / 100 //生成音量
          ) //生成音量
        if (!ok) { //生成失败
            await new Promise(resolve => setTimeout(resolve, 1000))
            continue
        }
        // 插入到头部
        playList.value.unshift({ //入队
            content: content, //内容
            filename: newFileName, //文件名
            play_mode: "serial", //播放模式
        })    
        EnterBarrageUserName.value = "" //清空
        EnterBarrageContent.value = "" //清空
    } while (start.value);
}
// 更标准的Fisher-Yates洗牌算法
function shuffleArray(array) {
    for (let i = array.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [array[i], array[j]] = [array[j], array[i]];
    }
    return array;
}
// let generate_wav_api_runing = false // 运行在
let generateLock = Promise.resolve(); // 初始化为已解决的Promise

// http://127.0.0.1:7073/generate_wav
// 生成wav文件
const generate_wav_api = async (_text:string,
    _language:string,
    _filename:string,
    _speaker_wav:string,
    _speed: number,
    _volume: number) => {
    const myLock = generateLock; // 获取当前锁状态
    let releaseLock;
    generateLock = new Promise(resolve => releaseLock = resolve); // 创建新锁
    
    await myLock; // 等待之前的锁释放
    try {
    // generate_wav_api_runing = true;
        const response = await fetch('http://127.0.0.1:7073/generate_wav', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                text: _text,
                language: _language,
                filename : _filename,
                speaker_wav: _speaker_wav,
                speed: _speed,// 1,
                volume: _volume// 0.5
            }),
        });
      // generate_wav_api_runing = false;
      return response.ok;
    }finally {
        releaseLock(); // 释放锁
    }
    
}

// 播放任务
const play_task_voice_api = async (_filename:string,play_mode:string) => {
    const response = await fetch('http://127.0.0.1:7073/play_task_voice', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
            filename : _filename,
            play_mode : play_mode
         }),
    });
    await response.json();
}

const humanVoiceOptions = ref([])

const getHumanVoiceFiles = async () => {
  try {
    const response = await fetch('http://127.0.0.1:7073/get_human_voice_files', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })
    const data = await response.json()
    // 转换为下拉选项格式
    humanVoiceOptions.value = (data.files || []).map(name => ({
      label: name,
      value: name
    }))
  } catch (error) {
    console.error('获取人声音色文件失败:', error)
  }
}


// 模型调用API
const DisposableSendApi = async (model,parameters,user_content,system_prompt,supplier_name) => {
  const response = await fetch('http://127.0.0.1:7072/disposable_send', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
              model: model, 
              parameters: parameters,
              user_content: user_content,
              system_prompt: system_prompt,
              supplier_name:supplier_name,
              }),
        });
        console.log(response);
        
        const data = await response.text();
        if (data.includes('</think>')) {
          // 去除 <think>
          const endIndex = data.indexOf('</think>');
          // 截取endIndex 过后的字符串
          const result = data.substring(endIndex + 7); // 跳过 </think> 标签
          const result2 = result.replace(/\n/g, '');
          return result2;
        }
        return data;
}

// 替换文本 替换变量
const ReplaceText= (text) => {
  let newText = text
  if (text.includes('{语种}')){ //是否包含
    newText = newText.replace('{语种}', selectedLanguageLabel.value) //替换
  }
  return newText;
} 

</script>

<style scoped lang="scss">

.gd-bg{
  background-color: v-bind(themeThinkBg);
}

.live-console {
  // height: calc(100% - 130px);
  height: calc(100vh - 20px); padding-bottom: 20px;
  .three-column-layout {
    padding: 0px  12px;
    display: grid;
    // grid-template-columns: repeat(3, 1fr);
    grid-template-columns: 1fr 3fr;  // 中间列宽度是两侧的1.5倍
    gap: 8px;
    // height: 50%;
    
    .n-card {
      // padding: 12px 0px;
      // height: 100%;
      display: flex;
      flex-direction: column;
      &-header {
        flex-shrink: 0;
        padding: 8px 12px;  
      }
      
      &-content {
        flex: 1;
        overflow: hidden;
        display: flex;
        flex-direction: column;
      }
    }
  }
  
  .control-panel {
    grid-column: 1 / 3;
  }
  
  .status-row {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 10px;
  }
  
  .settings-row {
    margin-top: 10px;
  }
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }
  
  .message-item {
    padding: 3px 0;
    // border-bottom: 1px solid #f0f0f0;
    
    .username {
      font-weight: bold;
      margin-right: 8px;
      color: var(--bt-theme-color);
    }
  }
  
  .message-input {
    padding: 5px;
    background: #f8fafc;
    border-radius: 8px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    
    .n-input-group {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .n-input {
        flex: 1;
        min-width: 0;
      }
      
      .n-button {
        flex-shrink: 0;
        margin: 0;
      }
    }
    
    .mode-switch {
      width: 40px;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    
    .voice-button {
      flex: 1;
      padding: 0 16px;
      min-width: 0;
      
      span {
        margin-left: 8px;
        white-space: nowrap;
      }
    }
  }
}

.control-blocks {
  // display: grid;
  // gap: 10px;  // 减小间距
}

.control-block {
  padding: 8px;  // 减小内边距
  border-radius: 4px;
  
  .block-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 4px;  // 减小间距
    
    h4 {
      font-size: 14px;  // 减小标题字号
      margin: 0;
    }
  }
  
  .block-content {
    padding-top: 4px;  // 减小间距
  }
}

.audio-blocks {
  display: grid;
  gap: 12px;
}

.audio-block {
  padding: 0px;
  // border: 1px solid #f0f0f0;
  border-radius: 4px;
  
  .audio-info {
    margin-bottom: 1px;
    
    h4 {
      margin: 0 0 4px 0;
      font-size: 14px;
    }
    
    p {
      margin: 0;
      font-size: 12px;
      color: #666;
    }
  }
  
  audio {
    width: 100%;
    height: 24px;
  }
}
.settings-bar {
  display: flex;
  flex-wrap: wrap;  // 允许换行
  align-items: center;
  gap: 10px;  // 统一间距
  padding: 12px;
  background: var(--n-color-embedded);
  border-radius: var(--n-border-radius);
  
  .n-input-group {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

</style>
