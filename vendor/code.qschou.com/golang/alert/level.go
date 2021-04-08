package alert

// alert_level 告警级别
type Level int64

// 告警级别档位常量
const (
	Highest Level = 5 // 最高级别关注
	High    Level = 4 // 高级别关注
	Middle  Level = 3 // 中等级别关注
	Low     Level = 2 // 低级别关注
	Lowest  Level = 1 // 最低级别关注
)

// 告警级别档位描述
const (
	HighestAttention = "最高级别关注【5｜级别1~5】"
	HighAttention    = "高级别关注【4｜级别1~5】"
	MiddleAttention  = "中等级别关注【3｜级别1~5】"
	LowAttention     = "低级别关注【2｜级别1~5】"
	LowestAttention  = "最低级别关注【1｜级别1~5】"
)

// 告警级别和描述的映射
var level2Attention = map[Level]string{
	Highest: HighestAttention,
	High:    HighAttention,
	Middle:  MiddleAttention,
	Low:     LowAttention,
	Lowest:  LowestAttention,
}

// 获取Level对应的Attention
func GetAttentionByLevel(level Level) string {
	desc, ok := level2Attention[level]
	if !ok {
		return LowestAttention
	}
	return desc
}
