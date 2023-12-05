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

movie_list=['深圳 金逸影城 中心店', '深圳 星美国际影城 北城店',
            '北京 橙天嘉禾影城 购物广场店', '北京 万达影城 中心店',
            '上海 大地影院 北城店', '深圳 金逸影城 西城店', '上海 中影国际影城 中心店',
            '广州 万达影城 西城店', '广州 星美国际影城 东城店', '深圳 大地影院 西城店',
            '上海 星美国际影城 购物广场店', '深圳 金逸影城 北城店', '北京 大地影院 北城店',
            '广州 万达影城 欢乐园店', '上海 橙天嘉禾影城 东城店', '北京 金逸影城 北城店',
            '深圳 橙天嘉禾影城 北城店', '上海 金逸影城 西城店', '深圳 大地影院 欢乐园店',
            '上海 万达影城 中心店']


