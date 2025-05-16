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
   <div class="live-console" :style="{backgroundColor: themeThinkBg}">
    <!-- 用户消息卡片 -->
    <div style="padding: 12px;">
      <n-card 
        title="用户消息" 
        :segmented="{content: true, footer: true}"
        header-style="padding:10px;font-size:14px"
        footer-style="padding:10px" 
        content-style="padding:10px;height:100%"
      >
        <!-- 控制区域：音量、语速、人声选择 -->
        <div class="control-block">
          <div style="display: flex; align-items: center; gap: 16px; width: 100%;">
            <div style="display: flex; align-items: center; gap: 8px; width: 30%;">
              <div style="width: 80px; font-size: 12px;">音量: {{ announcementVolume }}%</div>
              <n-slider
                v-model:value="announcementVolume"
                :min="10"
                :max="100"
                :step="1"
                style="flex: 1;"
                tooltip
                placement="bottom"
                :format-tooltip="(value) => `音量: ${value}%`"
                aria-label="消息音量调整"
              />
            </div>
            <div style="display: flex; align-items: center; gap: 8px; width: 30%;">
              <div style="width: 80px; font-size: 12px;">音速: {{ announcementSpeed }}x</div>
              <n-slider
                v-model:value="announcementSpeed"
                :min="0.5"
                :max="2"
                :step="0.01"
                style="flex: 1;"
                tooltip
                placement="bottom"
                :format-tooltip="(value) => `音速: ${value}x`"
                aria-label="消息语速调整"
              />
            </div>
            <div style="display: flex; align-items: center; gap: 8px; width: 40%;">
              <div style="width: 80px; font-size: 12px;">人声</div>
              <n-select 
                v-model:value="announcementVoice"
                :options="humanVoiceOptions"
                style="flex: 1;"
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
        <!-- 消息输入和语音控制区域 -->
        <div class="message-input">
          <div class="voice-shortcut" style="margin-bottom: 8px; display: flex; align-items: center; gap: 8px;">
            <n-button
              text
              size="small"
              @click="showShortcutModal = true"
              style="padding: 0;"
              aria-label="设置按键说话快捷键"
            >
              <template #icon>
                <i class="i-tdesign:keyboard w-16 h-16"></i>
              </template>
              按键说话快捷键：{{ voiceHotkey || '未设置' }}
            </n-button>
            <!-- 切换输入模式按钮 -->
            <n-button
              text
              size="small"
              @click="toggleInputMode"
              style="padding: 0;"
              aria-label="切换输入模式"
            >
              <template #icon>
                <i :class="isVoiceMode ? 'i-tdesign:keyboard w-16 h-16' : 'i-tdesign:microphone-1 w-16 h-16'"></i>
              </template>
              {{ isVoiceMode ? '文本输入' : '语音输入' }}
            </n-button>
          </div>
          <n-input-group>
            <!-- 文本输入模式 -->
            <n-input 
              v-if="!isVoiceMode"
              ref="messageInput"
              v-model:value="announcementMessage"
              type="textarea"
              :placeholder="`输入消息后按 Enter 发送`"
              :autosize="{ minRows: 2, maxRows: 5 }"
              @keyup.enter.ctrl.prevent="announcementMessage += '\n'"
              @keyup.enter.exact="sendUserMessage"
              style="flex: 1; min-width: 0;"
              aria-label="输入消息"
            />
            <!-- 语音输入模式 -->
            <n-button 
              v-if="isVoiceMode"
              :type="isRecording ? 'warning' : 'default'"
              :class="{ 'recording': isRecording }"
              style="flex: 1; min-width: 0; padding: 0 16px;"
              @click="handleVoiceButtonClick"
              @keydown.prevent.space="handleVoiceButtonClick"
              :aria-pressed="isRecording"
              :aria-label="isRecording ? '停止录音' : '开始录音'"
            >
              <template #icon>
                <i :class="isRecording ? 'i-tdesign:stop w-20 h-20' : 'i-tdesign:microphone-1 w-20 h-20'"></i>
              </template>
              <span>{{ isRecording ? '录音中' : `按 ${voiceHotkey || 'Space'} 录音` }}</span>
            </n-button>
            <!-- 发送按钮 -->
            <n-button 
              type="primary" 
              @click="sendUserMessage"
              :disabled="(!announcementMessage.trim() && !isVoiceMode) || !announcementVoice"
              aria-label="发送消息"
            >
              <template #icon>
                <i class="i-tdesign:send w-20 h-20"></i>
              </template>
              发送
            </n-button>
          </n-input-group>
          <!-- 快捷键设置弹窗 -->
          <n-modal v-model:show="showShortcutModal">
            <n-card style="width: 400px" title="设置按键说话快捷键" aria-modal="true">
              <n-input
                v-model:value="tempHotkey"
                placeholder="按下键盘按键以设置录音快捷键"
                readonly
                @keydown.stop.prevent="handleHotkeySet"
                aria-label="输入录音快捷键"
              />
              <template #footer>
                <div style="text-align: right;">
                  <n-button @click="showShortcutModal = false" style="margin-right: 8px;" aria-label="取消">取消</n-button>
                  <n-button type="primary" @click="saveHotkey" aria-label="确认">确定</n-button>
                </div>
              </template>
            </n-card>
          </n-modal>
        </div>
      </n-card>
    </div>
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
import { log } from 'mermaid/dist/logger.js'
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
const selectedAudioDriver = ref('')
const selectedSpeechModel = ref('')
const selectedSpeechModel2 = ref('')
const handleLanguageChange = async (value) => {
  try {
    const response = await fetch('http://192.168.1.10:7073/set_config', {
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
const loading2 = ref(false)
const initializeSpeechModel = async () => {
  // 加上logding
  if (selectedSpeechModel.value == "" || selectedSpeechModel2.value == "") {
    // message.error("请选择语音模型")
    return
  }
  loading2.value = true

  try {
   const response = await fetch('http://192.168.1.10:7074/set_model', {
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

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown);
  document.removeEventListener('keyup', handleKeyUp);
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
    const response = await fetch('http://192.168.1.10:7073/get_sound_cards', {
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
const speedModelOptions2 = ref([])

const getSpeedModels = async () => {
  try {
    const response = await fetch('http://192.168.1.10:7074/get_model_filenames', {
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
  fetchModules()
  getMames()
  getPromptRewrite()
  getAudioDevices()
  getHumanVoiceFiles()
  getSpeedModels()
  // 功能：组件挂载时设置事件监听器和初始焦点
  // 逻辑：添加 keydown 监听器，聚焦 textarea
  document.addEventListener('keydown', handleKeyDown);
  document.addEventListener('keyup', handleKeyUp);
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
  socket.value = io('ws://192.168.1.10:7073', {
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
          if (module.trigger_conditions.includes("SendGift")) {
            SendGift(module,newuuid)
          }
          //点赞 Like
          if (module.trigger_conditions.includes("Like")) {
            includesLike(module,newuuid)
          }
          //进入直播间 EnterLiveRoom
          if (module.trigger_conditions.includes("EnterLiveRoom")) {
            includesEnterLiveRoom(module,newuuid)
          }
        }
        
        //     TriggerShareRoom     = "ShareRoom"     // 分享直播间
//     TriggerFollowRoom    = "FollowRoom"    // 关注直播间
    }
    playListConsumption()
}
const intervalTime = ref(0) // 默认1000毫秒 1秒=1000毫秒
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
        while (playList.value.length >= 15) { //等待队列 消费完探入
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

// 进入直播间 模块
const EnterLiveRoomUserName = ref("")  //进入直播间用户名

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


//进入直播间 EnterLiveRoom
const includesEnterLiveRoom= async (module,newuuil) => {
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
           let prompt = await  ReplaceText(promptText.value) //提示词 赋值变量
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
            const  apidata =  await DisposableSendApi(
              model_api.value,
              parameters_api.value,
              ai_user_content, //文本
              prompt, //提示词
              supplierName_api.value, //供应商名称
            )
            console.log("apidata:",apidata);
            content = apidata
        if (module.retAi) { //是否改写
            if (model_api.value == "") { //是否改写
              console.log("请选择模型");
              await new Promise(resolve => setTimeout(resolve, 4000))
              continue
            }
           let prompt = await  ReplaceText(promptText.value) //提示词 赋值变量
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

// http://192.168.1.10:7073/generate_wav
// 生成wav文件
const temperature = ref(0.6)
const top_p = ref(0.8)
const top_k = ref(50)
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
        const response = await fetch('http://192.168.1.10:7074/generate_wav', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                // text: _text,
                // language: _language,
                // filename : _filename,
                // speaker_wav: _speaker_wav,
                // speed: _speed,// 1,
                // volume: _volume// 0.5
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
      // generate_wav_api_runing = false;
      return response.ok;
    }finally {
        releaseLock(); // 释放锁
    }
    
}










// 播放任务
const play_task_voice_api = async (_filename:string,play_mode:string) => {
    const response = await fetch('http://192.168.1.10:7073/play_task_voice', {
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

// --- 快捷键管理 ---
// 有效快捷键列表：限制用户可选的键，防止冲突
const validHotkeys = [
  'Space','Tab', 'Backspace',
  'KeyA', 'KeyB', 'KeyC', 'KeyD', 'KeyE', 'KeyF', 'KeyG', 'KeyH', 'KeyI', 'KeyJ',
  'KeyK', 'KeyL', 'KeyM', 'KeyN', 'KeyO', 'KeyP', 'KeyQ', 'KeyR', 'KeyS', 'KeyT',
  'KeyU', 'KeyV', 'KeyW', 'KeyX', 'KeyY', 'KeyZ',
  'Digit0', 'Digit1', 'Digit2', 'Digit3', 'Digit4', 'Digit5', 'Digit6', 'Digit7', 'Digit8', 'Digit9'
];


// 占位方法：处理录音按钮点击和快捷键触发
const handleVoiceButtonClick = () => {
  // 功能：切换录音状态并调用后台 API（由您实现）
  // 使用场景：点击录音按钮或按下快捷键时触发
  // 注意：isRecording 用于更新界面，您需根据 API 响应更新状态
    message.error("按下")
    return
  isRecording.value = !isRecording.value;
  // 示例实现（请替换为您的后台 API 调用）：
  // try {
  //   await fetch('http://your-api-endpoint/record', {
  //     method: 'POST',
  //     headers: { 'Content-Type': 'application/json' },
  //     body: JSON.stringify({ isRecording: isRecording.value })
  //   });
  //   message.info(isRecording.value ? '录音开始' : '录音停止');
  // } catch (error) {
  //   message.error('录音失败');
  //   isRecording.value = false;
  // }
};

const humanVoiceOptions = ref([])// 人声选项
const showShortcutModal = ref(false)// 快捷键设置弹窗的可见性
const tempHotkey = ref('')// 弹窗中临时存储的快捷键
const get_human_voice_files_text_map = ref({}) // 音色文件map
// 快捷键相关
const voiceHotkey = ref('Space'); // 当前快捷键，默认空格键
const isRecording = ref(false); // 录音状态
const handleVoiceButtonDown = () => {
  // 按下快捷键时触发
  isRecording.value = true;
  //message.info('开始录音');
};
const handleVoiceButtonUp = () => {
  // 松开快捷键时触发
  isRecording.value = false;
  if (announcementMessage.value.trim()) {
    message.warning('松开发送');
    sendVoiceMessage(); // 调用独立的后台发送方法
  } else {
    message.warning('请输入消息内容后再松开快捷键');
  }
};

const sendVoiceMessage = async () => {
  // 占位方法：独立的后台发送逻辑
  message.info('按键说话已发送：');
  
};


// textarea 引用，用于检测焦点状态
const messageInput = ref<HTMLInputElement | null>(null); // 明确指定类型为 HTMLInputElement 或 null
// --- 响应式状态 ---
// 消息输入相关
const announcementMessage = ref(''); // 用户输入的消息文本
const announcementVolume = ref(50); // 音量：10-100%，默认 50%
const announcementSpeed = ref(1.0); // 语速：0.5-2x，默认 1x
const announcementVoice = ref(''); // 选择的人声
//发送消息
const sendUserMessage = async () => {
  if (!announcementMessage.value.trim() || !announcementVoice.value) {
    message.error("发送消息请请输入消息并选择人声")
    return
  }
  try {
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
    
    playList.value.push({
      content: announcementMessage.value,
      filename: newFileName,
      play_mode: "serial"
    })
    message.success("消息已发送")
  } catch (error) {
    console.error("发送用户消息失败:", error)
    message.error("发送失败，请重试")
  }
  announcementMessage.value = ''
}





// --- 快捷键触发 ---
const handleKeyDown = (e: KeyboardEvent) => {
  // 功能：监听快捷键按下，触发录音或准备发送消息
  // 逻辑：仅在输入框未聚焦且弹窗未打开时，匹配 voiceHotkey 并调用 handleVoiceButtonDown
  if (document.activeElement !== messageInput.value && !showShortcutModal.value) {
    if (e.code === voiceHotkey.value && !e.ctrlKey && !e.altKey && !e.shiftKey) {
      e.preventDefault();
      handleVoiceButtonDown();
    }
  }
};
// --- 快捷键松开 ---
const handleKeyUp = (e: KeyboardEvent) => {
  // 功能：监听快捷键松开，触发发送消息
  // 逻辑：仅在输入框未聚焦且弹窗未打开时，匹配 voiceHotkey 并调用 handleVoiceButtonUp
  if (document.activeElement !== messageInput.value && !showShortcutModal.value) {
    if (e.code === voiceHotkey.value) {
      e.preventDefault();
      handleVoiceButtonUp();
    }
  }
};





// 处理弹窗中的快捷键设置
const handleHotkeySet = (e: KeyboardEvent) => {
  // 功能：捕获用户按下的键，验证是否有效
  // 逻辑：检查键是否在 validHotkeys 中，更新 tempHotkey
  // 注意：使用 e.code 而非 e.key，确保跨浏览器一致性
  const key = e.code;
  if (validHotkeys.includes(key)) {
    tempHotkey.value = key;
  } else {
    message.warning('请选择有效快捷键（如字母、数字）');
  }
};

// 保存快捷键
const saveHotkey = () => {
  // 功能：将临时快捷键保存为当前快捷键，关闭弹窗
  // 逻辑：检查 tempHotkey 是否非空，更新 voiceHotkey
  // 注意：可扩展为将快捷键保存到本地存储或后台
  if (tempHotkey.value) {
    voiceHotkey.value = tempHotkey.value;
    showShortcutModal.value = false;
    message.success(`快捷键已设置为 ${tempHotkey.value}`);
    tempHotkey.value = '';
  } else {
    message.error('请先选择一个快捷键');
  }
};

const getHumanVoiceFiles = async () => {
  try {
    const response = await fetch('http://192.168.1.10:7073/get_human_voice_files', {
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
const ReplaceText= async (text) => {
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
  if (text.includes('{分享直播间用户名}')){ //是否包含
    // newText = newText.replace('{分享直播间用户名}', selectedLanguageLabel.value) //替换
  }
  if (text.includes('{关注直播用户名}')){ //是否包含
    // newText = newText.replace('{关注直播用户名}', selectedLanguageLabel.value) //替换
  }
  if (text.includes('{弹幕内容}')){ //是否包含
    newText = newText.replace('{弹幕内容}', EnterBarrageContent.value) //替换
  }
  if (text.includes('{弹幕用户}')){ //是否包含
    newText = newText.replace('{弹幕用户}', EnterBarrageUserName.value) //替换
  }
  try {
    const response = await fetch('http://192.168.1.10:7073/get_combined_time_info', {
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

  return newText;
} 






</script>

<style scoped lang="scss">
.voice-controls .recording {
  animation: pulse 1.5s infinite;
  background-color: #fef0f0;
  border-color: #fde2e2;
  color: #f56c6c;
  position: relative;
}

@keyframes pulse {
  0% { 
    opacity: 1;
    transform: scale(1);
  }
  50% { 
    opacity: 0.8;
    transform: scale(1.03);
  }
  100% { 
    opacity: 1;
    transform: scale(1);
  }
}

.voice-controls .n-button {
  transition: all 0.3s ease;
}

.voice-controls .n-button:not(.recording):hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.voice-controls .n-button:not(.recording):active {
  transform: translateY(1px);
}

.gd-bg{
  background-color: v-bind(themeThinkBg);
}
// .n-slider-handle {
//   --n-handle-size: 5px !important;
//   width: var(--n-handle-size) !important;
//   height: var(--n-handle-size) !important;
//   transform: translateX(-50%) translateY(-50%) !important;
// }
.live-console {
  // height: calc(100% - 130px);
  height: calc(100vh - 20px); padding-bottom: 20px;
  .three-column-layout {
    padding: 0px  12px;
    display: grid;
    // grid-template-columns: repeat(3, 1fr);
    grid-template-columns: 2fr 3fr;  // 中间列宽度是两侧的1.5倍
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
  // padding: 12px;
  background: var(--n-color-embedded);
  border-radius: var(--n-border-radius);
  
  .n-input-group {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

</style>
