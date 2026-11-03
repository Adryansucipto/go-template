package scheduler

import (
	"context"
	"go-template/internal/controller/scheduler/outbox"
	"go-template/internal/controller/scheduler/session"
	"go-template/pkg/config"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type Controller struct {
	dig.In

	Config            config.Config
	OutboxRelayJob    outbox.RelayJob
	SessionCleanupJob session.CleanupJob
	Logger            *zap.Logger
}

func (c *Controller) Start(ctx context.Context) {
	if !c.Config.Scheduler.Enabled {
		c.Logger.Info("[Scheduler] Scheduler is disabled in config")
		return
	}

	loc := time.Local
	if c.Config.Scheduler.Timezone != "" {
		if parsedLoc, err := time.LoadLocation(c.Config.Scheduler.Timezone); err == nil {
			loc = parsedLoc
		} else {
			c.Logger.Warn("[Scheduler] invalid timezone in config, falling back to Local",
				zap.String("timezone", c.Config.Scheduler.Timezone),
				zap.Error(err),
			)
		}
	}

	// Support both standard 5-field cron ("*/5 * * * *") and 6-field with seconds ("*/5 * * * * *" or "@every 5s")
	cronEngine := cron.New(
		cron.WithLocation(loc),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
	)

	// 1. Register outbox_relay job from YAML
	if jobCfg, ok := c.Config.Scheduler.Jobs["outbox_relay"]; ok && jobCfg.Enabled {
		_, err := cronEngine.AddFunc(jobCfg.Cron, func() {
			c.OutboxRelayJob.ProcessBatch(ctx, jobCfg.BatchSize)
		})
		if err != nil {
			c.Logger.Error("[Scheduler] failed to schedule outbox_relay job", zap.String("cron", jobCfg.Cron), zap.Error(err))
		} else {
			c.Logger.Info("[Scheduler] registered job: outbox_relay", zap.String("schedule", jobCfg.Cron))
		}
	}

	// 2. Register cleanup_expired_sessions job from YAML
	if jobCfg, ok := c.Config.Scheduler.Jobs["cleanup_expired_sessions"]; ok && jobCfg.Enabled {
		_, err := cronEngine.AddFunc(jobCfg.Cron, func() {
			c.SessionCleanupJob.Run(ctx)
		})
		if err != nil {
			c.Logger.Error("[Scheduler] failed to schedule cleanup_expired_sessions job", zap.String("cron", jobCfg.Cron), zap.Error(err))
		} else {
			c.Logger.Info("[Scheduler] registered job: cleanup_expired_sessions", zap.String("schedule", jobCfg.Cron))
		}
	}

	cronEngine.Start()
	c.Logger.Info("[Scheduler] Cron scheduler started", zap.String("timezone", loc.String()))

	go func() {
		<-ctx.Done()
		c.Logger.Info("[Scheduler] Stopping cron scheduler gracefully...")
		stopCtx := cronEngine.Stop()
		<-stopCtx.Done()
		c.Logger.Info("[Scheduler] Cron scheduler stopped cleanly")
	}()
}

