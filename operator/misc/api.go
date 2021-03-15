package misc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	v1Core "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

var LISTEN_ADDRESS = GetEnvValue("INVENTA_OPERATOR_LISTEN_ADDRESS", "127.0.0.1")
var AUTH_TENANT_ID = GetEnvValue("INVENTA_OPERATOR_AUTH_TENANT_ID", "-1")
var AUTH_CLIENT_ID = GetEnvValue("INVENTA_OPERATOR_AUTH_CLIENT_ID", "-1")


func InitApi(store *InMemoryStore, enableAuth bool) {
	var provider *oidc.Provider
	if enableAuth {
		newProvider, err := oidc.NewProvider(context.Background(), fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", AUTH_TENANT_ID))
		if err != nil {
			log.Fatal(err)
		}
		provider = newProvider
	}


	authMiddleware := authenticationMiddleware{
		ClientID: AUTH_CLIENT_ID,
		Provider: provider,
	}

	addr := fmt.Sprintf("%s:8090", LISTEN_ADDRESS)

	r := mux.NewRouter()
	app := App{
		Router: r,
		Store:  store,
	}
	if enableAuth {
		r.Handle("/api/get-all", authMiddleware.Middleware(http.HandlerFunc(app.GetAll)))
	} else {
		r.Handle("/api/get-all", http.HandlerFunc(app.GetAll))
	}

	fmt.Printf("HTTP server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, handlers.LoggingHandler(os.Stdout, r)); err != nil {
		log.Fatal(err)
	}
}

type App struct {
	Router *mux.Router
	Store  *InMemoryStore
}

func (a *App) GetAll(w http.ResponseWriter, r *http.Request) {
	payload, err := json.Marshal(StoreToGetAllResponse(a.Store))
	if err != nil {
		log.Println("Unable to serialise InMemoryStore")
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(payload)
}

type GetAllResponse struct {
	Ingress []IngressDto
	Service []ServiceDto
}

func StoreToGetAllResponse(store *InMemoryStore) GetAllResponse {
	payload := GetAllResponse{
		Ingress: []IngressDto{},
		Service: []ServiceDto{},
	}

	for _, v := range store.Services {
		val := ServiceDto{
			Kind:       v.Kind,
			ApiVersion: v.APIVersion,
			Metadata:   v.ObjectMeta,
			Spec:       v.Spec,
			Status:     v.Status,
		}

		val.Metadata.ManagedFields = []v1.ManagedFieldsEntry{}
		delete(val.Metadata.Annotations, "kubectl.kubernetes.io/last-applied-configuration")

		payload.Service = append(payload.Service, val)
	}

	for _, v := range store.Ingresses {
		val := IngressDto{
			Kind:       v.Kind,
			ApiVersion: v.APIVersion,
			Metadata:   v.ObjectMeta,
			Spec:       v.Spec,
			Status:     v.Status,
		}

		val.Metadata.ManagedFields = []v1.ManagedFieldsEntry{}
		delete(val.Metadata.Annotations, "kubectl.kubernetes.io/last-applied-configuration")

		payload.Ingress = append(payload.Ingress, val)
	}

	return payload
}

type IngressDto struct {
	Kind       string
	ApiVersion string
	Metadata   v1.ObjectMeta
	Spec       v1beta1.IngressSpec
	Status     v1beta1.IngressStatus
}

type ServiceDto struct {
	Kind       string
	ApiVersion string
	Metadata   v1.ObjectMeta
	Spec       v1Core.ServiceSpec
	Status     v1Core.ServiceStatus
}

type authenticationMiddleware struct {
	ClientID string
	Provider *oidc.Provider
}

func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	var verifier = amw.Provider.Verifier(&oidc.Config{ClientID: amw.ClientID})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization") //Authorization: Bearer a7ydfs87afasd8f990
		splitToken := strings.Split(reqToken, "Bearer")
		if len(splitToken) != 2 {
			http.Error(w, "Token doesn't seem right", http.StatusUnauthorized)
			return
		}

		reqToken = strings.TrimSpace(splitToken[1])

		idToken, err := verifier.Verify(r.Context(), reqToken)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unable to verify token", http.StatusUnauthorized)
			return
		}

		var claims struct {
			Emails []string `json:"emails"`
		}
		if err := idToken.Claims(&claims); err != nil {
			fmt.Println(err)
			http.Error(w, "Unable to retrieve claims", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
