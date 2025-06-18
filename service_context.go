package sctx

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	DevEnv     = "dev"
	ProdEnv    = "prod"
	StagingEnv = "staging"
)

type Option func(*serviceContext)

type Component interface {
	ID() string
	InitFlags()
	Activate(ServiceContext) error
	Stop() error
}

type ServiceContext interface {
	Logger(prefix string) Logger
}

type serviceContext struct {
	name       string
	env        string
	components []Component
	store      map[string]Component
	cmdLine    *AppFlagSet
	logger     Logger
}

func NewServiceContext(opts ...Option) ServiceContext {
	sv := &serviceContext{
		store: make(map[string]Component),
	}

	sv.components = []Component{defaultLogger}

	for _, opt := range opts {
		opt(sv)
	}

	sv.initFlags()

	sv.cmdLine = newFlagSet(sv.name, flag.CommandLine)

	sv.parseFlags()

	sv.logger = defaultLogger.GetLogger("serviceContext")

	return sv
}

func (s *serviceContext) initFlags() {
	flag.StringVar(&s.env, "app-env", DevEnv, "Env for service. Ex: dev | stg | prd")

	for _, c := range s.components {
		c.InitFlags()
	}
}

func (s *serviceContext) Get(id string) (interface{}, bool) {
	c, ok := s.store[id]
	if !ok {
		return nil, false
	}

	return c, true
}

func (s *serviceContext) MustGet(id string) interface{} {
	c, ok := s.Get(id)

	if !ok {
		panic(fmt.Sprintf("can not get: %s\n", id))
	}
	return c
}

func (s *serviceContext) Load() error {
	s.logger.Info("Service context is loading ...")

	for _, c := range s.components {
		if err := c.Activate(s); err != nil {
			return err
		}
	}

	return nil
}

func (s *serviceContext) Logger(prefix string) Logger {
	return defaultLogger.GetLogger(prefix)
}

func (s *serviceContext) Stop() error {
	s.logger.Infoln("Stopping service context")

	for i := range s.components {
		if err := s.components[i].Stop(); err != nil {
			return err
		}
	}

	s.logger.Infoln("service context stopped")

	return nil
}

func (s *serviceContext) GetName() string { return s.name }
func (s *serviceContext) EnvName() string { return s.env }
func (s *serviceContext) OutEnv()         { s.cmdLine.GetSampleEnvs() }

func WithName(name string) Option {
	return func(s *serviceContext) { s.name = name }
}

func WithComponent(c Component) Option {
	return func(s *serviceContext) {
		if _, ok := s.store[c.ID()]; !ok {
			s.components = append(s.components, c)
			s.store[c.ID()] = c
		}
	}
}

func (s *serviceContext) parseFlags() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	_, err := os.Stat(envFile)
	if err == nil {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Fatalf("Loading env(%s): %s", envFile, err.Error())
		}
	} else if envFile != ".env" {
		log.Fatalf("Loading env(%s): %s", envFile, err.Error())
	}

	s.cmdLine.Parse([]string{})
}
