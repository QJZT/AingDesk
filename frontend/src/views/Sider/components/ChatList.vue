<template>
    <div class="recent-comunication">
        <!-- 最近对话列表 -->
        <ul class="recent-list">
            <li class="create-comunication" @click="makeNewChat">
                <div class="flex items-center" style="height: 100%;">
                    <i class="i-tdesign:chat-add w-16 h-16 mr-10 ml-8 text-[var(--bt-tit-color-secondary)]"></i>
                    <div class="comu-title">{{ $t("新建对话") }}</div>
                </div>
            </li>

            <li :class="{ active: currentContextId == item.context_id }" @click.stop="handleChoose($event, item)"
                v-for="item in chatList" :key="item.context_id">
                <div class="flex items-center" style="height: 100%;" v-if="!item.agent_info">
                    <i class="i-tdesign:chat w-16 h-16 mr-10 ml-8 text-[var(--bt-tit-color-secondary)]"></i>
                    <div class="comu-title">{{ item.title ? item.title : $t("对话内容") }}</div>
                </div>
                <div class="flex items-center" style="height: 100%;" v-else>
                    <span v-if="item.agent_info.icon" class="mr-10 ml-8">{{ item.agent_info.icon
                        }}</span>
                    <span v-else class="mr-10 ml-8">😀</span>
                    <div class="comu-title">{{ item.agent_info.agent_title }}</div>
                </div>

                <n-popselect trigger="click" :options='[
                    { label: $t("删除对话"), value: "delChat" },
                    { label: $t("修改标题"), value: "modifyTitle" }
                ]' :on-update:value="(val: any) => doChatOperateSelect(val, item.context_id)">
                    <div class="flex justify-center items-center" style="height: 100%; padding: 0 8px;">
                        <i class="i-common:more-operation w-16 h-16"></i>
                    </div>
                </n-popselect>
            </li>
        </ul>
    </div>
</template>

<script setup lang="ts">
import { handleChoose, doChatOperateSelect } from "../controller"
import { get_chat_list, makeNewChat, doFold,cleanAllChats } from "@/views/Sider/controller"
import { getSiderStoreData } from '../store';
import { useI18n } from "vue-i18n";
const { t:$t } = useI18n();

const {
    chatList,
    currentContextId,
} = getSiderStoreData()
</script>

<style scoped lang="scss">
@use "@/assets/base";

.recent-comunication {
    width: 100%;

    .recent-list {
        @include base.recent-list-style;
    }
}
</style>