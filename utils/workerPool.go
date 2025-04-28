package utils

import (
	"fmt"
	"sync"
)

type WorkerPool struct {
	workers  int
	jobQueue chan func()
	wg       sync.WaitGroup
}

// NewWorkerPool 创建一个新的工作池
func NewWorkerPool(workers int) *WorkerPool {
	pool := &WorkerPool{
		workers:  workers,
		jobQueue: make(chan func()),
	}
	return pool
}

// Run 启动线程池
func (pool *WorkerPool) Run() {
	for i := 0; i < pool.workers; i++ {
		go func(workerID int) {
			for job := range pool.jobQueue {
				fmt.Printf("Worker %d processing job\n", workerID)
				job() // 执行工作
			}
		}(i)
	}
}

// Submit 提交一个任务到线程池
func (pool *WorkerPool) Submit(job func()) {
	pool.jobQueue <- job
	pool.wg.Add(1)
}

// Wait 等待所有任务完成
func (pool *WorkerPool) Wait() {
	pool.wg.Wait()
}

// Close 关闭工作池
func (pool *WorkerPool) Close() {
	close(pool.jobQueue)
}
