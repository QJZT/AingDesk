// src/api/moduleApi.ts
import axios from "axios";

// 定义模块数据类型（基于你的测试数据）
export interface Module {
    ID: number;
    ModuleName: string;
    IntervalTimeStart: number;
    IntervalTimeEnd: number;
    TriggerConditions: string[];
    ReadStep: string;
    ScriptContent: string;
    IsModelRewrite: boolean;
    RewriteFrequency: number;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: string | null;
    ModuleConfigID: number;
    OrderNum: number;
}

// 创建独立的 axios 实例
const moduleInstance = axios.create({
    baseURL: "http://localhost:7072", // 你的后端地址
    timeout: 10000,
});

// GET 请求方法
async function get<T>(url: string, params: any = {}): Promise<T> {
    const response = await moduleInstance.get(url, { params });
    return response.data;
}

/**
 * @description 获取模块列表
 * @returns 模块列表
 */
export async function getBaseModules(): Promise<Module[]> {
    try {
        const data = await get<Module[]>("/base-modules", { user_id: 1 }); // 固定 user_id 为 1
        return data;
    } catch (err) {
        console.error("获取模块数据失败:", err);
        throw err; // 直接抛出错误，交给调用方处理
    }
}