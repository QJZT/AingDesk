<template>
    <div class="layout-sider-wrapper">
        <!-- logo部分 -->
        <div class="logo">
            <div class=logo-left>
                <!-- <n-image :src="logo" object-fit="cover" class="h-30" preview-disabled /> -->
                <span class="text-transparent bg-clip-text bg-gradient-to-r from-purple-500 via-pink-500 to-red-500 font-bold text-xl drop-shadow-lg">AI 无人直播</span>
            </div>
            <!-- <div>
                <i class="i-common:fold w-18 h-18  cursor-pointer" @click="doFold"></i>
            </div> -->
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
            <n-button type="primary" ghost style="width:100%" @click="">
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
const { t: $t } = useI18n()

const { siderBg } = getGlobalStore()
const {
    themeMode,
} = getSoftSettingsStoreData()
const { currentView } = getSiderStoreData()
/**
 * @description 计算不同模式下logo的图片
 */

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