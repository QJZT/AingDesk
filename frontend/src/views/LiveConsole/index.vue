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
        <n-text depth="3" style="margin-right: 8px;">语种设置</n-text>
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
      
      <div class="setting-item"  @click="getAudioDevices">
        <n-text depth="3" style="margin-right: 8px;">音频驱动</n-text>
        <n-select 
          v-model:value="selectedAudioDriver"
          :options="audioDeviceOptions"
          placeholder="选择驱动"
          style="width: 340px;"
          @update:value="handleLanguageChange"
        />
      </div>

      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">语音模型(主)</n-text>
        <n-select 
          v-model:value="selectedSpeechModel"
          :options="speedModelOptions"
          :loading="loading2"
          placeholder="选择驱动"
          style="width: 180px;"
          @update:value="initializeSpeechModel"
        />
      </div>
    
      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">语音模型(辅助)</n-text>
        <n-select 
          v-model:value="selectedSpeechModel2"
          :options="speedModelOptions2"
          :loading="loading2"
          placeholder="选择驱动"
          style="width: 180px;"
          @update:value="initializeSpeechModel"
        />
      </div>
     
    <div class="settings-bar">
      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">Temperature ({{ temperature }})<n-popover trigger="click">
          <template #trigger>
            <i class="i-tdesign:help-circle w-16 h-16 ml-2 cursor-pointer"></i>
          </template>
          <span>控制生成文本的随机性(0.1-1.5)，值越高结果越随机</span>
        </n-popover></n-text>
        <n-slider
          v-model:value="temperature"
          :min="0.1"
          :max="1.5"
          :step="0.1"
          style="width: 150px;"
          :format-tooltip="(value) => `Temperature: ${value}`"
        />
      </div>

      <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;">TopP ({{ top_p }}) <n-popover trigger="click">
          <template #trigger>
            <i class="i-tdesign:help-circle mt-5 cursor-pointer text-[var(--n-text-color)]"></i>
          </template>
          <span>控制生成文本的多样性(0.1-1)，值越低结果越保守</span>
        </n-popover></n-text>
      
        <n-slider
          v-model:value="top_p"
          :min="0.1"
          :max="1"
          :step="0.05"
          style="width: 150px;"
          :format-tooltip="(value) => `Top P: ${value}`"
        />
      </div>
      <div class="setting-item">
        <div style="display: flex; align-items: center;">
          <n-text depth="3" style="margin-right: 8px;">TopK ({{ top_k }})<n-popover trigger="click">
            <template #trigger>
              <i class="i-tdesign:help-circle mt-5 cursor-pointer text-[var(--n-text-color)]"></i>
            </template>
            <span>限制生成时考虑的词汇数量(1-100)，值越低结果越可预测</span>
          </n-popover></n-text>
        </div>
        <n-input-number
          v-model:value="top_k"
          :min="1"
          :max="100"
          :step="1"
          style="width: 100px;"
        />
      </div>
      <div class="setting-item">
        <div style="display: flex; align-items: center;">
          <n-text depth="3" style="margin-right: 8px;">并发数<n-popover trigger="click">
            <template #trigger>
              <i class="i-tdesign:help-circle mt-5 cursor-pointer text-[var(--n-text-color)]"></i>
            </template>
            <span>取决与显卡性能，并行执行次数</span>
          </n-popover></n-text>
        </div>
        <n-input-number
          v-model:value="MAX_CONCURRENT"
          :min="1"
          :max="100"
          :step="1"
          style="width: 100px;"
        />
      </div>
      <!-- speedModelOptions2 -->
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
          <!-- <n-button ghost>
            弹幕提示词
          </n-button> -->
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
      <div class="setting-item">
      <n-text depth="3" style="margin-right: 8px;">间隔时间(ms) </n-text>
        <n-input-number 
          v-model:value="intervalTime"
          :min="0"
          :max="500000"
          :step="100"
          style="width: 100px"
        />
      </div>
       <!-- 新增直播链接输入和启动按钮 -->
       <div class="setting-item">
        <n-text depth="3" style="margin-right: 8px;"></n-text>
        <n-input-group>
          <n-button 
            v-if="!start"
            type="primary" 
            @click="registerModules"
            :loading="loading"
            style="width: 220px;"
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
            style="width: 220px;"
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
      <n-card 
        title="模块" 
        style="" 
        :segmented="{content: true,footer:true}"
        header-style="padding:10px;font-size:14px"
        footer-style="padding:10px" 
        content-style="padding:10px;height:100%"
      >
        <!-- 新增刷新按钮 -->
        <template #header-extra>
          <n-button @click="fetchModules" size="small" :loading="loadingModules">
            <!-- <template #icon>
              <i class="i-tdesign:refresh"></i>
            </template> -->
            刷新
          </n-button>
        </template>
        <n-infinite-scroll style="height: 480px;overflow-y: auto;background: var(--n-color-embedded);
      border-radius: var(--n-border-radius);"  :distance="10">
          <div class="">
            <div class="control-block" v-for="(block, index) in modules" :key="index" style="">
              <n-divider v-if="index != 0"  style="font-size: 15px;height: 1px;margin: 12px 0;"></n-divider>
            <div class="block-header">
              <h4>{{ block.module_name }}</h4>
            <!-- <n-switch v-model:value="block.isActive"  size="small" @update:value="handleVolumeChange(block)"/> -->
            <n-switch v-model:value="block.retAi"  size="small" @update:value="handleVolumeChange(block)" v-if="block.module_type == 'base'">
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
                :min="10"
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
            <div style="display: flex; align-items: center; gap: 8px; width: 100%;" v-if="block.module_type == 'base'">
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
            <div style="display: flex; align-items: center; gap: 8px; width: 100%;" v-if="block.module_type == 'base'">
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
                @click="getHumanVoiceFiles"
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
                    <template #header-extra>
                      <div class="setting-item" >
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
        </template>
        <!-- <template #header-extra>
          <n-switch 
            v-model:value="autoReadMode" 
          >
            <template #checked>自动朗读</template>
            <template #unchecked>关闭朗读</template>
          </n-switch>
        </template> -->
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
            filter: 'brightness(1.2) contrast(1.2)',
            maxWidth: '100%'
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
            style="white-space: nowrap; overflow: hidden; text-overflow: ellipsis; display: flex; align-items: center; gap: 4px;max-width: 50vw;"
          >
          <!-- getMessageColor -->
          <span class="content" style="line-height: 1;width: 45vw;"  :style="{
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
        
        弹幕检测：{{ EnterBarrageContent }}
        <n-input type="text" v-model:value="EnterBarrageContent"/>
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
    
    
   <!-- 发送消息卡片 -->
        <!-- 用户消息和语音输入卡片的容器 -->
        <div class="message-voice-container">
      <!-- 用户消息卡片 -->
      <n-card 
        title="用户消息" 
        :segmented="{content: true, footer: true}"
        header-style="padding:10px;font-size:14px"
        footer-style="padding:10px" 
        content-style="padding:10px;"
        class="message-card"
      >
        <div class="control-block">
          <div class="control-row">
            <div class="control-item">
              <div class="control-label">音量: {{ announcementVolume }}%</div>
              <n-slider
                v-model:value="announcementVolume"
                :min="10"
                :max="100"
                :step="1"
                class="control-slider"
                tooltip
                placement="bottom"
                :format-tooltip="(value) => `音量: ${value}%`"
                aria-label="消息音量调整"
              />
            </div>
            <div class="control-item">
              <div class="control-label">语速: {{ announcementSpeed }}x</div>
              <n-slider
                v-model:value="announcementSpeed"
                :min="0.5"
                :max="2"
                :step="0.01"
                class="control-slider"
                tooltip
                placement="bottom"
                :format-tooltip="(value) => `语速: ${value}x`"
                aria-label="消息语速调整"
              />
            </div>
            <div class="control-item">
              <div class="control-label">人声</div>
              <n-select 
                v-model:value="announcementVoice"
                :options="humanVoiceOptions"
                class="control-select"
                placeholder="选择人声"
                clearable
                size="small"
                @click="getHumanVoiceFiles"
                aria-label="选择消息人声"
              />
            </div>
          </div>
        </div>
        <div style="height: 10px;"></div>
        <div class="message-input">
          <n-input-group>
            <n-input 
              v-model:value="announcementMessage"
              type="textarea"
              :placeholder="`输入消息后按 Enter 发送`"
              :autosize="{ minRows: 2, maxRows: 5 }"
              @keyup.enter.ctrl.prevent="announcementMessage += '\n'"
              @keyup.enter.exact="sendUserMessage"
              style="flex: 1; min-width: 0;"
              aria-label="输入消息"
            />
            <n-button 
              type="primary" 
              @click="sendUserMessage"
              :disabled="!announcementMessage.trim() || !announcementVoice"
              aria-label="发送消息"
              class="send-button"
            >
              <template #icon>
                <i class="i-tdesign:send w-20 h-20"></i>
              </template>
              发送
            </n-button>
          </n-input-group>
        </div>
      </n-card>

      <!-- 语音输入卡片 -->
      <n-card 
        title="语音输入" 
        :segmented="{content: true, footer: true}"
        header-style="padding:10px;font-size:14px"
        footer-style="padding:10px" 
        content-style="padding:10px;"
        class="voice-card"
      >
        <n-input-group class="voice-input-group">
          <n-button 
            :type="isRecording ? 'warning' : 'default'"
            :class="{ 'recording': isRecording }"
            class="voice-button"
            @mousedown="startVoiceInput"
            @mouseup="stopVoiceInput"
            @mouseleave="isRecording && stopVoiceInput()"
            :focusable="false"
            :aria-pressed="isRecording"
            :aria-label="isRecording ? '停止录音' : '开始录音'"
          >
            <template #icon>
              <i :class="isRecording ? 'i-tdesign:stop w-18 h-18' : 'i-tdesign:microphone-1 w-18 h-18'"></i>
            </template>
            {{ isRecording ? '录音中' : '语音' }}
          </n-button>
          <div class="mic-select">
            <n-text depth="3" class="mic-label">麦克风：</n-text>
            <n-select 
              v-model:value="selectedMicrophoneDriver"
              :options="audioDeviceOptions"
              placeholder="选择麦克风"
              class="mic-select-input"
              @update:value="handleLanguageChange"
              @click="getMicrophoneDevices"
            />
          </div>
        </n-input-group>
      </n-card>
    </div>



    <!-- 弹窗 -->
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
import KnowledgeChoosePanel from "@/views/KnowleadgeStore/components/KnowledgeChoosePanel.vue";
import { getKnowledgeStoreData } from "../KnowleadgeStore/store";
import { getHeaderStoreData } from '../Header/store';
import { getSoftSettingsStoreData } from '@/views/SoftSettings/store';
const {  showModelList } = getHeaderStoreData()
import { message, useDialog } from "@/utils/naive-tools"
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
const selectedAudioDriver = ref('')// 音频驱动
const selectedSpeechModel = ref('')
const selectedSpeechModel2 = ref('')
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
        soundCard:selectedAudioDriver.value,
        backgroundSoundCard: "耳机 (Q38-2 Stereo)",
        microphone: selectedMicrophoneDriver.value,//麦克风
       })
    })
    const data = await response.json()
    console.log('语言设置成功:', data)
  } catch (error) {
    console.error('语言设置失败:', error)
  }
}
const loading2 = ref(false)
const initializeSpeechModel = async () => {
  // 加上logding
  if (selectedSpeechModel.value == "" || selectedSpeechModel2.value == "") {
    // message.error("请选择语音模型")
    return
  }
  loading2.value = true

  try {
   const response = await fetch('http://127.0.0.1:7074/set_model', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        gpt_model_path: selectedSpeechModel.value,
        sovits_model_path: selectedSpeechModel2.value
       })
    })
    const data = await response.json()
    console.log('语音模型初始化成功:', data)
    message.success("模型初始化成功")
   }
   catch (error) {
    console.error('语音模型初始化失败:', error) 
  }finally {
    loading2.value = false
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
{ label: '简体中文', value: 'zh' },

  // { label: '西班牙语', value: 'es' },
  // { label: '法语', value: 'fr' },
  // { label: '德语', value: 'de' },
  // { label: '意大利语', value: 'it' },
  // { label: '葡萄牙语', value: 'pt' },
  // { label: '波兰语', value: 'pl' },
  // { label: '土耳其语', value: 'tr' },
  // { label: '俄语', value: 'ru' },
  // { label: '荷兰语', value: 'nl' },
  // { label: '捷克语', value: 'cs' },
  // { label: '阿拉伯语', value: 'ar' },
  // { label: '日语', value: 'ja' },
  // { label: '匈牙利语', value: 'hu' },
  // { label: '韩语', value: 'ko' }
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
  // console.log("block----------------------------");
  // console.log(block);
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

  //getMicrophoneDevices 麦克风  selectedMicrophoneDriver
  const selectedMicrophoneDriver = ref('');// 麦克风使用的音频驱动
  const getMicrophoneDevices = async () => {
    try {
      const response = await fetch('http://127.0.0.1:7073/get_sound_cards', {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' }
      });
      const data = await response.json();
      audioDeviceOptions.value = data.map(item => ({
        label: item,
        value: item
      }));
      const cable = audioDeviceOptions.value.find(opt => opt.label === 'CABLE Output (VB-Audio Virtual Cable)');
      if (cable) {
        selectedMicrophoneDriver.value = cable.value;
      }
    } catch (error) {
      console.error('获取音频设备失败:', error);
    }
  };




const speedModelOptions = ref([])
const speedModelOptions2 = ref([])

const getSpeedModels = async () => {
  try {
    const response = await fetch('http://127.0.0.1:7074/get_model_filenames', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      } 
    })
    const data = await response.json()
    
    speedModelOptions.value = (data.data.gpt_filenames || []).map(item => ({
      label: item,
      value: item
    }))
    speedModelOptions2.value = (data.data.sovits_filenames || []).map(item => ({
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
const loadingModules = ref(false)
const fetchModules = async () => {
  loadingModules.value = true
  try {
    console.log('正在获取模块数据...')
    let  newList   =[]
    const response = await fetch('http://127.0.0.1:7072/base-modules')
    let list = await response.json()
    console.log("lis23132t:",list);
    for (const  [index, module] of list.entries()) {
      console.log("module.id:",module.id);
      let  kv = await getModulesKv(module.id)
      console.log("kv:",kv);
      module.volume = kv.volume || 50
      module.isActive= kv.isActive || false
      module.retAi= kv.retAi || false
      module.speed = kv.speed || 1.0
      module.selectedNameId = kv.selectedNameId || ""
      await setModulesKv(module.id,module)
      newList.push(module)
    }
    modules.value=newList
  } catch (error) {
    console.error('获取模块数据失败:', error)
  } finally {
    setTimeout(() => {
      loadingModules.value = false
    }, 800)
  }
}

onMounted(() => {
  fetchModules()      // 获取模块列表
  getMames()          // 获取名称列表
  getPromptRewrite()  // 获取提示词重写配置
  getAudioDevices()   // 获取音频设备列表
  getHumanVoiceFiles() // 获取人声文件列表
  getSpeedModels()    // 获取语速模型列表
  // 功能：组件挂载时设置事件监听器和初始焦点
  // 逻辑：添加 keydown 监听器，聚焦 textarea
  document.addEventListener('keydown', handleKeyDown);
  nextTick(() => {
    messageInput.value?.$.el.focus();
  });

  // setInterval(fetchModules, 25000)
  // return
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
  // WebSocket 连接事件处理
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
    //弹幕回答题
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
      EnterLiveRoomUserName.value = message.u //进入直播间的用户名
      
      if (!selectedFilters.value.includes("JoinEvent")){
        return;
      }
              console.log("JoinEvent:",message);
              addMessage({content: message.s,data_types: "JoinEvent"})
    });

    // socket.value.on('ShareEvent', function(message) { //分享
    //   if (!selectedFilters.value.includes("ShareEvent")){
    //     return;
    //   }
    //           console.log("ShareEvent:",message);
    //           addMessage({content: message.s,data_types: "ShareEvent"})
    // });
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

    socket.value.on('ShareEvent', function(message) { // 分享直播间事件
      ShareRoomUserName.value = message.u; // 更新分享用户名
      if (!selectedFilters.value.includes("ShareEvent")) {
          return;
      }
      console.log("ShareEvent:", message);
      addMessage({ content: message.s, data_types: "ShareEvent" });
    });

    socket.value.on('FollowEvent', function(message) { // 关注直播间事件
        FollowRoomUserName.value = message.u; // 更新关注用户名
        if (!selectedFilters.value.includes("FollowEvent")) {
            return;
        }
        console.log("FollowEvent:", message);
        addMessage({ content: message.s, data_types: "FollowEvent" });
    });
    // 模拟 ShareEvent
  // setTimeout(() => {
  //   ShareRoomUserName.value = "UserA";
  //   addMessage({ content: "UserA 分享了直播间", data_types: "ShareEvent" });
  //   console.log("Simulated ShareEvent: UserA");
  // }, 21000);

  // // 模拟 FollowEvent
  // setTimeout(() => {
  //   FollowRoomUserName.value = "UserB";
  //   addMessage({ content: "UserB 关注了直播间", data_types: "FollowEvent" });
  //   console.log("Simulated FollowEvent: UserB");
  // }, 10000);

    // socket.value.on('LiveEndEvent', function(message) { //直播结束
    //           console.log("LiveEndEvent:",message);
    //           addMessage({content: message})
    // });
    // socket.value.on('DisconnectEvent', function(message) { //断开连接
    //           console.log("DisconnectEvent:",message);
    //           addMessage({content: message})
    // });
})

onUnmounted(() => {
    document.removeEventListener('keydown', handleKeyDown);
    if (socket.value) {
      socket.value.disconnect();
      socket.value = null;
    }
  });

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
const startUUID = ref("")//启动开关
//注册模块，

const registerModules = async () => {
    //启动之前检查变量做友好提示
    if (selectedModel.value == "") {
      message.error("请选择模型")
      return 
    }
    playList.value = []
    uesPlayList.value = []
    // await fetchModules()
    // return
    //循环模块列表 干活
    const newuuid = crypto.randomUUID()

    start.value = true
    startUUID.value = newuuid;
    for (const module of modules.value) {
//     TriggerSceneLoop     = "SceneLoop"     // 控场循环
//     TriggerIntervalLoop  = "IntervalLoop"  // 间隔循环：按照指定间隔时间循环执行
//     TriggerBarrageComment = "BarrageComment" // 弹幕评论：是否自动发送弹幕评论
//     TriggerSendGift      = "SendGift"      // 送礼物：是否自动发送礼物
//     TriggerLike          = "Like"          // 点赞：是否自动点赞
//     TriggerEnterLiveRoom = "EnterLiveRoom" // 进入直播间：是否自动进入直播间
//     TriggerShareRoom     = "ShareRoom"     // 分享直播间
//     TriggerFollowRoom    = "FollowRoom"    // 关注直播间
        //循环读取
        if (module.module_type == "base") { 
          if (module.trigger_conditions.includes("SceneLoop")) { // 执行循环：是否自动循环执行任务
            SceneLoop(module,newuuid)
          }
          if (module.trigger_conditions.includes("IntervalLoop")) { //间隔循环：按照指定间隔时间循环执行
            SceneLoop(module,newuuid)
          }
          //弹幕评论 BarrageComment
          if (module.trigger_conditions.includes("BarrageComment")) {// 弹幕
              BarrageComment(module,newuuid)
          }
          //送礼物 SendGift
          if (module.trigger_conditions.includes("SendGift")) { // 送礼物：是否自动发送礼物
            SendGift(module,newuuid)
          }
          //点赞 Like
          if (module.trigger_conditions.includes("Like")) {// 点赞：是否自动点赞
            includesLike(module,newuuid)
          }
          //进入直播间 EnterLiveRoom
          if (module.trigger_conditions.includes("EnterLiveRoom")) {// 进入直播间：是否自动进入直播间
            includesEnterLiveRoom(module,newuuid)
          }
          // 分享直播间 ShareRoom
          if (module.trigger_conditions.includes("ShareRoom")) { // 分享直播间：是否自动分享直播间
                ShareRoom(module, newuuid);
            }
            // 关注直播间 FollowRoom
            if (module.trigger_conditions.includes("FollowRoom")) { // 关注直播间：是否自动关注直播间
                FollowRoom(module, newuuid);
            }
          
        }
        
        //     TriggerShareRoom     = "ShareRoom"     // 分享直播间
//     TriggerFollowRoom    = "FollowRoom"    // 关注直播间
    }
    playListConsumption()
}

const intervalTime = ref(0) // 默认1000毫秒 1秒=1000毫秒
//消费playList
//定时播报
const playListConsumption= async () => {
    do {
        if (playList.value.length > 0) { //队列消费
          await new Promise(resolve => setTimeout(resolve, intervalTime.value))   // 等待1000毫秒
          console.log("kai:",playList.value);
            const item = playList.value.shift() //出队
            console.log("jie:",playList.value);
            if (item) {
                // uesPlayList.value.push(item) //已播放列表
                // if (uesPlayList.value.length > 100) {
                //     uesPlayList.value.shift() //删除第一条
                // }
                let data = await play_task_voice_api(item.filename, item.play_mode); //播放
                // await new Promise(resolve => setTimeout(resolve, data.duration * 1000)) //等待
            }
        }
        await new Promise(resolve => setTimeout(resolve, 100))   
    } while (start.value);
}

//循环模块处理
const SceneLoop= async (module,newuuil) => {
    let index = 0 //当前播放索引
    const script_content_len = module.script_content.length //脚本长度
    do {
      console.log("SceneLoop module:",module);
      
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
            console.log("randomTime:",randomTime);
            
            await new Promise(resolve => setTimeout(resolve, randomTime * 1000)) //等待
        }

        let content = module.script_content[index] 
        while (playList.value.length >= 5) { //等待队列 消费完探入
            await new Promise(resolve => setTimeout(resolve, 1000))
        }
        // 是否改写
        if (module.retAi) { //是否改写
          if (model_api.value == "") { //是否改写
            console.log("请选择模型");
            await new Promise(resolve => setTimeout(resolve, 4000))
            continue
          }
          let prompt = await ReplaceText(promptText.value) //提示词 赋值变量
          content = await ReplaceText(content) //内容赋值变量
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
        if (module.trigger_conditions.includes("IntervalLoop")){ //定时触发的 插队前面
            playList.value.unshift({ //入队
              content: "["+module.module_name+"]"+content, //内容
              filename: newFileName, //文件名
              play_mode: "serial", //播放模式
          })   
        }else {
          playList.value.push({ //入队
            content: "["+module.module_name+"]"+content, //内容
            filename: newFileName, //文件名
            play_mode: "serial", //播放模式
          })    
        }
        if (index == script_content_len - 1) { //循环
            index = 0
        } else {
            index = index + 1
        }
    } while (start.value && newuuil == startUUID.value);
}


// 点赞 Like
const EnterSupportRoomUserName = ref("")  //点赞直播间 用户名
const includesLike= async (module,newuuil) => {
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
        let content = await ReplaceText(module.script_content[index] ) //内容赋值变量
        // 添加 AI 改写逻辑
        if (module.retAi) { // 是否启用 AI 改写
              if (model_api.value == "") { // 检查是否选择了模型
                console.log("请选择模型");
                await new Promise(resolve => setTimeout(resolve, 4000));
                continue;
              }
              let prompt = await ReplaceText(promptText.value); // 提示词赋值变量 
              const apidata = await DisposableSendApi(
                model_api.value,
                parameters_api.value,
                content, // 文本
                prompt, // 提示词
                supplierName_api.value // 供应商名称
              );
              console.log("apidata:", apidata);
              content = apidata; // 替换为 AI 改写后的内容
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
        // 插入到头部
        playList.value.unshift({ //入队
            content: "["+module.module_name+"]"+content, //内容
            filename: newFileName, //文件名
            play_mode: "serial", //播放模式
        })    
        EnterSupportRoomUserName.value = "" //清空
        if (index == script_content_len - 1) { //循环
            index = 0
        } else {
            index = index + 1
        }
    } while (start.value && newuuil == startUUID.value);
}

//送礼物数据
const EnterGiftUserName = ref("")  //送礼物 用户名
const EnterGiftNum = ref("")  //送礼物数量
const EnterGiftGoodsName = ref("")  //送的礼物名字
const SendGift= async (module,newuuil) => {
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
        if (module.retAi) { //是否改写
            if (model_api.value == "") { //是否改写
              console.log("请选择模型");
              await new Promise(resolve => setTimeout(resolve, 4000))
              continue
            }
           let prompt = await ReplaceText(promptText.value) //提示词 赋值变量
           let new_content = await ReplaceText(module.script_content[index]) //内容赋值变量
           const  apidata =  await DisposableSendApi(
            model_api.value,
            parameters_api.value,
            new_content, //文本
            prompt, //提示词
            supplierName_api.value, //供应商名称
          )
          console.log("apidata:",apidata);
          content = apidata  //替换
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
        // 插入到头部
        playList.value.unshift({ //入队
            content: "["+module.module_name+"]"+content, //内容
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
    } while (start.value && newuuil == startUUID.value);
}

const EnterLiveRoomUserName = ref(""); // 进入直播间用户名，初始为空字符串

const includesEnterLiveRoom = async (module, newuuil) => {
    let index = 0; // 当前播放索引，跟踪脚本内容的当前位置
    const script_content_len = module.script_content.length; // 脚本内容的总长度，用于循环控制

    do {
        if (!module.isActive) { 
            // 如果模块未激活，等待1秒后继续下一次循环，避免无谓的资源占用
            await new Promise(resolve => setTimeout(resolve, 1000));
            continue;
        }

        if (module.read_step == "random") { // 检查是否为随机读取模式
            if (index == 0) { 
                // 仅在索引为0时（即循环开始）打乱脚本内容数组，确保随机性
                module.script_content = shuffleArray(module.script_content);
            }
        }

        if (module.interval_time_start != 0 && module.interval_time_end != 0) { 
            // 如果设置了间隔时间范围，生成随机间隔时间（单位：秒）
            const randomTime = Math.floor(Math.random() * (module.interval_time_end - module.interval_time_start + 1)) + module.interval_time_start;
            await new Promise(resolve => setTimeout(resolve, randomTime * 1000)); // 等待随机时间
        }

        while (EnterLiveRoomUserName.value == "") { 
            // 循环等待直到有用户进入直播间（EnterLiveRoomUserName.value 被 WebSocket 更新）
            await new Promise(resolve => setTimeout(resolve, 1000)); // 每1秒检查一次
        }

        let content = module.script_content[index]; 
        // 初始从脚本内容中取值，准备替换变量

        // 始终调用 ReplaceText 替换变量（包括 {进入直播间用户名}），无论是否启用AI改写
        content = await ReplaceText(content); // 替换所有变量，如 {进入直播间用户名}

        if (module.retAi) { // 如果启用了AI改写
            if (model_api.value == "") { 
                // 如果未选择模型，记录日志并等待4秒后继续，避免无效操作
                console.log("请选择模型");
                await new Promise(resolve => setTimeout(resolve, 4000));
                continue;
            }
            let prompt = await ReplaceText(promptText.value); // 获取并替换提示词中的变量
            const apidata = await DisposableSendApi(
                model_api.value,
                parameters_api.value,
                content, // 使用已替换变量的 content
                prompt, // 替换后的提示词
                supplierName_api.value, // 供应商名称
            );
            console.log("apidata:", apidata); // 记录AI返回的数据
            content = apidata; // 用AI处理后的内容替换原始内容
        }

        const newFileName = crypto.randomUUID() + "_" + Date.now() + '.wav'; // 生成唯一文件名
        let ok = await generate_wav_api(
            content, // 当前内容（已替换变量，可能经AI改写）
            selectedLanguage.value || "en", // 语种，默认为英语
            newFileName, // 生成的文件名
            module.selectedNameId, // 音色文件名
            module.speed, // 语速
            module.volume / 100 // 音量（转换为0-1范围）
        );
        if (!ok) { // 如果音频生成失败
            await new Promise(resolve => setTimeout(resolve, 1000)); // 等待1秒后重试
            continue;
        }

        // 插入到播放列表头部
        playList.value.unshift({
            content: "[" + module.module_name + "]" + content, // 内容前添加模块名称
            filename: newFileName, // 文件名
            play_mode: "serial", // 串行播放模式
        });

        // 清空变量，确保下次循环等待新用户
        EnterLiveRoomUserName.value = ""; // 清空用户名
        EnterGiftNum.value = ""; // 清空礼物数量
        EnterGiftGoodsName.value = ""; // 清空礼物名称

        if (index == script_content_len - 1) { // 如果到达脚本末尾
            index = 0; // 重置索引，从头开始
        } else {
            index = index + 1; // 移动到下一条脚本
        }
    } while (start.value && newuuil == startUUID.value); // 持续运行直到停止或UUID变更
};


const ShareRoomUserName = ref(""); // 分享直播间用户名，初始为空字符串
const ShareRoom = async (module, newuuid) => {
    let index = 0; // 当前播放索引，跟踪脚本内容的当前位置
    const script_content_len = module.script_content.length; // 脚本内容的总长度，用于循环控制

    do {
        if (!module.isActive) { 
            // 如果模块未激活，等待1秒后继续下一次循环，避免无谓的资源占用
            await new Promise(resolve => setTimeout(resolve, 1000));
            continue;
        }

        if (module.read_step == "random") { // 检查是否为随机读取模式
            if (index == 0) { 
                // 仅在索引为0时（即循环开始）打乱脚本内容数组，确保随机性
                module.script_content = shuffleArray(module.script_content);
            }
        }

        if (module.interval_time_start != 0 && module.interval_time_end != 0) { 
            // 如果设置了间隔时间范围，生成随机间隔时间（单位：秒）
            const randomTime = Math.floor(Math.random() * (module.interval_time_end - module.interval_time_start + 1)) + module.interval_time_start;
            await new Promise(resolve => setTimeout(resolve, randomTime * 1000)); // 等待随机时间
        }

        while (ShareRoomUserName.value == "") { 
            // 循环等待直到有用户分享直播间（ShareRoomUserName.value 被 WebSocket 更新）
            await new Promise(resolve => setTimeout(resolve, 1000)); // 每1秒检查一次
        }

        let content = module.script_content[index]; 
        // 初始从脚本内容中取值，准备替换变量

        // 始终调用 ReplaceText 替换变量（包括 {分享直播间用户名}），无论是否启用AI改写
        content = await ReplaceText(content); // 替换所有变量，如 {分享直播间用户名}

        if (module.retAi) { // 如果启用了AI改写
            if (model_api.value == "") { 
                // 如果未选择模型，记录日志并等待4秒后继续，避免无效操作
                console.log("请选择模型");
                await new Promise(resolve => setTimeout(resolve, 4000));
                continue;
            }
            let prompt = await ReplaceText(promptText.value); // 获取并替换提示词中的变量
            const apidata = await DisposableSendApi(
                model_api.value,
                parameters_api.value,
                content, // 使用已替换变量的 content
                prompt, // 替换后的提示词
                supplierName_api.value, // 供应商名称
            );
            console.log("apidata:", apidata); // 记录AI返回的数据
            content = apidata; // 用AI处理后的内容替换原始内容
        }

        const newFileName = crypto.randomUUID() + "_" + Date.now() + '.wav'; // 生成唯一文件名
        let ok = await generate_wav_api(
            content, // 当前内容（已替换变量，可能经AI改写）
            selectedLanguage.value || "en", // 语种，默认为英语
            newFileName, // 生成的文件名
            module.selectedNameId, // 音色文件名
            module.speed, // 语速
            module.volume / 100 // 音量（转换为0-1范围）
        );
        if (!ok) { // 如果音频生成失败
            await new Promise(resolve => setTimeout(resolve, 1000)); // 等待1秒后重试
            continue;
        }

        // 插入到播放列表头部
        playList.value.unshift({
            content: "[" + module.module_name + "]" + content, // 内容前添加模块名称
            filename: newFileName, // 文件名
            play_mode: "serial", // 串行播放模式
        });

        // 清空变量，确保下次循环等待新用户
        ShareRoomUserName.value = ""; // 清空分享直播间用户名
        EnterGiftNum.value = ""; // 清空礼物数量
        EnterGiftGoodsName.value = ""; // 清空礼物名称

        if (index == script_content_len - 1) { // 如果到达脚本末尾
            index = 0; // 重置索引，从头开始
        } else {
            index = index + 1; // 移动到下一条脚本
        }
    } while (start.value && newuuid == startUUID.value); // 持续运行直到停止或UUID变更
};


const FollowRoomUserName = ref(""); // 关注直播间用户名，初始为空字符串
const FollowRoom = async (module, newuuid) => {
    let index = 0; // 当前播放索引，跟踪脚本内容的当前位置
    const script_content_len = module.script_content.length; // 脚本内容的总长度，用于循环控制

    do {
        if (!module.isActive) { 
            // 如果模块未激活，等待1秒后继续下一次循环，避免无谓的资源占用
            await new Promise(resolve => setTimeout(resolve, 1000));
            continue;
        }

        if (module.read_step == "random") { // 检查是否为随机读取模式
            if (index == 0) { 
                // 仅在索引为0时（即循环开始）打乱脚本内容数组，确保随机性
                module.script_content = shuffleArray(module.script_content);
            }
        }

        if (module.interval_time_start != 0 && module.interval_time_end != 0) { 
            // 如果设置了间隔时间范围，生成随机间隔时间（单位：秒）
            const randomTime = Math.floor(Math.random() * (module.interval_time_end - module.interval_time_start + 1)) + module.interval_time_start;
            await new Promise(resolve => setTimeout(resolve, randomTime * 1000)); // 等待随机时间
        }

        while (FollowRoomUserName.value == "") { 
            // 循环等待直到有用户关注直播间（FollowRoomUserName.value 被 WebSocket 更新）
            await new Promise(resolve => setTimeout(resolve, 1000)); // 每1秒检查一次
        }

        let content = module.script_content[index]; 
        // 初始从脚本内容中取值，准备替换变量

        // 始终调用 ReplaceText 替换变量（包括 {关注直播用户名}），无论是否启用AI改写
        content = await ReplaceText(content); // 替换所有变量，如 {关注直播用户名}

        if (module.retAi) { // 如果启用了AI改写
            if (model_api.value == "") { 
                // 如果未选择模型，记录日志并等待4秒后继续，避免无效操作
                console.log("请选择模型");
                await new Promise(resolve => setTimeout(resolve, 4000));
                continue;
            }
            let prompt = await ReplaceText(promptText.value); // 获取并替换提示词中的变量
            const apidata = await DisposableSendApi(
                model_api.value,
                parameters_api.value,
                content, // 使用已替换变量的 content
                prompt, // 替换后的提示词
                supplierName_api.value, // 供应商名称
            );
            console.log("apidata:", apidata); // 记录AI返回的数据
            content = apidata; // 用AI处理后的内容替换原始内容
        }

        const newFileName = crypto.randomUUID() + "_" + Date.now() + '.wav'; // 生成唯一文件名
        let ok = await generate_wav_api(
            content, // 当前内容（已替换变量，可能经AI改写）
            selectedLanguage.value || "en", // 语种，默认为英语
            newFileName, // 生成的文件名
            module.selectedNameId, // 音色文件名
            module.speed, // 语速
            module.volume / 100 // 音量（转换为0-1范围）
        );
        if (!ok) { // 如果音频生成失败
            await new Promise(resolve => setTimeout(resolve, 1000)); // 等待1秒后重试
            continue;
        }

        // 插入到播放列表头部
        playList.value.unshift({
            content: "[" + module.module_name + "]" + content, // 内容前添加模块名称
            filename: newFileName, // 文件名
            play_mode: "serial", // 串行播放模式
        });

        // 清空变量，确保下次循环等待新用户
        FollowRoomUserName.value = ""; // 清空关注直播间用户名
        EnterGiftNum.value = ""; // 清空礼物数量
        EnterGiftGoodsName.value = ""; // 清空礼物名称

        if (index == script_content_len - 1) { // 如果到达脚本末尾
            index = 0; // 重置索引，从头开始
        } else {
            index = index + 1; // 移动到下一条脚本
        }
    } while (start.value && newuuid == startUUID.value); // 持续运行直到停止或UUID变更
};



//弹幕评论数据
const EnterBarrageUserName = ref("")  //弹幕评论 用户名
const EnterBarrageContent = ref("")  //弹幕评论 内容
const BarrageComment= async (module,newuuil) => {
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
            await new Promise(resolve => setTimeout(resolve, 2000))
        }
        if (model_api.value == "") { //是否改写
            console.log("请选择模型");
            await new Promise(resolve => setTimeout(resolve, 4000))
            continue
        }
        let content = "" // 
        let prompt = await ReplaceText(module.script_content[0]) //提示词
          let ai_user_content = await  ReplaceText(EnterBarrageContent.value) //弹幕内容作为ai输入
            const  apidata2 =  await DisposableSendApi(
              model_api.value,
              parameters_api.value,
              ai_user_content, //文本
              prompt, //提示词
              supplierName_api.value, //供应商名称
            )
            console.log("apidata:",apidata2);
            content = apidata2
        if (module.retAi) { //是否改写
            if (model_api.value == "") { //是否改写
              console.log("请选择模型");
              await new Promise(resolve => setTimeout(resolve, 4000))
              continue
            }
           let prompt = await  ReplaceText(promptText.value) //提示词 赋值变量
           let new_content = await ReplaceText(content) //内容赋值变量
           const  apidata =  await DisposableSendApi(
            model_api.value,
            parameters_api.value,
            new_content, //文本
            prompt, //提示词
            supplierName_api.value, //供应商名称
          )
          console.log("apidata:",apidata);
          content = apidata  //替换
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
        // 插入到头部
        playList.value.unshift({ //插队
            content: "["+module.module_name+"]"+content, //内容
            filename: newFileName, //文件名
            play_mode: "serial", //播放模式
        })    
        EnterBarrageUserName.value = "" //清空
        EnterBarrageContent.value = "" //清空
    } while (start.value && newuuil == startUUID.value);
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
const temperature = ref(0.6)
const top_p = ref(0.8)
const top_k = ref(50)

const MAX_CONCURRENT = ref(1); // 最大并发数设为1，即串行执行
let activeRequests = 0;
let queue: Array<() => Promise<void>> = [];
const generate_wav_api = async (_text:string,
    _language:string,
    _filename:string,
    _speaker_wav:string,
    _speed: number,
    _volume: number) => {
      return new Promise((resolve, reject) => {
    const execute = async () => {
      if (activeRequests >= MAX_CONCURRENT.value || queue.length === 0) return;
      
      activeRequests++;
      const task = queue.shift();
      try {
        await task?.();
        activeRequests--;
        execute(); // 处理下一个任务
      } catch (err) {
        activeRequests--;
        execute();
        throw err;
      }
    };
    queue.push(async () => {
      try {
        const response = await fetch('http://127.0.0.1:7074/generate_wav', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                "refer_wav_path": _speaker_wav,
                "prompt_text": get_human_voice_files_text_map.value[_speaker_wav],
                "prompt_language": "zh",
                "text": _text,
                "text_language": _language,//目标音频
                "cut_punc": "",
                "top_k": top_k.value, // *
                "top_p": top_p.value, // *
                "temperature": temperature.value,// 采样温度，控制生成的随机性，值越低越保守 // *
                "speed": _speed,
                "filename": _filename,
                "volume" : _volume,//
                "sample_rate" : 22050//
            }),
        });
        resolve(response.ok);
      } catch (err) {
        reject(err);
      }
    });

    execute();
  });
    // try {
    // // generate_wav_api_runing = true;
    //     const response = await fetch('http://127.0.0.1:7074/generate_wav', {
    //         method: 'POST',
    //         headers: {
    //             'Content-Type': 'application/json',
    //         },
    //         body: JSON.stringify({ 
    //             "refer_wav_path": _speaker_wav,
    //             "prompt_text": get_human_voice_files_text_map.value[_speaker_wav],
    //             "prompt_language": "zh",
    //             "text": _text,
    //             "text_language": _language,//目标音频
    //             "cut_punc": "",
    //             "top_k": top_k.value, // *
    //             "top_p": top_p.value, // *
    //             "temperature": temperature.value,// 采样温度，控制生成的随机性，值越低越保守 // *
    //             "speed": _speed,
    //             "filename": _filename,
    //             "volume" : _volume,//
    //             "sample_rate" : 22050//
    //         }),
    //     });
    //   // generate_wav_api_runing = false;
    //   return response.ok;
    // }finally {
    // }
    
}










// 播放任务
const play_task_voice_api = async (_filename:string,play_mode:string) => {
  console.log(`请求参数: filename=${_filename}, play_mode=${play_mode}`);
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











// textarea 引用，用于检测焦点状态
const messageInput = ref<HTMLInputElement | null>(null); // 明确指定类型为 HTMLInputElement 或 null
// --- 响应式状态 ---
// 消息输入相关
const announcementMessage = ref(''); // 用户输入的消息文本
const announcementVolume = ref(50); // 音量：10-100%，默认 50%
const announcementSpeed = ref(1.0); // 语速：0.5-2x，默认 1x
const announcementVoice = ref(''); // 选择的人声

const sendUserMessage = async () => {
  if (!announcementMessage.value.trim() || !announcementVoice.value) {
    message.error("发送消息请请输入消息并选择人声")
    return
  }
  try {
      // 如果配置了 AI 改写
  if (model_api.value) {
        let prompt = await ReplaceText(promptText.value); // 提示词
        let ai_user_content = await ReplaceText(announcementMessage.value); // 用户输入的消息文本

        const apidata = await DisposableSendApi(
          model_api.value,
          parameters_api.value,
          ai_user_content, // 文本
          prompt, // 提示词
          supplierName_api.value // 供应商名称
        );

        if (apidata) {
          announcementMessage.value = apidata; // 使用 AI 处理后的文本更新消息
        } else {
          message.error("AI 处理失败");
          return;
        }
      }

    
    const newFileName = crypto.randomUUID() + "_" + Date.now() + '.wav'
    const ok = await generate_wav_api(
      announcementMessage.value,//文本
      selectedLanguage.value || "en",//语种
      newFileName,//生成文件名
      announcementVoice.value,
      announcementSpeed.value,//生成速度
      announcementVolume.value / 100//生成音量
    )
    if (!ok) {
      message.error("音频生成失败")
      return
    }
    
    playList.value.unshift({
      content: "[插话]"+announcementMessage.value,
      filename: newFileName, // 文件名
      play_mode: "serial" // 播放模式
    })
    message.success("消息已发送")
  } catch (error) {
    console.error("发送用户消息失败:", error)
    message.error("发送失败，请重试")
  }
  announcementMessage.value = ''
}

const isRecording = ref(false);


  //录制开始
  const startVoiceInput = async () => {
    // 检查是否已经在录音状态，如果是，则直接返回
    if (isRecording.value) return;
  
    // 将录音状态设置为 true，表示开始录音
    isRecording.value = true;
  
    // 提示用户开始录音
    message.info('开始语音输入');
  
    try {
      // 发起 POST 请求到后端 API，开始录音
      const response = await fetch('http://127.0.0.1:7073/start_recording', {
        method: 'POST', // 请求方法为 POST
        headers: { 'Content-Type': 'application/json' }, // 设置请求头，表明发送的是 JSON 数据
        body: JSON.stringify({
          microphone: selectedMicrophoneDriver.value,//麦克风
        }) // 发送一个空的 JSON 对象作为请求体
      });
  
      // 检查响应状态，如果响应状态不是 2xx，则抛出错误
      if (!response.ok) {
        throw new Error('Start recording failed');
      }
  
      // 解析响应数据
      const data = await response.json();
      console.log('Start recording response:', data);
    } catch (error) {
      // 如果发生错误，提示用户开始录音失败
      message.error('开始录音失败');
      console.error('Start recording error:', error);
  
      // 提示用户检查网络连接和链接的合法性
      message.warning('请检查网络连接和链接的合法性，适当重试。');
      isRecording.value = false; // 将录音状态恢复为 false
    }
  };
  //录制结束
  const stopVoiceInput = async () => {
  // 如果当前没有在录音状态，直接返回，不做任何操作
  if (!isRecording.value) return;

  // 将录音状态设置为 false，表示停止录音
  isRecording.value = false;

  // 提示用户停止录音
  message.info('停止语音输入');

  try {
    // 发起 POST 请求到后端 API，停止录音
    const response = await fetch('http://127.0.0.1:7073/stop_recording', {
      method: 'POST', // 请求方法为 POST
      headers: { 'Content-Type': 'application/json' }, // 设置请求头，表明发送的是 JSON 数据
      body: JSON.stringify({})
    });

    // 检查响应状态，如果响应状态不是 2xx，则抛出错误
    if (!response.ok) {
      throw new Error('Stop recording failed');
    }

    // 解析响应数据
    const data = await response.json();
    console.log('Stop recording response:', data);

    // 检查返回的 message 字段是否表示成功
    if (data.message === "Recording stopped and processed") {
      // 如果返回的 message 表示成功，处理成功逻辑
      message.success('语音已转录'); // 提示用户语音已成功转录

      // 如果返回的 filename 存在，将其添加到播放列表
      if (data.filename) {
        playList.value.unshift({
          content: "[语音消息]", // 内容
          filename: data.filename, // 文件名
          play_mode: "serial" // 播放模式
        });
      }
    } else {
      // 如果返回的 message 不表示成功，抛出错误
      throw new Error('Stop recording failed');
    }
  } catch (error) {
    // 如果发生错误，提示用户停止录音失败
    message.error('停止录音失败');
    console.error('Stop recording error:', error);

    // 提示用户检查网络连接和链接的合法性
    message.warning('请检查网络连接和链接的合法性，适当重试。');
  }
};

  
  const handleKeyDown = (e: KeyboardEvent) => {
    if (
      document.activeElement?.tagName !== 'INPUT' &&
      document.activeElement?.tagName !== 'TEXTAREA'
    ) {
      if (e.code === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        sendUserMessage();
      }
    }
  };








const humanVoiceOptions = ref([]);
const get_human_voice_files_text_map = ref({});
const getHumanVoiceFiles = async () => {
  try {
    const response = await fetch('http://127.0.0.1:7073/get_human_voice_files', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })
    const data = await response.json()
    get_human_voice_files_text_map.value = data.text_map || {}
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
// const DisposableSendApi = async (model,parameters,user_content,system_prompt,supplier_name) => {
//   const response = await fetch('http://127.0.0.1:7072/disposable_send', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json',
//             },
//             body: JSON.stringify({ 
//               model: model, 
//               parameters: parameters,
//               user_content: user_content,
//               system_prompt: system_prompt,
//               supplier_name:supplier_name,
//               }),
//         });
//         console.log(response);
        
//         const data = await response.text();
//         if (data.includes('</think>')) {
//           // 去除 <think>
//           const endIndex = data.indexOf('</think>');
//           // 截取endIndex 过后的字符串
//           const result = data.substring(endIndex + 7); // 跳过 </think> 标签
//           const result2 = result.replace(/\n/g, '');
//           return result2;
//         }
//         return data;
// }

const DisposableSendApi = async (model, parameters, user_content, system_prompt, supplier_name) => {
    try {
        // 记录请求参数
        console.log('请求参数:', {
            model,
            parameters,
            user_content,
            system_prompt,
            supplier_name,
            trigger_conditions: ['BarrageComment'], // 确保触发弹幕
        });

        // 设置 60 秒超时
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 60000);

        const response = await fetch('http://127.0.0.1:7072/disposable_send', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                model,
                parameters,
                user_content,
                system_prompt,
                supplier_name,
                trigger_conditions: ['BarrageComment'], // 添加弹幕触发条件
            }),
            signal: controller.signal, // 支持超时中止
        });

        clearTimeout(timeoutId); // 清除超时定时器

        // 记录响应状态
        console.log('响应状态:', response.status, 'OK:', response.ok);

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`HTTP 错误: ${response.status}, 详情: ${errorText}`);
        }

        // 限制响应大小（10MB）
        const data = await response.text();
        if (data.length > 10 * 1024 * 1024) {
            throw new Error('响应数据过大');
        }

        // 处理 <think> 标签
        if (data.includes('</think>')) {
            const endIndex = data.indexOf('</think>');
            const result = data.substring(endIndex + 7).replace(/\n/g, '');
            return result;
        }

        return data;
    } catch (error) {
        console.error('请求失败:', error);
        throw error; // 抛出错误供测试捕获
    }
};



