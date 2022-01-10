package snowflake

// ID生成器
import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

// 初始化，snowflake算法，随机生成user_id
// 时间戳：startTime 从什么时候开始，比如2020-07-01开始，能用到往后的69年
// machineID： 工作机器的id
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	// 按照指定格式解析时间戳
	// 2006-01-02：是go语言诞生的时间
	// 转换成：2020-07-01 00:00:00 +0000 UTC这种格式
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	// 将startTime的时间作为起始时间，单位是毫秒
	// 会将st转城1593561600000（int64）这种格式
	// 会把sf.Epoch设定的值更换，已到达从startTime从该年开始的目的
	sf.Epoch = st.UnixNano() / 1000000
	// 生成雪花节点，传入一个工作机器节点id
	node, err = sf.NewNode(machineID)
	return
}

// 生成雪花ID
func GenID() int64 {
	// 将雪花ID转成int64类型
	return node.Generate().Int64()
}
