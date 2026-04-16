#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import sys, json

def main():
    # 从命令行参数获取 JSON
    if len(sys.argv) < 2:
        print("Error: missing arguments")
        return
    print(sys.argv[1])
    print(sys.argv[2])
    print(sys.argv[3])
    city = sys.argv[3]
    # 模拟天气数据
    weather_data = {
        "北京": "晴，25°C，湿度40%",
        "上海": "多云，22°C，湿度65%",
        "广州": "阵雨，28°C，湿度80%",
    }
    result = weather_data.get(city, f"抱歉，未找到{city}的天气信息")
    print(result)

if __name__ == "__main__":
    main()