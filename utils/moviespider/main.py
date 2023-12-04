import requests
import json

url = "https://api.wmdb.tv/api/v1/top?type=Douban&skip=0&limit=250&lang=Cn"

payload = {}

headers = {
    'Accept': 'application/json, text/plain, */*',
    'Origin': 'https://www.wmdb.tv',
    'Referer': 'https://www.wmdb.tv/',
    'Sec-Fetch-Dest': 'empty',
    'Sec-Fetch-Mode': 'cors',
    'Sec-Fetch-Site': 'cross-site'
}

response = requests.request("GET", url, headers=headers, data=payload)

# 将响应文本解析为 JSON
response_data = json.loads(response.text)

# 将 JSON 数据写入 data.txt 文件
with open('data.txt', 'w', encoding='utf-8') as f:
    json.dump(response_data, f, ensure_ascii=False, indent=4)

print("Data saved to data.txt")