// const DisposableSendApi = async (model,parameters,user_content,system_prompt,supplier_name) => {
//     // 打印请求开始时间
//     const startTime = Date.now();
//     console.log('请求开始时间:', new Date(startTime).toLocaleString());
//     console.log('请求参数:111111111111111111111111111111111', {
//         model,
//         parameters,
//         user_content,
//         system_prompt,
//         supplier_name
//     });

//     try {
//         const response = await fetch('http://127.0.0.1:7072/disposable_send', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json',
//             },
//             body: JSON.stringify({
//                 model: model,
//                 parameters: parameters,
//                 user_content: user_content,
//                 system_prompt: system_prompt,
//                 supplier_name:supplier_name,
//             }),
//         });

//         // 打印响应状态
//         console.log('响应状态:', {
//             status: response.status,
//             statusText: response.statusText,
//             headers: Object.fromEntries(response.headers.entries())
//         });

//         const data = await response.text();
//         // 打印原始响应数据
//         console.log('原始响应数据:', data);

//         let result = data;
//         if (data.includes('</think>')) {
//             // 去除 <think>
//             const endIndex = data.indexOf('</think>');
//             // 截取endIndex 过后的字符串
//             result = data.substring(endIndex + 7); // 跳过 </think> 标签
//             result = result.replace(/\n/g, '');
//         }

