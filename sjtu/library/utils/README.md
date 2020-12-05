### 上海交通大学图书馆图书信息搜索

使用方式
```bash
go run main.go -req xxx
# use -h option for help
```
命令行参数
```
-filename string
      保存文件名 (default "lib.json")
-req string
      搜索关键字
```
运行可得json格式文件，如下
```json
[
  {
     "book_id": "TP312GO/7 2020",
     "title": {
       "url": "http://opac.lib.sjtu.edu.cn:80/F/4HPGYIVGLRI897S1GXRPPVTKTMM1KFN8ECQ3N7D3I34JRMLCRI-01796?func=full-set-set&set_number=032450&set_entry=000007&format=999",
       "name": "Go专家编程"
     },
     "author": "任洪彩",
     "year": "2020",
     "info": {
       "url": "http://opac.lib.sjtu.edu.cn:80/F/4HPGYIVGLRI897S1GXRPPVTKTMM1KFN8ECQ3N7D3I34JRMLCRI-01797?func=item-global&doc_library=SJT01&doc_number=001160817&year=&volume=&sub_library=LTSKJ ",
       "name": "主馆(2/0)"
     },
     "type": "BK"
   },
   {
     "book_id": "TP312/J57 2020",
     "title": {
       "url": "http://opac.lib.sjtu.edu.cn:80/F/4HPGYIVGLRI897S1GXRPPVTKTMM1KFN8ECQ3N7D3I34JRMLCRI-01800?func=full-set-set&set_number=032450&set_entry=000008&format=999",
       "name": "Go并发编程实战"
     },
     "author": "汪明",
     "year": "2020",
     "info": {
       "url": "http://opac.lib.sjtu.edu.cn:80/F/4HPGYIVGLRI897S1GXRPPVTKTMM1KFN8ECQ3N7D3I34JRMLCRI-01801?func=item-global&doc_library=SJT01&doc_number=001170746&year=&volume=&sub_library=LTSKJ ",
       "name": "主馆(2/0)"
     },
     "type": "BK"
   },
   ...
]
```