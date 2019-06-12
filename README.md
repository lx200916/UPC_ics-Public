# Introduction
一个集各位大佬之长的整合版本并兼容最新教务系统的xls课表转ics小工具，因为用纯Python重写，消除了之前VBA泄露信息的隐患，将仓库开源。Python萌新菜鸡，追求又不是不能用。代码很烂，欢迎拍砖.
# Credits
@AndyLiu24
ChanJH@chanjh.com
@yupotian

# How To：
- 1.在教务系统» 我的课表 » 学期理论课表下载（打印）课表（未合并），得到学生个人课表+学号.xls
- 2.执行app.py,按照提示访问locallost:5000或IP:5000，选择并上传课表，浏览器会自动开始下载class.ics。
- 3.自行想办法导入手机日历。

# Warning
- 1.请先安装Python 3.4以上版本，并导入xlrd，flask模块（pip install xlrd或pip3 install xlrd），否则会无法执行程序或闪退
- 2.导入时请选择谷歌日历（不需要登录账号，用谷歌日历打开ics文件，即可自动导入系统日历）
- 3.对于iOS设备导入可能会麻烦些，建议参考http://chanjh.com/post/software/0031 中提到的方法导入。（辣鸡iOS）
- 4.再次强调！Python版本非常重要，Python3.4以上中imp被替换为importlib，程序已作出更改。
# Changelog
- 1.已知在Outlook及Flyme日历中会出现课程时间早一个小时的问题，理论上已经修复，请自测。
- 2.加入了老师名称及地点字段（不一定可以显示，视所使用的日历APP）
- 3.加入了一般来看没什么卵用的服务器功能，纯粹是萌新对Flask好奇驱使.....HTML页面位于/static，欢迎看不下去的大佬自行修改，反正做的很烂就是啦！（又不是不能用🤣）
## 3.修复了当晚上三节连堂时产生两个课表问题（注意！！！仅修复了9,10,11节连堂的问题，如果有其他的连堂情况请勿使用！！）
- 4.欢迎反馈BUG
- 5.再次感谢@AndyLiu24、ChanJH@chanjh.com大佬
- 6.以及@yupotian小姐姐的测试，RUA！
# TODO
- 1.也许会发布剥离Flask的版本但是目前没有需要，所以咕咕咕。
- 2.增强安全性和可靠性
- 3.重写HTML
- 4.To Be Continued....