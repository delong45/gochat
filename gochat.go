package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"encoding/json"
	"bytes"
	"os"
	"strings"
	"strconv"
)

const (
	VERSION     = "0.0.1"
	WXURL       = "http://in.qyapi.weixin.qq.com/cgi-bin/"
	WXINURL     = "http://in.qyapi.weixin.qq.com/cgi-bin/tencent/"
	CORPID      = ""
	SECRET      = ""
	CONTENT     = "This is gochat robot, just for testing, ignore it.."
	CHATID      = "wxzc"
)

type CreateChat struct {
	Name	 string    `json:"name"`
	UserList []string  `json:"userlist"`
}

type UpdateChat struct {
	Chatid	     string    `json:"chatid"`
	Name         string    `json:"name"`
	AddUserList  []string  `json:"add_user_list"`
	DelUserList  []string  `json:"del_user_list"`
}

type TextChat struct {
	Receiver struct {
		Type	string  `json:"type"`
		Id	    string  `json:"id"`
	} `json:"receiver"`
	Msgtype		string  `json:"msgtype"`
	Text struct {
		Content string  `json:"content"`
	} `json:"text"`
}

type NameList struct {
	NameList    []string `json:"name_list"`
}

type ReturnVal struct {
	Errcode		int		`json:"errcode"`
	Errmsg	    string  `json:"errmsg"`
}

type ResultToken struct {
	Errcode     int     `json:"errcode"`
	Errmsg      string  `json:"errmsg"`
	AccessToken string  `json:"access_token"`
	ExpiresIn   int     `json:"expires_in"`
}

type ResultCreateChat struct {
	Errcode     int     `json:"errcode"`
	Errmsg      string  `json:"errmsg"`
	Chatid      string  `json:"chatid"`
}

type UserInfo struct {
	Userid      string  `json:"userid"`
	Name        string  `json:"name"`
}

type ResultUserid struct {
	Errcode     int        `json:"errcode"`
	Errmsg      string     `json:"errmsg"`
	UserList    []UserInfo `json:"user_list"`
}

type StaffInfo struct {
	Name		string
	Phone		string
}

type Config struct {
	Chatid              string
	StaffList           []StaffInfo
	NoticeDuty          string
	NoticePerson        string
	NoticeDailyReport   string
	NoticeWeekReport    string
}

var Conf Config

func readContent(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func writeContent(file string, content string) {
	data := []byte(content)
	err := ioutil.WriteFile(file, data, 0644)
	if err != nil {
		panic(err)
	}
}

func readConfig(conf string) {
	file, _ := os.Open(conf)
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.Decode(&Conf)
}

func getToken(corpid string, secret string) (string, error) {
	baseurl := WXURL + "gettoken"
	u, _ := url.Parse(baseurl)
	q := u.Query()
	q.Set("corpid", corpid)
	q.Set("corpsecret", secret)
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String());
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var val ResultToken
	json.Unmarshal(body, &val)
	if val.Errcode > 0 {
		return "", fmt.Errorf(val.Errmsg)
	}
	return val.AccessToken, nil
}

