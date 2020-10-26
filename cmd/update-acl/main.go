package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/jessevdk/go-flags"
)

// Options contains the command line options
type Options struct {
	WhatIsMyIPServiceURL string `short:"s" long:"whatismyip-url" default:"https://h.borud.no/whatismyip" description:"WhatIsMyIP service URL" required:"yes"`
	NetworkACLID         string `short:"a" long:"network-acl-id" default:"" description:"Network ACL ID" required:"yes"`
	Port                 int64  `short:"p" long:"port" default:"22" description:"Port number" required:"yes"`
	Protocol             string `short:"t" long:"protocol-type" default:"6" description:"6 is TCP, 17 is UDP" required:"yes"`
	RuleNumber           int64  `short:"n" long:"rule-number" description:"ACL rule number" required:"yes"`
	AWSRegion            string `short:"r" long:"aws-region" default:"eu-north-1" description:"AWS region" required:"yes"`
	DryRun               bool   `short:"d" long:"dry-run" description:"Dry run only"`
}

var opt Options
var parser = flags.NewParser(&opt, flags.Default)

func main() {
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Printf(">>> %v", err)
		os.Exit(1)
	}

	// Make sure rule numbers are in the 200 range
	if opt.RuleNumber > 299 || opt.RuleNumber < 200 {
		log.Fatalf("Rule number (-n) should be between 200 and 299")
	}

	// Get my IP address
	myIP := whatIsMyIP()
	if myIP == nil {
		log.Fatalf("IP address parse error")
	}
	log.Printf("My ip is '%s'", myIP)

	// Update ACL
	updateACL(myIP)
}

// whatIsMyIP polls the whatismyip service to figure out what my IP is
func whatIsMyIP() net.IP {
	resp, err := http.Get(opt.WhatIsMyIPServiceURL)
	if err != nil {
		log.Fatalf("Unable to look up my IP: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read my IP response: %v", err)
	}

	return net.ParseIP(string(body))
}

// updateACL updates the Network ACL
func updateACL(ip net.IP) {

	input := &ec2.ReplaceNetworkAclEntryInput{
		CidrBlock:    aws.String(fmt.Sprintf("%s/32", ip)),
		Egress:       aws.Bool(false),
		NetworkAclId: aws.String(opt.NetworkACLID),
		PortRange: &ec2.PortRange{
			From: aws.Int64(opt.Port),
			To:   aws.Int64(opt.Port),
		},
		Protocol:   aws.String(opt.Protocol),
		RuleAction: aws.String("allow"),
		RuleNumber: aws.Int64(opt.RuleNumber),
	}

	if opt.DryRun {
		log.Printf("Would update ACL: %v", input)
		return
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String(opt.AWSRegion)})
	if err != nil {
		log.Fatalf("Unable to create new AWS session: %v", err)
	}

	svc := ec2.New(sess)

	_, err = svc.ReplaceNetworkAclEntry(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Fatalf("AWS Error: %v", aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return
	}

	log.Printf("Updated ACL: %v", input)
}
