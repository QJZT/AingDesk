import { ElectronEgg } from 'ee-core';
import { Lifecycle } from './preload/lifecycle';
import { preload } from './preload';
import { totalService } from './service/total';
import mcp from './controller/mcp';
const { spawn } = require('child_process');
import { pub } from './class/public';
import * as path from 'path';
// New app
const app = new ElectronEgg();

// Register lifecycle
const life = new Lifecycle();
app.register("ready", life.ready);
app.register("electron-app-ready", life.electronAppReady);
app.register("window-ready", life.windowReady);
app.register("before-close", life.beforeClose);
// Register preload
app.register("preload", preload);


// Register service
setTimeout(() => {
    // 分享服务
    const { shareService } = require('./service/share');
    const {mcpService} = require('./service/mcp');
    const shareIdPrefix = shareService.generateUniquePrefix();
    let socket = shareService.connectToCloudServer(shareIdPrefix);
    shareService.startReconnect(socket,shareIdPrefix);

    // RAG后台任务
    const { RagTask } = require('./rag/rag_task');
    let ragTaskObj = new RagTask()
    ragTaskObj.parseTask()

    // 创建索引
    ragTaskObj.switchToCosineIndex()

    // 同步MCP服务器列表
    mcpService.sync_cloud_mcp()
    
}, 1000);

// 启动统计服务
totalService.start();


let goProcess;
function startGoService() {
    const goExePath = path.resolve(pub.get_resource_path(), 'exe/control-go.exe');
    const dbPath = path.join(pub.get_data_path(), 'sqlite-data.db');
    goProcess = spawn(goExePath, ['-p', dbPath], {
      cwd: path.dirname(goExePath)
    });
  
    goProcess.stdout.on('data', (data) => {
      console.log(`control-go服务输出: ${data}`);
    });
  
    goProcess.stderr.on('data', (data) => {
      console.error(`control-go服务错误: ${data}`);
    });
  
    goProcess.on('close', (code) => {
      console.log(`control-go服务退出，代码 ${code}`);
    });
}
startGoService();


let xttsProcess;
function startXttsService() {
    const xttsExePath = path.resolve(pub.get_resource_path(), 'exe/xtts.exe');
    const xttsPath = path.join(pub.get_data_path(), 'xtts');
    xttsProcess = spawn(xttsExePath, ['-p', xttsPath], {
      cwd: path.dirname(xttsExePath)
    });
  
    xttsProcess.stdout.on('data', (data) => {
      console.log(`xtts服务输出: ${data}`);
    });
  
    xttsProcess.stderr.on('data', (data) => {
      console.error(`xtts服务错误: ${data}`);
    });
  
    xttsProcess.on('close', (code) => {
      console.log(`xtts服务退出，代码 ${code}`);
    });
}
startXttsService();


app.register("before-close", () => {
    life.beforeClose();
    if (goProcess && !goProcess.killed) {
      goProcess.kill();
    }
    if (xttsProcess && !xttsProcess.killed) {
      xttsProcess.kill();
    }
});  
// Run
app.run();