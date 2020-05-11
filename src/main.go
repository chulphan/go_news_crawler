package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/korean"
)

type Article struct {
	title     string
	pressDate string
	publisher string
	content   string
}

const (
	baseURL = "https://search.naver.com/search.naver"
	// url     = "https://search.naver.com/search.naver?where=news&query=%EC%82%BC%EC%84%B1%EC%A0%84%EC%9E%90&sm=tab_opt&sort=0&photo=0&field=0&reporter_article=&pd=3&ds=2020.04.25&de=2020.04.26&docid=&nso=so%3Ar%2Cp%3Afrom20200425to20200426%2Ca%3Aall&mynews=0&refresh_start=0&related=0"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("검색어를 입력해주세요: ")
	searchWord, _ := reader.ReadString('\n')
	searchWord = strings.TrimSpace(searchWord)
	fmt.Println(searchWord)

	fmt.Print("검색 시작날짜 (2020.05.09): ")
	startSearchDate, _ := reader.ReadString('\n')
	startSearchDate = strings.TrimSpace(startSearchDate)
	fmt.Println(startSearchDate)

	fmt.Print("검색 끝날짜 (2020.05.09): ")
	endSearchDate, _ := reader.ReadString('\n')
	endSearchDate = strings.TrimSpace(endSearchDate)
	fmt.Println(endSearchDate)

	// initialUrl := strings.Join(strings.Split("https://search.naver.com/search.naver?where=news&query="+searchWord+"&sm=tab_opt&sort=0&photo=0&field=0&reporter_article=&pd=3&reporter_article=&pd=3&docid=&mynews=0&refresh_start=0&related=0\n", "\n"), "")

	initialUrl := "https://search.naver.com/search.naver?where=news&query=" + searchWord + "&sm=tab_opt&sort=0&photo=0&field=0&reporter_article=&pd=3&reporter_article=&pd=3&docid=&mynews=0&refresh_start=0&related=0"

	if searchWord == "" {
		fmt.Println("검색어 없음.. 종료")
	}

	if startSearchDate != "" && endSearchDate != "" {
		joinStartSearchDate := strings.Join(strings.Split(startSearchDate, "."), "")
		joinEndSearchDate := strings.Join(strings.Split(endSearchDate, "."), "")

		ds := "&ds=" + startSearchDate
		de := "&de=" + endSearchDate
		nso := "&nso=so:r,p:from" + joinStartSearchDate
		to := "to" + joinEndSearchDate
		all := ",a:all"

		dateQueryStrArr := []string{ds, de, nso, to, all}

		containDateQuery := strings.Join(dateQueryStrArr, "")

		initialUrl = strings.Join([]string{initialUrl, containDateQuery}, "")
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", initialUrl, nil)

	logError(err)

	req.Header.Add("cookie", "npic=wn2ZumCp+gAZ1ZduZFFrpAoucn7Cbf+Ihhd2issZhiNvuPRGzk/3fHzMCVskjJbqCA==; NNB=HTWUMJ4UFKTVY; MM_NEW=1; NFS=2; nx_open_so=1; NRTK=ag#all_gr#1_ma#-2_si#0_en#0_sp#0; _ga=GA1.1.1148334420.1576898593; _ga_7VKFYR6RV1=GS1.1.1579407327.1.1.1579407330.0; ASID=d3269c69000001708ff282340000005b; nx_ssl=2; nid_inf=-1160397537; NID_AUT=x30e3uTQa7m81jVkS1tDfDzVYp5ya1nNglGzMrem2nWriL0VdTI/+TGJfyzlfnz0; NID_JKL=h9sJuzScUEI4Z1QiscTvLdOWr1xPcKRNLITwxiaN7Mw=; NID_SES=AAABjfTypiVl3erpwe894jMOArNRjScl/ll1WziISA/56X8U1OMPCl7UaWQoQseADaMu02u7lFUn2DGdsoO1iwXzKFgrWWF6O8/uqZcZTlVcehM5PZkVzpCg1aFD056XH1PLnRR6PiGWj2sbeVJxfRy7ztKq5y0x35C5BSklUXfxlbWvVjqsH/W32zMW5eq6etaR9e+FSpJ0/3pQVR0YGTKm9b/YTBq2+0t9CzYX7+hQ1POOHMDRpP4NTfxVEBUFaopnAT/RPGUQL1xaExjwXY7yoLGkvq9yJ5/KuRPTFHKp4bCdfIyVxN38da2PtlTZHDQoxGG5eBSz70JIWCOrCpF3PPMF/9aXiCMJQpvzdhI2iuHRNS93E6KhwLutgq3DAZxqDNMyyyYtCgMI/E23er9i0oc6c0xkwRQDRQNI4HY7P6pWjYFm1AIqm47kMTDUmkuzdV4P/inRP3ZD8vM/ihuSoOS8/2crk5WphPxL4tbveMXw57rl1hpIbZ06XVkrTxwYN4nfPc0fv+VpQvvrmfF8B60=; _naver_usersession_=W6Uce5E3G7uKGoBHdSXAAQ==; page_uid=UqSrJsprvxsssd6knzGssssssvl-455711; nx_mson=1")

	res, err := client.Do(req)

	logError(err)

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	logError(err)

	var nextLinkBox []string
	var linkBox []string

	fmt.Println(baseURL)

	// lastAnchorText := doc.Find(".paging a").Last().Text()

	idx := 0

	for {

		if doc.Find(".paging a").Size() == 0 {
			fmt.Println("검색결과가 없습니다")
			break
		}

		if idx == 0 {
			nextPage := doc.Find(".paging a").Last()
			if nextPage.Text() == "다음페이지" {
				link, _ := nextPage.Attr("href")
				nextLinkBox = append(nextLinkBox, link)
			}
		}

		sub_client := &http.Client{}
		urlArr := []string{baseURL, nextLinkBox[idx]}
		next_url := strings.Join(urlArr[:], "")

		_req, _err := http.NewRequest("GET", next_url, nil)

		logError(err)

		_req.Header.Add("cookie", "npic=wn2ZumCp+gAZ1ZduZFFrpAoucn7Cbf+Ihhd2issZhiNvuPRGzk/3fHzMCVskjJbqCA==; NNB=HTWUMJ4UFKTVY; MM_NEW=1; NFS=2; nx_open_so=1; NRTK=ag#all_gr#1_ma#-2_si#0_en#0_sp#0; _ga=GA1.1.1148334420.1576898593; _ga_7VKFYR6RV1=GS1.1.1579407327.1.1.1579407330.0; ASID=d3269c69000001708ff282340000005b; nx_ssl=2; nid_inf=-1160397537; NID_AUT=x30e3uTQa7m81jVkS1tDfDzVYp5ya1nNglGzMrem2nWriL0VdTI/+TGJfyzlfnz0; NID_JKL=h9sJuzScUEI4Z1QiscTvLdOWr1xPcKRNLITwxiaN7Mw=; NID_SES=AAABjfTypiVl3erpwe894jMOArNRjScl/ll1WziISA/56X8U1OMPCl7UaWQoQseADaMu02u7lFUn2DGdsoO1iwXzKFgrWWF6O8/uqZcZTlVcehM5PZkVzpCg1aFD056XH1PLnRR6PiGWj2sbeVJxfRy7ztKq5y0x35C5BSklUXfxlbWvVjqsH/W32zMW5eq6etaR9e+FSpJ0/3pQVR0YGTKm9b/YTBq2+0t9CzYX7+hQ1POOHMDRpP4NTfxVEBUFaopnAT/RPGUQL1xaExjwXY7yoLGkvq9yJ5/KuRPTFHKp4bCdfIyVxN38da2PtlTZHDQoxGG5eBSz70JIWCOrCpF3PPMF/9aXiCMJQpvzdhI2iuHRNS93E6KhwLutgq3DAZxqDNMyyyYtCgMI/E23er9i0oc6c0xkwRQDRQNI4HY7P6pWjYFm1AIqm47kMTDUmkuzdV4P/inRP3ZD8vM/ihuSoOS8/2crk5WphPxL4tbveMXw57rl1hpIbZ06XVkrTxwYN4nfPc0fv+VpQvvrmfF8B60=; _naver_usersession_=W6Uce5E3G7uKGoBHdSXAAQ==; page_uid=UqSrJsprvxsssd6knzGssssssvl-455711; nx_mson=1")

		_res, _err := sub_client.Do(_req)

		logError(_err)

		defer _res.Body.Close()

		_doc, _err := goquery.NewDocumentFromReader(_res.Body)

		logError(_err)

		if _doc.Find(".paging a").Last().Text() != "다음페이지" {
			break
		}

		link, _ := _doc.Find(".paging a").Last().Attr("href")

		nextLinkBox = append(nextLinkBox, link)
		idx += 1
	}

	// goRoutine => 위에서 가져온 페이지 중 각 페이지 내에 기사 링크를 가져옴.
	var nextLinkWg sync.WaitGroup
	// linkChan := make(chan string)
	for _, link := range nextLinkBox {
		nextLinkWg.Add(1)
		go getAllLinkInPage(&linkBox, link, &nextLinkWg)
	}

	nextLinkWg.Wait()

	// 위의 goroutine이 끝나면 이것도 goroutine 으로 기사 내용을 긁어온다.
	var wg sync.WaitGroup

	for _, link := range linkBox {
		wg.Add(1)
		go getArticleContent(link, &wg)
	}

	wg.Wait()
}

func getAllLinkInPage(linkBox *[]string, link string, nextLinkWg *sync.WaitGroup) {

	defer nextLinkWg.Done()

	pageURL := baseURL + link

	pageClient := &http.Client{}

	pageReq, pageErr := http.NewRequest("GET", pageURL, nil)

	logError(pageErr)

	pageRes, pageResErr := pageClient.Do(pageReq)

	logError(pageResErr)

	defer pageRes.Body.Close()

	pageDoc, pageDocErr := goquery.NewDocumentFromReader(pageRes.Body)

	logError(pageDocErr)

	pageDoc.Find(".type01 .txt_inline a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		if link != "#" {
			*linkBox = append(*linkBox, link)
			// linkChan <- link
		}
	})

	time.Sleep(time.Second)
}

