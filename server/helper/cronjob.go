package helper

import (
	"context"
	"graded-challenge-2-server/service"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CronJob struct {
	serverService service.ServerService
}

func NewCronJob(serverService service.ServerService) *CronJob {
	return &CronJob{serverService}
}

func (cj *CronJob) UpdateBookStatus() {
	cronJob := cron.New()

	cronJob.AddFunc("@every 5s", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := cj.serverService.UpdateBookStatus(ctx, &emptypb.Empty{})
		if err != nil {
			log.Println("Error updating books status: ", err)
		}

		if data.Counter <= 0 {
			log.Println("No books updated. All books are up to date.")
		} else {
			log.Printf("%d books updated to late.\n", data.Counter)
		}
	})

	cronJob.Start()
}
