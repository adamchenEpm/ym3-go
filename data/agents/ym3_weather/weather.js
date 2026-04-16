#!/usr/bin/env node
// -*- coding: utf-8 -*-

const fs = require('fs');

/**
 * 核心逻辑：从标准输入(stdin)读取JSON，处理后输出JSON
 * 这样可以规避 Windows/Linux 命令行转义字符的差异
 */
function main() {
    try {
        // 1. 同步读取标准输入 (文件描述符 0)
        // 使用 utf8 编码确保中文不乱码
        const inputData = fs.readFileSync(0, 'utf8').trim();

        if (!inputData) {
            console.log(JSON.stringify({ status: "error", message: "No input data provided" }));
            return;
        }

        // 2. 解析 Go 传过来的 JSON 参数
        const args = JSON.parse(inputData);
        const city = args.city || "";

        // 模拟企业级 Skill 的数据源
        const weatherData = {
            "北京": "晴，25°C，湿度40%",
            "上海": "多云，22°C，湿度65%",
            "广州": "阵雨，28°C，湿度80%",
        };

        const result = weatherData[city] || `抱歉，未找到${city}的天气信息`;

        // 3. 构造统一的返回格式
        // 生产环境务必返回标准 JSON，方便 Go 端的结构体(Struct)反序列化
        const response = {
            status: "success",
            data: {
                location: city,
                weather: result,
                source: "node_runtime"
            },
            timestamp: new Date().toISOString()
        };

        // 使用 process.stdout.write 确保输出纯净，没有多余的换行符
        process.stdout.write(JSON.stringify(response));

    } catch (err) {
        // 错误处理：即使脚本崩溃，也要返回 JSON 格式的错误信息
        const errorRes = {
            status: "error",
            message: err.message
        };
        process.stdout.write(JSON.stringify(errorRes));
        process.exit(1); // 非正常退出
    }
}

main();