package job

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/auho/go-toolkit/flow/action/work"
	"github.com/auho/go-toolkit/flow/flow"
	"github.com/auho/go-toolkit/flow/storage/redis/source"
	"github.com/auho/go-toolkit/flow/task"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type jobConfig struct {
	addr        string
	auth        string
	amount      int64
	pageSize    int64
	poolSize    int
	concurrency int
}

type Job struct {
}

func newJob() *Job {
	return &Job{}
}

func (j *Job) Command(workers ...Worker[string]) *cobra.Command {
	return j.build(workers)
}

func (j *Job) Run(workers ...Worker[string]) {
	cmd := j.build(workers)
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func (j *Job) run(jc *jobConfig, workers []Worker[string]) error {
	fmt.Println("please enter password for redis: ")
	_passwordByte, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		_passwordByte = []byte{}
	}
	jc.auth = string(_passwordByte)

	fmt.Println("addr: ", jc.addr)
	fmt.Println("auth: ", strings.Repeat("*", len(jc.auth)))
	fmt.Println("amount: ", jc.amount)
	fmt.Println("page size: ", jc.pageSize)
	fmt.Println("pool size: ", jc.poolSize)
	fmt.Println("worker concurrency: ", jc.concurrency)

	s, err := source.NewScan(source.Config{
		Amount:   jc.amount,
		PageSize: jc.pageSize,
		Options: &redis.Options{
			Network:  "tcp",
			Addr:     jc.addr,
			Password: jc.auth,
			PoolSize: jc.poolSize,
		},
	})

	if err != nil {
		return err
	}

	var opts []flow.Option[string]
	opts = append(opts, flow.WithSource[string](s))

	timeMark := time.Now().Format("20060102-150405")
	for _, worker := range workers {
		worker.Init(task.WithConcurrency(jc.concurrency))
		worker.WithRedisSource(s.GetClient())
		worker.WithTimeMark(timeMark)
		opts = append(opts, flow.WithActor[string](work.NewActor[string](worker)))
	}

	err = flow.RunFlow[string](opts...)
	if err != nil {
		return err
	}

	return nil
}

func (j *Job) build(workers []Worker[string]) *cobra.Command {
	var jc *jobConfig
	var cmd = &cobra.Command{
		Use: "redis-job",
		RunE: func(cmd *cobra.Command, args []string) error {
			return j.run(jc, workers)
		},
	}

	cmd.Flags().StringVar(&jc.addr, "addr", "", "addr of redis")
	cmd.Flags().Int64Var(&jc.amount, "amount", 0, "amount of job")
	cmd.Flags().Int64Var(&jc.pageSize, "page", 100, "page size of redis scan")
	cmd.Flags().IntVar(&jc.poolSize, "pool", 10, "pool size of redis client")
	cmd.Flags().IntVar(&jc.concurrency, "concurrency", 4, "concurrency of job")

	return cmd
}

func Command(workers ...Worker[string]) *cobra.Command {
	return newJob().Command(workers...)
}

func Run(workers ...Worker[string]) {
	newJob().Run(workers...)
}
