package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/hibare/DomainHQ/internal/api/handler"
	"github.com/hibare/DomainHQ/internal/config"
	commonMiddleware "github.com/hibare/GoCommon/pkg/http/middleware"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	app App
)

const (
	testAPIKey = "test-key"
)

func setEnv() {
	os.Setenv("DB_USERNAME", "john")
	os.Setenv("DB_PASSWORD", "pwd0123456789")
	os.Setenv("DB_NAME", "domain_hq_test")
	os.Setenv("API_KEYS", testAPIKey)
}

func unsetEnv() {
	os.Unsetenv("DB_USERNAME")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("API_KEYS")
}

func TruncateTables(db *gorm.DB) {
	tables, err := db.Migrator().GetTables()

	if err != nil {
		panic(err)
	}

	for _, table := range tables {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
	}
}

func TestMain(m *testing.M) {
	setEnv()
	config.LoadConfig()
	app.Init()
	TruncateTables(app.DB)
	code := m.Run()
	os.Exit(code)
	unsetEnv()
}

func TestHome(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	app.Router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := "Good to see you"
	assert.Equal(t, expectedBody, w.Body.String())
}

func TestWebFinger(t *testing.T) {
	testCases := []struct {
		Name         string
		URL          string
		ExpectStatus int
	}{
		{
			Name:         "URL without trailing slash (fail)",
			URL:          "/.well-known/webfinger",
			ExpectStatus: http.StatusUnprocessableEntity,
		},
		{
			Name:         "URL with trailing slash (fail)",
			URL:          "/.well-known/webfinger/",
			ExpectStatus: http.StatusUnprocessableEntity,
		},
		{
			Name:         "Domain allowed",
			URL:          "/.well-known/webfinger?resource=acct:test@example.com",
			ExpectStatus: http.StatusOK,
		},
		{
			Name:         "Domain not allowed",
			URL:          "/.well-known/webfinger?resource=acct:test@example1.com",
			ExpectStatus: http.StatusForbidden,
		},
		{
			Name:         "Invalid request",
			URL:          "/.well-known/webfinger?resource=test@example1.com",
			ExpectStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest("GET", tc.URL, nil)
			assert.NoError(t, err)
			app.Router.ServeHTTP(w, r)
			assert.Equal(t, tc.ExpectStatus, w.Code)

			if tc.ExpectStatus == http.StatusOK {
				expectedBody := handler.WebFingerResponse{
					Subject: "acct:test@example.com",
					Links: []handler.Link{
						{
							Rel:  handler.REL,
							Href: config.Current.WebFinger.Resource,
						},
					},
				}
				responseBody := handler.WebFingerResponse{}
				err = json.NewDecoder(w.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.Equal(t, expectedBody, responseBody)
			}
		})
	}
}

