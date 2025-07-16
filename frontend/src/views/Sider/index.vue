<template>
    <div class="layout-sider-wrapper">
        <!-- logo部分 -->
        <div class="logo">
            <div class=logo-left>
                <!-- <n-image :src="logo" object-fit="cover" class="h-30" preview-disabled /> -->
                <span class="text-transparent bg-clip-text bg-gradient-to-r from-purple-500 via-pink-500 to-red-500 font-bold text-xl drop-shadow-lg">AI 无人直播</span>
            </div>
        </div>
        
        <!-- 激活码有效期显示 -->
        <div class="flex justify-center items-center mt-2" style="margin-top: 8px;">
            <span class="text-sm text-gray-600">有效期：</span>
            <span class="text-sm ml-2" :class="expiryTextClass">{{ expiryText }}</span>
            
        </div>
        <!-- 新建对话按钮 -->
        <div class="flex justify-center items-center">
            <n-button type="primary" ghost style="width:100%"  @click="setCurrentView('LiveConsole')">
                <!-- <template #icon>
                    <i class="i-tdesign:chat-add w-16 h-16"></i>
                </template> -->
                直播控制台
            </n-button>
        </div>

         <!-- 新建对话按钮 -->
         <div class="flex justify-center items-center">
            <n-button type="primary" ghost style="width:100%" @click="setCurrentView('ModuleConfig')">
                <!-- <template #icon>
                    <i class="i-tdesign:chat-add w-16 h-16"></i>
                </template> -->
               模块配置
            </n-button>
        </div>

        <!-- 新建对话按钮 -->
        <div class="flex justify-center items-center">
            <n-button type="primary" ghost style="width:100%"@click="setCurrentView('DataListView')">
                <!-- <template #icon>
                    <i class="i-tdesign:chat-add w-16 h-16"></i>
                </template> -->
                数据面板
            </n-button>
        </div>
        <div class="flex justify-center items-center">
            <n-button type="primary" ghost style="width:100%" @click="setCurrentView('TutorialDocs')">
                <!-- <template #icon>
                    <i class="i-tdesign:chat-add w-16 h-16"></i>
                </template> -->
                教程文档
            </n-button>
        </div>
         <div class="sider-divider"></div>

        <div class="recent-header">
            <span class="ml-8 text-[var(--bt-notice-text-color)]"> LLM模型测试
            </span>
            <n-popover trigger="hover">
                <template #trigger>
                    <i class="i-ci:bar-top w-14 h-14 cursor-pointer" @click="cleanAllChats"></i>
                </template>
                {{ $t("清空对话") }}
            </n-popover>
        </div>

        <div class="sider-wrapper" style="overflow: hidden;" >
            <n-scrollbar :style="{ height: '100%' }">
                <div class="sider-top">
                    <ChatList />
                </div>
            </n-scrollbar>
            <!-- <div class="sider-divider"></div> -->
        </div>
        
        <!-- <div class="recent-header">
            <span class=" text-[var(--bt-notice-text-color)] flex justify-start items-center ml-10">{{ $t("知识库")
                }}</span>
        </div> -->

        <!-- 知识库 -->
        <!-- <div class="sider-wrapper" style="overflow: hidden; gap:10px">
            <n-scrollbar :style="{ height: '100%' }">
                <div class="sider-top">
                    <KnowledgeList />
                </div>
            </n-scrollbar>
            <div class="sider-divider"></div>
        </div> -->
        
        
        <!-- <div class="recent-header">
            <span class=" text-[var(--bt-notice-text-color)] flex justify-start items-center ml-10">
                {{ $t("知识库")}} 
            </span>
        </div> -->
        <n-tabs type="line" animated :size="'small'" justify-content="space-evenly" style=""
        @update:value="handleUpdateValue">
            <n-tab-pane name="zsk" tab="知识库" >
            </n-tab-pane>
            <n-tab-pane name="ypk" tab="人声库">
            </n-tab-pane>
        </n-tabs>

        <!-- 知识库 -->
        <div class="sider-wrapper" v-show="tab_pane == 'zsk'" style="overflow: hidden; gap:10px">
            <n-scrollbar :style="{ height: '100%' }">
                <div class="sider-top">
                    <KnowledgeList />
                </div>
            </n-scrollbar>
        </div>
        <div class="sider-wrapper" v-show="tab_pane == 'ypk'" style="overflow: hidden; gap:10px">
            <n-scrollbar :style="{ height: '100%' }">
                <div class="sider-top">
                    <AudioLibraryList />
                </div>
            </n-scrollbar>
            <!-- <n-scrollbar :style="{ height: '100%' }">
                <div class="sider-top">
                    <KnowledgeList />
                </div>
            </n-scrollbar> -->
        </div>

        <!-- 侧边栏下部分 -->
        <div class="sider-bottom">
            <div class="sider-divider"></div>
            <SiderBottom />
        </div>
    </div>

    <!-- 删除对话问询 -->
    <RemoveChatConfirm />

    <!-- 修改对话标题确认框 -->
    <ModifyChatConfirm />

    <!-- 智能体 -->
    <Agent />

    <!-- 优化知识库进度 -->
    <OptimizeProgress />
