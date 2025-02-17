package seqpool

import (
	"context"
	"runtime"
	"sync"
)

type IRunnable interface {
	IsStop() bool // 是否停止
}

type group[K comparable, V IRunnable] struct {
	key     K          // group的key
	val     V          // group的上下文
	lock    sync.Mutex // 处理任务添加的并发问题
	isInit  sync.Once
	running bool         // 是否正在运行
	stopped bool         // 是否停止
	tasks   []func(K, V) // 任务队列
}

// tryWakeUp 尝试唤醒
func (g *group[K, V]) tryWakeUp(fn func(K, V)) bool {
	// return g.running.CompareAndSwap(false, true)
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.stopped {
		return false
	}
	g.tasks = append(g.tasks, fn)
	if g.running {
		return false // 已经在运行了
	}
	g.running = true // 标记为正在运行, 并通知上层可以执行
	return true
}

// trySleep 返回是否休眠
func (g *group[K, V]) trySleep(isStop bool) bool {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.stopped = isStop
	if isStop {
		g.tasks = nil // 清空任务
	}
	g.running = len(g.tasks) > 0 // 有任务就继续执行
	return !g.running            // 没有任务就休眠
}

// tryPopTask 尝试弹出任务
func (g *group[K, V]) tryPopTask() (fn func(K, V)) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if len(g.tasks) == 0 {
		return nil
	}
	fn = g.tasks[0]
	g.tasks[0] = nil
	g.tasks = g.tasks[1:]
	return fn
}

// init 初始化
func (g *group[K, V]) init(initFn func(K) V) {
	g.isInit.Do(func() { g.val = initFn(g.key) })
}

type Pool[K comparable, V IRunnable] struct {
	wg        sync.WaitGroup
	groups    sync.Map
	ctx       context.Context
	stop      func()
	seqChan   chan *group[K, V] // 待执行的group
	groupInit func(K) V
}

func NewPool[K comparable, V IRunnable](ctx context.Context, groupInit func(K) V) *Pool[K, V] {
	ctx, stop := context.WithCancel(ctx)
	p := &Pool[K, V]{
		ctx:       ctx,
		stop:      stop,
		seqChan:   make(chan *group[K, V], 16), // TODO 任务等待队列的大小
		groupInit: groupInit,
	}
	p.start()
	return p
}

// getGroup 获取group
func (p *Pool[K, V]) getGroup(key K) *group[K, V] {
	g, ok := p.groups.Load(key)
	if !ok {
		return nil
	}
	gg, _ := g.(*group[K, V])
	return gg
}

// addGroup 添加group
func (p *Pool[K, V]) addGroup(key K) *group[K, V] {
	gg, _ := p.groups.LoadOrStore(key, &group[K, V]{key: key})
	return gg.(*group[K, V])
}

// delGroup 删除group
func (p *Pool[K, V]) delGroup(key K, g *group[K, V]) {
	p.groups.CompareAndDelete(key, g) // 避免删除新创建的group
}

// Submit 提交任务
func (p *Pool[K, V]) Submit(key K, task func(K, V)) {
	group := p.getGroup(key)
	if group == nil {
		group = p.addGroup(key)
	}
	group.init(p.groupInit) // sync.Once保证不会重复初始化
	if group.tryWakeUp(task) {
		p.seqChan <- group
	}
}

// Stop 停止
func (p *Pool[K, V]) Stop() {
	p.stop()
	p.wg.Wait()
}

// start 启动
func (p *Pool[K, V]) start() {
	total := runtime.NumCPU()
	p.wg.Add(total)
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-p.ctx.Done():
					return
				case gg := <-p.seqChan:
					g := p.getGroup(gg.key)
					if g != gg {
						// 不是同一个group, 这里可能是group已经被删除了
						continue
					}
					func(group *group[K, V]) {
						defer func() {
							if err := recover(); err != nil {
								// TODO log
							}
						}()
						var task = group.tryPopTask()
						if task != nil {
							task(group.key, group.val)
						}
						isStop := group.val.IsStop()
						if isStop {
							p.delGroup(group.key, group) // 删除group
						}
						if !group.trySleep(isStop) {
							p.seqChan <- group // 加入到队列中准备下次执行
						}
					}(g)
				}
			}
		}()
	}
}
