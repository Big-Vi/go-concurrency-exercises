//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
	"sync"
)

var wg sync.WaitGroup

func producer(stream Stream, tweetread chan *Tweet) {
	
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(tweetread)
			return
		}
		tweetread <- tweet
	}
}

func consumer(tweet *Tweet) {
	if tweet.IsTalkingAboutGo() {
		fmt.Println(tweet.Username, "\ttweets about golang")
	} else {
		fmt.Println(tweet.Username, "\tdoes not tweet about golang")
	}
	wg.Done()
}

func main() {
	start := time.Now()
	
	stream := GetMockStream()

	tweetread := make(chan *Tweet, len(stream.tweets))

	// Producer
	go producer(stream, tweetread)

	for tweet := range tweetread {
		wg.Add(1)
		go consumer(tweet)
	}

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
