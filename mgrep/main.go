package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/RemyJohnny/mgrep/worker"
	"github.com/RemyJohnny/mgrep/worklist"
	"github.com/alexflint/go-arg"
)

// Recursive function to explore all the file and directories within the path specified
// and populate filepaths to the workList
func exploreDir(wl *worklist.WorkList,path string){
	entries,err := os.ReadDir(path)
	if err != nil{
		fmt.Println("exploreDir error:",err)
		return
	}

	for _,entry := range entries{
		if entry.IsDir(){
			nextPath := filepath.Join(path,entry.Name())
			exploreDir(wl,nextPath)
		}else{
			fPath := filepath.Join(path,entry.Name())
			wl.Add(worklist.NewJob(fPath))
		}
	}
}

// args struct
var args struct{
	SearchKey string `arg:"positional,required"`
	SearchDir string `arg:"positional"`
}

// function used in goroutine to search a file and place the result in the result channel if there is a match.
// a worker
func work(wg *sync.WaitGroup,wl *worklist.WorkList,key string,results chan<- worker.Result){
	defer wg.Done()
	for{
		workEntry := wl.Next()
		if workEntry.Path != ""{
			workerResult := worker.SearchFile(workEntry.Path, key)
			if workerResult != nil{
				for _,r := range workerResult.Inner{
					results <- r
				}
			}
		}else{
			return
		}
	}
}

//prints result from the the result channel
func printResults(wg *sync.WaitGroup,signal <-chan int, results <-chan worker.Result){
	matchesFound := 0
	for{
		select{
		case result := <- results:
			fmt.Println(worker.ResultTempl(&result))
			matchesFound++
		case <- signal:
			if len(results) == 0{
				wg.Done()
				fmt.Printf("\n%v matches found\n",matchesFound)
				return
			}
		}
	}

}

func main(){
	arg.MustParse(&args)

	//wait Group for workers
	var workersWg sync.WaitGroup

	// creates a buffered channel for workList
	wl := worklist.New(100)

	results := make(chan worker.Result,100)

	//number of workers i.e number of goroutine spawned 
	workersNum := 10

	workersWg.Add(1)
	// goroutine to explore directory
	go func(){
		defer workersWg.Done()
		exploreDir(&wl,args.SearchDir)
		wl.Finalize(workersNum)
	}()

	//spawns goroutine according to the number or workers specified(workersNum)
	for i := 0; i < workersNum; i++{
		workersWg.Add(1)
		go work(&workersWg,&wl,args.SearchKey,results)
	}

	//channel used to signal result printer that all work is done
	blockWorkersWg := make(chan int)

	//goroutine used to wait for workers to avoid blocking the result printer
	go func(){
		workersWg.Wait()
		//signals the printer that all work is done
		close(blockWorkersWg)
	}()

	//printer's wait group
	var printWg sync.WaitGroup
	printWg.Add(1)
	// goroutine used for printing
	go printResults(&printWg,blockWorkersWg,results)

	//waits for the printer and exits the program
	printWg.Wait()
}