<template>
    <div class="recent-comunication">
        <!-- 列表 -->
        <ul class="recent-list">
            <li @click.stop="createNewKnowledgeStore" :class="{ 'add-knowledge': addingKnowledge }">
                <div class="flex items-center" style="height: 100%;">
                    <i class="i-tdesign:folder-add w-16 h-16 mr-10 ml-8 text-[var(--bt-tit-color-secondary)]"></i>
                    <div class="comu-title">{{ $t("新建音频") }}</div>
                </div>
            </li>
            <li :class="[{ active: item.id == activeKnowledge }]" @click="openKnowledgeStore(item)"
                v-for="item in list">
                <div class="flex items-center" style="height: 100%;">
                    <i class="i-tdesign:folder w-16 h-16 mr-10 ml-8 text-[var(--bt-tit-color-secondary)]"></i>
                    <div class="comu-title">#{{ item.id }}  {{ item.name  }}</div>
                </div>

                <n-popselect trigger="click" :options='options'
                    :on-update:value="(val: any) => dealPopOperation(val, item)">
                    <div class="flex justify-center items-center" style="height: 100%; padding: 0 8px;" @click.stop>
                        <i class="i-common:more-operation w-16 h-16"></i>
                    </div>
                </n-popselect>
            </li>
        </ul>

    </div>
</template>

<script setup lang="tsx">
import { openKnowledgeStore, createNewKnowledgeStore } from "@/views/KnowleadgeStore/controller"
import { dealPopOperation } from "@/views/Sider/controller"
import { getKnowledgeStoreData } from '@/views/KnowleadgeStore/store';
import { useI18n } from "vue-i18n";
const { t: $t } = useI18n()
const { knowledgeList, addingKnowledge, activeKnowledge, } = getKnowledgeStoreData()

const list = ref([]); // []ItemType[];

onMounted(async () => {
    try {
        const response = await fetch('http://127.0.0.1:7072/names', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (!response.ok) {
            throw new Error('Failed to fetch data');
        }
        const data = await response.json();
        interface ItemType {
            name: string;
            id: any;
        }
        list.value = data.map((item: ItemType) => ({
            name: item.name,
            id: item.id,
        }));
    } catch (error) {
        console.error('Error fetching knowledge list:', error);
    }
});


const options = [
    { label: $t("删除"), value: "delChat" },
    // { label: $t("修改"), value: "modifyTitle" },
    // {
    //     value: "optimization",
    //     label: $t("优化"),
    //     render(info: any) {
    //         console.log(info)
    //         return <n-tooltip trigger="hover">
    //             {
    //                 {
    //                     trigger: () => info.node,
    //                     default:()=><div class="w-200">{$t("增强索引,并释放多余的空间占用")}</div>
    //                 }
    //             }
    //         </n-tooltip>
    //     }
    // }
]

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