//         // 打印处理后的结果
//         console.log('处理后的结果:', result);

//         // 计算并打印请求耗时
//         const endTime = Date.now();
//         console.log('请求结束时间:', new Date(endTime).toLocaleString());
//         console.log('请求总耗时:', endTime - startTime, 'ms');

//         return result;
//     } catch (error) {
//         // 打印错误信息
//         console.error('请求发生错误:', {
//             message: error.message,
//             stack: error.stack
//         });
//         throw error;
//     }
// }

// 模型调用API
// const DisposableSendApi = async (model,parameters,user_content,system_prompt,supplier_name) => {
//   try {
//     // 参数验证
//     if (!model) throw new Error('模型名称不能为空');
//     if (!user_content) throw new Error('用户内容不能为空');

//     console.log('开始调用模型API，参数：', {
//       model,
//       parameters,
//       user_content,
//       system_prompt,
//       supplier_name
//     });

//     // 设置超时时间为30秒
//     const controller = new AbortController();
//     const timeoutId = setTimeout(() => controller.abort(), 30000);

//     try {
//       const response = await fetch('http://127.0.0.1:7072/disposable_send', {
//         method: 'POST',
//         headers: {
//           'Content-Type': 'application/json',
//         },
//         body: JSON.stringify({ 
//           model: model, 
//           parameters: parameters,
//           user_content: user_content,
//           system_prompt: system_prompt,
//           supplier_name:supplier_name,
//         }),
//         signal: controller.signal
//       });

