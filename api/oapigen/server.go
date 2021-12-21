// Package oapigen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.3 DO NOT EDIT.
package oapigen

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Triggers a dryrun for a provided task
	// (POST /v1/dryrun_tasks)
	CreateDryRunTask(w http.ResponseWriter, r *http.Request)
	// Creates a new task
	// (POST /v1/tasks)
	CreateTask(w http.ResponseWriter, r *http.Request, params CreateTaskParams)
	// Deletes a task by name
	// (DELETE /v1/tasks/{name})
	DeleteTaskByName(w http.ResponseWriter, r *http.Request, name string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// CreateDryRunTask operation middleware
func (siw *ServerInterfaceWrapper) CreateDryRunTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateDryRunTask(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// CreateTask operation middleware
func (siw *ServerInterfaceWrapper) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params CreateTaskParams

	// ------------- Optional query parameter "run" -------------
	if paramValue := r.URL.Query().Get("run"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "run", r.URL.Query(), &params.Run)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter run: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateTask(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// DeleteTaskByName operation middleware
func (siw *ServerInterfaceWrapper) DeleteTaskByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameter("simple", false, "name", chi.URLParam(r, "name"), &name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter name: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteTaskByName(w, r, name)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL     string
	BaseRouter  chi.Router
	Middlewares []MiddlewareFunc
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/dryrun_tasks", wrapper.CreateDryRunTask)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/tasks", wrapper.CreateTask)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/v1/tasks/{name}", wrapper.DeleteTaskByName)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"AP//KO8e6QgfAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
