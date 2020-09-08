package main

import (
	"encoding/json"
	"fmt"
	"github.com/extrame/xls"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	_ "net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)
import _ "github.com/gin-gonic/gin"

var (
	t  = 1        // 总票数
	mu sync.Mutex // 互斥锁
	n  sync.WaitGroup
)

type Time struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
type aiinfo struct {
	Data struct {
		CsId int `json:"csId"`
	}
}

type xiaoai struct {
	Weeks    string `json:"weeks"`
	Name     string `json:"name"`
	Position string `json:"position"`
	Sections string `json:"sections"`
	Teacher  string `json:"teacher"`
	Style    string `json:"style"`
	Day      int    `json:"day"`
	CsId     int    `json:"csId"`
	UserId   int    `json:"userId"`
	Devid    string `json:"deviceId"`
}
type Class struct {
	ClassName string `json:"className"`
	Week      struct {
		StartWeek int `json:"startWeek"`
		EndWeek   int `json:"endWeek"`
	} `json:"week"`
	Weekday   int    `json:"weekday"`
	ClassTime int    `json:"classTime"`
	Classroom string `json:"classroom"`
	Teacher   string `json:"teacher"`
	UID       []string
	Date      []string
	Repeat    int
}

var Weekd = []string{"SU", "MO", "TU", "WE", "TH", "FR", "SA"}
var icsString = "BEGIN:VCALENDAR\nMETHOD:PUBLISH\nVERSION:2.0\nX-WR-CALNAME:课程表\nPRODID:-//Apple Inc.//Mac OS X 10.14//EN\nX-APPLE-CALENDAR-COLOR:#FC4208\nX-WR-TIMEZONE:Asia/Shanghai\nCALSCALE:GREGORIAN\n"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var DONE_reminder = ""
var DONE_EventUID = ""
var DONE_UnitUID = ""
var DONE_CreatedTime = ""
var DONE_ALARMUID = ""

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")
	r.StaticFS("/static", http.Dir("./static"))
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.NoRoute(func(context *gin.Context) {
		context.HTML(404, "404.html", nil)
	})
	r.POST("/", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		path, _ := os.Getwd()
		path = filepath.Join(path, "upload", strconv.FormatInt(time.Now().Unix(), 10)+RandStringBytes(2))
		_ = os.MkdirAll(path, os.ModePerm)
		_ = c.SaveUploadedFile(file, path+"/kb.xls")
		//xlss(path)
		ics1(path)
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=class.ics")) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.File(path + "/class.ics")
		return

	})
	r.GET("/set", func(c *gin.Context) {
		c.HTML(200, "12.html", nil)
		return
	})
	r.GET("/mi", func(c *gin.Context) {
		c.HTML(200, "1.html", nil)
		return
	})
	r.POST("/mi", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		path, _ := os.Getwd()
		path = filepath.Join(path, "upload", strconv.FormatInt(time.Now().Unix(), 10)+RandStringBytes(2))
		_ = os.MkdirAll(path, os.ModePerm)
		_ = c.SaveUploadedFile(file, path+"/kb.xls")
		//xlss(path)
		userid, _ := strconv.Atoi(c.Request.Form["miid"][0])
		fmt.Print(userid)
		mu.Lock()

		t = icsm(path, userid)

		defer mu.Unlock()

		return

	})
	r.POST("/set", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		path, _ := os.Getwd()
		path = filepath.Join(path, "upload", strconv.FormatInt(time.Now().Unix(), 10)+RandStringBytes(2))
		_ = os.MkdirAll(path, os.ModePerm)
		_ = c.SaveUploadedFile(file, path+"/kb.xls")
		//xlss(path)
		date := c.Request.Form["date"][0]
		reminder, _ := strconv.Atoi(c.Request.Form["reminder"][0])
		ics(path, date, reminder-1, 0)

		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=class.ics")) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.File(path + "/class.ics")
		return
	})
	_ = r.Run("0.0.0.0:5000")

	//st, _ := json.Marshal(classlist)
	//fmt.Print(string(st))
}
func gettime() []Time {
	file, _ := os.Open("conf_classTime.json")
	var timelist []Time
	buf := make([]byte, 1024)
	total, _ := file.Read(buf)
	fmt.Println(json.Unmarshal(buf[:total], &timelist))
	return timelist
}

