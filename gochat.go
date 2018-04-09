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
)

const (
	VERSION		= "0.0.1"
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
	Chatid	     string	   `json:"chatid"`
	Name         string    `json:"name"`
	AddUserList  []string  `json:"add_user_list"`
	DelUserList  []string  `json:"del_user_list"`
}

type TextChat struct {
	Receiver struct {
		Type	string  `json:"type"`
		Id	    string  `json:"id"`
	} `json:"receiver"`
	Msgtype		string `json:"msgtype"`
	Text struct {
		Content string `json:"content"`
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
	fmt.Println(string(b))

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
	fmt.Println(string(result))

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
	fmt.Println(string(b))

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
	fmt.Println(string(result))

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
	fmt.Println(string(b))

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
	fmt.Println(string(result))

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
	fmt.Println(string(b))

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
	fmt.Println(string(result))

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
	fmt.Println(token)

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
	fmt.Println(token)

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

func sendMsgTX(chatType string, id string) {
	token, err := getToken(CORPID, SECRET)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

	err = sendMsg(token, chatType, id, CONTENT)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	ifCreate := flag.Bool("c", false, "Create Chat")
	ifUpdate := flag.Bool("u", false, "Update Chat")
	ifSend   := flag.Bool("s", false, "Send Message")

	id       := flag.String("id", "", "Chat ID")
	name     := flag.String("name", "", "Chat Name")
	chatType := flag.String("type", "group", "Chat Type")
	userList := flag.String("userlist", "", "User List")
	addList  := flag.String("addlist", "", "Add User List")
	delList  := flag.String("dellist", "", "Del User List")

	flag.Parse()

	if *ifCreate {
		createChatTX(*name, *userList)
		os.Exit(0)
	}

	if *ifUpdate {
		updateChatTX(*id, *name, *addList, *delList)
		os.Exit(0)
	}

	if *ifSend {
		sendMsgTX(*chatType, *id)
		os.Exit(0)
	}

	flag.PrintDefaults()
}
