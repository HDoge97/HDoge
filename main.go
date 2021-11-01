package main

import (
	"creditbank/sdkInit"
	"creditbank/service"
	"creditbank/web"
	"creditbank/web/controller"

	// "encoding/hex"
	"encoding/json"
	"fmt"

	// "log"
	"os"
	// "github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	// "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	cc_name    = "simplecc"
	cc_version = "1.0.0"
)

func main() {
	// init orgs information
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/creditbank/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},

		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/creditbank/fixtures/channel-artifacts/Org2MSPanchors.tx",
		},

		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org3",
			OrgMspId:      "Org3MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/creditbank/fixtures/channel-artifacts/Org3MSPanchors.tx",
		},

		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org4",
			OrgMspId:      "Org4MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/creditbank/fixtures/channel-artifacts/Org4MSPanchors.tx",
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

	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Println(">> create chaincode lifecycle error: %v", err)
		os.Exit(-1)
	}

	// invoke chaincode set status
	fmt.Println(">> 通过链码外部服务设置链码状态......")

	edu := service.Education{
		Name:           "张三",
		Gender:         "男",
		Nation:         "汉",
		EntityID:       "101",
		Place:          "北京",
		BirthDay:       "1991年01月01日",
		EnrollDate:     "2009年9月",
		GraduationDate: "2013年7月",
		SchoolName:     "中国政法大学",
		Major:          "社会学",
		QuaType:        "普通",
		Length:         "四年",
		Mode:           "普通全日制",
		Level:          "本科",
		Graduation:     "毕业",
		CertNo:         "111",
		Photo:          "/static/photo/11.png",
	}

	serviceSetup, err := service.InitService(info.ChaincodeID, info.ChannelID[0], info.Orgs[0], sdk)
	if err != nil {
		fmt.Println()
		os.Exit(-1)
	}
	msg, err := serviceSetup.SaveEdu(edu)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}

	result, err := serviceSetup.FindEduInfoByEntityID("101")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var edu service.Education
		json.Unmarshal(result, &edu)
		fmt.Println("根据身份证号码查询信息成功：")
		fmt.Println(edu)
	}

	Channelsum := len(info.ChannelID)
	fmt.Println("通道数量 ：", Channelsum)

	for i, _ := range info.ChannelID {
		blh, curblock, preblock, err := service.QBlock_one(info.ChannelID[i], info.Orgs[0], sdk)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("channel", i+1, "当前区块高度 :", blh)
			fmt.Println("当前区块Hash :")
			fmt.Println(curblock)
			fmt.Println("前一区块Hash :")
			fmt.Println(preblock, "\n")

		}
	}
	for i, _ := range info.ChannelID {
		// cfgID, cfgNum, cfgAnchHost, cfgAnchPort, cfgAnchOrg, err := service.ChannelInfo(info.ChannelID[i], info.Orgs[0], sdk)
		cfgID, cfgNum, err := service.ChannelInfo(info.ChannelID[i], info.Orgs[0], sdk)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("通道ID : ", cfgID)
			fmt.Println("区块编号 : ", cfgNum)
			// fmt.Println("组织1锚节点:\n  主机: ", cfgAnchHost)
			// fmt.Println("端口: ", cfgAnchPort)
			// fmt.Println("机构: ", cfgAnchOrg, "\n")

		}
	}

	// //以下测试获得区块信息：
	// options_user := fabsdk.WithUser(info.OrdererAdminUser)
	// options_org := fabsdk.WithOrg(info.OrdererOrgName)
	// clientChannelContext := sdk.ChannelContext(info.ChannelID, options_user, options_org)
	// client, err := ledger.New(clientChannelContext)

	// inf, err := client.QueryInfo()
	// // inf, err := client.QueryInfo()
	// // if err != nil {
	// // 	log.Fatalf("查询区块链概况: %v\n", err)
	// // 	return
	// // }

	// fmt.Printf("区块高度:\n%v\n", inf.BCI.Height)
	// fmt.Printf("当前区块Hash:\n%v\n", hex.EncodeToString(inf.BCI.CurrentBlockHash))
	// fmt.Printf("前一区块Hash:\n%v\n", hex.EncodeToString(inf.BCI.PreviousBlockHash))

	// // -------------------- 第1种方式： 根据哈希查询区块 ----------------
	// block, err := client.QueryBlockByHash(inf.BCI.CurrentBlockHash)
	// if err != nil {
	// 	log.Fatalf("查询区块信息失败: %v\n", err)
	// 	return
	// }
	// fmt.Printf("区块编号: %v\n", block.Header.Number)
	// fmt.Printf("区块Hash:\n%v\n", hex.EncodeToString(block.Header.DataHash))

	// // -------------------- 第1种方式： 根据块号查询区块 ----------------
	// blockNumber := inf.BCI.Height - 1
	// block, err = client.QueryBlock(blockNumber)

	// fmt.Printf("区块编号: %v\n", block.Header.Number)
	// fmt.Printf("区块Hash:\n%v\n", hex.EncodeToString(block.Header.DataHash))

	// cfg, err := client.QueryConfig()
	// fmt.Printf("通道名称: %v\n", cfg.ID())
	// fmt.Printf("区块个数: %v\n", cfg.BlockNumber())
	// fmt.Printf("锚节点:\n  主机:%v\n  端口:%v\n  机构:%v\n", cfg.AnchorPeers()[0].Host, cfg.AnchorPeers()[0].Port, cfg.AnchorPeers()[0].Org)
	// //以上测试获得区块信息

	app := controller.Application{
		Setup: serviceSetup,
	}
	web.WebStart(app)
}
