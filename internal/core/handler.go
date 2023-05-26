package core // 声明了这个文件所在的包名为 core

import (
	"encoding/json" // 导入 encoding/json 包，用于处理 JSON 数据
	"github.com/Aliceikkk/ViaGenshin/internal/mapper" // 导入 mapper 包，在代码中使用 mapper.Protocol 类型
)

func (s *Session) HandlePacket(from, to mapper.Protocol, name string, data []byte) ([]byte, error) {
	// 定义了 Session 结构体（s）的 HandlePacket 方法，并传入了参数 from, to, name, data
	// mapper.Protocol 是一个自定义的类型，用于表示通讯协议
	// HandlePacket 根据 name 的值确定要调用哪个方法进行处理，具体实现在后面
	switch name {
	case "GetPlayerTokenReq":
		return s.OnGetPlayerTokenReq(from, to, data) // 调用 OnGetPlayerTokenReq 方法处理请求
	case "GetPlayerTokenRsp":
		return s.OnGetPlayerTokenRsp(from, to, data) // 调用 OnGetPlayerTokenRsp 方法处理请求
	case "UnionCmdNotify":
		return s.OnUnionCmdNotify(from, to, data) // 调用 OnUnionCmdNotify 方法处理请求
	case "ClientAbilityChangeNotify":
		return s.OnClientAbilityChangeNotify(from, to, data) // 调用 OnClientAbilityChangeNotify 方法处理请求
	case "AbilityInvocationsNotify":
		return s.OnAbilityInvocationsNotify(from, to, data) // 调用 OnAbilityInvocationsNotify 方法处理请求
	case "CombatInvocationsNotify":
		return s.OnCombatInvocationsNotify(from, to, data) // 调用 OnCombatInvocationsNotify 方法处理请求
	}
	return data, nil // 当 name 不是预期值时，返回原始数据和 nil
}

type UnionCmdNotify struct {
	CmdList []*UnionCmd `json:"cmdList"` // 定义结构体 UnionCmdNotify，表示联合指令通知
}

type UnionCmd struct {
	MessageID uint16 `json:"messageId"` // 定义结构体 UnionCmd，表示联合指令
	Body      []byte `json:"body"` // 用 byte 数组表示联合指令体
}

func (s *Session) OnUnionCmdNotify(from, to mapper.Protocol, data []byte) ([]byte, error) {
	// 定义了 Session 结构体（s）的 OnUnionCmdNotify 方法，并传入了参数 from, to, data
	notify := new(UnionCmdNotify) // 创建 UnionCmdNotify 类型的结构体变量 notify
	err := json.Unmarshal(data, notify) // 将传入的 data 进行解析，解析出的结果保存在 notify 中
	if err != nil { // 如果解析过程中出现错误
		return data, err // 直接返回解析前的原始数据和错误信息
	}
	for _, cmd := range notify.CmdList { // 遍历联合指令通知中所有的联合指令
		name := s.mapping.CommandNameMap[from][cmd.MessageID] // 获取联合指令对应的名称，即 CommandNameMap 中 from 协议下 cmd 对应的名称
		cmd.MessageID = s.mapping.CommandPairMap[from][to][cmd.MessageID] // 更新联合指令的 MessageID，即 CommandPairMap 中 from 协议转换到 to 协议下 cmd 对应的 MessageID
		cmd.Body, err = s.ConvertPacketByName(from, to, name, cmd.Body) // 将联合指令体按照相应的协议进行转换，得到转换后的联合指令体
		if err != nil { // 如果转换过程中出现错误
			return data, err // 直接返回解析前的原始数据和错误信息
		}
	}
	return json.Marshal(notify) // 将处理后的结果以 JSON 的形式进行返回
}
