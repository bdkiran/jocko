package fsm

import (
	"fmt"

	"github.com/bdkiran/nolan/nolan/structs"
	"go.uber.org/zap"
)

func init() {
	registerCommand(structs.RegisterNodeRequestType, (*FSM).applyRegisterNode)
	registerCommand(structs.DeregisterNodeRequestType, (*FSM).applyDeregisterNode)
	registerCommand(structs.RegisterTopicRequestType, (*FSM).applyRegisterTopic)
	registerCommand(structs.DeregisterTopicRequestType, (*FSM).applyDeregisterTopic)
	registerCommand(structs.RegisterPartitionRequestType, (*FSM).applyRegisterPartition)
	registerCommand(structs.DeregisterPartitionRequestType, (*FSM).applyDeregisterPartition)
	registerCommand(structs.RegisterGroupRequestType, (*FSM).applyRegisterGroup)
}

func (c *FSM) applyRegisterGroup(buf []byte, index uint64) interface{} {
	var req structs.RegisterGroupRequest
	if err := structs.Decode(buf, &req); err != nil {
		zap.S().Panicf("failed to decode request: %v", err)
	}

	if err := c.state.EnsureGroup(index, &req.Group); err != nil {
		zap.S().Errorf("EnsureNode error: %s", err)
		return err
	}

	return nil
}

func (c *FSM) applyRegisterNode(buf []byte, index uint64) interface{} {
	var req structs.RegisterNodeRequest
	if err := structs.Decode(buf, &req); err != nil {
		zap.S().Panicf("failed to decode request: %v", err)
	}

	if err := c.state.EnsureNode(index, &req.Node); err != nil {
		zap.S().Errorf("EnsureNode error: %s", err)
		return err
	}

	return nil
}

func (c *FSM) applyDeregisterNode(buf []byte, index uint64) interface{} {
	var req structs.DeregisterNodeRequest
	if err := structs.Decode(buf, &req); err != nil {
		panic(fmt.Errorf("failed to decode request: %v", err))
	}

	if err := c.state.DeleteNode(index, req.Node.Node); err != nil {
		zap.S().Errorf("DeleteNode error: %s", err)
		return err
	}

	return nil
}

func (c *FSM) applyRegisterTopic(buf []byte, index uint64) interface{} {
	var req structs.RegisterTopicRequest
	if err := structs.Decode(buf, &req); err != nil {
		zap.S().Panicf("failed to decode request: %v", err)
	}

	if err := c.state.EnsureTopic(index, &req.Topic); err != nil {
		zap.S().Errorf("EnsureTopic error: %s", err)
		return err
	}

	return nil
}

func (c *FSM) applyDeregisterTopic(buf []byte, index uint64) interface{} {
	var req structs.DeregisterTopicRequest
	if err := structs.Decode(buf, &req); err != nil {
		zap.S().Panicf("failed to decode request: %v", err)
	}

	if err := c.state.DeleteTopic(index, req.Topic.Topic); err != nil {
		zap.S().Errorf("DeleteTopic error: %s", err)
		return err
	}

	return nil
}

func (c *FSM) applyRegisterPartition(buf []byte, index uint64) interface{} {
	var req structs.RegisterPartitionRequest
	if err := structs.Decode(buf, &req); err != nil {
		zap.S().Panicf("failed to decode request: %v", err)
	}

	if err := c.state.EnsurePartition(index, &req.Partition); err != nil {
		zap.S().Errorf("EnsurePartition error: %s", err)
		return err
	}

	return nil
}

func (c *FSM) applyDeregisterPartition(buf []byte, index uint64) interface{} {
	var req structs.DeregisterPartitionRequest
	if err := structs.Decode(buf, &req); err != nil {
		zap.S().Panicf("failed to decode request: %v", err)
	}

	if err := c.state.DeletePartition(index, req.Partition.Topic, req.Partition.Partition); err != nil {
		zap.S().Errorf("DeletePartition error: %s", err)
		return err
	}

	return nil
}
