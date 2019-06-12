# coding: utf-8
#!/usr/bin/python

import importlib
import time, datetime
import json
from random import Random
import sys
import xlrd
import re
import copy
import getopt, sys
flag=0

__author__ = 'ChanJH'
__site__ = 'chanjh.com'

checkFirstWeekDate = 0
checkReminder = 1

YES = 0
NO = 1

DONE_firstWeekDate = time.time()
DONE_reminder = ""
DONE_EventUID = ""
DONE_UnitUID = ""
DONE_CreatedTime = ""
DONE_ALARMUID = ""


classTimeList = {44}
classInfoList = []

def main():
	global path
	path=sys.argv[1]
	basicSetting();
	uniteSetting();
	classInfoHandle();
	icsCreateAndSave();


def save(string):
     f = open(path+"/class.ics", 'wb')
     f.write(string.encode("utf-8"))
     f.close()

def icsCreateAndSave():
	icsString = "BEGIN:VCALENDAR\nMETHOD:PUBLISH\nVERSION:2.0\nX-WR-CALNAME:课程表\nPRODID:-//Apple Inc.//Mac OS X 10.14//EN\nX-APPLE-CALENDAR-COLOR:#FC4208\nX-WR-TIMEZONE:Asia/Shanghai\nCALSCALE:GREGORIAN\n"
	global classTimeList, DONE_ALARMUID, DONE_UnitUID
	eventString = ""
	for classInfo in classInfoList :
		i = int((classInfo["classTime"]))

		className = classInfo["className"]+"|"+list(classTimeList.values())[i]["name"]+"|"+classInfo["classroom"]
		endTime = list(classTimeList.values())[i]["endTime"]
		startTime = list(classTimeList.values())[i]["startTime"]
		teacher = classInfo["teacher"]
		room = classInfo["classroom"]
		index = 0
		for date in classInfo["date"]:
			eventString = eventString+"BEGIN:VEVENT\nCREATED:"+classInfo["CREATED"]
			eventString = eventString+"\nUID:"+classInfo["UID"][index]
			eventString = eventString+"\nDTEND;TZID=Asia/Shanghai:"+date+"T"+endTime
			eventString = eventString+"00\nTRANSP:OPAQUE\nX-APPLE-TRAVEL-ADVISORY-BEHAVIOR:AUTOMATIC\nSUMMARY:"+className+"\nDESCRIPTION:"+teacher + "\nLOCATION:"+room
			eventString = eventString+"\nDTSTART;TZID=Asia/Shanghai:"+date+"T"+startTime+"00"
			eventString = eventString+"\nDTSTAMP:"+DONE_CreatedTime
			eventString = eventString+"\nSEQUENCE:0\nBEGIN:VALARM\nX-WR-ALARMUID:"+DONE_ALARMUID
			eventString = eventString+"\nUID:"+DONE_UnitUID
			eventString = eventString+"\nTRIGGER:"+DONE_reminder
			eventString = eventString+"\nDESCRIPTION:"+teacher+"\nACTION:DISPLAY\nEND:VALARM\nEND:VEVENT\n"

			index += 1
	icsString = icsString + eventString + "END:VCALENDAR"
	save(icsString)
	print("icsCreateAndSave")

def classInfoHandle():
	global classInfoList
	global DONE_firstWeekDate
	i = 0

	for classInfo in classInfoList :
		# 具体日期计算出来

		startWeek = json.dumps(classInfo["week"]["startWeek"])
		endWeek = json.dumps(classInfo["week"]["endWeek"])
		weekday = float(json.dumps(classInfo["weekday"]))
		
		dateLength = float((int(startWeek) - 1) * 7)
		startDate = datetime.datetime.fromtimestamp(int(time.mktime(DONE_firstWeekDate))) + datetime.timedelta(days = dateLength + weekday - 1)
		string = startDate.strftime('%Y%m%d')

		dateLength = float((int(endWeek) - 1) * 7)
		endDate = datetime.datetime.fromtimestamp(int(time.mktime(DONE_firstWeekDate))) + datetime.timedelta(days = dateLength + weekday - 1)
		
		date = startDate
		dateList = []
		dateList.append(string)
		i = NO
		while (i):
			date = date + datetime.timedelta(days = 7.0)
			if(date > endDate):
				i = YES
			else:
				string = date.strftime('%Y%m%d')
				dateList.append(string)
		classInfo["date"] = dateList

		# 设置 UID
		global DONE_CreatedTime, DONE_EventUID
		CreateTime()
		classInfo["CREATED"] = DONE_CreatedTime
		classInfo["DTSTAMP"] = DONE_CreatedTime
		UID_List = []
		for date  in dateList:
			UID_List.append(UID_Create())
		classInfo["UID"] = UID_List
	print("classInfoHandle")

