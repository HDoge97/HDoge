package controller

import (
	"creditbank/sdkInit"
	"creditbank/service"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var cuser User

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	ShowView(w, r, "login.html", nil)
}

func (app *Application) Qublock(w http.ResponseWriter, r *http.Request) {
	const (
		cc_name    = "simplecc"
		cc_version = "1.0.0"
	)

	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/creditbank/fixtures/channel-artifacts/Org3_1MSPanchors.tx",
		},

		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/creditbank/fixtures/channel-artifacts/Org3_2MSPanchors.tx",
		},

		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org3",
			OrgMspId:      "Org3MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/creditbank/fixtures/channel-artifacts/Org3_3MSPanchors.tx",
		},

		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org4",
			OrgMspId:      "Org4MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/creditbank/fixtures/channel-artifacts/Org3_4MSPanchors.tx",
		},
	}

	// init sdk env info
	info := sdkInit.SdkEnvInfo{
		ChannelID: []string{"channel1", "channel2", "channel3"},
		ChannelConfig: []string{os.Getenv("GOPATH") + "/src/creditbank/fixtures/channel-artifacts/channel1.tx",
			os.Getenv("GOPATH") + "/src/creditbank/fixtures/channel-artifacts/channel2.tx",
			os.Getenv("GOPATH") + "/src/creditbank/fixtures/channel-artifacts/channel3.tx"},
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    os.Getenv("GOPATH") + "/src/creditbank/chaincode/",
		ChaincodeVersion: cc_version,
	}

	// sdk setup
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	Channelsum := len(info.ChannelID)
	fmt.Fprintln(w, "通道数量 ：", Channelsum, "\n")

	for i, _ := range info.ChannelID {
		blh, curblock, preblock, err := service.QBlock_one(info.ChannelID[i], info.Orgs[0], sdk)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Fprintln(w, "channel", i+1, "当前区块高度 :", blh, "\n")
			fmt.Fprintln(w, "当前区块Hash :")
			fmt.Fprintln(w, curblock, "\n")
			fmt.Fprintln(w, "前一区块Hash :")
			fmt.Fprintln(w, preblock, "\n\n")
		}
	}

	for i, _ := range info.ChannelID {

		cfgID, cfgNum, err := service.ChannelInfo(info.ChannelID[i], info.Orgs[0], sdk)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Fprintln(w, "通道ID : ", cfgID, "\n")
			fmt.Fprintln(w, "区块编号 : ", cfgNum, "\n\n")
			// fmt.Println("锚节点:\n  主机:")
			// fmt.Println(cfgAnchHost)
			// fmt.Println("端口:")
			// fmt.Println(cfgAnchPort)
			// fmt.Println("机构:")
			// fmt.Println(cfgAnchOrg)

		}
	}

	// data := &struct {
	// 	CurrentUser User
	// }{
	// 	CurrentUser: cuser,
	// }

	ShowView(w, r, "qblock.html", nil)
}

func (app *Application) Index(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "index.html", nil)
}

func (app *Application) Help(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
	}{
		CurrentUser: cuser,
	}
	ShowView(w, r, "help.html", data)
}

func (app *Application) Queryhelp(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
	}{
		CurrentUser: cuser,
	}
	ShowView(w, r, "queryhelp.html", data)
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")

	var flag bool
	for _, user := range users {
		if user.LoginName == loginName && user.Password == password {
			cuser = user
			flag = true
			break
		}
	}

	data := &struct {
		CurrentUser User
		Flag        bool
	}{
		CurrentUser: cuser,
		Flag:        false,
	}

	if flag {
		// 登录成功
		ShowView(w, r, "index.html", data)
	} else {
		// 登录失败
		data.Flag = true
		data.CurrentUser.LoginName = loginName
		ShowView(w, r, "login.html", data)
	}
}

// 用户登出
func (app *Application) LoginOut(w http.ResponseWriter, r *http.Request) {
	cuser = User{}
	ShowView(w, r, "login.html", nil)
}

// 显示添加信息页面
func (app *Application) AddEduShow(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "addEdu.html", data)
}

