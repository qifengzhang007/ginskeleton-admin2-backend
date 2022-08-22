package test

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/utils/rabbitmq/hello_world"
	"goskeleton/app/utils/rabbitmq/publish_subscribe"
	"goskeleton/app/utils/rabbitmq/routing"
	"goskeleton/app/utils/rabbitmq/topics"
	"goskeleton/app/utils/rabbitmq/work_queue"
	_ "goskeleton/bootstrap"
	"os"
	"strconv"
	"testing"
)

//  消息队列（rabbitmq）在线文档地址：https://www.yuque.com/xiaofensinixidaouxiang/bkfhct/tkcuc8
// 延迟消息队列在线文档地址：https://www.yuque.com/xiaofensinixidaouxiang/bkfhct/grroyv
// 本篇的单元测试提供的是非延迟消息队列的测试，只要学会单元测试提供的示例，延迟队列也是非常简单的，参考在线文档即可

// 1.HelloWorld 模式
func TestRabbitMqHelloWorldProducer(t *testing.T) {

	helloProducer, err := hello_world.CreateProducer()
	if err != nil {
		t.Errorf("HelloWorld单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}
	var res bool
	for i := 0; i < 10; i++ {
		str := fmt.Sprintf("%d_HelloWorld开始发送消息测试", i+1)
		res = helloProducer.Send(str)
		//time.Sleep(time.Second * 1)
	}

	helloProducer.Close() // 消息投递结束，必须关闭连接
	// 总共发送了10条消息，我们简单判断一下最后一条消息返回的结果
	if res {
		t.Log("消息发送OK")
	} else {
		t.Errorf("HelloWorld模式消息发送失败")
	}
}

// 消费者
func TestMqHelloWorldConsumer(t *testing.T) {

	consumer, err := hello_world.CreateConsumer()
	if err != nil {
		t.Errorf("HelloWorld单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}

	consumer.OnConnectionError(func(err *amqp.Error) {
		t.Errorf(my_errors.ErrorsRabbitMqReconnectFail+"，%s\n", err.Error())
	})

	consumer.Received(func(receivedData string) {

		t.Logf("HelloWorld回调函数处理消息：--->%s\n", receivedData)
	})
}

// 2.WorkQueue模式
func TestRabbitMqWorkQueueProducer(t *testing.T) {

	producer, _ := work_queue.CreateProducer()
	var res bool
	for i := 0; i < 10; i++ {
		str := fmt.Sprintf("%d_WorkQueue开始发送消息测试", i+1)
		res = producer.Send(str)
		//time.Sleep(time.Second * 1)
	}

	producer.Close() // 消息投递结束，必须关闭连接

	if res {
		t.Logf("消息发送OK")
	} else {
		t.Errorf("WorkQueue模式消息发送失败")
	}
}

// 消费者
func TestMqWorkQueueConsumer(t *testing.T) {

	consumer, err := work_queue.CreateConsumer()
	if err != nil {
		t.Errorf("WorkQueue单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}

	consumer.OnConnectionError(func(err *amqp.Error) {
		t.Errorf(my_errors.ErrorsRabbitMqReconnectFail + ", %s" + err.Error())
	})

	consumer.Received(func(receivedData string) {

		t.Logf("WorkQueue回调函数处理消息：--->%s\n", receivedData)
	})
}

// 3.PublishSubscribe 发布、订阅模式模式
func TestRabbitMqPublishSubscribeProducer(t *testing.T) {

	producer, err := publish_subscribe.CreateProducer()
	if err != nil {
		t.Errorf("WorkQueue 单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}
	var res bool
	for i := 0; i < 10; i++ {
		str := fmt.Sprintf("%d_PublishSubscribe开始发送消息测试", i+1)
		// 参数二： 消息延迟的毫秒数，只有创建的对象是延迟模式该参数才有效
		res = producer.Send(str, 1000)
		fmt.Println(str, res)
		//time.Sleep(time.Second * 2)
	}

	producer.Close() // 消息投递结束，必须关闭连接

	if res {
		t.Log("消息发送OK")
	} else {
		t.Errorf("PublishSubscribe 模式消息发送失败")
	}
}

// 消费者
func TestRabbitMqPublishSubscribeConsumer(t *testing.T) {

	consumer, err := publish_subscribe.CreateConsumer()
	if err != nil {
		t.Errorf("PublishSubscribe单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}

	consumer.OnConnectionError(func(err *amqp.Error) {
		t.Errorf(my_errors.ErrorsRabbitMqReconnectFail + "，%s\n" + err.Error())
	})

	consumer.Received(func(receivedData string) {

		t.Logf("PublishSubscribe回调函数处理消息：--->%s\n", receivedData)
	})
}

// Routing 路由模式
func TestRabbitMqRoutingProducer(t *testing.T) {

	producer, err := routing.CreateProducer()

	if err != nil {
		t.Errorf("Routing单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}
	var res bool
	var key string
	for i := 1; i <= 20; i++ {

		//  将 偶数 和  奇数 分发到不同的key，消费者端，启动两个也各自处理偶数和奇数
		if i%2 == 0 {
			key = "key_even" //  偶数键
		} else {
			key = "key_odd" //  奇数键
		}

		//strData := fmt.Sprintf("%d_Routing_%s, 开始发送消息测试"+time.Now().Format("2006-01-02 15:04:05"), i, key)
		// 参数三： 消息延迟的毫秒数，只有创建的对象是延迟模式该参数才有效
		res = producer.Send(key, strconv.Itoa(i)+"- Routing开始发送消息测试", 10000)
		//time.Sleep(time.Second * 1)
	}

	producer.Close() // 消息投递结束，必须关闭连接

	if res {
		t.Logf("消息发送OK")
	} else {
		t.Errorf("Routing 模式消息发送失败")
	}
}

// 消费者
func TestRabbitMqRoutingConsumer(t *testing.T) {
	consumer, err := routing.CreateConsumer()

	if err != nil {
		t.Errorf("Routing单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}

	consumer.OnConnectionError(func(err *amqp.Error) {
		t.Errorf(my_errors.ErrorsRabbitMqReconnectFail + "， %s\n" + err.Error())
	})
	// 通过route_key 匹配指定队列的消息来处理
	consumer.Received("key_even", func(receivedData string) {
		fmt.Println("处理偶数的回调函数 ---> 收到消息内容: " + receivedData)
		//	t.Logf("处理偶数的回调函数：--->收到消息时间：%s - 消息内容：%s\n", time.Now().Format("2006-01-02 15:04:05"), receivedData)
	})
}

// topics 模式
func TestRabbitMqTopicsProducer(t *testing.T) {

	producer, err := topics.CreateProducer()
	if err != nil {
		t.Errorf("Routing单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}
	var res bool
	var key string
	for i := 1; i <= 10; i++ {

		//  将 偶数 和  奇数 分发到不同的key，消费者端，启动两个也各自处理偶数和奇数
		if i%2 == 0 {
			key = "key.even" //  偶数键
		} else {
			key = "key.odd" //  奇数键
		}
		strData := fmt.Sprintf("%d_Routing_%s, 开始发送消息测试", i, key)
		// 参数三： 消息延迟的毫秒数，只有创建的对象是延迟模式该参数才有效
		res = producer.Send(key, strData, 10000)
		//time.Sleep(time.Second * 1)
	}

	producer.Close() // 消息投递结束，必须关闭连接

	if res {
		t.Logf("消息发送OK")
	} else {
		t.Errorf("topics 模式消息发送失败")
	}
	//Output: 消息发送OK
}

// 消费者
func TestRabbitMqTopicsConsumer(t *testing.T) {
	consumer, err := topics.CreateConsumer()

	if err != nil {
		t.Errorf("Routing单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}

	consumer.OnConnectionError(func(err *amqp.Error) {
		t.Errorf(my_errors.ErrorsRabbitMqReconnectFail + "， %s\n" + err.Error())
	})
	// 通过route_key 模糊匹配队列路由键的消息来处理
	consumer.Received("#.odd", func(receivedData string) {

		t.Logf("模糊匹配偶数键：--->%s\n", receivedData)
	})
}
