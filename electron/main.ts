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

//     // 监听消息
//     let retryCount = 0;
//     const maxRetries = 20;

//     const pingInterval = setInterval(() => {
//       const http = require('http');
// const { app: electronApp } = require('electron');
//       const options = {
//           hostname: 'localhost',
//           port: 7072, // 假设服务运行在8080端口
//           path: '/ping',
//           method: 'GET',
//           headers: {
//               'Content-Type': 'application/json'
//           }
//       };

//       const req = http.request(options, (res) => {
//           console.log(`Ping接口响应状态码: ${res.statusCode}`);
//           res.on('data', (data) => {
//             clearInterval(pingInterval);
//           });
//           retryCount = 0; // 成功时重置计数器
//       });

//       req.on('error', (error) => {
//           console.error(`Ping err: ${error.message}`);
//           retryCount++;
//           if (retryCount >= maxRetries) {
//               clearInterval(pingInterval);
//               console.error('Ping err electronApp.quit');
//               // 由于 ElectronEgg 上不存在 quit 属性，推测使用 electron 的 app 模块来退出应用
//               electronApp.quit();
//           }
//       });

//       req.end();
//   }, 2000); // 每2秒调用一次
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