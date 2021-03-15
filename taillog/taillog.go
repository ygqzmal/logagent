package taillog

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"text.com/day11/logagent/kafka"
)

//专门从日志文件收集日志的模块

//TailTask: 一个日志收集的任务
type TailTask struct {
	path string
	topic string
	instance *tail.Tail
	//为了能实现退t.run()
	ctx context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:       path,
		topic:      topic,
		ctx:	    ctx,
		cancelFunc: cancel,
	}
	tailObj.init() //根据路径去打开对应的日志
	return
}

func (t *TailTask) init() {
	config := tail.Config{
		ReOpen:true,				//重新打开
		Follow:true,				//是否跟随
		Location: &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件的哪个地方开始读
		MustExist:false,			//文件不存在是否报错
		Poll:true,
	}
	//用config这个配置去打开文件
	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
	}
	//当goroutine执行的函数退出的时候, goroutine就结束了
	go t.run() //直接去采集日志发送到Kafka
}

func (t *TailTask) run() {
	for {
		select {
		//context执行cancel函数后, 会给ctx.Done通道发送消息
		case <- t.ctx.Done():
			fmt.Printf("tail task: %s_%s 结束了...\n", t.path, t.topic)
			return
		case line := <- t.instance.Lines ://从tailObj的通道中一行一行的读取日志数据
			//3.2 发往Kafka
			//kafka.SendToKafka(t.topic, line.Text) //函数调函数，同步操作，有点慢
			//解决: 先把日志数据发到一个通道中
			kafka.SendToChan(t.topic, line.Text)
			//kafka那个包中有单独的goroutine去取日志数据发送到kafka

		}
	}
}

