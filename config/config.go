package config

//为了和ini文件中的内容对应上
//这里需要添加tag标签
//例如在 Address string 中 Address是大写开头对应不上,所以添加tag
//如果一样，就不需要

type AppConf struct {
	KafkaConf  		`ini:"kafka"`
	EtcdConf		`ini:"etcd"`
}

type EtcdConf struct {
	Key string 		`ini:"collect_log_key"`
	Address string	`ini:"address"`
	Timeout int 	`ini:"timeout"`
}


type KafkaConf struct {
	Address string	`ini:"address"`
	ChanMaxSize int `ini:"chan_max_size"`
}


//--- unused ---
type TaillogConf struct {
	FileName string	`ini:"filename"`
}
