package box

import (
	"context"
	"fmt"
	"net/http"

	"github.com/snowmetas/cafe-go/log"
	"github.com/snowmetas/cafe-go/utils"
	//nolint:gosec
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application interface {
	Name() string
	Run()
	AddBundle(bundles ...Bundle)
}

type BaseApplication struct {
	Container

	name string
	ctx  context.Context

	// config
	config appConfig

	// lifecycle hooks
	beforeStart []func(ctx context.Context) error
	afterStart  []func(ctx context.Context) error
	beforeStop  []func(ctx context.Context) error
	afterStop   []func(ctx context.Context) error
}

type appContextKeyType string

const appContextKey appContextKeyType = "app-context"

func NewApplication(opts ...Option) *BaseApplication {
	defaults := getDefaults()
	customOptions := &options{}
	for _, opt := range opts {
		opt(customOptions)
	}

	ctx := utils.DerefCtx(customOptions.ctx, context.Background())
	appName := utils.DerefString(customOptions.appName, defaults.appName)

	app := &BaseApplication{
		Container: *New(),
		name:      appName,
		config: appConfig{
			warnMetric:   defaultWarnLogMetric(appName),
			errorMetric:  defaultErrorLogMetric(appName),
			profilerPort: customOptions.profilerPort,
			enableConfig: customOptions.withConfig,
			includePaths: customOptions.includePaths,
		},
		ctx: ctx,

		// hooks
		beforeStart: customOptions.beforeStart,
		afterStart:  customOptions.afterStart,
		beforeStop:  customOptions.beforeStop,
		afterStop:   customOptions.afterStop,
	}

	app.initLog()

	return app
}

func AppFromContext(ctx context.Context) Application {
	v := ctx.Value(appContextKey)
	if v != nil {
		if app, ok := v.(Application); ok {
			return app
		}
	}
	log.Error(ctx, "Not a valid cafe context")
	return nil
}

func (app *BaseApplication) Name() string {
	return app.name
}

func (app *BaseApplication) runBeforeStart() {
	if err := RunUntilError(app.ctx, app.beforeStart); err != nil {
		log.Error(app.ctx, "Error run before start hook:", err)
	}
}

func (app *BaseApplication) runAfterStart() {
	if err := RunUntilError(app.ctx, app.afterStart); err != nil {
		log.Error(app.ctx, "Error run after start hook:", err)
	}
}

func (app *BaseApplication) runBeforeStop() {
	if err := RunUntilError(app.ctx, app.beforeStop); err != nil {
		log.Error(app.ctx, "Error run before stopAll hook:", err)
	}
}

func (app *BaseApplication) runAfterStop() {
	if err := RunUntilError(app.ctx, app.afterStop); err != nil {
		log.Error(app.ctx, "Error run after stopAll hook:", err)
	}
}

func (app *BaseApplication) initLog() {
	//if zae.IsDevelopEnv() {
	log.SetLevel(log.DebugLevel)
	//}
}

func (app *BaseApplication) Run() {
	log.Infof("Run cafe application,name=%s", app.name)

	// 这个时候才知道应用的 application 对象是什么
	app.ctx = context.WithValue(app.ctx, appContextKey, app)

	if app.config.profilerPort != nil {
		go func() {
			err := http.ListenAndServe(fmt.Sprintf(":%d", *app.config.profilerPort), nil)
			if err != nil {
				log.Error(app.ctx, "Start pprof error: ", err)
			}
		}()
	}

	// start all
	app.runBeforeStart()

	finishCtx := app.StartAll(app.ctx)

	app.runAfterStart()

	// wait for shutdown or done
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	select {
	case <-finishCtx.Done():
		log.Info(app.ctx, "All bundle finished!")
	case <-shutdownSignal:
		log.Info(app.ctx, "Shutdown signal received")
	}

	// stop all
	app.runBeforeStop()

	ctx := app.StopAll(app.ctx)

	shutdownTimeout := time.After(30 * time.Second)
	select {
	case <-ctx.Done():
		log.Info(app.ctx, "Application stopped")
	case <-shutdownTimeout:
		log.Info(app.ctx, "Shutdown timeout, force stop application")
	}

	app.runAfterStop()

	log.Info(app.ctx, "Bye!")
}