</template>

<script setup lang="ts">
import { getGlobalStore } from "@/stores/global"
import OptimizeProgress from "@/views/KnowleadgeStore/components/OptimizeProgress.vue"
import ChatList from "./components/ChatList.vue"
import KnowledgeList from "./components/KnowledgeList.vue";
import AudioLibraryList from "./components/AudioLibraryList.vue";
import SiderBottom from "./components/SiderBottom.vue";
import RemoveChatConfirm from "./components/RemoveChatConfirm.vue";
import ModifyChatConfirm from "./components/ModifyChatConfirm.vue";
import { get_chat_list, makeNewChat, doFold,cleanAllChats } from "@/views/Sider/controller"
import { getSoftSettingsStoreData } from "../SoftSettings/store";
import Agent from "@/views/Agent/index.vue";
import logoImage from "@/assets/images/logo.png"
import logoDark from "@/assets/images/logo-dark.png"
import { getSiderStoreData } from "../Sider/store/index.ts";
import { useI18n } from "vue-i18n";
import { set } from "@vueuse/core";
import { knowledgeIsClose } from "@/views/KnowleadgeStore/controller"
import { useDialog } from "@/utils/naive-tools"
import { onMounted, onUnmounted, ref, computed } from 'vue'

const { t: $t } = useI18n()

const { siderBg } = getGlobalStore()
const {
    themeMode,
} = getSoftSettingsStoreData()
const { currentView } = getSiderStoreData()

// 激活码有效期相关状态
const expiryText = ref('检查中...')
const expiryDays = ref(0)
const checkTimer = ref<NodeJS.Timeout | null>(null)

// 计算有效期文本样式
const expiryTextClass = computed(() => {
    if (expiryDays.value <= 0) {
        return 'text-red-500 font-bold'
    } else if (expiryDays.value <= 7) {
        return 'text-orange-500 font-bold'
    } else {
        return 'text-green-600'
    }
})

/**
 * 检查激活码有效期
 */
/**
 * 检查激活码有效期
 */
/**
 * 检查激活码有效期
 */
/**
 * 检查激活码有效期
 */
/**
 * 检查激活码有效期
 */
async function checkActivationExpiry() {
    try {
        const response = await fetch('http://127.0.0.1:7072/get_kv', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ key: 'client_config' })
        })
        
        if (response.ok) {
            const data = await response.json()
            console.log('API返回完整数据:', data)
            
            // 根据实际返回的数据结构处理
            if (data.kv && data.kv.expires_at) {
                const expiresAt = data.kv.expires_at
                console.log('解析的expires_at:', expiresAt)
                
                const expiryDate = new Date(expiresAt)
                const now = new Date()
                const diffTime = expiryDate.getTime() - now.getTime()
                const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
                
                expiryDays.value = diffDays
                console.log('计算的剩余天数:', diffDays)
                
                if (diffTime <= 0) {
                    expiryText.value = '已过期'
                    showExpiredWarning()
                } else if (diffTime < 24 * 60 * 60 * 1000) {
                    // 不足一天时显示详细时间
                    const diffHours = Math.floor(diffTime / (1000 * 60 * 60))
                    const diffMinutes = Math.floor((diffTime % (1000 * 60 * 60)) / (1000 * 60))
                    
                    if (diffHours > 0) {
                        expiryText.value = `${diffHours}小时${diffMinutes}分钟后过期`
                    } else {
                        expiryText.value = `${diffMinutes}分钟后过期`
                    }
                    showExpiryWarning(diffDays)
                } else if (diffDays <= 7) {
                    expiryText.value = `${diffDays}天后过期`
                    showExpiryWarning(diffDays)
                } else {
                    expiryText.value = `${diffDays}天后过期`
                }
            } else {
                console.log('未找到kv.expires_at字段')
                expiryText.value = '无法获取有效期'
            }
        } else {
            console.log('HTTP请求失败:', response.status, response.statusText)
            expiryText.value = '检查失败'
        }
    } catch (error) {
        console.error('检查激活码有效期失败:', error)
        expiryText.value = '检查失败'
    }
}

