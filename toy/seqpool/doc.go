/*

	1. 按照group进行分组, 同一个group的任务不会并发执行, 不同group的任务会并发执行
	2. 每个group的任务会按照添加的顺序执行
	3. 不同group的任务需要保证一定的公平性, 避免某个group的任务一直得不到执行

*/

package seqpool
