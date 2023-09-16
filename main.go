package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type UserProfile struct {
	ID       int
	Comments []string
	Likes    int
	Friends  []int
}

func main() {
	start := time.Now()
	userProfile, err := handleGetUserProfile(10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(userProfile)
	fmt.Println("fetching user profile took", time.Since(start))
}

type Response struct {
	data any
	err  error
}

func handleGetUserProfile(id int) (*UserProfile, error) {

	//make response channel
	var (
		respch = make(chan Response, 3) //make buffered channel with 3 buffer item
		wg     = &sync.WaitGroup{}      //waitgroup to cordinate and manage go routine
	)

	//schedule a go routines for the asynchronous methods - after channel have been created
	go getComments(id, respch, wg)
	go getLikes(id, respch, wg)
	go getFriends(id, respch, wg)

	//adding 3 to the wait group (indicating no of work to be done) - and decrement afer each work is done
	wg.Add(3)

	//block untill work is done that is 0 remaining in the workgroup to unblock
	wg.Wait()

	//close the channel so the range would know to breakout of the loop
	close(respch)

	userProfile := &UserProfile{}
	//we going to range over our response channel
	for resp := range respch {
		if resp.err != nil {
			return nil, resp.err
		}
		switch msg := resp.data.(type) {
		case int:
			userProfile.Likes = msg
		case []int:
			userProfile.Friends = msg
		case []string:
			userProfile.Comments = msg
		}

	}
	return userProfile, nil
}

func getComments(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 200)
	comments := []string{
		"Hey, that was great",
		"Yeah buddy",
		"Ow, I didnt know that",
	}
	respch <- Response{
		data: comments,
		err:  nil,
	}

	//work is done
	wg.Done()
}

func getLikes(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 200)
	respch <- Response{
		data: 13,
		err:  nil,
	}

	//work is done
	wg.Done()
}

func getFriends(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 100)
	friendsIds := []int{11, 34, 60, 74, 20}
	respch <- Response{
		data: friendsIds,
		err:  nil,
	}
	//work is done
	wg.Done()

}
