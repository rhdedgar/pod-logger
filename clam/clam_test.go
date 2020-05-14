package clam_test

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/rhdedgar/pod-logger/clam"
	"github.com/rhdedgar/pod-logger/config"
)

type clientJSON struct {
	clientMessage string
}

var _ = Describe("Clam", func() {
	config.AppSecrets.OAPIURL = "http://localhost:8080"
	config.AppSecrets.OAPIToken = "exampletdapitoken"
	config.AppSecrets.TDAPIURL = "http://localhost:8080/api/url/"
	config.AppSecrets.TDAPIUser = "exampletdapiuser"
	config.AppSecrets.TDAPIToken = "exampletdapitoken"

	var (
		testNS   = "namespacetodelete"
		testUser = "usertoban"
		testBan  = "test_ban_category"
		e        = echo.New()
	)

	BeforeEach(func() {
		go func() {
			e = echo.New()

			e.Use(middleware.Logger())
			e.Use(middleware.Recover())

			e.DELETE("/apis/project.openshift.io/v1/projects/:namespace", deleteNS)
			e.PUT("/api/url/:username/ban", putBanUser)

			e.Use(middleware.Logger())

			e.Logger.Info(e.Start(":8080"))
		}()
	})

	AfterEach(func() {
		//e.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	})

	Describe("DeleteNS", func() {
		Context("Valid JSON returned from DeleteNS API", func() {
			It("Should return HTTP status code 200", func() {
				status, err := DeleteNS(testNS)

				Expect(err).To(BeNil())
				Expect(status).To(Equal(200))
			})
		})
	})

	Describe("BanUser", func() {
		Context("Valid JSON returned from BanUser API", func() {
			It("Should return HTTP status code 200", func() {
				status, err := BanUser(testUser, testBan)

				Expect(err).To(BeNil())
				Expect(status).To(Equal(200))
			})
		})
	})
})

// postDeleteNS mimicks a namespace deletion request response from the OpenShift API. It's accessed with:
// DELETE /apis/project.openshift.io/v1/projects/:namespace
func deleteNS(c echo.Context) error {
	j := &clientJSON{
		clientMessage: "deleteNS JSON returned successfully.",
	}

	return c.JSON(http.StatusOK, j)
}

// putBanUser mimicks a ban request response from the takedown API. It's accessed with:
// PUT /api/url/:username/ban
func putBanUser(c echo.Context) error {
	j := &clientJSON{
		clientMessage: "putBanUser JSON returned successfully.",
	}

	return c.JSON(http.StatusOK, j)
}