func ics1(path string) {
	ics(path, "20200907", 1, 0)
}
func icsm(path string, types int) int {
	ics(path, "20200907", 1, types)
	return 1
}
func ics(path string, date string, reminder int, types int) {
	var csid int
	if types != 0 {
		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://i.ai.mi.com/course/relation", strings.NewReader("{\n	\"deviceId\": \"76578245c6a6aabcb963ac1953930ddd\",\n	\"userId\": "+strconv.Itoa(types)+",\n	\"bind\": 1,\n	\"sync\": 2\n}"))

		if err != nil {
			fmt.Println(err)
		}
		req.Header.Add("Host", "i.ai.mi.com")
		req.Header.Add("Cookie", "_ga=GA1.2.1780394317.1591374303; _gid=GA1.2.739547335.1591374303; _gat_gtag_UA_148568844_1=1")
		req.Header.Add("access-control-allow-origin", "true")
		req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 10; ONEPLUS A6000 Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.96 Mobile Safari/537.36")
		req.Header.Add("content-type", "application/json")
		req.Header.Add("accept", "*/*")
		req.Header.Add("origin", "https://i.ai.mi.com")
		req.Header.Add("x-requested-with", "com.miui.voiceassist")
		req.Header.Add("sec-fetch-site", "same-origin")
		req.Header.Add("sec-fetch-mode", "cors")
		req.Header.Add("sec-fetch-dest", "empty")
		req.Header.Add("referer", "https://i.ai.mi.com/h5/precache/ai-schedule/")
		req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")

		res, err := client.Do(req)
		defer res.Body.Close()
		_, err = ioutil.ReadAll(res.Body)
		req, err = http.NewRequest("GET", "https://i.ai.mi.com/course/setting?deviceId=76578245c6a6aabcb963ac1953930ddd&userId="+strconv.Itoa(types), nil)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Add("Host", "i.ai.mi.com")
		req.Header.Add("Cookie", "_ga=GA1.2.1648326863.1591374801; _gid=GA1.2.1697428902.1591374801")
		req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 10; ONEPLUS A6000 Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.96 Mobile Safari/537.36")
		req.Header.Add("accept", "*/*")
		req.Header.Add("x-requested-with", "com.miui.voiceassist")
		req.Header.Add("sec-fetch-site", "same-origin")
		req.Header.Add("sec-fetch-mode", "cors")
		req.Header.Add("sec-fetch-dest", "empty")
		req.Header.Add("referer", "https://i.ai.mi.com/h5/precache/ai-schedule/")
		req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")

		res, err = client.Do(req)
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		var v aiinfo
		_ = json.Unmarshal(body, &v)

		firstweek, _ := time.Parse("20060102", date)
		c := time.Now()
		csid = v.Data.CsId
		payload := strings.NewReader("{\"csId\":" + strconv.Itoa(csid) + ",\"userId\":" + strconv.Itoa(types) + ",\"deviceId\":\"76578245c6a6aabcb963ac1953930ddd\",\"presentWeek\":" + strconv.Itoa(int(math.Floor(c.Sub(firstweek).Hours()/(24*7))+1)) + ",\"totalWeek\":20,\"isWeekend\":1,\"morningNum\":2,\"afternoonNum\":2,\"nightNum\":2,\"school\":\"\",\"confirm\":0,\"sections\":[{\"section\":1,\"startTime\":\"08:00\",\"endTime\":\"09:50\"},{\"section\":2,\"startTime\":\"10:10\",\"endTime\":\"12:00\"},{\"section\":3,\"startTime\":\"14:00\",\"endTime\":\"15:50\"},{\"section\":4,\"startTime\":\"16:10\",\"endTime\":\"18:00\"},{\"section\":5,\"startTime\":\"19:00\",\"endTime\":\"20:50\"},{\"section\":6,\"startTime\":\"21:00\",\"endTime\":\"21:50\"}],\"serverTime\":1591459498104}")

		req, err = http.NewRequest("PUT", "https://i.ai.mi.com/course/setting", payload)

		if err != nil {
			fmt.Println(err)
		}
		req.Header.Add("Host", "i.ai.mi.com")
		req.Header.Add("Cookie", "_ga=GA1.2.1780394317.1591374303; _gid=GA1.2.739547335.1591374303; _gat_gtag_UA_148568844_1=1")
		req.Header.Add("access-control-allow-origin", "true")
		req.Header.Add("origin", "https://i.ai.mi.com")
		req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 10; ONEPLUS A6000 Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.96 Mobile Safari/537.36")
		req.Header.Add("content-type", "application/json")
		req.Header.Add("accept", "*/*")
		req.Header.Add("x-requested-with", "com.miui.voiceassist")
		req.Header.Add("sec-fetch-site", "same-origin")
		req.Header.Add("sec-fetch-mode", "cors")
		req.Header.Add("sec-fetch-dest", "empty")
		req.Header.Add("referer", "https://i.ai.mi.com/h5/precache/ai-schedule/")
		req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")

		_, _ = client.Do(req)

	}

	var classlist []Class
	reminderList := []string{"-PT10M", "-PT30M", "-PT1H", "-PT2H", "-P1D"}
	DONE_ALARMUID = RandStringBytes(30) + "&UPC.edu"
	DONE_UnitUID = RandStringBytes(20) + "&UPC.edu"
	if xlFile, err := xls.Open(path+"/kb.xls", "utf-8"); err == nil {
		//fmt.Println(xlFile.GetSheet(0).Row(5).Col(0))
		sheet1 := xlFile.GetSheet(0)
		for i := 3; i <= (int(sheet1.MaxRow)); i++ {
			row := sheet1.Row(i)
			for j := 1; j <= 7; j++ {
				infos := strings.Split(row.Col(j), "\n")
				//fmt.Print(infos)

				for k := 1; k <= len(infos)/6; k++ {
					time := infos[3+6*k-6]
					var i1, i3 int
					var classtime int
					if infos[4+6*k-6] == "[03-04-05节]" {
						classtime = 7
					} else {
						classtime = i - 3
					}
					if strings.Contains(time, "(") {
						time = strings.Split(time, "(")[0]

					}
					if strings.Contains(time, "-") {
						if strings.Contains(time, ",") {
							week := strings.Split(time, ",")
							for _, i2 := range week {
								if strings.Contains(i2, "-") {
									i1, _ = strconv.Atoi(strings.Split(i2, "-")[0])
									i3, _ = strconv.Atoi(strings.Split(i2, "-")[1])
									class := Class{
										ClassName: infos[1+6*k-6],
										Week: struct {
											StartWeek int `json:"startWeek"`
											EndWeek   int `json:"endWeek"`
										}{i1, i3},
										Weekday:   j - 1,
										ClassTime: classtime,
										Classroom: infos[5+6*k-6],
										Teacher:   infos[2+6*k-6],
										Repeat:    1,
									}
									classlist = append(classlist, class)

								} else {
									i1, _ = strconv.Atoi(i2)
									i3, _ = strconv.Atoi(i2)
									class := Class{
										ClassName: infos[1+6*k-6],
										Week: struct {
											StartWeek int `json:"startWeek"`
											EndWeek   int `json:"endWeek"`
										}{i1, i3},
										Weekday:   j - 1,
										ClassTime: classtime,
										Classroom: infos[5+6*k-6],
										Teacher:   infos[2+6*k-6],
									}
									classlist = append(classlist, class)
								}

							}

						} else {
							i1, _ = strconv.Atoi(strings.Split(time, "-")[0])
							i3, _ = strconv.Atoi(strings.Split(time, "-")[1])
							class := Class{
								ClassName: infos[1+6*k-6],
								Week: struct {
									StartWeek int `json:"startWeek"`
									EndWeek   int `json:"endWeek"`
								}{i1, i3},
								Weekday:   j - 1,
								ClassTime: classtime,
								Classroom: infos[5+6*k-6],
								Teacher:   infos[2+6*k-6],
								Repeat:    1,
							}
							classlist = append(classlist, class)

						}

					} else if strings.Contains(time, ",") {
						if time == "3,5,7,9,11,13,15,17" {
							class := Class{
								ClassName: infos[1+6*k-6],
								Week: struct {
									StartWeek int `json:"startWeek"`
									EndWeek   int `json:"endWeek"`
								}{3, 17},
								Weekday:   j - 1,
								ClassTime: classtime,
								Classroom: infos[5+6*k-6],
								Teacher:   infos[2+6*k-6],
								Repeat:    2,
							}
							classlist = append(classlist, class)

						} else if time == "2,4,6,8,10,12,14,16" {

							class := Class{
								ClassName: infos[1+6*k-6],
								Week: struct {
									StartWeek int `json:"startWeek"`
									EndWeek   int `json:"endWeek"`
								}{2, 16},
								Weekday:   j - 1,
								ClassTime: classtime,
								Classroom: infos[5+6*k-6],
								Teacher:   infos[2+6*k-6],
								Repeat:    3,
							}
							classlist = append(classlist, class)
						} else {
							week := strings.Split(time, ",")
							for _, i2 := range week {
								if strings.Contains(i2, "-") {
									i1, _ = strconv.Atoi(strings.Split(time, "-")[0])
									i3, _ = strconv.Atoi(strings.Split(time, "-")[1])
									class := Class{
										ClassName: infos[1+6*k-6],
										Week: struct {
											StartWeek int `json:"startWeek"`
											EndWeek   int `json:"endWeek"`
										}{i1, i3},
										Weekday:   j - 1,
										ClassTime: classtime,
										Classroom: infos[5+6*k-6],
										Teacher:   infos[2+6*k-6],
										Repeat:    1,
									}
									classlist = append(classlist, class)

								} else {
									i1, _ = strconv.Atoi(i2)
									i3, _ = strconv.Atoi(i2)
									class := Class{
										ClassName: infos[1+6*k-6],
										Week: struct {
											StartWeek int `json:"startWeek"`
											EndWeek   int `json:"endWeek"`
										}{i1, i3},
										Weekday:   j - 1,
										ClassTime: classtime,
										Classroom: infos[5+6*k-6],
										Teacher:   infos[2+6*k-6],
									}
									classlist = append(classlist, class)
								}

							}

						}
					} else {
						i1, _ = strconv.Atoi(time)
						i3, _ = strconv.Atoi(time)
						class := Class{
							ClassName: infos[1+6*k-6],
							Week: struct {
								StartWeek int `json:"startWeek"`
								EndWeek   int `json:"endWeek"`
							}{i1, i3},
							Weekday:   j - 1,
							ClassTime: classtime,
							Classroom: infos[5+6*k-6],
							Teacher:   infos[2+6*k-6],
						}
						classlist = append(classlist, class)
					}

				}

			}

		}
	}
	times := gettime()
	var DONE_firstWeekDate = date
	firstweek, _ := time.Parse("20060102", DONE_firstWeekDate)
	eventString := ""
	for _, class := range classlist {

		//if class.Weekday==0 {
		//	class.Weekday=-1
		//}
		if types != 0 {
			xiaomi(class, csid, types)
		} else {

			duration, _ := time.ParseDuration(strconv.Itoa((class.Week.StartWeek*7+class.Weekday-8)*24) + "h")
			startdate := firstweek.Add(duration)
			if class.Repeat == 0 {
				duration, _ = time.ParseDuration(strconv.Itoa((class.Week.EndWeek*7+class.Weekday-8)*24) + "h")
				enddate := firstweek.Add(duration)
				flag := true
				duration, _ = time.ParseDuration("168h")
				class.Date = append(class.Date, startdate.Format("20060102"))
				class.UID = append(class.UID, RandStringBytes(20)+"&UPC.edu")
				for flag {
					startdate = startdate.Add(duration)
					if startdate.After(enddate) {
						flag = false
					} else {
						class.Date = append(class.Date, startdate.Format("20060102"))
						class.UID = append(class.UID, RandStringBytes(20)+"&UPC.edu")

					}

				}

				for i, dates := range class.Date {
					eventString = eventString + "BEGIN:VEVENT\nCREATED:20190327T075414Z\nUID:" + class.UID[i]
					eventString += "\nDTEND;TZID=Asia/Shanghai:" + dates + "T" + times[class.ClassTime].EndTime
					eventString += "00\nTRANSP:OPAQUE\nX-APPLE-TRAVEL-ADVISORY-BEHAVIOR:AUTOMATIC\nSUMMARY:" + class.ClassName + "|" + times[class.ClassTime].Name + "\nDESCRIPTION:" + class.Teacher + "\nLOCATION:" + class.Classroom
					eventString += "\nDTSTART;TZID=Asia/Shanghai:" + dates + "T" + times[class.ClassTime].StartTime + "00"
					eventString = eventString + "\nDTSTAMP:20190327T075414Z"
					eventString = eventString + "\nSEQUENCE:0\nBEGIN:VALARM\nX-WR-ALARMUID:" + DONE_ALARMUID
					eventString = eventString + "\nUID:" + DONE_UnitUID
					eventString = eventString + "\nTRIGGER:" + reminderList[reminder]
					eventString = eventString + "\nDESCRIPTION:" + class.Teacher + "\nACTION:DISPLAY\nEND:VALARM\nEND:VEVENT\n"

				}
			} else if class.Repeat == 1 {
				eventString = eventString + "BEGIN:VEVENT\nCREATED:20190327T075414Z\nUID:" + RandStringBytes(20) + "&UPC.edu"
				eventString += "\nDTEND;TZID=Asia/Shanghai:" + startdate.Format("20060102") + "T" + times[class.ClassTime].EndTime
				eventString += "00\nTRANSP:OPAQUE\nX-APPLE-TRAVEL-ADVISORY-BEHAVIOR:AUTOMATIC\nSUMMARY:" + class.ClassName + "|" + times[class.ClassTime].Name + "\nDESCRIPTION:" + class.Teacher + "\nLOCATION:" + class.Classroom
				eventString += "\nDTSTART;TZID=Asia/Shanghai:" + startdate.Format("20060102") + "T" + times[class.ClassTime].StartTime + "00"
				eventString = eventString + "\nDTSTAMP:20190327T075414Z\nRRULE:FREQ=WEEKLY;COUNT=" + strconv.Itoa(class.Week.EndWeek-class.Week.StartWeek+1)
				eventString = eventString + "\nSEQUENCE:0\nBEGIN:VALARM\nX-WR-ALARMUID:" + DONE_ALARMUID
				eventString = eventString + "\nUID:" + DONE_UnitUID
				eventString = eventString + "\nTRIGGER:" + reminderList[reminder]
				eventString = eventString + "\nDESCRIPTION:" + class.Teacher + "\nACTION:DISPLAY\nEND:VALARM\nEND:VEVENT\n"

			} else {
				eventString = eventString + "BEGIN:VEVENT\nCREATED:20190327T075414Z\nUID:" + RandStringBytes(20) + "&UPC.edu"
				eventString += "\nDTEND;TZID=Asia/Shanghai:" + startdate.Format("20060102") + "T" + times[class.ClassTime].EndTime
				eventString += "00\nTRANSP:OPAQUE\nX-APPLE-TRAVEL-ADVISORY-BEHAVIOR:AUTOMATIC\nSUMMARY:" + class.ClassName + "|" + times[class.ClassTime].Name + "\nDESCRIPTION:" + class.Teacher + "\nLOCATION:" + class.Classroom
				eventString += "\nDTSTART;TZID=Asia/Shanghai:" + startdate.Format("20060102") + "T" + times[class.ClassTime].StartTime + "00"
				eventString = eventString + "\nDTSTAMP:20190327T075414Z\nRRULE:FREQ=WEEKLY;COUNT=8;INTERVAL=2;"
				eventString = eventString + "\nSEQUENCE:0\nBEGIN:VALARM\nX-WR-ALARMUID:" + DONE_ALARMUID
				eventString = eventString + "\nUID:" + DONE_UnitUID
				eventString = eventString + "\nTRIGGER:" + reminderList[reminder]
				eventString = eventString + "\nDESCRIPTION:" + class.Teacher + "\nACTION:DISPLAY\nEND:VALARM\nEND:VEVENT\n"

			}

		}
	}

	if types != 0 {
		clean(csid, types)
	} else {
		icsString1 := icsString + eventString + "END:VCALENDAR"

		out, _ := os.OpenFile(path+"/class.ics", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)

		fmt.Print(fmt.Fprint(io.Writer(out), icsString1))

	}

}
func clean(csid int, userid int) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://i.ai.mi.com/course/irrelevant", strings.NewReader("{\"deviceId\":\"76578245c6a6aabcb963ac1953930ddd\",\"userId\":"+strconv.Itoa(userid)+"}"))

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Host", "i.ai.mi.com")
	req.Header.Add("Cookie", "_ga=GA1.2.1780394317.1591374303; _gid=GA1.2.739547335.1591374303; _gat_gtag_UA_148568844_1=1")
	req.Header.Add("access-control-allow-origin", "true")
	req.Header.Add("origin", "https://i.ai.mi.com")
	req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 10; ONEPLUS A6000 Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.96 Mobile Safari/537.36")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "*/*")
	req.Header.Add("x-requested-with", "com.miui.voiceassist")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://i.ai.mi.com/h5/precache/ai-schedule/")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	_, _ = client.Do(req)
	req, err = http.NewRequest("DELETE", "https://i.ai.mi.com/course/courseInfo?deviceId=76578245c6a6aabcb963ac1953930ddd&userId="+strconv.Itoa(userid)+"&csId="+strconv.Itoa(csid), strings.NewReader("[object Object]"))

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Host", "i.ai.mi.com")
	req.Header.Add("Cookie", "_ga=GA1.2.1780394317.1591374303; _gid=GA1.2.739547335.1591374303; _gat_gtag_UA_148568844_1=1")
	req.Header.Add("origin", "https://i.ai.mi.com")
	req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 10; ONEPLUS A6000 Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.96 Mobile Safari/537.36")
	req.Header.Add("content-type", "text/plain;charset=UTF-8")
	req.Header.Add("accept", "*/*")
	req.Header.Add("x-requested-with", "com.miui.voiceassist")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://i.ai.mi.com/h5/precache/ai-schedule/")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")

	_, _ = client.Do(req)

}
func xiaomi(class Class, csid int, userid int) {
	ai := xiaoai{Name: class.ClassName, Position: class.Classroom, Teacher: class.Teacher, Style: "{\"color\":\"#7D7AEA\",\"background\":\"#E5E5FB\"}", Sections: strconv.Itoa(class.ClassTime + 1), Day: class.Weekday, CsId: csid, UserId: userid, Devid: "76578245c6a6aabcb963ac1953930ddd"}
	if ai.Day == 0 {
		ai.Day = 7
	}
	tmp := ""
	for i := class.Week.StartWeek; i < class.Week.EndWeek; i++ {
		//slice1 = append(slice1, strconv.Itoa(i))
		tmp += strconv.Itoa(i) + ","
	}
	tmp += strconv.Itoa(class.Week.EndWeek)
	ai.Weeks = tmp
	buf, _ := json.Marshal(ai)
	url := "https://i.ai.mi.com/course/courseInfo"
	method := "POST"
	fmt.Print(string(buf))
	payload := strings.NewReader(string(buf))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Host", "i.ai.mi.com")
	req.Header.Add("Cookie", "_ga=GA1.2.1780394317.1591374303; _gid=GA1.2.739547335.1591374303; _gat_gtag_UA_148568844_1=1")
	req.Header.Add("access-control-allow-origin", "true")
	req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 10; ONEPLUS A6000 Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.96 Mobile Safari/537.36")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "*/*")
	req.Header.Add("origin", "https://i.ai.mi.com")
	req.Header.Add("x-requested-with", "com.miui.voiceassist")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://i.ai.mi.com/h5/precache/ai-schedule/")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

}
