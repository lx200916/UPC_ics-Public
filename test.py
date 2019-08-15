import xlrd
import re
import json
import copy
import sys
sys.setrecursionlimit(100000)
def infoold():
    data = xlrd.open_workbook("kb.xlsm")
    table = data.sheets()[0]
    global classinfos
    classinfos = []

    for i in range(1, table.ncols):
        for j in range(3, table.nrows - 1):
            info = table.col(i)[j].value.split()
            # print(str(info))

    # print(len(info))
            if (str(info) != '[]'):
                clasinfo = {}
                clasinfo['className'] = info[0]
                clasinfo['weekday'] = 7 if i == 1 else i - 1
                clasinfo['classroom'] = info[2].split(']')[-1]
                # clasinfo['classTime']=re.findall(r'\[.]',info[2])
                # print(info[2].split(']')[-1])
                clasinfo['teacher'] = info[1]
                clasinfo['classTime'] = re.findall(r'\[(.*)\]', str(info[2]))[0].replace('-', '').replace('节', '')
                weeks = info[2].split('[')[0]
                clasinfo['week'] = {}
                # print(weeks)
                if '-' in weeks:
                    clasinfo['week']['startWeek'] = weeks[0:weeks.find('-')]
                    clasinfo['week']['endWeek'] = weeks[weeks.find('-') + 1:weeks.find('周')]
                    classinfos.append(clasinfo)
                elif ',' in weeks:
                    week = weeks.replace('周', '').split(',')
                    for each in week:
                        classinfo = copy.deepcopy(clasinfo)
                        if '-' in weeks:
                            classinfo['week']['startWeek'] = weeks[0:weeks.find('-')]
                            classinfo['week']['endWeek'] = weeks[weeks.find('-') + 1:weeks.find('周')]
                            classinfos.append(classinfo)
                        else:
                            classinfo['week']['startWeek'] = each
                            classinfo['week']['endWeek'] = each
                            classinfos.append(classinfo)
                else:
                    each = weeks.replace('周', '')
                    clasinfo['week']['startWeek'] = each
                    clasinfo['week']['endWeek'] = each
                    classinfos.append(clasinfo)
            if (len(info) > 4):
                for i in range(1,int(len(info)/4)):
                    info[0:3]=info[4*i:4*i+3]
                    infoold()
                    print("reload")



def infonew():
    data = xlrd.open_workbook("122.xls")
    table = data.sheets()[0]

    global classinfos
    classinfos=[]

    for i in range(1, table.ncols):
        for j in range(3, table.nrows - 1):
            info = table.col(i)[j].value.split()
            # print(str(info))
            clasinfo = {}
            if (str(info) != '[]'):
                clasinfo['className'] = info[0]
                clasinfo['weekday'] = 7 if i == 1 else i - 1
                clasinfo['classroom'] = info[3]

                clasinfo['teacher'] = info[1]
                clasinfo['classTime'] = j - 2
                weeks = info[2].split('[')[0]
                clasinfo['week'] = {}
                if '-' in weeks:
                    clasinfo['week']['startWeek'] = weeks[0:weeks.find('-')]
                    clasinfo['week']['endWeek'] = weeks[weeks.find('-') + 1:]
                    classinfos.append(clasinfo)
                elif ',' in weeks:
                    week = weeks.split(',')
                    for each in week:
                        classinfo = copy.deepcopy(clasinfo)
                        if '-' in weeks:
                            classinfo['week']['startWeek'] = weeks[0:weeks.find('-')]
                            classinfo['week']['endWeek'] = weeks[weeks.find('-') + 1:]
                            classinfos.append(classinfo)
                        else:
                            classinfo['week']['startWeek'] = each
                            classinfo['week']['endWeek'] = each
                            classinfos.append(classinfo)
                else:
                    each = weeks
                    clasinfo['week']['startWeek'] = each
                    clasinfo['week']['endWeek'] = each
                    classinfos.append(clasinfo)
            if (len(info) > 5):
                for i in range(1, int(len(info) / 5)):
                    info = info[4:]
                    print(len(info))
                    infonew()



if __name__ == '__main__':

    infonew()

    print(classinfos)




