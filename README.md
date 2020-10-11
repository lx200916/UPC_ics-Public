# Credits
@AndyLiu24
ChanJH@chanjh.com
@yupotian
# Notes
此版本已下线,新版本使用Gin框架重写,见 `Golang` 分支

[![State-of-the-art Shitcode](https://img.shields.io/static/v1?label=State-of-the-art&message=Shitcode&color=7B5804)](https://github.com/trekhleb/state-of-the-art-shitcode)


# Warning
- 1.请先安装Python 3.4以上版本，并导入xlrd模块（pip install xlrd或pip3 install xlrd），否则会无法执行程序或excelReader.py闪退
- 2.导入时请选择谷歌日历（不需要登录账号，用谷歌日历打开ics文件，即可自动导入系统日历）
- 3. 对于iOS设备导入可能会麻烦些，建议参考http://chanjh.com/post/software/0031 中提到的方法导入。（辣鸡iOS）
- 4.再次强调！Python版本非常重要，Python3.4以上中imp被替换为importlib，程序已作出更改。
# Changelog
- 1.已知在Outlook及Flyme日历中会出现课程时间早一个小时的问题，理论上已经修复，请自测。
- 2.加入了老师名称及地点字段（不一定可以显示，视所使用的日历APP）
## 3.修复了当晚上三节连堂时产生两个课表问题（注意！！！仅修复了9,10,11节连堂的问题，如果有其他的连堂情况请勿使用！！）
- 4.欢迎反馈BUG
- 5.再次感谢@AndyLiu24、ChanJH@chanjh.com大佬
- 6.以及@yupotian小姐姐的测试，RUA！
