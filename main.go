package main

import (
	"log"
	"net/http"
	"os"
	"server-monitoring/file"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DateLayOut      = "2006-01-02"
	TimeStampLayout = "2006-01-02 15:04:05"
)

var (
	servers = []string{
		"https://google.com",
		"https://youtube.com",
		"https://facebook.com",
		"https://baidu.com",
		"https://wikipedia.org",
		"https://yahoo.com",
		"https://tmall.com",
		"https://amazon.com",
		"https://twitter.com",
		"https://live.com",
		"https://instagram.com",
	}

	//ReportColumnHeader Csv file column headers
	ReportColumnHeader = []string{
		"Server URL",
		"Response Status",
		"Response Time",
	}
	//ReportColumnHeader Csv file column headers
	ErrorReportColumnHeader = []string{
		"Server URL",
		"Response Status",
		"Error",
	}
)

type Report struct {
	Server         string
	ResponseStatus int
	ResponseTime   float64
	Err            error
}

func main() {

	duration, _ := time.ParseDuration(os.Getenv("DURATION"))
	for {
		log.Println("starting job")
		pingserver()
		time.Sleep(duration)
	}
}

func pingserver() {

	var (
		reschan = make(chan Report)
		report  []Report
		wg      = sync.WaitGroup{}
		err     error
	)

	for _, i := range servers {
		wg.Add(1)
		go func(url string) {
			var statuscode int
			start := time.Now()
			defer func() {
				wg.Done()
				reschan <- Report{
					Server:         url,
					ResponseStatus: statuscode,
					ResponseTime:   time.Since(start).Seconds(),
					Err:            err,
				}
			}()
			response, er := http.Get(url)
			if er != nil {
				log.Printf("Error Received while calling %s: %s", url, er.Error())
				err = er
			} else {
				defer response.Body.Close()
				statuscode = response.StatusCode
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < len(servers); i++ {
		select {
		case c := <-reschan:
			report = append(report, c)
		}
	}
	log.Println("Generating Reports")
	Writereports(report)
}

func Writereports(report []Report) {

	var (
		successreport []string
		errReport     []string
	)

	SuccessRecords := [][]string{
		ReportColumnHeader,
	}

	SuccessRecordsMoretime := [][]string{
		ReportColumnHeader,
	}

	ErrorRecords := [][]string{
		ErrorReportColumnHeader,
	}

	WritePath := os.Getenv("WritePath")

	for _, k := range report {
		if k.ResponseStatus != http.StatusBadGateway && k.ResponseStatus != 0 {
			successreport = append(successreport, k.Server)
			successreport = append(successreport, strconv.Itoa(k.ResponseStatus))
			successreport = append(successreport, strconv.FormatFloat(k.ResponseTime, 'f', -1, 64))

			if k.ResponseTime > 1 {
				SuccessRecordsMoretime = append(SuccessRecordsMoretime, successreport)
			} else {
				SuccessRecords = append(SuccessRecords, successreport)
			}
			successreport = []string{}
		} else {
			errReport = append(errReport, k.Server)
			errReport = append(errReport, strconv.Itoa(k.ResponseStatus))
			if k.Err != nil {
				errReport = append(errReport, k.Err.Error())
			}
			ErrorRecords = append(ErrorRecords, errReport)
			errReport = []string{}
		}
	}

	fileNameSplit := strings.Split(time.Now().Format(DateLayOut), "-")
	formattedFileName := strings.Join(fileNameSplit, "/") + "/"

	SuccessFileName := WritePath + "/" + formattedFileName + "/" + time.Now().Format(TimeStampLayout) + "_avaliable_lesstime.csv"
	if err := file.WriteOutputFile(SuccessFileName, SuccessRecords); err != nil {
		log.Printf("Error Generating available file : %s", err.Error())
		return
	}

	MoretimeFileName := WritePath + "/" + formattedFileName + "/" + time.Now().Format(TimeStampLayout) + "_avaliable_moretime.csv"
	if err := file.WriteOutputFile(MoretimeFileName, SuccessRecordsMoretime); err != nil {
		log.Printf("Error Generating available file : %s", err.Error())
		return
	}

	ErrorFileName := WritePath + "/" + formattedFileName + "/" + time.Now().Format(TimeStampLayout) + "_unavailable.csv"
	if err := file.WriteOutputFile(ErrorFileName, ErrorRecords); err != nil {
		log.Printf("Error Generating unavailable file : %s", err.Error())
		return
	}
	log.Printf("Files Generated successfully")
}
