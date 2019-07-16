// +build unit

package userapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/centrifuge/go-centrifuge/documents"
	"github.com/centrifuge/go-centrifuge/documents/invoice"
	"github.com/centrifuge/go-centrifuge/errors"
	"github.com/centrifuge/go-centrifuge/extensions"
	"github.com/centrifuge/go-centrifuge/extensions/funding"
	"github.com/centrifuge/go-centrifuge/httpapi/coreapi"
	"github.com/centrifuge/go-centrifuge/jobs"
	testingconfig "github.com/centrifuge/go-centrifuge/testingutils/config"
	testingdocuments "github.com/centrifuge/go-centrifuge/testingutils/documents"
	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/centrifuge/go-centrifuge/utils/byteutils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CreateFundingAgreement(t *testing.T) {
	getHTTPReqAndResp := func(ctx context.Context, b io.Reader) (*httptest.ResponseRecorder, *http.Request) {
		return httptest.NewRecorder(), httptest.NewRequest("POST", "/documents/{document_id}/funding_agreements", b).WithContext(ctx)
	}
	// empty document_id and invalid id
	rctx := chi.NewRouteContext()
	rctx.URLParams.Keys = make([]string, 1, 1)
	rctx.URLParams.Values = make([]string, 1, 1)
	rctx.URLParams.Keys[0] = "document_id"
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)
	h := handler{}
	for _, id := range []string{"", "invalid"} {
		rctx.URLParams.Values[0] = id
		w, r := getHTTPReqAndResp(ctx, nil)
		h.CreateFundingAgreement(w, r)
		assert.Equal(t, w.Code, http.StatusBadRequest)
		assert.Contains(t, w.Body.String(), coreapi.ErrInvalidDocumentID.Error())
	}

	// empty body
	rctx.URLParams.Values[0] = byteutils.HexBytes(utils.RandomSlice(32)).String()
	w, r := getHTTPReqAndResp(ctx, nil)
	h.CreateFundingAgreement(w, r)
	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "unexpected end of JSON input")

	// creation failed
	d, err := json.Marshal(map[string]interface{}{
		"data": funding.Data{},
	})
	assert.NoError(t, err)
	fundingSrv := new(funding.MockService)
	fundingSrv.On("CreateFundingAgreement", mock.Anything, mock.Anything, mock.Anything).Return(nil, jobs.NilJobID(), errors.New("failed to create funding agreement")).Once()
	h.srv = Service{fundingSrv: fundingSrv}
	w, r = getHTTPReqAndResp(ctx, bytes.NewReader(d))
	h.CreateFundingAgreement(w, r)
	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "failed to create funding agreement")

	// failed response conversion
	m := new(testingdocuments.MockModel)
	m.On("ID").Return(utils.RandomSlice(32))
	m.On("CurrentVersion").Return(utils.RandomSlice(32))
	m.On("Author").Return(nil, errors.New("somerror"))
	m.On("Timestamp").Return(nil, errors.New("somerror"))
	m.On("NFTs").Return(nil)
	m.On("GetCollaborators", mock.Anything).Return(documents.CollaboratorsAccess{}, errors.New("failed to get collaborators")).Once()
	fundingSrv.On("CreateFundingAgreement", mock.Anything, mock.Anything, mock.Anything).Return(m, jobs.NewJobID(), nil)
	w, r = getHTTPReqAndResp(ctx, bytes.NewReader(d))
	h.CreateFundingAgreement(w, r)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
	assert.Contains(t, w.Body.String(), "failed to get collaborators")

	// success
	m.On("GetCollaborators", mock.Anything).Return(documents.CollaboratorsAccess{}, nil)
	fundingSrv.On("GetDataAndSignatures", mock.Anything, mock.Anything, mock.Anything).Return(funding.Data{}, nil, nil)
	w, r = getHTTPReqAndResp(ctx, bytes.NewReader(d))
	h.CreateFundingAgreement(w, r)
	assert.Equal(t, w.Code, http.StatusAccepted)
	fundingSrv.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestHandler_GetFundingAgreements(t *testing.T) {
	getHTTPReqAndResp := func(ctx context.Context) (*httptest.ResponseRecorder, *http.Request) {
		return httptest.NewRecorder(), httptest.NewRequest("GET", "/documents/{document_id}/funding_agreements", nil).WithContext(ctx)
	}
	// empty document_id and invalid id
	rctx := chi.NewRouteContext()
	rctx.URLParams.Keys = make([]string, 1, 1)
	rctx.URLParams.Values = make([]string, 1, 1)
	rctx.URLParams.Keys[0] = "document_id"
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)
	h := handler{}
	for _, id := range []string{"", "invalid"} {
		rctx.URLParams.Values[0] = id
		w, r := getHTTPReqAndResp(ctx)
		h.GetFundingAgreements(w, r)
		assert.Equal(t, w.Code, http.StatusBadRequest)
		assert.Contains(t, w.Body.String(), coreapi.ErrInvalidDocumentID.Error())
	}

	// missing Doc
	id := utils.RandomSlice(32)
	rctx.URLParams.Values[0] = byteutils.HexBytes(id).String()
	docSrv := new(testingdocuments.MockService)
	docSrv.On("GetCurrentVersion", id).Return(nil, errors.New("doc not found")).Once()
	h.srv.coreAPISrv = newCoreAPIService(docSrv)
	w, r := getHTTPReqAndResp(ctx)
	h.GetFundingAgreements(w, r)
	assert.Equal(t, w.Code, http.StatusNotFound)

	// failed conversion
	m := new(testingdocuments.MockModel)
	m.On("ID").Return(utils.RandomSlice(32))
	m.On("CurrentVersion").Return(utils.RandomSlice(32))
	m.On("Author").Return(nil, errors.New("somerror"))
	m.On("Timestamp").Return(nil, errors.New("somerror"))
	m.On("NFTs").Return(nil)
	m.On("GetCollaborators", mock.Anything).Return(documents.CollaboratorsAccess{}, errors.New("failed to get collaborators")).Once()
	docSrv.On("GetCurrentVersion", id).Return(m, nil)
	w, r = getHTTPReqAndResp(ctx)
	h.GetFundingAgreements(w, r)
	assert.Equal(t, w.Code, http.StatusInternalServerError)

	// success
	m.On("GetCollaborators", mock.Anything).Return(documents.CollaboratorsAccess{}, nil)
	m.On("AttributeExists", mock.Anything).Return(false)
	w, r = getHTTPReqAndResp(ctx)
	h.GetFundingAgreements(w, r)
	assert.Equal(t, w.Code, http.StatusOK)
	docSrv.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestHandler_GetFundingAgreement(t *testing.T) {
	getHTTPReqAndResp := func(ctx context.Context) (*httptest.ResponseRecorder, *http.Request) {
		return httptest.NewRecorder(), httptest.NewRequest("GET", "/documents/{document_id}/funding_agreements/{agreement_id}", nil).WithContext(ctx)
	}
	// empty document_id and invalid id
	rctx := chi.NewRouteContext()
	rctx.URLParams.Keys = make([]string, 2, 2)
	rctx.URLParams.Values = make([]string, 2, 2)
	rctx.URLParams.Keys[0] = "document_id"
	rctx.URLParams.Keys[1] = "agreement_id"
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)
	h := handler{}
	for _, id := range []string{"", "invalid"} {
		rctx.URLParams.Values[0] = id
		w, r := getHTTPReqAndResp(ctx)
		h.GetFundingAgreement(w, r)
		assert.Equal(t, w.Code, http.StatusBadRequest)
		assert.Contains(t, w.Body.String(), coreapi.ErrInvalidDocumentID.Error())
	}

	id := utils.RandomSlice(32)
	rctx.URLParams.Values[0] = byteutils.HexBytes(id).String()
	for _, id := range []string{"", "invalid"} {
		rctx.URLParams.Values[1] = id
		w, r := getHTTPReqAndResp(ctx)
		h.GetFundingAgreement(w, r)
		assert.Equal(t, w.Code, http.StatusBadRequest)
		assert.Contains(t, w.Body.String(), ErrInvalidAgreementID.Error())
	}

	// missing Doc
	fundingID := hexutil.Encode(utils.RandomSlice(32))
	rctx.URLParams.Values[1] = fundingID
	docSrv := new(testingdocuments.MockService)
	docSrv.On("GetCurrentVersion", id).Return(nil, errors.New("doc not found")).Once()
	h.srv.coreAPISrv = newCoreAPIService(docSrv)
	w, r := getHTTPReqAndResp(ctx)
	h.GetFundingAgreement(w, r)
	assert.Equal(t, w.Code, http.StatusNotFound)

	// failed response conversion
	fundingSrv := new(funding.MockService)
	h.srv.fundingSrv = fundingSrv
	m := new(testingdocuments.MockModel)
	docSrv.On("GetCurrentVersion", id).Return(m, nil)
	m.On("ID").Return(utils.RandomSlice(32))
	m.On("CurrentVersion").Return(utils.RandomSlice(32))
	m.On("Author").Return(nil, errors.New("somerror"))
	m.On("Timestamp").Return(nil, errors.New("somerror"))
	m.On("NFTs").Return(nil)
	m.On("GetCollaborators", mock.Anything).Return(documents.CollaboratorsAccess{}, errors.New("failed to get collaborators")).Once()
	w, r = getHTTPReqAndResp(ctx)
	h.GetFundingAgreement(w, r)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
	assert.Contains(t, w.Body.String(), "failed to get collaborators")

	// success
	m.On("GetCollaborators", mock.Anything).Return(documents.CollaboratorsAccess{}, nil)
	fundingSrv.On("GetDataAndSignatures", mock.Anything, mock.Anything, mock.Anything).Return(funding.Data{}, nil, nil)
	w, r = getHTTPReqAndResp(ctx)
	h.GetFundingAgreement(w, r)
	assert.Equal(t, w.Code, http.StatusOK)
	fundingSrv.AssertExpectations(t)
	m.AssertExpectations(t)
	fundingSrv.AssertExpectations(t)
	docSrv.AssertExpectations(t)
}

