
// ideally: N^(1/3) processes, N^(2/3) distinguished points

package main

import (
	"math/rand"
	"fmt"
	"math"
	"sort"
	//"crypto/sha256" uncomment for implementing hash function
)


//constants that all the servers can share
var N float64 = math.Exp2(20)
var M int = int(math.Pow(N,0.666))

var B block



type block struct{
	//TODO fill this in
}

type triplet struct{
	Start_location int
	End_location int
	Length int 
}

type stored_triplets []triplet

//overridden functions needed for sorting triplets by length
func (s stored_triplets) Len() int {
    return len(s)
}
func (s stored_triplets) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
//flip the order of less to get a reverse sort
func (s stored_triplets) Less(i, j int) bool {
    return s[i].Length > s[j].Length
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



/*
Method to start a server spinning to recieve triplets. For 
each triplet recieved, the server adds that triplet to its
internal data structure, indexed by the end location
*/
func (sv *server) start(){
	fmt.Println("starting server", sv.me)	
	
	//each time we get a new triplet, store it on this server
	for trip := range(sv.peers[sv.me]){
		sv.stored_triplets[trip.End_location] = append(sv.stored_triplets[trip.End_location], trip)
	}	
}

/*
Method to have a server contruct triplets. See paper for implementation
details.
*/
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

	//signal to main that this construct triplets instance has finished
	sv.reply_channel <- triplet{}
}



/*
Method to check for collison in a single server's stored triplets
*/
func (sv *server) checkForCollisions(){
	for endpoint, triplets := range(sv.stored_triplets){
		if len(triplets) >= 3{
			//sort triplets
			sort.Sort(stored_triplets(triplets))

			c_length := triplets[0].Length
			index := -1
			current_values := make([]triplet)
			for c_length > 0{
				for index + 1 < len(triplets) && triplets[index + 1].Length == c_length
			}
		}
	}
}

	

/*
Method to get a random start value for hashing from the
server's unique seeded generator
*/
func (sv *server) getRandomStart() int{
	return int(sv.rand_generator.Float64()*N)
}



/*
Given a nonce, compute the hash of this nonce with the 
context of block. Note that the block B, is stored as a global
variable to reduce memory
*/
func Hash(nonce int) int{


	//TODO: implement this hashing based on the global block B

	return 0
}


/*
Method to get a block from file. Returns a new block minus
the 3 nonces
*/
func pullBlockFromServer() block {

	//TODO: implement this
	return block{}
}


/*
Method to send a block back to file.
*/
func sendToServer(B block){

	//TODO: implement this

}


/*
Method to add 3 nonces to a block
*/
func addNoncesToBlock(B block, trip triplet){
	//TODO: implement this
}


func main(){

	//get block from server, stored as global variable
	B = pullBlockFromServer()

	//TODO: set N and M based on difficulty of puzzle
	//N = 0
	//M = 0

	//calculate the number of servers based on N
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

	// have servers make triplets
	for i:= 0; i < num_servers; i ++{
		go servers[i].construct_triplets()
	}

	fmt.Println("Waiting to finish constructing triplets")

	//wait for construct_triplets() methods to return
	count := 0
	for count < int(num_servers) {
		<- reply_channel
		count += 1
	}

	//once we are done creating triplets, have each server try to find collisions in its set
	for i:= 0; i < num_servers; i++{
		go servers[i].checkForCollisions()
	}

	fmt.Println("Waiting for collison checking")

	//if we find a triplet, add it to the block and send the block back to the server
	for ret := range(reply_channel){
		addNoncesToBlock(B, ret)
		sendToServer(B)
	}
}



