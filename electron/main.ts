import { ElectronEgg } from 'ee-core';
import { Lifecycle } from './preload/lifecycle';
import { preload } from './preload';
import { totalService } from './service/total';
import mcp from './controller/mcp';
const { spawn } = require('child_process');
import { pub } from './class/public';
import * as path from 'path';
import { ipcMain } from 'electron';
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
    return new Promise<void>((resolve, reject) => {
        const pythonExePath = path.resolve(pub.get_resource_path(), 'exe/py7073-env/python.exe');
        const pythonScriptPath = path.resolve(pub.get_resource_path(), 'exe/py7073-code/run.py');
        const dataPath = path.join(pub.get_resource_path(), 'exe/py7074-code/data');
        
        console.log('=== py7073 启动调试信息 ===');
        console.log('Python路径:', pythonExePath);
        console.log('脚本路径:', pythonScriptPath);
        console.log('数据路径:', dataPath);
        console.log('Python文件存在:', fs.existsSync(pythonExePath));
        console.log('脚本文件存在:', fs.existsSync(pythonScriptPath));
        
        // 检查文件是否存在
        if (!fs.existsSync(pythonExePath)) {
            console.error('Python可执行文件不存在:', pythonExePath);
            reject(new Error(`Python可执行文件不存在: ${pythonExePath}`));
            return;
        }
        if (!fs.existsSync(pythonScriptPath)) {
            console.error('Python脚本文件不存在:', pythonScriptPath);
            reject(new Error(`Python脚本文件不存在: ${pythonScriptPath}`));
            return;
        }
        
        if (!fs.existsSync(dataPath)) {
            fs.mkdirSync(dataPath, { recursive: true });
            console.log('已创建目录:', dataPath);
        }
        
        console.log('开始启动py7073进程...');
        console.log('启动py7073指令:', pythonExePath, pythonScriptPath, '-p', dataPath);
        
        try {
            py7073Process = spawn(pythonExePath, [pythonScriptPath, '-p', dataPath], {
                cwd: path.dirname(pythonExePath)
            });
            
            console.log('py7073进程已spawn，PID:', py7073Process.pid);
            
            // 监听进程退出
            py7073Process.on('exit', (code, signal) => {
                console.log(`py7073进程退出，退出码: ${code}, 信号: ${signal}`);
                if (code !== 0) {
                    console.error(`py7073进程异常退出，退出码: ${code}, 信号: ${signal}`);
                }
            });
            
            // 监听进程错误
            py7073Process.on('error', (error) => {
                console.error('py7073进程启动失败:', error);
                reject(new Error(`py7073进程启动失败: ${error.message}`));
            });
            
            // 启用stdout监听以便调试
            py7073Process.stdout.on('data', (data) => {
                console.log(`[py7073 stdout]: ${data}`);
            });
            
            py7073Process.stderr.on('data', (data) => {
                console.error(`[py7073 stderr]: ${data}`);
            });
            
            // 简化启动检查：等待3秒后认为启动成功
            console.log('等待3秒后认为py7073启动成功...');
            
            // 添加进程状态检查
            let processExited = false;
            let resolved = false;
            let timeoutId;
            
            // 监听进程退出
            py7073Process.on('exit', (code, signal) => {
                processExited = true;
                console.log(`py7073进程退出，退出码: ${code}, 信号: ${signal}`);
                if (timeoutId) {
                    clearTimeout(timeoutId);
                }
                if (!resolved && code !== 0) {
                    console.error(`py7073进程异常退出，退出码: ${code}, 信号: ${signal}`);
                    resolved = true;
                    reject(new Error(`py7073进程异常退出，退出码: ${code}`));
                }
            });
            
            // 监听进程错误
            py7073Process.on('error', (error) => {
                console.error('py7073进程启动失败:', error);
                if (timeoutId) {
                    clearTimeout(timeoutId);
                }
                if (!resolved) {
                    resolved = true;
                    reject(new Error(`py7073进程启动失败: ${error.message}`));
                }
            });
            
            // 强制刷新控制台输出
            console.log('等待3秒后认为py7073启动成功...');
            process.stdout.write(''); // 强制刷新
            
            timeoutId = setTimeout(() => {
                console.log('=== setTimeout回调开始执行 ===');
                console.log('processExited状态:', processExited);
                console.log('resolved状态:', resolved);
                
                if (!resolved) {
                    if (!processExited) {
                        console.log('py7073 服务启动完成（等待模式）');
                        console.log('py7073进程状态 - PID:', py7073Process.pid, '已退出:', py7073Process.killed);
                        resolved = true;
                        resolve();
                    } else {
                        console.log('py7073进程已退出，启动失败');
                        resolved = true;
                        reject(new Error('py7073进程已退出'));
                    }
                }
                console.log('=== setTimeout回调执行完成 ===');
            }, 3000);
            
        } catch (error) {
            console.error('spawn py7073进程时发生异常:', error);
            reject(error);
        }
    });
}       