def UID_Create():
	return random_str(20) + "&Chanjh.com"


def CreateTime():
	# 生成 CREATED
	global DONE_CreatedTime
	date = datetime.datetime.now().strftime("%Y%m%dT%H%M%S")
	DONE_CreatedTime = date + "Z"
	# 生成 UID
	# global DONE_EventUID
	# DONE_EventUID = random_str(20) + "&Chanjh.com"

	print("CreateTime")

def uniteSetting():
	# 
	global DONE_ALARMUID
	DONE_ALARMUID = random_str(30) + "&Chanjh.com"
	# 
	global DONE_UnitUID
	DONE_UnitUID = random_str(20) + "&Chanjh.com"
	print("uniteSetting")

def setClassTime():
	data = []
	with open('conf_classTime.json', 'r',encoding='UTF-8') as f:
		data = json.load(f)
	global classTimeList
	classTimeList = data["classTime"]
	print("setclassTime")
	
def setClassInfo():
	data = xlrd.open_workbook(path+"/kb.xls")
	table = data.sheets()[0]
	global classinfos
	classinfos = []
	for i in range(1, table.ncols):
		for j in range(3, table.nrows - 1):
			info = table.col(i)[j].value.split()
			infonew(info,i,j)
	# elif(flag==2):
	# 	infoold()

	global classInfoList
	classInfoList=classinfos
	print(classInfoList)

# def infoold():
#     data = xlrd.open_workbook("kb.xlsm")
#     table = data.sheets()[0]
# 	global classinfos
#
#
#     for i in range(1, table.ncols):
#         for j in range(3, table.nrows - 1):
#             info = table.col(i)[j].value.split()
#             # print(str(info))
#
#     # print(len(info))
#             if (str(info) != '[]'):
#                 clasinfo = {}
#                 clasinfo['className'] = info[0]
#                 clasinfo['weekday'] = 7 if i == 1 else i - 1
#                 clasinfo['classroom'] = info[2].split(']')[-1]
#                 # clasinfo['classTime']=re.findall(r'\[.]',info[2])
#                 # print(info[2].split(']')[-1])
#                 clasinfo['teacher'] = info[1]
#                 clasinfo['classTime'] = re.findall(r'\[(.*)\]', str(info[2]))[0].replace('-', '').replace('节', '')
#                 weeks = info[2].split('[')[0]
#                 clasinfo['week'] = {}
#                 # print(weeks)
#                 if '-' in weeks:
#                     clasinfo['week']['startWeek'] = weeks[0:weeks.find('-')]
#                     clasinfo['week']['endWeek'] = weeks[weeks.find('-') + 1:weeks.find('周')]
#                     classinfos.append(clasinfo)
#                 elif ',' in weeks:
#                     week = weeks.replace('周', '').split(',')
#                     for each in week:
#                         classinfo = copy.deepcopy(clasinfo)
#                         if '-' in weeks:
#                             classinfo['week']['startWeek'] = weeks[0:weeks.find('-')]
#                             classinfo['week']['endWeek'] = weeks[weeks.find('-') + 1:weeks.find('周')]
#                             classinfos.append(classinfo)
#                         else:
#                             classinfo['week']['startWeek'] = each
#                             classinfo['week']['endWeek'] = each
#                             classinfos.append(classinfo)
#                 else:
#                     each = weeks.replace('周', '')
#                     clasinfo['week']['startWeek'] = each
#                     clasinfo['week']['endWeek'] = each
#                     classinfos.append(clasinfo)
#             if (len(info) > 4):
#                 for i in range(1,int(len(info)/4)):
#                     info[0:3]=info[4*i:4*i+3]
#                     infoold(info,i)
#                     print("reload")
#

def infonew(info,i,j):
	clasinfo = {}
	if (str(info) != '[]'):
		clasinfo['className'] = info[0]
		clasinfo['weekday'] = 7 if i == 1 else i - 1
		clasinfo['classroom'] = info[3]

		clasinfo['teacher'] = info[1]
		clasinfo['classTime'] = j - 3
		weeks = info[2].split('[')[0]
		clasinfo['week'] = {}
		if '-' in weeks:
			clasinfo['week']['startWeek'] = int(weeks[0:weeks.find('-')])
			clasinfo['week']['endWeek'] = int(weeks[weeks.find('-') + 1:])
			classinfos.append(clasinfo)
		elif ',' in weeks:
			week = weeks.split(',')
			for each in week:
				classinfo = copy.deepcopy(clasinfo)
				if '-' in weeks:
					classinfo['week']['startWeek'] = int(weeks[0:weeks.find('-')])
					classinfo['week']['endWeek'] = int(weeks[weeks.find('-') + 1:])
					classinfos.append(classinfo)
				else:
					classinfo['week']['startWeek'] = int(each)
					classinfo['week']['endWeek'] = int(each)
					classinfos.append(classinfo)
		else:
			each = int(weeks)
			clasinfo['week']['startWeek'] = each
			clasinfo['week']['endWeek'] = each
			classinfos.append(clasinfo)
		if (len(info) > 5):
			for k in range(1, int(len(info) / 4)):
					info = info[4:]
					print(len(info))
					infonew(info,i,j)


