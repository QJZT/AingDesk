<template>
  <div class="tutorial-docs">
    <n-card title="教程文档">
      <div class="docs-container">
        <n-scrollbar style="max-height: calc(100vh - 200px)">
          <div class="docs-content">
            <!-- 快速开始 -->
            <div class="doc-section">
              <h2>🚀 快速开始</h2>
              <div class="doc-item">
                <h3>1. 模型安装</h3>
                <div class="config-section">
                  <img src="@/assets/images/tutorial/model-install.png" alt="模型安装界面" class="tutorial-image" />
                  <div class="config-content">
                    <p>在开始使用前，需要先安装必要的AI模型和依赖。</p>
                    <ul>
                      <li><strong>语言模型：</strong>下载并配置大语言模型（如ChatGLM、Qwen等）</li>
                      <li><strong>语音模型：</strong>安装TTS（文本转语音）和ASR（语音识别）模型</li>
                      <li><strong>数字人模型：</strong>配置面部表情和口型同步模型</li>
                      <li><strong>依赖环境：</strong>确保Python环境和相关依赖包已正确安装</li>
                    </ul>
                    <div class="install-tips">
                      <n-alert type="info" title="安装提示">
                        <p>• 首次安装可能需要较长时间，请耐心等待</p>
                        <p>• 建议使用GPU加速以获得更好的性能</p>
                        <p>• 确保有足够的磁盘空间（至少10GB）</p>
                      </n-alert>
                    </div>
                  </div>
                </div>
              </div>
              
              <div class="doc-item">
                <h3>2. 基础配置</h3>
                <div class="config-section">
                  <img src="@/assets/images/tutorial/basic-config.png" alt="基础配置界面" class="tutorial-image" />
                  <div class="config-content">
                    <p>首先进入「模块配置」页面，配置基础模块和音频文件模块。</p>
                    <ul>
                      <li>设置文案内容和触发条件</li>
                      <li>配置音频文件路径</li>
                      <li>设定触发时间间隔</li>
                    </ul>
                  </div>
                </div>
              </div>
              
              <div class="doc-item">
                <h3>3. 直播控制</h3>
                <div class="config-section">
                  <img src="@/assets/images/tutorial/live-control-1.png" alt="直播控制台界面" class="tutorial-image" />
                  <img src="@/assets/images/tutorial/live-control-3.png" alt="直播控制台界面1" class="tutorial-image" />
                  <div class="config-content">
                    <p>在「直播控制台」中可以实时控制直播状态。</p>
                    <ul>
                      <li>启动/停止直播</li>
                      <li>实时监控直播数据</li>
                      <li>手动触发特定模块</li>
                    </ul>
                  </div>
                </div>
                
                <div class="config-section">
                  <img src="@/assets/images/tutorial/live-control-2.png" alt="提示词设置界面" class="tutorial-image" />
                  <div class="config-content">
                    <p>通过提示词设置功能，可以自定义AI的回复风格和内容。</p>
                    <ul>
                      <li>自定义AI角色设定</li>
                      <li>配置回复风格和语调</li>
                      <li>设置特定场景的回复模板</li>
                    </ul>
                  </div>
                </div>
              </div>
              
              <div class="doc-item">
                <h3>4. 数据监控</h3>
                <p>通过「数据面板」查看详细的直播数据和统计信息。</p>
                <div class="config-section">
                  <img src="@/assets/images/tutorial/data-monitoring.png" alt="数据监控面板" class="config-image" />
                  <div class="config-text">
                    <p>数据监控面板提供全面的直播数据分析：</p>
                    <ul>
                      <li>观看人数统计</li>
                      <li>互动数据分析</li>
                      <li>收益数据报表</li>
                      <li>实时流量监控</li>
                      <li>用户行为分析</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            <!-- 功能介绍 -->
            <div class="doc-section">
              <h2>📋 功能介绍</h2>
              
              <div class="doc-item">
                <h3>基础模块</h3>
                <p>文案配置是AI直播的核心功能，支持自动化生成对应文案。</p>
                <div class="feature-list">
                  <n-tag type="info">自动文案生成</n-tag>
                  <n-tag type="success">AI工作流对接</n-tag>
                  <n-tag type="warning">触发条件设置</n-tag>
                </div>
              </div>
              
              <div class="doc-item">
                <h3>音频文件模块</h3>
                <p>设定音频触发对接程序，实现音频与直播内容的同步。</p>
                <div class="feature-list">
                  <n-tag type="info">音频文件管理</n-tag>
                  <n-tag type="success">自动播放控制</n-tag>
                  <n-tag type="warning">音频质量优化</n-tag>
                </div>
              </div>
              
              <div class="doc-item">
                <h3>数字人模块</h3>
                <p>数字人对口型对接，实现更自然的直播效果。</p>
                <div class="feature-list">
                  <n-tag type="info">口型同步</n-tag>
                  <n-tag type="success">表情控制</n-tag>
                  <n-tag type="warning">动作配置</n-tag>
                </div>
              </div>
            </div>

            <!-- 常见问题 -->
            <div class="doc-section">
              <h2>❓ 常见问题</h2>
              
              <n-collapse>
                <n-collapse-item title="如何配置第一个直播模块？" name="1">
                  <p>1. 进入「模块配置」页面</p>
                  <p>2. 点击「基础模块」的 + 号</p>
                  <p>3. 填写模块名称和文案内容</p>
                  <p>4. 设置触发条件和时间间隔</p>
                  <p>5. 保存配置即可开始使用</p>
                </n-collapse-item>
                
                <n-collapse-item title="音频文件支持哪些格式？" name="2">
                  <p>目前支持以下音频格式：</p>
                  <ul>
                    <li>MP3 - 推荐使用</li>
                    <li>WAV - 高质量音频</li>
                    <li>AAC - 压缩音频</li>
                    <li>OGG - 开源格式</li>
                  </ul>
                </n-collapse-item>
                
                <n-collapse-item title="如何优化直播效果？" name="3">
                  <p>建议从以下几个方面优化：</p>
                  <ul>
                    <li>合理设置触发时间间隔，避免过于频繁</li>
                    <li>准备多样化的文案内容，提高观众体验</li>
                    <li>定期检查音频质量，确保清晰度</li>
                    <li>根据数据面板反馈调整策略</li>
                  </ul>
                </n-collapse-item>
                
                <n-collapse-item title="数据统计多久更新一次？" name="4">
                  <p>数据统计实时更新，包括：</p>
                  <ul>
                    <li>观看人数：每5秒更新</li>
                    <li>互动数据：实时更新</li>
                    <li>收益统计：每分钟更新</li>
                  </ul>
                </n-collapse-item>
              </n-collapse>
            </div>

            <!-- 联系支持 -->
            <div class="doc-section">
              <h2>📞 技术支持</h2>
              <div class="support-info">
                <p>如果您在使用过程中遇到问题，可以通过以下方式获取帮助：</p>
                <div class="contact-methods">
                  <n-card class="contact-card">
                    <div class="contact-item">
                      <i class="i-mdi:email w-20 h-20"></i>
                      <div>
                        <h4>邮件支持</h4>
                        <p>support@example.com</p>
                      </div>
                    </div>
                  </n-card>
                  
                  <n-card class="contact-card">
                    <div class="contact-item">
                      <i class="i-mdi:wechat w-20 h-20"></i>
                      <div>
                        <h4>微信群</h4>
                        <p>扫码加入用户交流群</p>
                      </div>
                    </div>
                  </n-card>
                  
                  <n-card class="contact-card">
                    <div class="contact-item">
                      <i class="i-mdi:github w-20 h-20"></i>
                      <div>
                        <h4>GitHub</h4>
                        <p>提交Issue或查看文档</p>
                      </div>
                    </div>
                  </n-card>
                </div>
              </div>
            </div>
          </div>
        </n-scrollbar>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// 可以在这里添加一些交互逻辑
