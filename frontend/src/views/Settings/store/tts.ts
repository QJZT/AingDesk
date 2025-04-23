import { defineStore, storeToRefs } from "pinia";
import { ref } from "vue";
import i18n from "@/lang";
import type { InstallProgress } from "@/views/Home/dto";

const useSettingsStoreTts = defineStore("settingsTts", () => {

    const settingsTtsShow = ref(false)
    // 机器配置信息
    const pcInfoTts = ref<Record<string, any>>({})
    // 可用模型列表
    const visibleModelListTts = ref<any[]>([])
    // 模型筛选关键字
    const modeTypeTts = ref("all")
    // 要安裝的模型名称
    const modelNameForInstallTts = ref<{ model: string; parameters: string }>({
        model: "",
        parameters: ""
    })
    // 安装模型的弹窗
    const installShowTts = ref(false)
    // 模型安装进度
    const modelInstallProgressTts = ref<InstallProgress>({
        status: 0,
        digest: "",
        total: 0,
        completed: 0,
        progress: 0,
        speed: 0
    })
    // 等待删除的大模型
    const modelForDelTts = ref("")
    // 删除大模型进度
    const modelDelLoadingTts = ref(false)
    // 删除模型问询
    const modelDelConfirmTts = ref(false)
    // 模型管理器安装进度
    const modelManagerInstallProgressTts = ref<InstallProgress>({
        status: 0,
        digest: "",
        total: 0,
        completed: 0,
        progress: 0,
        speed: 0
    })
    // 模型管理器安装提示
    const modelManagerInstallNoticeTts = ref("")
    // 模型管理器安装的位置
    const modelManagerInstallPathTts = ref("")
    // 模型管理器安装进度弹窗
    const modelManagerInstallProgresShowTts = ref(false)
    // 模型管理器安装问询
    const managerInstallConfirmTts = ref(false)
    // 选择需要安装的模型管理器
    const managerForInstallTts = ref("ollama")
    // 是否安装了模型管理器
    const isInstalledManagerTts = ref(false)
    // 下载展示文案
    const downloadTextTts = ref(i18n.global.t("正在连接，请稍候..."))
    // ollama接入地址
    const ollamaUrlTts = ref("")
    // 重新设定模型列表(安装ollama模型管理器后手动刷新一次数据)
    const isResetModelListTts = ref({
        status: false, // 是否刷新完成
        type: 0, // 0:默认，1: 重置模型列表
    })
    // 搜索过滤后的模型列表
    const filterListTts = ref<any[]>([])
    // 搜索模型关键字
    const searchTts = ref("")
    // 模型列表分页
    const paginationTts = ref({
        page: 1,
        pageSize: 10,
        showSizePicker: true,
        pageSizes: [10, 50, 100],
        onChange: (page: number) => {
            paginationTts.value.page = page
        },
        onUpdatePageSize: (pageSize: number) => {
            paginationTts.value.pageSize = pageSize
            paginationTts.value.page = 1
        }
    })
    return {
        settingsTtsShow,
        pcInfoTts,
        visibleModelListTts,
        modeTypeTts,
        modelNameForInstallTts,
        installShowTts,
        modelInstallProgressTts,
        modelForDelTts,
        modelDelLoadingTts,
        modelDelConfirmTts,
        modelManagerInstallProgressTts,
        modelManagerInstallNoticeTts,
        modelManagerInstallPathTts,
        modelManagerInstallProgresShowTts,
        managerInstallConfirmTts,
        managerForInstallTts,
        isInstalledManagerTts,
        downloadTextTts,
        ollamaUrlTts,
        isResetModelListTts,
        filterListTts,
        searchTts,
        paginationTts
    }
})

export function getSettingsStoreDataTts() {
    return storeToRefs(useSettingsStoreTts())
}