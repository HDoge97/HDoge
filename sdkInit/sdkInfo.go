package sdkInit

import (
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	// "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	// "github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	// "github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	// "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type OrgInfo struct {
	OrgAdminUser          string // like "Admin"
	OrgName               string // like "Org1"
	OrgMspId              string // like "Org1MSP"
	OrgUser               string // like "User1"
	orgMspClient          *mspclient.Client
	OrgAdminClientContext *contextAPI.ClientProvider
	OrgResMgmt            *resmgmt.Client
	OrgPeerNum            int
	//Peers                 []*fab.Peer
	OrgAnchorFile string // like ./channel-artifacts/Org2MSPanchors.tx

	// //以下为模板中的部分
	// client          []*channel.Client //1.chainhero 溯源 2.token
	// ledger          []*ledger.Client  //1.chainhero 溯源 2.token
	// admin           []*resmgmt.Client //1.chainhero 溯源 2.token
	// sdk             *fabsdk.FabricSDK //1.chainhero 溯源 2.token
	// event           []*event.Client   //1.chainhero 溯源 2.token
	// //以上为模板中的部分

}

type SdkEnvInfo struct {
	// 通道信息
	ChannelID     []string // like "simplecc"
	ChannelConfig []string // like os.Getenv("GOPATH") + "/src/github.com/hyperledger/fabric-samples/test-network/channel-artifacts/testchannel.tx"

	// 组织信息
	Orgs []*OrgInfo
	// 排序服务节点信息
	OrdererAdminUser     string // like "Admin"
	OrdererOrgName       string // like "OrdererOrg"
	OrdererEndpoint      string
	OrdererClientContext *contextAPI.ClientProvider
	// 链码信息
	ChaincodeID      string
	ChaincodeGoPath  string
	ChaincodePath    string
	ChaincodeVersion string
}
