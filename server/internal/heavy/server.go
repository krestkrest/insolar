//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package heavy

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/server/internal"
	"github.com/insolar/insolar/version"
)

type Server struct {
	cfgPath        string
	genesisCfgPath string
}

func New(cfgPath string, genesisCfgPath string) *Server {
	return &Server{
		cfgPath:        cfgPath,
		genesisCfgPath: genesisCfgPath,
	}
}

func (s *Server) Serve() {
	cfgHolder := configuration.NewHolder()
	var err error
	if len(s.cfgPath) != 0 {
		err = cfgHolder.LoadFromFile(s.cfgPath)
	} else {
		err = cfgHolder.Load()
	}
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err.Error())
	}

	b, err := ioutil.ReadFile(s.genesisCfgPath)
	if err != nil {
		log.Fatalf("failed to load genesis configuration from file: %v", s.genesisCfgPath)
	}
	var genesisCfg application.GenesisHeavyConfig
	err = json.Unmarshal(b, &genesisCfg)
	if err != nil {
		log.Fatalf("failed to pares genesis configuration from file: %v", s.genesisCfgPath)
	}

	cfg := &cfgHolder.Configuration

	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))

	ctx := context.Background()
	traceID := utils.RandTraceID() + "_main"
	ctx, inslog := inslogger.InitNodeLogger(ctx, cfg.Log, traceID, "", "")
	log.InitTicker()

	cmp, err := newComponents(ctx, *cfg, genesisCfg)
	fatal(ctx, err, "failed to create components")

	if cfg.Tracer.Jaeger.AgentEndpoint != "" {
		jaegerFlush := internal.Jaeger(ctx, cfg.Tracer.Jaeger, traceID, cmp.NodeRef, cmp.NodeRole)
		defer jaegerFlush()
	}

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	var waitChannel = make(chan bool)

	go func() {
		sig := <-gracefulStop
		inslog.Debug("caught sig: ", sig)

		inslog.Warn("GRACEFUL STOP APP")
		err = cmp.Stop(ctx)
		fatal(ctx, err, "failed to graceful stop components")
		close(waitChannel)
	}()

	err = cmp.Start(ctx)
	fatal(ctx, err, "failed to start components")
	fmt.Println("All components were started")
	<-waitChannel
}

func fatal(ctx context.Context, err error, message string) {
	if err == nil {
		return
	}
	inslogger.FromContext(ctx).Fatalf("%v: %v", message, err.Error())
}
