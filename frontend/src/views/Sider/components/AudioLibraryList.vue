<template>
    <div class="recent-comunication">
        <!-- 列表 -->
        <ul class="recent-list">
            <li @click.stop="createNewKnowledgeStore" :class="{ 'add-knowledge': addingKnowledge }">
                <div class="flex items-center" style="height: 100%;">
                    <i class="i-tdesign:folder-add w-16 h-16 mr-10 ml-8 text-[var(--bt-tit-color-secondary)]"></i>
                    <div class="comu-title">{{ $t("新建人声") }}</div>
                </div>
            </li>
            <!-- openKnowledgeStore(item) -->
            <li  @click="openFileModal(item)"
                v-for="item in list">
                <div class="flex items-center" style="height: 100%;">
                    <i class="i-tdesign:folder w-16 h-16 mr-10 ml-8 text-[var(--bt-tit-color-secondary)]"></i>
                    <div class="comu-title"> {{ item  }}</div>
                </div>
                <n-popselect trigger="click" :options='options'
                    :on-update:value="(val: any) => dealPopOperation(val, item)">
                    <div class="flex justify-center items-center" style="height: 100%; padding: 0 8px;" @click.stop>
                        <i class="i-common:more-operation w-16 h-16"></i>
                    </div>
                </n-popselect>
            </li>
        </ul>
        <!-- 新增的对话框 -->
        <n-modal v-model:show="showModal">
            <n-card
                style="width: 600px"
                title="新建人声"
                :bordered="false"
                size="huge"
                role="dialog"
                aria-modal="true"
            >
            <template #header-extra>
                <i class="i-tdesign:close-circle w-24 h-24 cursor-pointer text-[#909399]"
                    @click="showModal = false"></i>
            </template>
                请输入口播文案：
                <n-input v-model:value="koubowenan" placeholder="请输入口播文案" type="textarea" :autosize="{ minRows: 4, maxRows: 6 }"/>
                <div style="height: 20px;"></div>
                上传音频文件：
                <div class="upload-section">
                        <n-upload
                        :multiple="false"
                        directory-dnd
                        :show-file-list="true"
                        @change="handleFileSelect"
                        >
                            <n-upload-dragger>
                            <div style="margin-bottom: 12px">
                                <n-icon size="48" :depth="3">
                                <ArchiveIcon />
                                </n-icon>
                            </div>
                            <n-text style="font-size: 16px">
                                点击或者拖动文件到该区域来上传
                            </n-text>
                            <n-p depth="3" style="margin: 8px 0 0 0">
                                请不要上传敏感数据，比如你的银行卡号和密码，信用卡号有效期和安全码
                            </n-p>
                            </n-upload-dragger>
                        </n-upload>
                    </div>
                <template #footer>
                    <n-space justify="end">
                        <n-button @click="showModal = false">取消</n-button>
                        <n-button type="primary" @click="handleConfirm">确认</n-button>
                    </n-space>
                </template>
            </n-card>
        </n-modal>

        <!-- // 在模板中添加文件弹窗 -->
        <n-modal v-model:show="showFileModal">
            <n-card
                style="width: 800px"
                title="文件管理"
                :bordered="false"
                size="huge"
                role="dialog"
                aria-modal="true"
            >
            <template #header-extra>
                <i class="i-tdesign:close-circle w-24 h-24 cursor-pointer text-[#909399]"
                    @click="showFileModal = false"></i>
            </template>
                <div class="file-manager-container">
                    <!-- 上传区域 -->
                    <div class="upload-section">
                        <n-upload
                        action="http://127.0.0.1:7072/upload"
                        :data="{
                            id: currentItem.id
                        }"
                        multiple
                        directory-dnd
                        :on-change="handleUploadChange"
                        :on-error="handleUploadError"
                        :on-finish="handleUploadFinish"
                        :on-before-upload="handleBeforeUpload"
                        :on-download="handleDownload"
                        :on-remove="handleRemove"
                        :on-retry="handleRetry"
                        >
                            <n-upload-dragger>
                            <div style="margin-bottom: 12px">
                                <n-icon size="48" :depth="3">
                                <ArchiveIcon />
                                </n-icon>
                            </div>
                            <n-text style="font-size: 16px">
                                点击或者拖动文件到该区域来上传
                            </n-text>
                            <n-p depth="3" style="margin: 8px 0 0 0">
                                请不要上传敏感数据，比如你的银行卡号和密码，信用卡号有效期和安全码
                            </n-p>
                            </n-upload-dragger>
                        </n-upload>
                    </div>

                    <!-- 文件列表 -->
                    <div class="file-list-section">
                        <n-list bordered>
                            <n-list-item v-for="file in fileList" :key="file.id" class="file-item">
                                <div class="file-info" style="display: flex; align-items: center; gap: 8px;justify-content: space-between; ">
                                    <i class="i-tdesign:file"></i>
                                    <span class="file-name">{{ file.file_name }}</span>
                                      <n-space>
                                    <n-button text type="primary" @click="handleDownload(file)">
                                        <template #icon>
                                            <i class="i-tdesign:download"></i>
                                        </template>
                                        下载
                                    </n-button>
                                    <n-button text type="error" @click="handleDelete(file)">
                                        <template #icon>
                                            <i class="i-tdesign:delete"></i>
                                        </template>
                                        删除
                                    </n-button>
                                </n-space>
                                </div>
                             
                            </n-list-item>
                        </n-list>
                    </div>
                </div>
            </n-card>
        </n-modal>
    </div>
</template>

<script setup lang="tsx">
import { openKnowledgeStore } from "@/views/KnowleadgeStore/controller"
// import { dealPopOperation } from "@/views/Sider/controller"
import { getKnowledgeStoreData } from '@/views/KnowleadgeStore/store';
import { useI18n } from "vue-i18n";
import { message, useDialog } from "@/utils/naive-tools"
const { t: $t } = useI18n()
const { knowledgeList, addingKnowledge, activeKnowledge, } = getKnowledgeStoreData()