// 修改startPy7073Service函数
let py7074Process;
function startPy7074Service() {
    return new Promise<void>((resolve, reject) => {
        const pythonExePath = path.resolve(pub.get_resource_path(), 'exe/py7074-code/runtime/python.exe');
        const pythonScriptPath = path.resolve(pub.get_resource_path(), 'exe/py7074-code/api_qjzt.py');
        
        // 检查文件是否存在
        if (!fs.existsSync(pythonExePath)) {
            reject(new Error(`Python可执行文件不存在: ${pythonExePath}`));
            return;
        }
        if (!fs.existsSync(pythonScriptPath)) {
            reject(new Error(`Python脚本文件不存在: ${pythonScriptPath}`));
            return;
        }
        
        console.log('启动py7074指令:', pythonExePath, pythonScriptPath, '-a', "0.0.0.0", '-p', "7074");
        
        py7074Process = spawn(pythonExePath, [pythonScriptPath, '-a', "0.0.0.0", '-p', "7074"], {
            cwd: path.resolve(pub.get_resource_path(), 'exe/py7074-code'),
            env: {
                ...process.env,
                PYTHONPATH: path.resolve(pub.get_resource_path(), 'exe/py7074-code'),
                PYTHONIOENCODING: 'utf-8'
            }
        });
        
        // 监听进程退出和错误
        py7074Process.on('exit', (code, signal) => {
            if (code !== 0) {
                console.error(`py7074进程异常退出，退出码: ${code}, 信号: ${signal}`);
            }
        });
        
        py7074Process.on('error', (error) => {
            console.error('py7074进程启动失败:', error);
            reject(new Error(`py7074进程启动失败: ${error.message}`));
        });
        
        py7074Process.stderr.on('data', (data) => {
            console.error(`[py7074 stderr]: ${data}`);
        });
        
        // 设置超时机制
        const startTimeout = setTimeout(() => {
            clearInterval(pingInterval);
            reject(new Error('py7074服务启动超时（60秒）'));
        }, 60000);
        
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
                    clearTimeout(startTimeout);
                    console.log('py7074 run ok');
                    resolve();
                }
            });
            req.on('error', (error) => {
                console.log(`py7074 ping失败: ${error.message}`);
            });
            req.on('timeout', () => {
                req.destroy();
            });
            req.end();
        }, 2000);
    });
}


