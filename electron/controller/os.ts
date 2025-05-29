import fs from 'fs';
import path from 'path';
// 在文件顶部确保已导入 shell
import { app as electronApp, dialog, shell, ipcMain } from 'electron';

// 在文件末尾添加
ipcMain.handle('open-external-url', async (event, url) => {
  try {
    await shell.openExternal(url);
    return { success: true };
  } catch (error) {
    console.error(`打开URL错误: ${error}`);
    throw error;
  }
});
import { windowService } from '../service/os/window';
import { exec, spawn } from 'child_process';

/**
 * example
 * @class
 */
class OsController {

  /**
   * All methods receive two parameters
   * @param args Parameters transmitted by the frontend
   * @param event - Event are only available during IPC communication. For details, please refer to the controller documentation
   */

  /**
   * Message prompt dialog box
   */
  messageShow(): string {
    dialog.showMessageBoxSync({
      type: 'info', // "none", "info", "error", "question" 或者 "warning"
      title: 'Custom Title',
      message: 'Customize message content',
      detail: 'Other additional information'
    })
  
    return 'Opened the message box';
  }

  /**
   * Message prompt and confirmation dialog box
   */
  messageShowConfirm(): string {
    const res = dialog.showMessageBoxSync({
      type: 'info',
      title: 'Custom Title',
      message: 'Customize message content',
      detail: 'Other additional information',
      cancelId: 1, // Index of buttons used to cancel dialog boxes
      defaultId: 0, // Set default selected button
      buttons: ['confirm', 'cancel'], 
    })
    let data = (res === 0) ? 'click the confirm button' : 'click the cancel button';
  
    return data;
  }

  /**
   * Select Directory
   */
  selectFolder() {
    const filePaths = dialog.showOpenDialogSync({
      properties: ['openDirectory', 'createDirectory']
    });

    if (!filePaths) {
      return ""
    }

    return filePaths[0];
  } 

  /**
   * open directory
   */
  openDirectory(args: { id: any }): boolean {
    const { id } = args;
    if (!id) {
      return false;
    }
    let dir = '';
    if (path.isAbsolute(id)) {
      dir = id;
    } else {
      dir = electronApp.getPath(id);
    }

    shell.openPath(dir);
    return true;
  }

  /**
   * Select Picture
   */
  selectPic(): string | null {
    const filePaths = dialog.showOpenDialogSync({
      title: 'select pic',
      properties: ['openFile'],
      filters: [
        { name: 'Images', extensions: ['jpg', 'png', 'gif'] },
      ]
    });
    if (!filePaths) {
      return null
    }
    
    try {
      const data = fs.readFileSync(filePaths[0]);
      const pic =  'data:image/jpeg;base64,' + data.toString('base64');
      return pic;
    } catch (err) {
      console.error(err);
      return null;
    }
  }   

  /**
   * Open a new window
   */
  createWindow(args: any): any {
    const wcid = windowService.createWindow(args);
    return wcid;
  }
  
  /**
   * Get Window contents id
   */
  getWCid(args: any): any {
    const wcid = windowService.getWCid(args);
    return wcid;
  }

  /**
   * Realize communication between two windows through the transfer of the main process
   */
  window1ToWindow2(args: any): void {
    windowService.communicate(args);
    return;
  }

  /**
   * Realize communication between two windows through the transfer of the main process
   */
  window2ToWindow1(args: any): void {
    windowService.communicate(args);
    return;
  }

  /**
   * Create system notifications
   */
  sendNotification(args: { title?: string; subtitle?: string; body?: string; silent?: boolean }, event: any): boolean {
    const { title, subtitle, body, silent} = args;

    const options: any = {};
    if (title) {
      options.title = title;
    }
    if (subtitle) {
      options.subtitle = subtitle;
    }
    if (body) {
      options.body = body;
    }
    if (silent !== undefined) {
      options.silent = silent;
    }
    windowService.createNotification(options, event);

    return true
  }   
}
OsController.toString = () => '[class OsController]';

export default OsController;

// 在现有的ipcMain处理程序中添加
// 修改 execute-bat-file 处理函数
ipcMain.handle('execute-bat-file', async (event, batFilePath) => {
  return new Promise((resolve, reject) => {
    // 确保批处理文件存在
    if (!fs.existsSync(batFilePath)) {
      return reject(new Error(`批处理文件不存在: ${batFilePath}`));
    }
    
    // 获取批处理文件所在目录作为工作目录
    const workingDir = path.dirname(batFilePath);
    
    // 使用 spawn 替代 exec
    const childProcess = spawn('cmd.exe', ['/c', batFilePath], {
      windowsVerbatimArguments: true,
      detached: true,         // 分离子进程
      windowsHide: false,     // 显示窗口，便于调试
      cwd: workingDir,        // 设置工作目录
      env: process.env,       // 传递完整的环境变量
      shell: true             // 使用shell
    });
    
    let stdoutData = '';
    let stderrData = '';
    
    // 收集输出信息
    if (childProcess.stdout) {
      childProcess.stdout.on('data', (data) => {
        const dataStr = data.toString();
        stdoutData += dataStr;
        console.log(`批处理输出: ${dataStr}`);
      });
    }
    
    if (childProcess.stderr) {
      childProcess.stderr.on('data', (data) => {
        const dataStr = data.toString();
        stderrData += dataStr;
        console.error(`批处理错误: ${dataStr}`);
      });
    }
    
    childProcess.on('error', (error) => {
      console.error(`执行批处理文件错误: ${error}`);
      reject(error);
    });
    
    childProcess.on('exit', (code) => {
      console.log(`批处理进程退出，退出码: ${code}`);
      if (code === 0) {
        resolve({ success: true, stdout: stdoutData, stderr: stderrData });
      } else {
        reject(new Error(`进程退出，退出码: ${code}\n${stderrData}`));
      }
    });
  });
});
    
