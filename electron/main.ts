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

// // Register lifecycle
const life = new Lifecycle();
// app.register("ready", life.ready);
// app.register("electron-app-ready", life.electronAppReady);
// app.register("window-ready", life.windowReady);
// app.register("before-close", life.beforeClose);
// // Register preload
// app.register("preload", preload);

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
    return new Promise<void>((resolve) => {
        const goExePath = path.resolve(pub.get_resource_path(), 'exe/control-go.exe');
        const dbPath = path.join(pub.get_data_path(), 'sqlite-data.db');
        const goDataPath = path.join(pub.get_data_path(), 'control-go-file');
        const p2pPath = path.join(pub.get_data_path(), 'control-go-file-audio-modules');
        goProcess = spawn(goExePath, ['-p', dbPath,'-p1', goDataPath ,"-p2" , p2pPath], {
            cwd: path.dirname(goExePath)
        });
        console.log('启动go指令:', '-p', dbPath,'-p1', goDataPath ,"-p2" , p2pPath);
        goProcess.stdout.on('data', (data) => {
            console.log(`[py7072 out]: ${data}`);
        });
        const pingInterval = setInterval(() => {
            const http = require('http');
            const options = {
                hostname: 'localhost',
                port: 7072,
                path: '/ping',
                method: 'GET',
                timeout: 5000
            };
            const req = http.request(options, (res) => {
                if (res.statusCode === 200) {
                    clearInterval(pingInterval);
                    console.log('7072 run ok');
                    resolve();
                }
            });
            req.on('error', () => {});
            req.end();
        }, 1000);
    });
}


// 先不注册生命周期事件，等服务启动完成后再注册
const fs = require('fs');

// 修改startPy7073Service函数
let py7073Process;
function startPy7073Service() {
    return new Promise<void>((resolve) => {
        const pythonExePath = path.resolve(pub.get_resource_path(), 'exe/py7073-env/python.exe');
        const pythonScriptPath = path.resolve(pub.get_resource_path(), 'exe/py7073-code/run.py');
        const dataPath = path.join(pub.get_resource_path(), 'exe/py7074-code/data');
        if (!fs.existsSync(dataPath)) { //创建目录
            fs.mkdirSync(dataPath, { recursive: true });
            console.log('已创建目录:', dataPath);
        }
        // D:\\androidWork\\AingDesk\\build\\extraResources\\exe\\py7073-data\\wav\\task_voice\\59ce4dff-0a43-411d-8334-283e097a12ff_1747418025042.wav"
        py7073Process = spawn(pythonExePath, [pythonScriptPath, '-p', dataPath], {
            cwd: path.dirname(pythonExePath)
        });
        // py7073Process.stdout.on('data', (data) => {
        //     console.log(`[py7073 stdout]: ${data}`);
        // });
        py7073Process.stderr.on('data', (data) => {
            console.error(`[py7073 stderr]: ${data}`);
        });
        console.log('启动py7073指令:', pythonExePath, pythonScriptPath, '-p', dataPath);
        const pingInterval = setInterval(() => {
            const http = require('http');
            const options = {
                hostname: 'localhost',
                port: 7073,
                path: '/ping',
                method: 'GET',
                timeout: 5000
            };
            const req = http.request(options, (res) => {
                if (res.statusCode === 200) {
                    clearInterval(pingInterval);
                    console.log('py7073 run ok');
                    resolve();
                }
            });
            req.on('error', () => {});
            req.end();
        }, 1000);
    });
}


