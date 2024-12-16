package digest

import (
	"net/http"
	"time"

	"github.com/ProtoconNet/mitum-credential/types"
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

func (hd *Handlers) handleCredentialService(w http.ResponseWriter, r *http.Request) {
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleCredentialServiceInGroup(contract)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleCredentialServiceInGroup(contract string) (interface{}, error) {
	switch design, err := CredentialService(hd.database, contract); {
	case err != nil:
		return nil, util.ErrNotFound.WithMessage(err, "credential design, contract %s", contract)
	case design == nil:
		return nil, util.ErrNotFound.Errorf("credential design, contract %s", contract)
	default:
		hal, err := hd.buildCredentialServiceHal(contract, *design)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildCredentialServiceHal(contract string, design types.Design) (cdigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDIDService, "contract", contract)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(design, cdigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleCredential(w http.ResponseWriter, r *http.Request) {
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	templateID, err, status := cdigest.ParseRequest(w, r, "template_id")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	credentialID, err, status := cdigest.ParseRequest(w, r, "credential_id")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleCredentialInGroup(contract, templateID, credentialID)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleCredentialInGroup(contract, templateID, credentialID string) (interface{}, error) {
	switch credential, isActive, err := Credential(hd.database, contract, templateID, credentialID); {
	case err != nil:
		return nil, util.ErrNotFound.WithMessage(err, "credential by contract %s, template %s, id %s", contract, templateID, credentialID)
	case credential == nil:
		return nil, util.ErrNotFound.Errorf("credential by contract %s, template %s, id %s", contract, templateID, credentialID)
	default:
		hal, err := hd.buildCredentialHal(contract, *credential, isActive)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildCredentialHal(
	contract string,
	credential types.Credential,
	isActive bool,
) (cdigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathDIDCredential,
		"contract", contract,
		"template_id", credential.TemplateID(),
		"credential_id", credential.CredentialID(),
	)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(
		struct {
			Credential types.Credential `json:"credential"`
			IsActive   bool             `json:"is_active"`
		}{Credential: credential, IsActive: isActive},
		cdigest.NewHalLink(h, nil),
	)

	return hal, nil
}

func (hd *Handlers) handleCredentials(w http.ResponseWriter, r *http.Request) {
	limit := cdigest.ParseLimitQuery(r.URL.Query().Get("limit"))
	offset := cdigest.ParseStringQuery(r.URL.Query().Get("offset"))
	reverse := cdigest.ParseBoolQuery(r.URL.Query().Get("reverse"))

	cachekey := cdigest.CacheKey(
		r.URL.Path, cdigest.StringOffsetQuery(offset),
		cdigest.StringBoolQuery("reverse", reverse),
	)

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	templateID, err, status := cdigest.ParseRequest(w, r, "template_id")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		i, filled, err := hd.handleCredentialsInGroup(contract, templateID, offset, reverse, limit)

		return []interface{}{i, filled}, err
	})

	if err != nil {
		hd.Log().Err(err).Str("Issuer", contract).Msg("failed to get credentials")
		cdigest.HTTP2HandleError(w, err)

		return
	}

	var b []byte
	var filled bool
	{
		l := v.([]interface{})
		b = l[0].([]byte)
		filled = l[1].(bool)
	}

	cdigest.HTTP2WriteHalBytes(hd.encoder, w, b, http.StatusOK)

	if !shared {
		expire := hd.expireNotFilled
		if len(offset) > 0 && filled {
			expire = time.Minute
		}

		cdigest.HTTP2WriteCache(w, cachekey, expire)
	}
}

func (hd *Handlers) handleCredentialsInGroup(
	contract, templateID string,
	offset string,
	reverse bool,
	l int64,
) ([]byte, bool, error) {
	var limit int64
	if l < 0 {
		limit = hd.itemsLimiter("service-credentials")
	} else {
		limit = l
	}

	var vas []cdigest.Hal
	if err := CredentialsByServiceTemplate(
		hd.database, contract, templateID, reverse, offset, limit,
		func(credential types.Credential, isActive bool, st base.State) (bool, error) {
			hal, err := hd.buildCredentialHal(contract, credential, isActive)
			if err != nil {
				return false, err
			}
			vas = append(vas, hal)

			return true, nil
		},
	); err != nil {
		return nil, false, util.ErrNotFound.WithMessage(err, "credentials by contract %s, template %s", contract, templateID)
	} else if len(vas) < 1 {
		return nil, false, util.ErrNotFound.Errorf("credentials by contract %s, template %s", contract, templateID)
	}

	i, err := hd.buildCredentialsHal(contract, templateID, vas, offset, reverse)
	if err != nil {
		return nil, false, err
	}

	b, err := hd.encoder.Marshal(i)
	return b, int64(len(vas)) == limit, err
}

func (hd *Handlers) buildCredentialsHal(
	contract, templateID string,
	vas []cdigest.Hal,
	offset string,
	reverse bool,
) (cdigest.Hal, error) {
	baseSelf, err := hd.combineURL(
		HandlerPathDIDCredentials,
		"contract", contract,
		"template_id", templateID,
	)
	if err != nil {
		return nil, err
	}

	self := baseSelf
	if len(offset) > 0 {
		self = cdigest.AddQueryValue(baseSelf, cdigest.StringOffsetQuery(offset))
	}
	if reverse {
		self = cdigest.AddQueryValue(baseSelf, cdigest.StringBoolQuery("reverse", reverse))
	}

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(vas, cdigest.NewHalLink(self, nil))

	h, err := hd.combineURL(HandlerPathDIDService, "contract", contract)
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("service", cdigest.NewHalLink(h, nil))

	var nextOffset string

	if len(vas) > 0 {
		va, ok := vas[len(vas)-1].Interface().(struct {
			Credential types.Credential `json:"credential"`
			IsActive   bool             `json:"is_active"`
		})
		if !ok {
			return nil, errors.Errorf("failed to build credentials hal")
		}
		nextOffset = va.Credential.CredentialID()
	}

	if len(nextOffset) > 0 {
		next := baseSelf
		next = cdigest.AddQueryValue(next, cdigest.StringOffsetQuery(nextOffset))

		if reverse {
			next = cdigest.AddQueryValue(next, cdigest.StringBoolQuery("reverse", reverse))
		}

		hal = hal.AddLink("next", cdigest.NewHalLink(next, nil))
	}

	hal = hal.AddLink("reverse", cdigest.NewHalLink(cdigest.AddQueryValue(baseSelf, cdigest.StringBoolQuery("reverse", !reverse)), nil))

	return hal, nil
}

func (hd *Handlers) handleHolderCredential(w http.ResponseWriter, r *http.Request) {
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	holder, err, status := cdigest.ParseRequest(w, r, "holder")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleHolderCredentialsInGroup(contract, holder)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleHolderCredentialsInGroup(contract, holder string) (interface{}, error) {
	var did string
	switch d, err := HolderDID(hd.database, contract, holder); {
	case err != nil:
		return nil, util.ErrNotFound.WithMessage(err, "DID by contract %s, holder %s", contract, holder)
	case d == "":
		return nil, util.ErrNotFound.Errorf("DID by contract %s, holder %s", contract, holder)
	default:
		did = d
	}

	var vas []cdigest.Hal
	if err := CredentialsByServiceHolder(
		hd.database, contract, holder,
		func(credential types.Credential, isActive bool, st base.State) (bool, error) {
			hal, err := hd.buildCredentialHal(contract, credential, isActive)
			if err != nil {
				return false, err
			}
			vas = append(vas, hal)

			return true, nil
		},
	); err != nil {
		return nil, util.ErrNotFound.WithMessage(err, "credentials by contract %s, holder %s", contract, holder)
	} else if len(vas) < 1 {
		return nil, util.ErrNotFound.Errorf("credentials by contract %s, holder %s", contract, holder)
	}
	hal, err := hd.buildHolderDIDCredentialsHal(contract, holder, did, vas)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(hal)
}

func (hd *Handlers) buildHolderDIDCredentialsHal(
	contract, holder, did string,
	vas []cdigest.Hal,
) (cdigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDIDHolder, "contract", contract, "holder", holder)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(
		struct {
			DID         string        `json:"did"`
			Credentials []cdigest.Hal `json:"credentials"`
		}{
			DID:         did,
			Credentials: vas,
		}, cdigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleTemplate(w http.ResponseWriter, r *http.Request) {
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	templateID, err, status := cdigest.ParseRequest(w, r, "template_id")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleTemplateInGroup(contract, templateID)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleTemplateInGroup(contract, templateID string) (interface{}, error) {
	switch template, err := Template(hd.database, contract, templateID); {
	case err != nil:
		return nil, util.ErrNotFound.WithMessage(err, "template by contract %s, template %s", contract, templateID)
	case template == nil:
		return nil, util.ErrNotFound.Errorf("template by contract %s, template %s", contract, templateID)
	default:
		hal, err := hd.buildTemplateHal(contract, templateID, *template)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildTemplateHal(
	contract, templateID string,
	template types.Template,
) (cdigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathDIDTemplate,
		"contract", contract,
		"template_id", templateID,
	)
	if err != nil {
		return nil, err
	}

	hal := cdigest.NewBaseHal(template, cdigest.NewHalLink(h, nil))

	return hal, nil
}