func convertToUserid(token string, nameList []string) ([]UserInfo, error){
	baseurl := WXINURL + "user/convert_to_userid"
	u, _ := url.Parse(baseurl)
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()

	names := NameList {
		NameList: nameList }
	b, _ := json.Marshal(names)
	log.Println(string(b))

	body := bytes.NewBuffer([]byte(b))
	res, err := http.Post(u.String(), "application/json;charset=utf-8", body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	log.Println(string(result))

	var val ResultUserid
	json.Unmarshal(result, &val)
	if val.Errcode > 0 {
		return nil, fmt.Errorf(val.Errmsg)
	}
	return val.UserList, nil
}

func createChat(token string, name string, userList []string) (string, error) {
	baseurl := WXINURL + "chat/create"
	u, _ := url.Parse(baseurl)
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()

	chat := CreateChat {
		Name:     name,
		UserList: userList }
	b, _ := json.Marshal(chat)
	log.Println(string(b))

	body := bytes.NewBuffer([]byte(b))
	res, err := http.Post(u.String(), "application/json;charset=utf-8", body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	log.Println(string(result))

	var val ResultCreateChat
	json.Unmarshal(result, &val)
	if val.Errcode > 0 {
		return "", fmt.Errorf(val.Errmsg)
	}

	return val.Chatid, nil
}

func updateChat(token string, id string, name string, addUserList []string, delUserList []string) error {
	baseurl := WXINURL + "chat/update"
	u, _ := url.Parse(baseurl)
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()

	chat := UpdateChat {
		Chatid:      id,
		Name:        name,
		AddUserList: addUserList,
		DelUserList: delUserList }
	b, _ := json.Marshal(chat)
	log.Println(string(b))

	body := bytes.NewBuffer([]byte(b))
	res, err := http.Post(u.String(), "application/json;charset=utf-8", body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Println(string(result))

	var val ReturnVal
	json.Unmarshal(result, &val)
	if val.Errcode > 0 {
		return fmt.Errorf(val.Errmsg)
	}

	return nil
}

func sendMsg(token string, chatType string, id string, content string) error {
	baseurl := WXINURL + "chat/send"
	u, _ := url.Parse(baseurl)
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()

	var chat TextChat
	chat.Receiver.Type = chatType
	chat.Receiver.Id = id
	chat.Msgtype = "text"
	chat.Text.Content = content
	b, _ := json.Marshal(chat)
	log.Println(string(b))

	body := bytes.NewBuffer([]byte(b))
	res, err := http.Post(u.String(), "application/json;charset=utf-8", body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Println(string(result))

	var val ReturnVal
	json.Unmarshal(result, &val)
	if val.Errcode > 0 {
		return fmt.Errorf(val.Errmsg)
	}

	return nil
}

func createChatTX(name string, userList string) {
	token, err := getToken(CORPID, SECRET)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(token)

	users, err := convertToUserid(token, strings.Split(userList, ","))
	if err != nil {
		log.Fatal(err)
	}
	var list []string
	for _, user := range users {
		list = append(list, user.Userid)
	}

	chatId, err := createChat(token, name, list)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ChatId:" + chatId)
	return
}

func updateChatTX(id string, name string, addList string, delList string) {
	token, err := getToken(CORPID, SECRET)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(token)

	addUsers, err := convertToUserid(token, strings.Split(addList, ","))
	if err != nil {
		log.Fatal(err)
	}
	var addlist []string
	for _, addUser := range addUsers {
		addlist = append(addlist, addUser.Userid)
	}

	delUsers, err := convertToUserid(token, strings.Split(delList, ","))
	if err != nil {
		log.Fatal(err)
	}
	var dellist []string
	for _, delUser := range delUsers {
		dellist = append(dellist, delUser.Userid)
	}

	err = updateChat(token, id, name, addlist, dellist)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func sendMsgTX(content string, chatType string, id string) {
	token, err := getToken(CORPID, SECRET)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Token:" + token)

	err = sendMsg(token, chatType, id, content)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func getUserNo() int {
	num := readContent("./record.log")
	num = strings.TrimSpace(num)
	n, _ := strconv.Atoi(num)
	return n
}

func setUserNo(n int) {
	num := strconv.Itoa(n)
	writeContent("./record.log", num)
}

func getUserid(isLast bool) string {
	token, err := getToken(CORPID, SECRET)
	if err != nil {
		log.Fatal(err)
	}

	n := getUserNo()
	index := n % len(Conf.StaffList)
	name := Conf.StaffList[index].Name

	var list []string
	list = append(list, name)
	users, err := convertToUserid(token, list)
	if err != nil {
		log.Fatal(err)
	}
	id := users[0].Userid
	log.Println("Get staff: " + users[0].Userid)
	log.Println("Get staff: " + users[0].Name)
	log.Println("Get staff: " + Conf.StaffList[index].Phone)

	if isLast {
		n++
		setUserNo(n)
	}

	return id
}

func getStaffInfo() (string, string) {
	n := getUserNo()
	index := n % len(Conf.StaffList)
	name := Conf.StaffList[index].Name
	phone := Conf.StaffList[index].Phone
	info := name + ", " + phone
	cid := Conf.Chatid

	return info, cid
}

func getSessionInfo(category string) (string, string, string) {
	if category == "" {
		log.Fatal("Category of message to send is empty")
	}

	var content string
	var ctype string
	var cid string
	var info string

	if category == "duty" {
		content = Conf.NoticeDuty
		ctype = "group"
		info, cid = getStaffInfo()
		content += info
	} else if category == "person" {
		content = Conf.NoticePerson
		ctype = "single"
		cid = getUserid(false)
	} else if category == "daily" {
		content = Conf.NoticeDailyReport
		ctype = "single"
		cid = getUserid(false)
	} else if category == "week" {
		content = Conf.NoticeWeekReport
		ctype = "single"
		cid = getUserid(true)
	} else {
		log.Fatal("Wrong category type")
	}

	return content, ctype, cid
}

func main() {
	ifCreate := flag.Bool("c", false, "Create Chat")
	ifUpdate := flag.Bool("u", false, "Update Chat")
	ifSend   := flag.Bool("s", false, "Send Message")

	// 1. duty 2. person 3. daily 4. week
	category := flag.String("category", "", "Category of message to send")

	id       := flag.String("id", "", "Chat ID")
	name     := flag.String("name", "", "Chat Name")
	chatType := flag.String("type", "", "Chat Type")
	userList := flag.String("userlist", "", "User List")
	addList  := flag.String("addlist", "", "Add User List")
	delList  := flag.String("dellist", "", "Del User List")

	flag.Parse()

	readConfig("./conf.json")

	if *ifCreate {
		createChatTX(*name, *userList)
		os.Exit(0)
	}

	if *ifUpdate {
		updateChatTX(*id, *name, *addList, *delList)
		os.Exit(0)
	}

	if *ifSend {
		content := CONTENT
		ctype   := *chatType
		cid     := *id
		if cid == "" || ctype == "" {
			content, ctype, cid = getSessionInfo(*category)
		}
		sendMsgTX(content, ctype, cid)
		os.Exit(0)
	}

	flag.PrintDefaults()
}
