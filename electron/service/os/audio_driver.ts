import { spawn } from 'child_process';
import * as path from 'path';
import { pub } from '../../class/public';

class AudioDriverService {
  private installerPath: string;

  constructor() {
    // 音频驱动安装包的路径
    this.installerPath = path.resolve(pub.get_resource_path(), 'exe/VBCABLE_Driver_Pack45/VBCABLE_Setup_x64.exe');
  }

  /**
   * 启动音频驱动安装程序
   */
  async installAudioDriver(): Promise<{ success: boolean; message: string }> {
    try {
      console.log('启动音频驱动安装程序:', this.installerPath);

      // 直接启动安装程序
      const child = spawn(this.installerPath, [], {
        detached: true,
        stdio: 'ignore'
      });

      // 分离进程
      child.unref();

      return {
        success: true,
        message: '音频驱动安装程序已启动'
      };
    } catch (error) {
      console.error('启动音频驱动安装程序失败:', error);
      return {
        success: false,
        message: `启动失败: ${error.message}`
      };
    }
  }
}

export const audioDriverService = new AudioDriverService();