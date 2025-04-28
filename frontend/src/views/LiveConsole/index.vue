<template>
  <div class="live-console">
    <!-- 设置区域 -->
    <div class="settings-bar">
      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">时区设置</n-text>
        <n-select 
          v-model:value="selectedTimezone"
          :options="timezoneOptions"
          placeholder="选择时区"
          style="width: 250px;"
        />
      </div>
      
      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">语言设置</n-text>
        <n-select 
          v-model:value="selectedLanguage"
          :options="languageOptions"
          placeholder="选择语种"
          style="width: 150px;"
        />
      </div>
      
      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">音频驱动</n-text>
        <n-select 
          v-model:value="selectedAudioDriver"
          :options="audioDriverOptions"
          placeholder="选择驱动"
          style="width: 100px;"
        />
      </div>
      
      <!-- 新增直播链接输入和启动按钮 -->
      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">直播弹幕接入</n-text>
        <n-input 
          v-model:value="streamUrl"
          placeholder="输入直播间链接"
        >
          <template #suffix>
            <n-button 
              type="primary" 
              @click="toggleLive"
              :loading="loading"
              text
            >
              {{ isLive ? '停止接入' : '接入弹幕' }}
            </n-button>
          </template>
        </n-input>
      </div>
       <!-- 新增直播链接输入和启动按钮 -->
       <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">确认启动</n-text>
        <n-input-group>
          <!-- <n-input 
            v-model:value="streamUrl"
            placeholder="输入RTMP直播链接"
            style="width: 300px;"
          /> -->
          <n-button 
            type="primary" 
            @click="toggleLive"
            :loading="loading"
          >
            {{ isLive ? '停止音频直播' : '启动音频直播' }}
          </n-button>
        </n-input-group>
      </div>
    </div>

    <!-- 三列布局 -->
    <div class="three-column-layout">
      <n-card title="模块">
        <n-infinite-scroll style="height: calc(100vh - 260px);padding-right: 15px;" :distance="10">
          <div class="control-blocks">
          <div class="control-block" v-for="(block, index) in controlBlocks" :key="index" style=" background-color: #f8fafc;">
            <div class="block-header">
              <h4>{{ block.title }}</h4>
              <div style="display: flex; align-items: center; gap: 8px;">
              <!-- <n-select 
                v-model:value="block.selectedVoice"
                :options="block.voiceOptions.map(opt => ({ label: opt, value: opt }))"
                style="width: 90px;"
                placeholder="选择人声"
                clearable
                 size="small"
              /> -->
              
              <n-switch v-model:value="block.isActive"  size="small" />
              </div>
              
            </div>
            <div style="height: 5px;"></div>
            <div style="display: flex; align-items: center; gap: 8px; width: 100%;">
              <div style="width: 30%;font-size: 12px;">
                选择人声
            </div>
            <n-select 
                v-model:value="block.selectedVoice"
                :options="block.voiceOptions.map(opt => ({ label: opt, value: opt }))"
                style="width: 70%;"
                placeholder="选择人声"
                clearable
                 size="small"
              />
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
              />
            </div>  
              
            <!-- <div class="block-content" v-if="block.isActive"> -->
              <!-- <n-select 
                v-model:value="block.selectedVoice"
                :options="block.voiceOptions.map(opt => ({ label: opt, value: opt }))"
                style="margin-bottom: 8px;"
                placeholder="选择人声"
                clearable
              /> -->
              <!-- 这里放置每个功能方块的具体内容 -->
              <!-- <n-input v-if="block.type === 'input'" v-model:value="block.value" :placeholder="block.placeholder" /> -->
              <!-- 其他类型的功能组件 -->
            <!-- </div> -->
          </div>
        </div>
        </n-infinite-scroll>
      </n-card>

      <!-- 观众互动 -->
      <n-card title="观众互动">
        <template #header-extra>
          <n-switch 
            v-model:value="autoReadMode" 
          >
            <template #checked>自动朗读</template>
            <template #unchecked>关闭朗读</template>
          </n-switch>
        </template>
        <n-infinite-scroll style="height: calc(100vh - 260px);padding-right: 15px;" :distance="10">
          <div 
            v-for="(msg, index) in messages" 
            :key="index" 
            class="message-item"
            style="white-space: nowrap; overflow: hidden; text-overflow: ellipsis; display: flex; align-items: center; gap: 4px;"
          >
          <n-button 
            text
            @click=""
            style="padding: 0; min-width: 20px; display: flex; align-items: center;"
          >
            <template #icon>
              <i class="i-tdesign:play-circle w-15 h-15" style="vertical-align: middle;color: #888;"></i>
            </template>
          </n-button>
          <span class="username" style="line-height: 1;">{{ msg.user }}:</span>
          <span class="content" style="line-height: 1;">{{ msg.content }}</span>
        </div>
        </n-infinite-scroll>
        <div class="message-input">
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
        </div>
      </n-card>

      <!-- 直播数据 -->
      <n-card title="音频数据">
        <n-infinite-scroll style="height: calc(100vh - 260px);padding-right: 15px;" :distance="10">
          <div class="audio-blocks">
          <div class="audio-block" v-for="(audio, index) in audioFiles" :key="index" style="background-color: #f8fafc;">
            <div class="audio-info">
              <div style="display: flex; justify-content: space-between; align-items: center;">
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
              <p>{{ audio.description }}</p>
            </div>
            <audio controls :src="audio.url"></audio>
          </div>
        </div>
        </n-infinite-scroll>   
      </n-card>
        <!-- <div class="stats-grid">
          <n-statistic label="观看人数" :value="viewerCount" />
          <n-statistic label="点赞数" :value="likeCount" />
          <n-statistic label="评论数" :value="commentCount" />
          <n-statistic label="直播时长" :value="duration" />
        </div>
      </n-card> -->
    </div>

    <!-- 互动消息区域 -->
    
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
const autoReadMode = ref(false)
const controlBlocks = ref([
  {
    title: '推流控制',
    type: 'switch',
    isActive: true,
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认',
    volume: 50, // 默认音量
    speed: 1.0 // 默认音速
  },
  {
    title: '画质调节',
    type: 'select',
    isActive: false,
    options: ['高清', '标清', '流畅'],
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认',
    volume: 50, // 默认音量
    speed: 1.0 // 默认音速
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认',
    volume: 50, // 默认音量
    speed: 1.0 // 默认音速
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  },
  {
    title: '直播公告',
    type: 'input',
    isActive: false,
    placeholder: '输入直播公告内容',
    voiceOptions: ['默认', '男声', '女声', '卡通'],
    selectedVoice: '默认'
  }
])
const isLive = ref(false)
const loading = ref(false)
const streamUrl = ref('rtmp://your-stream-url.com/live')
const viewerCount = ref(0)
const likeCount = ref(0)
const commentCount = ref(0)
const duration = ref('00:00:00')
const messages = ref([
  { user: '用户1', content: '直播内容很棒！' },
  { user: '用户2', content: '什么时候开始？' }
])
const newMessage = ref('')
const isVoiceMode = ref(false)
const audioFiles = ref([
  {
    title: '背景音乐1',
    description: '轻松愉轻松愉快的背景音乐轻松愉快的背景音乐轻松愉快的背景音乐轻松愉快的背景音乐轻松愉快的背景音乐轻松愉快的背景音乐快的背景音乐',
    url: '/audio/bgm1.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  },
  {
    title: '开场音乐',
    description: '直播开场专用音乐',
    url: '/audio/opening.mp3'
  }
])
function toggleLive() {
  loading.value = true
  setTimeout(() => {
    isLive.value = !isLive.value
    loading.value = false
  }, 1000)
}

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

// 新增状态
const selectedTimezone = ref('')
const selectedLanguage = ref('')
const selectedAudioDriver = ref('')

const timezoneOptions = [
  { label: 'GMT-12 国际日期变更线西', value: 'GMT-12' },
  { label: 'GMT-11 萨摩亚标准时间', value: 'GMT-11' },
  { label: 'GMT-10 夏威夷-阿留申标准时间', value: 'GMT-10' },
  { label: 'GMT-9 阿拉斯加标准时间', value: 'GMT-9' },
  { label: 'GMT-8 太平洋标准时间', value: 'GMT-8' },
  { label: 'GMT-7 山地标准时间', value: 'GMT-7' },
  { label: 'GMT-6 中部标准时间', value: 'GMT-6' },
  { label: 'GMT-5 东部标准时间', value: 'GMT-5' },
  { label: 'GMT-4 大西洋标准时间', value: 'GMT-4' },
  { label: 'GMT-3 巴西利亚标准时间', value: 'GMT-3' },
  { label: 'GMT-2 南乔治亚和南桑威奇群岛时间', value: 'GMT-2' },
  { label: 'GMT-1 亚速尔群岛标准时间', value: 'GMT-1' },
  { label: 'GMT+0 格林威治标准时间', value: 'GMT+0' },
  { label: 'GMT+1 中欧标准时间', value: 'GMT+1' },
  { label: 'GMT+2 东欧标准时间', value: 'GMT+2' },
  { label: 'GMT+3 莫斯科标准时间', value: 'GMT+3' },
  { label: 'GMT+4 阿布扎比标准时间', value: 'GMT+4' },
  { label: 'GMT+5 巴基斯坦标准时间', value: 'GMT+5' },
  { label: 'GMT+6 孟加拉国标准时间', value: 'GMT+6' },
  { label: 'GMT+7 印度支那时间', value: 'GMT+7' },
  { label: 'GMT+8 北京时间', value: 'GMT+8' },
  { label: 'GMT+9 日本标准时间', value: 'GMT+9' },
  { label: 'GMT+10 澳大利亚东部标准时间', value: 'GMT+10' },
  { label: 'GMT+11 所罗门群岛时间', value: 'GMT+11' },
  { label: 'GMT+12 斐济标准时间', value: 'GMT+12' },
  { label: 'GMT+13 汤加标准时间', value: 'GMT+13' },
  { label: 'GMT+14 基里巴斯线岛时间', value: 'GMT+14' }
]

const languageOptions = [
{ label: '英语:en', value: 'en' },
  { label: '西班牙语:es', value: 'es' },
  { label: '法语:fr', value: 'fr' },
  { label: '德语:de', value: 'de' },
  { label: '意大利语:it', value: 'it' },
  { label: '葡萄牙语:pt', value: 'pt' },
  { label: '波兰语:pl', value: 'pl' },
  { label: '土耳其语:tr', value: 'tr' },
  { label: '俄语:ru', value: 'ru' },
  { label: '荷兰语:nl', value: 'nl' },
  { label: '捷克语:cs', value: 'cs' },
  { label: '阿拉伯语:ar', value: 'ar' },
  { label: '简体中文:zh-cn', value: 'zh-cn' },
  { label: '日语:ja', value: 'ja' },
  { label: '匈牙利语:hu', value: 'hu' },
  { label: '韩语:ko', value: 'ko' }
]

const audioDriverOptions = [
  { label: '默认驱动', value: 'default' },
  { label: 'ASIO', value: 'asio' },
  { label: 'WASAPI', value: 'wasapi' }
]
</script>

<style scoped lang="scss">
.live-console {
  height: calc(100% - 130px);
  
  .three-column-layout {
    display: grid;
    // grid-template-columns: repeat(3, 1fr);
    grid-template-columns: 1fr 1.5fr 1fr;  // 中间列宽度是两侧的1.5倍
    gap: 8px;
    height: 100%;
    
    .n-card {
      padding: 12px 0px;
      height: 100%;
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
    padding: 8px 0;
    border-bottom: 1px solid #f0f0f0;
    
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
  margin-top: 20px;
}

.control-blocks {
  display: grid;
  gap: 8px;  // 减小间距
}

.control-block {
  padding: 8px;  // 减小内边距
  border: 1px solid #f0f0f0;
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
  padding: 12px;
  border: 1px solid #f0f0f0;
  border-radius: 4px;
  
  .audio-info {
    margin-bottom: 8px;
    
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
  margin-bottom: 12px;
  
  .n-input-group {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}
</style>
