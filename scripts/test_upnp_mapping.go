package main

import (
	"fmt"
	"net"
	"time"

	"github.com/huin/goupnp/dcps/internetgateway1"
)

func main() {
	fmt.Println("=== UPnP 端口映射测试 ===")
	fmt.Println()

	// 发现设备
	clients, _, err := internetgateway1.NewWANIPConnection1Clients()
	if err != nil || len(clients) == 0 {
		fmt.Println("❌ 未找到UPnP设备:", err)
		return
	}
	client := clients[0]
	fmt.Println("✅ 找到UPnP IGD v1设备")

	// 获取本地IP
	localIP, _ := getLocalIP()
	fmt.Printf("本地IP: %s\n", localIP)

	// 获取外部IP
	externalIP, err := client.GetExternalIPAddress()
	if err != nil {
		fmt.Printf("⚠️  获取外部IP失败: %v\n", err)
	} else {
		fmt.Printf("外部IP: %s\n", externalIP)
	}

	// 测试端口映射
	testPort := uint16(26656)
	fmt.Printf("\n3. 添加端口映射 (TCP %d -> %s:%d)...\n", testPort, localIP, testPort)

	err = client.AddPortMapping(
		"",           // NewRemoteHost
		testPort,     // ExternalPort
		"TCP",        // Protocol
		testPort,     // InternalPort
		localIP,      // InternalClient
		true,         // Enabled
		"ShareToken Test", // Description
		3600,         // LeaseDuration (1 hour)
	)
	if err != nil {
		fmt.Printf("❌ 添加端口映射失败: %v\n", err)
		return
	}
	fmt.Println("✅ 端口映射添加成功!")

	// 验证映射
	fmt.Println("\n4. 端口映射已创建:")
	fmt.Printf("   外部: %s:%d (TCP)\n", externalIP, testPort)
	fmt.Printf("   内部: %s:%d\n", localIP, testPort)

	// 保持映射10秒
	fmt.Println("\n5. 等待10秒后清理...")
	time.Sleep(10 * time.Second)

	// 删除映射
	fmt.Println("6. 删除端口映射...")
	err = client.DeletePortMapping("", testPort, "TCP")
	if err != nil {
		fmt.Printf("⚠️  删除端口映射失败: %v\n", err)
	} else {
		fmt.Println("✅ 端口映射已删除")
	}

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println()
	fmt.Println("注意: 如果你的外部IP是内网IP（如192.168.x.x），")
	fmt.Println("     说明有双层NAT，外部用户仍然无法直接访问。")
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no local IP found")
}
