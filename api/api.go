package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"

	"github.com/hashicorp/consul-terraform-sync/driver"
	"github.com/hashicorp/consul-terraform-sync/event"
	"github.com/hashicorp/consul-terraform-sync/logging"
)

const (
	defaultAPIVersion = "v1"

	// StatusSuccessful is the successful status. This is determined based on status
	// type.
	//
	// Task Status: Determined by the success of a task updating. The 5 most
	// recent task updates are stored as an ‘event’ in CTS. A task is successful
	// when the most recent stored event is successful.
	StatusSuccessful = "successful"

	// StatusErrored is the errored status. This is determined based on status
	// type.
	//
	// Task Status: Determined by the success of a task updating. The 5 most
	// recent task updates are stored as an ‘event’ in CTS. A task is errored
	// when the most recent stored event is not successful but all prior stored
	// events are successful.
	StatusErrored = "errored"

	// StatusCritical is the critical status. This is determined based on status
	// type.
	//
	// Task Status: Determined by the success of a task updating. The 5 most
	// recent task updates are stored as an ‘event’ in CTS. A task is critical
	// when the most recent stored event is not successful and at least one prior
	// stored event is all not successful.
	StatusCritical = "critical"

	// StatusUnknown is when the status is unknown. This is determined
	// based on status type.
	//
	// Task Status: Determined by the success of a task updating. The 5 most
	// recent task updates are stored as an ‘event’ in CTS. A task is
	// unknown when no event data has been collected yet.
	StatusUnknown = "unknown"

	logSystemName = "api"
)

// API supports api requests to the cts binary
type API struct {
	store   *event.Store
	drivers *driver.Drivers
	port    int
	version string
	srv     *http.Server
}

// NewAPI create a new API object
func NewAPI(store *event.Store, drivers *driver.Drivers, port int) *API {
	//mux := http.NewServeMux()
	logger := logging.Global().Named(logSystemName)
	r := chi.NewRouter()

	// add the middleware, must be done first
	r.Use(withLogging)

	// retrieve overall status
	r.Mount(fmt.Sprintf("/%s/%s", defaultAPIVersion, overallStatusPath),
		newOverallStatusHandler(store, drivers, defaultAPIVersion))

	// retrieve all task statuses
	r.Mount(fmt.Sprintf("/%s/%s", defaultAPIVersion, taskStatusPath),
		newTaskStatusHandler(store, drivers, defaultAPIVersion))

	// crud task
	r.Mount(fmt.Sprintf("/%s/%s/", defaultAPIVersion, taskPath),
		newTaskHandler(store, drivers, defaultAPIVersion))

	// only want swagger validation for the openapi generated endpoints
	r.Group(func(r chi.Router) {
		swagger, err := GetSwagger()
		if err != nil {
			logger.Error("Error loading swagger spec", "error", err)
			panic("this should never error")
		}

		// Clear out the servers array in the swagger spec, that skips validating
		// that server names match. We don't know how this thing will be run.
		swagger.Servers = nil

		// Use our validation middleware to check all requests against the
		// OpenAPI schema.
		r.Use(middleware.OapiRequestValidator(swagger))

		taskHandler := NewTaskHandler()
		HandlerFromMux(taskHandler, r)
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	return &API{
		port:    port,
		drivers: drivers,
		store:   store,
		version: defaultAPIVersion,
		srv:     srv,
	}
}

// Serve starts up and handles shutdown for the http server to serve
// API requests
func (api *API) Serve(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(1)

	logger := logging.Global().Named(logSystemName)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := api.srv.Shutdown(ctxShutDown); err != nil {
					logger.Error("error stopping api server", "error", err)
				} else {
					logger.Info("shutdown api server")
				}
				return
			}
		}
	}()

	logger.Info("starting server", "port", api.port)
	if err := api.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("error serving api", "port", api.port, "error", err)
		return err
	}

	// wait for shutdown
	wg.Wait()
	return ctx.Err()
}

// jsonResponse adds the return response for handlers. Returns if json encode
// errored. Option to check error or add responses to jsonResponse test to
// test json encoding
func jsonResponse(w http.ResponseWriter, code int, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(response)
}

func jsonErrorResponse(ctx context.Context, w http.ResponseWriter, code int, err error) {
	err = jsonResponse(w, code, NewErrorResponse(err))
	if err != nil {
		logging.FromContext(ctx).Named(logSystemName).Error("error, could not generate json error response",
			"error", err)
	}
}

// getTaskName retrieves the taskname from the url. Returns empty string if no
// taskname is specified
func getTaskName(reqPath, apiPath, version string) (string, error) {
	taskPathNoID := fmt.Sprintf("/%s/%s", version, apiPath)
	if reqPath == taskPathNoID {
		return "", nil
	}

	taskPathWithID := taskPathNoID + "/"
	taskName := strings.TrimPrefix(reqPath, taskPathWithID)
	if invalid := strings.ContainsRune(taskName, '/'); invalid {
		return "", fmt.Errorf("unsupported path '%s'. request must be format "+
			"'path/{task-name}'. task name cannot have '/ ' and api "+
			"does not support further resources", reqPath)
	}

	return taskName, nil
}
