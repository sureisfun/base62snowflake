# Base62-Snowflake

I don't expect this to be very useful to many, but I wanted to use
a shorter Snowflake ID and I couldn't find any implementations using
a base 62 alphabet, so I scrabbled this one together. 


## Basic Usage
```
import (
  "fmt"
  "github.com/sureisfun/base62snowflake/snowflake"
)

func main() {
  fmt.Println(base62snowflake.GetSnowflakeID())
}
```


