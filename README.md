blackrock
===


Blackrock is an integer shuffling algorithm that is based on the [implementation from massscan](https://github.com/robertdavidgraham/masscan/blob/master/src/crypto-blackrock2.c)


It is especially useful for any workload that would use iterative actions, but does not want to cause load by doing actions in an iterative way. For example, internet scanning.


Blackrock assures that the output number will appear once, and only once during a full cycle.

## Usage

You can use blackrock to scan the IPv4 space fully, without sending too many packets to a single network at a time:


```go
func exampleBlackrock() {
   br := blackrock.Init(math.MaxUint32, 5, 4)
   for i := 0; i < sizeOfArray; i++ {
       scanTarget := br.Shuffle(i)
       targetIP := intToIP(scanTarget)
       doSomethingWithNetIP(targetIP)
   }
}

func intToIP(in int) net.IP {
   buf := bytes.Buffer{}
   binary.Write(buf, binary.BigEndian, int32(in))
   return buf.Bytes()[:4]
}
```