func TestHandler_UpdateFundingAgreement(t *testing.T) {
	getHTTPReqAndResp := func(ctx context.Context, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
		return httptest.NewRecorder(), httptest.NewRequest("PUT", "/documents/{document_id}/funding_agreements/{agreement_id}", body).WithContext(ctx)
	}
	// empty document_id and invalid id
	rctx := chi.NewRouteContext()
	rctx.URLParams.Keys = make([]string, 2, 2)
	rctx.URLParams.Values = make([]string, 2, 2)
	rctx.URLParams.Keys[0] = "document_id"
	rctx.URLParams.Keys[1] = "agreement_id"
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)
	h := handler{}
	for _, id := range []string{"", "invalid"} {
		rctx.URLParams.Values[0] = id
		w, r := getHTTPReqAndResp(ctx, nil)
		h.UpdateFundingAgreement(w, r)
		assert.Equal(t, w.Code, http.StatusBadRequest)
		assert.Contains(t, w.Body.String(), coreapi.ErrInvalidDocumentID.Error())
	}

	id := utils.RandomSlice(32)
	rctx.URLParams.Values[0] = byteutils.HexBytes(id).String()
	for _, id := range []string{"", "invalid"} {
		rctx.URLParams.Values[1] = id
		w, r := getHTTPReqAndResp(ctx, nil)
		h.UpdateFundingAgreement(w, r)
		assert.Equal(t, w.Code, http.StatusBadRequest)
		assert.Contains(t, w.Body.String(), ErrInvalidAgreementID.Error())
	}

	// empty body
	fundingID := utils.RandomSlice(32)
	rctx.URLParams.Values[1] = hexutil.Encode(fundingID)
	w, r := getHTTPReqAndResp(ctx, nil)
	h.UpdateFundingAgreement(w, r)
	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "unexpected end of JSON input")

	// update failed
	fundingSrv := new(funding.MockService)
	data := funding.CreateData()
	rctx.URLParams.Values[1] = data.AgreementID
	fundingID, err := hexutil.Decode(data.AgreementID)
	assert.NoError(t, err)
	fundingSrv.On("UpdateFundingAgreement", mock.Anything, id, fundingID, mock.Anything).Return(nil, nil, errors.New("failed to update")).Once()
	h.srv.fundingSrv = fundingSrv
	d, err := json.Marshal(map[string]interface{}{
		"data": data,
	})
	assert.NoError(t, err)
	w, r = getHTTPReqAndResp(ctx, bytes.NewReader(d))
	h.UpdateFundingAgreement(w, r)
	assert.Equal(t, w.Code, http.StatusNotFound)
	assert.Contains(t, w.Body.String(), "failed to update")

	// success
	inv, _ := invoice.CreateInvoiceWithEmbedCD(t, testingconfig.CreateAccountContext(t, cfg), did, nil)
	attrs, err := extensions.CreateAttributesList(inv, *data, "funding_agreement[{IDX}].", funding.AttrFundingLabel)
	assert.NoError(t, err)
	err = inv.AddAttributes(documents.CollaboratorsAccess{}, false, attrs...)
	assert.NoError(t, err)
	fundingSrv.On("UpdateFundingAgreement", mock.Anything, id, fundingID, mock.Anything).Return(inv, jobs.NewJobID(), nil)
	fundingSrv.On("GetDataAndSignatures", mock.Anything, mock.Anything, mock.Anything).Return(*data, nil, nil)
	w, r = getHTTPReqAndResp(ctx, bytes.NewReader(d))
	h.UpdateFundingAgreement(w, r)
	assert.Equal(t, w.Code, http.StatusAccepted)
	assert.Contains(t, w.Body.String(), data.AgreementID)
	fundingSrv.AssertExpectations(t)
}