/**
 * 显示即将过期警告
 */
function showExpiryWarning(days: number) {
    const dialogInstance = useDialog({
        title: '激活码即将过期',
        content: `您的激活码将在 ${days} 天后过期，请及时续费以免影响使用。`,
        positiveText: '我知道了',
        style: {
            width: '400px',
            maxWidth: '90vw'
        },
        closable: true,
        maskClosable: true,
        onOk: () => {
            // 关闭弹窗
            dialogInstance.destroy()
        }
    })
}

/**
 * 显示已过期警告
 */
function showExpiredWarning() {
    const dialogInstance = useDialog({
        title: '激活码已过期',
        content: '您的激活码已过期，请联系管理员续费。',
        positiveText: '我知道了',
        style: {
            width: '400px',
            maxWidth: '90vw'
        },
        closable: true,
        maskClosable: true,
        onOk: () => {
            // 关闭弹窗
            dialogInstance.destroy()
        }
    })
}

/**
 * 启动定时检查
 */
function startPeriodicCheck() {
    // 立即检查一次
    checkActivationExpiry()
    
    // 每小时检查一次
    checkTimer.value = setInterval(() => {
        checkActivationExpiry()
    }, 60 * 60 * 1000)
}

/**
 * 停止定时检查
 */
function stopPeriodicCheck() {
    if (checkTimer.value) {
        clearInterval(checkTimer.value)
        checkTimer.value = null
    }
}

// 组件挂载时启动检查
onMounted(() => {
    startPeriodicCheck()
})

// 组件卸载时清理定时器
onUnmounted(() => {
    stopPeriodicCheck()
})

const setCurrentView = (viewName: string) => {
    set(currentView, viewName)
    knowledgeIsClose()
}
const logo = computed(() => {
    if (themeMode.value == "light") {
        return logoImage
    } else {
        return logoDark
    }
})

/**
 * @description 获取历史对话列表
 */
get_chat_list()

const tab_pane = ref("zsk")
const handleUpdateValue =  (value: string) => {
        tab_pane.value = value;
        console.log(tab_pane.value);

}

</script>

<style scoped lang="scss">
@use "@/assets/base";

.layout-sider-wrapper {
    display: grid;
    // grid-template-rows: 50px 22px 2fr 22px 1fr 140px;
    grid-template-rows: 50px 50px 50px 50px 50px 25px 22px 1fr 35px  1fr 110px;
    height: 100%;
    box-sizing: border-box;
    padding: var(--bt-pd-small);
    background: v-bind(siderBg);

    .logo {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .logo-left {
            display: flex;
            gap: 10px;
            justify-content: flex-start;
            align-items: center;

            span {
                font-size: 18px;
                font-weight: bold;
            }
        }
    }

    .sider-wrapper {
        --hover-bg: #e3e3e3;
        --item-border-radius: 10px;

        box-sizing: border-box;
        @include base.row-flex-between;
        gap: var(--bt-mg-normal);
        flex-direction: column;

        .sider-top {
            width: 100%;
        }
    }

    .sider-bottom {
        width: 100%;
    }
}





.create-comunication {
    @include base.comu-list-item;
}

.recent-header {
    @include base.row-flex-between;
    margin: var(--bt-mg-small) 0;
    box-sizing: border-box;
}

.sider-divider {
    height: 1px;
    width: calc(100% - var(--bt-pd-normal)*2);
    background-color: rgba(0, 0, 0, .12);
    margin: var(--bt-mg-small) auto;
}

.knowledge-store-list {
    .create-comunication {
        @include base.comu-list-item;

        .comu-tit {
            @include base.comu-tit;
        }
    }

    .knowledge-list {
        @include base.recent-list-style;
    }
}
</style>