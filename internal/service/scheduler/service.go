package scheduler

import (
	"context"
	eventQueueJob "go-template/internal/service/scheduler/event_queue"
	"go-template/internal/service/scheduler/session"
	"go-template/pkg/config"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type Service struct {
	dig.In

	Config             config.Config
	EventQueueRelayJob eventQueueJob.RelayJob
	SessionCleanupJob  session.CleanupJob
	Logger             *zap.Logger
}

func (s *Service) Start(ctx context.Context) {
	if !s.Config.Scheduler.Enabled {
		s.Logger.Info("[SchedulerService] Scheduler is disabled in config")
		return
	}

	loc := time.Local
	if s.Config.Scheduler.Timezone != "" {
		if parsedLoc, err := time.LoadLocation(s.Config.Scheduler.Timezone); err == nil {
			loc = parsedLoc
		} else {
			s.Logger.Warn("[SchedulerService] invalid timezone in config, falling back to Local",
				zap.String("timezone", s.Config.Scheduler.Timezone),
				zap.Error(err),
			)
		}
	}

	cronEngine := cron.New(
		cron.WithLocation(loc),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
	)

	// 1. Register event_queue_relay job from YAML
	if jobCfg, ok := s.Config.Scheduler.Jobs["event_queue_relay"]; ok && jobCfg.Enabled {
		_, err := cronEngine.AddFunc(jobCfg.Cron, func() {
			s.EventQueueRelayJob.ProcessBatch(ctx, jobCfg.BatchSize)
		})
		if err != nil {
			s.Logger.Error("[SchedulerService] failed to schedule event_queue_relay job", zap.String("cron", jobCfg.Cron), zap.Error(err))
		} else {
			s.Logger.Info("[SchedulerService] registered job: event_queue_relay", zap.String("schedule", jobCfg.Cron))
		}
	}

	// 2. Register cleanup_expired_sessions job from YAML
	if jobCfg, ok := s.Config.Scheduler.Jobs["cleanup_expired_sessions"]; ok && jobCfg.Enabled {
		_, err := cronEngine.AddFunc(jobCfg.Cron, func() {
			s.SessionCleanupJob.Run(ctx)
		})
		if err != nil {
			s.Logger.Error("[SchedulerService] failed to schedule cleanup_expired_sessions job", zap.String("cron", jobCfg.Cron), zap.Error(err))
		} else {
			s.Logger.Info("[SchedulerService] registered job: cleanup_expired_sessions", zap.String("schedule", jobCfg.Cron))
		}
	}

	cronEngine.Start()
	s.Logger.Info("[SchedulerService] Cron scheduler started", zap.String("timezone", loc.String()))

	go func() {
		<-ctx.Done()
		s.Logger.Info("[SchedulerService] Stopping cron scheduler gracefully...")
		stopCtx := cronEngine.Stop()
		<-stopCtx.Done()
		s.Logger.Info("[SchedulerService] Cron scheduler stopped cleanly")
	}()
}
