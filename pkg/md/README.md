# API接口访问签名算法文档

| **版本号** | **日期** | **修订要点** | **备注** |
| ---------- | -------- | ------------ |  -------- |
| V2.2.1   | 2022/1/7 | 新建         |            |

## 一、目的

对接口访问需要进行校验签名，提供签名算法方便开发人员快速接入。

## 二、签名算法

签名的生成要素主要由  应用`appKey`,应用`appSecret`,随机数`rand`,时间戳（单位秒）`timestamp`按照一定的排列规则组成；生成的签名不能重复使用。

签名值生成过程

- 1.获取必要的`appKey`和 `appSecret`,由贵州天机信息科技有限公司提供

- 2.随机字符串4-6位的字符串

- 3.生成当前时间戳秒

- 4.按照一定规则进行排序得到一个plain字符串

    -  排序方式：`appKey=%s&appSecret=%s&rand=%s&timestamp=%s`

- 5.最后使用`HmacSHA256`对`plain`字符串进行消息摘要,使用`appSecret`为消息摘要的密钥

- 6.在http请求头需要携带签名；`x-appKey`应用appKey，`x-signature`签名，`x-timestamp`时间戳，`x-rand`随机数。

## 三、算法参考

### （一）JAVA

-  1.随机数生成

```java
 public static final String BASE = "abcdefghijklmnopqrstuvwxyz0123456789";

     /**
     * 随机生成4-6 位字符串
     */
    public static String randStr() {
        StringBuilder sb = new StringBuilder();
        int length = (int) (Math.random() * 3 + 4);
        for (int i = 0; i < length; i++) {
            int number = (int) (Math.random() *BASE.length());
            sb.append(BASE.charAt(number));
        }
        return sb.toString();
    }
```


- 2.时间戳 秒

```java
	
	    /**
     * 获取当前时间戳
     *
     * @return 时间戳字符串
     */
    public static String timestampStr() {
        long second = Instant.now().getEpochSecond();
        return String.valueOf(second);
    }
```
- 3.签名
```java
public static String getSign(String appKey, String appSecret, String rand,String timestamp) throws Exception {
    
        
        // 有随机数的签名
        String raw = "appKey=%s&appSecret=%s&rand=%s&timestamp=%s";
        String plain = String.format(raw, appKey, appSecret,rand, timestamp);
        
        SecretKeySpec secretKeySpec = new SecretKeySpec(appSecret.getBytes(StandardCharsets.UTF_8), "HmacSHA256");
        Mac mac = Mac.getInstance("HmacSHA256");
        mac.init(secretKeySpec);
        byte[] bytes = mac.doFinal(plain.getBytes());
        return byte2HexString(bytes);
    }
  
    public static String byte2HexString(byte[] bytes) {
        StringBuilder stringBuffer = new StringBuilder();
        for (int i = 0; i < bytes.length; ++i) {
            String temp = Integer.toHexString(bytes[i] & 255);
            if (temp.length() == 1) {
                stringBuffer.append("0");
            }
            stringBuffer.append(temp);
        }
        return stringBuffer.toString();
    }
```

### （二）Golang

- 1.随机数
```go


import (
	mrand "math/rand"
	"bytes"
	"crypto/rand"
	"math/big"
)

// BASE 随机字符串
const BASE = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandStr() string {
	var randStr string
	b := bytes.NewBufferString(BASE)
	length := 4 + mrand.Intn(3)
	bigInt := big.NewInt(int64(b.Len()))
	for i := 0; i < length; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		randStr += string(BASE[randomInt.Int64()])
	}
	return randStr
}
```

- 2.签名
```go
package main

import "fmt"
import "crypto/hmac"
import "crypto/sha256"
import "encoding/hex"

func CreateAccessSign(appKey, appSecret, rand,timestamp string) string {
	    // 有随机数的签名
	plain := fmt.Sprintf("appKey=%s&appSecret=%s&rand=%s&timestamp=%s", appKey, appSecret, rand, timestamp)
        key := []byte(appSecret)
        h := hmac.New(sha256.New, key)
        h.Write([]byte(plain))
       return hex.EncodeToString(h.Sum(nil))
}


```
### （三）Python

- 签名
```python
    import hmac, time, random
    from hashlib import sha256
    
    
    def hamc_sha256(key, plain):
        sign = hmac.new(key, plain, sha256).hexdigest()
        return sign
    
    
    if __name__ == '__main__':
        appKey = 'c7btj206n88j466jth10' #应用appKey
        appSecret = 'c7btj706n88j4edermd0' #应用密钥
        rand = random.randint(100000, 900000) # 6位随机数
        timestamp = int(time.time())
        raw = "appKey={}&appSecret={}&rand={}&timestamp={}".format(appKey,appSecret, rand, timestamp)
        sign = hamc_sha256(appSecret.encode(), raw.encode())
        print(sign)

```



### （四）JS

- 签名
```js
var appKey = '[接入方 appKey]';
var appSecret = '[接入方 appSecret]';
var timestamp = Math.round(new Date().getTime()/1000);//获取秒数时间戳
var rand = Math.round(100000 + Math.random() * 900000);//6位随机数
var magic = `appKey=${appKey}&appSecret=${appSecret}&rand=${rand}&timestamp=${timestamp}`;
var signature = CryptoJS.HmacSHA256(magic, appSecret).toString();
```

## 四、demo

以Java为例;提供一个访问接口的demo。

- example
```java

    @Test
    public void signTest() throws Exception {
        //应用appKey
        String appKey = "c7btj206n88j466jth10";
        //应用密钥
        String appSecret = "c7btj706n88j4edermd0";
        //随机数
        String rand = randStr();
        //时间戳
        String timestamp = timestampStr();
        //生成签名
        String sign = getSign(appKey, appSecret, rand, timestamp);
        // 使用http访问接口
        // 请求头携带  x-appKey，x-signature，x-timestamp，x-rand 参数
        OkHttpClient client = new OkHttpClient();
        String json = "";
        String url = "http://192.168.1.22:9808/api/v1/xxx";
        RequestBody body = RequestBody.create(json, MediaType.get("application/json; charset=utf-8"));
        Request request = new Request.Builder()
                .addHeader("x-appKey", appKey)
                .addHeader("x-signature", sign)
                .addHeader("x-timestamp", timestamp)
                .addHeader("x-rand", rand)
                .url(url)
                .post(body)
                .build();
        try (Response response = client.newCall(request).execute()) {
            String bodyStr = response.body().string();
            //打印响应结果
            System.out.println(bodyStr);
        }

    }
```

```xml
<!-- https://mvnrepository.com/artifact/xin.altitude.cms.common/ucode-cms-common -->
<dependency>
  <groupId>xin.altitude.cms.common</groupId>
  <artifactId>ucode-cms-common</artifactId>
  <version>1.2.9</version>
</dependency>

```