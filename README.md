<p align="center">
    <b>Golib</b> เป็นโปรเจ็คที่รวบรวมเอา Package ที่มีการใช้งานบ่อยๆ มาแพ็ครวมกัน และทำ function พร้อม configuration พื้นฐานให้สามารถนำไปใช้งานได้ง่ายๆ
</p>

[![Go Doc](https://img.shields.io/badge/%F0%9F%93%9A%20godoc-pkg-00ACD7.svg?color=00ACD7&style=flat)](https://pkg.go.dev/github.com/patcharp/golib/v2?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/patcharp/golib)](https://goreportcard.com/report/github.com/patcharp/golib)

## ⚙️ Installation

ใช้ได้กับ Go version `1.21` ขึ้นไป ซึ่งสามารถดาวโหลด และติดตั้งได้ที่นี่ [Go Download](https://golang.org/dl/)

หลังจากนั้นก็เริ่มต้นโปรเจ็คใหม่ของคุณ และใช้คำสั่ง `go get` เพื่อติดตั้ง golib ใช้งาน

```bash
go get -u github.com/patcharp/golib/v2
```

## 🎯 สารบัญ package ในโปรเจ็ค
- [cache](https://github.com/patcharp/golib/tree/master/cache) สำหรับการเชื่อมต่อหา Redis และ Key-Value อื่นๆ ที่ใช้ Protocol มาตรฐานของ Redis มีทั้งการเชื่อมต่อแบบ node เดี่ยว และแบบ cluster
- [crontab](https://github.com/patcharp/golib/tree/master/crontab) สำหรับสร้าง crontab service ในตัว Go Application
- [crypto](https://github.com/patcharp/golib/tree/master/crypto) ใช้เกี่ยวกับการเข้ารหัสข้อมูล ทั้ง shared key และ public key และการเข้ารหัสแบบ default ของ crypto-js ใน node.js
- [database](https://github.com/patcharp/golib/tree/master/database) สำหรับการเชื่อมต่อ database ชนิดต่างโดยใช้ [GORM](https://gorm.io/) มี function keepalive สำหรับ automatic reconnection ได้เอง
- [geomap](https://github.com/patcharp/golib/tree/master/geomap) สำหรับคำนวน geomap แบบใช้ polygon และ พิกัด latitude, longitude
- [hashing](https://github.com/patcharp/golib/tree/master/hashing) ใช้ hash ข้อมูลแบบมาตรฐาน โดยใช้ algorithm ของ SHA256, SHA1
- [helper](https://github.com/patcharp/golib/tree/master/helper) function ช่วยเหลือต่างๆ ประกอบไปด้วยการแปลงวันเดือนปี การ response api body ของ gofiber และการแปลงวันเดือนปีให้เป็นภาษาไทย
- [imagik](https://github.com/patcharp/golib/tree/master/imagik) สำหรับการประมวลผลรูปภาพ (image processing) และการ grab รูปภาพจาก website
- [lokilog](https://github.com/patcharp/golib/tree/master/lokilog) เป็น log output ของ logrus เพื่อ push ขึ้นไปยัง [Grafana Loki](https://grafana.com/oss/loki/) 
- [mq](https://github.com/patcharp/golib/tree/master/mq) สำหรับการเชื่อมต่อหา RabbitMQ โดยใช้ Protocol AMPQ
- [one](https://github.com/patcharp/golib/tree/master/one) package ช่วยเหลือสำหรับ INET One Platform ซึ่งตอนนี้ประกอบไปด้วย OneChat, CMP และ OneID
- [requests](https://github.com/patcharp/golib/tree/master/requests) function สำหรับ Call API โดยใช้แรงบันดาลใจมาจาก lib requests ของ Python ที่สามารถ Customize parameter ต่างๆ ได้เอง
- [server](https://github.com/patcharp/golib/tree/master/server) เป็น package สำหรับการทำ web application โดยใช้ [Go Fiber](https://github.com/gofiber/fiber) และเครื่องมือช่วยเหลือต่างๆ เพื่อให้การพัฒนา web application ง่ายขึ้น
- [util](https://github.com/patcharp/golib/tree/master/util) function utilities ต่างๆ ที่จะรวมรวบเครื่องไม้เครื่องมือมาไว้ใช้งานประกอบไปด้วย
  - httputil สำหรับ web application ไม่ว่าจะ header และ function อำนวยความสะดวก
  - function common ต่างๆ เช่นกัน getenv, atoi, atof, contain

## 👀 ตัวอย่างการใช้งาน

ตัวอย่างการใช้งาน package หรือ function สำหรับอำนวยความสะดวกต่างๆ ของ package นี้

#### 📖 **Web Api โดยใช้ package server และ helper**

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/patcharp/golib/v2/server"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	s := server.New(server.Config{
		Host: "127.0.0.1",
		Port: "5000",
		HealthCheck: true,
		RequestId: true, 
	})

	s.App().Get("/api/hello", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"message": "Hello GoLib v2",
		})
	})

	s.App().Get("/api/hello_v2", func(ctx *fiber.Ctx) error {
		return helper.HttpOk(ctx, fiber.Map{
			"message": "Hello GoLib v2",
		})
	})

	if err := s.Run(); err != nil {
		logrus.Errorln("Start server error ->", err)
	}
}
```

#### 📖 **การเรียก API โดยใช้ package requests**

```go
package main

import (
	"github.com/patcharp/golib/v2/requests"
	"github.com/sirupsen/logrus"
)

func main() {
	url := "https://www.google.com"
	resp, err := requests.Get(
		url, // url
		nil, // header
		nil, // body
		0,   // timeout (second)
	)
	if err != nil {
		logrus.Errorln("Request to", url, "error ->", err)
		return
	}

	logrus.Infoln("Response code:", resp.Code)
	logrus.Infoln("Response body:", string(resp.Body))
}
```

#### 📖 **การ Connect ไป Database โดยใช้ MySQL หรือ MariaDB พร้อมทั้ง Query ข้อมูลจาก database**

```go
package main

import (
  "errors"
  "github.com/patcharp/golib/v2/database"
  "github.com/sirupsen/logrus"
  "gorm.io/gorm"
)

type User struct {
  Id   int64
  Name string
}

func main() {
  db := database.NewMySqlWithConfig(database.MySQLConfig{
    Host:         "127.0.0.1",
    Port:         "3306",
    Username:     "username",
    Password:     "password",
    DatabaseName: "db_name",
    DebugMode:    true,
  },
  )

  if err := db.Connect(); err != nil {
    logrus.Errorln("Connect to database error ->", err)
    return
  }
  defer db.Close()

  // Query all users in table
  var users []User
  if err := db.Ctx().Find(&users).Error; err != nil {
    logrus.Errorln("Query all user in users table error ->", err)
    return
  }
  logrus.Infoln("=== Query all users in table ===")
  for _, u := range users {
    logrus.Infoln("ID:", u.Id, "\tName:", u.Name)
  }

  // Query specific user where id equal to 10
  var user User
  if err := db.Ctx().Where("id=?", 10).Take(&user).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      logrus.Errorln("User who has id = 10 was not found")
      return
    }
    // Unknown error occur
    logrus.Errorln("Query specific user error ->", err)
    return
  }
  logrus.Infoln("=== Specific user who has id equal to 10 is ===")
  logrus.Infoln("ID:", user.Id, "\tName:", user.Name)
}
```

## 👍 Contribute

ถ้าคุณชอบที่สามารถนำไปใช้งานได้ในโครงการของคุณ และต้องการขอบคุณ กรุณากดให้ [GitHub Star](https://github.com/patcharp/golib/stargazers) แก่เราด้วย เพื่อเป็นกำลังใจให้พัฒนาต่อไป

## 🚀 Release log
*v2.0.10 - 8 Dec 2023*
- go.mod
  - 🐞 แก้ไข 3rd party MySQL package มีปัญหา error หลังจากที่ import ไปใช้
  - 🐞 Update 3rd part package version
  - ⚡️ ปรับ Go version ให้เป็น Go 1.21
- crypto
  - ✨ เพิ่มฟังชัน AESByteEncrypt, AESByteDecrypt สำหรับเข้ารหัสข้อมูลประเภท Byte
  - ✨ เพิ่มฟังชัน EncryptTokenWithSign, DecryptTokenWithVerifySign เพื่อแนบ Signature และ Verify Signature
  - ⚡️ ปรับฟัง EncryptToken, DecryptToken ให้ไม่ต้อง Verify Signature และแนบ Signature
  - 🐞 แก้ไขฟังชัน VerifySignedByRSAKey ให้ verify signature ได้
  - 🐞 แก้ไขฟังชัน EncodeJWTAccessToken ให้ใช้ RS256 แทนของเดิมที่ key ยาวเกิน
- database
  - ✨ เพิ่มฟังชัน SetDebug สำหรับ set debug database
  - 🐞 แก้ไข Gorm logger ให้ log severity เหมือนกันกับ database config
- imagik
  - ✨ เพิ่มฟังชันสำหรับ New struct => NewImagikFromFile, NewImagikFromByte, NewImagikFromUrl
  - ✨ เพิ่มฟังชัน ResizeWithFilter ให้สามารถกำหนด filter algorithm ได้เอง
  - ✨ เพิ่มฟังชัน ExportAsFileWithQuality ให้สามารถกำหนด JPEG Quality ได้
- server
  - ⚡️ ปรับปรุง ErrorHandler ของ DefaultConfig ให้เป็นตาม DefaultErrorHandler
  - 🐞 แก้ไข Skipper ให้สามารถใช้งานได้ตามที่ตั้งค่าไป
*v2.0.9 - 29 Aug 2023*
- crypto
  - ✨ เพิ่ม function ให้สามารถ hash password โดยกำหนดรอบได้เอง
- database
  - 🐞 แก้ไข Debug parameter ให้สามารถ enable debug ได้จริง
  - ⚡️ ปรับปรุง Gorm connection และเพิ่ม New function ของ MySQL โดยเฉพาะ
  - ⚡️ ปรับปรุง Log output ของ gorm ให้ใช้โดยใช้ Logrus
  - ✨ เพิ่ม Model type ที่ใช้ ksuid แทนการใช้ uuid
- helper
  - ✨ เพิ่ม Error 413 Entity too large, 429 Too many request error
  - ✨ เพิ่ม Type Today สำหรับต้องการใช้งานเรียกวันที่ปัจจุบัน และ set เวลาเป็น 00:00:00 
- mq
  - ⚡️ ปรับปรุง Lib rabbitmq ใหม่ทั้งหมด และไปใช้ amqp มาตรฐานของ rabbitmq
  - 🐞 แก้ไข Connection ที่ไม่สามารถ reconnect ไปหา server เองได้
- server
  - ✨ เพิ่ม request id middleware ให้เรียกใช้ได้
  - ✨ เพิ่ม log middleware ให้ไปใช้ Logrus
  - ✨ เพิ่ม route /api/-/health สำหรับ healthcheck และไม่ทำการ print access log ให้เรียกใช้ได้
  - 🐞 แก้ไข default server error ให้แสดงผลให้ถูกต้อง
  - 🐞 แก้ไขให้ fiber print stacktrace หาก service เกิด crash และ return server error ออกไปหา client
- httputil
  - ✨ เพิ่ม function สำหรับดึงค่า ksuid จาก fiber context param
  - 🛠 เปลี่ยน function สำหรับดึงค่า uid จาก fiber context param ให้สื่อความหมาย
- util
  - 🚧 ลบ function time ที่ซ้ำกับ package today

## ‍💻 Code Contributors

[![](https://avatars.githubusercontent.com/u/40089397?s=32&v=4)](https://github.com/patcharp)

## ⚠️ License

Copyright (c) 2021 Patcharapong and Contributors โดยที่ GoLib นี้เป็น Open Source ที่ทุกคนสามารถนำไปใช้งานได้ฟรี ภายใต้อนุสัญญา [MIT License](https://github.com/patcharp/golib/blob/master/LICENSE) ตาม package ต้นทางที่เราได้นำมาใช้งาน
