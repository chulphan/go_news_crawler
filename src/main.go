package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/korean"
)

// func getPages(url) {

// }

func main() {
	// url := "https://search.naver.com/search.naver?where=news&sm=tab_jum&query=%EC%82%BC%EC%84%B1%EC%A0%84%EC%9E%90+IR"
	baseURL := "https://search.naver.com/search.naver"
	url := "https://search.naver.com/search.naver?where=news&query=%EC%82%BC%EC%84%B1%EC%A0%84%EC%9E%90&sm=tab_opt&sort=0&photo=0&field=0&reporter_article=&pd=3&ds=2020.04.25&de=2020.04.26&docid=&nso=so%3Ar%2Cp%3Afrom20200425to20200426%2Ca%3Aall&mynews=0&refresh_start=0&related=0"

	// res, err := http.Get(url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("cookie", "npic=wn2ZumCp+gAZ1ZduZFFrpAoucn7Cbf+Ihhd2issZhiNvuPRGzk/3fHzMCVskjJbqCA==; NNB=HTWUMJ4UFKTVY; MM_NEW=1; NFS=2; nx_open_so=1; NRTK=ag#all_gr#1_ma#-2_si#0_en#0_sp#0; _ga=GA1.1.1148334420.1576898593; _ga_7VKFYR6RV1=GS1.1.1579407327.1.1.1579407330.0; ASID=d3269c69000001708ff282340000005b; nx_ssl=2; nid_inf=-1160397537; NID_AUT=x30e3uTQa7m81jVkS1tDfDzVYp5ya1nNglGzMrem2nWriL0VdTI/+TGJfyzlfnz0; NID_JKL=h9sJuzScUEI4Z1QiscTvLdOWr1xPcKRNLITwxiaN7Mw=; NID_SES=AAABjfTypiVl3erpwe894jMOArNRjScl/ll1WziISA/56X8U1OMPCl7UaWQoQseADaMu02u7lFUn2DGdsoO1iwXzKFgrWWF6O8/uqZcZTlVcehM5PZkVzpCg1aFD056XH1PLnRR6PiGWj2sbeVJxfRy7ztKq5y0x35C5BSklUXfxlbWvVjqsH/W32zMW5eq6etaR9e+FSpJ0/3pQVR0YGTKm9b/YTBq2+0t9CzYX7+hQ1POOHMDRpP4NTfxVEBUFaopnAT/RPGUQL1xaExjwXY7yoLGkvq9yJ5/KuRPTFHKp4bCdfIyVxN38da2PtlTZHDQoxGG5eBSz70JIWCOrCpF3PPMF/9aXiCMJQpvzdhI2iuHRNS93E6KhwLutgq3DAZxqDNMyyyYtCgMI/E23er9i0oc6c0xkwRQDRQNI4HY7P6pWjYFm1AIqm47kMTDUmkuzdV4P/inRP3ZD8vM/ihuSoOS8/2crk5WphPxL4tbveMXw57rl1hpIbZ06XVkrTxwYN4nfPc0fv+VpQvvrmfF8B60=; _naver_usersession_=W6Uce5E3G7uKGoBHdSXAAQ==; page_uid=UqSrJsprvxsssd6knzGssssssvl-455711; nx_mson=1")

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err)
	}

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

		if _err != nil {
			log.Fatal("error in for loop")
		}

		_req.Header.Add("cookie", "npic=wn2ZumCp+gAZ1ZduZFFrpAoucn7Cbf+Ihhd2issZhiNvuPRGzk/3fHzMCVskjJbqCA==; NNB=HTWUMJ4UFKTVY; MM_NEW=1; NFS=2; nx_open_so=1; NRTK=ag#all_gr#1_ma#-2_si#0_en#0_sp#0; _ga=GA1.1.1148334420.1576898593; _ga_7VKFYR6RV1=GS1.1.1579407327.1.1.1579407330.0; ASID=d3269c69000001708ff282340000005b; nx_ssl=2; nid_inf=-1160397537; NID_AUT=x30e3uTQa7m81jVkS1tDfDzVYp5ya1nNglGzMrem2nWriL0VdTI/+TGJfyzlfnz0; NID_JKL=h9sJuzScUEI4Z1QiscTvLdOWr1xPcKRNLITwxiaN7Mw=; NID_SES=AAABjfTypiVl3erpwe894jMOArNRjScl/ll1WziISA/56X8U1OMPCl7UaWQoQseADaMu02u7lFUn2DGdsoO1iwXzKFgrWWF6O8/uqZcZTlVcehM5PZkVzpCg1aFD056XH1PLnRR6PiGWj2sbeVJxfRy7ztKq5y0x35C5BSklUXfxlbWvVjqsH/W32zMW5eq6etaR9e+FSpJ0/3pQVR0YGTKm9b/YTBq2+0t9CzYX7+hQ1POOHMDRpP4NTfxVEBUFaopnAT/RPGUQL1xaExjwXY7yoLGkvq9yJ5/KuRPTFHKp4bCdfIyVxN38da2PtlTZHDQoxGG5eBSz70JIWCOrCpF3PPMF/9aXiCMJQpvzdhI2iuHRNS93E6KhwLutgq3DAZxqDNMyyyYtCgMI/E23er9i0oc6c0xkwRQDRQNI4HY7P6pWjYFm1AIqm47kMTDUmkuzdV4P/inRP3ZD8vM/ihuSoOS8/2crk5WphPxL4tbveMXw57rl1hpIbZ06XVkrTxwYN4nfPc0fv+VpQvvrmfF8B60=; _naver_usersession_=W6Uce5E3G7uKGoBHdSXAAQ==; page_uid=UqSrJsprvxsssd6knzGssssssvl-455711; nx_mson=1")

		_res, _err := sub_client.Do(_req)

		if _err != nil {
			log.Fatal("error in for loop after sub_client")
		}

		defer _res.Body.Close()

		_doc, _err := goquery.NewDocumentFromReader(_res.Body)

		if _err != nil {
			log.Fatal("error in for loop after goquery")
		}

		if _doc.Find(".paging a").Last().Text() != "다음페이지" {
			break
		}

		link, _ := _doc.Find(".paging a").Last().Attr("href")

		nextLinkBox = append(nextLinkBox, link)
		idx += 1
	}

	fmt.Println("nextLink ", nextLinkBox)

	for _, v := range nextLinkBox {
		pageURL := baseURL + v

		pageClient := &http.Client{}

		pageReq, pageErr := http.NewRequest("GET", pageURL, nil)

		if pageErr != nil {
			log.Fatal(pageErr)
		}

		pageRes, pageResErr := pageClient.Do(pageReq)

		if pageResErr != nil {
			log.Fatal(pageResErr)
		}

		defer pageRes.Body.Close()

		pageDoc, pageDocErr := goquery.NewDocumentFromReader(pageRes.Body)

		if pageDocErr != nil {
			log.Fatal(pageDocErr)
		}

		pageDoc.Find(".type01 .txt_inline a").Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr("href")
			if link != "#" {
				linkBox = append(linkBox, link)
			}
		})

	}

	for _, link := range linkBox {
		contentClient := &http.Client{}

		contentReq, contentReqErr := http.NewRequest("GET", link, nil)

		if contentReqErr != nil {
			log.Fatal(contentReq)
		}

		contentRes, contentResErr := contentClient.Do(contentReq)

		if contentResErr != nil {
			log.Fatal(contentResErr)
		}

		defer contentRes.Body.Close()

		// utfContentBody, utfContentBodyErr := iconv.NewReader(contentRes.Body, "euc-kr", "utf-8")

		// if utfContentBodyErr != nil {
		// 	log.Fatal(utfContentBodyErr)
		// }

		contentDoc, contentDocErr := goquery.NewDocumentFromReader(contentRes.Body)

		fmt.Println(contentDoc)

		if contentDocErr != nil {
			log.Fatal(contentDocErr)
		}

		// fmt.Println(contentDoc.Find("#articleTitle").Text())
		title := contentDoc.Find("#articleTitle").Text()
		pressDate := contentDoc.Find(".t11").First().Text()
		publisher, _ := contentDoc.Find(".press_logo img").Attr("title")
		content := contentDoc.Find("#articleBodyContents").Text()
		utf8Title, _ := decodeToKOR(title)
		utf8PressDate, _ := decodeToKOR(pressDate)
		utf8Pulisher, _ := decodeToKOR(publisher)
		utf8Content, _ := decodeToKOR(content)
		fmt.Println(utf8Title)
		fmt.Println(utf8PressDate)
		fmt.Println(utf8Pulisher)
		fmt.Println(utf8Content)

		_, err := os.Stat("C:\\crawling_result\\2020-05-02")

		if os.IsNotExist(err) {
			errDir := os.MkdirAll("C:\\crawling_result\\2020-05-02", os.ModeDir)
			if errDir != nil {
				log.Fatal(err)
			}
		}

		file, createFileErr := os.Create("C:/crawling_result/2020-05-02/" + utf8Title + ".txt")
		if createFileErr != nil {
			log.Fatal(createFileErr)
		}

		defer file.Close()

		_err := ioutil.WriteFile("C:/crawling_result/2020-05-02/"+utf8Title+".txt", []byte(utf8Title+"\n"+utf8PressDate+"\n"+utf8Pulisher+"\n"+utf8Content), os.FileMode(644))

		if _err != nil {
			log.Fatal(err)
		}

		// _, fileContentErr := file.WriteString(utf8Title + "\n" + utf8PressDate + "\n" + utf8Content)

		// if fileContentErr != nil {
		// 	log.Fatal(fileContentErr)
		// }

	}

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
