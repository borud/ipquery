# ipquery

This project contains two binaries.  One is a picoservice to tell you
what your IP address is.  The other is a utility to set your Network
ACL.

This is useful if you do not have a fixed IP address and you
periodically need to update your network ACL to reflect your current
IP address.

## update-acl

Update ACL figures out what your externally visible IP address is by
calling the `ipquery` picoservice.  By default the `update-acl`
program uses my personal service, but if you are going to use this I
would strongly suggest you set up your own since I can't guarantee
the stability of this service.

You can view the command line options by running `update-acl -h`:

    $ bin/update-acl -h
    Usage:
      update-acl [OPTIONS]
    
    Application Options:
      -s, --whatismyip-url= WhatIsMyIP service URL (default: https://h.borud.no/whatismyip)
      -a, --network-acl-id= Network ACL ID
      -p, --port=           Port number (default: 22)
      -t, --protocol-type=  6 is TCP, 17 is UDP (default: 6)
      -n, --rule-number=    ACL rule number
      -r, --aws-region=     AWS region (default: eu-north-1)
      -d, --dry-run         Dry run only
    
    Help Options:
      -h, --help            Show this help message

### Using scripts

One way to make updating ACLs a bit more convenient is to create
shellscripts that contain the configurations for your different
environments.  Here is an example:

    #!/bin/sh
    #
    #
	WHATISMYIP_URL=<your whatismyip service URL>
    NETWORK_ACL_ID=<your network acl id>
    AWS_REGION=<your region>
    PORT=<your port>
    RULE_ID=<your rule id>
    
    update-acl \
	    -s $WHATISMYIP_URL \
    	-a $NETWORK_ACL_ID \
    	-r $AWS_REGION \
    	-p $PORT \
    	-n $RULE_ID

### Credentials

Since we use the AWS Go SDK you have to set up your credentials so
please have a look at this page on how to get this done:
https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html