def setFirstWeekDate(firstWeekDate):
	global DONE_firstWeekDate
	DONE_firstWeekDate = time.strptime(firstWeekDate,'%Y%m%d')
	print("setFirstWeekDate:",DONE_firstWeekDate)

def setReminder(reminder):
	global DONE_reminder
	reminderList = ["-PT10M","-PT30M","-PT1H","-PT2H","-P1D"]
	if(reminder == "1"):
		DONE_reminder = reminderList[0]
	elif(reminder == "2"):
		DONE_reminder = reminderList[1]
	elif(reminder == "3"):
		DONE_reminder = reminderList[2]
	elif(reminder == "4"):
		DONE_reminder = reminderList[3]
	elif(reminder == "5"):
		DONE_reminder = reminderList[4]
	else:
		DONE_reminder = "NULL"


	print("setReminder",reminder)

def checkReminder(reminder):
	# TODO

	print("checkReminder:",reminder)
	List = ["0","1","2","3","4","5"]
	for num in List:
		if (reminder == num):
			return YES
	return NO

def checkFirstWeekDate(firstWeekDate):
	# 长度判断
	if(len(firstWeekDate) != 8):
		return NO;
	
	year = firstWeekDate[0:4]
	month = firstWeekDate[4:6]
	date = firstWeekDate[6:8]
	dateList = [31,29,31,30,31,30,31,31,30,31,30,31]

	# 年份判断
	if(int(year) < 1970):
		return NO
	# 月份判断
	if(int(month) == 0 or int(month) > 12):
		return NO;
	# 日期判断
	if(int(date) > dateList[int(month)-1]):
		return NO;

	print("checkFirstWeekDate:",firstWeekDate)
	return YES

def basicSetting():
	info = "欢迎使用课程表生成工具。\n接下来你需要设置一些基础的信息方便生成数据\n"
	inform = "请选择课表来源：\n1.教务系统学生课表“我的课表”下载\n2.教务系统“各类课表查询”查询下载\n"
	flag=1


	
	info = "请设置第一周的星期一日期(如：20190225):\n"
	firstWeekDate = '20190225'
	checkInput(checkFirstWeekDate, firstWeekDate)
	
	info = "正在配置上课时间信息……\n"

	try :
		setClassTime()

	except :
		sys_exit()

	info = "正在配置课堂信息……\n"

	try :
		setClassInfo()

	except :
		sys_exit()

	info = "正在配置提醒功能，请输入数字选择提醒时间\n【0】不提醒\n【1】上课前 10 分钟提醒\n【2】上课前 30 分钟提醒\n【3】上课前 1 小时提醒\n【4】上课前 2 小时提醒\n【5】上课前 1 天提醒\n"
	reminder = '2'
	checkInput(checkReminder, reminder)
def checkInput(checkType, input):
	if(checkType == checkFirstWeekDate):
		if (checkFirstWeekDate(input)):
			info = "输入有误，请重新输入第一周的星期一日期(如：20160905):\n"
			firstWeekDate = input(info)
			checkInput(checkFirstWeekDate, firstWeekDate)
		else:
			setFirstWeekDate(input)
	elif(checkType == checkReminder):
		if(checkReminder(input)):
			info = "输入有误，请重新输入\n【1】上课前 10 分钟提醒\n【2】上课前 30 分钟提醒\n【3】上课前 1 小时提醒\n【4】上课前 2 小时提醒\n【5】上课前 1 天提醒\n"
			reminder ='2'
			# reminder = input(info)
			checkInput(checkReminder, reminder)
		else:
			setReminder(input)

	else:
		print("程序出错了……")
		exit()

def random_str(randomlength):
    str = ''
    chars = 'AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789'
    length = len(chars) - 1
    random = Random()
    for i in range(randomlength):
        str+=chars[random.randint(0, length)]
    return str
def sys_exit():
	print("配置文件错误，请检查。\n")
	sys.exit()
importlib.reload(sys);

main()