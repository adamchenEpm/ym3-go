#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import sys
import json
import io

# 【核心修复 1】：强制标准输入输出使用 UTF-8 编码，忽略系统默认的 GBK
sys.stdin = io.TextIOWrapper(sys.stdin.buffer, encoding='utf-8')
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

def main():
    try:
        # 现在读取的就是正确的 UTF-8 了
        input_data = sys.stdin.read().strip()

        if not input_data:
            return

        args = json.loads(input_data)
        city = args.get("city", "")

        weather_data = {
            "北京": "晴，25°C",
            "上海": "多云，22°C",
            "广州": "阵雨，28°C",
        }

        result = weather_data.get(city, f"抱歉，未找到{city}的天气信息")

        response = {
            "status": "success",
            "data": result
        }

        # 【核心修复 2】：ensure_ascii=False 保证中文原样输出，配合上面的 stdout 修改
        print(json.dumps(response, ensure_ascii=False))

    except Exception as e:
        print(json.dumps({"status": "error", "message": str(e)}, ensure_ascii=False))

if __name__ == "__main__":
    main()