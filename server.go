package main

import (
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"fmt"
	"time"
	"context"
	"log"
    // "go.mongodb.org/mongo-driver/mongo/readpref"
	//"net/http"
	//"encoding/json"
	//"strconv"
)

type Booking struct {
	State string `bson:"state" json:"state"`
	District string `bson:"district" json:"district"`
	Clinic string `bson:"clinic" json:"clinic"` 
	Service string `bson:"service" json:"service"` 
	CloseDays []int `bson:"closeDays" json:"closeDays"`
	HalfDays []int `bson:"halfDays" json:"halfDays"`
	StartOpHr int `bson:"startOpHr" json:"startOpHr"`
	EndOpHr int `bson:"EndOpHr" json:"EndOpHr"`
	StartOpHrHalfDay int `bson:"startOpHrHalfDay" json:"startOpHrHalfDay"`
	EndOpHrHalfDay int `bson:"EndOpHrHalfDay" json:"EndOpHrHalfDay"`
	PublicHolMonth []int `bson:"publicHolMonth" json:"publicHolMonth"`
	StaffDaily int `bson:"staffDaily" json:"staffDaily"` 
	AvgConsultTime int `bson:"avgConsultTime" json:"avgConsultTime"` 
	MonthlyOpSchedule []DailyOpSchedule`bson:"monthlyOpSchedule" json:"monthlyOpSchedule"` 
}

type DailyOpSchedule struct {
	Date string `bson:"date" json:"date"`
	IsHalfDay int `bson:"isHalfDay" json:"isHalfDay"`
	StaffPerDay []int `bson:"staffPerDay" json:"staffPerDay"`
	QueuesCapPerDay []int `bson:"queuesCapPerDay" json:"queuesCapPerDay"` 
	QueuesPerDay []QueuePerHr `bson:"queuesPerDay" json:"queuesPerDay"` 
}

type QueuePerHr struct {
	PatientIds []int `bson:"patientIds" json:"patientIds"`
	BookingReasons []int `bson:"bookingReasons" json:"bookingReasons"`
}

func initOpSchedule() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer func() {
		if err=client.Disconnect(ctx); err!=nil {
			panic(err)
		}
		}()
		
	mytcaDB := client.Database("test")
	bookingColl := mytcaDB.Collection("booking")

	queuePerHr := QueuePerHr{
							PatientIds: []int{}, 
							BookingReasons: []int{},
				}

	dailyOpSchedule := DailyOpSchedule{
							Date: "2020-08-28",
							IsHalfDay: 1,
							StaffPerDay: []int{4, 2},
							QueuesCapPerDay: []int{36, 36},
							QueuesPerDay: []QueuePerHr{queuePerHr},
				}

    res, err := bookingColl.UpdateOne(
		ctx, 
		bson.M{"clinic" : "kk_maran"},
		bson.D{
			{"$push", bson.D{{"monthlyOpSchedule", dailyOpSchedule}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v MongoDB Documents\n", res.ModifiedCount)
}

func initOpSchedule2() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer func() {
		if err=client.Disconnect(ctx); err!=nil {
			panic(err)
		}
		}()
		
	mytcaDB := client.Database("test")
	bookingColl := mytcaDB.Collection("booking")

	staff := 5
	avgConsultTime := 10

	year := 2020
	var month time.Month = 8
	t := time.Date(year, month, 1, 0,0,0,0,time.UTC)
	lastDayOfMonth := time.Date(year, month+1, 0, 0,0,0,0,time.UTC).Day()
	monthlyOpSchedule := []DailyOpSchedule{}

	//24-hour format
	startOpHr, endOpHr := 8, 17
	startOpHrHalfDay, endOpHrHalfDay := 8, 13
	
	for i:=1; i<=lastDayOfMonth; i++ {
		if(t.Weekday()==0) {
			//It's close day, so do nothing.
		} else if(t.Weekday()==6) {
			//It's a half-day
			staffPerDay := []int{} 
			queuesCapPerDay := []int{} 
			queuesPerDay := []QueuePerHr{}
			
			for j:=startOpHrHalfDay; j < endOpHrHalfDay; j++ {
				queueCapPerHr := staff * 60 / avgConsultTime
				queuePerHr := QueuePerHr{
					PatientIds: []int{}, 
					BookingReasons: []int{},
					}

				staffPerDay = append(staffPerDay, staff)
				queuesCapPerDay = append(queuesCapPerDay, queueCapPerHr)
				queuesPerDay = append(queuesPerDay, queuePerHr)
			}

			dailyOpSchedule := DailyOpSchedule{
				Date: t.String()[:10],
				IsHalfDay: 1,
				StaffPerDay: staffPerDay,
				QueuesCapPerDay: queuesCapPerDay,
				QueuesPerDay: queuesPerDay,
				}
			monthlyOpSchedule = append(monthlyOpSchedule, dailyOpSchedule)
		} else {
			//It's a full-day
			staffPerDay := []int{} 
			queuesCapPerDay := []int{} 
			queuesPerDay := []QueuePerHr{}
			
			for j:=startOpHr; j < endOpHr; j++ {
				queueCapPerHr := staff * 60 / avgConsultTime
				queuePerHr := QueuePerHr{
					PatientIds: []int{}, 
					BookingReasons: []int{},
					}

				staffPerDay = append(staffPerDay, staff)
				queuesCapPerDay = append(queuesCapPerDay, queueCapPerHr)
				queuesPerDay = append(queuesPerDay, queuePerHr)
			}

			dailyOpSchedule := DailyOpSchedule{
				Date: t.String()[:10],
				IsHalfDay: 0,
				StaffPerDay: staffPerDay,
				QueuesCapPerDay: queuesCapPerDay,
				QueuesPerDay: queuesPerDay,
				}
			monthlyOpSchedule = append(monthlyOpSchedule, dailyOpSchedule)
		}

		t = t.AddDate(0, 0, 1)
	}

    res, err := bookingColl.UpdateOne(
		ctx, 
		bson.M{"clinic" : "kk_maran"},
		bson.D{
			{"$set", bson.D{{"monthlyOpSchedule", monthlyOpSchedule}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v MongoDB Documents\n", res.ModifiedCount)
}

func main() {
	fmt.Println("Hello mytca-booking-go_test!")
	initOpSchedule2()
}
	