//       clearTimeout(timeoutId); // 清除超时计时器

//       // 检查响应状态
//       if (!response.ok) {
//         throw new Error(`API请求失败: ${response.status} ${response.statusText}`);
//       }

//       console.log('API响应状态：', response.status, response.statusText);
      
//       const data = await response.text();
//       if (!data) {
//         throw new Error('API返回数据为空');
//       }

//       console.log('API原始响应数据：', data);

//       if (data.includes('</think>')) {
//         // 去除 <think>
//         const endIndex = data.indexOf('</think>');
//         // 截取endIndex 过后的字符串
//         const result = data.substring(endIndex + 7); // 跳过 </think> 标签
//         const result2 = result.replace(/
// /g, '');
//         console.log('处理后的响应数据：', result2);
//         return result2;
//       }

//       console.log('未经处理的响应数据：', data);
//       return data;

//     } catch (error) {
//       clearTimeout(timeoutId); // 确保清除超时计时器
//       if (error.name === 'AbortError') {
//         throw new Error('请求超时，请稍后重试');
//       }
//       throw error; // 重新抛出其他错误
//     }

//   } catch (error) {
//     console.error('模型调用出错：', error);
//     throw new Error(`模型调用失败: ${error.message}`);
//   }
// }

// 替换文本 替换变量
const ReplaceText= async (text) => {
  console.log("原始文本：", text);
    console.log("当前进入直播间用户名：", EnterLiveRoomUserName.value); // 添加日志

  let newText = text
  if (text.includes('{语种}')){ //是否包含
    newText = newText.replace('{语种}', selectedLanguageLabel.value) //替换
  }
  if (text.includes('{送礼用户名}')){ //是否包含
    newText = newText.replace('{送礼用户名}', EnterGiftUserName.value) //替换
  }
  if (text.includes('{礼物名称}')){ //是否包含
    newText = newText.replace('{礼物名称}', EnterGiftGoodsName.value) //替换
  }
  if (text.includes('{礼物数量}')){ //是否包含
    newText = newText.replace('{礼物数量}', EnterGiftNum.value) //替换
  }
  if (text.includes('{点赞用户名}')){ //是否包含
    newText = newText.replace('{点赞用户名}', EnterSupportRoomUserName.value) //替换
  }
  if (text.includes('{进入直播间用户名}')){ //是否包含
    newText = newText.replace('{进入直播间用户名}', EnterLiveRoomUserName.value) //替换
  }
  if (text.includes('{分享直播间用户名}')){ //是否包含
    newText = newText.replace('{分享直播间用户名}', ShareRoomUserName.value);
  }
  if (text.includes('{关注直播用户名}')){ //是否包含
    newText = newText.replace('{关注直播用户名}', FollowRoomUserName.value);
  }
  if (text.includes('{弹幕内容}')){ //是否包含
    newText = newText.replace('{弹幕内容}', EnterBarrageContent.value) //替换
  }
  if (text.includes('{弹幕用户}')){ //是否包含
    newText = newText.replace('{弹幕用户}', EnterBarrageUserName.value) //替换
  }
  try {
    const response = await fetch('http://127.0.0.1:7073/get_combined_time_info', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })
    const data = await response.json()
    for (const key in data) {
      if (text.includes('{'+key+'}')){ //是否包含
        newText = newText.replace('{'+key+'}',data[key]) //替换
      }
    }
  } catch (error) {
  }
  console.log("替换后的文本：", newText); // 打印替换后的文本
  return newText;
} 