func getArticleContent(link string, wg *sync.WaitGroup) {

	defer wg.Done()

	contentClient := &http.Client{}

	contentReq, contentReqErr := http.NewRequest("GET", link, nil)

	logError(contentReqErr)

	contentRes, contentResErr := contentClient.Do(contentReq)

	logError(contentResErr)

	defer contentRes.Body.Close()

	contentDoc, contentDocErr := goquery.NewDocumentFromReader(contentRes.Body)

	fmt.Println(contentDoc)

	logError(contentDocErr)

	validFileRegEx := "[:\\\\/%*?:|\"<>]"

	validFileReg, _ := regexp.Compile(validFileRegEx)

	// fmt.Println(contentDoc.Find("#articleTitle").Text())
	title := contentDoc.Find("#articleTitle").Text()
	pressDate := contentDoc.Find(".t11").First().Text()
	publisher, _ := contentDoc.Find(".press_logo img").Attr("title")
	content := contentDoc.Find("#articleBodyContents").Text()
	utf8Title, _ := decodeToKOR(title)
	utf8Title = validFileReg.ReplaceAllString(utf8Title, "")
	utf8PressDate, _ := decodeToKOR(pressDate)
	utf8Pulisher, _ := decodeToKOR(publisher)
	utf8Content, _ := decodeToKOR(content)

	_currentDate := time.Now()

	currentDate := _currentDate.Format("2006-01-02")

	_, err := os.Stat("C:\\crawling_result\\" + currentDate)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll("C:\\crawling_result\\"+currentDate, os.ModeDir)
		logError(errDir)
	}

	_err := ioutil.WriteFile("C:/crawling_result/"+currentDate+"/"+utf8Title+"_"+currentDate+".txt", []byte(utf8Title+"\n"+link+"\n"+utf8PressDate+"\n"+utf8Pulisher+"\n"+utf8Content), os.FileMode(644))

	logError(_err)
}

func decodeToKOR(text string) (string, error) {
	dst := make([]byte, len(text)*2)
	tr := korean.EUCKR.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)

	if err != nil {
		return text, err
	}

	return string(dst[:nDst]), nil
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
