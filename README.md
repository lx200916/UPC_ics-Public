# Credits
@AndyLiu24
ChanJH@chanjh.com
@yupotian


[![State-of-the-art Shitcode](https://img.shields.io/static/v1?label=State-of-the-art&message=Shitcode&color=7B5804)](https://github.com/trekhleb/state-of-the-art-shitcode)

# How To：
- 1.在教务系统-课表查询下载课表（未合并），得到kb.xls，打开(教务系统使用的文件格式不规范，因此打开时出现错误是正常的，忽略即可）
- 2.将kb.xls的sheet1的整个表格复制到项目中kb.xlsm的sheet1中，请保持复制前后单元格行列的对应位置相同，kb.xlsm中已附示例。
- 3.选择Excel程序中 "视图"--"宏"，选择"Sheet1.schedule"（也可以继续执行"Sheet1.Unique",进行去重操作），执行。（如弹出安全警告，请"允许编辑"或"启用宏"）
- 4.查看sheet2中的课程信息是否是本人的，如果不是请重新执行宏。
- 5.先执行  excelReader.py（在弹出窗口中查看xlsm文件格式是否正确，如果你确定kb.xlsm文件复制前后格式相同，直接输入0）
- 6.执行main.py 输入信息后即可生成class.ics文件。
- 7.自行想办法导入手机日历。

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
