import { ref } from "vue";
interface PlayItem {
    content: any;
    filename: string;
    play_mode: string;
  }
const playList = ref<PlayItem[]>([])//播放列表
const start = ref(false)//启动开关

//注册模块，
const registerModules = (modules) => {
    //循环模块列表 干活
    start.value = true
    for (const module of modules) {
        //循环读取
        if (module.trigger_conditions.includes("ExecuteLoop")) {
            ExecuteLoop(module)
        }
        //弹幕评论

        //送礼物

        //进入直播间
    }
    playListConsumption()
}

//消费playList
const playListConsumption= async () => {
    do {
        if (playList.value.length > 0) { //队列消费
            const item = playList.value.shift() //出队
            if (item) {
                await play_task_voice_api(item.filename, item.play_mode); //播放
            }
        }
        await new Promise(resolve => setTimeout(resolve, 100))   
    } while (start.value);
}

//循环模块处理
const ExecuteLoop= async (module) => {
    let index = 0 //当前播放索引
    const script_content_len = module.script_content.length //脚本长度
    do {
        if (module.read_step == "random") { //随机
            if (index == 0) { 
                module.script_content = shuffleArray(module.script_content)//将数组 打乱
            }
        }
        let content = module.script_content[index] 
        while (playList.value.length >= 2) { //等待队列 消费完探入
            await new Promise(resolve => setTimeout(resolve, 100))
        }
        const newFileName = Date.now() + '.wav' //生成文件名
        await generate_wav_api(
            content, //文本
            "en",  //语种
            newFileName, //生成文件名
            "", //音色文件名
            content.speed, //生成速度
            content.volume) //生成音量
        
        playList.value.push({ //入队
            content: content, //内容
            filename: newFileName, //文件名
            play_mode: "serial", //播放模式
        })    
        if (index == script_content_len - 1) { //循环
            index = 0
        } else {
            index = index + 1
        }
    } while (start.value);
}
// 更标准的Fisher-Yates洗牌算法
function shuffleArray(array) {
    for (let i = array.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [array[i], array[j]] = [array[j], array[i]];
    }
    return array;
}

// http://192.168.1.10:7073/generate_wav
// 生成wav文件
const generate_wav_api = async (_text:string,
    _language:string,
    _filename:string,
    _speaker_wav:string,
    _speed: number,
    _volume: number) => {
    const response = await fetch('http://192.168.1.10:7073/generate_wav', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
            text: _text,
            language: _language,
            filename : _filename,
            speaker_wav: _speaker_wav,
            speed: _speed,// 1,
            volume: _volume// 0.5
         }),
    });
    await response.json();
}

// http://192.168.1.10:7073/play_task_voice
// 播放任务
const play_task_voice_api = async (_filename:string,play_mode:string) => {
    const response = await fetch('http://192.168.1.10:7073/play_task_voice', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
            filename : _filename,
            play_mode : play_mode
         }),
    });
    await response.json();
}

export {
    playList,
    registerModules,
}