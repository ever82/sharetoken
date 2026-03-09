package main

import (
	"fmt"

	"github.com/huin/goupnp/dcps/internetgateway1"
	"github.com/huin/goupnp/dcps/internetgateway2"
)

func main() {
	fmt.Println("=== UPnP 设备发现测试 ===")
	fmt.Println()

	// Try IGD v2 first
	fmt.Println("1. 搜索 UPnP IGD v2...")
	clients2, _, err := internetgateway2.NewWANIPConnection2Clients()
	if err == nil && len(clients2) > 0 {
		fmt.Println("   ✅ 找到 IGD v2 设备!")
		ip, err := clients2[0].GetExternalIPAddress()
		if err == nil {
			fmt.Printf("   外部IP: %s\n", ip)
		}
		return
	}
	if err != nil {
		fmt.Printf("   IGD v2 错误: %v\n", err)
	}

	// Try IGD v1
	fmt.Println("2. 搜索 UPnP IGD v1...")
	clients1, _, err := internetgateway1.NewWANIPConnection1Clients()
	if err == nil && len(clients1) > 0 {
		fmt.Println("   ✅ 找到 IGD v1 设备!")
		ip, err := clients1[0].GetExternalIPAddress()
		if err == nil {
			fmt.Printf("   外部IP: %s\n", ip)
		}
		return
	}
	if err != nil {
		fmt.Printf("   IGD v1 错误: %v\n", err)
	}

	// Try PPP
	fmt.Println("3. 搜索 UPnP PPP 连接...")
	pppClients, _, err := internetgateway1.NewWANPPPConnection1Clients()
	if err == nil && len(pppClients) > 0 {
		fmt.Println("   ✅ 找到 PPP 连接!")
		return
	}
	if err != nil {
		fmt.Printf("   PPP 错误: %v\n", err)
	}

	fmt.Println()
	fmt.Println("❌ 未找到任何UPnP设备")
	fmt.Println()
	fmt.Println("可能原因:")
	fmt.Println("   - 路由器UPnP功能未开启")
	fmt.Println("   - 路由器不支持UPnP")
	fmt.Println("   - 网络防火墙阻止了SSDP发现")
	fmt.Println()
	fmt.Println("解决方法:")
	fmt.Println("   1. 登录路由器管理界面 http://192.168.3.1")
	fmt.Println("   2. 找到 NAT/UPnP 设置")
	fmt.Println("   3. 开启 UPnP 功能")
}