// 添加信息
func (app *Application) AddEdu(w http.ResponseWriter, r *http.Request) {

	edu := service.Education{
		Name:           r.FormValue("name"),
		Gender:         r.FormValue("gender"),
		Nation:         r.FormValue("nation"),
		EntityID:       r.FormValue("entityID"),
		Place:          r.FormValue("place"),
		BirthDay:       r.FormValue("birthDay"),
		EnrollDate:     r.FormValue("enrollDate"),
		GraduationDate: r.FormValue("graduationDate"),
		SchoolName:     r.FormValue("schoolName"),
		Major:          r.FormValue("major"),
		QuaType:        r.FormValue("quaType"),
		Length:         r.FormValue("length"),
		Mode:           r.FormValue("mode"),
		Level:          r.FormValue("level"),
		Graduation:     r.FormValue("graduation"),
		CertNo:         r.FormValue("certNo"),
		Photo:          r.FormValue("photo"),
	}

	app.Setup.SaveEdu(edu)
	/*transactionID, err := app.Setup.SaveEdu(edu)

	data := &struct {
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Flag:true,
		Msg:"",
	}

	if err != nil {
		data.Msg = err.Error()
	}else{
		data.Msg = "信息添加成功:" + transactionID
	}*/

	//ShowView(w, r, "addEdu.html", data)
	r.Form.Set("certNo", edu.CertNo)
	r.Form.Set("name", edu.Name)
	app.FindCertByNoAndName(w, r)
}

func (app *Application) QueryPage(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "query.html", data)
}

// 根据证书编号与姓名查询信息
func (app *Application) FindCertByNoAndName(w http.ResponseWriter, r *http.Request) {
	certNo := r.FormValue("certNo")
	name := r.FormValue("name")
	result, err := app.Setup.FindEduByCertNoAndName(certNo, name)
	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	fmt.Println("根据证书编号与姓名查询信息成功：")
	fmt.Println(edu)

	data := &struct {
		Edu         service.Education
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		Edu:         edu,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     false,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryResult.html", data)
}

func (app *Application) QueryPage2(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "query2.html", data)
}

// 根据身份证号码查询信息
func (app *Application) FindByID(w http.ResponseWriter, r *http.Request) {
	entityID := r.FormValue("entityID")
	result, err := app.Setup.FindEduInfoByEntityID(entityID)
	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	data := &struct {
		Edu         service.Education
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		Edu:         edu,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     true,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "queryResult.html", data)
}

// 修改/添加新信息
func (app *Application) ModifyShow(w http.ResponseWriter, r *http.Request) {
	// 根据证书编号与姓名查询信息
	certNo := r.FormValue("certNo")
	name := r.FormValue("name")
	result, err := app.Setup.FindEduByCertNoAndName(certNo, name)

	var edu = service.Education{}
	json.Unmarshal(result, &edu)

	data := &struct {
		Edu         service.Education
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		Edu:         edu,
		CurrentUser: cuser,
		Flag:        true,
		Msg:         "",
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}

	ShowView(w, r, "modify.html", data)
}

// 修改/添加新信息
func (app *Application) Modify(w http.ResponseWriter, r *http.Request) {
	edu := service.Education{
		Name:           r.FormValue("name"),
		Gender:         r.FormValue("gender"),
		Nation:         r.FormValue("nation"),
		EntityID:       r.FormValue("entityID"),
		Place:          r.FormValue("place"),
		BirthDay:       r.FormValue("birthDay"),
		EnrollDate:     r.FormValue("enrollDate"),
		GraduationDate: r.FormValue("graduationDate"),
		SchoolName:     r.FormValue("schoolName"),
		Major:          r.FormValue("major"),
		QuaType:        r.FormValue("quaType"),
		Length:         r.FormValue("length"),
		Mode:           r.FormValue("mode"),
		Level:          r.FormValue("level"),
		Graduation:     r.FormValue("graduation"),
		CertNo:         r.FormValue("certNo"),
		Photo:          r.FormValue("photo"),
	}

	//transactionID, err := app.Setup.ModifyEdu(edu)
	app.Setup.ModifyEdu(edu)

	/*data := &struct {
		Edu service.Education
		CurrentUser User
		Msg string
		Flag bool
	}{
		CurrentUser:cuser,
		Flag:true,
		Msg:"",
	}

	if err != nil {
		data.Msg = err.Error()
	}else{
		data.Msg = "新信息添加成功:" + transactionID
	}

	ShowView(w, r, "modify.html", data)
	*/

	r.Form.Set("certNo", edu.CertNo)
	r.Form.Set("name", edu.Name)
	app.FindCertByNoAndName(w, r)
}
