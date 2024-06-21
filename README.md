使用方法： 1.编译WalletServer、ChainServer，得到WalletServer.exe、ChainServer.exe
2.运行WalletServer.exe，访问http://localhost:8089/CreateNewWallet 得到钥匙串、地址
3.配置.env，格式为nodeX="url"，其中X为任意数字
4，使用多个终端，输入 .\ChainServer.exe -nodeName=nodeX启动对应节点，其中X为第3步配置数字
5.向WalletServer请求签名，格式为：
http://localhost:8089/Sign
{
"sender_addr":"1EFf5VAjfqxy6YmHiUz7nq7MxXX8Sm3ZLQ",
"receiver_addr":"test",
"value":1,
"Info":"nothing",
"private_key":"4e835a25b947e7dedc61f358364c4cf3fa2d86219ca47f15e3984567ebfd57c7"
}
6.使用签名、地址、交易信息向节点请求加入交易池，格式为： http://localhost:9090/PushTransaction {
"sender_addr":"1EFf5VAjfqxy6YmHiUz7nq7MxXX8Sm3ZLQ",
"receiver_addr":"test",
"value":1,
"info":"nothing",
"public_key":"756d43539ad53fe7e41ec86da696cbb4884747d239c7b56b3174d44e511f7c2b50ac277bd2eb838e568baab1c24d47a914630ed6ef4de839f612f45320575180",
"signature":"2a7548580c726423f13ea7b32b15fcfc3dc4486fb4964dd54ff8d1f4b8a7ed3fac81286e306d106aecbeff406e5f4c2a2a24d81cc0a357a4ce985e2e1e2da83f"
}