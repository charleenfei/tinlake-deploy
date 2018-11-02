package cmd

import (
	"io/ioutil"

	"github.com/centrifuge/go-centrifuge/identity"
	"github.com/spf13/cobra"
)

var centrifugeIdString string
var purpose string

var createIdentityCmd = &cobra.Command{
	Use:   "createidentity",
	Short: "creates identity with signing key as p2p id against ethereum",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		//cmd requires a config file
		readConfigFile()
		baseBootstrap()
		identityService := identity.EthereumIdentityService{}
		var centrifugeId identity.CentID
		var err error
		if centrifugeIdString == "" {
			centrifugeId = identity.RandomCentID()
		} else {
			centrifugeId, err = identity.CentIDFromString(centrifugeIdString)
			if err != nil {
				panic(err)
			}
		}
		_, confirmations, err := identityService.CreateIdentity(centrifugeId)
		if err != nil {
			panic(err)
		}
		watchIdentity := <-confirmations
		log.Infof("Identity created [%s]", watchIdentity.Identity.CentID().String())
		// We need a way to return the identity created so it can be read by an automated process as well
		// when id autogenerated
		id := []byte("{\"id\": \"" + centrifugeId.String() + "\"}")
		err = ioutil.WriteFile("newidentity.json", id, 0644)
		if err != nil {
			panic(err)
		}
		log.Infof("Identity created [%s]", watchIdentity.Identity.CentID())
	},
}

//We should support multiple types of keys to add, at the moment only keyPurpose 1 - PeerID/Signature/Encryption
var addKeyCmd = &cobra.Command{
	Use:   "addkey",
	Short: "add a signing key as p2p id against ethereum",
	Long:  "add a signing key as p2p id against ethereum",
	Run: func(cmd *cobra.Command, args []string) {
		//cmd requires a config file
		readConfigFile()

		baseBootstrap()

		var purposeInt int

		switch purpose {
		case "p2p":
			purposeInt = identity.KeyPurposeP2P
		case "sign":
			purposeInt = identity.KeyPurposeSigning
		case "ethauth":
			purposeInt = identity.KeyPurposeEthMsgAuth
		default:
			panic("Option not supported")
		}

		err := identity.AddKeyFromConfig(purposeInt)
		if err != nil {
			panic(err)
		}

		return
	},
}

func init() {
	createIdentityCmd.Flags().StringVarP(&centrifugeIdString, "centrifugeid", "i", "", "Centrifuge ID")
	addKeyCmd.Flags().StringVarP(&centrifugeIdString, "centrifugeid", "i", "", "Centrifuge ID")
	addKeyCmd.Flags().StringVarP(&purpose, "purpose", "p", "", "Key Purpose [p2p|sign|ethauth]")
	rootCmd.AddCommand(createIdentityCmd)
	rootCmd.AddCommand(addKeyCmd)
}