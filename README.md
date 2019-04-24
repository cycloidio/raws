# Raws: AWS Reader [![Build Status](https://travis-ci.org/cycloidio/raws.svg?branch=master)](https://travis-ci.org/cycloidio/raws) [![Coverage Status](https://coveralls.io/repos/github/cycloidio/raws/badge.svg)](https://coveralls.io/github/cycloidio/raws)

## What is Raws?

Raws is a golang project helping to get information from AWS.

It currently provides simplicity - one package vs multitude in AWS - as well as multi-region management - all calls are done for each selected region(s).
Region's parameter also supports globbing, thus allowing to fetch data from all eu with: 'eu-\*' or all eu-west with 'eu-west-\*'

Currently only a couple of the most used information is gathered, but adding extra calls should not be complicated, as they all have the same logic.

Any contributions are welcome!

**IMPORTANT** we are still experimenting the usage of this library, hence the public interface isn't stable as we have to see that the methods signatures fulfill the main goal of the library which is to simplify the AWS SDK to gather information. Because of this, the repo contains tags which define each version using [Semantic Versioning convention](https://semver.org/).

## Getting started

### Import the library
To get started, you can download/include the library to your code:
```go
import 	"github.com/cycloidio/raws"
```

### Create a reader
```go
var config *aws.Config = nil
var accessKey string = "xxxxxxxxxxxxxxxx"
var secretKey string = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
var region []string = []string{"eu-*"}

c, err := raws.NewAWSReader(accessKey, secretKey, region, config)
if err != nil {
	fmt.Printf("Error while getting NewConnector: %s\n", err.Error())
	return err
}
```

### Start making call

Errors are intentionally ignored in this example, no inputs are provided to those calls, even though one could.

```go
elbs, _ := c.GetLoadBalancersV2(nil)
fmt.Println(elbs)
instances, _ := c.GetInstances(nil)
fmt.Println(instances)
vpcs, _ := c.GetVpcs(nil)
fmt.Println(vpcs)
```

You can also take a look at the [example file](example/main.go).

### Contribute

We use a custom generation tool located on `cmd/main.go` which basically uses a list of function definitions (`cmd/functions.go`) to generate the wrappers for those,
if you want to add a call to the AWS API you have to add it to that list and if the implementation fits the template it'll be automatically generated/implemented.

If it does not fit the template you'll have to implement it manually, an example is the `s3downloader.go`.

To generate the code just run `make generate`.

### Enjoy
That's it! Nothing more, nothing less.

## Notes

### YOUR data
By default the library only returns data that belongs to you, therefore snapshots, AMI, etc are only the one that you owned and not all available objects.

This could be fixed later on depending on the needs.

### Tags everywhere?
Because the library currently simply make the call as a forwarder, it does not provide more complex calls, to return more complex data. Due to that, there are also elements to keep in mind, some calls relative to load balancer, or RDS return only the objects without tags, other calls need to be done to get those tags per resource. 

## License

Please see [LICENSE](LICENSE).

