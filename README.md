# fiiicoin-adapter

fiiicoin-adapter适配了openwallet.AssetsAdapter接口，给应用提供了底层的区块链协议支持。

## 如何测试

openwtester包下的测试用例已经集成了openwallet钱包体系，创建conf文件，新建FIII.ini文件，编辑如下内容：

```ini

# node api url
serverAPI = "http://127.0.0.1:1005"
# Is network test?
isTestNet = false

```

## 资料介绍

### 官网

https://fiii.io/

### 区块浏览器

https://explorer.fiii.io/

### github

https://github.com/FiiiLab

### 适配资料

#### 地址编码的相关代码或代码链接

https://github.com/FiiiLab/FiiiCoin/blob/master/Node/Shared/FiiiChain.Consensus/AccountIdHelper.cs


#### rpc api文档的链接

https://documenter.getpostman.com/view/3484128/RzfmDmFx

#### 交易单序列化算法

https://github.com/FiiiLab/FiiiCoin/blob/master/Node/Shared/FiiiChain.Messages/TransactionMsg.cs

#### 交易单签名算法

https://github.com/FiiiLab/FiiiCoin/blob/master/Node/Shared/FiiiChain.Consensus/Script.cs