// 修改startPy7073Service函数
let py7074Process;
function startPy7074Service() {
    return new Promise<void>((resolve) => {
      // .\runtime2\python.exe api.py -a 0.0.0.0 -p 7074 -s D:\pythonWork\GPT-SoVITS-v3\data\models\ -g D:\pythonWork\GPT-SoVITS-v3\data\models\ 
        const pythonExePath = path.resolve(pub.get_resource_path(), 'exe/py7074-code/runtime/python.exe');
        const pythonScriptPath = path.resolve(pub.get_resource_path(), 'exe/py7074-code/api.py');
        const dataPathS = path.join(pub.get_resource_path(), 'exe/py7074-code/data/models/');
        const dataPathG = path.join(pub.get_resource_path(), 'exe/py7074-code/data/models/');
//         [py7074 stderr]: Traceback (most recent call last):
//   File "D:\androidWork\AingDesk\build\extraResources\exe\py7074-code\api.py", line 163, in <module>
//     from feature_extractor import cnhubert
// ModuleNotFoundError: No module named 'feature_extractor'
        // 打印启动指令
        console.log('启动py7074指令:', pythonExePath, pythonScriptPath, '-a', "0.0.0.0", '-p', "7074", '-s', dataPathS, "-g", dataPathG);
        py7074Process = spawn(pythonExePath, [pythonScriptPath, '-a', "0.0.0.0" ,'-p', "7074",'-s',dataPathS , "-g",dataPathG], {
            cwd: path.resolve(pub.get_resource_path(), 'exe/py7074-code')
        });
        
        // py7074Process.stdout.on('data', (data) => {
        //     console.log(`[py7074 stdout]: ${data}`);
        // });
        py7074Process.stderr.on('data', (data) => {
            console.error(`[py7074 stderr]: ${data}`);
        });
        const pingInterval = setInterval(() => {
            const http = require('http');
            const options = {
                hostname: 'localhost',
                port: 7074,
                path: '/ping',
                method: 'GET',
                timeout: 5000
            };
            const req = http.request(options, (res) => {
                if (res.statusCode === 200) {
                    clearInterval(pingInterval);
                    console.log('py7074 run ok');
                    resolve();
                }
            });
            req.on('error', () => {});
            req.end();
        }, 1000);
    });
}

// let xttsProcess;
// function startXttsService() {
//     const xttsExePath = path.resolve(pub.get_resource_path(), 'exe/xtts.exe');
//     const xttsPath = path.join(pub.get_data_path(), 'xtts');
//     xttsProcess = spawn(xttsExePath, ['-p', xttsPath], {
//       cwd: path.dirname(xttsExePath)
//     });
  
//     xttsProcess.stdout.on('data', (data) => {
//       console.log(`xtts服务输出: ${data}`);
//     });
  
//     xttsProcess.stderr.on('data', (data) => {
//       console.error(`xtts服务错误: ${data}`);
//     });
  
//     xttsProcess.on('close', (code) => {
//       console.log(`xtts服务退出，代码 ${code}`);
//     });
// }
// startXttsService();

// Run

// 修改应用启动流程
async function initializeApp() {
  try {
    //   const { BrowserWindow } = require('electron');
    //   // 创建加载提示窗口
    //   const loadingWindow = new BrowserWindow({
    //       width: 300,
    //       height: 200,
    //       frame: false,
    //       transparent: true,
    //       alwaysOnTop: true,
    //       webPreferences: {
    //           nodeIntegration: true
    //       }
    //   });
    //   loadingWindow.loadFile('loading.html');
      
      app.register("before-close", () => {
        life.beforeClose();
        if (goProcess && !goProcess.killed) {
          goProcess.kill();
        }
        if (py7073Process && !py7073Process.killed) {
          py7073Process.kill();
        }
        if (py7074Process && !py7074Process.killed) {
          py7074Process.kill();
        }
      });  
      console.log('py7072 run ...');
      await startGoService();//启动go
      //console.log('py7073 run ...');
      //await startPy7073Service(); //启动 7073
      //console.log('py7074 run ...');
      //await startPy7074Service();//启动 7074
      
      // 关闭加载窗口
    //   loadingWindow.close();
      
      // 服务启动完成后再注册生命周期和GUI相关逻辑
      app.register("ready", life.ready);
      app.register("electron-app-ready", life.electronAppReady);
      app.register("window-ready", life.windowReady);
      app.register("before-close", life.beforeClose);
      app.register("preload", preload);
      
      // 启动其他服务和GUI
      app.run();
  } catch (error) {
      console.error('应用初始化失败:', error);
      require('electron').app.quit();
  }
}

initializeApp();