package main

import (
	"os"
	"log"
	"github.com/Nik-U/pbc"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	params := pbc.GenerateA(160, 512)
  	pairing := params.NewPairing()
 	//g := pairing.NewG2().Rand()

	aggregate := pairing.NewG2()
	res := pairing.NewG2()
	//res1 := pairing.NewG2()
	res2 := pairing.NewG2()
	res3 := pairing.NewG2()
	res4 := pairing.NewG2()
	res5 := pairing.NewG2()
	res6 := pairing.NewG2()
	

	file, err := os.Open("result.txt")	
	check(err)
	defer file.Close()

	/*file1, err := os.Open("result1.txt")	
	check(err)
	defer file1.Close()*/

	file2, err := os.Open("result2.txt")	
	check(err)
	defer file2.Close()

	file3, err := os.Open("result3.txt")	
	check(err)
	defer file3.Close()

	file4, err := os.Open("result4.txt")	
	check(err)
	defer file4.Close()

	file5, err := os.Open("result5.txt")	
	check(err)
	defer file5.Close()

	file6, err := os.Open("result6.txt")	
	check(err)
	defer file6.Close()

	_, err = file.Read(res.Bytes())
	check(err)
	aggregate = pairing.NewG2().Add(aggregate, res)

	/*_, err = file1.Read(res1.Bytes())
	check(err)
	aggregate = pairing.NewG2().Add(aggregate, res1)*/

	_, err = file2.Read(res2.Bytes())
	check(err)
	aggregate = pairing.NewG2().Add(aggregate, res2)

	_, err = file3.Read(res3.Bytes())
	check(err)
	aggregate = pairing.NewG2().Add(aggregate, res3)

	_, err = file4.Read(res4.Bytes())
	check(err)
	aggregate = pairing.NewG2().Add(aggregate, res4)

	_, err = file5.Read(res5.Bytes())
	check(err)
	aggregate = pairing.NewG2().Add(aggregate, res5)

	_, err = file6.Read(res6.Bytes())
	check(err)
	aggregate = pairing.NewG2().Add(aggregate, res6)
	
	aggregate_sum := aggregate.Bytes()

	log.Printf("Aggregate %v", aggregate_sum)

}
