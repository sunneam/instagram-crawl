package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
	"whymeins.go/Random"
	"whymeins.go/model"
)


var (
	wFile *bufio.Writer

	f *os.File

	dir, _ =os.Getwd()

	//首页第一次分页地址
	murl ="https://www.veryins.com/t/?next=0&tag=1"

	//分页地址prefix
	purl="https://www.veryins.com/t/?"

	//个人主页分页prefix
	surl="https://www.veryins.com/user/post?"

	//爬取主页面总页数
	num=3
)


func main()  {
	//爬取接口
	http.HandleFunc("/veryins",StartVeryin)

	//拉起8082端口
	err :=http.ListenAndServe(":8082",nil)

	if err!=nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func StartVeryin(w http.ResponseWriter,req *http.Request)  {

	//判断文件是否存在
	FileCheck()

	wFile.WriteString("程序开始运行:"+time.Now().Format("2006-01-02 15:04:05")+"\n")

	defer f.Close()

	CNum:=0
	//初始化数据库
	model.Init()

	//主页收集器
	c:=colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
		colly.AllowURLRevisit(),
	)

	//禁用长连接
	c.WithTransport(&http.Transport{
		DisableKeepAlives:true,
	})

	//主页分页开启两条线程
	q1,_:=queue.New(
		2,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	//个人主页开启5条线程
	q2,_:=queue.New(
		5,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	//个人主页12条收集器
	d:=c.Clone()

	//主页爬取
	GetDataOnIndex(c,d,&CNum)

	//个人主页爬取
	GetDataInPerson(d)

	c.Visit("https://www.veryins.com/")

	//异步等待
	c.Wait();d.Wait()

	//启动线程
	q1.Run(c);q2.Run(d)

	wFile.WriteString(fmt.Sprintf("程序运行结束:"+time.Now().Format("2006-01-02 15:04:05"+"\n")))

	//更新缓存
	wFile.Flush()

	CNum=0

	f.Close()

	fmt.Println("爬取数据结束")
}

//获取主页前12条数据并访问个人主页
func GetDataOnIndex(c *colly.Collector,d *colly.Collector,Cnum *int)  {
	c.OnHTML(".item", func(el *colly.HTMLElement) {
		//主页面爬取完成
		mainpage:=model.MainPage{}
		mainpage.Code="https://www.veryins.com/p/"+el.ChildAttr(".img-wrap","data-code")
		mainpage.ThumbnailSrc=el.ChildAttr(".img-wrap","data-src")
		mainpage.Caption=el.ChildAttr(".img-wrap img","alt")
		mainpage.Username="https://www.veryins.com"+el.ChildAttr(".item-body a","href")
		mainpage.Data=el.ChildText(".item-body .likes span")

		//如果username相同
		user:=model.MainPage{}
		model.Db.C("bill").Find(map[string]string{"username":mainpage.Username}).One(&user)

		if len(user.Username)<1{

			//该用户未爬取过则爬取该用户所有作品
			model.Db.C("bill").Insert(mainpage)
			d.Visit(mainpage.Username)

		}else{

			//该用户已经爬取过,但该数据是新增的
			u:=model.MainPage{}
			model.Db.C("bill").Find(map[string]string{"username":mainpage.Username,"code":mainpage.Code}).One(&u)

			if len(u.Code)<1 {
				//插入父表
				model.Db.C("bill").Insert(mainpage)

				//插入子表
				p:=model.MainPage{
					Caption:mainpage.Caption,
					Code:mainpage.Code,
					ProfilePicUrl:mainpage.ProfilePicUrl,
					MediaType:mainpage.MediaType,
					ThumbnailSrc:mainpage.ThumbnailSrc,
					From:mainpage.Username,
				}
				model.Db.C("person").Insert(p)
			}

		}


	})

	//主页爬取完12条开始接口分页爬取
	c.OnResponse(func(r *colly.Response) {
		if *Cnum<1 {
			//第一次访问api接口
			c.Post(murl, map[string]string{})
			*Cnum++
		}else{
				//主页数据解析数据
				m:=model.MainPage{}
				j:=model.R{}
				json.Unmarshal(r.Body,&j)
				next:=j.Next
				for _,v:=range j.Media{
					m.Code=v.Code
					m.Caption=v.Caption
					m.Data=v.Date
					m.MediaType=v.MediaType
					m.Username="https://www.veryins.com/"+v.Owner.Username
					m.ProfilePicUrl=v.Owner.ProfilePicUrl
					m.FullName=v.Owner.FullName
					m.ThumbnailSrc=v.ThumbnailSrc

					//如果username相同
					user:=model.MainPage{}
					model.Db.C("bill").Find(map[string]string{"username":m.Username}).One(&user)

					if len(user.Username)<1{
						//该用户未爬取过则爬取该用户所有作品
						model.Db.C("bill").Insert(m)
						d.Visit(m.Username)

					}else{

						//该用户已经爬取过,但该数据是新增的
						u:=model.MainPage{}
						model.Db.C("bill").Find(map[string]string{"username":m.Username,"code":m.Code}).One(&u)

						if len(u.Code)<1 {
							//插入父表
							model.Db.C("bill").Insert(m)

							//插入子表
							p:=model.MainPage{
								Caption:m.Caption,
								Code:m.Code,
								ProfilePicUrl:m.ProfilePicUrl,
								MediaType:m.MediaType,
								ThumbnailSrc:m.ThumbnailSrc,
								From:m.Username,
							}

							model.Db.C("person").Insert(p)
						}
					}

				}

				//开始访问下一页
				nexturl:=fmt.Sprintf(purl+"next=%v&tag=1",next)

				if next<float64(num){
					c.Post(nexturl, map[string]string{})
				}
		}

	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent",Random.RandomString())
	})

}

//获取个人主页前12条数据
func GetDataInPerson(d *colly.Collector){

	d.OnHTML(".item", func(e *colly.HTMLElement) {
			mainpage:=model.MainPage{}
			mainpage.Code="https://www.veryins.com/p/"+e.ChildAttr(".img-wrap","data-code")
			mainpage.ThumbnailSrc=e.ChildAttr(".img-wrap","data-src")
			mainpage.Caption=e.ChildAttr(".img-wrap img","alt")
			mainpage.Username="https://www.veryins.com"+e.ChildAttr(".item-body a","href")
			mainpage.From=e.Request.URL.String()
			mainpage.Data=e.ChildText(".item-body .likes span")
			model.Db.C("person").Insert(mainpage)
	})

	//前12个人主页以后的api第一次调用
	d.OnResponse(func(r *colly.Response) {
		//解析地址,如果地址包含from则是个人主页api调用,否则就是主页触发的前个人主页前12个
		v,_:=url.ParseQuery(r.Request.URL.String())
		from:=v.Get("from")
		if len(from)==0 {
			reg:=regexp.MustCompile(`next-cursor=.*==`)
			StrNext:=reg.FindString(string(r.Body))
			next:=strings.Replace(StrNext,"next-cursor=\"","",-1)
			reg=regexp.MustCompile(`<div id="username".*data-fullname`)
			aa:=reg.FindString(string(r.Body))
			reg=regexp.MustCompile(`[0-9a-fA-F]{32}`)
			uid:=reg.FindString(aa)
			//类型为1
			tag:=1
			show:=fmt.Sprintf(surl+"next=%v&uid=%v&tag=%v&from=%v",next,uid,tag,r.Request.URL.String())
			d.Post(show,map[string]string{})
		}else{
			//取出地址中的uid
			v,_:=url.ParseQuery(r.Request.URL.String())
			uid:=v.Get("uid")
			//获取参数
			tone:=model.Tone{}
			json.Unmarshal(r.Body,&tone)
			//取出下一页
			has:=tone.PageInfo.HasNextPage
			next:=tone.PageInfo.EndCursor
			//取出地址来源
			from:=v.Get("from")
			//类型
			tag:=1
			//将数据append到son数组
			for _,v:=range tone.Nodes{
				p:=model.MainPage{
					Caption:v.Caption,
					Code:v.Code,
					ProfilePicUrl:v.DisplaySrc,
					MediaType:v.MediaType,
					ThumbnailSrc:v.ThumbnailSrc,
					From:from,
				}
				model.Db.C("person").Insert(p)
			}

			show:=fmt.Sprintf(surl+"next=%v&uid=%v&tag=%v&from=%v",next,uid,tag,from)
			if has==true {
				d.Post(show, map[string]string{})
			}
		}

	})

	d.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent",Random.RandomString())
	})
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//封装文件

func FileCheck()  {
	var err error
	//初始化日志文件路径
	Logfile:=dir+"/logs/"+time.Now().Format("2006-01-02")+".log"

	//检查文件是否存在
	if checkFileIsExist(Logfile) {
		f, _= os.OpenFile(Logfile, os.O_APPEND, 0666)
	} else {
		fmt.Println("正在创建文件",Logfile)
		f, err= os.Create(Logfile)
		if err!=nil {
			fmt.Println("创建文件错误,",err)
		}
	}
	//封装buffer
	wFile =bufio.NewWriter(f)
}



