/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"go-icp-checker/checker"
	"go-icp-checker/utils"
)

var (
	unitName string
	oYaml    string // yaml 输出结果
	oJson    string // json 输出结果
	oCsv     string // csvg 输出结果
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "检测备案信息",
	Long:  `查询输入对象的备案信息，`,
	Run: func(cmd *cobra.Command, args []string) {
		infos, err := CheckUnitInfo(unitName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// 这里来解决输出文件的格式
		for _, info := range infos {
			fmt.Println(utils.Prettify(info))

		}

	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVarP(&unitName, "unitName", "u", "", "主域名，公司名，备案号")
	_ = checkCmd.MarkFlagRequired("unitName")

	checkCmd.Flags().StringVarP(&oYaml, "oyaml", "", "", "yaml 格式输出")
	checkCmd.Flags().StringVarP(&oJson, "ojson", "", "", "json 格式输出")
	checkCmd.Flags().StringVarP(&oCsv, "oCsv", "", "", "Csv 格式输出")

}

// CheckUnitInfo 查询域名/公司名/备案号对应的备案信息
// 北京无忧创想信息技术有限公司
// 51cto.com.cn
// 京ICP备09067568号-5
// 京ICP备09067568号
func CheckUnitInfo(unitName string) ([]checker.UnitInfo, error) {
	client := checker.NewIcpClient()
	err := client.GetCookies()
	if err != nil {
		return nil, err
	}
	err = client.GetToken()
	if err != nil {
		return nil, err
	}

	uuid, distance, err := client.ImageVerify()
	if err != nil {
		return nil, err
	}

	_, err = client.GetSign(uuid, distance)
	if err != nil {
		return nil, err
	}

	domainInfo, err := client.GetIcpInfo(unitName)
	if err != nil {
		return nil, err
	}
	if domainInfo.Params.Total <= 0 {
		return nil, errors.New("没有备案信息")
	}

	return domainInfo.Params.List, nil
}
