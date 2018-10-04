// +build unit

package invoice

import (
	"encoding/hex"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/centrifuge/centrifuge-protobufs/gen/go/coredocument"
	"github.com/centrifuge/centrifuge-protobufs/gen/go/invoice"
	"github.com/centrifuge/go-centrifuge/centrifuge/documents"
	"github.com/centrifuge/go-centrifuge/centrifuge/identity"
	clientinvoicepb "github.com/centrifuge/go-centrifuge/centrifuge/protobufs/gen/go/invoice"
	"github.com/centrifuge/go-centrifuge/centrifuge/testingutils/documents"
	"github.com/centrifuge/go-centrifuge/centrifuge/tools"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestInvoice_FromCoreDocuments_invalidParameter(t *testing.T) {
	invoiceModel := &InvoiceModel{}

	emptyCoreDocument := &coredocumentpb.CoreDocument{}
	err := invoiceModel.UnpackCoreDocument(emptyCoreDocument)
	assert.Error(t, err, "it should not be possible to init a empty core document")

	err = invoiceModel.UnpackCoreDocument(nil)
	assert.Error(t, err, "it should not be possible to init a empty core document")

	invalidEmbeddedData := &any.Any{TypeUrl: "invalid"}
	coreDocument := &coredocumentpb.CoreDocument{EmbeddedData: invalidEmbeddedData}
	err = invoiceModel.UnpackCoreDocument(coreDocument)
	assert.Error(t, err, "it should not be possible to init invalid typeUrl")

}

func TestInvoice_InitCoreDocument_successful(t *testing.T) {
	invoiceModel := &InvoiceModel{}

	coreDocument := testinginvoice.CreateCDWithEmbeddedInvoice(t, testinginvoice.CreateInvoiceData())
	err := invoiceModel.UnpackCoreDocument(coreDocument)
	assert.Nil(t, err, "valid coredocument shouldn't produce an error")
}

func TestInvoice_InitCoreDocument_invalidCentId(t *testing.T) {
	invoiceModel := &InvoiceModel{}

	coreDocument := testinginvoice.CreateCDWithEmbeddedInvoice(t, invoicepb.InvoiceData{
		Recipient:   tools.RandomSlice(identity.CentIDByteLength + 1),
		Sender:      tools.RandomSlice(identity.CentIDByteLength),
		Payee:       tools.RandomSlice(identity.CentIDByteLength),
		GrossAmount: 42,
	})
	err := invoiceModel.UnpackCoreDocument(coreDocument)
	assert.Error(t, err, "invalid centID should produce an error")

}

func TestInvoice_CoreDocument_successful(t *testing.T) {
	invoiceModel := &InvoiceModel{}

	//init model with a CoreDoc
	coreDocument := testinginvoice.CreateCDWithEmbeddedInvoice(t, testinginvoice.CreateInvoiceData())
	invoiceModel.UnpackCoreDocument(coreDocument)

	returnedCoreDocument, err := invoiceModel.PackCoreDocument()
	assert.Nil(t, err, "transformation from invoice to CoreDoc failed")

	assert.Equal(t, coreDocument.EmbeddedData, returnedCoreDocument.EmbeddedData, "embeddedData should be the same")
	assert.Equal(t, coreDocument.EmbeddedDataSalts, returnedCoreDocument.EmbeddedDataSalts, "embeddedDataSalt should be the same")
}

func TestInvoice_ModelInterface(t *testing.T) {
	var i interface{} = &InvoiceModel{}
	_, ok := i.(documents.Model)
	assert.True(t, ok, "model interface not implemented correctly for invoiceModel")
}

func TestInvoice_Type(t *testing.T) {
	var model documents.Model
	model = &InvoiceModel{}
	assert.Equal(t, model.Type(), reflect.TypeOf(&InvoiceModel{}), "InvoiceType not correct")
}

func TestInvoice_JSON(t *testing.T) {
	invoiceModel := &InvoiceModel{}

	//init model with a CoreDoc
	coreDocument := testinginvoice.CreateCDWithEmbeddedInvoice(t, testinginvoice.CreateInvoiceData())
	invoiceModel.UnpackCoreDocument(coreDocument)

	jsonBytes, err := invoiceModel.JSON()
	assert.Nil(t, err, "marshal to json didn't work correctly")
	assert.True(t, json.Valid(jsonBytes), "json format not correct")

	err = invoiceModel.FromJSON(jsonBytes)
	assert.Nil(t, err, "unmarshal JSON didn't work correctly")

	receivedCoreDocument, err := invoiceModel.PackCoreDocument()
	assert.Nil(t, err, "JSON unmarshal damaged invoice variables")
	assert.Equal(t, receivedCoreDocument.EmbeddedData, coreDocument.EmbeddedData, "JSON unmarshal damaged invoice variables")
}

func TestInvoiceModel_UnpackCoreDocument(t *testing.T) {
	var model documents.Model = new(InvoiceModel)
	var err error

	// nil core doc
	err = model.UnpackCoreDocument(nil)
	assert.Error(t, err, "unpack must fail")

	// embed data missing
	err = model.UnpackCoreDocument(new(coredocumentpb.CoreDocument))
	assert.Error(t, err, "unpack must fail due to missing embed data")

	// successful
	coreDocument := testinginvoice.CreateCDWithEmbeddedInvoice(t, testinginvoice.CreateInvoiceData())
	err = model.UnpackCoreDocument(coreDocument)
	assert.Nil(t, err, "valid core document with embedded invoice shouldn't produce an error")

	receivedCoreDocument, err := model.PackCoreDocument()
	assert.Nil(t, err, "model should be able to return the core document with embedded invoice")

	assert.Equal(t, coreDocument.EmbeddedData, receivedCoreDocument.EmbeddedData, "embeddedData should be the same")
	assert.Equal(t, coreDocument.EmbeddedDataSalts, receivedCoreDocument.EmbeddedDataSalts, "embeddedDataSalt should be the same")
}

func TestInvoiceModel_getClientData(t *testing.T) {
	invData := testinginvoice.CreateInvoiceData()
	inv := new(InvoiceModel)
	err := inv.loadFromP2PProtobuf(&invData)
	assert.Nil(t, err, "must not error out")

	data, err := inv.getClientData()
	assert.Nil(t, err, "must not error out")
	assert.NotNil(t, data, "invoice data should not be nil")
	assert.Equal(t, data.GrossAmount, invData.GrossAmount, "gross amount must match")
	assert.Equal(t, data.Recipient, hex.EncodeToString(inv.Recipient[:]), "recipient should match")
	assert.Equal(t, data.Sender, hex.EncodeToString(inv.Sender[:]), "sender should match")
	assert.Equal(t, data.Payee, hex.EncodeToString(inv.Payee[:]), "payee should match")
}

func TestInvoiceModel_InitInvoiceInput(t *testing.T) {
	// fail recipient
	data := &clientinvoicepb.InvoiceData{
		Sender:    "some number",
		Payee:     "some payee",
		Recipient: "some recipient",
		ExtraData: "some data",
	}
	inv := new(InvoiceModel)
	err := inv.InitInvoiceInput(&clientinvoicepb.InvoiceCreatePayload{Data: data})
	assert.Error(t, err, "must return err")
	assert.Contains(t, err.Error(), "failed to decode recipient")

	data.Recipient = "010203040506"
	err = inv.InitInvoiceInput(&clientinvoicepb.InvoiceCreatePayload{Data: data})
	assert.Contains(t, err.Error(), "failed to decode sender")

	data.Sender = "010203060506"
	err = inv.InitInvoiceInput(&clientinvoicepb.InvoiceCreatePayload{Data: data})
	assert.Contains(t, err.Error(), "failed to decode payee")

	data.Payee = "010203030405"
	err = inv.InitInvoiceInput(&clientinvoicepb.InvoiceCreatePayload{Data: data})
	assert.Contains(t, err.Error(), "failed to decode extra data")

	data.ExtraData = "010203020301"
	collabs := []string{"010102040506", "some id"}
	err = inv.InitInvoiceInput(&clientinvoicepb.InvoiceCreatePayload{Data: data, Collaborators: collabs})
	assert.Contains(t, err.Error(), "failed to decode collaborator")

	collabs = []string{"010102040506", "010203020302"}
	err = inv.InitInvoiceInput(&clientinvoicepb.InvoiceCreatePayload{Data: data, Collaborators: collabs})
	assert.Nil(t, err, "must be nil")
	assert.Equal(t, inv.Sender[:], []byte{1, 2, 3, 6, 5, 6})
	assert.Equal(t, inv.Payee[:], []byte{1, 2, 3, 3, 4, 5})
	assert.Equal(t, inv.Recipient[:], []byte{1, 2, 3, 4, 5, 6})
	assert.Equal(t, inv.ExtraData[:], []byte{1, 2, 3, 2, 3, 1})
	assert.Equal(t, inv.Collaborators, []identity.CentID{{1, 1, 2, 4, 5, 6}, {1, 1, 2, 4, 5, 6}, {1, 2, 3, 2, 3, 2}})
}