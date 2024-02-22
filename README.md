# XLSWeather

实现类似于 openweathermap 的接口和返回信息，但从一个 XLSX 电子表格文件中读取。强制从 XLSX 电子表格文件中的时间视为相对日期的虚拟气象数据，XLSX 电子表格文件可以记录多天日期。

## 系统要求

- Linux：glibc 2.17 或更高版本
- macOS：10.14 Mojave 或更高版本
- Windows：10 或更高版本

## 命令行参数

`xlsweather <命令行参数>`

- `-f`: XLSX 文件路径 **(必须)**
- `-r`: 每次 GET 访问接口时，都重新读取 XLSX 电子表格文件而不是从内存中读取。
- `-d`: 基准日期 `YYYYMMDD` ，为空则为当前日期。
- `-l`: HTTP 接口所使用的 `<IP>:<端口号>` ，不提供 IP 则允许所有 IP。默认为 `127.0.0.1:80` 。
- `-u`: HTTP 接口的 URI 。默认为 `/data/2.5/weather` 。
- `-a`: 限制只有指定的几个 APPID 才能访问，使用英文逗号分隔。留空则不限制。
- `-t`: 强制按指定时间提供数据，格式示例: `"2006-01-02 15:04:05"` 。
- `-v`: 显示详细信息用于调试。
- `-rd`: 反转风向数据。
- `-tc`: 强制客户端时区为指定的 IANA 时区名称，例如 `Europe/Paris` 。
- `-ts`: 强制 XLSX 文件时区为指定的 IANA 时区名称，例如 `Asia/Tokyo` 。
- `-host`: 启动时临时添加一条项目到 hosts 文件中，结束时删除。例如 `"127.0.0.1 api.openweathermap.org"`

示例: `xlsweather -f testdata.xlsx -v`

按 `Ctrl+C` 可中止应用程序

## GET 请求参数

- `lat`: 纬度 (用于确定使用哪个时区的时间进行判断)
- `lon`: 经度 (用于确定使用哪个时区的时间进行判断)
- `APPID`: 32位 APPID
- `mode`: 返回数据类型，目前只支持 `xml`
- `units`: 度量单位，目前只支持 `metric`

示例: `http://127.0.0.1/data/2.5/weather?lat=48.0061&lon=0.1996&APPID=GGEkzWHqaaua3pdyRjzp7RiwTkvEpimV&mode=xml&units=metric`

## XLSX 电子表格文件格式要求

参见 [示例电子表格文件 testdata.xlsx](testdata.xlsx) 。

XLSX 电子表格文件中的天数部分，将会根据第一行和后续的时间重新进行计算。在程序启动时，会显示实际读入内存的数据集。

## LICENSE

Copyright (c) 2024 KagurazakaYashi XLSWeather is licensed under Mulan PSL v2. You can use this software according to the terms and conditions of the Mulan PSL v2. You may obtain a copy of Mulan PSL v2 at: http://license.coscl.org.cn/MulanPSL2 THIS SOFTWARE IS PROVIDED ON AN “AS IS” BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE. See the Mulan PSL v2 for more details.
