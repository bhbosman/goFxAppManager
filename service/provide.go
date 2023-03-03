package service

import (
	"context"
	"encoding/json"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type FxManagerState struct {
	Instances []InstateState `json:"instances,omitempty"`
}

type InstateState struct {
	Name   string `json:"name,omitempty"`
	Active bool   `json:"active,omitempty"`
}

func InvokeFxManagerStartAll() fx.Option {
	return fx.Options(
		fx.Invoke(
			func(
				params struct {
					fx.In
					Lifecycle           fx.Lifecycle
					FxManagerService    IFxManagerService
					ConfigurationFolder string `name:"ConfigurationFolder"`
					Logger              *zap.Logger
				},
			) error {
				fileName := filepath.Join(params.ConfigurationFolder, "FxManagerState.json")
				hook := fx.Hook{
					OnStart: func(ctx context.Context) error {
						file, fileErr := os.Open(fileName)
						if fileErr == nil {
							decoder := json.NewDecoder(file)
							fxManagerState := FxManagerState{}
							decodeErr := decoder.Decode(&fxManagerState)
							if decodeErr == nil {
								for _, instance := range fxManagerState.Instances {
									if instance.Active {
										err := params.FxManagerService.Start(ctx, instance.Name)
										if err != nil {
											params.Logger.Error("Issue starting service",
												zap.String("name", instance.Name),
												zap.Error(err))
										}
									}
								}
							}
						}
						return nil
					},
					OnStop: func(ctx context.Context) error {
						state, err := params.FxManagerService.GetState()
						if err != nil {
							return err
						}
						file, fileErr := os.Create(fileName)
						if fileErr != nil {
							params.Logger.Error("error creating file",
								zap.Error(err),
								zap.String("fileName", fileName))
						} else {
							defer func() {
								_ = file.Close()
							}()
							fxManagerState := FxManagerState{
								Instances: make([]InstateState, len(state)),
							}
							for i, s := range state {
								fxManagerState.Instances[i] = InstateState{
									Name:   s,
									Active: true,
								}
							}
							encoder := json.NewEncoder(file)
							encoder.SetIndent("", "\t")
							err = encoder.Encode(&fxManagerState)
							params.Logger.Error("error in encoding", zap.Error(err))

							for _, instance := range fxManagerState.Instances {
								if instance.Active {
									err := params.FxManagerService.Stop(ctx, instance.Name)
									if err != nil {
										params.Logger.Error("Issue stopping service",
											zap.String("name", instance.Name),
											zap.Error(err))
									}
								}
							}
						}

						return nil
					},
				}
				params.Lifecycle.Append(hook)
				return nil
			},
		),
	)
}

func ProvideFxManager() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Name:  "",
				Group: "",
				Target: func(
					params struct {
						fx.In
						OnData             OnDataCallback
						ApplicationContext context.Context `name:"Application"`
						PubSub             *pubsub.PubSub  `name:"Application"`
						Logger             *zap.Logger
						GoFunctionCounter  GoFunctionCounter.IService
					},
				) (IFxManagerService, error) {
					return newService(
						params.ApplicationContext,
						params.OnData,
						params.Logger,
						params.GoFunctionCounter,
					)
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						ApplicationName string                       `name:"ApplicationName"`
						PubSub          *pubsub.PubSub               `name:"Application"`
						FnApps          []messages.CreateAppCallback `group:"Apps"`
						Logger          *zap.Logger
					},
				) OnDataCallback {
					return func(applicationContext context.Context) (IFxManagerData, error) {
						return newData(
							applicationContext,
							params.FnApps,
							params.PubSub,
							params.Logger.Named("FxServiceData"),
						)
					}
				},
			},
		),
		fx.Invoke(
			func(
				params struct {
					fx.In
					Lifecycle        fx.Lifecycle
					FxManagerService IFxManagerService
				},
			) error {
				hook := fx.Hook{
					OnStart: func(ctx context.Context) error {
						return params.FxManagerService.OnStart(ctx)
					},
					OnStop: params.FxManagerService.OnStop,
				}
				params.Lifecycle.Append(hook)
				return nil
			},
		),
	)
}