</script>

<style scoped lang="scss">
  .live-console {
    height: calc(100vh - 20px);
    padding-bottom: 20px;
    background-color: v-bind(themeThinkBg);
  }
  
  .message-voice-container {
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  
  .message-card, .voice-card {
    background: var(--n-color-embedded);
    border-radius: var(--n-border-radius);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
  }
  
  .message-card:hover, .voice-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
  
  .message-input {
    padding: 8px;
    background: #f8fafc;
    border-radius: 8px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }
  
  .message-input .n-input-group {
    display: flex;
    align-items: flex-start;
    gap: 8px;
  }
  
  .message-input .n-input {
    flex: 1;
    min-width: 0;
    border-radius: 6px;
  }
  
  .message-input .send-button {
    height: 50px;
    padding: 0 16px;
    border-radius: 6px;
    transition: all 0.3s ease;
  }
  
  .message-input .send-button:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  .message-input .send-button:active {
    transform: translateY(1px);
  }
  
  .voice-button {
    height: 40px;
    padding: 0 16px;
    border-radius: 6px;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100px;
    margin-right: 12px;
  }
  
  .voice-button:not(.recording):hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  .voice-button:not(.recording):active {
    transform: translateY(1px);
  }
  
  .voice-button.recording {
    animation: recordingPulse 1.2s infinite;
    background-color: #fef0f0;
    border-color: #fde2e2;
    color: #f56c6c;
  }
  
  @keyframes recordingPulse {
    0% {
      transform: scale(1);
      opacity: 1;
      box-shadow: 0 0 0 0 rgba(245, 108, 108, 0.4);
    }
    50% {
      transform: scale(1.05);
      opacity: 0.9;
      box-shadow: 0 0 0 8px rgba(245, 108, 108, 0.2);
    }
    100% {
      transform: scale(1);
      opacity: 1;
      box-shadow: 0 0 0 0 rgba(245, 108, 108, 0);
    }
  }
  
  .control-row {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    width: 100%;
  }
  
  .control-item {
    flex: 1;
    min-width: 200px;
    display: flex;
    align-items: center;
    gap: 8px;
  }
  
  .control-label {
    width: 80px;
    font-size: 12px;
    white-space: nowrap;
  }
  
  .control-slider {
    flex: 1;
  }
  
  .control-select {
    flex: 1;
  }
  
  .mic-select {
    display: flex;
    align-items: center;
    flex: 1;
  }
  
  .mic-label {
    margin-right: 8px;
    white-space: nowrap;
  }
  
  .mic-select-input {
    flex: 1;
    min-width: 200px;
  }
  
  .three-column-layout {
    padding: 0px 12px;
    display: grid;
    grid-template-columns: 2fr 3fr;
    gap: 8px;
  }
  
  .settings-bar {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 10px;
    background: var(--n-color-embedded);
    border-radius: var(--n-border-radius);
  }
  
  .control-block {
    padding: 8px;
    border-radius: 4px;
  }
  
  .message-card .n-card-header,
  .voice-card .n-card-header {
    border-bottom: 1px solid #e5e7eb;
    padding-bottom: 8px;
    font-weight: 600;
  }
  
  .control-block {
    padding: 12px;
  }
  
  .control-block div {
    font-size: 13px;
  }
  
  .live-console .three-column-layout .n-card {
    display: flex;
    flex-direction: column;
  }
  
  .live-console .three-column-layout .n-card-header {
    flex-shrink: 0;
    padding: 8px 12px;
  }
  
  .live-console .three-column-layout .n-card-content {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }
  
  .message-item {
    padding: 3px 0;
  }
  
  .message-input .mode-switch {
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .message-input .voice-button {
    flex: 1;
    padding: 0 16px;
    min-width: 0;
  }
  
  .message-input .voice-button span {
    margin-left: 8px;
    white-space: nowrap;
  }
  
  .control-block .block-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 4px;
  }
  
  .control-block .block-header h4 {
    font-size: 14px;
    margin: 0;
  }
  
  .control-block .block-content {
    padding-top: 4px;
  }
  
  .settings-bar .n-input-group {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  
  .mic-select {
    display: flex;
    align-items: center;
    flex-grow: 1;
  }
  </style>
