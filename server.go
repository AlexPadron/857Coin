
// ideally: N^(1/3) processes, N^(2/3) distinguished points

package main

import (
	"math/rand"
	"fmt"
	"math"
)

type triplet struct{
	//TODO: fill this out 
	// start location, end location, number of steps (all ints)
}

type server struct{
	peers []chan triplet
	me int 
	num_processors int 
	stored_triplets []triplet //switch to map from end to []triplet
	rand_generator *rand.Rand
	reply_channel chan triplet
}

/*
Function to create a new server
@param peers list of peers
@param me: the server's index into peers
@param num_processors: number of processors
*/
func Make(peers []chan triplet, me int, num_processors int, reply_channel chan triplet) *server{

	//store all variable into server state
	sv := &server{}
	sv.peers = peers
	sv.me = me
	sv.num_processors = num_processors

	sv.rand_generator = rand.New(rand.NewSource(int64(sv.me)))

	return sv
}

func (sv *server) start(){
	fmt.Println("starting server", sv.me)

	//TODO: implement this
	
	// two threads:
	// thread one -- computes this server's chain 
	go sv.construct_triplets()
	
	
	for trip := range(peers[me]){
		// store triplets recieved, based on end value
		// check if the list (of triples with this end value) has length >= 3
		// if so, see if we have a collision
		// 	start with the triple with longest length, hash until length becomes the next largest length (length is decreasing as you hash)
		//	if they're different, keep going; if not, drop one of them (since they're the same)
		//	continue hasing both of these until the shortest length is reached. at that point, if the three values are all different, 
		//		hash all three while keeping the previous values. if at any point all three become the same, 
		//		return their previous values as the collision
		// 	(if at some point two of them become the same but the third one remains different, we haven't found a collision,
		//		but we still want to keep all three triples in the list
		
		// if you have a collision, announce this and send it out
		// if not, prune the list of triples (so the next time you check is more efficient)
		
	}
	
}


func (sv *server) construct_triplets(){
	// randomly select start location from space of possible hashes, start
	// s = start, length = 0
	// keep hashing s until s is a distinguished point (less than N^(2/3)) ; length += 1
	// 	or until length = 20 * N / M, where M is the number of distinguished points
	// if we stopped at a distinguished point
	// 	end = where s is now
	// 	create a triple: start, end, length
	// 	send this triple to the server with number (end % (num_processes))
	
	
}



func main(){
	N := math.Exp2(20)
	num_servers := int(math.Pow(N,0.333))

	fmt.Println("Running server to solve puzzle with N:", N)

	//channel to push replies
	reply_channel := make(chan triplet)

	//list of channels to give to servers
	channels := make([]chan triplet, num_servers)

	//make list of channels
	for i := 0; i < num_servers; i++{
		channels[i] = make(chan triplet)
	}

	fmt.Println("Done initializing channels")

	//list of servers 
	servers := make([]server, num_servers)

	//make a fuckton of servers
	for i := 0; i < num_servers; i++{
		servers[i] = *Make(channels, i, num_servers, reply_channel)
	}

	//start all servers
	for i:= 0; i < num_servers; i ++{
		servers[i].start()
	}

	fmt.Println("Starting all servers, waiting for results")

	//print out anything the servers return 
	for ret := range(reply_channel){
		fmt.Println("got", ret)
	}
}
