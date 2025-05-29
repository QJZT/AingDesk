import undici from 'undici';
const { close } = undici;

// DisposableSendApi 压力测试
const testDisposableSend = async () => {
    const testParams = {
        model: 'qwen2.5',
        parameters: '7b',
        user_content: '内置3种游戏模式！挑战模式、音乐模式、宝宝模式，全家都能玩到停不下来！',
        system_prompt: '你是一个专业的全球化直播文案优化助手。请按照以下要求处理和改善我提供的口播文案：\n\n1. 语言转换：…种：\n   - 简体中文\n\n请直接返回出改善后的对应语种文案，不要任何特殊符号，不要任何额外说明。',
        supplier_name: 'ollama',
        trigger_conditions: ['BarrageComment'], // 添加弹幕触发条件
    };

    const startTime = Date.now();
    const totalRequests = 100000; // 总请求数
    const results = [];
    let successCount = 0;
    let failCount = 0;
    const batchSize = 10; // 每批并发请求数
    const batchDelay = 100; // 批次间延迟（毫秒）

    console.log(`开始压力测试，总请求数: ${totalRequests}, 批次大小: ${batchSize}`);

    // 分批发送请求
    for (let i = 0; i < totalRequests; i += batchSize) {
        const batchPromises = [];
        for (let j = 0; j < batchSize && i + j < totalRequests; j++) {
            const requestId = i + j + 1;
            batchPromises.push(
                (async () => {
                    try {
                        const controller = new AbortController();
                        const timeoutId = setTimeout(() => controller.abort(), 60000); // 60秒超时

                        const response = await fetch('http://127.0.0.1:7072/disposable_send', {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json',
                            },
                            body: JSON.stringify({
                                ...testParams,
                                user_content: `${testParams.user_content} (请求 ${requestId})`, // 区分请求
                            }),
                            signal: controller.signal,
                        });

                        clearTimeout(timeoutId);

                        if (!response.ok) {
                            const errorText = await response.text();
                            throw new Error(`HTTP 错误: ${response.status}, 详情: ${errorText}`);
                        }

                        const data = await response.text();
                        if (data.length > 10 * 1024 * 1024) {
                            throw new Error('响应数据过大');
                        }

                        let result = data;
                        if (data.includes('</think>')) {
                            const endIndex = data.indexOf('</think>');
                            result = data.substring(endIndex + 7).replace(/\n/g, '');
                        }

                        results.push({
                            requestId,
                            status: 'success',
                            data: result,
                        });
                        successCount++;
                        console.log(`请求 ${requestId} 成功`);
                    } catch (error) {
                        results.push({
                            requestId,
                            status: 'error',
                            error: error.message,
                        });
                        failCount++;
                        console.error(`请求 ${requestId} 失败:`, error);
                    }
                })()
            );
        }

        // 等待当前批次完成
        await Promise.all(batchPromises);
        // 批次间延迟
        if (i + batchSize < totalRequests) {
            await new Promise(resolve => setTimeout(resolve, batchDelay));
        }
    }

    const endTime = Date.now();
    const totalTime = endTime - startTime;

    // 输出测试结果
    console.log('\n测试结果统计：');
    console.log(`总请求数: ${totalRequests}`);
    console.log(`成功请求数: ${successCount}`);
    console.log(`失败请求数: ${failCount}`);
    console.log(`总耗时: ${totalTime}ms`);
    console.log(`平均响应时间: ${(totalTime / totalRequests).toFixed(2)}ms`);
    console.log(`成功率: ${((successCount / totalRequests) * 100).toFixed(2)}%`);

    return {
        totalRequests,
        successCount,
        failCount,
        totalTime,
        averageTime: totalTime / totalRequests,
        successRate: (successCount / totalRequests) * 100,
        results,
    };
};

// 运行测试
testDisposableSend()
    .then(result => {
        console.log('\n详细测试结果：', result);
    })
    .catch(error => {
        console.error('测试执行失败：', error);
    })
