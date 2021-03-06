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
	fmt.Println(">> ??????????????????????????????????????????......")

	edu := service.Education{
		Name:           "??????",
		Gender:         "???",
		Nation:         "???",
		EntityID:       "101",
		Place:          "??????",
		BirthDay:       "1991???01???01???",
		EnrollDate:     "2009???9???",
		GraduationDate: "2013???7???",
		SchoolName:     "??????????????????",
		Major:          "?????????",
		QuaType:        "??????",
		Length:         "??????",
		Mode:           "???????????????",
		Level:          "??????",
		Graduation:     "??????",
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
		fmt.Println("??????????????????, ???????????????: " + msg)
	}

	result, err := serviceSetup.FindEduInfoByEntityID("101")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var edu service.Education
		json.Unmarshal(result, &edu)
		fmt.Println("??????????????????????????????????????????")
		fmt.Println(edu)
	}

	Channelsum := len(info.ChannelID)
	fmt.Println("???????????? ???", Channelsum)

	for i, _ := range info.ChannelID {
		blh, curblock, preblock, err := service.QBlock_one(info.ChannelID[i], info.Orgs[0], sdk)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("channel", i+1, "?????????????????? :", blh)
			fmt.Println("????????????Hash :")
			fmt.Println(curblock)
			fmt.Println("????????????Hash :")
			fmt.Println(preblock, "\n")

		}
	}
	for i, _ := range info.ChannelID {
		// cfgID, cfgNum, cfgAnchHost, cfgAnchPort, cfgAnchOrg, err := service.ChannelInfo(info.ChannelID[i], info.Orgs[0], sdk)
		cfgID, cfgNum, err := service.ChannelInfo(info.ChannelID[i], info.Orgs[0], sdk)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("??????ID : ", cfgID)
			fmt.Println("???????????? : ", cfgNum)
			// fmt.Println("??????1?????????:\n  ??????: ", cfgAnchHost)
			// fmt.Println("??????: ", cfgAnchPort)
			// fmt.Println("??????: ", cfgAnchOrg, "\n")

		}
	}

	// //?????????????????????????????????
	// options_user := fabsdk.WithUser(info.OrdererAdminUser)
	// options_org := fabsdk.WithOrg(info.OrdererOrgName)
	// clientChannelContext := sdk.ChannelContext(info.ChannelID, options_user, options_org)
	// client, err := ledger.New(clientChannelContext)

	// inf, err := client.QueryInfo()
	// // inf, err := client.QueryInfo()
	// // if err != nil {
	// // 	log.Fatalf("?????????????????????: %v\n", err)
	// // 	return
	// // }

	// fmt.Printf("????????????:\n%v\n", inf.BCI.Height)
	// fmt.Printf("????????????Hash:\n%v\n", hex.EncodeToString(inf.BCI.CurrentBlockHash))
	// fmt.Printf("????????????Hash:\n%v\n", hex.EncodeToString(inf.BCI.PreviousBlockHash))

	// // -------------------- ???1???????????? ???????????????????????? ----------------
	// block, err := client.QueryBlockByHash(inf.BCI.CurrentBlockHash)
	// if err != nil {
	// 	log.Fatalf("????????????????????????: %v\n", err)
	// 	return
	// }
	// fmt.Printf("????????????: %v\n", block.Header.Number)
	// fmt.Printf("??????Hash:\n%v\n", hex.EncodeToString(block.Header.DataHash))

	// // -------------------- ???1???????????? ???????????????????????? ----------------
	// blockNumber := inf.BCI.Height - 1
	// block, err = client.QueryBlock(blockNumber)

	// fmt.Printf("????????????: %v\n", block.Header.Number)
	// fmt.Printf("??????Hash:\n%v\n", hex.EncodeToString(block.Header.DataHash))

	// cfg, err := client.QueryConfig()
	// fmt.Printf("????????????: %v\n", cfg.ID())
	// fmt.Printf("????????????: %v\n", cfg.BlockNumber())
	// fmt.Printf("?????????:\n  ??????:%v\n  ??????:%v\n  ??????:%v\n", cfg.AnchorPeers()[0].Host, cfg.AnchorPeers()[0].Port, cfg.AnchorPeers()[0].Org)
	// //??????????????????????????????

	app := controller.Application{
		Setup: serviceSetup,
	}
	web.WebStart(app)
}
