import xlrd
import copy
data=xlrd.open_workbook("122.xls")
table=data.sheets()[0]
j=2

classinfos=[]
def infonew():
    for i in range(1,table.ncols):
        for j in range(3,table.nrows-1):
            info=table.col(i)[j].value.split()
            # print(str(info))
            clasinfo = {}
            if (str(info) != '[]'):
                clasinfo['className'] = info[0]
                clasinfo['weekday'] = 7 if i == 1 else i - 1
                clasinfo['classroom'] = info[3]

                clasinfo['teacher'] = info[1]
                clasinfo['classTime'] = j-2
                weeks = info[2].split('[')[0]
                clasinfo['week'] = {}
                if '-' in weeks:
                    clasinfo['week']['startWeek'] = weeks[0:weeks.find('-')]
                    clasinfo['week']['endWeek'] = weeks[weeks.find('-') +1:]
                    classinfos.append(clasinfo)
                elif ',' in weeks:
                    week = weeks.split(',')
                    for each in week:
                        classinfo=copy.deepcopy(clasinfo)
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
                print(info[0])