func TestGPGKeyAdd(t *testing.T) {
	testCases := []struct {
		Name         string
		URL          string
		ExpectStatus int
		APIKey       string
		KeyText      string
	}{
		{
			Name:         "Successful Add",
			URL:          "/pks/add",
			ExpectStatus: http.StatusOK,
			APIKey:       testAPIKey,
			KeyText:      "-----BEGIN%20PGP%20PUBLIC%20KEY%20BLOCK-----%0A%0AmQINBGSxkwYBEAC3uQxR24dmrn3Xa9R0TRreET4RXssYjwVJdWnmg2YgBliv6Xm2%0AXpKbHnUikjbzA1DbKyKY6GYtuSxCRUanZAFEjtpQMzi%2FcM3CvPPtniTdgzVtGdPN%0AxtQ7EvzL6GgXIYq1DTpu2Tvd6VuZTlPMOyrlCN9ejIITQjbtn3G5fK%2BRHrMN6Eve%0AN0bksVqh2FaKg%2BI%2BmvKegP6SNH1TLe8m9OjxJSOVOBMqZPxDewFpxvqLxyHZpPKs%0AHtlSKK4Q%2Fk%2FYR%2BeHKbYhVJncQchAVBIhtNz%2BFdd5bCFEZeZQjQ2IdTG41mN9tcCZ%0AtoquEGdDsrzLa7nzzB2MjgsSusSxZLtEYcAQrUvtxCDRBLUoDaVdk7jr0%2B3YQ7Un%0AfwDr2HsjbLMniuTTW3N22x%2FWXim1JdXQ7169Y8ADUa4%2BPNHIwz8%2FXkI9f3w0m24D%0AS9nnukKLn54YyPSacw0S6gAQ3JcNfXf3%2BdUpfCKdYSDNHdQUfNvWl3kndyxMTibl%0AXI5qmfua08aVcr2X1MCrG92yGXPSnbLCfqag3b52l9LIO2RdsPjGLCFchd3IzzIE%0A3VIibWI%2FCC9F0s9H5nZZbIfc%2BMcxmNSVug0j7l7Vi3CDSsoMfwxGL976XoqeuW59%0Ab%2BmIElNoSEYz8%2BEbkGTahnxv2vbK0l717XbKBQTcID1XJ1r%2B6mZDEzLJVQARAQAB%0AtCtFeGFtcGxlIChleGFtcGxlIGtleSkgPGV4YW1wbGVAZXhhbXBsZS5jb20%2BiQJO%0ABBMBCgA4FiEEIqN6mnDjllFX4WAH%2FgZrBLRNoNMFAmSxkwYCGwMFCwkIBwIGFQoJ%0ACAsCBBYCAwECHgECF4AACgkQ%2FgZrBLRNoNOYkg%2F%2FVsEnEp6BGgtlu3BHzI6n%2Bvf8%0AzmjpjS8%2FE34SrupeXw7Nzurpl2T8yifUP2LFj5LCA1NV3bUItwqWB87OUvEuB2RM%0AxYbKasw4eQJxy3U9FGOk9iOeUmbBD8DlGw58uBL47ukpKvhj%2BvBt6z7Q5RQPE4Wx%0AfyS9h%2BDsVrAQPfljC1O2IuITs3DuXp3CtGt8ARinkclfV9sdzBxILErEXSiktK16%0A1RCegM9%2F19iRitwD8EK8o7SEMw4vmZyR%2BkENfcKjj2WvF6LWhuCMeP7y0U5rC346%0AWZ9Phchz%2FS5beeNdlbrYicqSk2%2BFyc%2BHd2YmcgIX8utLqrpAzgT85G3nXAstDY%2BX%0AGmrYwG7JNXMJ0rYYrQN3tYmu%2FL8ossgVL7L46HBfuVCpYS1i7iL7EkjyRXW4NLyP%0AG2oLYnTJj%2BdOOfGDdJM1ocQbxauQTZ%2F7ibzbHsna9aRM2Cfcic%2FLUHNmaOL%2FZGlu%0AD%2Fd29IzwocOcYVilNP%2Bch7hXnST%2FpsE1m97M2u3XYAKDCqIeShBFVehUnvmzNktB%0Afqr%2FpsMaWMmarXY4k4KEUruoasM45K0HR2oqM9hY%2BzsNEogcRKHzi0c8OIDXQ01w%0AQReDSinqu67xN8QA9GoOpNT9VB8%2BEG4nTvsBNWsYm3lZBRHJHNzkHIdkLudwjk6l%0A5eNEkXYEGGz8N3%2B7BB%2B5Ag0EZLGTBgEQAMsekxi3WRTsp477Z78qWSjrxZlw0yDP%0A9sTBhsiMXhXp2y0bUKTb4uFVHYgCUJBnMAr%2F74m6s%2Fna5nETmT5hjehHV7Pmw9uP%0Aa128%2F43Jc1Nol6A81J%2BzT3W4zFAsbaOyLS8q5stGaiCnLh30FVGez%2Fcs%2FmZeLk%2F1%0AIVZ%2FV1CZJpwqIh8ca1H1WzaWsxlYgxJLTJMWYcr3JK6tkrcpBzyuBCp%2BQ6cpJepq%0AAedDgrofZXuXzPify1VquBPhGgO9zV%2BZxPDgFaGlAmm0JZ3V02wTNIkKsr1vIzei%0AExmuk7EFqDT89%2By2AbZLdFtKt%2BDkOdljaGdUaoDqGUcxGoFL%2BN77RQGPpRKUsizX%0AUELnylBwHgu6ncvTsn0ouX%2FnALpoYduC7GkvVba3tXuHEJkBH7B%2Fv0cPMTKl%2B0Ep%0AtoXDBiCCJ5O5JoM44DgmKSmhyrDa4GHJpLWR7wkYNryVM17RP3Ukw3rLXfVCYllT%0ABrCvxPN9xwLHeCORiR%2BC1yL9Kn125RiCXyQa7H9APJGgSx%2FmbCeaJesYBTfJwjT4%0ApNO4np6q3CarK%2FIutOfd8duYOuRVkJxisBN0XHY%2BQW2FDASNKwIcEbgwQA7%2BM8RA%0Ai%2FlM07uYcRwbGKSEyGp7ksMRi4Lf%2BuqjKQe%2FeDNAXNSB203Xhm2X5hR3PjpuiBdq%0AKaeIAsxZ3cdxABEBAAGJAjYEGAEKACAWIQQio3qacOOWUVfhYAf%2BBmsEtE2g0wUC%0AZLGTBgIbDAAKCRD%2BBmsEtE2g09x5D%2F4ybLo6Y%2Fpj%2FqZtAzHsL0V5jZyKqBf2M0FV%0Awev3iyoqERveAjgfpzha%2BKTc8Q6sB4d5qPqM%2B57UEGnOVYce3QZEslSwPUOhFaKG%0AqtqCHyGcs%2BhwpVxZZ9vGdLA5aezljiqynhUpoYxhhpw2JUwt1PqOutoPpmJMM2FT%0A3ekEO3ZMRh2eW9CigjWsoqFMuDbkIJ%2Fkwy3NDADX1UqSMaLYIHCstXUqgUm4FXnH%0A2T9lJKBu6tGrpSXd%2ByY2lyG3UIf1hVQ1m4DBEGgLzggpuBFmyfuMmq%2FhL5TLH41E%0AxLnITNINHAlm1TdMi%2BKelxKPvLwnlZRl3I0FgOZqctMVi7ZbZY%2BQeXg4JzhvsbWy%0ALwEpPXIQlCRQs9RMjFFzHR1bMAC3oP7s0lP8%2Bci3bhB4yd6omauZQGGerXlKkeNI%0AGqhAntToQP3OsxFVEj9vw7branRMjhjZcNbW4P4uA7hvAEGIOIcgU48kORez7MX5%0AHoU3qdEoIbJsxjFwz5jv3sR1N4cYhmO%2FPaEg%2Btb2uzgzkBIocG25xw6Mo1sOcpRm%0AHmexwn7h7Su9zrY2%2FQqupkHd9HpnYp6b2%2FKABn7eUIC99tRXQjuvo8LIoldhFUYk%0AkE63SZcnMlSEztUWYZUngX3Dj4eAQc4cZXj62dZtZVP5j%2FnKpzJe2dEAVzrqSyZC%0AKtQIWXTIGw%3D%3D%0A%3DoPyT%0A-----END%20PGP%20PUBLIC%20KEY%20BLOCK-----%0A",
		},
		{
			Name:         "Failed - 401",
			URL:          "/pks/add",
			ExpectStatus: http.StatusUnauthorized,
			APIKey:       "",
			KeyText:      "-----BEGIN%20PGP%20PUBLIC%20KEY%20BLOCK-----%0A%0AmQINBGSxkwYBEAC3uQxR24dmrn3Xa9R0TRreET4RXssYjwVJdWnmg2YgBliv6Xm2%0AXpKbHnUikjbzA1DbKyKY6GYtuSxCRUanZAFEjtpQMzi%2FcM3CvPPtniTdgzVtGdPN%0AxtQ7EvzL6GgXIYq1DTpu2Tvd6VuZTlPMOyrlCN9ejIITQjbtn3G5fK%2BRHrMN6Eve%0AN0bksVqh2FaKg%2BI%2BmvKegP6SNH1TLe8m9OjxJSOVOBMqZPxDewFpxvqLxyHZpPKs%0AHtlSKK4Q%2Fk%2FYR%2BeHKbYhVJncQchAVBIhtNz%2BFdd5bCFEZeZQjQ2IdTG41mN9tcCZ%0AtoquEGdDsrzLa7nzzB2MjgsSusSxZLtEYcAQrUvtxCDRBLUoDaVdk7jr0%2B3YQ7Un%0AfwDr2HsjbLMniuTTW3N22x%2FWXim1JdXQ7169Y8ADUa4%2BPNHIwz8%2FXkI9f3w0m24D%0AS9nnukKLn54YyPSacw0S6gAQ3JcNfXf3%2BdUpfCKdYSDNHdQUfNvWl3kndyxMTibl%0AXI5qmfua08aVcr2X1MCrG92yGXPSnbLCfqag3b52l9LIO2RdsPjGLCFchd3IzzIE%0A3VIibWI%2FCC9F0s9H5nZZbIfc%2BMcxmNSVug0j7l7Vi3CDSsoMfwxGL976XoqeuW59%0Ab%2BmIElNoSEYz8%2BEbkGTahnxv2vbK0l717XbKBQTcID1XJ1r%2B6mZDEzLJVQARAQAB%0AtCtFeGFtcGxlIChleGFtcGxlIGtleSkgPGV4YW1wbGVAZXhhbXBsZS5jb20%2BiQJO%0ABBMBCgA4FiEEIqN6mnDjllFX4WAH%2FgZrBLRNoNMFAmSxkwYCGwMFCwkIBwIGFQoJ%0ACAsCBBYCAwECHgECF4AACgkQ%2FgZrBLRNoNOYkg%2F%2FVsEnEp6BGgtlu3BHzI6n%2Bvf8%0AzmjpjS8%2FE34SrupeXw7Nzurpl2T8yifUP2LFj5LCA1NV3bUItwqWB87OUvEuB2RM%0AxYbKasw4eQJxy3U9FGOk9iOeUmbBD8DlGw58uBL47ukpKvhj%2BvBt6z7Q5RQPE4Wx%0AfyS9h%2BDsVrAQPfljC1O2IuITs3DuXp3CtGt8ARinkclfV9sdzBxILErEXSiktK16%0A1RCegM9%2F19iRitwD8EK8o7SEMw4vmZyR%2BkENfcKjj2WvF6LWhuCMeP7y0U5rC346%0AWZ9Phchz%2FS5beeNdlbrYicqSk2%2BFyc%2BHd2YmcgIX8utLqrpAzgT85G3nXAstDY%2BX%0AGmrYwG7JNXMJ0rYYrQN3tYmu%2FL8ossgVL7L46HBfuVCpYS1i7iL7EkjyRXW4NLyP%0AG2oLYnTJj%2BdOOfGDdJM1ocQbxauQTZ%2F7ibzbHsna9aRM2Cfcic%2FLUHNmaOL%2FZGlu%0AD%2Fd29IzwocOcYVilNP%2Bch7hXnST%2FpsE1m97M2u3XYAKDCqIeShBFVehUnvmzNktB%0Afqr%2FpsMaWMmarXY4k4KEUruoasM45K0HR2oqM9hY%2BzsNEogcRKHzi0c8OIDXQ01w%0AQReDSinqu67xN8QA9GoOpNT9VB8%2BEG4nTvsBNWsYm3lZBRHJHNzkHIdkLudwjk6l%0A5eNEkXYEGGz8N3%2B7BB%2B5Ag0EZLGTBgEQAMsekxi3WRTsp477Z78qWSjrxZlw0yDP%0A9sTBhsiMXhXp2y0bUKTb4uFVHYgCUJBnMAr%2F74m6s%2Fna5nETmT5hjehHV7Pmw9uP%0Aa128%2F43Jc1Nol6A81J%2BzT3W4zFAsbaOyLS8q5stGaiCnLh30FVGez%2Fcs%2FmZeLk%2F1%0AIVZ%2FV1CZJpwqIh8ca1H1WzaWsxlYgxJLTJMWYcr3JK6tkrcpBzyuBCp%2BQ6cpJepq%0AAedDgrofZXuXzPify1VquBPhGgO9zV%2BZxPDgFaGlAmm0JZ3V02wTNIkKsr1vIzei%0AExmuk7EFqDT89%2By2AbZLdFtKt%2BDkOdljaGdUaoDqGUcxGoFL%2BN77RQGPpRKUsizX%0AUELnylBwHgu6ncvTsn0ouX%2FnALpoYduC7GkvVba3tXuHEJkBH7B%2Fv0cPMTKl%2B0Ep%0AtoXDBiCCJ5O5JoM44DgmKSmhyrDa4GHJpLWR7wkYNryVM17RP3Ukw3rLXfVCYllT%0ABrCvxPN9xwLHeCORiR%2BC1yL9Kn125RiCXyQa7H9APJGgSx%2FmbCeaJesYBTfJwjT4%0ApNO4np6q3CarK%2FIutOfd8duYOuRVkJxisBN0XHY%2BQW2FDASNKwIcEbgwQA7%2BM8RA%0Ai%2FlM07uYcRwbGKSEyGp7ksMRi4Lf%2BuqjKQe%2FeDNAXNSB203Xhm2X5hR3PjpuiBdq%0AKaeIAsxZ3cdxABEBAAGJAjYEGAEKACAWIQQio3qacOOWUVfhYAf%2BBmsEtE2g0wUC%0AZLGTBgIbDAAKCRD%2BBmsEtE2g09x5D%2F4ybLo6Y%2Fpj%2FqZtAzHsL0V5jZyKqBf2M0FV%0Awev3iyoqERveAjgfpzha%2BKTc8Q6sB4d5qPqM%2B57UEGnOVYce3QZEslSwPUOhFaKG%0AqtqCHyGcs%2BhwpVxZZ9vGdLA5aezljiqynhUpoYxhhpw2JUwt1PqOutoPpmJMM2FT%0A3ekEO3ZMRh2eW9CigjWsoqFMuDbkIJ%2Fkwy3NDADX1UqSMaLYIHCstXUqgUm4FXnH%0A2T9lJKBu6tGrpSXd%2ByY2lyG3UIf1hVQ1m4DBEGgLzggpuBFmyfuMmq%2FhL5TLH41E%0AxLnITNINHAlm1TdMi%2BKelxKPvLwnlZRl3I0FgOZqctMVi7ZbZY%2BQeXg4JzhvsbWy%0ALwEpPXIQlCRQs9RMjFFzHR1bMAC3oP7s0lP8%2Bci3bhB4yd6omauZQGGerXlKkeNI%0AGqhAntToQP3OsxFVEj9vw7branRMjhjZcNbW4P4uA7hvAEGIOIcgU48kORez7MX5%0AHoU3qdEoIbJsxjFwz5jv3sR1N4cYhmO%2FPaEg%2Btb2uzgzkBIocG25xw6Mo1sOcpRm%0AHmexwn7h7Su9zrY2%2FQqupkHd9HpnYp6b2%2FKABn7eUIC99tRXQjuvo8LIoldhFUYk%0AkE63SZcnMlSEztUWYZUngX3Dj4eAQc4cZXj62dZtZVP5j%2FnKpzJe2dEAVzrqSyZC%0AKtQIWXTIGw%3D%3D%0A%3DoPyT%0A-----END%20PGP%20PUBLIC%20KEY%20BLOCK-----%0A",
		},
		{
			Name:         "Failed - 400",
			URL:          "/pks/add",
			ExpectStatus: http.StatusBadRequest,
			APIKey:       testAPIKey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			payload := strings.NewReader(fmt.Sprintf("keytext=%s", tc.KeyText))

			w := httptest.NewRecorder()
			r, err := http.NewRequest("POST", tc.URL, payload)
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Add(commonMiddleware.AuthHeaderName, tc.APIKey)
			assert.NoError(t, err)
			app.Router.ServeHTTP(w, r)
			assert.Equal(t, tc.ExpectStatus, w.Code)
		})
	}
}

