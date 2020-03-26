# Go routine

Go routine is a program to handle a pool of workers with go channel.


## Usage

Start the pool with a pool of 10 workers and 10 jobs that will be created.
Here is the output.
```bash
go run cmd/main.go -workers 10 -jobs 10
020/03/26 15:53:07 [INFO] starting HTTP server on port :9009
2020/03/26 15:53:07 [INFO] starting 10 jobs
2020/03/26 15:53:07 [INFO] starting 10 workers
2020/03/26 15:53:07 [workerID:2][email_job] Processed job in 0.023 seconds
2020/03/26 15:53:07 [workerID:3][email_job] Processed job in 0.020 seconds
2020/03/26 15:53:07 [workerID:4][email_job] Processed job in 0.021 seconds
2020/03/26 15:53:07 [workerID:5][email_job] Processed job in 0.021 seconds
2020/03/26 15:53:07 [workerID:9][notification_job] Processed job in 0.062 seconds
2020/03/26 15:53:07 [workerID:7][email_job] Processed job in 0.020 seconds
2020/03/26 15:53:07 [workerID:0][notification_job] Processed job in 0.062 seconds
2020/03/26 15:53:07 [workerID:1][notification_job] Processed job in 0.061 seconds
2020/03/26 15:53:07 [workerID:6][notification_job] Processed job in 0.062 seconds
2020/03/26 15:53:07 [workerID:8][notification_job] Processed job in 0.062 seconds
```

You can also push a random job to the pool with an endpoint
```bash
curl localhost:9009/new-job
```