const list = ref([]); // []ItemType[];

onMounted(async () => {
    await api_names();
    const intervalId = setInterval(api_names, 20000);
    onBeforeUnmount(() => clearInterval(intervalId));
});

const api_names = async () => {
    const response = await fetch('http://127.0.0.1:7073/get_human_voice_files', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })
    const data = await response.json()
    console.log(data.files);
    list.value = data.files;
    console.log("list.value:",list.value);
    // get_human_voice_files_text_map.value = data.text_map || {}
    // 转换为下拉选项格式
    // humanVoiceOptions.value = (data.files || []).map(name => ({
    //   label: name,
    //   value: name
    // }))
    // } catch (error) {
    //     console.error('获取人声音色文件失败:', error)
    // }
    // try {
    //     const response = await fetch('http://127.0.0.1:7072/names', {
    //         method: 'POST',
    //         headers: {
    //             'Content-Type': 'application/json',
    //         },
    //     });
    //     if (!response.ok) {
    //         throw new Error('Failed to fetch data');
    //     }
    //     const data = await response.json();
    //     interface ItemType {
    //         name: string;
    //         id: any;
    //     }
    //     list.value = data.map((item: ItemType) => ({
    //         name: item.name,
    //         id: item.id,
    //     }));
    // } catch (error) {
    //     console.error('Error fetching knowledge list:', error);
    // }
};

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

const showModal = ref(false);
const koubowenan = ref('');
const selectedFiles = ref([]); // 新增：用于暂存选择的文件

const handleFileSelect = ({ fileList: files }) => {
    if (files.length > 1) {
        message.error('只能提交一个文件') 
        return;
    }
    selectedFiles.value = files.map(file => file.file);
    console.log(selectedFiles.value);
};
const handleConfirm = async () => {
    console.log(selectedFiles.value[0]);
    try {
        const formData = new FormData();
        formData.append('text', koubowenan.value);
        if (selectedFiles.value.length > 0) {
            formData.append('file', selectedFiles.value[0]);
            formData.append('filename', selectedFiles.value[0].name); // 确保文件名也被发送
        }
        const response = await fetch('http://127.0.0.1:7073/convert_and_save_audio_text', {
            method: 'POST',
            body: formData,
        });

        if (!response.ok) {
            throw new Error('Failed to set audio name');
        }
        message.success("上传成功！")
        showModal.value = false;
        koubowenan.value = '';
        api_names() // 刷新列表?
        // 可以在这里添加刷新列表的逻辑
    } catch (error) {
        console.error('Error setting audio name:', error);
    }
};

const createNewKnowledgeStore = () => {
    showModal.value = true;
};

/***
 * @description 知识库操作
 */
 const dealPopOperation = async (val: string, item: any) => {
    if (val == "delChat") {
        try {
        const response = await fetch('http://127.0.0.1:7072/del_name', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ id: item.id }),
        });
        if (!response.ok) {
            throw new Error('Failed to set audio name');
        }
        api_names() // 刷新列表
    } catch (error) {
        console.error('Error setting audio name:', error);
    }
    }
}

const showFileModal = ref(false);
const currentItem = ref(null);
const fileList = ref([]);

const openFileModal = async (item) => {
    showFileModal.value = true;
    currentItem.value = item;
    await getFileList(item.id);
};

const getFileList = async (itemId) => {
    try {
        const response = await fetch(`http://127.0.0.1:7072/get_name_list`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ id: itemId }),
        });
        if (!response.ok) {
            throw new Error('Failed to fetch file list');
        }
        const data = await response.json();
        fileList.value = data.list; // 从返回的 JSON 中提取 list 字段   
    } catch (error) {
        console.error('Error fetching file list:', error);
    }
};


const handleUploadChange = (options) => {
    console.log('文件状态变化:', options);
};

const handleUploadError = (options) => {
    console.error('上传失败:', options);
};

const handleUploadFinish = (options) => {
    console.log('上传完成:', options);
    getFileList(currentItem.value.id); // 刷新文件列表
};

const handleBeforeUpload = (options) => {
    console.log('上传前检查:', options);
    return true;
};

const handleDownload = async (file) => {
    try {
        const response = await fetch(`http://127.0.0.1:7072/download/${file.associated_id}/${file.file_name}`);
        if (!response.ok) {
            throw new Error('Failed to download file');
        }
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = file.file_name;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
    } catch (error) {
        console.error('Error downloading file:', error);
    }
};
const handleRemove = (options) => {
    console.log('删除文件:', options);
    return true;
};

const handleRetry = (options) => {
    console.log('重试上传:', options);
    return true;
};

const handleDelete = async (file: any) => {
    try {
        const response = await fetch(`http://127.0.0.1:7072/delete_file/${file.associated_id}/${file.file_name}`);
        if (!response.ok) {
            throw new Error('Failed to delete file');
        }
        // 刷新文件列表
        getFileList(file.associated_id);
    } catch (error) {
        console.error('Error deleting file:', error);
    }
};

</script>

<style scoped lang="scss">
@use "@/assets/base";

.recent-comunication {
    width: 100%;

    .recent-list {
        @include base.recent-list-style;
    }
}

.file-manager-container {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.upload-section {
    display: flex;
    justify-content: center;
}

.file-list-section {
    max-height: 400px;
    overflow-y: auto;
}

.file-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
}

.file-info {
    display: flex;
    align-items: center;
    gap: 8px;
}

.file-name {
    font-size: 14px;
    color: var(--text-color);
}

.upload-button {
    width: 100%;
}
</style>