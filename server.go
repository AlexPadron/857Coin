
// ideally: N^(1/3) processes, N^(2/3) distinguished points

package main

import (
	"math/rand"
	"fmt"
	"math"
	//"crypto/sha256" uncomment for implementing hash function
)


var N float64 = math.Exp2(20)
var M int = int(math.Pow(N,0.666))

type triplet struct{
	Start_location int
	End_location int
	Length int 
}

type server struct{
	peers []chan triplet
	me int 
	num_processors int 
	stored_triplets map[int][]triplet //switch to map from end to []triplet
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

	sv.stored_triplets = make(map[int][]triplet)

	sv.rand_generator = rand.New(rand.NewSource(int64(sv.me)))

	return sv
}

func (sv *server) start(){
	fmt.Println("starting server", sv.me)	
	
	//each time we get a new triplet, store it on this server
	for trip := range(sv.peers[sv.me]){
		sv.stored_triplets[trip.End_location] = append(sv.stored_triplets[trip.End_location], trip)
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
	start := sv.getRandomStart()
	Lmax := 20*int(math.Pow(N,0.333))
	L := 0

	a := start
	for L < Lmax {
		L += 1
		a = Hash(a) 

		if a < M {
			sv.peers[int(math.Mod(float64(a), float64(N)))] <- triplet{start, a , L}
		}
	}
}


func (sv *server) checkForCollisions(){

	//TODO implement this, push collisions to reply channel
}

	

func (sv *server) getRandomStart() int{
	return int(sv.rand_generator.Float64()*N)
}

func Hash(a int) int{


	//TODO: implement this

	return 0
}


func main(){
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

	//make a fuckton of servers and start them spinning for triplets
	for i := 0; i < num_servers; i++{
		servers[i] = *Make(channels, i, num_servers, reply_channel)
		go servers[i].start()
	}

	//make triplets
	for i:= 0; i < num_servers; i ++{
		go servers[i].construct_triplets()
	}

	//once we are done creating triplets, have each server try to find collisions in its set
	for i:= 0; i < num_servers; i++{
		go servers[i].checkForCollisions()
	}

	//print out anything the servers return 
	for ret := range(reply_channel){
		fmt.Println("got triplet", ret)
	}
}
