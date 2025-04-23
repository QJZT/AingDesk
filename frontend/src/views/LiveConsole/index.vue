<template>
  <div class="live-console">
    <!-- 直播状态和控制区域 -->
    <div class="control-panel">
      <n-card title="直播控制台">
        <div class="status-row">
          <n-tag :type="isLive ? 'success' : 'error'" round>
            {{ isLive ? '直播中' : '未开播' }}
          </n-tag>
          <n-button 
            type="primary" 
            @click="toggleLive"
            :loading="loading"
          >
            {{ isLive ? '结束直播' : '开始直播' }}
          </n-button>
        </div>
        
        <!-- 直播设置 -->
        <div class="settings-row">
          <n-input-group>
            <n-input v-model:value="streamUrl" placeholder="推流地址" />
            <n-button type="info" @click="copyStreamUrl">复制</n-button>
          </n-input-group>
        </div>
      </n-card>
    </div>

    <!-- 数据统计面板 -->
    <div class="stats-panel">
      <n-card title="直播数据">
        <div class="stats-grid">
          <n-statistic label="观看人数" :value="viewerCount" />
          <n-statistic label="点赞数" :value="likeCount" />
          <n-statistic label="评论数" :value="commentCount" />
          <n-statistic label="直播时长" :value="duration" />
        </div>
      </n-card>
    </div>

    <!-- 互动消息区域 -->
    <div class="interaction-panel">
      <n-card title="观众互动">
        <n-scrollbar style="height: 300px">
          <div 
            v-for="(msg, index) in messages" 
            :key="index" 
            class="message-item"
          >
            <span class="username">{{ msg.user }}:</span>
            <span class="content">{{ msg.content }}</span>
          </div>
        </n-scrollbar>
        
        <div class="message-input">
          <n-input 
            v-model:value="newMessage" 
            placeholder="发送公告..."
            @keyup.enter="sendAnnouncement"
          >
            <template #suffix>
              <n-button text @click="sendAnnouncement">
                <i class="i-common:send w-20 h-20"></i>
              </n-button>
            </template>
          </n-input>
        </div>
      </n-card>
    </div>
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
</script>

<style scoped lang="scss">
.live-console {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: auto 1fr;
  gap: 20px;
//   padding: 20px;
  height: 100%;
  
  .control-panel {
    grid-column: 1 / 3;
  }
  
  .status-row {
    display: flex;
    align-items: center;
    gap: 20px;
    margin-bottom: 20px;
  }
  
  .settings-row {
    margin-top: 20px;
  }
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
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
    margin-top: 20px;
  }
}
</style>