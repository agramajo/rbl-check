# RBL-Check

Check IP in RBL via DNS

## Example

```
./rbl-check -list rbl.txt -ip 186.189.238.52 
FOUND 52.238.189.186.zen.spamhaus.org 127.0.0.11
FOUND 52.238.189.186.pbl.spamhaus.org 127.0.0.11

```

## Help

```
Usage of ./rbl-check:
  -c int
    	The amount of workers to use. (default 16)
  -ip string
    	The IP to check.
  -list string
    	The RBL list.
  -t int
    	The DNS request timeout. (default 5)
  -v	Verbose output.

```