const activeCollapse = ref(['1'])
</script>

<style scoped lang="scss">
.tutorial-docs {
  height: 100vh;
  padding: 20px;
  background: var(--n-color);
  
  .docs-container {
    .docs-content {
      .doc-section {
        margin-bottom: 40px;
        
        h2 {
          color: var(--n-text-color);
          font-size: 24px;
          font-weight: bold;
          margin-bottom: 20px;
          padding-bottom: 10px;
          border-bottom: 2px solid var(--n-border-color);
        }
        
        .doc-item {
          margin-bottom: 30px;
          padding: 20px;
          background: var(--n-card-color);
          border-radius: 8px;
          border: 1px solid var(--n-border-color);
          
          h3 {
            color: var(--n-text-color);
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 10px;
          }
          
          p {
            color: var(--n-text-color-2);
            line-height: 1.6;
            margin-bottom: 10px;
          }
          
          ul {
            color: var(--n-text-color-2);
            padding-left: 20px;
            
            li {
              margin-bottom: 5px;
              line-height: 1.5;
            }
          }
          
          .feature-list {
            display: flex;
            gap: 8px;
            flex-wrap: wrap;
            margin-top: 10px;
          }
        }
      }
      
      .config-section {
         display: flex;
         flex-direction: column;
         gap: 15px;
         margin: 15px 0;
       }
       
       .tutorial-image {
          width: 100%;
          height: auto;
          border-radius: 8px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }
       
       .config-content {
          width: 100%;
        }
        
        .install-tips {
          margin-top: 15px;
          
          .n-alert {
            p {
              margin: 5px 0;
              font-size: 14px;
            }
          }
        }
      
      .support-info {
        .contact-methods {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
          gap: 20px;
          margin-top: 20px;
          
          .contact-card {
            .contact-item {
              display: flex;
              align-items: center;
              gap: 15px;
              
              i {
                color: var(--n-primary-color);
              }
              
              h4 {
                color: var(--n-text-color);
                font-size: 16px;
                font-weight: 600;
                margin-bottom: 5px;
              }
              
              p {
                color: var(--n-text-color-2);
                font-size: 14px;
                margin: 0;
              }
            }
          }
        }
      }
    }
  }
}

@media (max-width: 768px) {
   .tutorial-docs {
     padding: 15px;
   }
 }
</style>