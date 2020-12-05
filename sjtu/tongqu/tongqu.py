import requests
import json
import time 
from requests.packages.urllib3.exceptions import InsecureRequestWarning
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)


conf = json.load(open("config.json", encoding="utf8"))

url = "https://tongqu.sjtu.edu.cn/api/act/sign"

header = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36 Edg/87.0.664.41",
    "Origin": "https://tongqu.sjtu.edu.cn",
    "Cookie": conf["cookie"],  # 登录之后的cookie
    "X-Requested-With": "XMLHttpRequest",
    "Host": "tongqu.sjtu.edu.cn"
}

# 配合其网站解析，需要转成字符串
conf["data"]["user_sign_info"] = json.dumps(conf["data"]["user_sign_info"])

start = time.time()
counter = 1
while True:
    res = requests.post(url, data=conf["data"], headers=header, verify=False)
    print("正在进行第{}次尝试".format(counter))
    counter = counter + 1
    if(res.json()["error"] == 0):
        end = time.time()
        print("共用时: {}s".format(end-start))
        break
    else:
      print(res.json()["msg"])