// 修改startPy7073Service函数
let py9872Process;
function startPy9872Service() {
    return new Promise<void>((resolve) => {
      // .\runtime2\python.exe api.py -a 0.0.0.0 -p 7074 -s D:\pythonWork\GPT-SoVITS-v3\data\models\ -g D:\pythonWork\GPT-SoVITS-v3\data\models\ 
        const pythonExePath = path.resolve(pub.get_resource_path(), 'exe/py7074-code/runtime/python.exe');
        const pythonScriptPath = path.resolve(pub.get_resource_path(), 'exe/py7074-code/GPT_SoVITS/inference_webui_fast.py');
        // const dataPathS = path.join(pub.get_resource_path(), 'exe/py7074-code/data/models/');
        // const dataPathG = path.join(pub.get_resource_path(), 'exe/py7074-code/data/models/');
//         [py7074 stderr]: Traceback (most recent call last):
//   File "D:\androidWork\无人直播\build\extraResources\exe\py7074-code\api.py", line 163, in <module>
//     from feature_extractor import cnhubert
// ModuleNotFoundError: No module named 'feature_extractor'
        // 打印启动指令
        console.log('启动py9872指令:', pythonExePath, pythonScriptPath,);
        py9872Process = spawn(pythonExePath, [pythonScriptPath,], {
            cwd: path.resolve(pub.get_resource_path(), 'exe/py7074-code'),
            env: {
                ...process.env,
                PYTHONPATH: path.resolve(pub.get_resource_path(), 'exe/py7074-code'),
                PYTHONIOENCODING: 'utf-8'
            }
        });
        
        // py7074Process.stdout.on('data', (data) => {
        //     console.log(`[py7074 stdout]: ${data}`);
        // });
        // py9872Process.stderr.on('data', (data) => {
        //     console.error(`[py7074 stderr]: ${data}`);
        // });
        const pingInterval = setInterval(() => {
            const http = require('http');
            const options = {
                hostname: 'localhost',
                port: 9872,
                path: '/',
                method: 'GET',
                timeout: 5000
            };
            const req = http.request(options, (res) => {
                if (res.statusCode === 200) {
                    clearInterval(pingInterval);
                    console.log('py9872 run ok');
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

// 添加服务控制的IPC处理程序
ipcMain.handle('start-py9872-service', async () => {
  try {
    if (py9872Process && !py9872Process.killed) {
      return { success: false, message: 'py9872服务已在运行中' };
    }
    await startPy9872Service();
    return { success: true, message: 'py9872服务启动成功' };
  } catch (error) {
    console.error('启动py9872服务失败:', error);
    return { success: false, message: `启动py9872服务失败: ${error.message}` };
  }
});

ipcMain.handle('stop-py9872-service', async () => {
  try {
    if (!py9872Process || py9872Process.killed) {
      return { success: false, message: 'py9872服务未运行' };
    }
    py9872Process.kill();
    py9872Process = null;
    return { success: true, message: 'py9872服务已停止' };
  } catch (error) {
    console.error('停止py9872服务失败:', error);
    return { success: false, message: `停止py9872服务失败: ${error.message}` };
  }
});

ipcMain.handle('start-py7074-service', async () => {
  try {
    if (py7074Process && !py7074Process.killed) {
      return { success: false, message: 'py7074服务已在运行中' };
    }
    await startPy7074Service();
    return { success: true, message: 'py7074服务启动成功' };
  } catch (error) {
    console.error('启动py7074服务失败:', error);
    return { success: false, message: `启动py7074服务失败: ${error.message}` };
  }
});

ipcMain.handle('stop-py7074-service', async () => {
  try {
    if (!py7074Process || py7074Process.killed) {
      return { success: false, message: 'py7074服务未运行' };
    }
    py7074Process.kill();
    py7074Process = null;
    return { success: true, message: 'py7074服务已停止' };
  } catch (error) {
    console.error('停止py7074服务失败:', error);
    return { success: false, message: `停止py7074服务失败: ${error.message}` };
  }
});


ipcMain.handle('get-py9872-service-status', async () => {
  const isRunning = py9872Process && !py9872Process.killed;
  return { 
    isRunning, 
    message: isRunning ? 'py9872服务正在运行' : 'py9872服务未运行' 
  };
});

// Run

// 修改应用启动流程
async function initializeApp() {
  try {
    // 检查是否是重启后的启动
    if (process.env.ACTIVATION_RESTART === 'true') {
      console.log('检测到激活重启，清除重启标志');
      delete process.env.ACTIVATION_RESTART;
    }
    
    console.log('开始启动Go服务...');
    
    // 注册进程退出时的清理逻辑
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
      if (py9872Process && !py9872Process.killed) {
        py9872Process.kill();
      }
    });
    
    // 首先启动Go服务(7072端口)，这是激活码验证的前提
    await startGoService();
    console.log('Go服务启动成功，现在可以进行激活码验证');
    
    // 等待一下确保服务完全启动
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    // 先检查是否已有有效的激活码
    const hasValidCode = await checkExistingActivationCode();
    
    let isActivated = hasValidCode;
    
    // 如果没有有效激活码，显示验证窗口
    if (!hasValidCode) {
        console.log('未发现有效激活码，显示验证窗口');
        isActivated = await showActivationWindow();
    }
    
    if (!isActivated) {
      console.log('激活码验证失败或用户取消，退出应用');
      require('electron').app.quit();
      return;
    }

    console.log('激活码验证成功，准备重启应用程序...');
    
    // 如果是首次激活，重启应用程序
    if (!hasValidCode) {
      console.log('首次激活成功，正在重启应用程序...');
      const { app } = require('electron');
      
      // 设置重启标志，避免无限重启
      process.env.ACTIVATION_RESTART = 'true';
      
      // 重启应用程序
      app.relaunch();
      app.quit();
      return;
    }
    
    console.log('激活码验证成功，开始启动其他后端服务...');
    
    // 启动其他Python服务
    console.log('正在启动 py7073 服务...');
    await startPy7073Service();
    console.log('py7073 服务启动成功!');
    
    console.log('正在启动 py7074 服务...');
    await startPy7074Service();
    console.log('py7074 服务启动成功!');
    
    console.log('正在启动 py9872 服务...');
    await startPy9872Service();
    console.log('py9872 服务启动成功!');
    
    console.log('所有服务启动完成，注册应用生命周期...');
    
    // 注册应用生命周期
    app.register("ready", life.ready);
    app.register("electron-app-ready", life.electronAppReady);
    app.register("window-ready", life.windowReady);
    app.register("preload", preload);
    
    // 启动主应用
    app.run();
    
    console.log('应用启动完成!');
    
  } catch (error) {
    console.error('应用初始化失败:', error);
    const { dialog } = require('electron');
    await dialog.showErrorBox('启动失败', `应用启动失败：${error.message}`);
    require('electron').app.quit();
  }
}

// 新的激活码验证窗口函数
async function showActivationWindow(): Promise<boolean> {
  return new Promise((resolve) => {
    const { BrowserWindow, ipcMain } = require('electron');
    
    // 创建激活码验证窗口作为主窗口
    const activationWindow = new BrowserWindow({
      width: 600,
      height: 500,
      center: true,
      resizable: false,
      minimizable: false,
      maximizable: false,
      closable: true,
      alwaysOnTop: true,
      title: '激活码验证',
      webPreferences: {
        nodeIntegration: true,
        contextIsolation: false
      }
    });
    
    // 阻止窗口关闭，只能通过验证成功或明确取消
    activationWindow.on('close', (event) => {
      event.preventDefault();
      // 显示确认对话框
      const { dialog } = require('electron');
      const result = dialog.showMessageBoxSync(activationWindow, {
        type: 'question',
        buttons: ['确认退出', '继续验证'],
        defaultId: 1,
        title: '确认退出',
        message: '您确定要退出应用吗？',
        detail: '退出后需要重新启动应用进行激活码验证。'
      });
      
      if (result === 0) {
        activationWindow.destroy();
        resolve(false);
      }
    });
    
    // 监听激活码验证结果
    ipcMain.once('activation-result', (event, success) => {
      activationWindow.destroy();
      resolve(success);
    });
    
    // 加载激活码验证页面
    const htmlContent = `
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <title>激活码验证</title>
        <style>
            body {
                font-family: 'Microsoft YaHei', sans-serif;
                margin: 0;
                padding: 20px;
                background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
                color: white;
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: center;
                min-height: 100vh;
                box-sizing: border-box;
            }
            .container {
                background: rgba(255, 255, 255, 0.1);
                padding: 40px;
                border-radius: 15px;
                backdrop-filter: blur(10px);
                box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
                text-align: center;
                max-width: 400px;
                width: 100%;
            }
            h1 {
                margin-bottom: 30px;
                font-size: 24px;
                font-weight: 300;
            }
            .input-group {
                margin-bottom: 20px;
                text-align: left;
            }
            label {
                display: block;
                margin-bottom: 8px;
                font-size: 14px;
                opacity: 0.9;
            }
            input {
                width: 100%;
                padding: 12px;
                border: none;
                border-radius: 8px;
                background: rgba(255, 255, 255, 0.2);
                color: white;
                font-size: 16px;
                box-sizing: border-box;
            }
            input::placeholder {
                color: rgba(255, 255, 255, 0.7);
            }
            input:focus {
                outline: none;
                background: rgba(255, 255, 255, 0.3);
            }
            .button-group {
                display: flex;
                gap: 15px;
                margin-top: 30px;
            }
            button {
                flex: 1;
                padding: 12px 20px;
                border: none;
                border-radius: 8px;
                font-size: 16px;
                cursor: pointer;
                transition: all 0.3s ease;
            }
            .btn-primary {
                background: #4CAF50;
                color: white;
            }
            .btn-primary:hover {
                background: #45a049;
            }
            .btn-secondary {
                background: rgba(255, 255, 255, 0.2);
                color: white;
            }
            .btn-secondary:hover {
                background: rgba(255, 255, 255, 0.3);
            }
            .message {
                margin-bottom: 20px;
                padding: 10px;
                border-radius: 5px;
                font-size: 14px;
            }
            .error {
                background: rgba(244, 67, 54, 0.3);
                border: 1px solid rgba(244, 67, 54, 0.5);
            }
            .loading {
                display: none;
                margin-top: 20px;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>激活码验证</h1>
            <div id="message" class="message" style="display: none;"></div>
            
            <div class="input-group">
                <label for="activationCode">请输入激活码：</label>
                <input type="text" id="activationCode" placeholder="请输入您的激活码" />
            </div>
            
            <div class="button-group">
                <button type="button" class="btn-secondary" onclick="cancelActivation()">退出</button>
                <button type="button" class="btn-primary" onclick="validateCode()">验证</button>
            </div>
            
            <div class="loading" id="loading">
                <p>正在验证激活码，请稍候...</p>
            </div>
        </div>
        
        <script>
            const { ipcRenderer } = require('electron');
            
            let isValidating = false;
            
            // 回车键验证
            document.getElementById('activationCode').addEventListener('keypress', function(e) {
                if (e.key === 'Enter' && !isValidating) {
                    validateCode();
                }
            });
            
            // 验证激活码
            async function validateCode() {
                if (isValidating) return;
                
                const code = document.getElementById('activationCode').value.trim();
                if (!code) {
                    showMessage('请输入激活码', 'error');
                    return;
                }
                
                isValidating = true;
                document.getElementById('loading').style.display = 'block';
                
                try {
                    const result = await validateActivationCode(code);
                    
                    if (result.success) {
                        showMessage('激活码验证成功！请重启应用...', 'success');
                        setTimeout(() => {
                            ipcRenderer.send('activation-result', true);
                        }, 1000);
                    } else {
                        showMessage(result.message || '激活码验证失败', 'error');
                        isValidating = false;
                        document.getElementById('loading').style.display = 'none';
                    }
                } catch (error) {
                    showMessage('验证过程中发生错误：' + error.message, 'error');
                    isValidating = false;
                    document.getElementById('loading').style.display = 'none';
                }
            }
            
            // 取消激活
            function cancelActivation() {
                ipcRenderer.send('activation-result', false);
            }
            
            // 显示消息
            function showMessage(text, type) {
                const messageEl = document.getElementById('message');
                messageEl.textContent = text;
                messageEl.className = 'message ' + type;
                messageEl.style.display = 'block';
            }
            
            // 激活码验证函数
            async function validateActivationCode(code) {
                const response = await fetch('http://127.0.0.1:7072/client_config', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ activation_code: code }),
                    timeout: 10000
                });
                
                const data = await response.json();
                
                return {
                    success: data.activation_valid === true,
                    message: data.message,
                    clientConfig: data
                };
            }
        </script>
    </body>
    </html>
    `;
    
    activationWindow.loadURL('data:text/html;charset=utf-8,' + encodeURIComponent(htmlContent));
    activationWindow.show();
  });
}

// 激活码验证函数
function validateActivationCode(activationCode?: string): Promise<{success: boolean, message: string, needInput: boolean, clientConfig?: any}> {
    return new Promise((resolve) => {
        const http = require('http');
        const postData = activationCode ? JSON.stringify({ activation_code: activationCode }) : '';
        
        const options = {
            hostname: '127.0.0.1',
            port: 7072,
            path: '/client_config',
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Content-Length': Buffer.byteLength(postData)
            },
            timeout: 10000
        };
        
        const req = http.request(options, (res) => {
            let data = '';
            res.on('data', (chunk) => {
                data += chunk;
            });
            
            res.on('end', () => {
                try {
                    const response = JSON.parse(data);
                    console.log('激活码验证响应:', response);
                    
                    if (response.activation_valid === true) {
                        // 验证成功
                        resolve({ 
                            success: true, 
                            message: response.message || '激活码验证成功', 
                            needInput: false,
                            clientConfig: response.client_config
                        });
                    } else {
                        // 验证失败，判断是否需要用户输入新激活码
                        const needInput = response.message && (
                            response.message.includes('激活码不存在') || 
                            response.message.includes('激活码为空') ||
                            response.message.includes('请填写激活码') ||  // 添加这个条件
                            response.message.includes('过期') ||
                            response.message.includes('设备不匹配') ||
                            response.message.includes('无效')
                        );
                        
                        resolve({ 
                            success: false, 
                            message: response.message || '激活码验证失败', 
                            needInput,
                            clientConfig: response.client_config
                        });
                    }
                } catch (error) {
                    console.error('解析激活码验证响应失败:', error);
                    resolve({ 
                        success: false, 
                        message: '激活码验证接口响应解析失败', 
                        needInput: false 
                    });
                }
            });
        });
        
        req.on('error', (error) => {
            console.error('激活码验证请求失败:', error);
            resolve({ 
                success: false, 
                message: '激活码验证接口请求失败', 
                needInput: false 
            });
        });
        
        req.on('timeout', () => {
            console.error('激活码验证请求超时');
            req.destroy();
            resolve({ 
                success: false, 
                message: '激活码验证接口请求超时', 
                needInput: false 
            });
        });
        
        if (postData) {
            req.write(postData);
        }
        req.end();
    });
}

// 显示激活码输入对话框 - 使用更稳定的实现
function showActivationCodeDialog(message: string, currentConfig?: any): Promise<string | null> {
    return new Promise(async (resolve) => {
        console.log('开始创建激活码输入对话框...');
        
        const { dialog, BrowserWindow } = require('electron');
        
        try {
            // 方法1：使用原生对话框（推荐）
            const result = await dialog.showMessageBox({
                type: 'question',
                buttons: ['确认', '取消'],
                defaultId: 0,
                title: '激活码验证',
                message: '请输入激活码',
                detail: message + (currentConfig && currentConfig.activation_code ? 
                    `\n\n当前激活码: ${currentConfig.activation_code}\n有效期: ${currentConfig.expires_at || '未知'}` : ''),
                checkboxLabel: '记住激活码',
                checkboxChecked: false
            });
            
            if (result.response === 0) {
                // 用户点击确认，显示输入框
                const inputResult = await dialog.showSaveDialog({
                    title: '输入激活码',
                    defaultPath: '',
                    filters: [{ name: 'Text', extensions: ['txt'] }],
                    properties: ['showOverwriteConfirmation']
                });
                
                // 这个方法不太理想，让我们使用方法2
            }
            
            // 方法2：创建简单的HTML文件然后加载
            const fs = require('fs');
            const path = require('path');
            const os = require('os');
            
            // 创建临时HTML文件
            const tempDir = os.tmpdir();
            const tempHtmlPath = path.join(tempDir, 'activation_dialog.html');
            
            const htmlContent = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>激活码验证</title>
    <style>
        body { 
            font-family: 'Microsoft YaHei', Arial, sans-serif; 
            padding: 30px; 
            margin: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .container { 
            background: white;
            padding: 40px;
            border-radius: 12px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.3);
            max-width: 500px;
            width: 100%;
            text-align: center;
        }
        .title {
            color: #333;
            margin-bottom: 20px;
            font-size: 24px;
            font-weight: bold;
        }
        .message { 
            margin: 20px 0; 
            color: #d32f2f; 
            font-size: 16px;
            line-height: 1.6;
            text-align: left;
            background: #fff3cd;
            padding: 15px;
            border-radius: 8px;
            border-left: 4px solid #ffc107;
        }
        .config-info {
            background: #e3f2fd;
            padding: 15px;
            margin: 15px 0;
            border-radius: 8px;
            border-left: 4px solid #2196f3;
            text-align: left;
            font-size: 14px;
        }
        #activationCode { 
            width: 100%; 
            padding: 15px; 
            margin: 20px 0; 
            border: 2px solid #ddd;
            border-radius: 8px;
            font-size: 16px;
            outline: none;
            box-sizing: border-box;
            transition: border-color 0.3s;
        }
        #activationCode:focus {
            border-color: #2196f3;
            box-shadow: 0 0 0 3px rgba(33, 150, 243, 0.1);
        }
        .button-group {
            margin-top: 30px;
        }
        button { 
            padding: 12px 30px; 
            margin: 0 10px; 
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 16px;
            font-weight: bold;
            transition: all 0.3s;
            min-width: 100px;
        }
        .btn-confirm { 
            background: #4caf50; 
            color: white; 
        }
        .btn-confirm:hover {
            background: #45a049;
            transform: translateY(-2px);
        }
        .btn-cancel { 
            background: #f44336; 
            color: white; 
        }
        .btn-cancel:hover {
            background: #da190b;
            transform: translateY(-2px);
        }
        .btn-confirm:disabled {
            background: #cccccc;
            cursor: not-allowed;
            transform: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="title">🔐 激活码验证</div>
        <div class="message">${message}</div>
        ${currentConfig && currentConfig.activation_code ? `
        <div class="config-info">
            <strong>📋 当前激活码信息：</strong><br>
            激活码: <code>${currentConfig.activation_code}</code><br>
            有效期: <code>${currentConfig.expires_at || '未知'}</code>
        </div>
        ` : ''}
        <div>
            <input type="text" id="activationCode" placeholder="请输入新的激活码" maxlength="100" autocomplete="off" />
        </div>
        <div class="button-group">
            <button class="btn-confirm" id="confirmBtn" onclick="confirmInput()">✅ 确认</button>
            <button class="btn-cancel" onclick="cancelInput()">❌ 取消</button>
        </div>
    </div>
    <script>
        const { ipcRenderer } = require('electron');
        
        let isSubmitted = false;
        
        function confirmInput() {
            if (isSubmitted) return;
            
            const codeInput = document.getElementById('activationCode');
            const confirmBtn = document.getElementById('confirmBtn');
            const code = codeInput.value.trim();
            
            if (!code) {
                codeInput.style.borderColor = '#f44336';
                codeInput.focus();
                
                // 显示错误提示
                let errorDiv = document.getElementById('error-msg');
                if (!errorDiv) {
                    errorDiv = document.createElement('div');
                    errorDiv.id = 'error-msg';
                    errorDiv.style.color = '#f44336';
                    errorDiv.style.fontSize = '14px';
                    errorDiv.style.marginTop = '10px';
                    codeInput.parentNode.appendChild(errorDiv);
                }
                errorDiv.textContent = '⚠️ 请输入激活码';
                
                setTimeout(() => {
                    if (errorDiv) errorDiv.remove();
                    codeInput.style.borderColor = '#ddd';
                }, 3000);
                return;
            }
            
            isSubmitted = true;
            confirmBtn.disabled = true;
            confirmBtn.textContent = '⏳ 验证中...';
            
            console.log('用户输入激活码:', code);
            ipcRenderer.send('activation-code-input', code);
        }
        
        function cancelInput() {
            if (isSubmitted) return;
            isSubmitted = true;
            console.log('用户取消输入');
            ipcRenderer.send('activation-code-input', null);
        }
        
        // 监听回车键
        document.getElementById('activationCode').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                confirmInput();
            }
        });
        
        // 页面加载完成后自动聚焦输入框
        window.addEventListener('DOMContentLoaded', function() {
            console.log('激活码对话框页面加载完成');
            setTimeout(() => {
                const input = document.getElementById('activationCode');
                if (input) {
                    input.focus();
                    input.select();
                }
            }, 200);
        });
        
        // 防止页面意外关闭
        window.addEventListener('beforeunload', function(e) {
            if (!isSubmitted) {
                e.preventDefault();
                e.returnValue = '';
            }
        });
    </script>
</body>
</html>
`;
            
            // 写入临时文件
            fs.writeFileSync(tempHtmlPath, htmlContent, 'utf8');
            console.log('临时HTML文件已创建:', tempHtmlPath);
            
            // 创建窗口
            const inputWindow = new BrowserWindow({
                width: 600,
                height: 500,
                modal: true,
                resizable: false,
                show: false,
                autoHideMenuBar: true,
                webPreferences: {
                    nodeIntegration: true,
                    contextIsolation: false,
                    webSecurity: false
                }
            });
            
            console.log('激活码对话框窗口已创建');
            
            let isResolved = false;
            
            // 监听输入结果
            const { ipcMain } = require('electron');
            const inputHandler = (event, code) => {
                if (isResolved) return;
                isResolved = true;
                
                console.log('收到用户输入:', code ? '有输入' : '无输入或取消');
                ipcMain.removeListener('activation-code-input', inputHandler);
                
                // 清理临时文件
                try {
                    if (fs.existsSync(tempHtmlPath)) {
                        fs.unlinkSync(tempHtmlPath);
                        console.log('临时HTML文件已删除');
                    }
                } catch (cleanupError) {
                    console.warn('清理临时文件失败:', cleanupError);
                }
                
                if (!inputWindow.isDestroyed()) {
                    inputWindow.close();
                }
                resolve(code);
            };
            
            ipcMain.once('activation-code-input', inputHandler);
            
            // 窗口关闭处理
            inputWindow.on('closed', () => {
                if (isResolved) return;
                isResolved = true;
                
                console.log('激活码对话框窗口已关闭');
                ipcMain.removeListener('activation-code-input', inputHandler);
                
                // 清理临时文件
                try {
                    if (fs.existsSync(tempHtmlPath)) {
                        fs.unlinkSync(tempHtmlPath);
                    }
                } catch (cleanupError) {
                    console.warn('清理临时文件失败:', cleanupError);
                }
                
                resolve(null);
            });
            
            // 加载HTML文件
            console.log('开始加载对话框内容...');
            inputWindow.loadFile(tempHtmlPath);
            
            // 内容加载完成后显示窗口
            inputWindow.webContents.once('did-finish-load', () => {
                if (isResolved) return;
                
                console.log('对话框内容加载完成，显示窗口');
                inputWindow.show();
                inputWindow.focus();
                inputWindow.center();
            });
            
            // 错误处理
            inputWindow.webContents.on('did-fail-load', (event, errorCode, errorDescription) => {
                if (isResolved) return;
                isResolved = true;
                
                console.error('激活码对话框加载失败:', errorCode, errorDescription);
                ipcMain.removeListener('activation-code-input', inputHandler);
                
                // 清理临时文件
                try {
                    if (fs.existsSync(tempHtmlPath)) {
                        fs.unlinkSync(tempHtmlPath);
                    }
                } catch (cleanupError) {
                    console.warn('清理临时文件失败:', cleanupError);
                }
                
                if (!inputWindow.isDestroyed()) {
                    inputWindow.close();
                }
                resolve(null);
            });
            
            // 超时处理
            setTimeout(() => {
                if (!isResolved && !inputWindow.isDestroyed()) {
                    console.log('激活码对话框超时，自动关闭');
                    isResolved = true;
                    ipcMain.removeListener('activation-code-input', inputHandler);
                    inputWindow.close();
                    resolve(null);
                }
            }, 120000); // 2分钟超时
            
        } catch (error) {
            console.error('创建激活码对话框时发生错误:', error);
            resolve(null);
        }
    });
}

// 检查现有验证码是否有效
async function checkExistingActivationCode(): Promise<boolean> {
    console.log('开始检查现有激活码...');
    try {
        const http = require('http');
        const postData = JSON.stringify({});
        const options = {
            hostname: '127.0.0.1',
            port: 7072,
            path: '/client_config',
            method: 'POST',
            timeout: 5000,
            headers: {
                'Content-Type': 'application/json',
                'Content-Length': Buffer.byteLength(postData)
            }
        };
        
        return new Promise((resolve) => {
            console.log('正在向Go服务请求客户端配置...');
            const req = http.request(options, (res) => {
                console.log(`Go服务响应状态码: ${res.statusCode}`);
                let data = '';
                res.on('data', (chunk) => {
                    data += chunk;
                });
                
                res.on('end', () => {
                    console.log('Go服务响应数据:', data);
                    try {
                        const response = JSON.parse(data);
                        console.log('解析后的响应:', response);
                        
                        // 根据新的数据格式检查activation_valid字段
                        if (response.activation_valid === true) {
                            console.log('设备已激活，跳过验证窗口');
                            console.log('激活信息:', response.client_config);
                            resolve(true);
                            return;
                        } else if (response.activation_valid === false) {
                            console.log('设备未激活，需要显示验证窗口');
                            console.log('消息:', response.message);
                            resolve(false);
                            return;
                        } else {
                            console.log('响应格式不正确，activation_valid字段缺失或格式错误');
                        }
                    } catch (error) {
                        console.log('解析响应失败:', error);
                    }
                    console.log('没有有效的激活状态，需要显示验证窗口');
                    resolve(false);
                });
            });
            
            req.on('error', (error) => {
                console.log('检查现有激活码失败:', error.message);
                console.log('将显示激活码验证窗口');
                resolve(false);
            });
            
            req.on('timeout', () => {
                req.destroy();
                console.log('检查现有激活码超时');
                console.log('将显示激活码验证窗口');
                resolve(false);
            });
            
            // 发送POST数据
            req.write(postData);
            req.end();
        });
    } catch (error) {
        console.log('检查现有激活码异常:', error);
        console.log('将显示激活码验证窗口');
        return false;
    }
}

// 在文件最后添加（替换直接调用 initializeApp）
const { app: electronApp } = require('electron');

// 等待 Electron 应用准备就绪后再初始化
electronApp.whenReady().then(() => {
  console.log('Electron 应用已准备就绪，开始初始化...');
  initializeApp().catch(error => {
    console.error('应用初始化失败:', error);
    electronApp.quit();
  });
});

// 处理所有窗口关闭的情况
electronApp.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    electronApp.quit();
  }
});

// 处理应用激活（macOS）
electronApp.on('activate', () => {
  // 在 macOS 上，当点击 dock 图标并且没有其他窗口打开时，
  // 通常会重新创建一个窗口
});