func TestGPGPubKeyLookup(t *testing.T) {
	testCases := []struct {
		Name         string
		URL          string
		ExpectStatus int
		Query        string
	}{
		{
			Name:         "Lookup by fingerprint - 200",
			URL:          "/pks/lookup",
			ExpectStatus: http.StatusOK,
			Query:        "op=get&search=0x22A37A9A70E3965157E16007FE066B04B44DA0D3",
		},
		{
			Name:         "Lookup by email - 200",
			URL:          "/pks/lookup",
			ExpectStatus: http.StatusOK,
			Query:        "op=get&search=example@example.com",
		},
		{
			Name:         "Lookup by Key ID - 200",
			URL:          "/pks/lookup",
			ExpectStatus: http.StatusOK,
			Query:        "op=get&search=FE066B04B44DA0D3",
		},
		{
			Name:         "Lookup by Key ID (Short) - 200",
			URL:          "/pks/lookup",
			ExpectStatus: http.StatusOK,
			Query:        "op=get&search=B44DA0D3",
		},
		{
			Name:         "Invalid op",
			URL:          "/pks/lookup",
			ExpectStatus: http.StatusBadRequest,
			Query:        "op=index&search=B44DA0D3",
		},
		{
			Name:         "No query",
			URL:          "/pks/lookup",
			ExpectStatus: http.StatusUnprocessableEntity,
			Query:        "",
		},
		{
			Name:         "Key not found",
			URL:          "/pks/lookup",
			ExpectStatus: http.StatusNotFound,
			Query:        "op=get&search=example@example.in",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", tc.URL, tc.Query), nil)
			assert.NoError(t, err)
			app.Router.ServeHTTP(w, r)
			assert.Equal(t, tc.ExpectStatus, w.Code)
		})
	}
}
