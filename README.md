# Counter Project Go #

A golang counter library that supports different values type. 

## Examples ##

```golang
c := counter.Counter{}

// The incrementation works in new goroutine
c.Inc("keyDuration", time.Duration(20))
c.Inc("keyInt", 1)

c.WaitForFinalizationOfIncrements()

d := c.Val("keyDuration").(time.Duration)
n := c.Val("keyInt").(int)

c.Clear("keyDuration")
c.Clear("keyInt")

fmt.Println(d, n)
```

## Types supported  ##

Counter library supports the types

```golang
int
time.Duration
```