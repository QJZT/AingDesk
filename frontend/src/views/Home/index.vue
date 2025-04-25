<template>
    <n-layout :has-sider="true" class="layout-wrapper">

        <!-- 侧边栏 -->
        <n-layout-sider :width="siderWidth" class="layout-sider">
            <Sider/>
        </n-layout-sider>


        <n-layout-sider :class="['layout-sider', { 'no-border': knowledgeSiderWidth == 0 }]"
            :width="knowledgeSiderWidth">
            <KnowledgeStore />
        </n-layout-sider>

        
        <n-layout>
            <!-- 页面A -->
             <!-- {{ currentView }} -->
             <!-- {{ currentView == 'LiveConsole' }} -->
            <n-layout-header class="layout-header" v-if="currentView == 'ChatContent'">
                <Header />
            </n-layout-header>
            <n-layout-content v-if="currentView === 'ChatContent'"  class="layout-content" style="padding:0">
                <ChatContent v-if="currentView === 'ChatContent'" />
            </n-layout-content>
            <LiveConsole v-show="currentView == 'LiveConsole'" />
            <ModuleConfig v-if="currentView == 'ModuleConfig'" />
        </n-layout>
    </n-layout>

    <!-- 欢迎界面 -->
    <Welcome />
    <!-- 第三方api弹窗 -->
    <ThirdPartyApi />

    <Settings />

    <!-- 更换目录时数据迁移进度 -->
    <n-modal v-model:show="dataPathChangeCheckShow">
        <n-card class="w-600">
            <template #header>
                <div class="flex justify-start items-center gap-1.25">
                    <i class="i-tdesign:error-circle w-20 h-20 text-[#E6A23C]"></i><span class="text-4">{{
                        pathChangeNotice }},{{ $t("请勿关闭软件,否则将导致软件异常") }}</span>
                </div>
            </template>
            <div class="flex justify-between items-center mb-5">
                <span>{{ dataPathChangeStatusValues.message || $t("等待复制") }}</span>
                <span>[{{ dataPathChangeStatusValues.fileCurrent }}/{{ dataPathChangeStatusValues.fileTotal }}]</span>
            </div>
            <n-progress type="line" :percentage="dataPathChangeStatusValues.percent" indicator-placement="inside"
                processing>
            </n-progress>
            <div class="flex justify-between items-center mt-20">
                <span>{{ $t("总大小") }}: {{ getByteUnit(dataPathChangeStatusValues.total) }}</span>
                <span>{{ $t("已复制") }}: {{ getByteUnit(dataPathChangeStatusValues.current) }}</span>
                <span>{{ $t("速度") }}: {{ getByteUnit(dataPathChangeStatusValues.speed) }}/s</span>
            </div>
        </n-card>
    </n-modal>

    <!-- 软件设置 -->
    <SoftSettings />
</template>

<script setup lang="ts">
import SoftSettings from "@/views/SoftSettings/index.vue";
import Settings from '@/views/Settings/index.vue';
import Sider from "../Sider/index.vue";
import Header from "@/views/Header/index.vue";
import ChatContent from "@/views/ChatContent/index.vue";
import ModuleConfig from '@/views/moduleConfig/index.vue'
import LiveConsole from '@/views/LiveConsole/index.vue'
import KnowledgeStore from "@/views/KnowleadgeStore/index.vue";
import Welcome from "./components/Welcome.vue";
import ThirdPartyApi from "@/views/ThirdPartyApi/index.vue";
import { get_languages } from "@/views/SoftSettings/controller"
import { useI18n } from "vue-i18n";
import { computed } from "vue";
import { getByteUnit } from "@/utils/tools"
import { getSiderStoreData } from "../Sider/store/index.ts";
import { getKnowledgeStoreData } from "../KnowleadgeStore/store/index.ts";
import { getSoftSettingsStoreData } from "../SoftSettings/store/index.ts";
import { nextTick } from 'vue';

const { siderWidth,currentView } = getSiderStoreData()
const { knowledgeSiderWidth, } = getKnowledgeStoreData()
const { dataPathChangeCheckShow, dataPathChangeStatusValues } = getSoftSettingsStoreData()
const { t: $t } = useI18n()
/**
 * @description 获取支持的语言和语言列表
 */
get_languages()

// 数据迁移文案计算属性
const pathChangeNotice = computed(() => {
    if (dataPathChangeStatusValues.value.status == 0) {
        return $t("准备迁移数据")
    } else if (dataPathChangeStatusValues.value.status == 1) {
        return $t("迁移数据中")
    }
})


// const currentView = ref('ChatContent') // 默认显示聊天内容
// const switchView = async (viewName: string) => {
//     currentView.value = '';
//     await nextTick();
//     currentView.value = viewName;
// }
// const checkService = async () => {
//   try {
//     const response = await fetch('http://127.0.0.1:7072/ping')
//     if (!response.ok) {
//       throw new Error('服务未正常启动')
//     }
//     console.log('服务已正常启动')
//   } catch (error) {
//     console.error('服务检查失败:', error)
//   }
// }
// checkService();

// const switchView = async (viewName: string) => {
//   await checkService()
//   currentView.value = ''
//   await nextTick()
//   currentView.value = viewName
// }

</script>

<style scoped lang="scss">
@use "@/assets/base";

.layout-wrapper {
    height: 100vh;
    background: var(--bt-bg);

    .layout-sider {
        transition: .2s;
        border-right: 1px solid rgba(0, 0, 0, .08);
    }

    .layout-content {
        padding: var(--bt-pd-normal);
        height: calc(100% - 67px);
    }

    .layout-header {
        padding: 0 var(--bt-pd-normal);
        border: 1px solid rgba(0, 0, 0, .12);
        border-left: none;
        border-top: none;
    }
}

.no-border {
    border-right: none !important
}
</style>