package command

import "flag"

// Command /*命令行参数解析*/
func Command() map[string]string {
	pEnvStr := flag.String("env","prod","runtime environment, dev|prod")
	pCfgStr := flag.String("cfg","app.conf","config file, example:app.conf")
	flag.Parse()
	result := make(map[string]string)
	result["env"] = *pEnvStr
	result["cfg"] = *pCfgStr
	return result
}