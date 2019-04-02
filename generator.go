package core

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"gopkg.in/robfig/cron.v2"
	"log"
	"strings"
	"time"
	"v2ray.com/core/app/stats/command"
)

type StatsGenerator struct {
	DB      *Database
	ApiHost string
	ApiPort uint
	cron    *cron.Cron
	conn    *grpc.ClientConn
}

func (sg *StatsGenerator) Start() {
	address := fmt.Sprintf("%s:%d", sg.ApiHost, sg.ApiPort)
	log.Printf("Connecting to v2ray api %s", address)
	tmp, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Can't connect with v2ray api server")
	}
	sg.conn = tmp
	c := command.NewStatsServiceClient(sg.conn)
	sg.cron = cron.New()
	_, err = sg.cron.AddFunc("0 */5 * * * *", func() {
		now := time.Now()
		resp, err := c.QueryStats(context.Background(), &command.QueryStatsRequest{
			Pattern: "",
			Reset_:  true,
		})
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range resp.GetStat() {
			tmp := strings.Split(v.Name, ">>>")
			traffic := Traffic{Date: now, DataType: tmp[0], Name: tmp[1], Type: tmp[3], Value: v.Value}
			sg.DB.Create(traffic)
		}
	})
	if err != nil {
		log.Fatalln(err)
	}
	sg.cron.Start()
}

func (sg *StatsGenerator) Stop() {
	sg.cron.Stop()
	err := sg.conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
