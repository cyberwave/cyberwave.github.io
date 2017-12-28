[原方](https://geeks.uniplaces.com/building-a-worker-pool-in-golang-1e6c0fdfd78c)

# 动机
最近，我们创建了一个处理Salesforce和我们的Core（Uniplaces主API）之间的所有事件的服务。

![photo](https://cdn-images-1.medium.com/max/1600/1*MAVJ_-gowsfvaftsytvUXQ.png)

每当有需要在Salesforce中表示的平台上的事件时，Core将发布消息到FIFO队列。该服务将轮询队列中的消息，创建一个有数据更改的负载并将其发送给Salesforce。

我们希望这个过程很快，以确保两个系统中的信息尽快同步，并可扩展以适应Uniplaces的快速增长。

线程池（也称为工作池）允许我们轻松地管理具有可配置大小的消费者线程的集合。

在这篇文章中，我将重点介绍我们在实现这个时所考虑的消费者/工作池。

# 构架
我们希望为工作池提供的主要功能是：

- **动态可伸缩性：** 工作者可以*告诉*池它是否实际正在处理消息，如果是，池将立即产生另一个工人，而不是在X秒钟内休眠。这允许池关于队列中的消息量是动态的。
- **优雅关闭：** 池进入关闭状态，没有新的工作人员被创建，当前正在运行的人可以在退出之前完成他们的工作。

考虑到这一点，我们的架构可以用下图来描述：
![架构图](https://cdn-images-1.medium.com/max/1600/1*eA_s0e-l4zchlAvSd5-Etw.png)  
即使这正好描述了我们想要的内容，但是如果我们按照图中所描述的那样来实现这个功能，则会导致非常混乱的代码，随处可见*ifs*和代码重复。

所以我们决定最好把它作为一个状态机来实现：
![状态机](https://cdn-images-1.medium.com/max/1600/1*--KUVBRkifrTyD-ZyEWzZQ.png)

这使我们有更清晰的代码，同时得到相同的结果。

*好的...用可爱的图表就够了，让我们来看看代码！*

# 实现
## 工作者
工作者由一个标识符和两个通道组成，以将其状态传达给池。
```golang
package worker

type ID struct {
	uuid.UUID
}

type Worker struct {
	ID                ID
	ProcessingChannel chan<- bool
	FinishedChannel   chan<- ID
}

func NewWorker(processingChannel chan<- bool, finishedChannel chan<- ID) Worker {
	return Worker{
		ID:                generateID(),
		ProcessingChannel: processingChannel,
		FinishedChannel:   finishedChannel,
	}
}

func generateID() ID {
	return ID{UUID: uuid.New()}
}
```
**worker.go**被托管在[GitHub](https://gist.github.com/jmdalmeida/cfe5369f0916370ee3496cee6220989a#file-worker-go)

## 消费者
每个工作者将调用这个 Consume 方法，它轮询队列并处理消息更新池的状态。
```golang
func (queue Queue) Consume(worker worker.Worker) error {
	// signal the pool that the worker has finished after the method is executed
	defer func() { worker.FinishedChannel <- worker.ID }()

	// poll the queue
	message, err := queue.service.Read(queue.name)
	if err != nil {
		// signal the pool there's no message being processed
		worker.ProcessingChannel <- false

		return err
	}

	if len(message.Messages) == 0 {
		// signal the pool there's no message being processed
		worker.ProcessingChannel <- false

		return nil
	}
	// signal the pool that a message will be processed
	worker.ProcessingChannel <- true
  
	// handle the message...
}
```
**consumer.go**托管在[GitHub](https://gist.github.com/jmdalmeida/b4faf303f0167b48ee476041171cdab0#file-consumer-go)

# 工作者池
## 状态
状态由每个周期（池的内部主循环）更新的转换结构表示，并且在某些情况下可以在状态之间持续信息。
```golang
package workerpool

const (
	stateInitial = iota
	stateMain
	stateExit
	stateLaunch
	stateProcessing
	stateFinish
	stateSleep
	stateWait
	stateTimeout
	stateQuit
)

type transition struct {
	state   int
	payload payload
}

type payload struct {
	workerID     worker.ID
	isProcessing bool
}
```
**state.go**托管在[github](https://gist.github.com/jmdalmeida/f6222e25465f328c353ae07d4a537288)

## 池
工作池使用 max workers，interval 和 timeout 参数进行初始化，这个参数允许跨环境的灵活性，以及​​一个关闭*对象*来控制当它退出时池如何处理。
```golang
type WorkerPool struct {
	maxWorkers        int
	interval          int
	timeout           int
	shutdown          *Shutdown
	processingChannel chan bool
	finishedChannel   chan worker.ID
	needsToQuit       bool
	activeWorkers     map[worker.ID]carbon.Carbon
	currentTransition transition
}

func NewWorkerPool(config Config, shutdown *Shutdown) WorkerPool {
	maxWorkers := config.Workers
	initialTransition := transition{
		state: stateInitial,
	}

	return WorkerPool{
		maxWorkers:        maxWorkers,
		interval:          config.Interval,
		timeout:           config.Timeout,
		shutdown:          shutdown,
		processingChannel: make(chan bool, maxWorkers),
		finishedChannel:   make(chan worker.ID, maxWorkers),
		activeWorkers:     map[worker.ID]carbon.Carbon{},
		currentTransition: initialTransition,
	}
}
```
**workerpool.go**托管于[github](https://gist.github.com/jmdalmeida/f1a54719a5c12d1db8c06c97219e1756#file-workerpool-go)

开始方法表示系统当前处理状态的主循环，这个方法接收由工作者调用的函数。
```golang
func (workerPool *WorkerPool) Start(actionHandler func(worker worker.Worker)) {
	for workerPool.currentTransition.state != stateExit {
		switch workerPool.currentTransition.state {
		case stateInitial:
			workerPool.processInitialState()
		case stateMain:
			workerPool.processMainState()
		case stateWait:
			workerPool.processWaitState()
		case stateSleep:
			workerPool.processSleepState()
		case stateLaunch:
			workerPool.processLaunchState(actionHandler)
		case stateQuit:
			workerPool.processQuitState()
		case stateTimeout:
			workerPool.processTimeoutState()
		case stateFinish:
			workerPool.processFinishState()
		case stateProcessing:
			workerPool.processProcessingState()
		default:
			panic("invalid state transition")
		}
	}
	workerPool.shutdown.doneChannel <- true
}
```
**workerpool.go**托管于[github](https://gist.github.com/jmdalmeida/7abbcbd79898227b600783d77f1384cc#file-workerpool-go)

接下来，我们有 *processLaunchState* 函数，它将负责创建一个新的工作者，并用 *actionHandler* 函数产生一个新的*goroutine*。

```golang
func (workerPool *WorkerPool) processLaunchState(actionHandler func(worker worker.Worker)) {
	newWorker := worker.NewWorker(workerPool.processingChannel, workerPool.finishedChannel)
	workerPool.activeWorkers[newWorker.ID] = *carbon.Now()

	go actionHandler(newWorker)

	workerPool.goToState(stateWait, nil)
}
```
**workerpool.go**托管于[github](https://gist.github.com/jmdalmeida/934cdd17a9fe216db2fb2557f8403d2d#file-workerpool-go)

最后，*processWaitState* 函数将监听来自工作者（已完成和正在处理的通道）或来自应用程序本身（shutdown channel），或者在出现错误，超时的事件。
```golang
func (workerPool *WorkerPool) processWaitState() {
	select {
	case finished := <-workerPool.finishedChannel:
		payload := &payload{
			workerID: finished,
		}
		workerPool.goToState(stateFinish, payload)
	case processing := <-workerPool.processingChannel:
		payload := &payload{
			isProcessing: processing,
		}
		workerPool.goToState(stateProcessing, payload)
	case <-workerPool.shutdown.initiateChannel:
		workerPool.goToState(stateQuit, nil)
	case <-time.After(time.Second * time.Duration(workerPool.timeout)):
		workerPool.goToState(stateTimeout, nil)
	}
}
```
**workerpool.go**被托管于[github](https://gist.github.com/jmdalmeida/bba479670998d3597cefa78e5b31ea77#file-workerpool-go)

*注：有些转换非常简单，所以为简单起见，省略了方法。*

## 关闭
关机的任务是让主线程等待来自系统的中断信号(**Ctrl + C**)，这个信号会 *ping* 这个池来启动关机，这意味着没有新的工作者被创建，并且已经运行的人员有时间完成。然后等待第二个 *中断* 信号（立即强制关闭），一个来自池的信号指示所有工人已经完成或超时，让主线程完成。
```golang
type Shutdown struct {
	initiateChannel chan bool
	doneChannel     chan bool
	timeout         time.Duration
}

func NewShutdown(timeout time.Duration) *Shutdown {
	return &Shutdown{
		initiateChannel: make(chan bool, 1),
		doneChannel:     make(chan bool, 1),
		timeout:         timeout,
	}
}

func (shutdown *Shutdown) WaitForSignal() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel

	logger.Info("received interrupt signal")
	shutdown.initiateChannel <- true
	select {
	case <-signalChannel:
		logger.Warning("forcing shutdown")
		os.Exit(1)
	case <-shutdown.doneChannel:
		logger.Info("cleanup successful, exiting")
	case <-time.After(time.Second * shutdown.timeout):
		logger.Info("cleanup timed out, exiting")
	}
}
```
**shutdown.go**被托管在[github](https://gist.github.com/jmdalmeida/9fb33035f4819167f875cfbdf77abbda#file-shutdown-go)

## 命令

该消费者从 CLI 上下文中执行。以下命令将在生产环境中启动消费者， 每个周期之间有 10 个最大工作者和 10 秒间隔：
```shell
GOENV=prod ./cli -w 10 -i 10
```

该命令创建如下：
```golang
func Command(interactor usecase.SalesforceInteractor) cli.Command {
	var shutdownTimeout int
	var config workerpool.Config

	return cli.Command{
		Name:  "consumer",
		Usage: "Consume Salesforce queue events from core",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:        "workers, w",
				Usage:       "Number of go workers available to launch",
				Destination: &config.Workers,
				Value:       2,
			},
			cli.IntFlag{
				Name:        "interval, i",
				Usage:       "Time interval to wait when a worker has nothing to process",
				Destination: &config.Interval,
				Value:       30,
			},
			cli.IntFlag{
				Name:        "timeout, t",
				Usage:       "Timeout for the worker pool to wait for worker events",
				Destination: &config.Timeout,
				Value:       60,
			},
			cli.IntFlag{
				Name:        "shutdown-timeout, s",
				Usage:       "Timeout for the worker pool cleanup before shutting down",
				Destination: &shutdownTimeout,
				Value:       60,
			},
		},
		Action: func(c *cli.Context) error {
			actionHandler := consume(interactor)

			shutdown := workerpool.NewShutdown(time.Duration(shutdownTimeout))
			workerPool := workerpool.NewWorkerPool(config, shutdown)
			go workerPool.Start(actionHandler)
			shutdown.WaitForSignal()

			return nil
		},
	}
}

func consume(interactor usecase.SalesforceInboundInteractor) func(worker worker.Worker) {
	return func(worker worker.Worker) {
		interactor.ConsumeEventFromQueue(worker)
	}
}
```
**command.go**被托管在[github](https://gist.github.com/jmdalmeida/e018c007a93db963e1cf2f50c2e388d5#file-command-go)

正如您在Action处理函数中看到的，WorkerPool是使用来自命令行调用和Shutdown处理程序的参数创建的。然后启动函数成为守护进程，主线程等待关闭信号。

# 未来的改进
这是这个系统的第一次迭代，所以有一些方面可以在未来得到改进。例如，一个线程池应该有线程预初始化，等待接收工作。这将避免让池总是创建和销毁线程，这是一个资源繁重的操作。另一个改进是将池提取到单独的供应商，并使其可用于不同于消耗队列的上下文。
